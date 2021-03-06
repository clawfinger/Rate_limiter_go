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

// CalendarClient is the client API for Calendar service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CalendarClient interface {
	Validate(ctx context.Context, in *LoginAttempt, opts ...grpc.CallOption) (*AttemptResult, error)
	DropStats(ctx context.Context, in *Stats, opts ...grpc.CallOption) (*OperationResult, error)
	AppBlacklist(ctx context.Context, in *IP, opts ...grpc.CallOption) (*OperationResult, error)
	RemoveBlacklist(ctx context.Context, in *IP, opts ...grpc.CallOption) (*OperationResult, error)
	AppWhitelist(ctx context.Context, in *IP, opts ...grpc.CallOption) (*OperationResult, error)
	RemoveWhitelist(ctx context.Context, in *IP, opts ...grpc.CallOption) (*OperationResult, error)
}

type calendarClient struct {
	cc grpc.ClientConnInterface
}

func NewCalendarClient(cc grpc.ClientConnInterface) CalendarClient {
	return &calendarClient{cc}
}

func (c *calendarClient) Validate(ctx context.Context, in *LoginAttempt, opts ...grpc.CallOption) (*AttemptResult, error) {
	out := new(AttemptResult)
	err := c.cc.Invoke(ctx, "/ratelimiter.Calendar/Validate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) DropStats(ctx context.Context, in *Stats, opts ...grpc.CallOption) (*OperationResult, error) {
	out := new(OperationResult)
	err := c.cc.Invoke(ctx, "/ratelimiter.Calendar/DropStats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) AppBlacklist(ctx context.Context, in *IP, opts ...grpc.CallOption) (*OperationResult, error) {
	out := new(OperationResult)
	err := c.cc.Invoke(ctx, "/ratelimiter.Calendar/AppBlacklist", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) RemoveBlacklist(ctx context.Context, in *IP, opts ...grpc.CallOption) (*OperationResult, error) {
	out := new(OperationResult)
	err := c.cc.Invoke(ctx, "/ratelimiter.Calendar/RemoveBlacklist", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) AppWhitelist(ctx context.Context, in *IP, opts ...grpc.CallOption) (*OperationResult, error) {
	out := new(OperationResult)
	err := c.cc.Invoke(ctx, "/ratelimiter.Calendar/AppWhitelist", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) RemoveWhitelist(ctx context.Context, in *IP, opts ...grpc.CallOption) (*OperationResult, error) {
	out := new(OperationResult)
	err := c.cc.Invoke(ctx, "/ratelimiter.Calendar/RemoveWhitelist", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CalendarServer is the server API for Calendar service.
// All implementations must embed UnimplementedCalendarServer
// for forward compatibility
type CalendarServer interface {
	Validate(context.Context, *LoginAttempt) (*AttemptResult, error)
	DropStats(context.Context, *Stats) (*OperationResult, error)
	AppBlacklist(context.Context, *IP) (*OperationResult, error)
	RemoveBlacklist(context.Context, *IP) (*OperationResult, error)
	AppWhitelist(context.Context, *IP) (*OperationResult, error)
	RemoveWhitelist(context.Context, *IP) (*OperationResult, error)
	mustEmbedUnimplementedCalendarServer()
}

// UnimplementedCalendarServer must be embedded to have forward compatible implementations.
type UnimplementedCalendarServer struct {
}

func (UnimplementedCalendarServer) Validate(context.Context, *LoginAttempt) (*AttemptResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validate not implemented")
}
func (UnimplementedCalendarServer) DropStats(context.Context, *Stats) (*OperationResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DropStats not implemented")
}
func (UnimplementedCalendarServer) AppBlacklist(context.Context, *IP) (*OperationResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppBlacklist not implemented")
}
func (UnimplementedCalendarServer) RemoveBlacklist(context.Context, *IP) (*OperationResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveBlacklist not implemented")
}
func (UnimplementedCalendarServer) AppWhitelist(context.Context, *IP) (*OperationResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppWhitelist not implemented")
}
func (UnimplementedCalendarServer) RemoveWhitelist(context.Context, *IP) (*OperationResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveWhitelist not implemented")
}
func (UnimplementedCalendarServer) mustEmbedUnimplementedCalendarServer() {}

// UnsafeCalendarServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CalendarServer will
// result in compilation errors.
type UnsafeCalendarServer interface {
	mustEmbedUnimplementedCalendarServer()
}

func RegisterCalendarServer(s grpc.ServiceRegistrar, srv CalendarServer) {
	s.RegisterService(&Calendar_ServiceDesc, srv)
}

func _Calendar_Validate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginAttempt)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).Validate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ratelimiter.Calendar/Validate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).Validate(ctx, req.(*LoginAttempt))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_DropStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Stats)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).DropStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ratelimiter.Calendar/DropStats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).DropStats(ctx, req.(*Stats))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_AppBlacklist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IP)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).AppBlacklist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ratelimiter.Calendar/AppBlacklist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).AppBlacklist(ctx, req.(*IP))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_RemoveBlacklist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IP)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).RemoveBlacklist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ratelimiter.Calendar/RemoveBlacklist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).RemoveBlacklist(ctx, req.(*IP))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_AppWhitelist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IP)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).AppWhitelist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ratelimiter.Calendar/AppWhitelist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).AppWhitelist(ctx, req.(*IP))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_RemoveWhitelist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IP)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).RemoveWhitelist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ratelimiter.Calendar/RemoveWhitelist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).RemoveWhitelist(ctx, req.(*IP))
	}
	return interceptor(ctx, in, info, handler)
}

// Calendar_ServiceDesc is the grpc.ServiceDesc for Calendar service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Calendar_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ratelimiter.Calendar",
	HandlerType: (*CalendarServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Validate",
			Handler:    _Calendar_Validate_Handler,
		},
		{
			MethodName: "DropStats",
			Handler:    _Calendar_DropStats_Handler,
		},
		{
			MethodName: "AppBlacklist",
			Handler:    _Calendar_AppBlacklist_Handler,
		},
		{
			MethodName: "RemoveBlacklist",
			Handler:    _Calendar_RemoveBlacklist_Handler,
		},
		{
			MethodName: "AppWhitelist",
			Handler:    _Calendar_AppWhitelist_Handler,
		},
		{
			MethodName: "RemoveWhitelist",
			Handler:    _Calendar_RemoveWhitelist_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc_api.proto",
}
