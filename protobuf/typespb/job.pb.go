// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.22.2
// source: gitlabexporter/protobuf/job.proto

package typespb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
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

type JobKind int32

const (
	JobKind_JOBKIND_UNSPECIFIED JobKind = 0
	JobKind_JOBKIND_BUILD       JobKind = 1
	JobKind_JOBKIND_BRIDGE      JobKind = 2
)

// Enum value maps for JobKind.
var (
	JobKind_name = map[int32]string{
		0: "JOBKIND_UNSPECIFIED",
		1: "JOBKIND_BUILD",
		2: "JOBKIND_BRIDGE",
	}
	JobKind_value = map[string]int32{
		"JOBKIND_UNSPECIFIED": 0,
		"JOBKIND_BUILD":       1,
		"JOBKIND_BRIDGE":      2,
	}
)

func (x JobKind) Enum() *JobKind {
	p := new(JobKind)
	*p = x
	return p
}

func (x JobKind) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (JobKind) Descriptor() protoreflect.EnumDescriptor {
	return file_gitlabexporter_protobuf_job_proto_enumTypes[0].Descriptor()
}

func (JobKind) Type() protoreflect.EnumType {
	return &file_gitlabexporter_protobuf_job_proto_enumTypes[0]
}

func (x JobKind) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use JobKind.Descriptor instead.
func (JobKind) EnumDescriptor() ([]byte, []int) {
	return file_gitlabexporter_protobuf_job_proto_rawDescGZIP(), []int{0}
}

type Job struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                 int64                `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name               string               `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Pipeline           *PipelineReference   `protobuf:"bytes,3,opt,name=pipeline,proto3" json:"pipeline,omitempty"`
	Ref                string               `protobuf:"bytes,4,opt,name=ref,proto3" json:"ref,omitempty"`
	RefPath            string               `protobuf:"bytes,5,opt,name=ref_path,json=refPath,proto3" json:"ref_path,omitempty"`
	Status             string               `protobuf:"bytes,6,opt,name=status,proto3" json:"status,omitempty"`
	FailureReason      string               `protobuf:"bytes,7,opt,name=failure_reason,json=failureReason,proto3" json:"failure_reason,omitempty"`
	Timestamps         *JobTimestamps       `protobuf:"bytes,8,opt,name=timestamps,proto3" json:"timestamps,omitempty"`
	QueuedDuration     *durationpb.Duration `protobuf:"bytes,9,opt,name=queued_duration,json=queuedDuration,proto3" json:"queued_duration,omitempty"`
	Duration           *durationpb.Duration `protobuf:"bytes,10,opt,name=duration,proto3" json:"duration,omitempty"`
	Coverage           float64              `protobuf:"fixed64,11,opt,name=coverage,proto3" json:"coverage,omitempty"`
	Stage              string               `protobuf:"bytes,12,opt,name=stage,proto3" json:"stage,omitempty"`
	Tags               []string             `protobuf:"bytes,13,rep,name=tags,proto3" json:"tags,omitempty"`
	AllowFailure       bool                 `protobuf:"varint,14,opt,name=allow_failure,json=allowFailure,proto3" json:"allow_failure,omitempty"`
	Manual             bool                 `protobuf:"varint,15,opt,name=manual,proto3" json:"manual,omitempty"`
	Retried            bool                 `protobuf:"varint,16,opt,name=retried,proto3" json:"retried,omitempty"`
	Retryable          bool                 `protobuf:"varint,17,opt,name=retryable,proto3" json:"retryable,omitempty"`
	Kind               JobKind              `protobuf:"varint,18,opt,name=kind,proto3,enum=gitlabexporter.protobuf.JobKind" json:"kind,omitempty"`
	DownstreamPipeline *PipelineReference   `protobuf:"bytes,19,opt,name=downstream_pipeline,json=downstreamPipeline,proto3,oneof" json:"downstream_pipeline,omitempty"`
	Runner             *RunnerReference     `protobuf:"bytes,20,opt,name=runner,proto3" json:"runner,omitempty"`
}

func (x *Job) Reset() {
	*x = Job{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gitlabexporter_protobuf_job_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Job) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Job) ProtoMessage() {}

