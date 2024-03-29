// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package path_apiv1

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

// StationsClient is the client API for Stations service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StationsClient interface {
	// Lists the metadata for all available stations.
	ListStations(ctx context.Context, in *ListStationsRequest, opts ...grpc.CallOption) (*ListStationsResponse, error)
	// Gets the metadata for a specific station.
	GetStation(ctx context.Context, in *GetStationRequest, opts ...grpc.CallOption) (*StationData, error)
	// Gets the posted train schedule for a station.
	GetStationSchedule(ctx context.Context, in *GetStationScheduleRequest, opts ...grpc.CallOption) (*GetStationScheduleResponse, error)
	// Gets the expected upcoming trains for the station using realtime data.
	GetUpcomingTrains(ctx context.Context, in *GetUpcomingTrainsRequest, opts ...grpc.CallOption) (*GetUpcomingTrainsResponse, error)
}

type stationsClient struct {
	cc grpc.ClientConnInterface
}

func NewStationsClient(cc grpc.ClientConnInterface) StationsClient {
	return &stationsClient{cc}
}

func (c *stationsClient) ListStations(ctx context.Context, in *ListStationsRequest, opts ...grpc.CallOption) (*ListStationsResponse, error) {
	out := new(ListStationsResponse)
	err := c.cc.Invoke(ctx, "/path_api.v1.Stations/ListStations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stationsClient) GetStation(ctx context.Context, in *GetStationRequest, opts ...grpc.CallOption) (*StationData, error) {
	out := new(StationData)
	err := c.cc.Invoke(ctx, "/path_api.v1.Stations/GetStation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stationsClient) GetStationSchedule(ctx context.Context, in *GetStationScheduleRequest, opts ...grpc.CallOption) (*GetStationScheduleResponse, error) {
	out := new(GetStationScheduleResponse)
	err := c.cc.Invoke(ctx, "/path_api.v1.Stations/GetStationSchedule", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stationsClient) GetUpcomingTrains(ctx context.Context, in *GetUpcomingTrainsRequest, opts ...grpc.CallOption) (*GetUpcomingTrainsResponse, error) {
	out := new(GetUpcomingTrainsResponse)
	err := c.cc.Invoke(ctx, "/path_api.v1.Stations/GetUpcomingTrains", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StationsServer is the server API for Stations service.
// All implementations should embed UnimplementedStationsServer
// for forward compatibility
type StationsServer interface {
	// Lists the metadata for all available stations.
	ListStations(context.Context, *ListStationsRequest) (*ListStationsResponse, error)
	// Gets the metadata for a specific station.
	GetStation(context.Context, *GetStationRequest) (*StationData, error)
	// Gets the posted train schedule for a station.
	GetStationSchedule(context.Context, *GetStationScheduleRequest) (*GetStationScheduleResponse, error)
	// Gets the expected upcoming trains for the station using realtime data.
	GetUpcomingTrains(context.Context, *GetUpcomingTrainsRequest) (*GetUpcomingTrainsResponse, error)
}

// UnimplementedStationsServer should be embedded to have forward compatible implementations.
type UnimplementedStationsServer struct {
}

func (UnimplementedStationsServer) ListStations(context.Context, *ListStationsRequest) (*ListStationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListStations not implemented")
}
func (UnimplementedStationsServer) GetStation(context.Context, *GetStationRequest) (*StationData, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStation not implemented")
}
func (UnimplementedStationsServer) GetStationSchedule(context.Context, *GetStationScheduleRequest) (*GetStationScheduleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStationSchedule not implemented")
}
func (UnimplementedStationsServer) GetUpcomingTrains(context.Context, *GetUpcomingTrainsRequest) (*GetUpcomingTrainsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUpcomingTrains not implemented")
}

// UnsafeStationsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StationsServer will
// result in compilation errors.
type UnsafeStationsServer interface {
	mustEmbedUnimplementedStationsServer()
}

func RegisterStationsServer(s grpc.ServiceRegistrar, srv StationsServer) {
	s.RegisterService(&Stations_ServiceDesc, srv)
}

func _Stations_ListStations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListStationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StationsServer).ListStations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/path_api.v1.Stations/ListStations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StationsServer).ListStations(ctx, req.(*ListStationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stations_GetStation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StationsServer).GetStation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/path_api.v1.Stations/GetStation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StationsServer).GetStation(ctx, req.(*GetStationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stations_GetStationSchedule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStationScheduleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StationsServer).GetStationSchedule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/path_api.v1.Stations/GetStationSchedule",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StationsServer).GetStationSchedule(ctx, req.(*GetStationScheduleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stations_GetUpcomingTrains_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUpcomingTrainsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StationsServer).GetUpcomingTrains(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/path_api.v1.Stations/GetUpcomingTrains",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StationsServer).GetUpcomingTrains(ctx, req.(*GetUpcomingTrainsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Stations_ServiceDesc is the grpc.ServiceDesc for Stations service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Stations_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "path_api.v1.Stations",
	HandlerType: (*StationsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListStations",
			Handler:    _Stations_ListStations_Handler,
		},
		{
			MethodName: "GetStation",
			Handler:    _Stations_GetStation_Handler,
		},
		{
			MethodName: "GetStationSchedule",
			Handler:    _Stations_GetStationSchedule_Handler,
		},
		{
			MethodName: "GetUpcomingTrains",
			Handler:    _Stations_GetUpcomingTrains_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "stations.proto",
}
