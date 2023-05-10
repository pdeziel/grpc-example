package hello_test

import (
	"context"
	"errors"
	"io"
	"sync"
	"testing"

	hello "github.com/pdeziel/grpc-example"
	"github.com/pdeziel/grpc-example/mock"
	"github.com/pdeziel/grpc-example/pb"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type clientTestSuite struct {
	suite.Suite
	server *mock.HelloService
	client *hello.Client
}

func TestClient(t *testing.T) {
	suite.Run(t, new(clientTestSuite))
}

func (s *clientTestSuite) SetupSuite() {
	// Start the mock service using the default bufconn
	s.server = mock.New(nil)

	// Connect a client to the mock
	require := s.Require()
	require.NoError(s.initClient(), "could not connect a client to the mock server")
}

func (s *clientTestSuite) TearDownSuite() {
	s.server.Shutdown()
}

func (s *clientTestSuite) AfterTest() {
	s.server.Reset()
}

// Create a new client from the mock service
func (s *clientTestSuite) initClient() (err error) {
	s.client = &hello.Client{}
	if err = s.client.ConnectMock(s.server, grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		return err
	}

	return nil
}

func (s *clientTestSuite) TestSayHello() {
	require := s.Require()

	// Configure the server mock
	s.server.OnSayHello = func(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
		return &pb.HelloReply{
			Greeting: "Bonjour",
		}, nil
	}

	// Make a unary request
	greeting, err := s.client.SayHello(context.Background(), "fr")
	require.NoError(err, "could not call the service")
	require.Equal("Bonjour", greeting)
}

func (s *clientTestSuite) TestSayClientStream() {
	require := s.Require()
	messages := []string{"Hello", "Bonjour", "Hola"}

	// Configure the server mock
	s.server.OnSayClientStream = func(ctx context.Context, stream pb.Hello_SayClientStreamServer) (err error) {
		reply := &pb.HelloManyReply{
			Greetings: make([]*pb.HelloReply, 0, len(messages)),
		}
		var sent int
		for {
			if _, err = stream.Recv(); err != nil {
				if errors.Is(err, io.EOF) {
					return stream.SendAndClose(reply)
				}
				return err
			}

			if sent >= len(messages) {
				return status.Error(codes.InvalidArgument, "no more messages")
			}

			reply.Greetings = append(reply.Greetings, &pb.HelloReply{
				Greeting: messages[sent],
			})
			sent++
		}
	}

	// Make a client stream request
	langs := make(chan string, len(messages))
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, message := range messages {
			langs <- message
		}
		close(langs)
	}()

	greetings, err := s.client.SayClientStream(context.Background(), langs)
	require.NoError(err, "could not call the service")
	require.Equal(messages, greetings)

	wg.Wait()
}

func (s *clientTestSuite) TestSayServerStream() {
	require := s.Require()
	messages := []string{"Hello", "Bonjour", "Hola"}

	// Configure the server mock
	s.server.OnSayServerStream = func(req *pb.HelloManyRequest, stream pb.Hello_SayServerStreamServer) (err error) {
		for _, message := range messages {
			if err = stream.Send(&pb.HelloReply{
				Greeting: message,
			}); err != nil {
				return err
			}
		}
		return nil
	}

	// Make a server stream request
	langs := []string{"en", "fr", "es"}
	greetings, err := s.client.SayServerStream(context.Background(), langs)
	require.NoError(err, "could not call the service")
	require.Equal(messages, greetings)
}

func (s *clientTestSuite) TestSayBidirectional() {
	require := s.Require()
	messages := []string{"Hello", "Bonjour", "Hola"}

	// Configure the server mock
	s.server.OnSayBidirectional = func(stream pb.Hello_SayBidirectionalServer) (err error) {
		var sent int
		for {
			if _, err = stream.Recv(); err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				return err
			}

			if sent >= len(messages) {
				return status.Error(codes.InvalidArgument, "no more messages")
			}

			if err = stream.Send(&pb.HelloReply{
				Greeting: messages[sent],
			}); err != nil {
				return err
			}
			sent++
		}
	}

	// Setup the send stream
	langs := make(chan string, len(messages))
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, message := range messages {
			langs <- message
		}
		close(langs)
	}()

	// Setup the receive stream
	responses := make(chan string, len(messages))
	greetings := make([]string, 0, len(messages))
	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range responses {
			greetings = append(greetings, msg)
		}
	}()

	err := s.client.SayBidirectional(context.Background(), langs, responses)
	require.NoError(err, "could not call the service")
	require.Equal(messages, greetings)
}