func (x *Job) ProtoReflect() protoreflect.Message {
	mi := &file_gitlabexporter_protobuf_job_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Job.ProtoReflect.Descriptor instead.
func (*Job) Descriptor() ([]byte, []int) {
	return file_gitlabexporter_protobuf_job_proto_rawDescGZIP(), []int{0}
}

func (x *Job) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Job) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Job) GetPipeline() *PipelineReference {
	if x != nil {
		return x.Pipeline
	}
	return nil
}

func (x *Job) GetRef() string {
	if x != nil {
		return x.Ref
	}
	return ""
}

func (x *Job) GetRefPath() string {
	if x != nil {
		return x.RefPath
	}
	return ""
}

func (x *Job) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Job) GetFailureReason() string {
	if x != nil {
		return x.FailureReason
	}
	return ""
}

func (x *Job) GetTimestamps() *JobTimestamps {
	if x != nil {
		return x.Timestamps
	}
	return nil
}

func (x *Job) GetQueuedDuration() *durationpb.Duration {
	if x != nil {
		return x.QueuedDuration
	}
	return nil
}

func (x *Job) GetDuration() *durationpb.Duration {
	if x != nil {
		return x.Duration
	}
	return nil
}

func (x *Job) GetCoverage() float64 {
	if x != nil {
		return x.Coverage
	}
	return 0
}

func (x *Job) GetStage() string {
	if x != nil {
		return x.Stage
	}
	return ""
}

func (x *Job) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *Job) GetAllowFailure() bool {
	if x != nil {
		return x.AllowFailure
	}
	return false
}

func (x *Job) GetManual() bool {
	if x != nil {
		return x.Manual
	}
	return false
}

func (x *Job) GetRetried() bool {
	if x != nil {
		return x.Retried
	}
	return false
}

func (x *Job) GetRetryable() bool {
	if x != nil {
		return x.Retryable
	}
	return false
}

func (x *Job) GetKind() JobKind {
	if x != nil {
		return x.Kind
	}
	return JobKind_JOBKIND_UNSPECIFIED
}

func (x *Job) GetDownstreamPipeline() *PipelineReference {
	if x != nil {
		return x.DownstreamPipeline
	}
	return nil
}

func (x *Job) GetRunner() *RunnerReference {
	if x != nil {
		return x.Runner
	}
	return nil
}

type JobTimestamps struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CreatedAt  *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	QueuedAt   *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=queued_at,json=queuedAt,proto3" json:"queued_at,omitempty"`
	StartedAt  *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=started_at,json=startedAt,proto3" json:"started_at,omitempty"`
	FinishedAt *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=finished_at,json=finishedAt,proto3" json:"finished_at,omitempty"`
	ErasedAt   *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=erased_at,json=erasedAt,proto3" json:"erased_at,omitempty"`
}

func (x *JobTimestamps) Reset() {
	*x = JobTimestamps{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gitlabexporter_protobuf_job_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobTimestamps) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobTimestamps) ProtoMessage() {}

func (x *JobTimestamps) ProtoReflect() protoreflect.Message {
	mi := &file_gitlabexporter_protobuf_job_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobTimestamps.ProtoReflect.Descriptor instead.
func (*JobTimestamps) Descriptor() ([]byte, []int) {
	return file_gitlabexporter_protobuf_job_proto_rawDescGZIP(), []int{1}
}

func (x *JobTimestamps) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *JobTimestamps) GetQueuedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.QueuedAt
	}
	return nil
}

func (x *JobTimestamps) GetStartedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.StartedAt
	}
	return nil
}

func (x *JobTimestamps) GetFinishedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.FinishedAt
	}
	return nil
}

func (x *JobTimestamps) GetErasedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ErasedAt
	}
	return nil
}

var File_gitlabexporter_protobuf_job_proto protoreflect.FileDescriptor

