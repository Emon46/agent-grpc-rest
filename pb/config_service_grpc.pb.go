// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: config_service.proto

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

// ConfigServiceClient is the client API for ConfigService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConfigServiceClient interface {
	GetWorkerConfig(ctx context.Context, in *GetWorkerConfigRequest, opts ...grpc.CallOption) (*WorkerConfigResponse, error)
	UpdateWorkerConfig(ctx context.Context, in *UpdateWorkerConfigRequest, opts ...grpc.CallOption) (*WorkerConfigResponse, error)
}

type configServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewConfigServiceClient(cc grpc.ClientConnInterface) ConfigServiceClient {
	return &configServiceClient{cc}
}

func (c *configServiceClient) GetWorkerConfig(ctx context.Context, in *GetWorkerConfigRequest, opts ...grpc.CallOption) (*WorkerConfigResponse, error) {
	out := new(WorkerConfigResponse)
	err := c.cc.Invoke(ctx, "/config.ConfigService/GetWorkerConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configServiceClient) UpdateWorkerConfig(ctx context.Context, in *UpdateWorkerConfigRequest, opts ...grpc.CallOption) (*WorkerConfigResponse, error) {
	out := new(WorkerConfigResponse)
	err := c.cc.Invoke(ctx, "/config.ConfigService/UpdateWorkerConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConfigServiceServer is the server API for ConfigService service.
// All implementations must embed UnimplementedConfigServiceServer
// for forward compatibility
type ConfigServiceServer interface {
	GetWorkerConfig(context.Context, *GetWorkerConfigRequest) (*WorkerConfigResponse, error)
	UpdateWorkerConfig(context.Context, *UpdateWorkerConfigRequest) (*WorkerConfigResponse, error)
	mustEmbedUnimplementedConfigServiceServer()
}

// UnimplementedConfigServiceServer must be embedded to have forward compatible implementations.
type UnimplementedConfigServiceServer struct {
}

func (UnimplementedConfigServiceServer) GetWorkerConfig(context.Context, *GetWorkerConfigRequest) (*WorkerConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWorkerConfig not implemented")
}
func (UnimplementedConfigServiceServer) UpdateWorkerConfig(context.Context, *UpdateWorkerConfigRequest) (*WorkerConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateWorkerConfig not implemented")
}
func (UnimplementedConfigServiceServer) mustEmbedUnimplementedConfigServiceServer() {}

// UnsafeConfigServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConfigServiceServer will
// result in compilation errors.
type UnsafeConfigServiceServer interface {
	mustEmbedUnimplementedConfigServiceServer()
}

func RegisterConfigServiceServer(s grpc.ServiceRegistrar, srv ConfigServiceServer) {
	s.RegisterService(&ConfigService_ServiceDesc, srv)
}

func _ConfigService_GetWorkerConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetWorkerConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServiceServer).GetWorkerConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/config.ConfigService/GetWorkerConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServiceServer).GetWorkerConfig(ctx, req.(*GetWorkerConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConfigService_UpdateWorkerConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateWorkerConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServiceServer).UpdateWorkerConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/config.ConfigService/UpdateWorkerConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServiceServer).UpdateWorkerConfig(ctx, req.(*UpdateWorkerConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ConfigService_ServiceDesc is the grpc.ServiceDesc for ConfigService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ConfigService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "config.ConfigService",
	HandlerType: (*ConfigServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetWorkerConfig",
			Handler:    _ConfigService_GetWorkerConfig_Handler,
		},
		{
			MethodName: "UpdateWorkerConfig",
			Handler:    _ConfigService_UpdateWorkerConfig_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "config_service.proto",
}