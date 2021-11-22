// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: raft_log.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Symbols defined in public import of google/protobuf/timestamp.proto.

type Timestamp = timestamppb.Timestamp

type RaftLog struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index      uint64                 `protobuf:"varint,1,opt,name=Index,proto3" json:"Index,omitempty"`
	Term       uint64                 `protobuf:"varint,2,opt,name=Term,proto3" json:"Term,omitempty"`
	Data       []byte                 `protobuf:"bytes,3,opt,name=Data,proto3" json:"Data,omitempty"`
	Extensions []byte                 `protobuf:"bytes,4,opt,name=Extensions,proto3" json:"Extensions,omitempty"`
	AppendedAt *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=AppendedAt,proto3" json:"AppendedAt,omitempty"`
	Type       uint32                 `protobuf:"varint,6,opt,name=Type,proto3" json:"Type,omitempty"`
}

func (x *RaftLog) Reset() {
	*x = RaftLog{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raft_log_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RaftLog) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RaftLog) ProtoMessage() {}

func (x *RaftLog) ProtoReflect() protoreflect.Message {
	mi := &file_raft_log_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RaftLog.ProtoReflect.Descriptor instead.
func (*RaftLog) Descriptor() ([]byte, []int) {
	return file_raft_log_proto_rawDescGZIP(), []int{0}
}

func (x *RaftLog) GetIndex() uint64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *RaftLog) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *RaftLog) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *RaftLog) GetExtensions() []byte {
	if x != nil {
		return x.Extensions
	}
	return nil
}

func (x *RaftLog) GetAppendedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.AppendedAt
	}
	return nil
}

func (x *RaftLog) GetType() uint32 {
	if x != nil {
		return x.Type
	}
	return 0
}

var File_raft_log_proto protoreflect.FileDescriptor

var file_raft_log_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x72, 0x61, 0x66, 0x74, 0x5f, 0x6c, 0x6f, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x02, 0x70, 0x62, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb7, 0x01, 0x0a, 0x07, 0x52, 0x61, 0x66, 0x74, 0x4c, 0x6f,
	0x67, 0x12, 0x14, 0x0a, 0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x44,
	0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x1e, 0x0a, 0x0a, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x0a, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12,
	0x3a, 0x0a, 0x0a, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x65, 0x64, 0x41, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x0a, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x65, 0x64, 0x41, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x54,
	0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x42,
	0x07, 0x5a, 0x05, 0x2e, 0x2f, 0x3b, 0x70, 0x62, 0x50, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_raft_log_proto_rawDescOnce sync.Once
	file_raft_log_proto_rawDescData = file_raft_log_proto_rawDesc
)

func file_raft_log_proto_rawDescGZIP() []byte {
	file_raft_log_proto_rawDescOnce.Do(func() {
		file_raft_log_proto_rawDescData = protoimpl.X.CompressGZIP(file_raft_log_proto_rawDescData)
	})
	return file_raft_log_proto_rawDescData
}

var file_raft_log_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_raft_log_proto_goTypes = []interface{}{
	(*RaftLog)(nil),               // 0: pb.RaftLog
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_raft_log_proto_depIdxs = []int32{
	1, // 0: pb.RaftLog.AppendedAt:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_raft_log_proto_init() }
func file_raft_log_proto_init() {
	if File_raft_log_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_raft_log_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RaftLog); i {
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
			RawDescriptor: file_raft_log_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_raft_log_proto_goTypes,
		DependencyIndexes: file_raft_log_proto_depIdxs,
		MessageInfos:      file_raft_log_proto_msgTypes,
	}.Build()
	File_raft_log_proto = out.File
	file_raft_log_proto_rawDesc = nil
	file_raft_log_proto_goTypes = nil
	file_raft_log_proto_depIdxs = nil
}
