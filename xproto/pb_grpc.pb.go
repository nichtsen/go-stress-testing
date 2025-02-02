// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package xproto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// ApiServerClient is the client API for ApiServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ApiServerClient interface {
	// HelloWorld 接口
	HelloWorld(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type apiServerClient struct {
	cc grpc.ClientConnInterface
}

func NewApiServerClient(cc grpc.ClientConnInterface) ApiServerClient {
	return &apiServerClient{cc}
}

func (c *apiServerClient) HelloWorld(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/xproto.ApiServer/HelloWorld", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ApiServerServer is the server API for ApiServer service.
// All implementations must embed UnimplementedApiServerServer
// for forward compatibility
type ApiServerServer interface {
	// HelloWorld 接口
	HelloWorld(context.Context, *Request) (*Response, error)
	mustEmbedUnimplementedApiServerServer()
}

// UnimplementedApiServerServer must be embedded to have forward compatible implementations.
type UnimplementedApiServerServer struct {
}

func (UnimplementedApiServerServer) HelloWorld(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HelloWorld not implemented")
}
func (UnimplementedApiServerServer) mustEmbedUnimplementedApiServerServer() {}

// UnsafeApiServerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ApiServerServer will
// result in compilation errors.
type UnsafeApiServerServer interface {
	mustEmbedUnimplementedApiServerServer()
}

func RegisterApiServerServer(s grpc.ServiceRegistrar, srv ApiServerServer) {
	s.RegisterService(&ApiServer_ServiceDesc, srv)
}

func _ApiServer_HelloWorld_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServerServer).HelloWorld(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/xproto.ApiServer/HelloWorld",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServerServer).HelloWorld(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// ApiServer_ServiceDesc is the grpc.ServiceDesc for ApiServer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ApiServer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "xproto.ApiServer",
	HandlerType: (*ApiServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HelloWorld",
			Handler:    _ApiServer_HelloWorld_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "xproto/pb.proto",
}
