// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: protobuf/locations.proto

package protobuf

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

// LocationsClient is the client API for Locations service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LocationsClient interface {
	// Sends a greeting
	Insert(ctx context.Context, in *NewLocation, opts ...grpc.CallOption) (*Empty, error)
}

type locationsClient struct {
	cc grpc.ClientConnInterface
}

func NewLocationsClient(cc grpc.ClientConnInterface) LocationsClient {
	return &locationsClient{cc}
}

func (c *locationsClient) Insert(ctx context.Context, in *NewLocation, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/gridgo.Locations/Insert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LocationsServer is the server API for Locations service.
// All implementations must embed UnimplementedLocationsServer
// for forward compatibility
type LocationsServer interface {
	// Sends a greeting
	Insert(context.Context, *NewLocation) (*Empty, error)
	mustEmbedUnimplementedLocationsServer()
}

// UnimplementedLocationsServer must be embedded to have forward compatible implementations.
type UnimplementedLocationsServer struct {
}

func (UnimplementedLocationsServer) Insert(context.Context, *NewLocation) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Insert not implemented")
}
func (UnimplementedLocationsServer) mustEmbedUnimplementedLocationsServer() {}

// UnsafeLocationsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LocationsServer will
// result in compilation errors.
type UnsafeLocationsServer interface {
	mustEmbedUnimplementedLocationsServer()
}

func RegisterLocationsServer(s grpc.ServiceRegistrar, srv LocationsServer) {
	s.RegisterService(&Locations_ServiceDesc, srv)
}

func _Locations_Insert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewLocation)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocationsServer).Insert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gridgo.Locations/Insert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocationsServer).Insert(ctx, req.(*NewLocation))
	}
	return interceptor(ctx, in, info, handler)
}

// Locations_ServiceDesc is the grpc.ServiceDesc for Locations service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Locations_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gridgo.Locations",
	HandlerType: (*LocationsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Insert",
			Handler:    _Locations_Insert_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protobuf/locations.proto",
}
