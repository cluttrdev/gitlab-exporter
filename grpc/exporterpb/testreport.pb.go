// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.22.2
// source: gitlabexporter/proto/models/testreport.proto

package exporterpb

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

type TestReport struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	PipelineId   int64   `protobuf:"varint,2,opt,name=pipeline_id,json=pipelineId,proto3" json:"pipeline_id,omitempty"`
	TotalTime    float64 `protobuf:"fixed64,3,opt,name=total_time,json=totalTime,proto3" json:"total_time,omitempty"`
	TotalCount   int64   `protobuf:"varint,4,opt,name=total_count,json=totalCount,proto3" json:"total_count,omitempty"`
	SuccessCount int64   `protobuf:"varint,5,opt,name=success_count,json=successCount,proto3" json:"success_count,omitempty"`
	FailedCount  int64   `protobuf:"varint,6,opt,name=failed_count,json=failedCount,proto3" json:"failed_count,omitempty"`
	SkippedCount int64   `protobuf:"varint,7,opt,name=skipped_count,json=skippedCount,proto3" json:"skipped_count,omitempty"`
	ErrorCount   int64   `protobuf:"varint,8,opt,name=error_count,json=errorCount,proto3" json:"error_count,omitempty"`
}

func (x *TestReport) Reset() {
	*x = TestReport{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gitlabexporter_proto_models_testreport_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestReport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestReport) ProtoMessage() {}

func (x *TestReport) ProtoReflect() protoreflect.Message {
	mi := &file_gitlabexporter_proto_models_testreport_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestReport.ProtoReflect.Descriptor instead.
func (*TestReport) Descriptor() ([]byte, []int) {
	return file_gitlabexporter_proto_models_testreport_proto_rawDescGZIP(), []int{0}
}

func (x *TestReport) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *TestReport) GetPipelineId() int64 {
	if x != nil {
		return x.PipelineId
	}
	return 0
}

func (x *TestReport) GetTotalTime() float64 {
	if x != nil {
		return x.TotalTime
	}
	return 0
}

func (x *TestReport) GetTotalCount() int64 {
	if x != nil {
		return x.TotalCount
	}
	return 0
}

func (x *TestReport) GetSuccessCount() int64 {
	if x != nil {
		return x.SuccessCount
	}
	return 0
}

func (x *TestReport) GetFailedCount() int64 {
	if x != nil {
		return x.FailedCount
	}
	return 0
}

func (x *TestReport) GetSkippedCount() int64 {
	if x != nil {
		return x.SkippedCount
	}
	return 0
}

func (x *TestReport) GetErrorCount() int64 {
	if x != nil {
		return x.ErrorCount
	}
	return 0
}

type TestSuite struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	TestreportId string  `protobuf:"bytes,2,opt,name=testreport_id,json=testreportId,proto3" json:"testreport_id,omitempty"`
	PipelineId   int64   `protobuf:"varint,3,opt,name=pipeline_id,json=pipelineId,proto3" json:"pipeline_id,omitempty"`
	Name         string  `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	TotalTime    float64 `protobuf:"fixed64,5,opt,name=total_time,json=totalTime,proto3" json:"total_time,omitempty"`
	TotalCount   int64   `protobuf:"varint,6,opt,name=total_count,json=totalCount,proto3" json:"total_count,omitempty"`
	SuccessCount int64   `protobuf:"varint,7,opt,name=success_count,json=successCount,proto3" json:"success_count,omitempty"`
	FailedCount  int64   `protobuf:"varint,8,opt,name=failed_count,json=failedCount,proto3" json:"failed_count,omitempty"`
	SkippedCount int64   `protobuf:"varint,9,opt,name=skipped_count,json=skippedCount,proto3" json:"skipped_count,omitempty"`
	ErrorCount   int64   `protobuf:"varint,10,opt,name=error_count,json=errorCount,proto3" json:"error_count,omitempty"`
}

func (x *TestSuite) Reset() {
	*x = TestSuite{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gitlabexporter_proto_models_testreport_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestSuite) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestSuite) ProtoMessage() {}

func (x *TestSuite) ProtoReflect() protoreflect.Message {
	mi := &file_gitlabexporter_proto_models_testreport_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestSuite.ProtoReflect.Descriptor instead.
func (*TestSuite) Descriptor() ([]byte, []int) {
	return file_gitlabexporter_proto_models_testreport_proto_rawDescGZIP(), []int{1}
}

func (x *TestSuite) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *TestSuite) GetTestreportId() string {
	if x != nil {
		return x.TestreportId
	}
	return ""
}

func (x *TestSuite) GetPipelineId() int64 {
	if x != nil {
		return x.PipelineId
	}
	return 0
}

func (x *TestSuite) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *TestSuite) GetTotalTime() float64 {
	if x != nil {
		return x.TotalTime
	}
	return 0
}

func (x *TestSuite) GetTotalCount() int64 {
	if x != nil {
		return x.TotalCount
	}
	return 0
}

func (x *TestSuite) GetSuccessCount() int64 {
	if x != nil {
		return x.SuccessCount
	}
	return 0
}

func (x *TestSuite) GetFailedCount() int64 {
	if x != nil {
		return x.FailedCount
	}
	return 0
}

func (x *TestSuite) GetSkippedCount() int64 {
	if x != nil {
		return x.SkippedCount
	}
	return 0
}

func (x *TestSuite) GetErrorCount() int64 {
	if x != nil {
		return x.ErrorCount
	}
	return 0
}

type TestCase struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string                   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	TestsuiteId    string                   `protobuf:"bytes,2,opt,name=testsuite_id,json=testsuiteId,proto3" json:"testsuite_id,omitempty"`
	TestreportId   string                   `protobuf:"bytes,3,opt,name=testreport_id,json=testreportId,proto3" json:"testreport_id,omitempty"`
	PipelineId     int64                    `protobuf:"varint,4,opt,name=pipeline_id,json=pipelineId,proto3" json:"pipeline_id,omitempty"`
	Status         string                   `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
	Name           string                   `protobuf:"bytes,6,opt,name=name,proto3" json:"name,omitempty"`
	Classname      string                   `protobuf:"bytes,7,opt,name=classname,proto3" json:"classname,omitempty"`
	File           string                   `protobuf:"bytes,8,opt,name=file,proto3" json:"file,omitempty"`
	ExecutionTime  float64                  `protobuf:"fixed64,9,opt,name=execution_time,json=executionTime,proto3" json:"execution_time,omitempty"`
	SystemOutput   string                   `protobuf:"bytes,10,opt,name=system_output,json=systemOutput,proto3" json:"system_output,omitempty"`
	StackTrace     string                   `protobuf:"bytes,11,opt,name=stack_trace,json=stackTrace,proto3" json:"stack_trace,omitempty"`
	AttachmentUrl  string                   `protobuf:"bytes,12,opt,name=attachment_url,json=attachmentUrl,proto3" json:"attachment_url,omitempty"`
	RecentFailures *TestCase_RecentFailures `protobuf:"bytes,13,opt,name=recent_failures,json=recentFailures,proto3" json:"recent_failures,omitempty"`
}

