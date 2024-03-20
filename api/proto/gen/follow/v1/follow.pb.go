// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: follow/v1/follow.proto

package followv1

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

type FollowRelation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Follower int64 `protobuf:"varint,2,opt,name=follower,proto3" json:"follower,omitempty"`
	Followee int64 `protobuf:"varint,3,opt,name=followee,proto3" json:"followee,omitempty"`
}

func (x *FollowRelation) Reset() {
	*x = FollowRelation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_follow_v1_follow_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FollowRelation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FollowRelation) ProtoMessage() {}

func (x *FollowRelation) ProtoReflect() protoreflect.Message {
	mi := &file_follow_v1_follow_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FollowRelation.ProtoReflect.Descriptor instead.
func (*FollowRelation) Descriptor() ([]byte, []int) {
	return file_follow_v1_follow_proto_rawDescGZIP(), []int{0}
}

func (x *FollowRelation) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *FollowRelation) GetFollower() int64 {
	if x != nil {
		return x.Follower
	}
	return 0
}

func (x *FollowRelation) GetFollowee() int64 {
	if x != nil {
		return x.Followee
	}
	return 0
}

type FollowRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Followee int64 `protobuf:"varint,1,opt,name=followee,proto3" json:"followee,omitempty"` // 被关注者
	Follower int64 `protobuf:"varint,2,opt,name=follower,proto3" json:"follower,omitempty"` // 关注者
}

func (x *FollowRequest) Reset() {
	*x = FollowRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_follow_v1_follow_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FollowRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FollowRequest) ProtoMessage() {}

func (x *FollowRequest) ProtoReflect() protoreflect.Message {
	mi := &file_follow_v1_follow_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FollowRequest.ProtoReflect.Descriptor instead.
func (*FollowRequest) Descriptor() ([]byte, []int) {
	return file_follow_v1_follow_proto_rawDescGZIP(), []int{1}
}

func (x *FollowRequest) GetFollowee() int64 {
	if x != nil {
		return x.Followee
	}
	return 0
}

func (x *FollowRequest) GetFollower() int64 {
	if x != nil {
		return x.Follower
	}
	return 0
}

type FollowResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *FollowResponse) Reset() {
	*x = FollowResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_follow_v1_follow_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FollowResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FollowResponse) ProtoMessage() {}

func (x *FollowResponse) ProtoReflect() protoreflect.Message {
	mi := &file_follow_v1_follow_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FollowResponse.ProtoReflect.Descriptor instead.
func (*FollowResponse) Descriptor() ([]byte, []int) {
	return file_follow_v1_follow_proto_rawDescGZIP(), []int{2}
}

type CancelFollowRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Followee int64 `protobuf:"varint,1,opt,name=followee,proto3" json:"followee,omitempty"` // 被关注者
	Follower int64 `protobuf:"varint,2,opt,name=follower,proto3" json:"follower,omitempty"` // 关注者
}

func (x *CancelFollowRequest) Reset() {
	*x = CancelFollowRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_follow_v1_follow_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CancelFollowRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CancelFollowRequest) ProtoMessage() {}

func (x *CancelFollowRequest) ProtoReflect() protoreflect.Message {
	mi := &file_follow_v1_follow_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CancelFollowRequest.ProtoReflect.Descriptor instead.
func (*CancelFollowRequest) Descriptor() ([]byte, []int) {
	return file_follow_v1_follow_proto_rawDescGZIP(), []int{3}
}

func (x *CancelFollowRequest) GetFollowee() int64 {
	if x != nil {
		return x.Followee
	}
	return 0
}

func (x *CancelFollowRequest) GetFollower() int64 {
	if x != nil {
		return x.Follower
	}
	return 0
}

type CancelFollowResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CancelFollowResponse) Reset() {
	*x = CancelFollowResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_follow_v1_follow_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CancelFollowResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CancelFollowResponse) ProtoMessage() {}

func (x *CancelFollowResponse) ProtoReflect() protoreflect.Message {
	mi := &file_follow_v1_follow_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CancelFollowResponse.ProtoReflect.Descriptor instead.
func (*CancelFollowResponse) Descriptor() ([]byte, []int) {
	return file_follow_v1_follow_proto_rawDescGZIP(), []int{4}
}

type GetFolloweeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Follower int64 `protobuf:"varint,1,opt,name=follower,proto3" json:"follower,omitempty"`
	Offset   int64 `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit    int64 `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *GetFolloweeRequest) Reset() {
	*x = GetFolloweeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_follow_v1_follow_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFolloweeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFolloweeRequest) ProtoMessage() {}

