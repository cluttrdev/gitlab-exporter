// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.22.2
// source: gitlabexporter/protobuf/metric.proto

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

type Metric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        []byte                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Iid       int64                  `protobuf:"varint,2,opt,name=iid,proto3" json:"iid,omitempty"`
	JobId     int64                  `protobuf:"varint,3,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
	Name      string                 `protobuf:"bytes,10,opt,name=name,proto3" json:"name,omitempty"`
	Labels    []*Metric_Label        `protobuf:"bytes,11,rep,name=labels,proto3" json:"labels,omitempty"`
	Value     float64                `protobuf:"fixed64,12,opt,name=value,proto3" json:"value,omitempty"`
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,13,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *Metric) Reset() {
	*x = Metric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gitlabexporter_protobuf_metric_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_gitlabexporter_protobuf_metric_proto_msgTypes[0]
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
	return file_gitlabexporter_protobuf_metric_proto_rawDescGZIP(), []int{0}
}

func (x *Metric) GetId() []byte {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *Metric) GetIid() int64 {
	if x != nil {
		return x.Iid
	}
	return 0
}

func (x *Metric) GetJobId() int64 {
	if x != nil {
		return x.JobId
	}
	return 0
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
		mi := &file_gitlabexporter_protobuf_metric_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric_Label) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric_Label) ProtoMessage() {}

func (x *Metric_Label) ProtoReflect() protoreflect.Message {
	mi := &file_gitlabexporter_protobuf_metric_proto_msgTypes[1]
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
	return file_gitlabexporter_protobuf_metric_proto_rawDescGZIP(), []int{0, 0}
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

var File_gitlabexporter_protobuf_metric_proto protoreflect.FileDescriptor

var file_gitlabexporter_protobuf_metric_proto_rawDesc = []byte{
	0x0a, 0x24, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78,
	0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x1a,
	0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x97, 0x02, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x69,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x69, 0x69, 0x64, 0x12, 0x15, 0x0a,
	0x06, 0x6a, 0x6f, 0x62, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6a,
	0x6f, 0x62, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3d, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65,
	0x6c, 0x73, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61,
	0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x52,
	0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x38, 0x0a,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x1a, 0x31, 0x0a, 0x05, 0x4c, 0x61, 0x62, 0x65, 0x6c,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x37, 0x5a, 0x35, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6c, 0x75, 0x74, 0x74, 0x72, 0x64,
	0x65, 0x76, 0x2f, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2d, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74,
	0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x79, 0x70, 0x65,
	0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gitlabexporter_protobuf_metric_proto_rawDescOnce sync.Once
	file_gitlabexporter_protobuf_metric_proto_rawDescData = file_gitlabexporter_protobuf_metric_proto_rawDesc
)

func file_gitlabexporter_protobuf_metric_proto_rawDescGZIP() []byte {
	file_gitlabexporter_protobuf_metric_proto_rawDescOnce.Do(func() {
		file_gitlabexporter_protobuf_metric_proto_rawDescData = protoimpl.X.CompressGZIP(file_gitlabexporter_protobuf_metric_proto_rawDescData)
	})
	return file_gitlabexporter_protobuf_metric_proto_rawDescData
}

var file_gitlabexporter_protobuf_metric_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_gitlabexporter_protobuf_metric_proto_goTypes = []interface{}{
	(*Metric)(nil),                // 0: gitlabexporter.protobuf.Metric
	(*Metric_Label)(nil),          // 1: gitlabexporter.protobuf.Metric.Label
	(*timestamppb.Timestamp)(nil), // 2: google.protobuf.Timestamp
}
var file_gitlabexporter_protobuf_metric_proto_depIdxs = []int32{
	1, // 0: gitlabexporter.protobuf.Metric.labels:type_name -> gitlabexporter.protobuf.Metric.Label
	2, // 1: gitlabexporter.protobuf.Metric.timestamp:type_name -> google.protobuf.Timestamp
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_gitlabexporter_protobuf_metric_proto_init() }
func file_gitlabexporter_protobuf_metric_proto_init() {
	if File_gitlabexporter_protobuf_metric_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gitlabexporter_protobuf_metric_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_gitlabexporter_protobuf_metric_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_gitlabexporter_protobuf_metric_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_gitlabexporter_protobuf_metric_proto_goTypes,
		DependencyIndexes: file_gitlabexporter_protobuf_metric_proto_depIdxs,
		MessageInfos:      file_gitlabexporter_protobuf_metric_proto_msgTypes,
	}.Build()
	File_gitlabexporter_protobuf_metric_proto = out.File
	file_gitlabexporter_protobuf_metric_proto_rawDesc = nil
	file_gitlabexporter_protobuf_metric_proto_goTypes = nil
	file_gitlabexporter_protobuf_metric_proto_depIdxs = nil
}
