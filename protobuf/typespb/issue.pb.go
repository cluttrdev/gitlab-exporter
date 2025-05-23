// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.22.2
// source: gitlabexporter/protobuf/issue.proto

package typespb

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

type IssueType int32

const (
	IssueType_ISSUE_TYPE_UNSPECIFIED IssueType = 0
	IssueType_ISSUE_TYPE_UNKNOWN     IssueType = 1
	IssueType_ISSUE_TYPE_ISSUE       IssueType = 2
	IssueType_ISSUE_TYPE_INCIDENT    IssueType = 3
	IssueType_ISSUE_TYPE_TEST_CASE   IssueType = 4
	IssueType_ISSUE_TYPE_REQUIREMENT IssueType = 5
	IssueType_ISSUE_TYPE_TASK        IssueType = 6
	IssueType_ISSUE_TYPE_TICKET      IssueType = 7
	IssueType_ISSUE_TYPE_OBJECTIVE   IssueType = 8
	IssueType_ISSUE_TYPE_KEY_RESULT  IssueType = 9
	IssueType_ISSUE_TYPE_EPIC        IssueType = 10
)

// Enum value maps for IssueType.
var (
	IssueType_name = map[int32]string{
		0:  "ISSUE_TYPE_UNSPECIFIED",
		1:  "ISSUE_TYPE_UNKNOWN",
		2:  "ISSUE_TYPE_ISSUE",
		3:  "ISSUE_TYPE_INCIDENT",
		4:  "ISSUE_TYPE_TEST_CASE",
		5:  "ISSUE_TYPE_REQUIREMENT",
		6:  "ISSUE_TYPE_TASK",
		7:  "ISSUE_TYPE_TICKET",
		8:  "ISSUE_TYPE_OBJECTIVE",
		9:  "ISSUE_TYPE_KEY_RESULT",
		10: "ISSUE_TYPE_EPIC",
	}
	IssueType_value = map[string]int32{
		"ISSUE_TYPE_UNSPECIFIED": 0,
		"ISSUE_TYPE_UNKNOWN":     1,
		"ISSUE_TYPE_ISSUE":       2,
		"ISSUE_TYPE_INCIDENT":    3,
		"ISSUE_TYPE_TEST_CASE":   4,
		"ISSUE_TYPE_REQUIREMENT": 5,
		"ISSUE_TYPE_TASK":        6,
		"ISSUE_TYPE_TICKET":      7,
		"ISSUE_TYPE_OBJECTIVE":   8,
		"ISSUE_TYPE_KEY_RESULT":  9,
		"ISSUE_TYPE_EPIC":        10,
	}
)

func (x IssueType) Enum() *IssueType {
	p := new(IssueType)
	*p = x
	return p
}

