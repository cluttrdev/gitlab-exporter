// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.22.2
// source: gitlabexporter/protobuf/mergerequest.proto

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

type MergeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ID of the merge request.
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// Internal ID of the merge request.
	Iid int64 `protobuf:"varint,2,opt,name=iid,proto3" json:"iid,omitempty"`
	// ID of the merge request project.
	ProjectId int64 `protobuf:"varint,3,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	// Timestamp of when the merge request was created.
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// Timestamp of when the merge request was updated.
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	// Timestamp of when the merge request merged.
	MergedAt *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=merged_at,json=mergedAt,proto3" json:"merged_at,omitempty"`
	// Timestamp of when the merge request was closed.
	ClosedAt *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=closed_at,json=closedAt,proto3" json:"closed_at,omitempty"`
	// ID of the merge request source project.
	SourceProjectId int64 `protobuf:"varint,8,opt,name=source_project_id,json=sourceProjectId,proto3" json:"source_project_id,omitempty"`
	// ID of the merge request target project.
	TargetProjectId int64 `protobuf:"varint,9,opt,name=target_project_id,json=targetProjectId,proto3" json:"target_project_id,omitempty"`
	// Source branch of the merge request.
	SourceBranch string `protobuf:"bytes,10,opt,name=source_branch,json=sourceBranch,proto3" json:"source_branch,omitempty"`
	// Target branch of the merge request.
	TargetBranch string `protobuf:"bytes,11,opt,name=target_branch,json=targetBranch,proto3" json:"target_branch,omitempty"`
	// Title of the merge request.
	Title string `protobuf:"bytes,12,opt,name=title,proto3" json:"title,omitempty"`
	// State of the merge request. Can be opened, closed, merged or locked.
	State string `protobuf:"bytes,13,opt,name=state,proto3" json:"state,omitempty"`
	// Detailed merge status of the merge request.
	DetailedMergeStatus string `protobuf:"bytes,14,opt,name=detailed_merge_status,json=detailedMergeStatus,proto3" json:"detailed_merge_status,omitempty"`
	// Indicates if the merge request is a draft.
	Draft bool `protobuf:"varint,15,opt,name=draft,proto3" json:"draft,omitempty"`
	// Indicates if merge request has conflicts and cannot merge.
	HasConflicts bool `protobuf:"varint,16,opt,name=has_conflicts,json=hasConflicts,proto3" json:"has_conflicts,omitempty"`
	// Error message shown when a merge has failed.
	MergeError string `protobuf:"bytes,17,opt,name=merge_error,json=mergeError,proto3" json:"merge_error,omitempty"`
	// References of the base SHA, the head SHA, and the start SHA for this merge request.
	DiffRefs *MergeRequestDiffRefs `protobuf:"bytes,18,opt,name=diff_refs,json=diffRefs,proto3" json:"diff_refs,omitempty"`
	// User who created this merge request.
	Author *User `protobuf:"bytes,19,opt,name=author,proto3" json:"author,omitempty"`
	// First assignee of the merge request.
	Assignee *User `protobuf:"bytes,20,opt,name=assignee,proto3" json:"assignee,omitempty"`
	// Assignees of the merge request.
	Assignees []*User `protobuf:"bytes,21,rep,name=assignees,proto3" json:"assignees,omitempty"`
	// Reviewers of the merge request.
	Reviewers []*User `protobuf:"bytes,22,rep,name=reviewers,proto3" json:"reviewers,omitempty"`
	// The user who merged this merge request, the user who set it to auto-merge, or null.
	MergeUser *User `protobuf:"bytes,23,opt,name=merge_user,json=mergeUser,proto3" json:"merge_user,omitempty"`
	// User who closed this merge request.
	CloseUser *User `protobuf:"bytes,24,opt,name=close_user,json=closeUser,proto3" json:"close_user,omitempty"`
	// Labels of the merge request.
	Labels []string `protobuf:"bytes,25,rep,name=labels,proto3" json:"labels,omitempty"`
	// Diff head SHA of the merge request.
	Sha string `protobuf:"bytes,26,opt,name=sha,proto3" json:"sha,omitempty"`
	// SHA of the merge request commit. Empty until merged.
	MergeCommitSha string `protobuf:"bytes,27,opt,name=merge_commit_sha,json=mergeCommitSha,proto3" json:"merge_commit_sha,omitempty"`
	// SHA of the squash commit. Empty until merged.
	SquashCommitSha string `protobuf:"bytes,28,opt,name=squash_commit_sha,json=squashCommitSha,proto3" json:"squash_commit_sha,omitempty"`
	// Number of changes made on the merge request.
	ChangesCount string `protobuf:"bytes,29,opt,name=changes_count,json=changesCount,proto3" json:"changes_count,omitempty"`
	// User notes count of the merge request.
	UserNotesCount int64 `protobuf:"varint,30,opt,name=user_notes_count,json=userNotesCount,proto3" json:"user_notes_count,omitempty"`
	// Number of upvotes for the merge request.
	Upvotes int64 `protobuf:"varint,31,opt,name=upvotes,proto3" json:"upvotes,omitempty"`
	// Number of downvotes for the merge request.
	Downvotes int64 `protobuf:"varint,32,opt,name=downvotes,proto3" json:"downvotes,omitempty"`
	// Pipeline running on the branch HEAD of the merge request.
	Pipeline *PipelineInfo `protobuf:"bytes,33,opt,name=pipeline,proto3" json:"pipeline,omitempty"`
	// Milestone of the merge request.
	Milestone *Milestone `protobuf:"bytes,34,opt,name=milestone,proto3" json:"milestone,omitempty"`
	// Web URL of the merge request.
	WebUrl string `protobuf:"bytes,35,opt,name=web_url,json=webUrl,proto3" json:"web_url,omitempty"`
}

