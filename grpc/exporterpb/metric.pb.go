// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.22.2
// source: gitlabexporter/proto/models/metric.proto

package exporterpb

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

type Metric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Labels    []*Metric_Label        `protobuf:"bytes,2,rep,name=labels,proto3" json:"labels,omitempty"`
	Value     float64                `protobuf:"fixed64,3,opt,name=value,proto3" json:"value,omitempty"`
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Job       *Metric_JobReference   `protobuf:"bytes,5,opt,name=job,proto3" json:"job,omitempty"`
}

func (x *Metric) Reset() {
	*x = Metric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gitlabexporter_proto_models_metric_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_gitlabexporter_proto_models_metric_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metric.ProtoReflect.Descriptor instead.
func (*Metric) Descriptor() ([]byte, []int) {
	return file_gitlabexporter_proto_models_metric_proto_rawDescGZIP(), []int{0}
}

func (x *Metric) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Metric) GetLabels() []*Metric_Label {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *Metric) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *Metric) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *Metric) GetJob() *Metric_JobReference {
	if x != nil {
		return x.Job
	}
	return nil
}

type Metric_Label struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name  string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Metric_Label) Reset() {
	*x = Metric_Label{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gitlabexporter_proto_models_metric_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric_Label) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric_Label) ProtoMessage() {}

func (x *Metric_Label) ProtoReflect() protoreflect.Message {
	mi := &file_gitlabexporter_proto_models_metric_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metric_Label.ProtoReflect.Descriptor instead.
func (*Metric_Label) Descriptor() ([]byte, []int) {
	return file_gitlabexporter_proto_models_metric_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Metric_Label) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Metric_Label) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type Metric_JobReference struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Metric_JobReference) Reset() {
	*x = Metric_JobReference{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gitlabexporter_proto_models_metric_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric_JobReference) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric_JobReference) ProtoMessage() {}

func (x *Metric_JobReference) ProtoReflect() protoreflect.Message {
	mi := &file_gitlabexporter_proto_models_metric_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metric_JobReference.ProtoReflect.Descriptor instead.
func (*Metric_JobReference) Descriptor() ([]byte, []int) {
	return file_gitlabexporter_proto_models_metric_proto_rawDescGZIP(), []int{0, 1}
}

func (x *Metric_JobReference) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Metric_JobReference) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_gitlabexporter_proto_models_metric_proto protoreflect.FileDescriptor

var file_gitlabexporter_proto_models_metric_proto_rawDesc = []byte{
	0x0a, 0x28, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x6d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1b, 0x67, 0x69, 0x74, 0x6c,
	0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xda, 0x02, 0x0a, 0x06, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x41, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62,
	0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x4c, 0x61, 0x62,
	0x65, 0x6c, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x42, 0x0a, 0x03, 0x6a, 0x6f,
	0x62, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62,
	0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x4a, 0x6f, 0x62,
	0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x03, 0x6a, 0x6f, 0x62, 0x1a, 0x31,
	0x0a, 0x05, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x1a, 0x32, 0x0a, 0x0c, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x36, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6c, 0x75, 0x74, 0x74, 0x72, 0x64, 0x65, 0x76, 0x2f, 0x67, 0x69,
	0x74, 0x6c, 0x61, 0x62, 0x2d, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2f, 0x67, 0x72,
	0x70, 0x63, 0x2f, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x70, 0x62, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gitlabexporter_proto_models_metric_proto_rawDescOnce sync.Once
	file_gitlabexporter_proto_models_metric_proto_rawDescData = file_gitlabexporter_proto_models_metric_proto_rawDesc
)

func file_gitlabexporter_proto_models_metric_proto_rawDescGZIP() []byte {
	file_gitlabexporter_proto_models_metric_proto_rawDescOnce.Do(func() {
		file_gitlabexporter_proto_models_metric_proto_rawDescData = protoimpl.X.CompressGZIP(file_gitlabexporter_proto_models_metric_proto_rawDescData)
	})
	return file_gitlabexporter_proto_models_metric_proto_rawDescData
}

var file_gitlabexporter_proto_models_metric_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_gitlabexporter_proto_models_metric_proto_goTypes = []interface{}{
	(*Metric)(nil),                // 0: gitlabexporter.proto.models.Metric
	(*Metric_Label)(nil),          // 1: gitlabexporter.proto.models.Metric.Label
	(*Metric_JobReference)(nil),   // 2: gitlabexporter.proto.models.Metric.JobReference
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
}
var file_gitlabexporter_proto_models_metric_proto_depIdxs = []int32{
	1, // 0: gitlabexporter.proto.models.Metric.labels:type_name -> gitlabexporter.proto.models.Metric.Label
	3, // 1: gitlabexporter.proto.models.Metric.timestamp:type_name -> google.protobuf.Timestamp
	2, // 2: gitlabexporter.proto.models.Metric.job:type_name -> gitlabexporter.proto.models.Metric.JobReference
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_gitlabexporter_proto_models_metric_proto_init() }
func file_gitlabexporter_proto_models_metric_proto_init() {
	if File_gitlabexporter_proto_models_metric_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gitlabexporter_proto_models_metric_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metric); i {
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
		file_gitlabexporter_proto_models_metric_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metric_Label); i {
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
		file_gitlabexporter_proto_models_metric_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metric_JobReference); i {
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
			RawDescriptor: file_gitlabexporter_proto_models_metric_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_gitlabexporter_proto_models_metric_proto_goTypes,
		DependencyIndexes: file_gitlabexporter_proto_models_metric_proto_depIdxs,
		MessageInfos:      file_gitlabexporter_proto_models_metric_proto_msgTypes,
	}.Build()
	File_gitlabexporter_proto_models_metric_proto = out.File
	file_gitlabexporter_proto_models_metric_proto_rawDesc = nil
	file_gitlabexporter_proto_models_metric_proto_goTypes = nil
	file_gitlabexporter_proto_models_metric_proto_depIdxs = nil
}