func (x IssueType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (IssueType) Descriptor() protoreflect.EnumDescriptor {
	return file_gitlabexporter_protobuf_issue_proto_enumTypes[0].Descriptor()
}

func (IssueType) Type() protoreflect.EnumType {
	return &file_gitlabexporter_protobuf_issue_proto_enumTypes[0]
}

func (x IssueType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use IssueType.Descriptor instead.
func (IssueType) EnumDescriptor() ([]byte, []int) {
	return file_gitlabexporter_protobuf_issue_proto_rawDescGZIP(), []int{0}
}

type IssueSeverity int32

const (
	IssueSeverity_ISSUE_SEVERITY_UNSPECIFIED IssueSeverity = 0
	IssueSeverity_ISSUE_SEVERITY_UNKNOWN     IssueSeverity = 1
	IssueSeverity_ISSUE_SEVERITY_LOW         IssueSeverity = 2
	IssueSeverity_ISSUE_SEVERITY_MEDIUM      IssueSeverity = 3
	IssueSeverity_ISSUE_SEVERITY_HIGH        IssueSeverity = 4
	IssueSeverity_ISSUE_SEVERITY_CRITICAL    IssueSeverity = 5
)

// Enum value maps for IssueSeverity.
var (
	IssueSeverity_name = map[int32]string{
		0: "ISSUE_SEVERITY_UNSPECIFIED",
		1: "ISSUE_SEVERITY_UNKNOWN",
		2: "ISSUE_SEVERITY_LOW",
		3: "ISSUE_SEVERITY_MEDIUM",
		4: "ISSUE_SEVERITY_HIGH",
		5: "ISSUE_SEVERITY_CRITICAL",
	}
	IssueSeverity_value = map[string]int32{
		"ISSUE_SEVERITY_UNSPECIFIED": 0,
		"ISSUE_SEVERITY_UNKNOWN":     1,
		"ISSUE_SEVERITY_LOW":         2,
		"ISSUE_SEVERITY_MEDIUM":      3,
		"ISSUE_SEVERITY_HIGH":        4,
		"ISSUE_SEVERITY_CRITICAL":    5,
	}
)

func (x IssueSeverity) Enum() *IssueSeverity {
	p := new(IssueSeverity)
	*p = x
	return p
}

func (x IssueSeverity) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (IssueSeverity) Descriptor() protoreflect.EnumDescriptor {
	return file_gitlabexporter_protobuf_issue_proto_enumTypes[1].Descriptor()
}

func (IssueSeverity) Type() protoreflect.EnumType {
	return &file_gitlabexporter_protobuf_issue_proto_enumTypes[1]
}

func (x IssueSeverity) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use IssueSeverity.Descriptor instead.
func (IssueSeverity) EnumDescriptor() ([]byte, []int) {
	return file_gitlabexporter_protobuf_issue_proto_rawDescGZIP(), []int{1}
}

type IssueState int32

const (
	IssueState_ISSUE_STATE_UNSPECIFIED IssueState = 0
	IssueState_ISSUE_STATE_UNKNOWN     IssueState = 1
	IssueState_ISSUE_STATE_OPENED      IssueState = 2
	IssueState_ISSUE_STATE_CLOSED      IssueState = 3
	IssueState_ISSUE_STATE_LOCKED      IssueState = 4
	IssueState_ISSUE_STATE_ALL         IssueState = 5
)

// Enum value maps for IssueState.
var (
	IssueState_name = map[int32]string{
		0: "ISSUE_STATE_UNSPECIFIED",
		1: "ISSUE_STATE_UNKNOWN",
		2: "ISSUE_STATE_OPENED",
		3: "ISSUE_STATE_CLOSED",
		4: "ISSUE_STATE_LOCKED",
		5: "ISSUE_STATE_ALL",
	}
	IssueState_value = map[string]int32{
		"ISSUE_STATE_UNSPECIFIED": 0,
		"ISSUE_STATE_UNKNOWN":     1,
		"ISSUE_STATE_OPENED":      2,
		"ISSUE_STATE_CLOSED":      3,
		"ISSUE_STATE_LOCKED":      4,
		"ISSUE_STATE_ALL":         5,
	}
)

func (x IssueState) Enum() *IssueState {
	p := new(IssueState)
	*p = x
	return p
}

func (x IssueState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (IssueState) Descriptor() protoreflect.EnumDescriptor {
	return file_gitlabexporter_protobuf_issue_proto_enumTypes[2].Descriptor()
}

func (IssueState) Type() protoreflect.EnumType {
	return &file_gitlabexporter_protobuf_issue_proto_enumTypes[2]
}

func (x IssueState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use IssueState.Descriptor instead.
func (IssueState) EnumDescriptor() ([]byte, []int) {
	return file_gitlabexporter_protobuf_issue_proto_rawDescGZIP(), []int{2}
}

type Issue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         int64             `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Iid        int64             `protobuf:"varint,2,opt,name=iid,proto3" json:"iid,omitempty"`
	Project    *ProjectReference `protobuf:"bytes,3,opt,name=project,proto3" json:"project,omitempty"`
	Timestamps *IssueTimestamps  `protobuf:"bytes,4,opt,name=timestamps,proto3" json:"timestamps,omitempty"`
	Title      string            `protobuf:"bytes,5,opt,name=title,proto3" json:"title,omitempty"`
	Labels     []string          `protobuf:"bytes,6,rep,name=labels,proto3" json:"labels,omitempty"`
	Type       IssueType         `protobuf:"varint,7,opt,name=type,proto3,enum=gitlabexporter.protobuf.IssueType" json:"type,omitempty"`
	Severity   IssueSeverity     `protobuf:"varint,8,opt,name=severity,proto3,enum=gitlabexporter.protobuf.IssueSeverity" json:"severity,omitempty"`
	State      IssueState        `protobuf:"varint,9,opt,name=state,proto3,enum=gitlabexporter.protobuf.IssueState" json:"state,omitempty"`
}

func (x *Issue) Reset() {
	*x = Issue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gitlabexporter_protobuf_issue_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Issue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Issue) ProtoMessage() {}

func (x *Issue) ProtoReflect() protoreflect.Message {
	mi := &file_gitlabexporter_protobuf_issue_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Issue.ProtoReflect.Descriptor instead.
func (*Issue) Descriptor() ([]byte, []int) {
	return file_gitlabexporter_protobuf_issue_proto_rawDescGZIP(), []int{0}
}

func (x *Issue) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Issue) GetIid() int64 {
	if x != nil {
		return x.Iid
	}
	return 0
}

func (x *Issue) GetProject() *ProjectReference {
	if x != nil {
		return x.Project
	}
	return nil
}

func (x *Issue) GetTimestamps() *IssueTimestamps {
	if x != nil {
		return x.Timestamps
	}
	return nil
}

func (x *Issue) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Issue) GetLabels() []string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *Issue) GetType() IssueType {
	if x != nil {
		return x.Type
	}
	return IssueType_ISSUE_TYPE_UNSPECIFIED
}

func (x *Issue) GetSeverity() IssueSeverity {
	if x != nil {
		return x.Severity
	}
	return IssueSeverity_ISSUE_SEVERITY_UNSPECIFIED
}

func (x *Issue) GetState() IssueState {
	if x != nil {
		return x.State
	}
	return IssueState_ISSUE_STATE_UNSPECIFIED
}

type IssueTimestamps struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	ClosedAt  *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=closed_at,json=closedAt,proto3" json:"closed_at,omitempty"`
}

func (x *IssueTimestamps) Reset() {
	*x = IssueTimestamps{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gitlabexporter_protobuf_issue_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IssueTimestamps) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IssueTimestamps) ProtoMessage() {}

func (x *IssueTimestamps) ProtoReflect() protoreflect.Message {
	mi := &file_gitlabexporter_protobuf_issue_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IssueTimestamps.ProtoReflect.Descriptor instead.
func (*IssueTimestamps) Descriptor() ([]byte, []int) {
	return file_gitlabexporter_protobuf_issue_proto_rawDescGZIP(), []int{1}
}

func (x *IssueTimestamps) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *IssueTimestamps) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *IssueTimestamps) GetClosedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ClosedAt
	}
	return nil
}