func (x *MergeRequest) Reset() {
	*x = MergeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gitlabexporter_protobuf_mergerequest_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MergeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MergeRequest) ProtoMessage() {}

func (x *MergeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gitlabexporter_protobuf_mergerequest_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MergeRequest.ProtoReflect.Descriptor instead.
func (*MergeRequest) Descriptor() ([]byte, []int) {
	return file_gitlabexporter_protobuf_mergerequest_proto_rawDescGZIP(), []int{0}
}

func (x *MergeRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *MergeRequest) GetIid() int64 {
	if x != nil {
		return x.Iid
	}
	return 0
}

func (x *MergeRequest) GetProjectId() int64 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *MergeRequest) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *MergeRequest) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *MergeRequest) GetMergedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.MergedAt
	}
	return nil
}

func (x *MergeRequest) GetClosedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ClosedAt
	}
	return nil
}

func (x *MergeRequest) GetSourceProjectId() int64 {
	if x != nil {
		return x.SourceProjectId
	}
	return 0
}

func (x *MergeRequest) GetTargetProjectId() int64 {
	if x != nil {
		return x.TargetProjectId
	}
	return 0
}

func (x *MergeRequest) GetSourceBranch() string {
	if x != nil {
		return x.SourceBranch
	}
	return ""
}

func (x *MergeRequest) GetTargetBranch() string {
	if x != nil {
		return x.TargetBranch
	}
	return ""
}

func (x *MergeRequest) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *MergeRequest) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *MergeRequest) GetDetailedMergeStatus() string {
	if x != nil {
		return x.DetailedMergeStatus
	}
	return ""
}

func (x *MergeRequest) GetDraft() bool {
	if x != nil {
		return x.Draft
	}
	return false
}

func (x *MergeRequest) GetHasConflicts() bool {
	if x != nil {
		return x.HasConflicts
	}
	return false
}

func (x *MergeRequest) GetMergeError() string {
	if x != nil {
		return x.MergeError
	}
	return ""
}

func (x *MergeRequest) GetDiffRefs() *MergeRequestDiffRefs {
	if x != nil {
		return x.DiffRefs
	}
	return nil
}

