// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.22.2
// source: gitlabexporter/protobuf/deployment.proto

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

type DeploymentStatus int32

const (
	DeploymentStatus_DEPLOYMENT_STATUS_UNSPECIFIED DeploymentStatus = 0
	DeploymentStatus_DEPLOYMENT_STATUS_CREATED     DeploymentStatus = 1
	DeploymentStatus_DEPLOYMENT_STATUS_RUNNING     DeploymentStatus = 2
	DeploymentStatus_DEPLOYMENT_STATUS_SUCCESS     DeploymentStatus = 3
	DeploymentStatus_DEPLOYMENT_STATUS_FAILED      DeploymentStatus = 4
	DeploymentStatus_DEPLOYMENT_STATUS_CANCELED    DeploymentStatus = 5
	DeploymentStatus_DEPLOYMENT_STATUS_SKIPPED     DeploymentStatus = 6
	DeploymentStatus_DEPLOYMENT_STATUS_BLOCKED     DeploymentStatus = 7
)

// Enum value maps for DeploymentStatus.
var (
	DeploymentStatus_name = map[int32]string{
		0: "DEPLOYMENT_STATUS_UNSPECIFIED",
		1: "DEPLOYMENT_STATUS_CREATED",
		2: "DEPLOYMENT_STATUS_RUNNING",
		3: "DEPLOYMENT_STATUS_SUCCESS",
		4: "DEPLOYMENT_STATUS_FAILED",
		5: "DEPLOYMENT_STATUS_CANCELED",
		6: "DEPLOYMENT_STATUS_SKIPPED",
		7: "DEPLOYMENT_STATUS_BLOCKED",
	}
	DeploymentStatus_value = map[string]int32{
		"DEPLOYMENT_STATUS_UNSPECIFIED": 0,
		"DEPLOYMENT_STATUS_CREATED":     1,
		"DEPLOYMENT_STATUS_RUNNING":     2,
		"DEPLOYMENT_STATUS_SUCCESS":     3,
		"DEPLOYMENT_STATUS_FAILED":      4,
		"DEPLOYMENT_STATUS_CANCELED":    5,
		"DEPLOYMENT_STATUS_SKIPPED":     6,
		"DEPLOYMENT_STATUS_BLOCKED":     7,
	}
)

func (x DeploymentStatus) Enum() *DeploymentStatus {
	p := new(DeploymentStatus)
	*p = x
	return p
}

func (x DeploymentStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DeploymentStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_gitlabexporter_protobuf_deployment_proto_enumTypes[0].Descriptor()
}

func (DeploymentStatus) Type() protoreflect.EnumType {
	return &file_gitlabexporter_protobuf_deployment_proto_enumTypes[0]
}

func (x DeploymentStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DeploymentStatus.Descriptor instead.
func (DeploymentStatus) EnumDescriptor() ([]byte, []int) {
	return file_gitlabexporter_protobuf_deployment_proto_rawDescGZIP(), []int{0}
}

type Deployment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          int64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Iid         int64                 `protobuf:"varint,2,opt,name=iid,proto3" json:"iid,omitempty"`
	Job         *JobReference         `protobuf:"bytes,3,opt,name=job,proto3" json:"job,omitempty"`
	Triggerer   *UserReference        `protobuf:"bytes,4,opt,name=triggerer,proto3" json:"triggerer,omitempty"`
	Environment *EnvironmentReference `protobuf:"bytes,5,opt,name=environment,proto3" json:"environment,omitempty"`
	Timestamps  *DeploymentTimestamps `protobuf:"bytes,6,opt,name=timestamps,proto3" json:"timestamps,omitempty"`
	Status      DeploymentStatus      `protobuf:"varint,7,opt,name=status,proto3,enum=gitlabexporter.protobuf.DeploymentStatus" json:"status,omitempty"`
	Ref         string                `protobuf:"bytes,8,opt,name=ref,proto3" json:"ref,omitempty"`
	Sha         string                `protobuf:"bytes,9,opt,name=sha,proto3" json:"sha,omitempty"`
}

func (x *Deployment) Reset() {
	*x = Deployment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gitlabexporter_protobuf_deployment_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Deployment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Deployment) ProtoMessage() {}

func (x *Deployment) ProtoReflect() protoreflect.Message {
	mi := &file_gitlabexporter_protobuf_deployment_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Deployment.ProtoReflect.Descriptor instead.
func (*Deployment) Descriptor() ([]byte, []int) {
	return file_gitlabexporter_protobuf_deployment_proto_rawDescGZIP(), []int{0}
}

func (x *Deployment) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Deployment) GetIid() int64 {
	if x != nil {
		return x.Iid
	}
	return 0
}

func (x *Deployment) GetJob() *JobReference {
	if x != nil {
		return x.Job
	}
	return nil
}

func (x *Deployment) GetTriggerer() *UserReference {
	if x != nil {
		return x.Triggerer
	}
	return nil
}

