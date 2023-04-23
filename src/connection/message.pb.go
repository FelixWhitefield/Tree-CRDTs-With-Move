// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.22.3
// source: message.proto

package connection

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

type PeerAddresses struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PeerAddrs []string `protobuf:"bytes,1,rep,name=peerAddrs,proto3" json:"peerAddrs,omitempty"`
}

func (x *PeerAddresses) Reset() {
	*x = PeerAddresses{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PeerAddresses) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PeerAddresses) ProtoMessage() {}

func (x *PeerAddresses) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PeerAddresses.ProtoReflect.Descriptor instead.
func (*PeerAddresses) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{0}
}

func (x *PeerAddresses) GetPeerAddrs() []string {
	if x != nil {
		return x.PeerAddrs
	}
	return nil
}

type PeerID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id []byte `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *PeerID) Reset() {
	*x = PeerID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PeerID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PeerID) ProtoMessage() {}

func (x *PeerID) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PeerID.ProtoReflect.Descriptor instead.
func (*PeerID) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{1}
}

func (x *PeerID) GetId() []byte {
	if x != nil {
		return x.Id
	}
	return nil
}

type OperationMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id []byte `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Op []byte `protobuf:"bytes,2,opt,name=op,proto3" json:"op,omitempty"` // Will be a protobuf encoded operation
}

func (x *OperationMsg) Reset() {
	*x = OperationMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OperationMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OperationMsg) ProtoMessage() {}

func (x *OperationMsg) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OperationMsg.ProtoReflect.Descriptor instead.
func (*OperationMsg) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{2}
}

func (x *OperationMsg) GetId() []byte {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *OperationMsg) GetOp() []byte {
	if x != nil {
		return x.Op
	}
	return nil
}

type OperationAck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id  []byte `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Ack bool   `protobuf:"varint,2,opt,name=ack,proto3" json:"ack,omitempty"`
}

func (x *OperationAck) Reset() {
	*x = OperationAck{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OperationAck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OperationAck) ProtoMessage() {}

func (x *OperationAck) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OperationAck.ProtoReflect.Descriptor instead.
func (*OperationAck) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{3}
}

func (x *OperationAck) GetId() []byte {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *OperationAck) GetAck() bool {
	if x != nil {
		return x.Ack
	}
	return false
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Message:
	//
	//	*Message_PeerID
	//	*Message_PeerAddresses
	//	*Message_Operation
	//	*Message_OperationAck
	Message isMessage_Message `protobuf_oneof:"message"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{4}
}

func (m *Message) GetMessage() isMessage_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *Message) GetPeerID() *PeerID {
	if x, ok := x.GetMessage().(*Message_PeerID); ok {
		return x.PeerID
	}
	return nil
}

func (x *Message) GetPeerAddresses() *PeerAddresses {
	if x, ok := x.GetMessage().(*Message_PeerAddresses); ok {
		return x.PeerAddresses
	}
	return nil
}

func (x *Message) GetOperation() *OperationMsg {
	if x, ok := x.GetMessage().(*Message_Operation); ok {
		return x.Operation
	}
	return nil
}

func (x *Message) GetOperationAck() *OperationAck {
	if x, ok := x.GetMessage().(*Message_OperationAck); ok {
		return x.OperationAck
	}
	return nil
}

type isMessage_Message interface {
	isMessage_Message()
}

type Message_PeerID struct {
	PeerID *PeerID `protobuf:"bytes,1,opt,name=peerID,proto3,oneof"`
}

type Message_PeerAddresses struct {
	PeerAddresses *PeerAddresses `protobuf:"bytes,2,opt,name=peerAddresses,proto3,oneof"`
}

type Message_Operation struct {
	Operation *OperationMsg `protobuf:"bytes,3,opt,name=operation,proto3,oneof"`
}

type Message_OperationAck struct {
	OperationAck *OperationAck `protobuf:"bytes,4,opt,name=operationAck,proto3,oneof"`
}

func (*Message_PeerID) isMessage_Message() {}

func (*Message_PeerAddresses) isMessage_Message() {}

func (*Message_Operation) isMessage_Message() {}

func (*Message_OperationAck) isMessage_Message() {}

var File_message_proto protoreflect.FileDescriptor