func (x *MergeRequest) GetAuthor() *User {
	if x != nil {
		return x.Author
	}
	return nil
}

func (x *MergeRequest) GetAssignee() *User {
	if x != nil {
		return x.Assignee
	}
	return nil
}

func (x *MergeRequest) GetAssignees() []*User {
	if x != nil {
		return x.Assignees
	}
	return nil
}

func (x *MergeRequest) GetReviewers() []*User {
	if x != nil {
		return x.Reviewers
	}
	return nil
}

func (x *MergeRequest) GetMergeUser() *User {
	if x != nil {
		return x.MergeUser
	}
	return nil
}

func (x *MergeRequest) GetCloseUser() *User {
	if x != nil {
		return x.CloseUser
	}
	return nil
}

func (x *MergeRequest) GetLabels() []string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *MergeRequest) GetSha() string {
	if x != nil {
		return x.Sha
	}
	return ""
}

func (x *MergeRequest) GetMergeCommitSha() string {
	if x != nil {
		return x.MergeCommitSha
	}
	return ""
}

func (x *MergeRequest) GetSquashCommitSha() string {
	if x != nil {
		return x.SquashCommitSha
	}
	return ""
}

func (x *MergeRequest) GetChangesCount() string {
	if x != nil {
		return x.ChangesCount
	}
	return ""
}

func (x *MergeRequest) GetUserNotesCount() int64 {
	if x != nil {
		return x.UserNotesCount
	}
	return 0
}

func (x *MergeRequest) GetUpvotes() int64 {
	if x != nil {
		return x.Upvotes
	}
	return 0
}

func (x *MergeRequest) GetDownvotes() int64 {
	if x != nil {
		return x.Downvotes
	}
	return 0
}

func (x *MergeRequest) GetPipeline() *PipelineInfo {
	if x != nil {
		return x.Pipeline
	}
	return nil
}

func (x *MergeRequest) GetMilestone() *Milestone {
	if x != nil {
		return x.Milestone
	}
	return nil
}

func (x *MergeRequest) GetWebUrl() string {
	if x != nil {
		return x.WebUrl
	}
	return ""
}

type MergeRequestDiffRefs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BaseSha  string `protobuf:"bytes,1,opt,name=base_sha,json=baseSha,proto3" json:"base_sha,omitempty"`
	HeadSha  string `protobuf:"bytes,2,opt,name=head_sha,json=headSha,proto3" json:"head_sha,omitempty"`
	StartSha string `protobuf:"bytes,3,opt,name=start_sha,json=startSha,proto3" json:"start_sha,omitempty"`
}

func (x *MergeRequestDiffRefs) Reset() {
	*x = MergeRequestDiffRefs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gitlabexporter_protobuf_mergerequest_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MergeRequestDiffRefs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MergeRequestDiffRefs) ProtoMessage() {}

func (x *MergeRequestDiffRefs) ProtoReflect() protoreflect.Message {
	mi := &file_gitlabexporter_protobuf_mergerequest_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MergeRequestDiffRefs.ProtoReflect.Descriptor instead.
func (*MergeRequestDiffRefs) Descriptor() ([]byte, []int) {
	return file_gitlabexporter_protobuf_mergerequest_proto_rawDescGZIP(), []int{1}
}

func (x *MergeRequestDiffRefs) GetBaseSha() string {
	if x != nil {
		return x.BaseSha
	}
	return ""
}

func (x *MergeRequestDiffRefs) GetHeadSha() string {
	if x != nil {
		return x.HeadSha
	}
	return ""
}

func (x *MergeRequestDiffRefs) GetStartSha() string {
	if x != nil {
		return x.StartSha
	}
	return ""
}

type Milestone struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Iid       int64                  `protobuf:"varint,2,opt,name=iid,proto3" json:"iid,omitempty"`
	ProjectId int64                  `protobuf:"varint,3,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	GroupId   int64                  `protobuf:"varint,4,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	StartDate *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=start_date,json=startDate,proto3" json:"start_date,omitempty"`
	DueDate   *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=due_date,json=dueDate,proto3" json:"due_date,omitempty"`
	Title     string                 `protobuf:"bytes,9,opt,name=title,proto3" json:"title,omitempty"`
	State     string                 `protobuf:"bytes,10,opt,name=state,proto3" json:"state,omitempty"`
	Expired   bool                   `protobuf:"varint,11,opt,name=expired,proto3" json:"expired,omitempty"`
	WebUrl    string                 `protobuf:"bytes,12,opt,name=web_url,json=webUrl,proto3" json:"web_url,omitempty"`
}

