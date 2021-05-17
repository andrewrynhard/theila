// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: rpc/resource.proto

package rpc

import (
	proto "github.com/golang/protobuf/proto"
	resource "github.com/talos-systems/talos/pkg/machinery/api/resource"
	common "github.com/talos-systems/theila/api/common"
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

type GetFromClusterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Resource *resource.GetRequest `protobuf:"bytes,1,opt,name=resource,proto3" json:"resource,omitempty"`
	// Context to use to get the resource.
	Context *common.Context `protobuf:"bytes,2,opt,name=context,proto3" json:"context,omitempty"`
	// Data source to use to get the resource.
	Source common.Source `protobuf:"varint,3,opt,name=source,proto3,enum=common.Source" json:"source,omitempty"`
}

func (x *GetFromClusterRequest) Reset() {
	*x = GetFromClusterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_resource_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFromClusterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFromClusterRequest) ProtoMessage() {}

func (x *GetFromClusterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_resource_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFromClusterRequest.ProtoReflect.Descriptor instead.
func (*GetFromClusterRequest) Descriptor() ([]byte, []int) {
	return file_rpc_resource_proto_rawDescGZIP(), []int{0}
}

func (x *GetFromClusterRequest) GetResource() *resource.GetRequest {
	if x != nil {
		return x.Resource
	}
	return nil
}

func (x *GetFromClusterRequest) GetContext() *common.Context {
	if x != nil {
		return x.Context
	}
	return nil
}

func (x *GetFromClusterRequest) GetSource() common.Source {
	if x != nil {
		return x.Source
	}
	return common.Source_Kubernetes
}

type GetFromClusterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Body should contain JSON encoded spec.
	Body string `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *GetFromClusterResponse) Reset() {
	*x = GetFromClusterResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_resource_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFromClusterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFromClusterResponse) ProtoMessage() {}

func (x *GetFromClusterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_resource_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFromClusterResponse.ProtoReflect.Descriptor instead.
func (*GetFromClusterResponse) Descriptor() ([]byte, []int) {
	return file_rpc_resource_proto_rawDescGZIP(), []int{1}
}

func (x *GetFromClusterResponse) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

type ListFromClusterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Resource *resource.ListRequest `protobuf:"bytes,1,opt,name=resource,proto3" json:"resource,omitempty"`
	// Context to use to get the resource.
	Context *common.Context `protobuf:"bytes,2,opt,name=context,proto3" json:"context,omitempty"`
	// Data source to use to get the resource.
	Source common.Source `protobuf:"varint,3,opt,name=source,proto3,enum=common.Source" json:"source,omitempty"`
	// Selectors allow filtering list results by labels.
	Selectors []string `protobuf:"bytes,4,rep,name=selectors,proto3" json:"selectors,omitempty"`
}

func (x *ListFromClusterRequest) Reset() {
	*x = ListFromClusterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_resource_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListFromClusterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListFromClusterRequest) ProtoMessage() {}

func (x *ListFromClusterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_resource_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListFromClusterRequest.ProtoReflect.Descriptor instead.
func (*ListFromClusterRequest) Descriptor() ([]byte, []int) {
	return file_rpc_resource_proto_rawDescGZIP(), []int{2}
}

func (x *ListFromClusterRequest) GetResource() *resource.ListRequest {
	if x != nil {
		return x.Resource
	}
	return nil
}

func (x *ListFromClusterRequest) GetContext() *common.Context {
	if x != nil {
		return x.Context
	}
	return nil
}

func (x *ListFromClusterRequest) GetSource() common.Source {
	if x != nil {
		return x.Source
	}
	return common.Source_Kubernetes
}

func (x *ListFromClusterRequest) GetSelectors() []string {
	if x != nil {
		return x.Selectors
	}
	return nil
}

type ListFromClusterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Messages should contain JSON encoded list spec.
	Messages []string `protobuf:"bytes,1,rep,name=messages,proto3" json:"messages,omitempty"`
}

func (x *ListFromClusterResponse) Reset() {
	*x = ListFromClusterResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_resource_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListFromClusterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListFromClusterResponse) ProtoMessage() {}

func (x *ListFromClusterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_resource_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListFromClusterResponse.ProtoReflect.Descriptor instead.
func (*ListFromClusterResponse) Descriptor() ([]byte, []int) {
	return file_rpc_resource_proto_rawDescGZIP(), []int{3}
}

func (x *ListFromClusterResponse) GetMessages() []string {
	if x != nil {
		return x.Messages
	}
	return nil
}

var File_rpc_resource_proto protoreflect.FileDescriptor

var file_rpc_resource_proto_rawDesc = []byte{
	0x0a, 0x12, 0x72, 0x70, 0x63, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x1a, 0x1d,
	0x74, 0x61, 0x6c, 0x6f, 0x73, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2f, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x74, 0x68, 0x65, 0x69, 0x6c, 0x61, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x9c, 0x01, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x46, 0x72, 0x6f, 0x6d, 0x43, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x30, 0x0a, 0x08,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x52, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x29,
	0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0f, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74,
	0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x12, 0x26, 0x0a, 0x06, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0e, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x22, 0x2c, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x46, 0x72, 0x6f, 0x6d, 0x43, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x62,
	0x6f, 0x64, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x22,
	0xbc, 0x01, 0x0a, 0x16, 0x4c, 0x69, 0x73, 0x74, 0x46, 0x72, 0x6f, 0x6d, 0x43, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x31, 0x0a, 0x08, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x52, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x29, 0x0a,
	0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x52,
	0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x12, 0x26, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0e, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x22, 0x35,
	0x0a, 0x17, 0x4c, 0x69, 0x73, 0x74, 0x46, 0x72, 0x6f, 0x6d, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x73, 0x32, 0xaf, 0x01, 0x0a, 0x16, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x48, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x1f, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x72, 0x6f, 0x6d, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x72, 0x6f, 0x6d, 0x43, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4b, 0x0a, 0x04, 0x4c, 0x69,
	0x73, 0x74, 0x12, 0x20, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x46, 0x72, 0x6f, 0x6d, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e,
	0x4c, 0x69, 0x73, 0x74, 0x46, 0x72, 0x6f, 0x6d, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x29, 0x5a, 0x27, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x61, 0x6c, 0x6f, 0x73, 0x2d, 0x73, 0x79, 0x73, 0x74,
	0x65, 0x6d, 0x73, 0x2f, 0x74, 0x68, 0x65, 0x69, 0x6c, 0x61, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x72,
	0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_resource_proto_rawDescOnce sync.Once
	file_rpc_resource_proto_rawDescData = file_rpc_resource_proto_rawDesc
)

func file_rpc_resource_proto_rawDescGZIP() []byte {
	file_rpc_resource_proto_rawDescOnce.Do(func() {
		file_rpc_resource_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_resource_proto_rawDescData)
	})
	return file_rpc_resource_proto_rawDescData
}

var file_rpc_resource_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_rpc_resource_proto_goTypes = []interface{}{
	(*GetFromClusterRequest)(nil),   // 0: resource.GetFromClusterRequest
	(*GetFromClusterResponse)(nil),  // 1: resource.GetFromClusterResponse
	(*ListFromClusterRequest)(nil),  // 2: resource.ListFromClusterRequest
	(*ListFromClusterResponse)(nil), // 3: resource.ListFromClusterResponse
	(*resource.GetRequest)(nil),     // 4: resource.GetRequest
	(*common.Context)(nil),          // 5: common.Context
	(common.Source)(0),              // 6: common.Source
	(*resource.ListRequest)(nil),    // 7: resource.ListRequest
}
var file_rpc_resource_proto_depIdxs = []int32{
	4, // 0: resource.GetFromClusterRequest.resource:type_name -> resource.GetRequest
	5, // 1: resource.GetFromClusterRequest.context:type_name -> common.Context
	6, // 2: resource.GetFromClusterRequest.source:type_name -> common.Source
	7, // 3: resource.ListFromClusterRequest.resource:type_name -> resource.ListRequest
	5, // 4: resource.ListFromClusterRequest.context:type_name -> common.Context
	6, // 5: resource.ListFromClusterRequest.source:type_name -> common.Source
	0, // 6: resource.ClusterResourceService.Get:input_type -> resource.GetFromClusterRequest
	2, // 7: resource.ClusterResourceService.List:input_type -> resource.ListFromClusterRequest
	1, // 8: resource.ClusterResourceService.Get:output_type -> resource.GetFromClusterResponse
	3, // 9: resource.ClusterResourceService.List:output_type -> resource.ListFromClusterResponse
	8, // [8:10] is the sub-list for method output_type
	6, // [6:8] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_rpc_resource_proto_init() }
func file_rpc_resource_proto_init() {
	if File_rpc_resource_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_resource_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFromClusterRequest); i {
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
		file_rpc_resource_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFromClusterResponse); i {
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
		file_rpc_resource_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListFromClusterRequest); i {
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
		file_rpc_resource_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListFromClusterResponse); i {
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
			RawDescriptor: file_rpc_resource_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_resource_proto_goTypes,
		DependencyIndexes: file_rpc_resource_proto_depIdxs,
		MessageInfos:      file_rpc_resource_proto_msgTypes,
	}.Build()
	File_rpc_resource_proto = out.File
	file_rpc_resource_proto_rawDesc = nil
	file_rpc_resource_proto_goTypes = nil
	file_rpc_resource_proto_depIdxs = nil
}