func (x *Deployment) GetEnvironment() *EnvironmentReference {
	if x != nil {
		return x.Environment
	}
	return nil
}

func (x *Deployment) GetTimestamps() *DeploymentTimestamps {
	if x != nil {
		return x.Timestamps
	}
	return nil
}

func (x *Deployment) GetStatus() DeploymentStatus {
	if x != nil {
		return x.Status
	}
	return DeploymentStatus_DEPLOYMENT_STATUS_UNSPECIFIED
}

func (x *Deployment) GetRef() string {
	if x != nil {
		return x.Ref
	}
	return ""
}

func (x *Deployment) GetSha() string {
	if x != nil {
		return x.Sha
	}
	return ""
}

type DeploymentTimestamps struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CreatedAt  *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	FinishedAt *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=finished_at,json=finishedAt,proto3" json:"finished_at,omitempty"`
	UpdatedAt  *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *DeploymentTimestamps) Reset() {
	*x = DeploymentTimestamps{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gitlabexporter_protobuf_deployment_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeploymentTimestamps) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeploymentTimestamps) ProtoMessage() {}

func (x *DeploymentTimestamps) ProtoReflect() protoreflect.Message {
	mi := &file_gitlabexporter_protobuf_deployment_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeploymentTimestamps.ProtoReflect.Descriptor instead.
func (*DeploymentTimestamps) Descriptor() ([]byte, []int) {
	return file_gitlabexporter_protobuf_deployment_proto_rawDescGZIP(), []int{1}
}

func (x *DeploymentTimestamps) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *DeploymentTimestamps) GetFinishedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.FinishedAt
	}
	return nil
}

func (x *DeploymentTimestamps) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

var File_gitlabexporter_protobuf_deployment_proto protoreflect.FileDescriptor

var file_gitlabexporter_protobuf_deployment_proto_rawDesc = []byte{
	0x0a, 0x28, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x67, 0x69, 0x74, 0x6c,
	0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x28, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f,
	0x72, 0x74, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x72, 0x65,
	0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb4,
	0x03, 0x0a, 0x0a, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a,
	0x03, 0x69, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x69, 0x69, 0x64, 0x12,
	0x37, 0x0a, 0x03, 0x6a, 0x6f, 0x62, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x67,
	0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65,
	0x6e, 0x63, 0x65, 0x52, 0x03, 0x6a, 0x6f, 0x62, 0x12, 0x44, 0x0a, 0x09, 0x74, 0x72, 0x69, 0x67,
	0x67, 0x65, 0x72, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x67, 0x69,
	0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65,
	0x6e, 0x63, 0x65, 0x52, 0x09, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x65, 0x72, 0x12, 0x4f,
	0x0a, 0x0b, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f,
	0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6e,
	0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e,
	0x63, 0x65, 0x52, 0x0b, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x12,
	0x4d, 0x0a, 0x0a, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x73, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f,
	0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x65,
	0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x73, 0x52, 0x0a, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x73, 0x12, 0x41,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x29,
	0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d,
	0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x65, 0x66, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x72, 0x65, 0x66, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x68, 0x61, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x73, 0x68, 0x61, 0x22, 0xc9, 0x01, 0x0a, 0x14, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x73, 0x12, 0x39,
	0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x3b, 0x0a, 0x0b, 0x66, 0x69, 0x6e,
	0x69, 0x73, 0x68, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x66, 0x69, 0x6e, 0x69,
	0x73, 0x68, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x2a, 0x8e, 0x02, 0x0a, 0x10, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x21, 0x0a, 0x1d, 0x44, 0x45, 0x50, 0x4c, 0x4f, 0x59,
	0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x53, 0x50,
	0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1d, 0x0a, 0x19, 0x44, 0x45, 0x50,
	0x4c, 0x4f, 0x59, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x43,
	0x52, 0x45, 0x41, 0x54, 0x45, 0x44, 0x10, 0x01, 0x12, 0x1d, 0x0a, 0x19, 0x44, 0x45, 0x50, 0x4c,
	0x4f, 0x59, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x52, 0x55,
	0x4e, 0x4e, 0x49, 0x4e, 0x47, 0x10, 0x02, 0x12, 0x1d, 0x0a, 0x19, 0x44, 0x45, 0x50, 0x4c, 0x4f,
	0x59, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x53, 0x55, 0x43,
	0x43, 0x45, 0x53, 0x53, 0x10, 0x03, 0x12, 0x1c, 0x0a, 0x18, 0x44, 0x45, 0x50, 0x4c, 0x4f, 0x59,
	0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x46, 0x41, 0x49, 0x4c,
	0x45, 0x44, 0x10, 0x04, 0x12, 0x1e, 0x0a, 0x1a, 0x44, 0x45, 0x50, 0x4c, 0x4f, 0x59, 0x4d, 0x45,
	0x4e, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x43, 0x41, 0x4e, 0x43, 0x45, 0x4c,
	0x45, 0x44, 0x10, 0x05, 0x12, 0x1d, 0x0a, 0x19, 0x44, 0x45, 0x50, 0x4c, 0x4f, 0x59, 0x4d, 0x45,
	0x4e, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x53, 0x4b, 0x49, 0x50, 0x50, 0x45,
	0x44, 0x10, 0x06, 0x12, 0x1d, 0x0a, 0x19, 0x44, 0x45, 0x50, 0x4c, 0x4f, 0x59, 0x4d, 0x45, 0x4e,
	0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x42, 0x4c, 0x4f, 0x43, 0x4b, 0x45, 0x44,
	0x10, 0x07, 0x42, 0x37, 0x5a, 0x35, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x63, 0x6c, 0x75, 0x74, 0x74, 0x72, 0x64, 0x65, 0x76, 0x2f, 0x67, 0x69, 0x74, 0x6c, 0x61,
	0x62, 0x2d, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_gitlabexporter_protobuf_deployment_proto_rawDescOnce sync.Once
	file_gitlabexporter_protobuf_deployment_proto_rawDescData = file_gitlabexporter_protobuf_deployment_proto_rawDesc
)