var file_message_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x05, 0x6d, 0x70, 0x72, 0x6f, 0x74, 0x22, 0x2d, 0x0a, 0x0d, 0x50, 0x65, 0x65, 0x72, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x65, 0x65, 0x72, 0x41,
	0x64, 0x64, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x70, 0x65, 0x65, 0x72,
	0x41, 0x64, 0x64, 0x72, 0x73, 0x22, 0x18, 0x0a, 0x06, 0x50, 0x65, 0x65, 0x72, 0x49, 0x44, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x69, 0x64, 0x22,
	0x2e, 0x0a, 0x0c, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x73, 0x67, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x0e, 0x0a, 0x02, 0x6f, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x6f, 0x70, 0x22,
	0x30, 0x0a, 0x0c, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x63, 0x6b, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x10, 0x0a, 0x03, 0x61, 0x63, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x03, 0x61, 0x63,
	0x6b, 0x22, 0xeb, 0x01, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x27, 0x0a,
	0x06, 0x70, 0x65, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e,
	0x6d, 0x70, 0x72, 0x6f, 0x74, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x49, 0x44, 0x48, 0x00, 0x52, 0x06,
	0x70, 0x65, 0x65, 0x72, 0x49, 0x44, 0x12, 0x3c, 0x0a, 0x0d, 0x70, 0x65, 0x65, 0x72, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x6d, 0x70, 0x72, 0x6f, 0x74, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x65, 0x73, 0x48, 0x00, 0x52, 0x0d, 0x70, 0x65, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x65, 0x73, 0x12, 0x33, 0x0a, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6d, 0x70, 0x72, 0x6f, 0x74, 0x2e,
	0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x73, 0x67, 0x48, 0x00, 0x52, 0x09,
	0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x39, 0x0a, 0x0c, 0x6f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x63, 0x6b, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x13, 0x2e, 0x6d, 0x70, 0x72, 0x6f, 0x74, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x41, 0x63, 0x6b, 0x48, 0x00, 0x52, 0x0c, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x41, 0x63, 0x6b, 0x42, 0x09, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42,
	0x3c, 0x5a, 0x3a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x46, 0x65,
	0x6c, 0x69, 0x78, 0x57, 0x68, 0x69, 0x74, 0x65, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x2f, 0x54, 0x72,
	0x65, 0x65, 0x2d, 0x43, 0x52, 0x44, 0x54, 0x73, 0x2d, 0x57, 0x69, 0x74, 0x68, 0x2d, 0x4d, 0x6f,
	0x76, 0x65, 0x2f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_message_proto_rawDescOnce sync.Once
	file_message_proto_rawDescData = file_message_proto_rawDesc
)

func file_message_proto_rawDescGZIP() []byte {
	file_message_proto_rawDescOnce.Do(func() {
		file_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_message_proto_rawDescData)
	})
	return file_message_proto_rawDescData
}

var file_message_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_message_proto_goTypes = []interface{}{
	(*PeerAddresses)(nil), // 0: mprot.PeerAddresses
	(*PeerID)(nil),        // 1: mprot.PeerID
	(*OperationMsg)(nil),  // 2: mprot.OperationMsg
	(*OperationAck)(nil),  // 3: mprot.OperationAck
	(*Message)(nil),       // 4: mprot.Message
}
var file_message_proto_depIdxs = []int32{
	1, // 0: mprot.Message.peerID:type_name -> mprot.PeerID
	0, // 1: mprot.Message.peerAddresses:type_name -> mprot.PeerAddresses
	2, // 2: mprot.Message.operation:type_name -> mprot.OperationMsg
	3, // 3: mprot.Message.operationAck:type_name -> mprot.OperationAck
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_message_proto_init() }
func file_message_proto_init() {
	if File_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PeerAddresses); i {
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
		file_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PeerID); i {
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
		file_message_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OperationMsg); i {
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
		file_message_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OperationAck); i {
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
		file_message_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
	file_message_proto_msgTypes[4].OneofWrappers = []interface{}{
		(*Message_PeerID)(nil),
		(*Message_PeerAddresses)(nil),
		(*Message_Operation)(nil),
		(*Message_OperationAck)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_message_proto_goTypes,
		DependencyIndexes: file_message_proto_depIdxs,
		MessageInfos:      file_message_proto_msgTypes,
	}.Build()
	File_message_proto = out.File
	file_message_proto_rawDesc = nil
	file_message_proto_goTypes = nil
	file_message_proto_depIdxs = nil
}
