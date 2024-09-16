// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.22.2
// source: gitlabexporter/protobuf/service/service.proto

package servicepb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	GitLabExporter_RecordBridges_FullMethodName                = "/gitlabexporter.protobuf.service.GitLabExporter/RecordBridges"
	GitLabExporter_RecordCommits_FullMethodName                = "/gitlabexporter.protobuf.service.GitLabExporter/RecordCommits"
	GitLabExporter_RecordJobs_FullMethodName                   = "/gitlabexporter.protobuf.service.GitLabExporter/RecordJobs"
	GitLabExporter_RecordMergeRequests_FullMethodName          = "/gitlabexporter.protobuf.service.GitLabExporter/RecordMergeRequests"
	GitLabExporter_RecordMergeRequestNoteEvents_FullMethodName = "/gitlabexporter.protobuf.service.GitLabExporter/RecordMergeRequestNoteEvents"
	GitLabExporter_RecordMetrics_FullMethodName                = "/gitlabexporter.protobuf.service.GitLabExporter/RecordMetrics"
	GitLabExporter_RecordPipelines_FullMethodName              = "/gitlabexporter.protobuf.service.GitLabExporter/RecordPipelines"
	GitLabExporter_RecordProjects_FullMethodName               = "/gitlabexporter.protobuf.service.GitLabExporter/RecordProjects"
	GitLabExporter_RecordSections_FullMethodName               = "/gitlabexporter.protobuf.service.GitLabExporter/RecordSections"
	GitLabExporter_RecordTestCases_FullMethodName              = "/gitlabexporter.protobuf.service.GitLabExporter/RecordTestCases"
	GitLabExporter_RecordTestReports_FullMethodName            = "/gitlabexporter.protobuf.service.GitLabExporter/RecordTestReports"
	GitLabExporter_RecordTestSuites_FullMethodName             = "/gitlabexporter.protobuf.service.GitLabExporter/RecordTestSuites"
	GitLabExporter_RecordTraces_FullMethodName                 = "/gitlabexporter.protobuf.service.GitLabExporter/RecordTraces"
	GitLabExporter_RecordUsers_FullMethodName                  = "/gitlabexporter.protobuf.service.GitLabExporter/RecordUsers"
)

// GitLabExporterClient is the client API for GitLabExporter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GitLabExporterClient interface {
	RecordBridges(ctx context.Context, in *RecordBridgesRequest, opts ...grpc.CallOption) (*RecordSummary, error)
	RecordCommits(ctx context.Context, in *RecordCommitsRequest, opts ...grpc.CallOption) (*RecordSummary, error)
	RecordJobs(ctx context.Context, in *RecordJobsRequest, opts ...grpc.CallOption) (*RecordSummary, error)
	RecordMergeRequests(ctx context.Context, in *RecordMergeRequestsRequest, opts ...grpc.CallOption) (*RecordSummary, error)
	RecordMergeRequestNoteEvents(ctx context.Context, in *RecordMergeRequestNoteEventsRequest, opts ...grpc.CallOption) (*RecordSummary, error)
	RecordMetrics(ctx context.Context, in *RecordMetricsRequest, opts ...grpc.CallOption) (*RecordSummary, error)
	RecordPipelines(ctx context.Context, in *RecordPipelinesRequest, opts ...grpc.CallOption) (*RecordSummary, error)
	RecordProjects(ctx context.Context, in *RecordProjectsRequest, opts ...grpc.CallOption) (*RecordSummary, error)
	RecordSections(ctx context.Context, in *RecordSectionsRequest, opts ...grpc.CallOption) (*RecordSummary, error)
	RecordTestCases(ctx context.Context, in *RecordTestCasesRequest, opts ...grpc.CallOption) (*RecordSummary, error)
	RecordTestReports(ctx context.Context, in *RecordTestReportsRequest, opts ...grpc.CallOption) (*RecordSummary, error)
	RecordTestSuites(ctx context.Context, in *RecordTestSuitesRequest, opts ...grpc.CallOption) (*RecordSummary, error)
	RecordTraces(ctx context.Context, in *RecordTracesRequest, opts ...grpc.CallOption) (*RecordSummary, error)
	RecordUsers(ctx context.Context, in *RecordUsersRequest, opts ...grpc.CallOption) (*RecordSummary, error)
}

type gitLabExporterClient struct {
	cc grpc.ClientConnInterface
}

func NewGitLabExporterClient(cc grpc.ClientConnInterface) GitLabExporterClient {
	return &gitLabExporterClient{cc}
}

