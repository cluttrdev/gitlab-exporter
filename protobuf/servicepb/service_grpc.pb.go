// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.22.2
// source: gitlabexporter/protobuf/service/service.proto

package servicepb

import (
	context "context"
	typespb "github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	GitLabExporter_RecordPipelines_FullMethodName   = "/gitlabexporter.protobuf.service.GitLabExporter/RecordPipelines"
	GitLabExporter_RecordJobs_FullMethodName        = "/gitlabexporter.protobuf.service.GitLabExporter/RecordJobs"
	GitLabExporter_RecordSections_FullMethodName    = "/gitlabexporter.protobuf.service.GitLabExporter/RecordSections"
	GitLabExporter_RecordBridges_FullMethodName     = "/gitlabexporter.protobuf.service.GitLabExporter/RecordBridges"
	GitLabExporter_RecordTestReports_FullMethodName = "/gitlabexporter.protobuf.service.GitLabExporter/RecordTestReports"
	GitLabExporter_RecordTestSuites_FullMethodName  = "/gitlabexporter.protobuf.service.GitLabExporter/RecordTestSuites"
	GitLabExporter_RecordTestCases_FullMethodName   = "/gitlabexporter.protobuf.service.GitLabExporter/RecordTestCases"
	GitLabExporter_RecordMetrics_FullMethodName     = "/gitlabexporter.protobuf.service.GitLabExporter/RecordMetrics"
	GitLabExporter_RecordTraces_FullMethodName      = "/gitlabexporter.protobuf.service.GitLabExporter/RecordTraces"
)

// GitLabExporterClient is the client API for GitLabExporter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GitLabExporterClient interface {
	RecordPipelines(ctx context.Context, opts ...grpc.CallOption) (GitLabExporter_RecordPipelinesClient, error)
	RecordJobs(ctx context.Context, opts ...grpc.CallOption) (GitLabExporter_RecordJobsClient, error)
	RecordSections(ctx context.Context, opts ...grpc.CallOption) (GitLabExporter_RecordSectionsClient, error)
	RecordBridges(ctx context.Context, opts ...grpc.CallOption) (GitLabExporter_RecordBridgesClient, error)
	RecordTestReports(ctx context.Context, opts ...grpc.CallOption) (GitLabExporter_RecordTestReportsClient, error)
	RecordTestSuites(ctx context.Context, opts ...grpc.CallOption) (GitLabExporter_RecordTestSuitesClient, error)
	RecordTestCases(ctx context.Context, opts ...grpc.CallOption) (GitLabExporter_RecordTestCasesClient, error)
	RecordMetrics(ctx context.Context, opts ...grpc.CallOption) (GitLabExporter_RecordMetricsClient, error)
	RecordTraces(ctx context.Context, opts ...grpc.CallOption) (GitLabExporter_RecordTracesClient, error)
}

type gitLabExporterClient struct {
	cc grpc.ClientConnInterface
}

func NewGitLabExporterClient(cc grpc.ClientConnInterface) GitLabExporterClient {
	return &gitLabExporterClient{cc}
}