func (x *Milestone) Reset() {
	*x = Milestone{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gitlabexporter_protobuf_mergerequest_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Milestone) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Milestone) ProtoMessage() {}

func (x *Milestone) ProtoReflect() protoreflect.Message {
	mi := &file_gitlabexporter_protobuf_mergerequest_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Milestone.ProtoReflect.Descriptor instead.
func (*Milestone) Descriptor() ([]byte, []int) {
	return file_gitlabexporter_protobuf_mergerequest_proto_rawDescGZIP(), []int{2}
}

func (x *Milestone) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Milestone) GetIid() int64 {
	if x != nil {
		return x.Iid
	}
	return 0
}

func (x *Milestone) GetProjectId() int64 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *Milestone) GetGroupId() int64 {
	if x != nil {
		return x.GroupId
	}
	return 0
}

func (x *Milestone) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Milestone) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *Milestone) GetStartDate() *timestamppb.Timestamp {
	if x != nil {
		return x.StartDate
	}
	return nil
}

func (x *Milestone) GetDueDate() *timestamppb.Timestamp {
	if x != nil {
		return x.DueDate
	}
	return nil
}

func (x *Milestone) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Milestone) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *Milestone) GetExpired() bool {
	if x != nil {
		return x.Expired
	}
	return false
}

func (x *Milestone) GetWebUrl() string {
	if x != nil {
		return x.WebUrl
	}
	return ""
}

var File_gitlabexporter_protobuf_mergerequest_proto protoreflect.FileDescriptor