var file_gitlabexporter_protobuf_job_proto_rawDesc = []byte{
	0x0a, 0x21, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x6a, 0x6f, 0x62, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x17, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72,
	0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x28, 0x67,
	0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xcd, 0x06, 0x0a, 0x03, 0x4a, 0x6f, 0x62, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x46, 0x0a, 0x08, 0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78,
	0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x50, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63,
	0x65, 0x52, 0x08, 0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x72,
	0x65, 0x66, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x72, 0x65, 0x66, 0x12, 0x19, 0x0a,
	0x08, 0x72, 0x65, 0x66, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x72, 0x65, 0x66, 0x50, 0x61, 0x74, 0x68, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x25, 0x0a, 0x0e, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x5f, 0x72, 0x65, 0x61, 0x73,
	0x6f, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72,
	0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x46, 0x0a, 0x0a, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x67, 0x69,
	0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4a, 0x6f, 0x62, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x73, 0x52, 0x0a, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x73, 0x12,
	0x42, 0x0a, 0x0f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x64, 0x5f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x0e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x64, 0x44, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x35, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x67, 0x65, 0x18,
	0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x61, 0x67, 0x73, 0x18, 0x0d, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73,
	0x12, 0x23, 0x0a, 0x0d, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x5f, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72,
	0x65, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x46, 0x61,
	0x69, 0x6c, 0x75, 0x72, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x61, 0x6e, 0x75, 0x61, 0x6c, 0x18,
	0x0f, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x6d, 0x61, 0x6e, 0x75, 0x61, 0x6c, 0x12, 0x18, 0x0a,
	0x07, 0x72, 0x65, 0x74, 0x72, 0x69, 0x65, 0x64, 0x18, 0x10, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07,
	0x72, 0x65, 0x74, 0x72, 0x69, 0x65, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65, 0x74, 0x72, 0x79,
	0x61, 0x62, 0x6c, 0x65, 0x18, 0x11, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x72, 0x65, 0x74, 0x72,
	0x79, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x34, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x12, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x20, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f,
	0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4a, 0x6f,
	0x62, 0x4b, 0x69, 0x6e, 0x64, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x60, 0x0a, 0x13, 0x64,
	0x6f, 0x77, 0x6e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x5f, 0x70, 0x69, 0x70, 0x65, 0x6c, 0x69,
	0x6e, 0x65, 0x18, 0x13, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61,
	0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x50, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x66, 0x65, 0x72,
	0x65, 0x6e, 0x63, 0x65, 0x48, 0x00, 0x52, 0x12, 0x64, 0x6f, 0x77, 0x6e, 0x73, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x50, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x88, 0x01, 0x01, 0x12, 0x40, 0x0a,
	0x06, 0x72, 0x75, 0x6e, 0x6e, 0x65, 0x72, 0x18, 0x14, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e,
	0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x52, 0x75, 0x6e, 0x6e, 0x65, 0x72, 0x52, 0x65,
	0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x06, 0x72, 0x75, 0x6e, 0x6e, 0x65, 0x72, 0x42,
	0x16, 0x0a, 0x14, 0x5f, 0x64, 0x6f, 0x77, 0x6e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x5f, 0x70,
	0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x22, 0xb4, 0x02, 0x0a, 0x0d, 0x4a, 0x6f, 0x62, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x73, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x12, 0x37, 0x0a, 0x09, 0x71, 0x75, 0x65, 0x75, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x08, 0x71, 0x75, 0x65, 0x75, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a,
	0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x3b, 0x0a, 0x0b, 0x66, 0x69, 0x6e, 0x69,
	0x73, 0x68, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x66, 0x69, 0x6e, 0x69, 0x73,
	0x68, 0x65, 0x64, 0x41, 0x74, 0x12, 0x37, 0x0a, 0x09, 0x65, 0x72, 0x61, 0x73, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x08, 0x65, 0x72, 0x61, 0x73, 0x65, 0x64, 0x41, 0x74, 0x2a, 0x49,
	0x0a, 0x07, 0x4a, 0x6f, 0x62, 0x4b, 0x69, 0x6e, 0x64, 0x12, 0x17, 0x0a, 0x13, 0x4a, 0x4f, 0x42,
	0x4b, 0x49, 0x4e, 0x44, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44,
	0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x4a, 0x4f, 0x42, 0x4b, 0x49, 0x4e, 0x44, 0x5f, 0x42, 0x55,
	0x49, 0x4c, 0x44, 0x10, 0x01, 0x12, 0x12, 0x0a, 0x0e, 0x4a, 0x4f, 0x42, 0x4b, 0x49, 0x4e, 0x44,
	0x5f, 0x42, 0x52, 0x49, 0x44, 0x47, 0x45, 0x10, 0x02, 0x42, 0x37, 0x5a, 0x35, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6c, 0x75, 0x74, 0x74, 0x72, 0x64, 0x65,
	0x76, 0x2f, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2d, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65,
	0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gitlabexporter_protobuf_job_proto_rawDescOnce sync.Once
	file_gitlabexporter_protobuf_job_proto_rawDescData = file_gitlabexporter_protobuf_job_proto_rawDesc
)