func (x *TestCase) Reset() {
	*x = TestCase{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gitlabexporter_proto_models_testreport_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestCase) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestCase) ProtoMessage() {}

func (x *TestCase) ProtoReflect() protoreflect.Message {
	mi := &file_gitlabexporter_proto_models_testreport_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestCase.ProtoReflect.Descriptor instead.
func (*TestCase) Descriptor() ([]byte, []int) {
	return file_gitlabexporter_proto_models_testreport_proto_rawDescGZIP(), []int{2}
}

func (x *TestCase) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *TestCase) GetTestsuiteId() string {
	if x != nil {
		return x.TestsuiteId
	}
	return ""
}

func (x *TestCase) GetTestreportId() string {
	if x != nil {
		return x.TestreportId
	}
	return ""
}

func (x *TestCase) GetPipelineId() int64 {
	if x != nil {
		return x.PipelineId
	}
	return 0
}

func (x *TestCase) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *TestCase) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *TestCase) GetClassname() string {
	if x != nil {
		return x.Classname
	}
	return ""
}

func (x *TestCase) GetFile() string {
	if x != nil {
		return x.File
	}
	return ""
}

func (x *TestCase) GetExecutionTime() float64 {
	if x != nil {
		return x.ExecutionTime
	}
	return 0
}