var file_gitlabexporter_protobuf_mergerequest_proto_rawDesc = []byte{
	0x0a, 0x2a, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x6d, 0x65, 0x72, 0x67, 0x65, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x67, 0x69,
	0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x26, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78,
	0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x22,
	0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xee, 0x0b, 0x0a, 0x0c, 0x4d, 0x65, 0x72, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x03, 0x69, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65,
	0x63, 0x74, 0x49, 0x64, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12,
	0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x37, 0x0a, 0x09, 0x6d, 0x65,
	0x72, 0x67, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x08, 0x6d, 0x65, 0x72, 0x67, 0x65,
	0x64, 0x41, 0x74, 0x12, 0x37, 0x0a, 0x09, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x08, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x64, 0x41, 0x74, 0x12, 0x2a, 0x0a, 0x11,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x2a, 0x0a, 0x11, 0x74, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x5f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0f, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x6a, 0x65,
	0x63, 0x74, 0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x62,
	0x72, 0x61, 0x6e, 0x63, 0x68, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x42, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x12, 0x23, 0x0a, 0x0d, 0x74, 0x61, 0x72,
	0x67, 0x65, 0x74, 0x5f, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0c, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x42, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x12, 0x14,
	0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74,
	0x69, 0x74, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x0d, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x32, 0x0a, 0x15, 0x64, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x5f, 0x6d, 0x65, 0x72, 0x67, 0x65, 0x5f, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x13, 0x64, 0x65, 0x74, 0x61, 0x69,
	0x6c, 0x65, 0x64, 0x4d, 0x65, 0x72, 0x67, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x14,
	0x0a, 0x05, 0x64, 0x72, 0x61, 0x66, 0x74, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x64,
	0x72, 0x61, 0x66, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x68, 0x61, 0x73, 0x5f, 0x63, 0x6f, 0x6e, 0x66,
	0x6c, 0x69, 0x63, 0x74, 0x73, 0x18, 0x10, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x68, 0x61, 0x73,
	0x43, 0x6f, 0x6e, 0x66, 0x6c, 0x69, 0x63, 0x74, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x6d, 0x65, 0x72,
	0x67, 0x65, 0x5f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x11, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x6d, 0x65, 0x72, 0x67, 0x65, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x4a, 0x0a, 0x09, 0x64, 0x69,
	0x66, 0x66, 0x5f, 0x72, 0x65, 0x66, 0x73, 0x18, 0x12, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2d, 0x2e,
	0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x72, 0x67, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x44, 0x69, 0x66, 0x66, 0x52, 0x65, 0x66, 0x73, 0x52, 0x08, 0x64, 0x69,
	0x66, 0x66, 0x52, 0x65, 0x66, 0x73, 0x12, 0x35, 0x0a, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x18, 0x13, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65,
	0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x12, 0x39, 0x0a,
	0x08, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x65, 0x18, 0x14, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1d, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x08,
	0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x65, 0x12, 0x3b, 0x0a, 0x09, 0x61, 0x73, 0x73, 0x69,
	0x67, 0x6e, 0x65, 0x65, 0x73, 0x18, 0x15, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x67, 0x69,
	0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x09, 0x61, 0x73, 0x73, 0x69,
	0x67, 0x6e, 0x65, 0x65, 0x73, 0x12, 0x3b, 0x0a, 0x09, 0x72, 0x65, 0x76, 0x69, 0x65, 0x77, 0x65,
	0x72, 0x73, 0x18, 0x16, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61,
	0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x09, 0x72, 0x65, 0x76, 0x69, 0x65, 0x77, 0x65,
	0x72, 0x73, 0x12, 0x3c, 0x0a, 0x0a, 0x6d, 0x65, 0x72, 0x67, 0x65, 0x5f, 0x75, 0x73, 0x65, 0x72,
	0x18, 0x17, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65,
	0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x09, 0x6d, 0x65, 0x72, 0x67, 0x65, 0x55, 0x73, 0x65, 0x72,
	0x12, 0x3c, 0x0a, 0x0a, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x18, 0x18,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70,
	0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x52, 0x09, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x16,
	0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x19, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06,
	0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x68, 0x61, 0x18, 0x1a, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x73, 0x68, 0x61, 0x12, 0x28, 0x0a, 0x10, 0x6d, 0x65, 0x72, 0x67,
	0x65, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x5f, 0x73, 0x68, 0x61, 0x18, 0x1b, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x6d, 0x65, 0x72, 0x67, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x53,
	0x68, 0x61, 0x12, 0x2a, 0x0a, 0x11, 0x73, 0x71, 0x75, 0x61, 0x73, 0x68, 0x5f, 0x63, 0x6f, 0x6d,
	0x6d, 0x69, 0x74, 0x5f, 0x73, 0x68, 0x61, 0x18, 0x1c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x73,
	0x71, 0x75, 0x61, 0x73, 0x68, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x53, 0x68, 0x61, 0x12, 0x23,
	0x0a, 0x0d, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x1d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x12, 0x28, 0x0a, 0x10, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x6e, 0x6f, 0x74, 0x65,
	0x73, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x1e, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x75,
	0x73, 0x65, 0x72, 0x4e, 0x6f, 0x74, 0x65, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x18, 0x0a,
	0x07, 0x75, 0x70, 0x76, 0x6f, 0x74, 0x65, 0x73, 0x18, 0x1f, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07,
	0x75, 0x70, 0x76, 0x6f, 0x74, 0x65, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x64, 0x6f, 0x77, 0x6e, 0x76,
	0x6f, 0x74, 0x65, 0x73, 0x18, 0x20, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x64, 0x6f, 0x77, 0x6e,
	0x76, 0x6f, 0x74, 0x65, 0x73, 0x12, 0x41, 0x0a, 0x08, 0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e,
	0x65, 0x18, 0x21, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62,
	0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x50, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x08,
	0x70, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x12, 0x40, 0x0a, 0x09, 0x6d, 0x69, 0x6c, 0x65,
	0x73, 0x74, 0x6f, 0x6e, 0x65, 0x18, 0x22, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x67, 0x69,
	0x74, 0x6c, 0x61, 0x62, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x69, 0x6c, 0x65, 0x73, 0x74, 0x6f, 0x6e, 0x65, 0x52,
	0x09, 0x6d, 0x69, 0x6c, 0x65, 0x73, 0x74, 0x6f, 0x6e, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x77, 0x65,
	0x62, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x23, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x77, 0x65, 0x62,
	0x55, 0x72, 0x6c, 0x22, 0x69, 0x0a, 0x14, 0x4d, 0x65, 0x72, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x44, 0x69, 0x66, 0x66, 0x52, 0x65, 0x66, 0x73, 0x12, 0x19, 0x0a, 0x08, 0x62,
	0x61, 0x73, 0x65, 0x5f, 0x73, 0x68, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x62,
	0x61, 0x73, 0x65, 0x53, 0x68, 0x61, 0x12, 0x19, 0x0a, 0x08, 0x68, 0x65, 0x61, 0x64, 0x5f, 0x73,
	0x68, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x68, 0x65, 0x61, 0x64, 0x53, 0x68,
	0x61, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x73, 0x68, 0x61, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x74, 0x61, 0x72, 0x74, 0x53, 0x68, 0x61, 0x22, 0xae,
	0x03, 0x0a, 0x09, 0x4d, 0x69, 0x6c, 0x65, 0x73, 0x74, 0x6f, 0x6e, 0x65, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03,
	0x69, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x69, 0x69, 0x64, 0x12, 0x1d,
	0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x19, 0x0a,
	0x08, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x07, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39,
	0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x44, 0x61, 0x74, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x64, 0x75, 0x65,
	0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x64, 0x75, 0x65, 0x44, 0x61, 0x74, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18,
	0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x64, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x65,
	0x78, 0x70, 0x69, 0x72, 0x65, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x77, 0x65, 0x62, 0x5f, 0x75, 0x72,
	0x6c, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x77, 0x65, 0x62, 0x55, 0x72, 0x6c, 0x42,
	0x37, 0x5a, 0x35, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6c,
	0x75, 0x74, 0x74, 0x72, 0x64, 0x65, 0x76, 0x2f, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2d, 0x65,
	0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gitlabexporter_protobuf_mergerequest_proto_rawDescOnce sync.Once
	file_gitlabexporter_protobuf_mergerequest_proto_rawDescData = file_gitlabexporter_protobuf_mergerequest_proto_rawDesc
)

