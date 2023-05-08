package hello

import (
	"context"
	"errors"
	"io"
	"sync"

	"github.com/pdeziel/grpc-example/mock"
	"github.com/pdeziel/grpc-example/pb"
	"google.golang.org/grpc"
)

// Client wraps a generated gRPC client with the connection
type Client struct {
	api pb.HelloClient
	cc  *grpc.ClientConn
}

// Create a new client from client options
func NewClient(endpoint string, opts ...grpc.DialOption) (c *Client, err error) {
	c = &Client{}

	// "Dial" the server, by default this is a non-blocking call which establishes a
	// connection in the background but doesn't do anything with it yet.
	if c.cc, err = grpc.Dial(endpoint, opts...); err != nil {
		return nil, err
	}
	c.api = pb.NewHelloClient(c.cc)
	return
}

// Close the connection
func (c *Client) Close() error {
	return c.cc.Close()
}

// Connect the mock to the client
func (c *Client) ConnectMock(mock *mock.HelloService, opts ...grpc.DialOption) (err error) {
	if c.api, err = mock.Client(context.Background(), opts...); err != nil {
		return err
	}

	return nil
}

func (c *Client) SayHello(ctx context.Context, langCode string) (_ string, err error) {
	req := &pb.HelloRequest{
		IsoLanguageCode: langCode,
	}

	var rep *pb.HelloReply
	if rep, err = c.api.SayHello(ctx, req); err != nil {
		return "", err
	}

	return rep.Greeting, nil
}

func (c *Client) SayClientStream(ctx context.Context, langs <-chan string) (greetings []string, err error) {
	var stream pb.Hello_SayClientStreamClient
	if stream, err = c.api.SayClientStream(ctx); err != nil {
		return nil, err
	}

	for langCode := range langs {
		req := &pb.HelloRequest{
			IsoLanguageCode: langCode,
		}

		if err = stream.Send(req); err != nil {
			return nil, err
		}
	}

	var rep *pb.HelloManyReply
	if rep, err = stream.CloseAndRecv(); err != nil {
		return nil, err
	}

	greetings = make([]string, 0, len(rep.Greetings))
	for _, greeting := range rep.Greetings {
		greetings = append(greetings, greeting.Greeting)
	}

	return greetings, nil
}

func (c *Client) SayServerStream(ctx context.Context, langCodes []string) (greetings []string, err error) {
	req := &pb.HelloManyRequest{
		IsoLanguageCodes: langCodes,
	}

	var stream pb.Hello_SayServerStreamClient
	if stream, err = c.api.SayServerStream(ctx, req); err != nil {
		return nil, err
	}

	for {
		var rep *pb.HelloReply
		if rep, err = stream.Recv(); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		greetings = append(greetings, rep.Greeting)
	}

	return greetings, nil
}

func (c *Client) SayBidirectional(ctx context.Context, langCodes []string) (greetings []string, err error) {
	var stream pb.Hello_SayBidirectionalClient
	if stream, err = c.api.SayBidirectional(ctx); err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, langCode := range langCodes {
			req := &pb.HelloRequest{
				IsoLanguageCode: langCode,
			}

			if err = stream.Send(req); err != nil {
				return
			}
		}

		stream.CloseSend()
	}()

	for {
		var rep *pb.HelloReply
		if rep, err = stream.Recv(); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		greetings = append(greetings, rep.Greeting)
	}

	// Wait for the sender goroutine to receive all the responses
	wg.Wait()

	return greetings, nil
}
