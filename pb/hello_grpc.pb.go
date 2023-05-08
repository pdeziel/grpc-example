// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// HelloClient is the client API for Hello service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HelloClient interface {
	// Unary RPC - one request and one response
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
	// Server side streaming RPC - one request and a stream of responses
	SayServerStream(ctx context.Context, in *HelloManyRequest, opts ...grpc.CallOption) (Hello_SayServerStreamClient, error)
	// Client side streaming RPC - a stream of requests and one response
	SayClientStream(ctx context.Context, opts ...grpc.CallOption) (Hello_SayClientStreamClient, error)
	// Bidirectional streaming RPC - a stream of requests and a stream of responses
	SayBidirectional(ctx context.Context, opts ...grpc.CallOption) (Hello_SayBidirectionalClient, error)
}

type helloClient struct {
	cc grpc.ClientConnInterface
}

func NewHelloClient(cc grpc.ClientConnInterface) HelloClient {
	return &helloClient{cc}
}

func (c *helloClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/hello.Hello/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *helloClient) SayServerStream(ctx context.Context, in *HelloManyRequest, opts ...grpc.CallOption) (Hello_SayServerStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Hello_ServiceDesc.Streams[0], "/hello.Hello/SayServerStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &helloSayServerStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Hello_SayServerStreamClient interface {
	Recv() (*HelloReply, error)
	grpc.ClientStream
}

type helloSayServerStreamClient struct {
	grpc.ClientStream
}

func (x *helloSayServerStreamClient) Recv() (*HelloReply, error) {
	m := new(HelloReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *helloClient) SayClientStream(ctx context.Context, opts ...grpc.CallOption) (Hello_SayClientStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Hello_ServiceDesc.Streams[1], "/hello.Hello/SayClientStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &helloSayClientStreamClient{stream}
	return x, nil
}

type Hello_SayClientStreamClient interface {
	Send(*HelloRequest) error
	CloseAndRecv() (*HelloManyReply, error)
	grpc.ClientStream
}

type helloSayClientStreamClient struct {
	grpc.ClientStream
}

func (x *helloSayClientStreamClient) Send(m *HelloRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *helloSayClientStreamClient) CloseAndRecv() (*HelloManyReply, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(HelloManyReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *helloClient) SayBidirectional(ctx context.Context, opts ...grpc.CallOption) (Hello_SayBidirectionalClient, error) {
	stream, err := c.cc.NewStream(ctx, &Hello_ServiceDesc.Streams[2], "/hello.Hello/SayBidirectional", opts...)
	if err != nil {
		return nil, err
	}
	x := &helloSayBidirectionalClient{stream}
	return x, nil
}

type Hello_SayBidirectionalClient interface {
	Send(*HelloRequest) error
	Recv() (*HelloReply, error)
	grpc.ClientStream
}

type helloSayBidirectionalClient struct {
	grpc.ClientStream
}

func (x *helloSayBidirectionalClient) Send(m *HelloRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *helloSayBidirectionalClient) Recv() (*HelloReply, error) {
	m := new(HelloReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// HelloServer is the server API for Hello service.
// All implementations must embed UnimplementedHelloServer
// for forward compatibility
type HelloServer interface {
	// Unary RPC - one request and one response
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
	// Server side streaming RPC - one request and a stream of responses
	SayServerStream(*HelloManyRequest, Hello_SayServerStreamServer) error
	// Client side streaming RPC - a stream of requests and one response
	SayClientStream(Hello_SayClientStreamServer) error
	// Bidirectional streaming RPC - a stream of requests and a stream of responses
	SayBidirectional(Hello_SayBidirectionalServer) error
	mustEmbedUnimplementedHelloServer()
}

// UnimplementedHelloServer must be embedded to have forward compatible implementations.
type UnimplementedHelloServer struct {
}

func (UnimplementedHelloServer) SayHello(context.Context, *HelloRequest) (*HelloReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (UnimplementedHelloServer) SayServerStream(*HelloManyRequest, Hello_SayServerStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method SayServerStream not implemented")
}
func (UnimplementedHelloServer) SayClientStream(Hello_SayClientStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method SayClientStream not implemented")
}
func (UnimplementedHelloServer) SayBidirectional(Hello_SayBidirectionalServer) error {
	return status.Errorf(codes.Unimplemented, "method SayBidirectional not implemented")
}
func (UnimplementedHelloServer) mustEmbedUnimplementedHelloServer() {}

// UnsafeHelloServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HelloServer will
// result in compilation errors.
type UnsafeHelloServer interface {
	mustEmbedUnimplementedHelloServer()
}

func RegisterHelloServer(s grpc.ServiceRegistrar, srv HelloServer) {
	s.RegisterService(&Hello_ServiceDesc, srv)
}

func _Hello_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HelloServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hello.Hello/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HelloServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Hello_SayServerStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(HelloManyRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(HelloServer).SayServerStream(m, &helloSayServerStreamServer{stream})
}

type Hello_SayServerStreamServer interface {
	Send(*HelloReply) error
	grpc.ServerStream
}

type helloSayServerStreamServer struct {
	grpc.ServerStream
}

func (x *helloSayServerStreamServer) Send(m *HelloReply) error {
	return x.ServerStream.SendMsg(m)
}

func _Hello_SayClientStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(HelloServer).SayClientStream(&helloSayClientStreamServer{stream})
}

type Hello_SayClientStreamServer interface {
	SendAndClose(*HelloManyReply) error
	Recv() (*HelloRequest, error)
	grpc.ServerStream
}

type helloSayClientStreamServer struct {
	grpc.ServerStream
}

func (x *helloSayClientStreamServer) SendAndClose(m *HelloManyReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *helloSayClientStreamServer) Recv() (*HelloRequest, error) {
	m := new(HelloRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Hello_SayBidirectional_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(HelloServer).SayBidirectional(&helloSayBidirectionalServer{stream})
}

type Hello_SayBidirectionalServer interface {
	Send(*HelloReply) error
	Recv() (*HelloRequest, error)
	grpc.ServerStream
}

type helloSayBidirectionalServer struct {
	grpc.ServerStream
}

func (x *helloSayBidirectionalServer) Send(m *HelloReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *helloSayBidirectionalServer) Recv() (*HelloRequest, error) {
	m := new(HelloRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Hello_ServiceDesc is the grpc.ServiceDesc for Hello service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Hello_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "hello.Hello",
	HandlerType: (*HelloServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Hello_SayHello_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SayServerStream",
			Handler:       _Hello_SayServerStream_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SayClientStream",
			Handler:       _Hello_SayClientStream_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "SayBidirectional",
			Handler:       _Hello_SayBidirectional_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "hello.proto",
}