func file_gitlabexporter_protobuf_mergerequest_proto_rawDescGZIP() []byte {
	file_gitlabexporter_protobuf_mergerequest_proto_rawDescOnce.Do(func() {
		file_gitlabexporter_protobuf_mergerequest_proto_rawDescData = protoimpl.X.CompressGZIP(file_gitlabexporter_protobuf_mergerequest_proto_rawDescData)
	})
	return file_gitlabexporter_protobuf_mergerequest_proto_rawDescData
}

var file_gitlabexporter_protobuf_mergerequest_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_gitlabexporter_protobuf_mergerequest_proto_goTypes = []interface{}{
	(*MergeRequest)(nil),          // 0: gitlabexporter.protobuf.MergeRequest
	(*MergeRequestDiffRefs)(nil),  // 1: gitlabexporter.protobuf.MergeRequestDiffRefs
	(*Milestone)(nil),             // 2: gitlabexporter.protobuf.Milestone
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
	(*User)(nil),                  // 4: gitlabexporter.protobuf.User
	(*PipelineInfo)(nil),          // 5: gitlabexporter.protobuf.PipelineInfo
}
var file_gitlabexporter_protobuf_mergerequest_proto_depIdxs = []int32{
	3,  // 0: gitlabexporter.protobuf.MergeRequest.created_at:type_name -> google.protobuf.Timestamp
	3,  // 1: gitlabexporter.protobuf.MergeRequest.updated_at:type_name -> google.protobuf.Timestamp
	3,  // 2: gitlabexporter.protobuf.MergeRequest.merged_at:type_name -> google.protobuf.Timestamp
	3,  // 3: gitlabexporter.protobuf.MergeRequest.closed_at:type_name -> google.protobuf.Timestamp
	1,  // 4: gitlabexporter.protobuf.MergeRequest.diff_refs:type_name -> gitlabexporter.protobuf.MergeRequestDiffRefs
	4,  // 5: gitlabexporter.protobuf.MergeRequest.author:type_name -> gitlabexporter.protobuf.User
	4,  // 6: gitlabexporter.protobuf.MergeRequest.assignee:type_name -> gitlabexporter.protobuf.User
	4,  // 7: gitlabexporter.protobuf.MergeRequest.assignees:type_name -> gitlabexporter.protobuf.User
	4,  // 8: gitlabexporter.protobuf.MergeRequest.reviewers:type_name -> gitlabexporter.protobuf.User
	4,  // 9: gitlabexporter.protobuf.MergeRequest.merge_user:type_name -> gitlabexporter.protobuf.User
	4,  // 10: gitlabexporter.protobuf.MergeRequest.close_user:type_name -> gitlabexporter.protobuf.User
	5,  // 11: gitlabexporter.protobuf.MergeRequest.pipeline:type_name -> gitlabexporter.protobuf.PipelineInfo
	2,  // 12: gitlabexporter.protobuf.MergeRequest.milestone:type_name -> gitlabexporter.protobuf.Milestone
	3,  // 13: gitlabexporter.protobuf.Milestone.created_at:type_name -> google.protobuf.Timestamp
	3,  // 14: gitlabexporter.protobuf.Milestone.updated_at:type_name -> google.protobuf.Timestamp
	3,  // 15: gitlabexporter.protobuf.Milestone.start_date:type_name -> google.protobuf.Timestamp
	3,  // 16: gitlabexporter.protobuf.Milestone.due_date:type_name -> google.protobuf.Timestamp
	17, // [17:17] is the sub-list for method output_type
	17, // [17:17] is the sub-list for method input_type
	17, // [17:17] is the sub-list for extension type_name
	17, // [17:17] is the sub-list for extension extendee
	0,  // [0:17] is the sub-list for field type_name
}