var File_gitlabexporter_protobuf_issue_proto protoreflect.FileDescriptor

var file_gitlabexporter_protobuf_issue_proto_rawDesc = []byte{
	0x0a, 0x23, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x69, 0x73, 0x73, 0x75, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70,
	0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x28, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e,
	0x63, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9d, 0x03, 0x0a, 0x05, 0x49, 0x73,
	0x73, 0x75, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x03, 0x69, 0x69, 0x64, 0x12, 0x43, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65,
	0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63,
	0x65, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x48, 0x0a, 0x0a, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28,
	0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x49, 0x73, 0x73, 0x75, 0x65, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x73, 0x52, 0x0a, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x61,
	0x62, 0x65, 0x6c, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65,
	0x6c, 0x73, 0x12, 0x36, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x22, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x49, 0x73, 0x73, 0x75, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x42, 0x0a, 0x08, 0x73, 0x65,
	0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x26, 0x2e, 0x67,
	0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x49, 0x73, 0x73, 0x75, 0x65, 0x53, 0x65, 0x76, 0x65,
	0x72, 0x69, 0x74, 0x79, 0x52, 0x08, 0x73, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x12, 0x39,
	0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x23, 0x2e,
	0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x49, 0x73, 0x73, 0x75, 0x65, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22, 0xc0, 0x01, 0x0a, 0x0f, 0x49, 0x73,
	0x73, 0x75, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x73, 0x12, 0x39, 0x0a,
	0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x12, 0x37, 0x0a, 0x09, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x08, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x64, 0x41, 0x74, 0x2a, 0x9a, 0x02, 0x0a,
	0x09, 0x49, 0x73, 0x73, 0x75, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x16, 0x49, 0x53,
	0x53, 0x55, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49,
	0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12, 0x49, 0x53, 0x53, 0x55, 0x45, 0x5f,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x01, 0x12, 0x14,
	0x0a, 0x10, 0x49, 0x53, 0x53, 0x55, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x49, 0x53, 0x53,
	0x55, 0x45, 0x10, 0x02, 0x12, 0x17, 0x0a, 0x13, 0x49, 0x53, 0x53, 0x55, 0x45, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x49, 0x4e, 0x43, 0x49, 0x44, 0x45, 0x4e, 0x54, 0x10, 0x03, 0x12, 0x18, 0x0a,
	0x14, 0x49, 0x53, 0x53, 0x55, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x54, 0x45, 0x53, 0x54,
	0x5f, 0x43, 0x41, 0x53, 0x45, 0x10, 0x04, 0x12, 0x1a, 0x0a, 0x16, 0x49, 0x53, 0x53, 0x55, 0x45,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x52, 0x45, 0x51, 0x55, 0x49, 0x52, 0x45, 0x4d, 0x45, 0x4e,
	0x54, 0x10, 0x05, 0x12, 0x13, 0x0a, 0x0f, 0x49, 0x53, 0x53, 0x55, 0x45, 0x5f, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x54, 0x41, 0x53, 0x4b, 0x10, 0x06, 0x12, 0x15, 0x0a, 0x11, 0x49, 0x53, 0x53, 0x55,
	0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x54, 0x49, 0x43, 0x4b, 0x45, 0x54, 0x10, 0x07, 0x12,
	0x18, 0x0a, 0x14, 0x49, 0x53, 0x53, 0x55, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4f, 0x42,
	0x4a, 0x45, 0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x08, 0x12, 0x19, 0x0a, 0x15, 0x49, 0x53, 0x53,
	0x55, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4b, 0x45, 0x59, 0x5f, 0x52, 0x45, 0x53, 0x55,
	0x4c, 0x54, 0x10, 0x09, 0x12, 0x13, 0x0a, 0x0f, 0x49, 0x53, 0x53, 0x55, 0x45, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x45, 0x50, 0x49, 0x43, 0x10, 0x0a, 0x2a, 0xb4, 0x01, 0x0a, 0x0d, 0x49, 0x73,
	0x73, 0x75, 0x65, 0x53, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x12, 0x1e, 0x0a, 0x1a, 0x49,
	0x53, 0x53, 0x55, 0x45, 0x5f, 0x53, 0x45, 0x56, 0x45, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x55, 0x4e,
	0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1a, 0x0a, 0x16, 0x49,
	0x53, 0x53, 0x55, 0x45, 0x5f, 0x53, 0x45, 0x56, 0x45, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x55, 0x4e,
	0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x01, 0x12, 0x16, 0x0a, 0x12, 0x49, 0x53, 0x53, 0x55, 0x45,
	0x5f, 0x53, 0x45, 0x56, 0x45, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x4c, 0x4f, 0x57, 0x10, 0x02, 0x12,
	0x19, 0x0a, 0x15, 0x49, 0x53, 0x53, 0x55, 0x45, 0x5f, 0x53, 0x45, 0x56, 0x45, 0x52, 0x49, 0x54,
	0x59, 0x5f, 0x4d, 0x45, 0x44, 0x49, 0x55, 0x4d, 0x10, 0x03, 0x12, 0x17, 0x0a, 0x13, 0x49, 0x53,
	0x53, 0x55, 0x45, 0x5f, 0x53, 0x45, 0x56, 0x45, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x48, 0x49, 0x47,
	0x48, 0x10, 0x04, 0x12, 0x1b, 0x0a, 0x17, 0x49, 0x53, 0x53, 0x55, 0x45, 0x5f, 0x53, 0x45, 0x56,
	0x45, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x43, 0x52, 0x49, 0x54, 0x49, 0x43, 0x41, 0x4c, 0x10, 0x05,
	0x2a, 0x9f, 0x01, 0x0a, 0x0a, 0x49, 0x73, 0x73, 0x75, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12,
	0x1b, 0x0a, 0x17, 0x49, 0x53, 0x53, 0x55, 0x45, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x55,
	0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x17, 0x0a, 0x13,
	0x49, 0x53, 0x53, 0x55, 0x45, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x55, 0x4e, 0x4b, 0x4e,
	0x4f, 0x57, 0x4e, 0x10, 0x01, 0x12, 0x16, 0x0a, 0x12, 0x49, 0x53, 0x53, 0x55, 0x45, 0x5f, 0x53,
	0x54, 0x41, 0x54, 0x45, 0x5f, 0x4f, 0x50, 0x45, 0x4e, 0x45, 0x44, 0x10, 0x02, 0x12, 0x16, 0x0a,
	0x12, 0x49, 0x53, 0x53, 0x55, 0x45, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x43, 0x4c, 0x4f,
	0x53, 0x45, 0x44, 0x10, 0x03, 0x12, 0x16, 0x0a, 0x12, 0x49, 0x53, 0x53, 0x55, 0x45, 0x5f, 0x53,
	0x54, 0x41, 0x54, 0x45, 0x5f, 0x4c, 0x4f, 0x43, 0x4b, 0x45, 0x44, 0x10, 0x04, 0x12, 0x13, 0x0a,
	0x0f, 0x49, 0x53, 0x53, 0x55, 0x45, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x41, 0x4c, 0x4c,
	0x10, 0x05, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x6f, 0x2e, 0x63, 0x6c, 0x75, 0x74, 0x74, 0x72, 0x2e,
	0x64, 0x65, 0x76, 0x2f, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2d, 0x65, 0x78, 0x70, 0x6f, 0x72,
	0x74, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x79, 0x70,
	0x65, 0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gitlabexporter_protobuf_issue_proto_rawDescOnce sync.Once
	file_gitlabexporter_protobuf_issue_proto_rawDescData = file_gitlabexporter_protobuf_issue_proto_rawDesc
)

