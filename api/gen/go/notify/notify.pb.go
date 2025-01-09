// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.26.1
// source: notify/notify.proto

package notifyv1

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
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

type NotifyType int32

const (
	NotifyType_moderate    NotifyType = 0
	NotifyType_significant NotifyType = 1
	NotifyType_alert       NotifyType = 2
)

// Enum value maps for NotifyType.
var (
	NotifyType_name = map[int32]string{
		0: "moderate",
		1: "significant",
		2: "alert",
	}
	NotifyType_value = map[string]int32{
		"moderate":    0,
		"significant": 1,
		"alert":       2,
	}
)

func (x NotifyType) Enum() *NotifyType {
	p := new(NotifyType)
	*p = x
	return p
}

func (x NotifyType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (NotifyType) Descriptor() protoreflect.EnumDescriptor {
	return file_notify_notify_proto_enumTypes[0].Descriptor()
}

func (NotifyType) Type() protoreflect.EnumType {
	return &file_notify_notify_proto_enumTypes[0]
}

func (x NotifyType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use NotifyType.Descriptor instead.
func (NotifyType) EnumDescriptor() ([]byte, []int) {
	return file_notify_notify_proto_rawDescGZIP(), []int{0}
}

type MailNotify struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Subject string `protobuf:"bytes,1,opt,name=subject,proto3" json:"subject,omitempty"`
	Body    string `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *MailNotify) Reset() {
	*x = MailNotify{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notify_notify_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MailNotify) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MailNotify) ProtoMessage() {}

func (x *MailNotify) ProtoReflect() protoreflect.Message {
	mi := &file_notify_notify_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MailNotify.ProtoReflect.Descriptor instead.
func (*MailNotify) Descriptor() ([]byte, []int) {
	return file_notify_notify_proto_rawDescGZIP(), []int{0}
}

func (x *MailNotify) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *MailNotify) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

type PhoneNotify struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Body string `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *PhoneNotify) Reset() {
	*x = PhoneNotify{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notify_notify_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PhoneNotify) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PhoneNotify) ProtoMessage() {}

func (x *PhoneNotify) ProtoReflect() protoreflect.Message {
	mi := &file_notify_notify_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PhoneNotify.ProtoReflect.Descriptor instead.
func (*PhoneNotify) Descriptor() ([]byte, []int) {
	return file_notify_notify_proto_rawDescGZIP(), []int{1}
}

func (x *PhoneNotify) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

type Channels struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Mail  *MailNotify  `protobuf:"bytes,1,opt,name=mail,proto3,oneof" json:"mail,omitempty"`
	Phone *PhoneNotify `protobuf:"bytes,2,opt,name=phone,proto3,oneof" json:"phone,omitempty"`
}

func (x *Channels) Reset() {
	*x = Channels{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notify_notify_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Channels) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Channels) ProtoMessage() {}

func (x *Channels) ProtoReflect() protoreflect.Message {
	mi := &file_notify_notify_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Channels.ProtoReflect.Descriptor instead.
func (*Channels) Descriptor() ([]byte, []int) {
	return file_notify_notify_proto_rawDescGZIP(), []int{2}
}

func (x *Channels) GetMail() *MailNotify {
	if x != nil {
		return x.Mail
	}
	return nil
}

func (x *Channels) GetPhone() *PhoneNotify {
	if x != nil {
		return x.Phone
	}
	return nil
}

type SendMessageReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID     string     `protobuf:"bytes,1,opt,name=userID,proto3" json:"userID,omitempty"`
	NotifyType NotifyType `protobuf:"varint,2,opt,name=notifyType,proto3,enum=notify.NotifyType" json:"notifyType,omitempty"`
	Channels   *Channels  `protobuf:"bytes,3,opt,name=channels,proto3" json:"channels,omitempty"`
}

func (x *SendMessageReq) Reset() {
	*x = SendMessageReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notify_notify_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendMessageReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMessageReq) ProtoMessage() {}

func (x *SendMessageReq) ProtoReflect() protoreflect.Message {
	mi := &file_notify_notify_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMessageReq.ProtoReflect.Descriptor instead.
func (*SendMessageReq) Descriptor() ([]byte, []int) {
	return file_notify_notify_proto_rawDescGZIP(), []int{3}
}

func (x *SendMessageReq) GetUserID() string {
	if x != nil {
		return x.UserID
	}
	return ""
}

func (x *SendMessageReq) GetNotifyType() NotifyType {
	if x != nil {
		return x.NotifyType
	}
	return NotifyType_moderate
}

func (x *SendMessageReq) GetChannels() *Channels {
	if x != nil {
		return x.Channels
	}
	return nil
}

type SendMessageResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Respond string `protobuf:"bytes,1,opt,name=respond,proto3" json:"respond,omitempty"`
}

func (x *SendMessageResp) Reset() {
	*x = SendMessageResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notify_notify_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendMessageResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMessageResp) ProtoMessage() {}

func (x *SendMessageResp) ProtoReflect() protoreflect.Message {
	mi := &file_notify_notify_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMessageResp.ProtoReflect.Descriptor instead.
func (*SendMessageResp) Descriptor() ([]byte, []int) {
	return file_notify_notify_proto_rawDescGZIP(), []int{4}
}

func (x *SendMessageResp) GetRespond() string {
	if x != nil {
		return x.Respond
	}
	return ""
}

type MailApproval struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Approval bool `protobuf:"varint,1,opt,name=approval,proto3" json:"approval,omitempty"`
}

func (x *MailApproval) Reset() {
	*x = MailApproval{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notify_notify_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MailApproval) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MailApproval) ProtoMessage() {}

func (x *MailApproval) ProtoReflect() protoreflect.Message {
	mi := &file_notify_notify_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MailApproval.ProtoReflect.Descriptor instead.
func (*MailApproval) Descriptor() ([]byte, []int) {
	return file_notify_notify_proto_rawDescGZIP(), []int{5}
}

func (x *MailApproval) GetApproval() bool {
	if x != nil {
		return x.Approval
	}
	return false
}

type PhoneApproval struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Approval bool `protobuf:"varint,1,opt,name=approval,proto3" json:"approval,omitempty"`
}

func (x *PhoneApproval) Reset() {
	*x = PhoneApproval{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notify_notify_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PhoneApproval) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PhoneApproval) ProtoMessage() {}

func (x *PhoneApproval) ProtoReflect() protoreflect.Message {
	mi := &file_notify_notify_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PhoneApproval.ProtoReflect.Descriptor instead.
func (*PhoneApproval) Descriptor() ([]byte, []int) {
	return file_notify_notify_proto_rawDescGZIP(), []int{6}
}

func (x *PhoneApproval) GetApproval() bool {
	if x != nil {
		return x.Approval
	}
	return false
}

type Preferences struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Mail  *MailApproval  `protobuf:"bytes,1,opt,name=mail,proto3,oneof" json:"mail,omitempty"`
	Phone *PhoneApproval `protobuf:"bytes,2,opt,name=phone,proto3,oneof" json:"phone,omitempty"`
}

func (x *Preferences) Reset() {
	*x = Preferences{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notify_notify_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Preferences) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Preferences) ProtoMessage() {}

func (x *Preferences) ProtoReflect() protoreflect.Message {
	mi := &file_notify_notify_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Preferences.ProtoReflect.Descriptor instead.
func (*Preferences) Descriptor() ([]byte, []int) {
	return file_notify_notify_proto_rawDescGZIP(), []int{7}
}

func (x *Preferences) GetMail() *MailApproval {
	if x != nil {
		return x.Mail
	}
	return nil
}

func (x *Preferences) GetPhone() *PhoneApproval {
	if x != nil {
		return x.Phone
	}
	return nil
}

type EditPreferencesReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID      string       `protobuf:"bytes,1,opt,name=userID,proto3" json:"userID,omitempty"`
	Preferences *Preferences `protobuf:"bytes,2,opt,name=preferences,proto3" json:"preferences,omitempty"`
}

func (x *EditPreferencesReq) Reset() {
	*x = EditPreferencesReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notify_notify_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EditPreferencesReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EditPreferencesReq) ProtoMessage() {}

func (x *EditPreferencesReq) ProtoReflect() protoreflect.Message {
	mi := &file_notify_notify_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EditPreferencesReq.ProtoReflect.Descriptor instead.
func (*EditPreferencesReq) Descriptor() ([]byte, []int) {
	return file_notify_notify_proto_rawDescGZIP(), []int{8}
}

func (x *EditPreferencesReq) GetUserID() string {
	if x != nil {
		return x.UserID
	}
	return ""
}

func (x *EditPreferencesReq) GetPreferences() *Preferences {
	if x != nil {
		return x.Preferences
	}
	return nil
}

type EditPreferencesResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Respond string `protobuf:"bytes,1,opt,name=respond,proto3" json:"respond,omitempty"`
}

func (x *EditPreferencesResp) Reset() {
	*x = EditPreferencesResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_notify_notify_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EditPreferencesResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EditPreferencesResp) ProtoMessage() {}

func (x *EditPreferencesResp) ProtoReflect() protoreflect.Message {
	mi := &file_notify_notify_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EditPreferencesResp.ProtoReflect.Descriptor instead.
func (*EditPreferencesResp) Descriptor() ([]byte, []int) {
	return file_notify_notify_proto_rawDescGZIP(), []int{9}
}

func (x *EditPreferencesResp) GetRespond() string {
	if x != nil {
		return x.Respond
	}
	return ""
}

var File_notify_notify_proto protoreflect.FileDescriptor

var file_notify_notify_proto_rawDesc = []byte{
	0x0a, 0x13, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x1a, 0x17, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4c, 0x0a, 0x0a, 0x4d, 0x61, 0x69, 0x6c, 0x4e, 0x6f,
	0x74, 0x69, 0x66, 0x79, 0x12, 0x21, 0x0a, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x07,
	0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x1b, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x04,
	0x62, 0x6f, 0x64, 0x79, 0x22, 0x2a, 0x0a, 0x0b, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x4e, 0x6f, 0x74,
	0x69, 0x66, 0x79, 0x12, 0x1b, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79,
	0x22, 0x7a, 0x0a, 0x08, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x12, 0x2b, 0x0a, 0x04,
	0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6e, 0x6f, 0x74,
	0x69, 0x66, 0x79, 0x2e, 0x4d, 0x61, 0x69, 0x6c, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x48, 0x00,
	0x52, 0x04, 0x6d, 0x61, 0x69, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x2e, 0x0a, 0x05, 0x70, 0x68, 0x6f,
	0x6e, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66,
	0x79, 0x2e, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x48, 0x01, 0x52,
	0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x88, 0x01, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x6d, 0x61,
	0x69, 0x6c, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x22, 0xb7, 0x01, 0x0a,
	0x0e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x12,
	0x31, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x19, 0xfa, 0x42, 0x16, 0x72, 0x14, 0x10, 0x01, 0x32, 0x10, 0x5e, 0x5b, 0x61, 0x2d, 0x7a, 0x41,
	0x2d, 0x5a, 0x30, 0x2d, 0x39, 0x5f, 0x2d, 0x5d, 0x2b, 0x24, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x44, 0x12, 0x3a, 0x0a, 0x0a, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x54, 0x79, 0x70, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x2e,
	0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x54, 0x79, 0x70, 0x65, 0x42, 0x06, 0xfa, 0x42, 0x03, 0x82,
	0x01, 0x00, 0x52, 0x0a, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x36,
	0x0a, 0x08, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x10, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65,
	0x6c, 0x73, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x08, 0x63, 0x68,
	0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x22, 0x2b, 0x0a, 0x0f, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x72, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x64, 0x22, 0x2a, 0x0a, 0x0c, 0x4d, 0x61, 0x69, 0x6c, 0x41, 0x70, 0x70, 0x72, 0x6f,
	0x76, 0x61, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x61, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x61, 0x6c, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x61, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x61, 0x6c, 0x22,
	0x2b, 0x0a, 0x0d, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x41, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x61, 0x6c,
	0x12, 0x1a, 0x0a, 0x08, 0x61, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x08, 0x61, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x61, 0x6c, 0x22, 0x81, 0x01, 0x0a,
	0x0b, 0x50, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x12, 0x2d, 0x0a, 0x04,
	0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6e, 0x6f, 0x74,
	0x69, 0x66, 0x79, 0x2e, 0x4d, 0x61, 0x69, 0x6c, 0x41, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x61, 0x6c,
	0x48, 0x00, 0x52, 0x04, 0x6d, 0x61, 0x69, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x30, 0x0a, 0x05, 0x70,
	0x68, 0x6f, 0x6e, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x6e, 0x6f, 0x74,
	0x69, 0x66, 0x79, 0x2e, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x41, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x61,
	0x6c, 0x48, 0x01, 0x52, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x88, 0x01, 0x01, 0x42, 0x07, 0x0a,
	0x05, 0x5f, 0x6d, 0x61, 0x69, 0x6c, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x70, 0x68, 0x6f, 0x6e, 0x65,
	0x22, 0x88, 0x01, 0x0a, 0x12, 0x45, 0x64, 0x69, 0x74, 0x50, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65,
	0x6e, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x12, 0x31, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x19, 0xfa, 0x42, 0x16, 0x72, 0x14, 0x10, 0x01,
	0x32, 0x10, 0x5e, 0x5b, 0x61, 0x2d, 0x7a, 0x41, 0x2d, 0x5a, 0x30, 0x2d, 0x39, 0x5f, 0x2d, 0x5d,
	0x2b, 0x24, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x12, 0x3f, 0x0a, 0x0b, 0x70, 0x72,
	0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x13, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x2e, 0x50, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65,
	0x6e, 0x63, 0x65, 0x73, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x0b,
	0x70, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x22, 0x2f, 0x0a, 0x13, 0x45,
	0x64, 0x69, 0x74, 0x50, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64, 0x2a, 0x36, 0x0a, 0x0a,
	0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0c, 0x0a, 0x08, 0x6d, 0x6f,
	0x64, 0x65, 0x72, 0x61, 0x74, 0x65, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x73, 0x69, 0x67, 0x6e,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x6e, 0x74, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x61, 0x6c, 0x65,
	0x72, 0x74, 0x10, 0x02, 0x32, 0x48, 0x0a, 0x06, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x12, 0x3e,
	0x0a, 0x0b, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x2e,
	0x6e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x17, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x2e, 0x53,
	0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x32, 0x53,
	0x0a, 0x05, 0x55, 0x73, 0x65, 0x72, 0x73, 0x12, 0x4a, 0x0a, 0x0f, 0x45, 0x64, 0x69, 0x74, 0x50,
	0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x12, 0x1a, 0x2e, 0x6e, 0x6f, 0x74,
	0x69, 0x66, 0x79, 0x2e, 0x45, 0x64, 0x69, 0x74, 0x50, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e,
	0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x1a, 0x1b, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x2e,
	0x45, 0x64, 0x69, 0x74, 0x50, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x42, 0x13, 0x5a, 0x11, 0x2e, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x3b,
	0x6e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_notify_notify_proto_rawDescOnce sync.Once
	file_notify_notify_proto_rawDescData = file_notify_notify_proto_rawDesc
)

func file_notify_notify_proto_rawDescGZIP() []byte {
	file_notify_notify_proto_rawDescOnce.Do(func() {
		file_notify_notify_proto_rawDescData = protoimpl.X.CompressGZIP(file_notify_notify_proto_rawDescData)
	})
	return file_notify_notify_proto_rawDescData
}

var file_notify_notify_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_notify_notify_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_notify_notify_proto_goTypes = []interface{}{
	(NotifyType)(0),             // 0: notify.NotifyType
	(*MailNotify)(nil),          // 1: notify.MailNotify
	(*PhoneNotify)(nil),         // 2: notify.PhoneNotify
	(*Channels)(nil),            // 3: notify.Channels
	(*SendMessageReq)(nil),      // 4: notify.SendMessageReq
	(*SendMessageResp)(nil),     // 5: notify.SendMessageResp
	(*MailApproval)(nil),        // 6: notify.MailApproval
	(*PhoneApproval)(nil),       // 7: notify.PhoneApproval
	(*Preferences)(nil),         // 8: notify.Preferences
	(*EditPreferencesReq)(nil),  // 9: notify.EditPreferencesReq
	(*EditPreferencesResp)(nil), // 10: notify.EditPreferencesResp
}
var file_notify_notify_proto_depIdxs = []int32{
	1,  // 0: notify.Channels.mail:type_name -> notify.MailNotify
	2,  // 1: notify.Channels.phone:type_name -> notify.PhoneNotify
	0,  // 2: notify.SendMessageReq.notifyType:type_name -> notify.NotifyType
	3,  // 3: notify.SendMessageReq.channels:type_name -> notify.Channels
	6,  // 4: notify.Preferences.mail:type_name -> notify.MailApproval
	7,  // 5: notify.Preferences.phone:type_name -> notify.PhoneApproval
	8,  // 6: notify.EditPreferencesReq.preferences:type_name -> notify.Preferences
	4,  // 7: notify.Notify.SendMessage:input_type -> notify.SendMessageReq
	9,  // 8: notify.Users.EditPreferences:input_type -> notify.EditPreferencesReq
	5,  // 9: notify.Notify.SendMessage:output_type -> notify.SendMessageResp
	10, // 10: notify.Users.EditPreferences:output_type -> notify.EditPreferencesResp
	9,  // [9:11] is the sub-list for method output_type
	7,  // [7:9] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_notify_notify_proto_init() }
func file_notify_notify_proto_init() {
	if File_notify_notify_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_notify_notify_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MailNotify); i {
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
		file_notify_notify_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PhoneNotify); i {
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
		file_notify_notify_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Channels); i {
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
		file_notify_notify_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendMessageReq); i {
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
		file_notify_notify_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendMessageResp); i {
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
		file_notify_notify_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MailApproval); i {
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
		file_notify_notify_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PhoneApproval); i {
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
		file_notify_notify_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Preferences); i {
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
		file_notify_notify_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EditPreferencesReq); i {
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
		file_notify_notify_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EditPreferencesResp); i {
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
	file_notify_notify_proto_msgTypes[2].OneofWrappers = []interface{}{}
	file_notify_notify_proto_msgTypes[7].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_notify_notify_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_notify_notify_proto_goTypes,
		DependencyIndexes: file_notify_notify_proto_depIdxs,
		EnumInfos:         file_notify_notify_proto_enumTypes,
		MessageInfos:      file_notify_notify_proto_msgTypes,
	}.Build()
	File_notify_notify_proto = out.File
	file_notify_notify_proto_rawDesc = nil
	file_notify_notify_proto_goTypes = nil
	file_notify_notify_proto_depIdxs = nil
}
