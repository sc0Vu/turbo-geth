// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.3
// source: remote/kv.proto

package remote

import (
	proto "github.com/golang/protobuf/proto"
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

type SeekRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BucketName    []byte `protobuf:"bytes,1,opt,name=bucketName,proto3" json:"bucketName,omitempty"`
	SeekKey       []byte `protobuf:"bytes,2,opt,name=seekKey,proto3" json:"seekKey,omitempty"` // streaming start from this key
	Prefix        []byte `protobuf:"bytes,3,opt,name=prefix,proto3" json:"prefix,omitempty"`   // streaming stops when see first key without given prefix
	StartSreaming bool   `protobuf:"varint,4,opt,name=startSreaming,proto3" json:"startSreaming,omitempty"`
}

func (x *SeekRequest) Reset() {
	*x = SeekRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_remote_kv_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SeekRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SeekRequest) ProtoMessage() {}

func (x *SeekRequest) ProtoReflect() protoreflect.Message {
	mi := &file_remote_kv_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SeekRequest.ProtoReflect.Descriptor instead.
func (*SeekRequest) Descriptor() ([]byte, []int) {
	return file_remote_kv_proto_rawDescGZIP(), []int{0}
}

func (x *SeekRequest) GetBucketName() []byte {
	if x != nil {
		return x.BucketName
	}
	return nil
}

func (x *SeekRequest) GetSeekKey() []byte {
	if x != nil {
		return x.SeekKey
	}
	return nil
}

func (x *SeekRequest) GetPrefix() []byte {
	if x != nil {
		return x.Prefix
	}
	return nil
}

func (x *SeekRequest) GetStartSreaming() bool {
	if x != nil {
		return x.StartSreaming
	}
	return false
}

type Pair struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Pair) Reset() {
	*x = Pair{}
	if protoimpl.UnsafeEnabled {
		mi := &file_remote_kv_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pair) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pair) ProtoMessage() {}

func (x *Pair) ProtoReflect() protoreflect.Message {
	mi := &file_remote_kv_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pair.ProtoReflect.Descriptor instead.
func (*Pair) Descriptor() ([]byte, []int) {
	return file_remote_kv_proto_rawDescGZIP(), []int{1}
}

func (x *Pair) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *Pair) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type PairKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	VSize uint64 `protobuf:"varint,2,opt,name=vSize,proto3" json:"vSize,omitempty"`
}

func (x *PairKey) Reset() {
	*x = PairKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_remote_kv_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PairKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PairKey) ProtoMessage() {}

func (x *PairKey) ProtoReflect() protoreflect.Message {
	mi := &file_remote_kv_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PairKey.ProtoReflect.Descriptor instead.
func (*PairKey) Descriptor() ([]byte, []int) {
	return file_remote_kv_proto_rawDescGZIP(), []int{2}
}

func (x *PairKey) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *PairKey) GetVSize() uint64 {
	if x != nil {
		return x.VSize
	}
	return 0
}

var File_remote_kv_proto protoreflect.FileDescriptor

var file_remote_kv_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2f, 0x6b, 0x76, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x05, 0x65, 0x74, 0x68, 0x64, 0x62, 0x22, 0x85, 0x01, 0x0a, 0x0b, 0x53, 0x65, 0x65,
	0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x62, 0x75, 0x63, 0x6b,
	0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x62, 0x75,
	0x63, 0x6b, 0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x65, 0x65, 0x6b,
	0x4b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x73, 0x65, 0x65, 0x6b, 0x4b,
	0x65, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x06, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x12, 0x24, 0x0a, 0x0d, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x53, 0x72, 0x65, 0x61, 0x6d, 0x69, 0x6e, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x0d, 0x73, 0x74, 0x61, 0x72, 0x74, 0x53, 0x72, 0x65, 0x61, 0x6d, 0x69, 0x6e, 0x67,
	0x22, 0x2e, 0x0a, 0x04, 0x50, 0x61, 0x69, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x22, 0x31, 0x0a, 0x07, 0x50, 0x61, 0x69, 0x72, 0x4b, 0x65, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x76, 0x53,
	0x69, 0x7a, 0x65, 0x32, 0x31, 0x0a, 0x02, 0x4b, 0x76, 0x12, 0x2b, 0x0a, 0x04, 0x53, 0x65, 0x65,
	0x6b, 0x12, 0x12, 0x2e, 0x65, 0x74, 0x68, 0x64, 0x62, 0x2e, 0x53, 0x65, 0x65, 0x6b, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0b, 0x2e, 0x65, 0x74, 0x68, 0x64, 0x62, 0x2e, 0x50, 0x61,
	0x69, 0x72, 0x28, 0x01, 0x30, 0x01, 0x42, 0x29, 0x0a, 0x10, 0x69, 0x6f, 0x2e, 0x74, 0x75, 0x72,
	0x62, 0x6f, 0x2d, 0x67, 0x65, 0x74, 0x68, 0x2e, 0x64, 0x62, 0x42, 0x02, 0x4b, 0x56, 0x50, 0x01,
	0x5a, 0x0f, 0x2e, 0x2f, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x3b, 0x72, 0x65, 0x6d, 0x6f, 0x74,
	0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_remote_kv_proto_rawDescOnce sync.Once
	file_remote_kv_proto_rawDescData = file_remote_kv_proto_rawDesc
)

func file_remote_kv_proto_rawDescGZIP() []byte {
	file_remote_kv_proto_rawDescOnce.Do(func() {
		file_remote_kv_proto_rawDescData = protoimpl.X.CompressGZIP(file_remote_kv_proto_rawDescData)
	})
	return file_remote_kv_proto_rawDescData
}

var file_remote_kv_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_remote_kv_proto_goTypes = []interface{}{
	(*SeekRequest)(nil), // 0: ethdb.SeekRequest
	(*Pair)(nil),        // 1: ethdb.Pair
	(*PairKey)(nil),     // 2: ethdb.PairKey
}
var file_remote_kv_proto_depIdxs = []int32{
	0, // 0: ethdb.Kv.Seek:input_type -> ethdb.SeekRequest
	1, // 1: ethdb.Kv.Seek:output_type -> ethdb.Pair
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_remote_kv_proto_init() }
func file_remote_kv_proto_init() {
	if File_remote_kv_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_remote_kv_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SeekRequest); i {
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
		file_remote_kv_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Pair); i {
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
		file_remote_kv_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PairKey); i {
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
			RawDescriptor: file_remote_kv_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_remote_kv_proto_goTypes,
		DependencyIndexes: file_remote_kv_proto_depIdxs,
		MessageInfos:      file_remote_kv_proto_msgTypes,
	}.Build()
	File_remote_kv_proto = out.File
	file_remote_kv_proto_rawDesc = nil
	file_remote_kv_proto_goTypes = nil
	file_remote_kv_proto_depIdxs = nil
}