func (c *gitLabExporterClient) RecordPipelines(ctx context.Context, opts ...grpc.CallOption) (GitLabExporter_RecordPipelinesClient, error) {
	stream, err := c.cc.NewStream(ctx, &GitLabExporter_ServiceDesc.Streams[0], GitLabExporter_RecordPipelines_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &gitLabExporterRecordPipelinesClient{stream}
	return x, nil
}

type GitLabExporter_RecordPipelinesClient interface {
	Send(*typespb.Pipeline) error
	CloseAndRecv() (*RecordSummary, error)
	grpc.ClientStream
}

type gitLabExporterRecordPipelinesClient struct {
	grpc.ClientStream
}

func (x *gitLabExporterRecordPipelinesClient) Send(m *typespb.Pipeline) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gitLabExporterRecordPipelinesClient) CloseAndRecv() (*RecordSummary, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(RecordSummary)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gitLabExporterClient) RecordJobs(ctx context.Context, opts ...grpc.CallOption) (GitLabExporter_RecordJobsClient, error) {
	stream, err := c.cc.NewStream(ctx, &GitLabExporter_ServiceDesc.Streams[1], GitLabExporter_RecordJobs_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &gitLabExporterRecordJobsClient{stream}
	return x, nil
}

type GitLabExporter_RecordJobsClient interface {
	Send(*typespb.Job) error
	CloseAndRecv() (*RecordSummary, error)
	grpc.ClientStream
}

type gitLabExporterRecordJobsClient struct {
	grpc.ClientStream
}

func (x *gitLabExporterRecordJobsClient) Send(m *typespb.Job) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gitLabExporterRecordJobsClient) CloseAndRecv() (*RecordSummary, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(RecordSummary)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gitLabExporterClient) RecordSections(ctx context.Context, opts ...grpc.CallOption) (GitLabExporter_RecordSectionsClient, error) {
	stream, err := c.cc.NewStream(ctx, &GitLabExporter_ServiceDesc.Streams[2], GitLabExporter_RecordSections_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &gitLabExporterRecordSectionsClient{stream}
	return x, nil
}

type GitLabExporter_RecordSectionsClient interface {
	Send(*typespb.Section) error
	CloseAndRecv() (*RecordSummary, error)
	grpc.ClientStream
}

type gitLabExporterRecordSectionsClient struct {
	grpc.ClientStream
}

func (x *gitLabExporterRecordSectionsClient) Send(m *typespb.Section) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gitLabExporterRecordSectionsClient) CloseAndRecv() (*RecordSummary, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(RecordSummary)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gitLabExporterClient) RecordBridges(ctx context.Context, opts ...grpc.CallOption) (GitLabExporter_RecordBridgesClient, error) {
	stream, err := c.cc.NewStream(ctx, &GitLabExporter_ServiceDesc.Streams[3], GitLabExporter_RecordBridges_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &gitLabExporterRecordBridgesClient{stream}
	return x, nil
}

type GitLabExporter_RecordBridgesClient interface {
	Send(*typespb.Bridge) error
	CloseAndRecv() (*RecordSummary, error)
	grpc.ClientStream
}

type gitLabExporterRecordBridgesClient struct {
	grpc.ClientStream
}

func (x *gitLabExporterRecordBridgesClient) Send(m *typespb.Bridge) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gitLabExporterRecordBridgesClient) CloseAndRecv() (*RecordSummary, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(RecordSummary)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gitLabExporterClient) RecordTestReports(ctx context.Context, opts ...grpc.CallOption) (GitLabExporter_RecordTestReportsClient, error) {
	stream, err := c.cc.NewStream(ctx, &GitLabExporter_ServiceDesc.Streams[4], GitLabExporter_RecordTestReports_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &gitLabExporterRecordTestReportsClient{stream}
	return x, nil
}

type GitLabExporter_RecordTestReportsClient interface {
	Send(*typespb.TestReport) error
	CloseAndRecv() (*RecordSummary, error)
	grpc.ClientStream
}

type gitLabExporterRecordTestReportsClient struct {
	grpc.ClientStream
}

func (x *gitLabExporterRecordTestReportsClient) Send(m *typespb.TestReport) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gitLabExporterRecordTestReportsClient) CloseAndRecv() (*RecordSummary, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(RecordSummary)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gitLabExporterClient) RecordTestSuites(ctx context.Context, opts ...grpc.CallOption) (GitLabExporter_RecordTestSuitesClient, error) {
	stream, err := c.cc.NewStream(ctx, &GitLabExporter_ServiceDesc.Streams[5], GitLabExporter_RecordTestSuites_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &gitLabExporterRecordTestSuitesClient{stream}
	return x, nil
}

type GitLabExporter_RecordTestSuitesClient interface {
	Send(*typespb.TestSuite) error
	CloseAndRecv() (*RecordSummary, error)
	grpc.ClientStream
}

type gitLabExporterRecordTestSuitesClient struct {
	grpc.ClientStream
}

func (x *gitLabExporterRecordTestSuitesClient) Send(m *typespb.TestSuite) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gitLabExporterRecordTestSuitesClient) CloseAndRecv() (*RecordSummary, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(RecordSummary)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gitLabExporterClient) RecordTestCases(ctx context.Context, opts ...grpc.CallOption) (GitLabExporter_RecordTestCasesClient, error) {
	stream, err := c.cc.NewStream(ctx, &GitLabExporter_ServiceDesc.Streams[6], GitLabExporter_RecordTestCases_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &gitLabExporterRecordTestCasesClient{stream}
	return x, nil
}

type GitLabExporter_RecordTestCasesClient interface {
	Send(*typespb.TestCase) error
	CloseAndRecv() (*RecordSummary, error)
	grpc.ClientStream
}

type gitLabExporterRecordTestCasesClient struct {
	grpc.ClientStream
}

func (x *gitLabExporterRecordTestCasesClient) Send(m *typespb.TestCase) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gitLabExporterRecordTestCasesClient) CloseAndRecv() (*RecordSummary, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(RecordSummary)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gitLabExporterClient) RecordMetrics(ctx context.Context, opts ...grpc.CallOption) (GitLabExporter_RecordMetricsClient, error) {
	stream, err := c.cc.NewStream(ctx, &GitLabExporter_ServiceDesc.Streams[7], GitLabExporter_RecordMetrics_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &gitLabExporterRecordMetricsClient{stream}
	return x, nil
}

type GitLabExporter_RecordMetricsClient interface {
	Send(*typespb.Metric) error
	CloseAndRecv() (*RecordSummary, error)
	grpc.ClientStream
}

type gitLabExporterRecordMetricsClient struct {
	grpc.ClientStream
}

func (x *gitLabExporterRecordMetricsClient) Send(m *typespb.Metric) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gitLabExporterRecordMetricsClient) CloseAndRecv() (*RecordSummary, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(RecordSummary)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gitLabExporterClient) RecordTraces(ctx context.Context, opts ...grpc.CallOption) (GitLabExporter_RecordTracesClient, error) {
	stream, err := c.cc.NewStream(ctx, &GitLabExporter_ServiceDesc.Streams[8], GitLabExporter_RecordTraces_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &gitLabExporterRecordTracesClient{stream}
	return x, nil
}

type GitLabExporter_RecordTracesClient interface {
	Send(*typespb.Trace) error
	CloseAndRecv() (*RecordSummary, error)
	grpc.ClientStream
}

type gitLabExporterRecordTracesClient struct {
	grpc.ClientStream
}

func (x *gitLabExporterRecordTracesClient) Send(m *typespb.Trace) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gitLabExporterRecordTracesClient) CloseAndRecv() (*RecordSummary, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(RecordSummary)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GitLabExporterServer is the server API for GitLabExporter service.
// All implementations must embed UnimplementedGitLabExporterServer
// for forward compatibility
type GitLabExporterServer interface {
	RecordPipelines(GitLabExporter_RecordPipelinesServer) error
	RecordJobs(GitLabExporter_RecordJobsServer) error
	RecordSections(GitLabExporter_RecordSectionsServer) error
	RecordBridges(GitLabExporter_RecordBridgesServer) error
	RecordTestReports(GitLabExporter_RecordTestReportsServer) error
	RecordTestSuites(GitLabExporter_RecordTestSuitesServer) error
	RecordTestCases(GitLabExporter_RecordTestCasesServer) error
	RecordMetrics(GitLabExporter_RecordMetricsServer) error
	RecordTraces(GitLabExporter_RecordTracesServer) error
	mustEmbedUnimplementedGitLabExporterServer()
}

// UnimplementedGitLabExporterServer must be embedded to have forward compatible implementations.
type UnimplementedGitLabExporterServer struct {
}

func (UnimplementedGitLabExporterServer) RecordPipelines(GitLabExporter_RecordPipelinesServer) error {
	return status.Errorf(codes.Unimplemented, "method RecordPipelines not implemented")
}
func (UnimplementedGitLabExporterServer) RecordJobs(GitLabExporter_RecordJobsServer) error {
	return status.Errorf(codes.Unimplemented, "method RecordJobs not implemented")
}
func (UnimplementedGitLabExporterServer) RecordSections(GitLabExporter_RecordSectionsServer) error {
	return status.Errorf(codes.Unimplemented, "method RecordSections not implemented")
}
func (UnimplementedGitLabExporterServer) RecordBridges(GitLabExporter_RecordBridgesServer) error {
	return status.Errorf(codes.Unimplemented, "method RecordBridges not implemented")
}
func (UnimplementedGitLabExporterServer) RecordTestReports(GitLabExporter_RecordTestReportsServer) error {
	return status.Errorf(codes.Unimplemented, "method RecordTestReports not implemented")
}
func (UnimplementedGitLabExporterServer) RecordTestSuites(GitLabExporter_RecordTestSuitesServer) error {
	return status.Errorf(codes.Unimplemented, "method RecordTestSuites not implemented")
}
func (UnimplementedGitLabExporterServer) RecordTestCases(GitLabExporter_RecordTestCasesServer) error {
	return status.Errorf(codes.Unimplemented, "method RecordTestCases not implemented")
}
func (UnimplementedGitLabExporterServer) RecordMetrics(GitLabExporter_RecordMetricsServer) error {
	return status.Errorf(codes.Unimplemented, "method RecordMetrics not implemented")
}
func (UnimplementedGitLabExporterServer) RecordTraces(GitLabExporter_RecordTracesServer) error {
	return status.Errorf(codes.Unimplemented, "method RecordTraces not implemented")
}
func (UnimplementedGitLabExporterServer) mustEmbedUnimplementedGitLabExporterServer() {}

// UnsafeGitLabExporterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GitLabExporterServer will
// result in compilation errors.
type UnsafeGitLabExporterServer interface {
	mustEmbedUnimplementedGitLabExporterServer()
}

func RegisterGitLabExporterServer(s grpc.ServiceRegistrar, srv GitLabExporterServer) {
	s.RegisterService(&GitLabExporter_ServiceDesc, srv)
}

func _GitLabExporter_RecordPipelines_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GitLabExporterServer).RecordPipelines(&gitLabExporterRecordPipelinesServer{stream})
}