func (c *gitLabExporterClient) RecordBridges(ctx context.Context, in *RecordBridgesRequest, opts ...grpc.CallOption) (*RecordSummary, error) {
	out := new(RecordSummary)
	err := c.cc.Invoke(ctx, GitLabExporter_RecordBridges_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitLabExporterClient) RecordCommits(ctx context.Context, in *RecordCommitsRequest, opts ...grpc.CallOption) (*RecordSummary, error) {
	out := new(RecordSummary)
	err := c.cc.Invoke(ctx, GitLabExporter_RecordCommits_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitLabExporterClient) RecordJobs(ctx context.Context, in *RecordJobsRequest, opts ...grpc.CallOption) (*RecordSummary, error) {
	out := new(RecordSummary)
	err := c.cc.Invoke(ctx, GitLabExporter_RecordJobs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitLabExporterClient) RecordMergeRequests(ctx context.Context, in *RecordMergeRequestsRequest, opts ...grpc.CallOption) (*RecordSummary, error) {
	out := new(RecordSummary)
	err := c.cc.Invoke(ctx, GitLabExporter_RecordMergeRequests_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitLabExporterClient) RecordMergeRequestNoteEvents(ctx context.Context, in *RecordMergeRequestNoteEventsRequest, opts ...grpc.CallOption) (*RecordSummary, error) {
	out := new(RecordSummary)
	err := c.cc.Invoke(ctx, GitLabExporter_RecordMergeRequestNoteEvents_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitLabExporterClient) RecordMetrics(ctx context.Context, in *RecordMetricsRequest, opts ...grpc.CallOption) (*RecordSummary, error) {
	out := new(RecordSummary)
	err := c.cc.Invoke(ctx, GitLabExporter_RecordMetrics_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitLabExporterClient) RecordPipelines(ctx context.Context, in *RecordPipelinesRequest, opts ...grpc.CallOption) (*RecordSummary, error) {
	out := new(RecordSummary)
	err := c.cc.Invoke(ctx, GitLabExporter_RecordPipelines_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitLabExporterClient) RecordProjects(ctx context.Context, in *RecordProjectsRequest, opts ...grpc.CallOption) (*RecordSummary, error) {
	out := new(RecordSummary)
	err := c.cc.Invoke(ctx, GitLabExporter_RecordProjects_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitLabExporterClient) RecordSections(ctx context.Context, in *RecordSectionsRequest, opts ...grpc.CallOption) (*RecordSummary, error) {
	out := new(RecordSummary)
	err := c.cc.Invoke(ctx, GitLabExporter_RecordSections_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitLabExporterClient) RecordTestCases(ctx context.Context, in *RecordTestCasesRequest, opts ...grpc.CallOption) (*RecordSummary, error) {
	out := new(RecordSummary)
	err := c.cc.Invoke(ctx, GitLabExporter_RecordTestCases_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitLabExporterClient) RecordTestReports(ctx context.Context, in *RecordTestReportsRequest, opts ...grpc.CallOption) (*RecordSummary, error) {
	out := new(RecordSummary)
	err := c.cc.Invoke(ctx, GitLabExporter_RecordTestReports_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitLabExporterClient) RecordTestSuites(ctx context.Context, in *RecordTestSuitesRequest, opts ...grpc.CallOption) (*RecordSummary, error) {
	out := new(RecordSummary)
	err := c.cc.Invoke(ctx, GitLabExporter_RecordTestSuites_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitLabExporterClient) RecordTraces(ctx context.Context, in *RecordTracesRequest, opts ...grpc.CallOption) (*RecordSummary, error) {
	out := new(RecordSummary)
	err := c.cc.Invoke(ctx, GitLabExporter_RecordTraces_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitLabExporterClient) RecordUsers(ctx context.Context, in *RecordUsersRequest, opts ...grpc.CallOption) (*RecordSummary, error) {
	out := new(RecordSummary)
	err := c.cc.Invoke(ctx, GitLabExporter_RecordUsers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GitLabExporterServer is the server API for GitLabExporter service.
// All implementations must embed UnimplementedGitLabExporterServer
// for forward compatibility
type GitLabExporterServer interface {
	RecordBridges(context.Context, *RecordBridgesRequest) (*RecordSummary, error)
	RecordCommits(context.Context, *RecordCommitsRequest) (*RecordSummary, error)
	RecordJobs(context.Context, *RecordJobsRequest) (*RecordSummary, error)
	RecordMergeRequests(context.Context, *RecordMergeRequestsRequest) (*RecordSummary, error)
	RecordMergeRequestNoteEvents(context.Context, *RecordMergeRequestNoteEventsRequest) (*RecordSummary, error)
	RecordMetrics(context.Context, *RecordMetricsRequest) (*RecordSummary, error)
	RecordPipelines(context.Context, *RecordPipelinesRequest) (*RecordSummary, error)
	RecordProjects(context.Context, *RecordProjectsRequest) (*RecordSummary, error)
	RecordSections(context.Context, *RecordSectionsRequest) (*RecordSummary, error)
	RecordTestCases(context.Context, *RecordTestCasesRequest) (*RecordSummary, error)
	RecordTestReports(context.Context, *RecordTestReportsRequest) (*RecordSummary, error)
	RecordTestSuites(context.Context, *RecordTestSuitesRequest) (*RecordSummary, error)
	RecordTraces(context.Context, *RecordTracesRequest) (*RecordSummary, error)
	RecordUsers(context.Context, *RecordUsersRequest) (*RecordSummary, error)
	mustEmbedUnimplementedGitLabExporterServer()
}

// UnimplementedGitLabExporterServer must be embedded to have forward compatible implementations.
type UnimplementedGitLabExporterServer struct {
}

func (UnimplementedGitLabExporterServer) RecordBridges(context.Context, *RecordBridgesRequest) (*RecordSummary, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordBridges not implemented")
}
func (UnimplementedGitLabExporterServer) RecordCommits(context.Context, *RecordCommitsRequest) (*RecordSummary, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordCommits not implemented")
}
func (UnimplementedGitLabExporterServer) RecordJobs(context.Context, *RecordJobsRequest) (*RecordSummary, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordJobs not implemented")
}
func (UnimplementedGitLabExporterServer) RecordMergeRequests(context.Context, *RecordMergeRequestsRequest) (*RecordSummary, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordMergeRequests not implemented")
}
func (UnimplementedGitLabExporterServer) RecordMergeRequestNoteEvents(context.Context, *RecordMergeRequestNoteEventsRequest) (*RecordSummary, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordMergeRequestNoteEvents not implemented")
}
func (UnimplementedGitLabExporterServer) RecordMetrics(context.Context, *RecordMetricsRequest) (*RecordSummary, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordMetrics not implemented")
}
func (UnimplementedGitLabExporterServer) RecordPipelines(context.Context, *RecordPipelinesRequest) (*RecordSummary, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordPipelines not implemented")
}
func (UnimplementedGitLabExporterServer) RecordProjects(context.Context, *RecordProjectsRequest) (*RecordSummary, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordProjects not implemented")
}
func (UnimplementedGitLabExporterServer) RecordSections(context.Context, *RecordSectionsRequest) (*RecordSummary, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordSections not implemented")
}
func (UnimplementedGitLabExporterServer) RecordTestCases(context.Context, *RecordTestCasesRequest) (*RecordSummary, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordTestCases not implemented")
}
func (UnimplementedGitLabExporterServer) RecordTestReports(context.Context, *RecordTestReportsRequest) (*RecordSummary, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordTestReports not implemented")
}
func (UnimplementedGitLabExporterServer) RecordTestSuites(context.Context, *RecordTestSuitesRequest) (*RecordSummary, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordTestSuites not implemented")
}
func (UnimplementedGitLabExporterServer) RecordTraces(context.Context, *RecordTracesRequest) (*RecordSummary, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordTraces not implemented")
}
func (UnimplementedGitLabExporterServer) RecordUsers(context.Context, *RecordUsersRequest) (*RecordSummary, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordUsers not implemented")
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

func _GitLabExporter_RecordBridges_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordBridgesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitLabExporterServer).RecordBridges(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitLabExporter_RecordBridges_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitLabExporterServer).RecordBridges(ctx, req.(*RecordBridgesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitLabExporter_RecordCommits_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordCommitsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitLabExporterServer).RecordCommits(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitLabExporter_RecordCommits_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitLabExporterServer).RecordCommits(ctx, req.(*RecordCommitsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitLabExporter_RecordJobs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordJobsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitLabExporterServer).RecordJobs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitLabExporter_RecordJobs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitLabExporterServer).RecordJobs(ctx, req.(*RecordJobsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitLabExporter_RecordMergeRequests_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordMergeRequestsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitLabExporterServer).RecordMergeRequests(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitLabExporter_RecordMergeRequests_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitLabExporterServer).RecordMergeRequests(ctx, req.(*RecordMergeRequestsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitLabExporter_RecordMergeRequestNoteEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordMergeRequestNoteEventsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitLabExporterServer).RecordMergeRequestNoteEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitLabExporter_RecordMergeRequestNoteEvents_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitLabExporterServer).RecordMergeRequestNoteEvents(ctx, req.(*RecordMergeRequestNoteEventsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitLabExporter_RecordMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordMetricsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitLabExporterServer).RecordMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitLabExporter_RecordMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitLabExporterServer).RecordMetrics(ctx, req.(*RecordMetricsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitLabExporter_RecordPipelines_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordPipelinesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitLabExporterServer).RecordPipelines(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitLabExporter_RecordPipelines_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitLabExporterServer).RecordPipelines(ctx, req.(*RecordPipelinesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitLabExporter_RecordProjects_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordProjectsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitLabExporterServer).RecordProjects(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitLabExporter_RecordProjects_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitLabExporterServer).RecordProjects(ctx, req.(*RecordProjectsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitLabExporter_RecordSections_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordSectionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitLabExporterServer).RecordSections(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitLabExporter_RecordSections_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitLabExporterServer).RecordSections(ctx, req.(*RecordSectionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitLabExporter_RecordTestCases_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordTestCasesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitLabExporterServer).RecordTestCases(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitLabExporter_RecordTestCases_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitLabExporterServer).RecordTestCases(ctx, req.(*RecordTestCasesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitLabExporter_RecordTestReports_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordTestReportsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitLabExporterServer).RecordTestReports(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitLabExporter_RecordTestReports_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitLabExporterServer).RecordTestReports(ctx, req.(*RecordTestReportsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitLabExporter_RecordTestSuites_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordTestSuitesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitLabExporterServer).RecordTestSuites(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitLabExporter_RecordTestSuites_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitLabExporterServer).RecordTestSuites(ctx, req.(*RecordTestSuitesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitLabExporter_RecordTraces_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordTracesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitLabExporterServer).RecordTraces(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitLabExporter_RecordTraces_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitLabExporterServer).RecordTraces(ctx, req.(*RecordTracesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitLabExporter_RecordUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitLabExporterServer).RecordUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitLabExporter_RecordUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitLabExporterServer).RecordUsers(ctx, req.(*RecordUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GitLabExporter_ServiceDesc is the grpc.ServiceDesc for GitLabExporter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GitLabExporter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gitlabexporter.protobuf.service.GitLabExporter",
	HandlerType: (*GitLabExporterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RecordBridges",
			Handler:    _GitLabExporter_RecordBridges_Handler,
		},
		{
			MethodName: "RecordCommits",
			Handler:    _GitLabExporter_RecordCommits_Handler,
		},
		{
			MethodName: "RecordJobs",
			Handler:    _GitLabExporter_RecordJobs_Handler,
		},
		{
			MethodName: "RecordMergeRequests",
			Handler:    _GitLabExporter_RecordMergeRequests_Handler,
		},
		{
			MethodName: "RecordMergeRequestNoteEvents",
			Handler:    _GitLabExporter_RecordMergeRequestNoteEvents_Handler,
		},
		{
			MethodName: "RecordMetrics",
			Handler:    _GitLabExporter_RecordMetrics_Handler,
		},
		{
			MethodName: "RecordPipelines",
			Handler:    _GitLabExporter_RecordPipelines_Handler,
		},
		{
			MethodName: "RecordProjects",
			Handler:    _GitLabExporter_RecordProjects_Handler,
		},
		{
			MethodName: "RecordSections",
			Handler:    _GitLabExporter_RecordSections_Handler,
		},
		{
			MethodName: "RecordTestCases",
			Handler:    _GitLabExporter_RecordTestCases_Handler,
		},
		{
			MethodName: "RecordTestReports",
			Handler:    _GitLabExporter_RecordTestReports_Handler,
		},
		{
			MethodName: "RecordTestSuites",
			Handler:    _GitLabExporter_RecordTestSuites_Handler,
		},
		{
			MethodName: "RecordTraces",
			Handler:    _GitLabExporter_RecordTraces_Handler,
		},
		{
			MethodName: "RecordUsers",
			Handler:    _GitLabExporter_RecordUsers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gitlabexporter/protobuf/service/service.proto",
}