func (x *GetFolloweeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_follow_v1_follow_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFolloweeRequest.ProtoReflect.Descriptor instead.
func (*GetFolloweeRequest) Descriptor() ([]byte, []int) {
	return file_follow_v1_follow_proto_rawDescGZIP(), []int{5}
}

func (x *GetFolloweeRequest) GetFollower() int64 {
	if x != nil {
		return x.Follower
	}
	return 0
}

func (x *GetFolloweeRequest) GetOffset() int64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *GetFolloweeRequest) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type GetFolloweeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FollowRelations []*FollowRelation `protobuf:"bytes,1,rep,name=follow_relations,json=followRelations,proto3" json:"follow_relations,omitempty"`
}

func (x *GetFolloweeResponse) Reset() {
	*x = GetFolloweeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_follow_v1_follow_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFolloweeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFolloweeResponse) ProtoMessage() {}

func (x *GetFolloweeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_follow_v1_follow_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFolloweeResponse.ProtoReflect.Descriptor instead.
func (*GetFolloweeResponse) Descriptor() ([]byte, []int) {
	return file_follow_v1_follow_proto_rawDescGZIP(), []int{6}
}

func (x *GetFolloweeResponse) GetFollowRelations() []*FollowRelation {
	if x != nil {
		return x.FollowRelations
	}
	return nil
}

type GetFollowerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Followee int64 `protobuf:"varint,1,opt,name=followee,proto3" json:"followee,omitempty"`
	Offset   int64 `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit    int64 `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *GetFollowerRequest) Reset() {
	*x = GetFollowerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_follow_v1_follow_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFollowerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFollowerRequest) ProtoMessage() {}

func (x *GetFollowerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_follow_v1_follow_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFollowerRequest.ProtoReflect.Descriptor instead.
func (*GetFollowerRequest) Descriptor() ([]byte, []int) {
	return file_follow_v1_follow_proto_rawDescGZIP(), []int{7}
}

func (x *GetFollowerRequest) GetFollowee() int64 {
	if x != nil {
		return x.Followee
	}
	return 0
}

func (x *GetFollowerRequest) GetOffset() int64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *GetFollowerRequest) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type GetFollowerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FollowRelations []*FollowRelation `protobuf:"bytes,1,rep,name=follow_relations,json=followRelations,proto3" json:"follow_relations,omitempty"`
}

func (x *GetFollowerResponse) Reset() {
	*x = GetFollowerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_follow_v1_follow_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFollowerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFollowerResponse) ProtoMessage() {}

func (x *GetFollowerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_follow_v1_follow_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFollowerResponse.ProtoReflect.Descriptor instead.
func (*GetFollowerResponse) Descriptor() ([]byte, []int) {
	return file_follow_v1_follow_proto_rawDescGZIP(), []int{8}
}

func (x *GetFollowerResponse) GetFollowRelations() []*FollowRelation {
	if x != nil {
		return x.FollowRelations
	}
	return nil
}

type FollowInfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Follower int64 `protobuf:"varint,1,opt,name=follower,proto3" json:"follower,omitempty"`
	Followee int64 `protobuf:"varint,2,opt,name=followee,proto3" json:"followee,omitempty"`
}

func (x *FollowInfoRequest) Reset() {
	*x = FollowInfoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_follow_v1_follow_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FollowInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FollowInfoRequest) ProtoMessage() {}

func (x *FollowInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_follow_v1_follow_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FollowInfoRequest.ProtoReflect.Descriptor instead.
func (*FollowInfoRequest) Descriptor() ([]byte, []int) {
	return file_follow_v1_follow_proto_rawDescGZIP(), []int{9}
}

func (x *FollowInfoRequest) GetFollower() int64 {
	if x != nil {
		return x.Follower
	}
	return 0
}

func (x *FollowInfoRequest) GetFollowee() int64 {
	if x != nil {
		return x.Followee
	}
	return 0
}

type FollowInfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FollowRelations *FollowRelation `protobuf:"bytes,1,opt,name=follow_relations,json=followRelations,proto3" json:"follow_relations,omitempty"`
}

