// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.21.12
// source: locker.proto

package locker

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
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

// Enumerated lock state.
type StatusResponse_LockState int32

const (
	StatusResponse_UNKNOWN  StatusResponse_LockState = 0
	StatusResponse_LOCKED   StatusResponse_LockState = 1
	StatusResponse_UNLOCKED StatusResponse_LockState = 2
)

// Enum value maps for StatusResponse_LockState.
var (
	StatusResponse_LockState_name = map[int32]string{
		0: "UNKNOWN",
		1: "LOCKED",
		2: "UNLOCKED",
	}
	StatusResponse_LockState_value = map[string]int32{
		"UNKNOWN":  0,
		"LOCKED":   1,
		"UNLOCKED": 2,
	}
)

func (x StatusResponse_LockState) Enum() *StatusResponse_LockState {
	p := new(StatusResponse_LockState)
	*p = x
	return p
}

func (x StatusResponse_LockState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (StatusResponse_LockState) Descriptor() protoreflect.EnumDescriptor {
	return file_locker_proto_enumTypes[0].Descriptor()
}

func (StatusResponse_LockState) Type() protoreflect.EnumType {
	return &file_locker_proto_enumTypes[0]
}

func (x StatusResponse_LockState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use StatusResponse_LockState.Descriptor instead.
func (StatusResponse_LockState) EnumDescriptor() ([]byte, []int) {
	return file_locker_proto_rawDescGZIP(), []int{5, 0}
}

// Request message for the Lock method.
type LockRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The user initiating the lock.
	User string `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	// The user's id
	Uid int32 `protobuf:"varint,2,opt,name=uid,proto3" json:"uid,omitempty"`
	// The tty (terminal) identifier.
	Tty string `protobuf:"bytes,3,opt,name=tty,proto3" json:"tty,omitempty"`
	// The session identifier.
	SessionId string `protobuf:"bytes,4,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	// List of allowed users.
	AllowedUsers []string `protobuf:"bytes,5,rep,name=allowed_users,json=allowedUsers,proto3" json:"allowed_users,omitempty"`
	// List of allowed groups.
	AllowedGroups []string `protobuf:"bytes,6,rep,name=allowed_groups,json=allowedGroups,proto3" json:"allowed_groups,omitempty"`
	// Reason for locking (optional—an empty string if not provided).
	Reason string `protobuf:"bytes,7,opt,name=reason,proto3" json:"reason,omitempty"`
	// Email address of the locking user
	Email string `protobuf:"bytes,8,opt,name=email,proto3" json:"email,omitempty"`
	// Whether to unlock on process exit.
	UnlockOnExit bool `protobuf:"varint,9,opt,name=unlock_on_exit,json=unlockOnExit,proto3" json:"unlock_on_exit,omitempty"`
	// Optional duration after which the lock should automatically unlock (e.g. "10m").
	UnlockTime *durationpb.Duration `protobuf:"bytes,10,opt,name=unlock_time,json=unlockTime,proto3" json:"unlock_time,omitempty"`
	// Optional duration after which the lock is released if idle.
	UnlockAfterIdle *durationpb.Duration `protobuf:"bytes,11,opt,name=unlock_after_idle,json=unlockAfterIdle,proto3" json:"unlock_after_idle,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *LockRequest) Reset() {
	*x = LockRequest{}
	mi := &file_locker_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LockRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LockRequest) ProtoMessage() {}

func (x *LockRequest) ProtoReflect() protoreflect.Message {
	mi := &file_locker_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LockRequest.ProtoReflect.Descriptor instead.
func (*LockRequest) Descriptor() ([]byte, []int) {
	return file_locker_proto_rawDescGZIP(), []int{0}
}

func (x *LockRequest) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *LockRequest) GetUid() int32 {
	if x != nil {
		return x.Uid
	}
	return 0
}

func (x *LockRequest) GetTty() string {
	if x != nil {
		return x.Tty
	}
	return ""
}

func (x *LockRequest) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

func (x *LockRequest) GetAllowedUsers() []string {
	if x != nil {
		return x.AllowedUsers
	}
	return nil
}

func (x *LockRequest) GetAllowedGroups() []string {
	if x != nil {
		return x.AllowedGroups
	}
	return nil
}

func (x *LockRequest) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

func (x *LockRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *LockRequest) GetUnlockOnExit() bool {
	if x != nil {
		return x.UnlockOnExit
	}
	return false
}

func (x *LockRequest) GetUnlockTime() *durationpb.Duration {
	if x != nil {
		return x.UnlockTime
	}
	return nil
}

func (x *LockRequest) GetUnlockAfterIdle() *durationpb.Duration {
	if x != nil {
		return x.UnlockAfterIdle
	}
	return nil
}

// Response message for the Lock method.
type LockResponse struct {
	state   protoimpl.MessageState `protogen:"open.v1"`
	Success bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	// An optional message (e.g. error description or confirmation).
	Message       string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LockResponse) Reset() {
	*x = LockResponse{}
	mi := &file_locker_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LockResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LockResponse) ProtoMessage() {}

func (x *LockResponse) ProtoReflect() protoreflect.Message {
	mi := &file_locker_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LockResponse.ProtoReflect.Descriptor instead.
func (*LockResponse) Descriptor() ([]byte, []int) {
	return file_locker_proto_rawDescGZIP(), []int{1}
}

func (x *LockResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *LockResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// Request message for the Unlock method (no parameters).
type UnlockRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UnlockRequest) Reset() {
	*x = UnlockRequest{}
	mi := &file_locker_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UnlockRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnlockRequest) ProtoMessage() {}

func (x *UnlockRequest) ProtoReflect() protoreflect.Message {
	mi := &file_locker_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnlockRequest.ProtoReflect.Descriptor instead.
func (*UnlockRequest) Descriptor() ([]byte, []int) {
	return file_locker_proto_rawDescGZIP(), []int{2}
}

// Response message for the Unlock method.
type UnlockResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UnlockResponse) Reset() {
	*x = UnlockResponse{}
	mi := &file_locker_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UnlockResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnlockResponse) ProtoMessage() {}

func (x *UnlockResponse) ProtoReflect() protoreflect.Message {
	mi := &file_locker_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnlockResponse.ProtoReflect.Descriptor instead.
func (*UnlockResponse) Descriptor() ([]byte, []int) {
	return file_locker_proto_rawDescGZIP(), []int{3}
}

func (x *UnlockResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *UnlockResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// Request message for the Status method (no parameters).
type StatusRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StatusRequest) Reset() {
	*x = StatusRequest{}
	mi := &file_locker_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusRequest) ProtoMessage() {}

func (x *StatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_locker_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusRequest.ProtoReflect.Descriptor instead.
func (*StatusRequest) Descriptor() ([]byte, []int) {
	return file_locker_proto_rawDescGZIP(), []int{4}
}

// Response message for the Status method.
type StatusResponse struct {
	state protoimpl.MessageState   `protogen:"open.v1"`
	State StatusResponse_LockState `protobuf:"varint,1,opt,name=state,proto3,enum=locker.StatusResponse_LockState" json:"state,omitempty"`
	// If the system is locked, include the details.
	User            string               `protobuf:"bytes,2,opt,name=user,proto3" json:"user,omitempty"`
	Uid             int32                `protobuf:"varint,3,opt,name=uid,proto3" json:"uid,omitempty"`
	Tty             string               `protobuf:"bytes,4,opt,name=tty,proto3" json:"tty,omitempty"`
	SessionId       string               `protobuf:"bytes,5,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	AllowedUsers    []string             `protobuf:"bytes,6,rep,name=allowed_users,json=allowedUsers,proto3" json:"allowed_users,omitempty"`
	AllowedGroups   []string             `protobuf:"bytes,7,rep,name=allowed_groups,json=allowedGroups,proto3" json:"allowed_groups,omitempty"`
	Reason          string               `protobuf:"bytes,8,opt,name=reason,proto3" json:"reason,omitempty"`
	Email           string               `protobuf:"bytes,9,opt,name=email,proto3" json:"email,omitempty"`
	UnlockOnExit    bool                 `protobuf:"varint,10,opt,name=unlock_on_exit,json=unlockOnExit,proto3" json:"unlock_on_exit,omitempty"`
	UnlockTime      *durationpb.Duration `protobuf:"bytes,11,opt,name=unlock_time,json=unlockTime,proto3" json:"unlock_time,omitempty"`
	UnlockAfterIdle *durationpb.Duration `protobuf:"bytes,12,opt,name=unlock_after_idle,json=unlockAfterIdle,proto3" json:"unlock_after_idle,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *StatusResponse) Reset() {
	*x = StatusResponse{}
	mi := &file_locker_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusResponse) ProtoMessage() {}

func (x *StatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_locker_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusResponse.ProtoReflect.Descriptor instead.
func (*StatusResponse) Descriptor() ([]byte, []int) {
	return file_locker_proto_rawDescGZIP(), []int{5}
}

func (x *StatusResponse) GetState() StatusResponse_LockState {
	if x != nil {
		return x.State
	}
	return StatusResponse_UNKNOWN
}

func (x *StatusResponse) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *StatusResponse) GetUid() int32 {
	if x != nil {
		return x.Uid
	}
	return 0
}

func (x *StatusResponse) GetTty() string {
	if x != nil {
		return x.Tty
	}
	return ""
}

func (x *StatusResponse) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

func (x *StatusResponse) GetAllowedUsers() []string {
	if x != nil {
		return x.AllowedUsers
	}
	return nil
}

func (x *StatusResponse) GetAllowedGroups() []string {
	if x != nil {
		return x.AllowedGroups
	}
	return nil
}

func (x *StatusResponse) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

func (x *StatusResponse) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *StatusResponse) GetUnlockOnExit() bool {
	if x != nil {
		return x.UnlockOnExit
	}
	return false
}

func (x *StatusResponse) GetUnlockTime() *durationpb.Duration {
	if x != nil {
		return x.UnlockTime
	}
	return nil
}

func (x *StatusResponse) GetUnlockAfterIdle() *durationpb.Duration {
	if x != nil {
		return x.UnlockAfterIdle
	}
	return nil
}

// Request message for the Authorize method.
type AuthorizeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Username      string                 `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AuthorizeRequest) Reset() {
	*x = AuthorizeRequest{}
	mi := &file_locker_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AuthorizeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthorizeRequest) ProtoMessage() {}

func (x *AuthorizeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_locker_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthorizeRequest.ProtoReflect.Descriptor instead.
func (*AuthorizeRequest) Descriptor() ([]byte, []int) {
	return file_locker_proto_rawDescGZIP(), []int{6}
}

func (x *AuthorizeRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

// Response message for the Authorize method.
type AuthorizeResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// True if the given username is authorized.
	Authorized bool `protobuf:"varint,1,opt,name=authorized,proto3" json:"authorized,omitempty"`
	// Optional explanation (e.g. reason for denial).
	Reason        string `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AuthorizeResponse) Reset() {
	*x = AuthorizeResponse{}
	mi := &file_locker_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AuthorizeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthorizeResponse) ProtoMessage() {}

func (x *AuthorizeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_locker_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthorizeResponse.ProtoReflect.Descriptor instead.
func (*AuthorizeResponse) Descriptor() ([]byte, []int) {
	return file_locker_proto_rawDescGZIP(), []int{7}
}

func (x *AuthorizeResponse) GetAuthorized() bool {
	if x != nil {
		return x.Authorized
	}
	return false
}

func (x *AuthorizeResponse) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

// Request email address associated with the user.
type EmailRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Username      string                 `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EmailRequest) Reset() {
	*x = EmailRequest{}
	mi := &file_locker_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EmailRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmailRequest) ProtoMessage() {}

func (x *EmailRequest) ProtoReflect() protoreflect.Message {
	mi := &file_locker_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmailRequest.ProtoReflect.Descriptor instead.
func (*EmailRequest) Descriptor() ([]byte, []int) {
	return file_locker_proto_rawDescGZIP(), []int{8}
}

func (x *EmailRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

// Response email address associated with the user.
type EmailResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Found         bool                   `protobuf:"varint,1,opt,name=found,proto3" json:"found,omitempty"`
	Email         string                 `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EmailResponse) Reset() {
	*x = EmailResponse{}
	mi := &file_locker_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EmailResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmailResponse) ProtoMessage() {}

func (x *EmailResponse) ProtoReflect() protoreflect.Message {
	mi := &file_locker_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmailResponse.ProtoReflect.Descriptor instead.
func (*EmailResponse) Descriptor() ([]byte, []int) {
	return file_locker_proto_rawDescGZIP(), []int{9}
}

func (x *EmailResponse) GetFound() bool {
	if x != nil {
		return x.Found
	}
	return false
}

func (x *EmailResponse) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

var File_locker_proto protoreflect.FileDescriptor

var file_locker_proto_rawDesc = string([]byte{
	0x0a, 0x0c, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x87, 0x03, 0x0a, 0x0b, 0x4c, 0x6f, 0x63, 0x6b, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03,
	0x74, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x74, 0x74, 0x79, 0x12, 0x1d,
	0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x23, 0x0a,
	0x0d, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x05,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x55, 0x73, 0x65,
	0x72, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x5f, 0x67, 0x72,
	0x6f, 0x75, 0x70, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0d, 0x61, 0x6c, 0x6c, 0x6f,
	0x77, 0x65, 0x64, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x61,
	0x73, 0x6f, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f,
	0x6e, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x24, 0x0a, 0x0e, 0x75, 0x6e, 0x6c, 0x6f, 0x63,
	0x6b, 0x5f, 0x6f, 0x6e, 0x5f, 0x65, 0x78, 0x69, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0c, 0x75, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x4f, 0x6e, 0x45, 0x78, 0x69, 0x74, 0x12, 0x3a, 0x0a,
	0x0b, 0x75, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x75,
	0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x45, 0x0a, 0x11, 0x75, 0x6e, 0x6c,
	0x6f, 0x63, 0x6b, 0x5f, 0x61, 0x66, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x6c, 0x65, 0x18, 0x0b,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x0f, 0x75, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x41, 0x66, 0x74, 0x65, 0x72, 0x49, 0x64, 0x6c, 0x65,
	0x22, 0x42, 0x0a, 0x0c, 0x4c, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0x0f, 0x0a, 0x0d, 0x55, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x44, 0x0a, 0x0e, 0x55, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x0f, 0x0a, 0x0d, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0xf6, 0x03, 0x0a,
	0x0e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x36, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x20,
	0x2e, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x4c, 0x6f, 0x63, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x75,
	0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x10, 0x0a,
	0x03, 0x74, 0x74, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x74, 0x74, 0x79, 0x12,
	0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x23,
	0x0a, 0x0d, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18,
	0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x55, 0x73,
	0x65, 0x72, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x5f, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0d, 0x61, 0x6c, 0x6c,
	0x6f, 0x77, 0x65, 0x64, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65,
	0x61, 0x73, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73,
	0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x24, 0x0a, 0x0e, 0x75, 0x6e, 0x6c, 0x6f,
	0x63, 0x6b, 0x5f, 0x6f, 0x6e, 0x5f, 0x65, 0x78, 0x69, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x0c, 0x75, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x4f, 0x6e, 0x45, 0x78, 0x69, 0x74, 0x12, 0x3a,
	0x0a, 0x0b, 0x75, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a,
	0x75, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x45, 0x0a, 0x11, 0x75, 0x6e,
	0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x61, 0x66, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x6c, 0x65, 0x18,
	0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x0f, 0x75, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x41, 0x66, 0x74, 0x65, 0x72, 0x49, 0x64, 0x6c,
	0x65, 0x22, 0x32, 0x0a, 0x09, 0x4c, 0x6f, 0x63, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x0b,
	0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x4c,
	0x4f, 0x43, 0x4b, 0x45, 0x44, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x55, 0x4e, 0x4c, 0x4f, 0x43,
	0x4b, 0x45, 0x44, 0x10, 0x02, 0x22, 0x2e, 0x0a, 0x10, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69,
	0x7a, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x4b, 0x0a, 0x11, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69,
	0x7a, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x61, 0x75,
	0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a,
	0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65,
	0x61, 0x73, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73,
	0x6f, 0x6e, 0x22, 0x2a, 0x0a, 0x0c, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x3b,
	0x0a, 0x0d, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x66, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05,
	0x66, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x32, 0xac, 0x02, 0x0a, 0x0d,
	0x4c, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x31, 0x0a,
	0x04, 0x4c, 0x6f, 0x63, 0x6b, 0x12, 0x13, 0x2e, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x4c,
	0x6f, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x6c, 0x6f, 0x63,
	0x6b, 0x65, 0x72, 0x2e, 0x4c, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x37, 0x0a, 0x06, 0x55, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x15, 0x2e, 0x6c, 0x6f, 0x63,
	0x6b, 0x65, 0x72, 0x2e, 0x55, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x16, 0x2e, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x55, 0x6e, 0x6c, 0x6f, 0x63,
	0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x06, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x15, 0x2e, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x6c, 0x6f, 0x63,
	0x6b, 0x65, 0x72, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x40, 0x0a, 0x09, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x12,
	0x18, 0x2e, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69,
	0x7a, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x6c, 0x6f, 0x63, 0x6b,
	0x65, 0x72, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x34, 0x0a, 0x05, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x14, 0x2e,
	0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x45, 0x6d, 0x61,
	0x69, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x27, 0x5a, 0x25, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x67, 0x72, 0x65, 0x77, 0x65, 0x6c,
	0x6c, 0x2f, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x3b, 0x6c, 0x6f, 0x63,
	0x6b, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_locker_proto_rawDescOnce sync.Once
	file_locker_proto_rawDescData []byte
)

func file_locker_proto_rawDescGZIP() []byte {
	file_locker_proto_rawDescOnce.Do(func() {
		file_locker_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_locker_proto_rawDesc), len(file_locker_proto_rawDesc)))
	})
	return file_locker_proto_rawDescData
}

var file_locker_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_locker_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_locker_proto_goTypes = []any{
	(StatusResponse_LockState)(0), // 0: locker.StatusResponse.LockState
	(*LockRequest)(nil),           // 1: locker.LockRequest
	(*LockResponse)(nil),          // 2: locker.LockResponse
	(*UnlockRequest)(nil),         // 3: locker.UnlockRequest
	(*UnlockResponse)(nil),        // 4: locker.UnlockResponse
	(*StatusRequest)(nil),         // 5: locker.StatusRequest
	(*StatusResponse)(nil),        // 6: locker.StatusResponse
	(*AuthorizeRequest)(nil),      // 7: locker.AuthorizeRequest
	(*AuthorizeResponse)(nil),     // 8: locker.AuthorizeResponse
	(*EmailRequest)(nil),          // 9: locker.EmailRequest
	(*EmailResponse)(nil),         // 10: locker.EmailResponse
	(*durationpb.Duration)(nil),   // 11: google.protobuf.Duration
}
var file_locker_proto_depIdxs = []int32{
	11, // 0: locker.LockRequest.unlock_time:type_name -> google.protobuf.Duration
	11, // 1: locker.LockRequest.unlock_after_idle:type_name -> google.protobuf.Duration
	0,  // 2: locker.StatusResponse.state:type_name -> locker.StatusResponse.LockState
	11, // 3: locker.StatusResponse.unlock_time:type_name -> google.protobuf.Duration
	11, // 4: locker.StatusResponse.unlock_after_idle:type_name -> google.protobuf.Duration
	1,  // 5: locker.LockerService.Lock:input_type -> locker.LockRequest
	3,  // 6: locker.LockerService.Unlock:input_type -> locker.UnlockRequest
	5,  // 7: locker.LockerService.Status:input_type -> locker.StatusRequest
	7,  // 8: locker.LockerService.Authorize:input_type -> locker.AuthorizeRequest
	9,  // 9: locker.LockerService.Email:input_type -> locker.EmailRequest
	2,  // 10: locker.LockerService.Lock:output_type -> locker.LockResponse
	4,  // 11: locker.LockerService.Unlock:output_type -> locker.UnlockResponse
	6,  // 12: locker.LockerService.Status:output_type -> locker.StatusResponse
	8,  // 13: locker.LockerService.Authorize:output_type -> locker.AuthorizeResponse
	10, // 14: locker.LockerService.Email:output_type -> locker.EmailResponse
	10, // [10:15] is the sub-list for method output_type
	5,  // [5:10] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_locker_proto_init() }
func file_locker_proto_init() {
	if File_locker_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_locker_proto_rawDesc), len(file_locker_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_locker_proto_goTypes,
		DependencyIndexes: file_locker_proto_depIdxs,
		EnumInfos:         file_locker_proto_enumTypes,
		MessageInfos:      file_locker_proto_msgTypes,
	}.Build()
	File_locker_proto = out.File
	file_locker_proto_goTypes = nil
	file_locker_proto_depIdxs = nil
}
