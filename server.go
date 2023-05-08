package hello

import (
	"context"
	"errors"
	"io"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/pdeziel/grpc-example/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Struct that implements the gRPC service
type Server struct {
	pb.UnimplementedHelloServer
	srv      *grpc.Server
	messages *Messages
	echan    chan error
}

// Create a new server
func NewServer(opts ...grpc.ServerOption) (s *Server, err error) {
	s = &Server{
		echan: make(chan error),
	}

	// Load the messages from the JSON file
	s.messages = &Messages{}
	if err = s.messages.Load("messages.json"); err != nil {
		return nil, err
	}

	s.srv = grpc.NewServer(opts...)
	pb.RegisterHelloServer(s.srv, s)
	return
}

// Start the server
func (s *Server) Serve(bindaddr string) (err error) {
	// Catch OS signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		s.echan <- s.Shutdown()
	}()

	// Listen for TCP requests
	var sock net.Listener
	if sock, err = net.Listen("tcp", bindaddr); err != nil {
		return err
	}

	go s.Run(sock)

	// Wait on the error channel until the server is stopped
	if err = <-s.echan; err != nil {
		return err
	}
	return nil
}

// Run the gRPC server on the socket. This is extracted into its own method to make it
// easier to inject sockets for testing.
func (s *Server) Run(sock net.Listener) {
	defer sock.Close()
	if err := s.srv.Serve(sock); err != nil {
		s.echan <- err
	}
}

// Stop the server
func (s *Server) Shutdown() error {
	s.srv.GracefulStop()
	return nil
}

// Unary RPC
func (s *Server) SayHello(ctx context.Context, req *pb.HelloRequest) (rep *pb.HelloReply, err error) {
	// Lookup the message from the request
	var msg string
	if msg, err = s.messages.Get(req.IsoLanguageCode); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	// Return the message
	rep = &pb.HelloReply{
		Greeting:        msg,
		IsoLanguageCode: req.IsoLanguageCode,
		CreatedAt:       time.Now().Format(time.RFC3339),
	}

	return rep, nil
}

// Client streaming RPC
func (s *Server) SayClientStream(stream pb.Hello_SayClientStreamServer) (err error) {
	reply := &pb.HelloManyReply{
		Greetings: make([]*pb.HelloReply, 0),
	}
	for {
		var req *pb.HelloRequest
		if req, err = stream.Recv(); err != nil {
			if errors.Is(err, io.EOF) {
				return stream.SendAndClose(reply)
			}

			return err
		}

		var msg string
		if msg, err = s.messages.Get(req.IsoLanguageCode); err != nil {
			return status.Error(codes.NotFound, err.Error())
		}

		reply.Greetings = append(reply.Greetings, &pb.HelloReply{
			Greeting:        msg,
			IsoLanguageCode: req.IsoLanguageCode,
			CreatedAt:       time.Now().Format(time.RFC3339),
		})
	}
}

// Server streaming RPC
func (s *Server) SayServerStream(req *pb.HelloManyRequest, stream pb.Hello_SayServerStreamServer) (err error) {
	for _, iso := range req.IsoLanguageCodes {
		var msg string
		if msg, err = s.messages.Get(iso); err != nil {
			return status.Error(codes.NotFound, err.Error())
		}

		if err = stream.Send(&pb.HelloReply{
			Greeting:        msg,
			IsoLanguageCode: iso,
			CreatedAt:       time.Now().Format(time.RFC3339),
		}); err != nil {
			return err
		}
	}

	return nil
}

// Bidirectional streaming RPC
func (s *Server) SayBidirectional(stream pb.Hello_SayBidirectionalServer) (err error) {
	for {
		var req *pb.HelloRequest
		if req, err = stream.Recv(); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}

			return err
		}

		var msg string
		if msg, err = s.messages.Get(req.IsoLanguageCode); err != nil {
			return status.Error(codes.NotFound, err.Error())
		}

		if err = stream.Send(&pb.HelloReply{
			Greeting:        msg,
			IsoLanguageCode: req.IsoLanguageCode,
			CreatedAt:       time.Now().Format(time.RFC3339),
		}); err != nil {
			return err
		}
	}
}