type GitLabExporter_RecordPipelinesServer interface {
	SendAndClose(*RecordSummary) error
	Recv() (*typespb.Pipeline, error)
	grpc.ServerStream
}

type gitLabExporterRecordPipelinesServer struct {
	grpc.ServerStream
}

func (x *gitLabExporterRecordPipelinesServer) SendAndClose(m *RecordSummary) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gitLabExporterRecordPipelinesServer) Recv() (*typespb.Pipeline, error) {
	m := new(typespb.Pipeline)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _GitLabExporter_RecordJobs_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GitLabExporterServer).RecordJobs(&gitLabExporterRecordJobsServer{stream})
}

type GitLabExporter_RecordJobsServer interface {
	SendAndClose(*RecordSummary) error
	Recv() (*typespb.Job, error)
	grpc.ServerStream
}

type gitLabExporterRecordJobsServer struct {
	grpc.ServerStream
}

func (x *gitLabExporterRecordJobsServer) SendAndClose(m *RecordSummary) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gitLabExporterRecordJobsServer) Recv() (*typespb.Job, error) {
	m := new(typespb.Job)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _GitLabExporter_RecordSections_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GitLabExporterServer).RecordSections(&gitLabExporterRecordSectionsServer{stream})
}

type GitLabExporter_RecordSectionsServer interface {
	SendAndClose(*RecordSummary) error
	Recv() (*typespb.Section, error)
	grpc.ServerStream
}