func (x *TestCase) GetSystemOutput() string {
	if x != nil {
		return x.SystemOutput
	}
	return ""
}

func (x *TestCase) GetStackTrace() string {
	if x != nil {
		return x.StackTrace
	}
	return ""
}

func (x *TestCase) GetAttachmentUrl() string {
	if x != nil {
		return x.AttachmentUrl
	}
	return ""
}

func (x *TestCase) GetRecentFailures() *TestCase_RecentFailures {
	if x != nil {
		return x.RecentFailures
	}
	return nil
}

type TestCase_RecentFailures struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count      int64  `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	BaseBranch string `protobuf:"bytes,2,opt,name=base_branch,json=baseBranch,proto3" json:"base_branch,omitempty"`
}

func (x *TestCase_RecentFailures) Reset() {
	*x = TestCase_RecentFailures{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gitlabexporter_proto_models_testreport_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestCase_RecentFailures) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestCase_RecentFailures) ProtoMessage() {}

func (x *TestCase_RecentFailures) ProtoReflect() protoreflect.Message {
	mi := &file_gitlabexporter_proto_models_testreport_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestCase_RecentFailures.ProtoReflect.Descriptor instead.
func (*TestCase_RecentFailures) Descriptor() ([]byte, []int) {
	return file_gitlabexporter_proto_models_testreport_proto_rawDescGZIP(), []int{2, 0}
}

func (x *TestCase_RecentFailures) GetCount() int64 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *TestCase_RecentFailures) GetBaseBranch() string {
	if x != nil {
		return x.BaseBranch
	}
	return ""
}

var File_gitlabexporter_proto_models_testreport_proto protoreflect.FileDescriptor

var file_gitlabexporter_proto_models_testreport_proto_rawDesc = []byte{
	0x0a, 0x2c, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x74, 0x65,
	0x73, 0x74, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1b,
	0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x22, 0x8b, 0x02, 0x0a, 0x0a,
	0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x69,
	0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0a, 0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x74,
	0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x09, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x6f,
	0x74, 0x61, 0x6c, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x73,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0c, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x12, 0x21, 0x0a, 0x0c, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x6b, 0x69, 0x70, 0x70, 0x65, 0x64, 0x5f, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x73, 0x6b, 0x69, 0x70,
	0x70, 0x65, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0xc3, 0x02, 0x0a, 0x09, 0x54, 0x65,
	0x73, 0x74, 0x53, 0x75, 0x69, 0x74, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x74, 0x65, 0x73, 0x74, 0x72,
	0x65, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x74, 0x65, 0x73, 0x74, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b,
	0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x0a, 0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x49, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64,
	0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x66, 0x61,
	0x69, 0x6c, 0x65, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x6b, 0x69,
	0x70, 0x70, 0x65, 0x64, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x0c, 0x73, 0x6b, 0x69, 0x70, 0x70, 0x65, 0x64, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1f,
	0x0a, 0x0b, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0a, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22,
	0x9d, 0x04, 0x0a, 0x08, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x21, 0x0a, 0x0c,
	0x74, 0x65, 0x73, 0x74, 0x73, 0x75, 0x69, 0x74, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x74, 0x65, 0x73, 0x74, 0x73, 0x75, 0x69, 0x74, 0x65, 0x49, 0x64, 0x12,
	0x23, 0x0a, 0x0d, 0x74, 0x65, 0x73, 0x74, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x74, 0x65, 0x73, 0x74, 0x72, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x70, 0x69, 0x70, 0x65, 0x6c,
	0x69, 0x6e, 0x65, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66,
	0x69, 0x6c, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0d, 0x65, 0x78, 0x65,
	0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x79,
	0x73, 0x74, 0x65, 0x6d, 0x5f, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12,
	0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x5f, 0x74, 0x72, 0x61, 0x63, 0x65, 0x18, 0x0b,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x54, 0x72, 0x61, 0x63, 0x65,
	0x12, 0x25, 0x0a, 0x0e, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x75,
	0x72, 0x6c, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68,
	0x6d, 0x65, 0x6e, 0x74, 0x55, 0x72, 0x6c, 0x12, 0x5d, 0x0a, 0x0f, 0x72, 0x65, 0x63, 0x65, 0x6e,
	0x74, 0x5f, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x73, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x34, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x54,
	0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x2e, 0x52, 0x65, 0x63, 0x65, 0x6e, 0x74, 0x46, 0x61,
	0x69, 0x6c, 0x75, 0x72, 0x65, 0x73, 0x52, 0x0e, 0x72, 0x65, 0x63, 0x65, 0x6e, 0x74, 0x46, 0x61,
	0x69, 0x6c, 0x75, 0x72, 0x65, 0x73, 0x1a, 0x47, 0x0a, 0x0e, 0x52, 0x65, 0x63, 0x65, 0x6e, 0x74,
	0x46, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1f,
	0x0a, 0x0b, 0x62, 0x61, 0x73, 0x65, 0x5f, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x62, 0x61, 0x73, 0x65, 0x42, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x42,
	0x36, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6c,
	0x75, 0x74, 0x74, 0x72, 0x64, 0x65, 0x76, 0x2f, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2d, 0x65,
	0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x65, 0x78, 0x70,
	0x6f, 0x72, 0x74, 0x65, 0x72, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gitlabexporter_proto_models_testreport_proto_rawDescOnce sync.Once
	file_gitlabexporter_proto_models_testreport_proto_rawDescData = file_gitlabexporter_proto_models_testreport_proto_rawDesc
)

func file_gitlabexporter_proto_models_testreport_proto_rawDescGZIP() []byte {
	file_gitlabexporter_proto_models_testreport_proto_rawDescOnce.Do(func() {
		file_gitlabexporter_proto_models_testreport_proto_rawDescData = protoimpl.X.CompressGZIP(file_gitlabexporter_proto_models_testreport_proto_rawDescData)
	})
	return file_gitlabexporter_proto_models_testreport_proto_rawDescData
}

var file_gitlabexporter_proto_models_testreport_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_gitlabexporter_proto_models_testreport_proto_goTypes = []interface{}{
	(*TestReport)(nil),              // 0: gitlabexporter.proto.models.TestReport
	(*TestSuite)(nil),               // 1: gitlabexporter.proto.models.TestSuite
	(*TestCase)(nil),                // 2: gitlabexporter.proto.models.TestCase
	(*TestCase_RecentFailures)(nil), // 3: gitlabexporter.proto.models.TestCase.RecentFailures
}
var file_gitlabexporter_proto_models_testreport_proto_depIdxs = []int32{
	3, // 0: gitlabexporter.proto.models.TestCase.recent_failures:type_name -> gitlabexporter.proto.models.TestCase.RecentFailures
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_gitlabexporter_proto_models_testreport_proto_init() }
func file_gitlabexporter_proto_models_testreport_proto_init() {
	if File_gitlabexporter_proto_models_testreport_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gitlabexporter_proto_models_testreport_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestReport); i {
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
		file_gitlabexporter_proto_models_testreport_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestSuite); i {
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
		file_gitlabexporter_proto_models_testreport_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestCase); i {
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
		file_gitlabexporter_proto_models_testreport_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestCase_RecentFailures); i {
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
			RawDescriptor: file_gitlabexporter_proto_models_testreport_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_gitlabexporter_proto_models_testreport_proto_goTypes,
		DependencyIndexes: file_gitlabexporter_proto_models_testreport_proto_depIdxs,
		MessageInfos:      file_gitlabexporter_proto_models_testreport_proto_msgTypes,
	}.Build()
	File_gitlabexporter_proto_models_testreport_proto = out.File
	file_gitlabexporter_proto_models_testreport_proto_rawDesc = nil
	file_gitlabexporter_proto_models_testreport_proto_goTypes = nil
	file_gitlabexporter_proto_models_testreport_proto_depIdxs = nil
}
