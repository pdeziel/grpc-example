package mock

import (
	"context"
	"sync"

	"github.com/pdeziel/grpc-example/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// This file contains an example of a gRPC service mock that can be used for testing
// client-side gRPC code. Test code can connect using a bufconn to avoid real network
// requests but still exercise the gRPC stack.

// These constants represent the endpoints of the gRPC service and map directly to the
// FullMethod in the generated gRPC code. They are intended to be used directly in
// tests to mock specific responses for endpoints.
const (
	SayHelloRPC = "/hello.Hello/SayHello"
)

var ErrUnavailable = status.Error(codes.Unavailable, "mock method has not been configured")

// New creates a new mock service that can be connected to using a bufconn. If no bufnet
// is provided, a default one is created.
func New(bufnet *Listener, opts ...grpc.ServerOption) *HelloService {
	if bufnet == nil {
		bufnet = NewBufConn()
	}

	hs := &HelloService{
		bufnet: bufnet,
		srv:    grpc.NewServer(opts...),
		Calls:  make(map[string]int),
	}

	pb.RegisterHelloServer(hs.srv, hs)
	go hs.srv.Serve(hs.bufnet.Sock())

	return hs
}

type HelloService struct {
	sync.RWMutex
	pb.UnimplementedHelloServer
	bufnet             *Listener
	srv                *grpc.Server
	client             pb.HelloClient
	Calls              map[string]int
	OnSayHello         func(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error)
	OnSayClientStream  func(ctx context.Context, stream pb.Hello_SayClientStreamServer) error
	OnSayServerStream  func(req *pb.HelloManyRequest, stream pb.Hello_SayServerStreamServer) error
	OnSayBidirectional func(stream pb.Hello_SayBidirectionalServer) error
}

// Create and connect a client to the mock server
func (s *HelloService) Client(ctx context.Context, opts ...grpc.DialOption) (client pb.HelloClient, err error) {
	if s.client == nil {
		if len(opts) == 0 {
			opts = make([]grpc.DialOption, 0, 1)
			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		}

		var cc *grpc.ClientConn
		if cc, err = s.bufnet.Connect(ctx, opts...); err != nil {
			return nil, err
		}
		s.client = pb.NewHelloClient(cc)
	}
	return s.client, nil
}

// Reset the client with the new dial options
func (s *HelloService) ResetClient(ctx context.Context, opts ...grpc.DialOption) (pb.HelloClient, error) {
	s.client = nil
	return s.Client(ctx, opts...)
}

// Shutdown the sever and cleanup (cannot be used after shutdown)
func (s *HelloService) Shutdown() {
	s.srv.GracefulStop()
	s.bufnet.Close()
}

// Reset the calls map and all associated handlers in preparation for a new test.
func (s *HelloService) Reset() {
	for key := range s.Calls {
		delete(s.Calls, key)
	}

	s.OnSayHello = nil
	s.OnSayClientStream = nil
	s.OnSayServerStream = nil
	s.OnSayBidirectional = nil
}

func (s *HelloService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	s.incrCall()
	if s.OnSayHello != nil {
		return s.OnSayHello(ctx, req)
	}

	return nil, ErrUnavailable
}

func (s *HelloService) SayClientStream(stream pb.Hello_SayClientStreamServer) error {
	s.incrCall()
	if s.OnSayClientStream != nil {
		return s.OnSayClientStream(stream.Context(), stream)
	}

	return ErrUnavailable
}

func (s *HelloService) SayServerStream(req *pb.HelloManyRequest, stream pb.Hello_SayServerStreamServer) error {
	s.incrCall()
	if s.OnSayServerStream != nil {
		return s.OnSayServerStream(req, stream)
	}

	return ErrUnavailable
}

func (s *HelloService) SayBidirectional(stream pb.Hello_SayBidirectionalServer) error {
	s.incrCall()
	if s.OnSayBidirectional != nil {
		return s.OnSayBidirectional(stream)
	}

	return ErrUnavailable
}

func (s *HelloService) incrCall() {
	s.Lock()
	defer s.Unlock()
	s.Calls[SayHelloRPC]++
}