type gitLabExporterRecordSectionsServer struct {
	grpc.ServerStream
}

func (x *gitLabExporterRecordSectionsServer) SendAndClose(m *RecordSummary) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gitLabExporterRecordSectionsServer) Recv() (*typespb.Section, error) {
	m := new(typespb.Section)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _GitLabExporter_RecordBridges_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GitLabExporterServer).RecordBridges(&gitLabExporterRecordBridgesServer{stream})
}

type GitLabExporter_RecordBridgesServer interface {
	SendAndClose(*RecordSummary) error
	Recv() (*typespb.Bridge, error)
	grpc.ServerStream
}

type gitLabExporterRecordBridgesServer struct {
	grpc.ServerStream
}

func (x *gitLabExporterRecordBridgesServer) SendAndClose(m *RecordSummary) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gitLabExporterRecordBridgesServer) Recv() (*typespb.Bridge, error) {
	m := new(typespb.Bridge)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _GitLabExporter_RecordTestReports_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GitLabExporterServer).RecordTestReports(&gitLabExporterRecordTestReportsServer{stream})
}

type GitLabExporter_RecordTestReportsServer interface {
	SendAndClose(*RecordSummary) error
	Recv() (*typespb.TestReport, error)
	grpc.ServerStream
}

type gitLabExporterRecordTestReportsServer struct {
	grpc.ServerStream
}

