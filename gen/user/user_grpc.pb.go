// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: user.proto

package user

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	UserService_LoginUser_FullMethodName    = "/user.UserService/LoginUser"
	UserService_RegisterUser_FullMethodName = "/user.UserService/RegisterUser"
)

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceClient interface {
	LoginUser(ctx context.Context, in *LoginMessageRequest, opts ...grpc.CallOption) (*LoginMessageResponse, error)
	RegisterUser(ctx context.Context, in *RegisterMessageRequest, opts ...grpc.CallOption) (*RegisterMessageResponse, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) LoginUser(ctx context.Context, in *LoginMessageRequest, opts ...grpc.CallOption) (*LoginMessageResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginMessageResponse)
	err := c.cc.Invoke(ctx, UserService_LoginUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) RegisterUser(ctx context.Context, in *RegisterMessageRequest, opts ...grpc.CallOption) (*RegisterMessageResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterMessageResponse)
	err := c.cc.Invoke(ctx, UserService_RegisterUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility.
type UserServiceServer interface {
	LoginUser(context.Context, *LoginMessageRequest) (*LoginMessageResponse, error)
	RegisterUser(context.Context, *RegisterMessageRequest) (*RegisterMessageResponse, error)
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedUserServiceServer struct{}

func (UnimplementedUserServiceServer) LoginUser(context.Context, *LoginMessageRequest) (*LoginMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginUser not implemented")
}
func (UnimplementedUserServiceServer) RegisterUser(context.Context, *RegisterMessageRequest) (*RegisterMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterUser not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}
func (UnimplementedUserServiceServer) testEmbeddedByValue()                     {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	// If the following call pancis, it indicates UnimplementedUserServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_LoginUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).LoginUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_LoginUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).LoginUser(ctx, req.(*LoginMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_RegisterUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).RegisterUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_RegisterUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).RegisterUser(ctx, req.(*RegisterMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LoginUser",
			Handler:    _UserService_LoginUser_Handler,
		},
		{
			MethodName: "RegisterUser",
			Handler:    _UserService_RegisterUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user.proto",
}