func file_gitlabexporter_protobuf_deployment_proto_rawDescGZIP() []byte {
	file_gitlabexporter_protobuf_deployment_proto_rawDescOnce.Do(func() {
		file_gitlabexporter_protobuf_deployment_proto_rawDescData = protoimpl.X.CompressGZIP(file_gitlabexporter_protobuf_deployment_proto_rawDescData)
	})
	return file_gitlabexporter_protobuf_deployment_proto_rawDescData
}

var file_gitlabexporter_protobuf_deployment_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_gitlabexporter_protobuf_deployment_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_gitlabexporter_protobuf_deployment_proto_goTypes = []interface{}{
	(DeploymentStatus)(0),         // 0: gitlabexporter.protobuf.DeploymentStatus
	(*Deployment)(nil),            // 1: gitlabexporter.protobuf.Deployment
	(*DeploymentTimestamps)(nil),  // 2: gitlabexporter.protobuf.DeploymentTimestamps
	(*JobReference)(nil),          // 3: gitlabexporter.protobuf.JobReference
	(*UserReference)(nil),         // 4: gitlabexporter.protobuf.UserReference
	(*EnvironmentReference)(nil),  // 5: gitlabexporter.protobuf.EnvironmentReference
	(*timestamppb.Timestamp)(nil), // 6: google.protobuf.Timestamp
}
var file_gitlabexporter_protobuf_deployment_proto_depIdxs = []int32{
	3, // 0: gitlabexporter.protobuf.Deployment.job:type_name -> gitlabexporter.protobuf.JobReference
	4, // 1: gitlabexporter.protobuf.Deployment.triggerer:type_name -> gitlabexporter.protobuf.UserReference
	5, // 2: gitlabexporter.protobuf.Deployment.environment:type_name -> gitlabexporter.protobuf.EnvironmentReference
	2, // 3: gitlabexporter.protobuf.Deployment.timestamps:type_name -> gitlabexporter.protobuf.DeploymentTimestamps
	0, // 4: gitlabexporter.protobuf.Deployment.status:type_name -> gitlabexporter.protobuf.DeploymentStatus
	6, // 5: gitlabexporter.protobuf.DeploymentTimestamps.created_at:type_name -> google.protobuf.Timestamp
	6, // 6: gitlabexporter.protobuf.DeploymentTimestamps.finished_at:type_name -> google.protobuf.Timestamp
	6, // 7: gitlabexporter.protobuf.DeploymentTimestamps.updated_at:type_name -> google.protobuf.Timestamp
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_gitlabexporter_protobuf_deployment_proto_init() }
func file_gitlabexporter_protobuf_deployment_proto_init() {
	if File_gitlabexporter_protobuf_deployment_proto != nil {
		return
	}
	file_gitlabexporter_protobuf_references_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_gitlabexporter_protobuf_deployment_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Deployment); i {
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
		file_gitlabexporter_protobuf_deployment_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeploymentTimestamps); i {
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
			RawDescriptor: file_gitlabexporter_protobuf_deployment_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_gitlabexporter_protobuf_deployment_proto_goTypes,
		DependencyIndexes: file_gitlabexporter_protobuf_deployment_proto_depIdxs,
		EnumInfos:         file_gitlabexporter_protobuf_deployment_proto_enumTypes,
		MessageInfos:      file_gitlabexporter_protobuf_deployment_proto_msgTypes,
	}.Build()
	File_gitlabexporter_protobuf_deployment_proto = out.File
	file_gitlabexporter_protobuf_deployment_proto_rawDesc = nil
	file_gitlabexporter_protobuf_deployment_proto_goTypes = nil
	file_gitlabexporter_protobuf_deployment_proto_depIdxs = nil
}
