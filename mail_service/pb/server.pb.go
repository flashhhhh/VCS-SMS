// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.12.4
// source: proto/server.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type EmptyRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EmptyRequest) Reset() {
	*x = EmptyRequest{}
	mi := &file_proto_server_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EmptyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmptyRequest) ProtoMessage() {}

func (x *EmptyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmptyRequest.ProtoReflect.Descriptor instead.
func (*EmptyRequest) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{0}
}

type AddressesResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Addresses     []*AddressInfo         `protobuf:"bytes,1,rep,name=addresses,proto3" json:"addresses,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AddressesResponse) Reset() {
	*x = AddressesResponse{}
	mi := &file_proto_server_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddressesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddressesResponse) ProtoMessage() {}

func (x *AddressesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddressesResponse.ProtoReflect.Descriptor instead.
func (*AddressesResponse) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{1}
}

func (x *AddressesResponse) GetAddresses() []*AddressInfo {
	if x != nil {
		return x.Addresses
	}
	return nil
}

type AddressInfo struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Address       string                 `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AddressInfo) Reset() {
	*x = AddressInfo{}
	mi := &file_proto_server_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddressInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddressInfo) ProtoMessage() {}

func (x *AddressInfo) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddressInfo.ProtoReflect.Descriptor instead.
func (*AddressInfo) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{2}
}

func (x *AddressInfo) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *AddressInfo) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type GetServerInformationRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	StartTime     int64                  `protobuf:"varint,1,opt,name=startTime,proto3" json:"startTime,omitempty"` // timestamp in unix format
	EndTime       int64                  `protobuf:"varint,2,opt,name=endTime,proto3" json:"endTime,omitempty"`     // timestamp in unix format
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetServerInformationRequest) Reset() {
	*x = GetServerInformationRequest{}
	mi := &file_proto_server_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetServerInformationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetServerInformationRequest) ProtoMessage() {}

func (x *GetServerInformationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetServerInformationRequest.ProtoReflect.Descriptor instead.
func (*GetServerInformationRequest) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{3}
}

func (x *GetServerInformationRequest) GetStartTime() int64 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

func (x *GetServerInformationRequest) GetEndTime() int64 {
	if x != nil {
		return x.EndTime
	}
	return 0
}

type GetServerInformationResponse struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	NumServers      int64                  `protobuf:"varint,1,opt,name=numServers,proto3" json:"numServers,omitempty"`
	NumOnServers    int64                  `protobuf:"varint,2,opt,name=numOnServers,proto3" json:"numOnServers,omitempty"`
	NumOffServers   int64                  `protobuf:"varint,3,opt,name=numOffServers,proto3" json:"numOffServers,omitempty"`
	MeanUptimeRatio float32                `protobuf:"fixed32,4,opt,name=meanUptimeRatio,proto3" json:"meanUptimeRatio,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *GetServerInformationResponse) Reset() {
	*x = GetServerInformationResponse{}
	mi := &file_proto_server_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetServerInformationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetServerInformationResponse) ProtoMessage() {}

func (x *GetServerInformationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetServerInformationResponse.ProtoReflect.Descriptor instead.
func (*GetServerInformationResponse) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{4}
}

func (x *GetServerInformationResponse) GetNumServers() int64 {
	if x != nil {
		return x.NumServers
	}
	return 0
}

func (x *GetServerInformationResponse) GetNumOnServers() int64 {
	if x != nil {
		return x.NumOnServers
	}
	return 0
}

func (x *GetServerInformationResponse) GetNumOffServers() int64 {
	if x != nil {
		return x.NumOffServers
	}
	return 0
}

func (x *GetServerInformationResponse) GetMeanUptimeRatio() float32 {
	if x != nil {
		return x.MeanUptimeRatio
	}
	return 0
}

var File_proto_server_proto protoreflect.FileDescriptor

const file_proto_server_proto_rawDesc = "" +
	"\n" +
	"\x12proto/server.proto\x12\x1dserver_administration_service\"\x0e\n" +
	"\fEmptyRequest\"]\n" +
	"\x11AddressesResponse\x12H\n" +
	"\taddresses\x18\x01 \x03(\v2*.server_administration_service.AddressInfoR\taddresses\"7\n" +
	"\vAddressInfo\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x18\n" +
	"\aaddress\x18\x02 \x01(\tR\aaddress\"U\n" +
	"\x1bGetServerInformationRequest\x12\x1c\n" +
	"\tstartTime\x18\x01 \x01(\x03R\tstartTime\x12\x18\n" +
	"\aendTime\x18\x02 \x01(\x03R\aendTime\"\xb2\x01\n" +
	"\x1cGetServerInformationResponse\x12\x1e\n" +
	"\n" +
	"numServers\x18\x01 \x01(\x03R\n" +
	"numServers\x12\"\n" +
	"\fnumOnServers\x18\x02 \x01(\x03R\fnumOnServers\x12$\n" +
	"\rnumOffServers\x18\x03 \x01(\x03R\rnumOffServers\x12(\n" +
	"\x0fmeanUptimeRatio\x18\x04 \x01(\x02R\x0fmeanUptimeRatio2\xa1\x02\n" +
	"\x1bServerAdministrationService\x12p\n" +
	"\x0fGetAllAddresses\x12+.server_administration_service.EmptyRequest\x1a0.server_administration_service.AddressesResponse\x12\x8f\x01\n" +
	"\x14GetServerInformation\x12:.server_administration_service.GetServerInformationRequest\x1a;.server_administration_service.GetServerInformationResponseB\x06Z\x04./pbb\x06proto3"

var (
	file_proto_server_proto_rawDescOnce sync.Once
	file_proto_server_proto_rawDescData []byte
)

func file_proto_server_proto_rawDescGZIP() []byte {
	file_proto_server_proto_rawDescOnce.Do(func() {
		file_proto_server_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_server_proto_rawDesc), len(file_proto_server_proto_rawDesc)))
	})
	return file_proto_server_proto_rawDescData
}

var file_proto_server_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_server_proto_goTypes = []any{
	(*EmptyRequest)(nil),                 // 0: server_administration_service.EmptyRequest
	(*AddressesResponse)(nil),            // 1: server_administration_service.AddressesResponse
	(*AddressInfo)(nil),                  // 2: server_administration_service.AddressInfo
	(*GetServerInformationRequest)(nil),  // 3: server_administration_service.GetServerInformationRequest
	(*GetServerInformationResponse)(nil), // 4: server_administration_service.GetServerInformationResponse
}
var file_proto_server_proto_depIdxs = []int32{
	2, // 0: server_administration_service.AddressesResponse.addresses:type_name -> server_administration_service.AddressInfo
	0, // 1: server_administration_service.ServerAdministrationService.GetAllAddresses:input_type -> server_administration_service.EmptyRequest
	3, // 2: server_administration_service.ServerAdministrationService.GetServerInformation:input_type -> server_administration_service.GetServerInformationRequest
	1, // 3: server_administration_service.ServerAdministrationService.GetAllAddresses:output_type -> server_administration_service.AddressesResponse
	4, // 4: server_administration_service.ServerAdministrationService.GetServerInformation:output_type -> server_administration_service.GetServerInformationResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_server_proto_init() }
func file_proto_server_proto_init() {
	if File_proto_server_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_server_proto_rawDesc), len(file_proto_server_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_server_proto_goTypes,
		DependencyIndexes: file_proto_server_proto_depIdxs,
		MessageInfos:      file_proto_server_proto_msgTypes,
	}.Build()
	File_proto_server_proto = out.File
	file_proto_server_proto_goTypes = nil
	file_proto_server_proto_depIdxs = nil
}