func file_gitlabexporter_protobuf_issue_proto_rawDescGZIP() []byte {
	file_gitlabexporter_protobuf_issue_proto_rawDescOnce.Do(func() {
		file_gitlabexporter_protobuf_issue_proto_rawDescData = protoimpl.X.CompressGZIP(file_gitlabexporter_protobuf_issue_proto_rawDescData)
	})
	return file_gitlabexporter_protobuf_issue_proto_rawDescData
}

var file_gitlabexporter_protobuf_issue_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_gitlabexporter_protobuf_issue_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_gitlabexporter_protobuf_issue_proto_goTypes = []interface{}{
	(IssueType)(0),                // 0: gitlabexporter.protobuf.IssueType
	(IssueSeverity)(0),            // 1: gitlabexporter.protobuf.IssueSeverity
	(IssueState)(0),               // 2: gitlabexporter.protobuf.IssueState
	(*Issue)(nil),                 // 3: gitlabexporter.protobuf.Issue
	(*IssueTimestamps)(nil),       // 4: gitlabexporter.protobuf.IssueTimestamps
	(*ProjectReference)(nil),      // 5: gitlabexporter.protobuf.ProjectReference
	(*timestamppb.Timestamp)(nil), // 6: google.protobuf.Timestamp
}
var file_gitlabexporter_protobuf_issue_proto_depIdxs = []int32{
	5, // 0: gitlabexporter.protobuf.Issue.project:type_name -> gitlabexporter.protobuf.ProjectReference
	4, // 1: gitlabexporter.protobuf.Issue.timestamps:type_name -> gitlabexporter.protobuf.IssueTimestamps
	0, // 2: gitlabexporter.protobuf.Issue.type:type_name -> gitlabexporter.protobuf.IssueType
	1, // 3: gitlabexporter.protobuf.Issue.severity:type_name -> gitlabexporter.protobuf.IssueSeverity
	2, // 4: gitlabexporter.protobuf.Issue.state:type_name -> gitlabexporter.protobuf.IssueState
	6, // 5: gitlabexporter.protobuf.IssueTimestamps.created_at:type_name -> google.protobuf.Timestamp
	6, // 6: gitlabexporter.protobuf.IssueTimestamps.updated_at:type_name -> google.protobuf.Timestamp
	6, // 7: gitlabexporter.protobuf.IssueTimestamps.closed_at:type_name -> google.protobuf.Timestamp
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_gitlabexporter_protobuf_issue_proto_init() }
func file_gitlabexporter_protobuf_issue_proto_init() {
	if File_gitlabexporter_protobuf_issue_proto != nil {
		return
	}
	file_gitlabexporter_protobuf_references_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_gitlabexporter_protobuf_issue_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Issue); i {
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
		file_gitlabexporter_protobuf_issue_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IssueTimestamps); i {
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
			RawDescriptor: file_gitlabexporter_protobuf_issue_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_gitlabexporter_protobuf_issue_proto_goTypes,
		DependencyIndexes: file_gitlabexporter_protobuf_issue_proto_depIdxs,
		EnumInfos:         file_gitlabexporter_protobuf_issue_proto_enumTypes,
		MessageInfos:      file_gitlabexporter_protobuf_issue_proto_msgTypes,
	}.Build()
	File_gitlabexporter_protobuf_issue_proto = out.File
	file_gitlabexporter_protobuf_issue_proto_rawDesc = nil
	file_gitlabexporter_protobuf_issue_proto_goTypes = nil
	file_gitlabexporter_protobuf_issue_proto_depIdxs = nil
}