func (x *gitLabExporterRecordTestReportsServer) SendAndClose(m *RecordSummary) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gitLabExporterRecordTestReportsServer) Recv() (*typespb.TestReport, error) {
	m := new(typespb.TestReport)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _GitLabExporter_RecordTestSuites_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GitLabExporterServer).RecordTestSuites(&gitLabExporterRecordTestSuitesServer{stream})
}

type GitLabExporter_RecordTestSuitesServer interface {
	SendAndClose(*RecordSummary) error
	Recv() (*typespb.TestSuite, error)
	grpc.ServerStream
}

type gitLabExporterRecordTestSuitesServer struct {
	grpc.ServerStream
}

func (x *gitLabExporterRecordTestSuitesServer) SendAndClose(m *RecordSummary) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gitLabExporterRecordTestSuitesServer) Recv() (*typespb.TestSuite, error) {
	m := new(typespb.TestSuite)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _GitLabExporter_RecordTestCases_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GitLabExporterServer).RecordTestCases(&gitLabExporterRecordTestCasesServer{stream})
}

type GitLabExporter_RecordTestCasesServer interface {
	SendAndClose(*RecordSummary) error
	Recv() (*typespb.TestCase, error)
	grpc.ServerStream
}

type gitLabExporterRecordTestCasesServer struct {
	grpc.ServerStream
}

func (x *gitLabExporterRecordTestCasesServer) SendAndClose(m *RecordSummary) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gitLabExporterRecordTestCasesServer) Recv() (*typespb.TestCase, error) {
	m := new(typespb.TestCase)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _GitLabExporter_RecordMetrics_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GitLabExporterServer).RecordMetrics(&gitLabExporterRecordMetricsServer{stream})
}

type GitLabExporter_RecordMetricsServer interface {
	SendAndClose(*RecordSummary) error
	Recv() (*typespb.Metric, error)
	grpc.ServerStream
}

type gitLabExporterRecordMetricsServer struct {
	grpc.ServerStream
}

func (x *gitLabExporterRecordMetricsServer) SendAndClose(m *RecordSummary) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gitLabExporterRecordMetricsServer) Recv() (*typespb.Metric, error) {
	m := new(typespb.Metric)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _GitLabExporter_RecordTraces_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GitLabExporterServer).RecordTraces(&gitLabExporterRecordTracesServer{stream})
}

type GitLabExporter_RecordTracesServer interface {
	SendAndClose(*RecordSummary) error
	Recv() (*typespb.Trace, error)
	grpc.ServerStream
}

type gitLabExporterRecordTracesServer struct {
	grpc.ServerStream
}

func (x *gitLabExporterRecordTracesServer) SendAndClose(m *RecordSummary) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gitLabExporterRecordTracesServer) Recv() (*typespb.Trace, error) {
	m := new(typespb.Trace)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GitLabExporter_ServiceDesc is the grpc.ServiceDesc for GitLabExporter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GitLabExporter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gitlabexporter.protobuf.service.GitLabExporter",
	HandlerType: (*GitLabExporterServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RecordPipelines",
			Handler:       _GitLabExporter_RecordPipelines_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "RecordJobs",
			Handler:       _GitLabExporter_RecordJobs_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "RecordSections",
			Handler:       _GitLabExporter_RecordSections_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "RecordBridges",
			Handler:       _GitLabExporter_RecordBridges_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "RecordTestReports",
			Handler:       _GitLabExporter_RecordTestReports_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "RecordTestSuites",
			Handler:       _GitLabExporter_RecordTestSuites_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "RecordTestCases",
			Handler:       _GitLabExporter_RecordTestCases_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "RecordMetrics",
			Handler:       _GitLabExporter_RecordMetrics_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "RecordTraces",
			Handler:       _GitLabExporter_RecordTraces_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "gitlabexporter/protobuf/service/service.proto",
}