package example

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ExampleClient is the clientStream API for Example service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ExampleClient interface {
	Echo(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	ClientStream(ctx context.Context, opts ...grpc.CallOption) (Example_ClientStreamClient, error)
	ServerStream(ctx context.Context, in *Request, opts ...grpc.CallOption) (Example_ServerStreamClient, error)
	BiStream(ctx context.Context, opts ...grpc.CallOption) (Example_BiStreamClient, error)
}

type exampleClient struct {
	cc grpc.ClientConnInterface
}

func NewExampleClient(cc grpc.ClientConnInterface) ExampleClient {
	return &exampleClient{cc}
}

func (c *exampleClient) Echo(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/ppzxc.grpc.example.Example/Echo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exampleClient) ClientStream(ctx context.Context, opts ...grpc.CallOption) (Example_ClientStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Example_serviceDesc.Streams[0], "/ppzxc.grpc.example.Example/ClientStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &exampleClientStreamClient{stream}
	return x, nil
}

type Example_ClientStreamClient interface {
	Send(*Request) error
	CloseAndRecv() (*Response, error)
	grpc.ClientStream
}

type exampleClientStreamClient struct {
	grpc.ClientStream
}

func (x *exampleClientStreamClient) Send(m *Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *exampleClientStreamClient) CloseAndRecv() (*Response, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *exampleClient) ServerStream(ctx context.Context, in *Request, opts ...grpc.CallOption) (Example_ServerStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Example_serviceDesc.Streams[1], "/ppzxc.grpc.example.Example/ServerStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &exampleServerStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Example_ServerStreamClient interface {
	Recv() (*Response, error)
	grpc.ClientStream
}

type exampleServerStreamClient struct {
	grpc.ClientStream
}

func (x *exampleServerStreamClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *exampleClient) BiStream(ctx context.Context, opts ...grpc.CallOption) (Example_BiStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Example_serviceDesc.Streams[2], "/ppzxc.grpc.example.Example/BiStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &exampleBiStreamClient{stream}
	return x, nil
}

type Example_BiStreamClient interface {
	Send(*Request) error
	Recv() (*Response, error)
	grpc.ClientStream
}

type exampleBiStreamClient struct {
	grpc.ClientStream
}

func (x *exampleBiStreamClient) Send(m *Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *exampleBiStreamClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ExampleServer is the server API for Example service.
type ExampleServer interface {
	Echo(context.Context, *Request) (*Response, error)
	ClientStream(Example_ClientStreamServer) error
	ServerStream(*Request, Example_ServerStreamServer) error
	BiStream(Example_BiStreamServer) error
}

// UnimplementedExampleServer can be embedded to have forward compatible implementations.
type UnimplementedExampleServer struct {
}

func (*UnimplementedExampleServer) Echo(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Echo not implemented")
}
func (*UnimplementedExampleServer) ClientStream(Example_ClientStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method ClientStream not implemented")
}
func (*UnimplementedExampleServer) ServerStream(*Request, Example_ServerStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method ServerStream not implemented")
}
func (*UnimplementedExampleServer) BiStream(Example_BiStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method BiStream not implemented")
}

func RegisterExampleServer(s *grpc.Server, srv ExampleServer) {
	s.RegisterService(&_Example_serviceDesc, srv)
}

func _Example_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExampleServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ppzxc.grpc.example.Example/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExampleServer).Echo(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Example_ClientStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ExampleServer).ClientStream(&exampleClientStreamServer{stream})
}

type Example_ClientStreamServer interface {
	SendAndClose(*Response) error
	Recv() (*Request, error)
	grpc.ServerStream
}

type exampleClientStreamServer struct {
	grpc.ServerStream
}

func (x *exampleClientStreamServer) SendAndClose(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *exampleClientStreamServer) Recv() (*Request, error) {
	m := new(Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Example_ServerStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ExampleServer).ServerStream(m, &exampleServerStreamServer{stream})
}

type Example_ServerStreamServer interface {
	Send(*Response) error
	grpc.ServerStream
}

type exampleServerStreamServer struct {
	grpc.ServerStream
}

func (x *exampleServerStreamServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func _Example_BiStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ExampleServer).BiStream(&exampleBiStreamServer{stream})
}

type Example_BiStreamServer interface {
	Send(*Response) error
	Recv() (*Request, error)
	grpc.ServerStream
}

type exampleBiStreamServer struct {
	grpc.ServerStream
}

func (x *exampleBiStreamServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *exampleBiStreamServer) Recv() (*Request, error) {
	m := new(Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Example_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ppzxc.grpc.example.Example",
	HandlerType: (*ExampleServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _Example_Echo_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ClientStream",
			Handler:       _Example_ClientStream_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "ServerStream",
			Handler:       _Example_ServerStream_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "BiStream",
			Handler:       _Example_BiStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "example.proto",
}