func (x *FollowInfoResponse) Reset() {
	*x = FollowInfoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_follow_v1_follow_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FollowInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FollowInfoResponse) ProtoMessage() {}

func (x *FollowInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_follow_v1_follow_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FollowInfoResponse.ProtoReflect.Descriptor instead.
func (*FollowInfoResponse) Descriptor() ([]byte, []int) {
	return file_follow_v1_follow_proto_rawDescGZIP(), []int{10}
}

func (x *FollowInfoResponse) GetFollowRelations() *FollowRelation {
	if x != nil {
		return x.FollowRelations
	}
	return nil
}

type GetFollowStaticsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid int64 `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
}

func (x *GetFollowStaticsRequest) Reset() {
	*x = GetFollowStaticsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_follow_v1_follow_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFollowStaticsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFollowStaticsRequest) ProtoMessage() {}

func (x *GetFollowStaticsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_follow_v1_follow_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFollowStaticsRequest.ProtoReflect.Descriptor instead.
func (*GetFollowStaticsRequest) Descriptor() ([]byte, []int) {
	return file_follow_v1_follow_proto_rawDescGZIP(), []int{11}
}

func (x *GetFollowStaticsRequest) GetUid() int64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

type GetFollowStaticsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Followees int64 `protobuf:"varint,1,opt,name=followees,proto3" json:"followees,omitempty"` // 关注数
	Followers int64 `protobuf:"varint,2,opt,name=followers,proto3" json:"followers,omitempty"` // 粉丝数
}

func (x *GetFollowStaticsResponse) Reset() {
	*x = GetFollowStaticsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_follow_v1_follow_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFollowStaticsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFollowStaticsResponse) ProtoMessage() {}

func (x *GetFollowStaticsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_follow_v1_follow_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFollowStaticsResponse.ProtoReflect.Descriptor instead.
func (*GetFollowStaticsResponse) Descriptor() ([]byte, []int) {
	return file_follow_v1_follow_proto_rawDescGZIP(), []int{12}
}

func (x *GetFollowStaticsResponse) GetFollowees() int64 {
	if x != nil {
		return x.Followees
	}
	return 0
}

func (x *GetFollowStaticsResponse) GetFollowers() int64 {
	if x != nil {
		return x.Followers
	}
	return 0
}

var File_follow_v1_follow_proto protoreflect.FileDescriptor

var file_follow_v1_follow_proto_rawDesc = []byte{
	0x0a, 0x16, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x2f, 0x76, 0x31, 0x2f, 0x66, 0x6f, 0x6c, 0x6c,
	0x6f, 0x77, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77,
	0x2e, 0x76, 0x31, 0x22, 0x58, 0x0a, 0x0e, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x6c,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65,
	0x72, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x08, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x22, 0x47, 0x0a,
	0x0d, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a,
	0x0a, 0x08, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x08, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x6f,
	0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x66, 0x6f,
	0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x22, 0x10, 0x0a, 0x0e, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x4d, 0x0a, 0x13, 0x43, 0x61, 0x6e, 0x63,
	0x65, 0x6c, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1a, 0x0a, 0x08, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x08, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x66,
	0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x66,
	0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x22, 0x16, 0x0a, 0x14, 0x43, 0x61, 0x6e, 0x63, 0x65,
	0x6c, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x5e, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65,
	0x72, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d,
	0x69, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x22,
	0x5b, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x44, 0x0a, 0x10, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77,
	0x5f, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x6f, 0x6c,
	0x6c, 0x6f, 0x77, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0f, 0x66, 0x6f, 0x6c,
	0x6c, 0x6f, 0x77, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x5e, 0x0a, 0x12,
	0x47, 0x65, 0x74, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x22, 0x5b, 0x0a, 0x13,
	0x47, 0x65, 0x74, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x44, 0x0a, 0x10, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x5f, 0x72, 0x65,
	0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e,
	0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77,
	0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0f, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77,
	0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x4b, 0x0a, 0x11, 0x46, 0x6f, 0x6c,
	0x6c, 0x6f, 0x77, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a,
	0x0a, 0x08, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x08, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x6f,
	0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x66, 0x6f,
	0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x22, 0x5a, 0x0a, 0x12, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x44, 0x0a, 0x10,
	0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x5f, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x2e,
	0x76, 0x31, 0x2e, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x0f, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x22, 0x2b, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x53,
	0x74, 0x61, 0x74, 0x69, 0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a,
	0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x75, 0x69, 0x64, 0x22,
	0x56, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x53, 0x74, 0x61, 0x74,
	0x69, 0x63, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x66,
	0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x66, 0x6f, 0x6c,
	0x6c, 0x6f, 0x77, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x66, 0x6f,
	0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x73, 0x32, 0xe3, 0x03, 0x0a, 0x0d, 0x46, 0x6f, 0x6c, 0x6c,
	0x6f, 0x77, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3d, 0x0a, 0x06, 0x46, 0x6f, 0x6c,
	0x6c, 0x6f, 0x77, 0x12, 0x18, 0x2e, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x2e, 0x76, 0x31, 0x2e,
	0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e,
	0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4f, 0x0a, 0x0c, 0x43, 0x61, 0x6e, 0x63,
	0x65, 0x6c, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x12, 0x1e, 0x2e, 0x66, 0x6f, 0x6c, 0x6c, 0x6f,
	0x77, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x46, 0x6f, 0x6c, 0x6c, 0x6f,
	0x77, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x66, 0x6f, 0x6c, 0x6c, 0x6f,
	0x77, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x46, 0x6f, 0x6c, 0x6c, 0x6f,
	0x77, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4c, 0x0a, 0x0b, 0x47, 0x65, 0x74,
	0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x12, 0x1d, 0x2e, 0x66, 0x6f, 0x6c, 0x6c, 0x6f,
	0x77, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77,
	0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4c, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x46, 0x6f,
	0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x12, 0x1d, 0x2e, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x2e,
	0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x2e, 0x76,
	0x31, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x49, 0x0a, 0x0a, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x49,
	0x6e, 0x66, 0x6f, 0x12, 0x1c, 0x2e, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x2e, 0x76, 0x31, 0x2e,
	0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1d, 0x2e, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x6f,
	0x6c, 0x6c, 0x6f, 0x77, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x5b, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x53, 0x74, 0x61,
	0x74, 0x69, 0x63, 0x73, 0x12, 0x22, 0x2e, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x2e, 0x76, 0x31,
	0x2e, 0x47, 0x65, 0x74, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x66, 0x6f, 0x6c, 0x6c, 0x6f,
	0x77, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x53, 0x74,
	0x61, 0x74, 0x69, 0x63, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x9b, 0x01,
	0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x2e, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x2e, 0x76, 0x31, 0x42,
	0x0b, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x38,
	0x67, 0x69, 0x74, 0x65, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4e, 0x6d, 0x6f, 0x65, 0x6b, 0x2f,
	0x6b, 0x69, 0x74, 0x62, 0x6f, 0x6f, 0x6b, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x2f, 0x76, 0x31, 0x3b,
	0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x46, 0x58, 0x58, 0xaa, 0x02,
	0x09, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x09, 0x46, 0x6f, 0x6c,
	0x6c, 0x6f, 0x77, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x15, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x5c,
	0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02,
	0x0a, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_follow_v1_follow_proto_rawDescOnce sync.Once
	file_follow_v1_follow_proto_rawDescData = file_follow_v1_follow_proto_rawDesc
)

func file_follow_v1_follow_proto_rawDescGZIP() []byte {
	file_follow_v1_follow_proto_rawDescOnce.Do(func() {
		file_follow_v1_follow_proto_rawDescData = protoimpl.X.CompressGZIP(file_follow_v1_follow_proto_rawDescData)
	})
	return file_follow_v1_follow_proto_rawDescData
}

var file_follow_v1_follow_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_follow_v1_follow_proto_goTypes = []interface{}{
	(*FollowRelation)(nil),           // 0: follow.v1.FollowRelation
	(*FollowRequest)(nil),            // 1: follow.v1.FollowRequest
	(*FollowResponse)(nil),           // 2: follow.v1.FollowResponse
	(*CancelFollowRequest)(nil),      // 3: follow.v1.CancelFollowRequest
	(*CancelFollowResponse)(nil),     // 4: follow.v1.CancelFollowResponse
	(*GetFolloweeRequest)(nil),       // 5: follow.v1.GetFolloweeRequest
	(*GetFolloweeResponse)(nil),      // 6: follow.v1.GetFolloweeResponse
	(*GetFollowerRequest)(nil),       // 7: follow.v1.GetFollowerRequest
	(*GetFollowerResponse)(nil),      // 8: follow.v1.GetFollowerResponse
	(*FollowInfoRequest)(nil),        // 9: follow.v1.FollowInfoRequest
	(*FollowInfoResponse)(nil),       // 10: follow.v1.FollowInfoResponse
	(*GetFollowStaticsRequest)(nil),  // 11: follow.v1.GetFollowStaticsRequest
	(*GetFollowStaticsResponse)(nil), // 12: follow.v1.GetFollowStaticsResponse
}
var file_follow_v1_follow_proto_depIdxs = []int32{
	0,  // 0: follow.v1.GetFolloweeResponse.follow_relations:type_name -> follow.v1.FollowRelation
	0,  // 1: follow.v1.GetFollowerResponse.follow_relations:type_name -> follow.v1.FollowRelation
	0,  // 2: follow.v1.FollowInfoResponse.follow_relations:type_name -> follow.v1.FollowRelation
	1,  // 3: follow.v1.FollowService.Follow:input_type -> follow.v1.FollowRequest
	3,  // 4: follow.v1.FollowService.CancelFollow:input_type -> follow.v1.CancelFollowRequest
	5,  // 5: follow.v1.FollowService.GetFollowee:input_type -> follow.v1.GetFolloweeRequest
	7,  // 6: follow.v1.FollowService.GetFollower:input_type -> follow.v1.GetFollowerRequest
	9,  // 7: follow.v1.FollowService.FollowInfo:input_type -> follow.v1.FollowInfoRequest
	11, // 8: follow.v1.FollowService.GetFollowStatics:input_type -> follow.v1.GetFollowStaticsRequest
	2,  // 9: follow.v1.FollowService.Follow:output_type -> follow.v1.FollowResponse
	4,  // 10: follow.v1.FollowService.CancelFollow:output_type -> follow.v1.CancelFollowResponse
	6,  // 11: follow.v1.FollowService.GetFollowee:output_type -> follow.v1.GetFolloweeResponse
	8,  // 12: follow.v1.FollowService.GetFollower:output_type -> follow.v1.GetFollowerResponse
	10, // 13: follow.v1.FollowService.FollowInfo:output_type -> follow.v1.FollowInfoResponse
	12, // 14: follow.v1.FollowService.GetFollowStatics:output_type -> follow.v1.GetFollowStaticsResponse
	9,  // [9:15] is the sub-list for method output_type
	3,  // [3:9] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_follow_v1_follow_proto_init() }
func file_follow_v1_follow_proto_init() {
	if File_follow_v1_follow_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_follow_v1_follow_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FollowRelation); i {
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
		file_follow_v1_follow_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FollowRequest); i {
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
		file_follow_v1_follow_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FollowResponse); i {
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
		file_follow_v1_follow_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CancelFollowRequest); i {
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
		file_follow_v1_follow_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CancelFollowResponse); i {
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
		file_follow_v1_follow_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFolloweeRequest); i {
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
		file_follow_v1_follow_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFolloweeResponse); i {
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
		file_follow_v1_follow_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFollowerRequest); i {
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
		file_follow_v1_follow_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFollowerResponse); i {
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
		file_follow_v1_follow_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FollowInfoRequest); i {
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
		file_follow_v1_follow_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FollowInfoResponse); i {
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
		file_follow_v1_follow_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFollowStaticsRequest); i {
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
		file_follow_v1_follow_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFollowStaticsResponse); i {
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
			RawDescriptor: file_follow_v1_follow_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_follow_v1_follow_proto_goTypes,
		DependencyIndexes: file_follow_v1_follow_proto_depIdxs,
		MessageInfos:      file_follow_v1_follow_proto_msgTypes,
	}.Build()
	File_follow_v1_follow_proto = out.File
	file_follow_v1_follow_proto_rawDesc = nil
	file_follow_v1_follow_proto_goTypes = nil
	file_follow_v1_follow_proto_depIdxs = nil
}