func file_gitlabexporter_protobuf_job_proto_rawDescGZIP() []byte {
	file_gitlabexporter_protobuf_job_proto_rawDescOnce.Do(func() {
		file_gitlabexporter_protobuf_job_proto_rawDescData = protoimpl.X.CompressGZIP(file_gitlabexporter_protobuf_job_proto_rawDescData)
	})
	return file_gitlabexporter_protobuf_job_proto_rawDescData
}

var file_gitlabexporter_protobuf_job_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_gitlabexporter_protobuf_job_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_gitlabexporter_protobuf_job_proto_goTypes = []interface{}{
	(JobKind)(0),                  // 0: gitlabexporter.protobuf.JobKind
	(*Job)(nil),                   // 1: gitlabexporter.protobuf.Job
	(*JobTimestamps)(nil),         // 2: gitlabexporter.protobuf.JobTimestamps
	(*PipelineReference)(nil),     // 3: gitlabexporter.protobuf.PipelineReference
	(*durationpb.Duration)(nil),   // 4: google.protobuf.Duration
	(*RunnerReference)(nil),       // 5: gitlabexporter.protobuf.RunnerReference
	(*timestamppb.Timestamp)(nil), // 6: google.protobuf.Timestamp
}
var file_gitlabexporter_protobuf_job_proto_depIdxs = []int32{
	3,  // 0: gitlabexporter.protobuf.Job.pipeline:type_name -> gitlabexporter.protobuf.PipelineReference
	2,  // 1: gitlabexporter.protobuf.Job.timestamps:type_name -> gitlabexporter.protobuf.JobTimestamps
	4,  // 2: gitlabexporter.protobuf.Job.queued_duration:type_name -> google.protobuf.Duration
	4,  // 3: gitlabexporter.protobuf.Job.duration:type_name -> google.protobuf.Duration
	0,  // 4: gitlabexporter.protobuf.Job.kind:type_name -> gitlabexporter.protobuf.JobKind
	3,  // 5: gitlabexporter.protobuf.Job.downstream_pipeline:type_name -> gitlabexporter.protobuf.PipelineReference
	5,  // 6: gitlabexporter.protobuf.Job.runner:type_name -> gitlabexporter.protobuf.RunnerReference
	6,  // 7: gitlabexporter.protobuf.JobTimestamps.created_at:type_name -> google.protobuf.Timestamp
	6,  // 8: gitlabexporter.protobuf.JobTimestamps.queued_at:type_name -> google.protobuf.Timestamp
	6,  // 9: gitlabexporter.protobuf.JobTimestamps.started_at:type_name -> google.protobuf.Timestamp
	6,  // 10: gitlabexporter.protobuf.JobTimestamps.finished_at:type_name -> google.protobuf.Timestamp
	6,  // 11: gitlabexporter.protobuf.JobTimestamps.erased_at:type_name -> google.protobuf.Timestamp
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_gitlabexporter_protobuf_job_proto_init() }
func file_gitlabexporter_protobuf_job_proto_init() {
	if File_gitlabexporter_protobuf_job_proto != nil {
		return
	}
	file_gitlabexporter_protobuf_references_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_gitlabexporter_protobuf_job_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Job); i {
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
		file_gitlabexporter_protobuf_job_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobTimestamps); i {
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
	file_gitlabexporter_protobuf_job_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_gitlabexporter_protobuf_job_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_gitlabexporter_protobuf_job_proto_goTypes,
		DependencyIndexes: file_gitlabexporter_protobuf_job_proto_depIdxs,
		EnumInfos:         file_gitlabexporter_protobuf_job_proto_enumTypes,
		MessageInfos:      file_gitlabexporter_protobuf_job_proto_msgTypes,
	}.Build()
	File_gitlabexporter_protobuf_job_proto = out.File
	file_gitlabexporter_protobuf_job_proto_rawDesc = nil
	file_gitlabexporter_protobuf_job_proto_goTypes = nil
	file_gitlabexporter_protobuf_job_proto_depIdxs = nil
}
