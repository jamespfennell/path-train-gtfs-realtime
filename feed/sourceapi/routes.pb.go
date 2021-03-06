// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.22.0-devel
// 	protoc        v3.6.0
// source: routes.proto

package path_api_v1

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type ListRoutesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Optional. The maximum number of elements to return for a single request.
	// If unspecified, the server will pick a reasonable default.
	PageSize int32 `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// Optional. The page token returned by the server in a previous call. Used
	// to get the next page.
	PageToken string `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
}

func (x *ListRoutesRequest) Reset() {
	*x = ListRoutesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_routes_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRoutesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRoutesRequest) ProtoMessage() {}

func (x *ListRoutesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_routes_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRoutesRequest.ProtoReflect.Descriptor instead.
func (*ListRoutesRequest) Descriptor() ([]byte, []int) {
	return file_routes_proto_rawDescGZIP(), []int{0}
}

func (x *ListRoutesRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListRoutesRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

type ListRoutesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The page of routes.
	Routes []*RouteData `protobuf:"bytes,1,rep,name=routes,proto3" json:"routes,omitempty"`
	// The page token used to request the next page. Empty/unspecified if there
	// are no more results.
	NextPageToken string `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
}

func (x *ListRoutesResponse) Reset() {
	*x = ListRoutesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_routes_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRoutesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRoutesResponse) ProtoMessage() {}

func (x *ListRoutesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_routes_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRoutesResponse.ProtoReflect.Descriptor instead.
func (*ListRoutesResponse) Descriptor() ([]byte, []int) {
	return file_routes_proto_rawDescGZIP(), []int{1}
}

func (x *ListRoutesResponse) GetRoutes() []*RouteData {
	if x != nil {
		return x.Routes
	}
	return nil
}

func (x *ListRoutesResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

type GetRouteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The route to get information about.
	Route Route `protobuf:"varint,1,opt,name=route,proto3,enum=path_api.v1.Route" json:"route,omitempty"`
}

func (x *GetRouteRequest) Reset() {
	*x = GetRouteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_routes_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRouteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRouteRequest) ProtoMessage() {}

func (x *GetRouteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_routes_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRouteRequest.ProtoReflect.Descriptor instead.
func (*GetRouteRequest) Descriptor() ([]byte, []int) {
	return file_routes_proto_rawDescGZIP(), []int{2}
}

func (x *GetRouteRequest) GetRoute() Route {
	if x != nil {
		return x.Route
	}
	return Route_ROUTE_UNSPECIFIED
}

type GetRouteScheduleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Route Route `protobuf:"varint,1,opt,name=route,proto3,enum=path_api.v1.Route" json:"route,omitempty"`
	// Optional. The maximum number of elements to return for a single request.
	// If unspecified, the server will pick a reasonable default.
	PageSize int32 `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// Optional. The page token returned by the server in a previous call. Used
	// to get the next page.
	PageToken string `protobuf:"bytes,3,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
}

func (x *GetRouteScheduleRequest) Reset() {
	*x = GetRouteScheduleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_routes_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRouteScheduleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRouteScheduleRequest) ProtoMessage() {}

func (x *GetRouteScheduleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_routes_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRouteScheduleRequest.ProtoReflect.Descriptor instead.
func (*GetRouteScheduleRequest) Descriptor() ([]byte, []int) {
	return file_routes_proto_rawDescGZIP(), []int{3}
}

func (x *GetRouteScheduleRequest) GetRoute() Route {
	if x != nil {
		return x.Route
	}
	return Route_ROUTE_UNSPECIFIED
}

func (x *GetRouteScheduleRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *GetRouteScheduleRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

type GetRouteScheduleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NextPageToken string `protobuf:"bytes,1,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
}

func (x *GetRouteScheduleResponse) Reset() {
	*x = GetRouteScheduleResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_routes_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRouteScheduleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRouteScheduleResponse) ProtoMessage() {}

func (x *GetRouteScheduleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_routes_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRouteScheduleResponse.ProtoReflect.Descriptor instead.
func (*GetRouteScheduleResponse) Descriptor() ([]byte, []int) {
	return file_routes_proto_rawDescGZIP(), []int{4}
}

func (x *GetRouteScheduleResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

// Data representing a single route in the PATH system.
type RouteData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The route this entry represents.
	Route Route `protobuf:"varint,1,opt,name=route,proto3,enum=path_api.v1.Route" json:"route,omitempty"`
	// The ID in the GTFS database of this route.
	Id string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	// The name (long name) of the route.
	Name string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	// The color (headsign color) of this route.
	Color string `protobuf:"bytes,4,opt,name=color,proto3" json:"color,omitempty"`
	// The collection of lines along this route.
	Lines []*RouteData_RouteLine `protobuf:"bytes,5,rep,name=lines,proto3" json:"lines,omitempty"`
}

func (x *RouteData) Reset() {
	*x = RouteData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_routes_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RouteData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RouteData) ProtoMessage() {}

func (x *RouteData) ProtoReflect() protoreflect.Message {
	mi := &file_routes_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RouteData.ProtoReflect.Descriptor instead.
func (*RouteData) Descriptor() ([]byte, []int) {
	return file_routes_proto_rawDescGZIP(), []int{5}
}

func (x *RouteData) GetRoute() Route {
	if x != nil {
		return x.Route
	}
	return Route_ROUTE_UNSPECIFIED
}

func (x *RouteData) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *RouteData) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RouteData) GetColor() string {
	if x != nil {
		return x.Color
	}
	return ""
}

func (x *RouteData) GetLines() []*RouteData_RouteLine {
	if x != nil {
		return x.Lines
	}
	return nil
}

// Represents a single line within this route (think direction of travel).
type RouteData_RouteLine struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The friendly name of this route line.
	DisplayName string `protobuf:"bytes,1,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	// The headsign displayed when a train is traveling along this line.
	Headsign string `protobuf:"bytes,2,opt,name=headsign,proto3" json:"headsign,omitempty"`
	// The direction of travel.
	Direction Direction `protobuf:"varint,3,opt,name=direction,proto3,enum=path_api.v1.Direction" json:"direction,omitempty"`
}

func (x *RouteData_RouteLine) Reset() {
	*x = RouteData_RouteLine{}
	if protoimpl.UnsafeEnabled {
		mi := &file_routes_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RouteData_RouteLine) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RouteData_RouteLine) ProtoMessage() {}

func (x *RouteData_RouteLine) ProtoReflect() protoreflect.Message {
	mi := &file_routes_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RouteData_RouteLine.ProtoReflect.Descriptor instead.
func (*RouteData_RouteLine) Descriptor() ([]byte, []int) {
	return file_routes_proto_rawDescGZIP(), []int{5, 0}
}

func (x *RouteData_RouteLine) GetDisplayName() string {
	if x != nil {
		return x.DisplayName
	}
	return ""
}

func (x *RouteData_RouteLine) GetHeadsign() string {
	if x != nil {
		return x.Headsign
	}
	return ""
}

func (x *RouteData_RouteLine) GetDirection() Direction {
	if x != nil {
		return x.Direction
	}
	return Direction_DIRECTION_UNSPECIFIED
}

var File_routes_proto protoreflect.FileDescriptor

var file_routes_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b,
	0x70, 0x61, 0x74, 0x68, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4f, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x52,
	0x6f, 0x75, 0x74, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09,
	0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61, 0x67,
	0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70,
	0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x6c, 0x0a, 0x12, 0x4c, 0x69, 0x73, 0x74,
	0x52, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e,
	0x0a, 0x06, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16,
	0x2e, 0x70, 0x61, 0x74, 0x68, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x6f, 0x75,
	0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x06, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x12, 0x26,
	0x0a, 0x0f, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6e, 0x65, 0x78, 0x74, 0x50, 0x61, 0x67,
	0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x3b, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x52, 0x6f, 0x75,
	0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x05, 0x72, 0x6f, 0x75,
	0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x70, 0x61, 0x74, 0x68, 0x5f,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x52, 0x05, 0x72, 0x6f,
	0x75, 0x74, 0x65, 0x22, 0x7f, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x53,
	0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x28,
	0x0a, 0x05, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e,
	0x70, 0x61, 0x74, 0x68, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x6f, 0x75, 0x74,
	0x65, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65,
	0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67,
	0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x61, 0x67, 0x65, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x42, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x52, 0x6f, 0x75, 0x74, 0x65,
	0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x26, 0x0a, 0x0f, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6e, 0x65, 0x78, 0x74, 0x50,
	0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0xaa, 0x02, 0x0a, 0x09, 0x52, 0x6f, 0x75,
	0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x28, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x70, 0x61, 0x74, 0x68, 0x5f, 0x61, 0x70, 0x69,
	0x2e, 0x76, 0x31, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x74, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x12, 0x36, 0x0a, 0x05, 0x6c, 0x69,
	0x6e, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x70, 0x61, 0x74, 0x68,
	0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x44, 0x61, 0x74,
	0x61, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x4c, 0x69, 0x6e, 0x65, 0x52, 0x05, 0x6c, 0x69, 0x6e,
	0x65, 0x73, 0x1a, 0x80, 0x01, 0x0a, 0x09, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x4c, 0x69, 0x6e, 0x65,
	0x12, 0x21, 0x0a, 0x0c, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x68, 0x65, 0x61, 0x64, 0x73, 0x69, 0x67, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x65, 0x61, 0x64, 0x73, 0x69, 0x67, 0x6e, 0x12,
	0x34, 0x0a, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x16, 0x2e, 0x70, 0x61, 0x74, 0x68, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31,
	0x2e, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x64, 0x69, 0x72, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x32, 0xd0, 0x02, 0x0a, 0x06, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x73,
	0x12, 0x61, 0x0a, 0x0a, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x12, 0x1e,
	0x2e, 0x70, 0x61, 0x74, 0x68, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f,
	0x2e, 0x70, 0x61, 0x74, 0x68, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x12, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0c, 0x12, 0x0a, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x6f, 0x75,
	0x74, 0x65, 0x73, 0x12, 0x5c, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x12,
	0x1c, 0x2e, 0x70, 0x61, 0x74, 0x68, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65,
	0x74, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e,
	0x70, 0x61, 0x74, 0x68, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x6f, 0x75, 0x74,
	0x65, 0x44, 0x61, 0x74, 0x61, 0x22, 0x1a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x14, 0x12, 0x12, 0x2f,
	0x76, 0x31, 0x2f, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x2f, 0x7b, 0x72, 0x6f, 0x75, 0x74, 0x65,
	0x7d, 0x12, 0x84, 0x01, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x53, 0x63,
	0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x24, 0x2e, 0x70, 0x61, 0x74, 0x68, 0x5f, 0x61, 0x70,
	0x69, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x53, 0x63, 0x68,
	0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x70,
	0x61, 0x74, 0x68, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x6f,
	0x75, 0x74, 0x65, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x23, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1d, 0x12, 0x1b, 0x2f, 0x76, 0x31,
	0x2f, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x2f, 0x7b, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x7d, 0x2f,
	0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_routes_proto_rawDescOnce sync.Once
	file_routes_proto_rawDescData = file_routes_proto_rawDesc
)

func file_routes_proto_rawDescGZIP() []byte {
	file_routes_proto_rawDescOnce.Do(func() {
		file_routes_proto_rawDescData = protoimpl.X.CompressGZIP(file_routes_proto_rawDescData)
	})
	return file_routes_proto_rawDescData
}

var file_routes_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_routes_proto_goTypes = []interface{}{
	(*ListRoutesRequest)(nil),        // 0: path_api.v1.ListRoutesRequest
	(*ListRoutesResponse)(nil),       // 1: path_api.v1.ListRoutesResponse
	(*GetRouteRequest)(nil),          // 2: path_api.v1.GetRouteRequest
	(*GetRouteScheduleRequest)(nil),  // 3: path_api.v1.GetRouteScheduleRequest
	(*GetRouteScheduleResponse)(nil), // 4: path_api.v1.GetRouteScheduleResponse
	(*RouteData)(nil),                // 5: path_api.v1.RouteData
	(*RouteData_RouteLine)(nil),      // 6: path_api.v1.RouteData.RouteLine
	(Route)(0),                       // 7: path_api.v1.Route
	(Direction)(0),                   // 8: path_api.v1.Direction
}
var file_routes_proto_depIdxs = []int32{
	5, // 0: path_api.v1.ListRoutesResponse.routes:type_name -> path_api.v1.RouteData
	7, // 1: path_api.v1.GetRouteRequest.route:type_name -> path_api.v1.Route
	7, // 2: path_api.v1.GetRouteScheduleRequest.route:type_name -> path_api.v1.Route
	7, // 3: path_api.v1.RouteData.route:type_name -> path_api.v1.Route
	6, // 4: path_api.v1.RouteData.lines:type_name -> path_api.v1.RouteData.RouteLine
	8, // 5: path_api.v1.RouteData.RouteLine.direction:type_name -> path_api.v1.Direction
	0, // 6: path_api.v1.Routes.ListRoutes:input_type -> path_api.v1.ListRoutesRequest
	2, // 7: path_api.v1.Routes.GetRoute:input_type -> path_api.v1.GetRouteRequest
	3, // 8: path_api.v1.Routes.GetRouteSchedule:input_type -> path_api.v1.GetRouteScheduleRequest
	1, // 9: path_api.v1.Routes.ListRoutes:output_type -> path_api.v1.ListRoutesResponse
	5, // 10: path_api.v1.Routes.GetRoute:output_type -> path_api.v1.RouteData
	4, // 11: path_api.v1.Routes.GetRouteSchedule:output_type -> path_api.v1.GetRouteScheduleResponse
	9, // [9:12] is the sub-list for method output_type
	6, // [6:9] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_routes_proto_init() }
func file_routes_proto_init() {
	if File_routes_proto != nil {
		return
	}
	file_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_routes_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListRoutesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_routes_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListRoutesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_routes_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRouteRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_routes_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRouteScheduleRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_routes_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRouteScheduleResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_routes_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RouteData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_routes_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RouteData_RouteLine); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_routes_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_routes_proto_goTypes,
		DependencyIndexes: file_routes_proto_depIdxs,
		MessageInfos:      file_routes_proto_msgTypes,
	}.Build()
	File_routes_proto = out.File
	file_routes_proto_rawDesc = nil
	file_routes_proto_goTypes = nil
	file_routes_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// RoutesClient is the client API for Routes service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RoutesClient interface {
	// Lists all routes within the PATH system.
	ListRoutes(ctx context.Context, in *ListRoutesRequest, opts ...grpc.CallOption) (*ListRoutesResponse, error)
	// Gets information about a single route.
	GetRoute(ctx context.Context, in *GetRouteRequest, opts ...grpc.CallOption) (*RouteData, error)
	// Gets the posted train schedule for a route.
	GetRouteSchedule(ctx context.Context, in *GetRouteScheduleRequest, opts ...grpc.CallOption) (*GetRouteScheduleResponse, error)
}

type routesClient struct {
	cc grpc.ClientConnInterface
}

func NewRoutesClient(cc grpc.ClientConnInterface) RoutesClient {
	return &routesClient{cc}
}

func (c *routesClient) ListRoutes(ctx context.Context, in *ListRoutesRequest, opts ...grpc.CallOption) (*ListRoutesResponse, error) {
	out := new(ListRoutesResponse)
	err := c.cc.Invoke(ctx, "/path_api.v1.Routes/ListRoutes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *routesClient) GetRoute(ctx context.Context, in *GetRouteRequest, opts ...grpc.CallOption) (*RouteData, error) {
	out := new(RouteData)
	err := c.cc.Invoke(ctx, "/path_api.v1.Routes/GetRoute", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *routesClient) GetRouteSchedule(ctx context.Context, in *GetRouteScheduleRequest, opts ...grpc.CallOption) (*GetRouteScheduleResponse, error) {
	out := new(GetRouteScheduleResponse)
	err := c.cc.Invoke(ctx, "/path_api.v1.Routes/GetRouteSchedule", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RoutesServer is the server API for Routes service.
type RoutesServer interface {
	// Lists all routes within the PATH system.
	ListRoutes(context.Context, *ListRoutesRequest) (*ListRoutesResponse, error)
	// Gets information about a single route.
	GetRoute(context.Context, *GetRouteRequest) (*RouteData, error)
	// Gets the posted train schedule for a route.
	GetRouteSchedule(context.Context, *GetRouteScheduleRequest) (*GetRouteScheduleResponse, error)
}

// UnimplementedRoutesServer can be embedded to have forward compatible implementations.
type UnimplementedRoutesServer struct {
}

func (*UnimplementedRoutesServer) ListRoutes(context.Context, *ListRoutesRequest) (*ListRoutesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRoutes not implemented")
}
func (*UnimplementedRoutesServer) GetRoute(context.Context, *GetRouteRequest) (*RouteData, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoute not implemented")
}
func (*UnimplementedRoutesServer) GetRouteSchedule(context.Context, *GetRouteScheduleRequest) (*GetRouteScheduleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRouteSchedule not implemented")
}

func RegisterRoutesServer(s *grpc.Server, srv RoutesServer) {
	s.RegisterService(&_Routes_serviceDesc, srv)
}

func _Routes_ListRoutes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRoutesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoutesServer).ListRoutes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/path_api.v1.Routes/ListRoutes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoutesServer).ListRoutes(ctx, req.(*ListRoutesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Routes_GetRoute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRouteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoutesServer).GetRoute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/path_api.v1.Routes/GetRoute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoutesServer).GetRoute(ctx, req.(*GetRouteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Routes_GetRouteSchedule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRouteScheduleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoutesServer).GetRouteSchedule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/path_api.v1.Routes/GetRouteSchedule",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoutesServer).GetRouteSchedule(ctx, req.(*GetRouteScheduleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Routes_serviceDesc = grpc.ServiceDesc{
	ServiceName: "path_api.v1.Routes",
	HandlerType: (*RoutesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListRoutes",
			Handler:    _Routes_ListRoutes_Handler,
		},
		{
			MethodName: "GetRoute",
			Handler:    _Routes_GetRoute_Handler,
		},
		{
			MethodName: "GetRouteSchedule",
			Handler:    _Routes_GetRouteSchedule_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "routes.proto",
}