func init() { file_gitlabexporter_protobuf_mergerequest_proto_init() }
func file_gitlabexporter_protobuf_mergerequest_proto_init() {
	if File_gitlabexporter_protobuf_mergerequest_proto != nil {
		return
	}
	file_gitlabexporter_protobuf_pipeline_proto_init()
	file_gitlabexporter_protobuf_user_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_gitlabexporter_protobuf_mergerequest_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MergeRequest); i {
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
		file_gitlabexporter_protobuf_mergerequest_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MergeRequestDiffRefs); i {
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
		file_gitlabexporter_protobuf_mergerequest_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Milestone); i {
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
			RawDescriptor: file_gitlabexporter_protobuf_mergerequest_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_gitlabexporter_protobuf_mergerequest_proto_goTypes,
		DependencyIndexes: file_gitlabexporter_protobuf_mergerequest_proto_depIdxs,
		MessageInfos:      file_gitlabexporter_protobuf_mergerequest_proto_msgTypes,
	}.Build()
	File_gitlabexporter_protobuf_mergerequest_proto = out.File
	file_gitlabexporter_protobuf_mergerequest_proto_rawDesc = nil
	file_gitlabexporter_protobuf_mergerequest_proto_goTypes = nil
	file_gitlabexporter_protobuf_mergerequest_proto_depIdxs = nil
}