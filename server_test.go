package hello_test

import (
	"context"
	"errors"
	"io"
	"testing"

	hello "github.com/pdeziel/grpc-example"
	"github.com/pdeziel/grpc-example/mock"
	"github.com/pdeziel/grpc-example/pb"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type serverTestSuite struct {
	suite.Suite
	server *hello.Server
	grpc   *mock.Listener
}

func TestServer(t *testing.T) {
	suite.Run(t, new(serverTestSuite))
}

func (s *serverTestSuite) SetupSuite() {
	require := s.Require()

	var err error
	s.server, err = hello.NewServer()
	require.NoError(err, "could not start the server for tests")

	// Create a bufconn listener to avoid network requests
	s.grpc = mock.NewBufConn()

	// Run the server directly
	go s.server.Run(s.grpc.Sock())
}

func (s *serverTestSuite) TearDownSuite() {
	s.Require().NoError(s.server.Shutdown(), "could not shutdown the server")
}

func (s *serverTestSuite) initClient(ctx context.Context) pb.HelloClient {
	cc, err := s.grpc.Connect(ctx, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err, "could not create client connection to the server")
	return pb.NewHelloClient(cc)
}

func (s *serverTestSuite) TestSayHello() {
	require := s.Require()

	// Create a client connection to the server
	client := s.initClient(context.Background())

	// Make a unary request
	req := &pb.HelloRequest{
		IsoLanguageCode: "fr",
	}
	rep, err := client.SayHello(context.Background(), req)
	require.NoError(err, "could not call the service")
	require.Equal("Bonjour", rep.Greeting, "unexpected greeting")
}

func (s *serverTestSuite) TestSayClientStream() {
	require := s.Require()

	// Create a client connection to the server
	client := s.initClient(context.Background())

	// Make a client streaming request
	stream, err := client.SayClientStream(context.Background())
	require.NoError(err, "could not call the service")

	// Send some hello requests from the client
	messages := []string{"en", "fr", "es"}
	for _, msg := range messages {
		require.NoError(stream.Send(&pb.HelloRequest{
			IsoLanguageCode: msg,
		}), "could not send a message")
	}

	expected := []string{"Hello", "Bonjour", "Hola"}

	// Close the stream from the client side
	rep, err := stream.CloseAndRecv()
	require.NoError(err, "could not close the stream")
	require.Len(rep.Greetings, len(expected), "unexpected number of greetings")
	for i, msg := range rep.Greetings {
		require.Equal(expected[i], msg.Greeting, "unexpected greeting")
	}
}

func (s *serverTestSuite) TestSayServerStream() {
	require := s.Require()

	// Create a client connection to the server
	client := s.initClient(context.Background())

	// Make a server streaming request
	stream, err := client.SayServerStream(context.Background(), &pb.HelloManyRequest{
		IsoLanguageCodes: []string{"en", "fr", "es"},
	})
	require.NoError(err, "could not call the service")

	expected := []string{"Hello", "Bonjour", "Hola"}

	// Receive greetings from the server
	actual := make([]string, 0, len(expected))
	for {
		var rep *pb.HelloReply
		if rep, err = stream.Recv(); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
		}

		require.NoError(err, "could not receive a message")
		actual = append(actual, rep.Greeting)
	}

	require.Equal(expected, actual, "unexpected greetings")
}

func (s *serverTestSuite) TestSayBidirectional() {
	require := s.Require()

	// Create a client connection to the server
	client := s.initClient(context.Background())

	// Make a bidirectional streaming request
	stream, err := client.SayBidirectional(context.Background())
	require.NoError(err, "could not call the service")

	// Send some hello requests from the client
	messages := []string{"en", "fr", "es"}
	expected := []string{"Hello", "Bonjour", "Hola"}
	actual := make([]string, 0, len(expected))
	for _, msg := range messages {
		require.NoError(stream.Send(&pb.HelloRequest{
			IsoLanguageCode: msg,
		}), "could not send a message")

		var rep *pb.HelloReply
		if rep, err = stream.Recv(); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
		}
		require.NoError(err, "could not receive a message from the server")
		actual = append(actual, rep.Greeting)
	}

	if err = stream.CloseSend(); err != nil {
		if errors.Is(err, io.EOF) {
			return
		}
	}

	require.Equal(expected, actual, "unexpected greetings")
}
