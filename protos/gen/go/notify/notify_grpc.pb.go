// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: notify/notify.proto

package notifyv1

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

// SendNotifyClient is the client API for SendNotify service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SendNotifyClient interface {
	SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error)
}

type sendNotifyClient struct {
	cc grpc.ClientConnInterface
}

func NewSendNotifyClient(cc grpc.ClientConnInterface) SendNotifyClient {
	return &sendNotifyClient{cc}
}

func (c *sendNotifyClient) SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error) {
	out := new(SendMessageResponse)
	err := c.cc.Invoke(ctx, "/notify.SendNotify/SendMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SendNotifyServer is the server API for SendNotify service.
// All implementations must embed UnimplementedSendNotifyServer
// for forward compatibility
type SendNotifyServer interface {
	SendMessage(context.Context, *SendMessageRequest) (*SendMessageResponse, error)
	mustEmbedUnimplementedSendNotifyServer()
}

// UnimplementedSendNotifyServer must be embedded to have forward compatible implementations.
type UnimplementedSendNotifyServer struct {
}

func (UnimplementedSendNotifyServer) SendMessage(context.Context, *SendMessageRequest) (*SendMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedSendNotifyServer) mustEmbedUnimplementedSendNotifyServer() {}

// UnsafeSendNotifyServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SendNotifyServer will
// result in compilation errors.
type UnsafeSendNotifyServer interface {
	mustEmbedUnimplementedSendNotifyServer()
}

func RegisterSendNotifyServer(s grpc.ServiceRegistrar, srv SendNotifyServer) {
	s.RegisterService(&SendNotify_ServiceDesc, srv)
}

func _SendNotify_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SendNotifyServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/notify.SendNotify/SendMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SendNotifyServer).SendMessage(ctx, req.(*SendMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SendNotify_ServiceDesc is the grpc.ServiceDesc for SendNotify service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SendNotify_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "notify.SendNotify",
	HandlerType: (*SendNotifyServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMessage",
			Handler:    _SendNotify_SendMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "notify/notify.proto",
}

// UsersClient is the client API for Users service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UsersClient interface {
	UserPreferences(ctx context.Context, in *UserPreferencesRequest, opts ...grpc.CallOption) (*UserPreferencesResponse, error)
}

type usersClient struct {
	cc grpc.ClientConnInterface
}

func NewUsersClient(cc grpc.ClientConnInterface) UsersClient {
	return &usersClient{cc}
}

func (c *usersClient) UserPreferences(ctx context.Context, in *UserPreferencesRequest, opts ...grpc.CallOption) (*UserPreferencesResponse, error) {
	out := new(UserPreferencesResponse)
	err := c.cc.Invoke(ctx, "/notify.Users/UserPreferences", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UsersServer is the server API for Users service.
// All implementations must embed UnimplementedUsersServer
// for forward compatibility
type UsersServer interface {
	UserPreferences(context.Context, *UserPreferencesRequest) (*UserPreferencesResponse, error)
	mustEmbedUnimplementedUsersServer()
}

// UnimplementedUsersServer must be embedded to have forward compatible implementations.
type UnimplementedUsersServer struct {
}

func (UnimplementedUsersServer) UserPreferences(context.Context, *UserPreferencesRequest) (*UserPreferencesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserPreferences not implemented")
}
func (UnimplementedUsersServer) mustEmbedUnimplementedUsersServer() {}

// UnsafeUsersServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UsersServer will
// result in compilation errors.
type UnsafeUsersServer interface {
	mustEmbedUnimplementedUsersServer()
}

func RegisterUsersServer(s grpc.ServiceRegistrar, srv UsersServer) {
	s.RegisterService(&Users_ServiceDesc, srv)
}

func _Users_UserPreferences_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserPreferencesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).UserPreferences(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/notify.Users/UserPreferences",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).UserPreferences(ctx, req.(*UserPreferencesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Users_ServiceDesc is the grpc.ServiceDesc for Users service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Users_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "notify.Users",
	HandlerType: (*UsersServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UserPreferences",
			Handler:    _Users_UserPreferences_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "notify/notify.proto",
}
