// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: protobuf/locations.proto

package protobuf

import (
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

// The request message containing the user's name.
type NewLocation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username  string  `protobuf:"bytes,1,opt,name=Username,proto3" json:"Username,omitempty"`
	Longitude float64 `protobuf:"fixed64,2,opt,name=Longitude,proto3" json:"Longitude,omitempty"`
	Latitude  float64 `protobuf:"fixed64,3,opt,name=Latitude,proto3" json:"Latitude,omitempty"`
}

func (x *NewLocation) Reset() {
	*x = NewLocation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_locations_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewLocation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewLocation) ProtoMessage() {}

func (x *NewLocation) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_locations_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewLocation.ProtoReflect.Descriptor instead.
func (*NewLocation) Descriptor() ([]byte, []int) {
	return file_protobuf_locations_proto_rawDescGZIP(), []int{0}
}

func (x *NewLocation) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *NewLocation) GetLongitude() float64 {
	if x != nil {
		return x.Longitude
	}
	return 0
}

func (x *NewLocation) GetLatitude() float64 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_locations_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_locations_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_protobuf_locations_proto_rawDescGZIP(), []int{1}
}

var File_protobuf_locations_proto protoreflect.FileDescriptor

var file_protobuf_locations_proto_rawDesc = []byte{
	0x0a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x6c, 0x6f, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x67, 0x72, 0x69, 0x64,
	0x67, 0x6f, 0x22, 0x63, 0x0a, 0x0b, 0x4e, 0x65, 0x77, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x4c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x09, 0x4c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x4c,
	0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x4c,
	0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x32, 0x3b, 0x0a, 0x09, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x2e, 0x0a,
	0x06, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x12, 0x13, 0x2e, 0x67, 0x72, 0x69, 0x64, 0x67, 0x6f,
	0x2e, 0x4e, 0x65, 0x77, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x0d, 0x2e, 0x67,
	0x72, 0x69, 0x64, 0x67, 0x6f, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x0c, 0x5a,
	0x0a, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_protobuf_locations_proto_rawDescOnce sync.Once
	file_protobuf_locations_proto_rawDescData = file_protobuf_locations_proto_rawDesc
)

func file_protobuf_locations_proto_rawDescGZIP() []byte {
	file_protobuf_locations_proto_rawDescOnce.Do(func() {
		file_protobuf_locations_proto_rawDescData = protoimpl.X.CompressGZIP(file_protobuf_locations_proto_rawDescData)
	})
	return file_protobuf_locations_proto_rawDescData
}

var file_protobuf_locations_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_protobuf_locations_proto_goTypes = []interface{}{
	(*NewLocation)(nil), // 0: gridgo.NewLocation
	(*Empty)(nil),       // 1: gridgo.Empty
}
var file_protobuf_locations_proto_depIdxs = []int32{
	0, // 0: gridgo.Locations.Insert:input_type -> gridgo.NewLocation
	1, // 1: gridgo.Locations.Insert:output_type -> gridgo.Empty
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_protobuf_locations_proto_init() }
func file_protobuf_locations_proto_init() {
	if File_protobuf_locations_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protobuf_locations_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewLocation); i {
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
		file_protobuf_locations_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
			RawDescriptor: file_protobuf_locations_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protobuf_locations_proto_goTypes,
		DependencyIndexes: file_protobuf_locations_proto_depIdxs,
		MessageInfos:      file_protobuf_locations_proto_msgTypes,
	}.Build()
	File_protobuf_locations_proto = out.File
	file_protobuf_locations_proto_rawDesc = nil
	file_protobuf_locations_proto_goTypes = nil
	file_protobuf_locations_proto_depIdxs = nil
}
