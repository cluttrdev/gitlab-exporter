package clickhouse

import "time"

type Project struct {
	Id          int64 `ch:"id"`
	NamespaceId int64 `ch:"namespace_id"`

	Name     string `ch:"name"`
	FullName string `ch:"full_name"`
	Path     string `ch:"path"`
	FullPath string `ch:"full_path"`

	Description string   `ch:"description"`
	Topics      []string `ch:"topics"`

	CreatedAt      float64 `ch:"created_at"`
	UpdatedAt      float64 `ch:"updated_at"`
	LastActivityAt float64 `ch:"last_activity_at"`

	JobArtifactsSize      int64 `ch:"job_artifacts_size"`
	ContainerRegistrySize int64 `ch:"container_registry_size"`
	LfsObjectsSize        int64 `ch:"lfs_objects_size"`
	PackagesSize          int64 `ch:"packages_size"`
	PipelineArtifactsSize int64 `ch:"pipeline_artifacts_size"`
	RepositorySize        int64 `ch:"repository_size"`
	SnippetsSize          int64 `ch:"snippets_size"`
	StorageSize           int64 `ch:"storage_size"`
	UploadsSize           int64 `ch:"uploads_size"`
	WikiSize              int64 `ch:"wiki_size"`

	ForksCount      int64 `ch:"forks_count"`
	StarsCount      int64 `ch:"stars_count"`
	CommitCount     int64 `ch:"commit_count"`
	OpenIssuesCount int64 `ch:"open_issues_count"`

	Archived   bool   `ch:"archived"`
	Visibility string `ch:"visibility"`

	DefaultBranch string `ch:"default_branch"`
}

type Pipeline struct {
	Id        int64 `ch:"id"`
	Iid       int64 `ch:"iid"`
	ProjectId int64 `ch:"project_id"`

	Name          string `ch:"name"`
	Ref           string `ch:"ref"`
	RefPath       string `ch:"ref_path"`
	Sha           string `ch:"sha"`
	Source        string `ch:"source"`
	Status        string `ch:"status"`
	FailureReason string `ch:"failure_reason"`

	CommittedAt float64 `ch:"committed_at"`
	CreatedAt   float64 `ch:"created_at"`
	UpdatedAt   float64 `ch:"updated_at"`
	StartedAt   float64 `ch:"started_at"`
	FinishedAt  float64 `ch:"finished_at"`

	QueuedDuration float64 `ch:"queued_duration"`
	Duration       float64 `ch:"duration"`

	Coverage float64 `ch:"coverage"`

	Warnings   bool `ch:"warnings"`
	YamlErrors bool `ch:"yaml_errors"`

	Child                     bool                `ch:"child"`
	UpstreamPipelineId        int64               `ch:"upstream_pipeline_id"`
	UpstreamPipelineIid       int64               `ch:"upstream_pipeline_iid"`
	UpstreamPipelineProjectId int64               `ch:"upstream_pipeline_project_id"`
	DownstreamPipelines       []pipelineReference `ch:"downstream_pipelines"`

	MergeRequestId        int64 `ch:"merge_request_id"`
	MergeRequestIid       int64 `ch:"merge_request_iid"`
	MergeRequestProjectId int64 `ch:"merge_request_project_id"`

	UserId int64 `ch:"user_id"`
}

type pipelineReference struct {
	Id        int64 `ch:"id"`
	Iid       int64 `ch:"iid"`
	ProjectId int64 `ch:"project_id"`
}

type Issue struct {
	Id        int64 `ch:"id"`
	Iid       int64 `ch:"iid"`
	ProjectId int64 `ch:"project_id"`

	CreatedAt float64 `ch:"created_at"`
	UpdatedAt float64 `ch:"updated_at"`
	ClosedAt  float64 `ch:"closed_at"`

	Title  string   `ch:"title"`
	Labels []string `ch:"labels"`

	Type     string `ch:"type"`
	Severity string `ch:"severity"`
	State    string `ch:"state"`
}

type Job struct {
	Id         int64 `ch:"id"`
	PipelineId int64 `ch:"pipeline_id"`
	ProjectId  int64 `ch:"project_id"`

	Name          string `ch:"name"`
	Ref           string `ch:"ref"`
	RefPath       string `ch:"ref_path"`
	Status        string `ch:"status"`
	FailureReason string `ch:"failure_reason"`
	ExitCode      int64  `ch:"exit_code"`

	CreatedAt  float64 `ch:"created_at"`
	QueuedAt   float64 `ch:"queued_at"`
	StartedAt  float64 `ch:"started_at"`
	FinishedAt float64 `ch:"finished_at"`
	ErasedAt   float64 `ch:"erased_at"`

	QueuedDuration float64 `ch:"queued_duration"`
	Duration       float64 `ch:"duration"`

	Coverage float64 `ch:"coverage"`

	Stage      string     `ch:"stage"`
	TagList    []string   `ch:"tag_list"`
	Properties [][]string `ch:"properties"`

	AllowFailure bool `ch:"allow_failure"`
	Manual       bool `ch:"manual"`
	Retried      bool `ch:"retried"`
	Retryable    bool `ch:"retryable"`

	Kind                        string `ch:"kind"`
	DownstreamPipelineId        int64  `ch:"downstream_pipeline_id"`
	DownstreamPipelineIid       int64  `ch:"downstream_pipeline_iid"`
	DownstreamPipelineProjectId int64  `ch:"downstream_pipeline_project_id"`

	RunnerId string `ch:"runner_id"`

	// deprecated
	Pipeline []any `ch:"pipeline"` // Tuple(id Int64, project_id Int64, ref String, sha String, status String)
}

type Section struct {
	Id         int64 `ch:"id"`
	JobId      int64 `ch:"job_id"`
	PipelineId int64 `ch:"pipeline_id"`
	ProjectId  int64 `ch:"project_id"`

	Name string `ch:"name"`

	StartedAt  float64 `ch:"started_at"`
	FinishedAt float64 `ch:"finished_at"`

	Duration float64 `ch:"duration"`

	// deprecated
	Job      []any `ch:"job"`      // Tuple(id Int64, name String, status String)
	Pipeline []any `ch:"pipeline"` // Tuple(id Int64, project_id Int64, ref String, sha String, status String)
}

type TestReport struct {
	Id         string `ch:"id"`
	JobId      int64  `ch:"job_id"`
	PipelineId int64  `ch:"pipeline_id"`
	ProjectId  int64  `ch:"project_id"`

	TotalTime    float64 `ch:"total_time"`
	TotalCount   int64   `ch:"total_count"`
	ErrorCount   int64   `ch:"error_count"`
	FailedCount  int64   `ch:"failed_count"`
	SkippedCount int64   `ch:"skipped_count"`
	SuccessCount int64   `ch:"success_count"`
}

type TestSuite struct {
	Id           string `ch:"id"`
	TestReportId string `ch:"testreport_id"`
	JobId        int64  `ch:"job_id"`
	PipelineId   int64  `ch:"pipeline_id"`
	ProjectId    int64  `ch:"project_id"`

	Name         string  `ch:"name"`
	TotalTime    float64 `ch:"total_time"`
	TotalCount   int64   `ch:"total_count"`
	ErrorCount   int64   `ch:"error_count"`
	FailedCount  int64   `ch:"failed_count"`
	SkippedCount int64   `ch:"skipped_count"`
	SuccessCount int64   `ch:"success_count"`

	Properties [][]string `ch:"properties"`
}

type TestCase struct {
	Id           string `ch:"id"`
	TestSuiteId  string `ch:"testsuite_id"`
	TestReportId string `ch:"testreport_id"`
	JobId        int64  `ch:"job_id"`
	PipelineId   int64  `ch:"pipeline_id"`
	ProjectId    int64  `ch:"project_id"`

	Status        string  `ch:"status"`
	Name          string  `ch:"name"`
	Classname     string  `ch:"classname"`
	File          string  `ch:"file"`
	ExecutionTime float64 `ch:"execution_time"`
	SystemOutput  string  `ch:"system_output"`
	AttachmentUrl string  `ch:"attachment_url"`

	Properties [][]string `ch:"properties"`

	ReportCreatedAt uint32 `ch:"report_created_at"`
}

type Metric struct {
	Id         string `ch:"id"`
	Iid        int64  `ch:"iid"`
	JobId      int64  `ch:"job_id"`
	PipelineId int64  `ch:"pipeline_id"`
	ProjectId  int64  `ch:"project_id"`

	Name      string            `ch:"name"`
	Labels    map[string]string `ch:"labels"`
	Value     float64           `ch:"value"`
	Timestamp int64             `ch:"timestamp"`
}

type MergeRequest struct {
	Id        int64 `ch:"id"`
	Iid       int64 `ch:"iid"`
	ProjectId int64 `ch:"project_id"`

	CreatedAt float64 `ch:"created_at"`
	UpdatedAt float64 `ch:"updated_at"`
	MergedAt  float64 `ch:"merged_at"`
	ClosedAt  float64 `ch:"closed_at"`

	Name        string   `ch:"name"`
	Title       string   `ch:"title"`
	Description string   `ch:"description"`
	Labels      []string `ch:"labels"`

	State       string `ch:"state"`
	MergeStatus string `ch:"merge_status"`
	MergeError  string `ch:"merge_error"`

	SourceProjectId int64  `ch:"source_project_id"`
	SourceBranch    string `ch:"source_branch"`
	TargetProjectId int64  `ch:"target_project_id"`
	TargetBranch    string `ch:"target_branch"`

	Additions   int64 `ch:"additions"`
	Changes     int64 `ch:"changes"`
	Deletions   int64 `ch:"deletions"`
	FileCount   int64 `ch:"file_count"`
	CommitCount int64 `ch:"commit_count"`

	BaseSha         string `ch:"base_sha"`
	HeadSha         string `ch:"head_sha"`
	StartSha        string `ch:"start_sha"`
	MergeCommitSha  string `ch:"merge_commit_sha"`
	RebaseCommitSha string `ch:"rebase_commit_sha"`

	AuthorId          int64    `ch:"author_id"`
	AuthorUsername    string   `ch:"author_username"`
	AuthorName        string   `ch:"author_name"`
	AssigneesId       []int64  `ch:"assignees_id"`
	AssigneesUsername []string `ch:"assignees_username"`
	AssigneesName     []string `ch:"assignees_name"`
	ReviewersId       []int64  `ch:"reviewers_id"`
	ReviewersUsername []string `ch:"reviewers_username"`
	ReviewersName     []string `ch:"reviewers_name"`
	ApproversId       []int64  `ch:"approvers_id"`
	ApproversUsername []string `ch:"approvers_username"`
	ApproversName     []string `ch:"approvers_name"`
	MergeUserId       int64    `ch:"merge_user_id"`
	MergeUserUsername string   `ch:"merge_user_username"`
	MergeUserName     string   `ch:"merge_user_name"`

	CommitShas []string `ch:"commit_shas"`

	Approved  bool `ch:"approved"`
	Conflicts bool `ch:"conflicts"`
	Draft     bool `ch:"draft"`
	Mergeable bool `ch:"mergeable"`

	MilestoneId        int64 `ch:"milestone_id"`
	MilestoneIid       int64 `ch:"milestone_iid"`
	MilestoneProjectId int64 `ch:"milestone_project_id"`
}

type MergeRequestCommit struct {
	Id              string `ch:"id"`
	MergeRequestId  int64  `ch:"mergerequest_id"`
	MergeRequestIid int64  `ch:"mergerequest_iid"`
	ProjectId       int64  `ch:"project_id"`

	Sha string `ch:"sha"`

	Title    string     `ch:"title"`
	Message  string     `ch:"message"`
	Trailers [][]string `ch:"trailers"`

	AuthorId       int64  `ch:"author_id"`
	AuthorUsername string `ch:"author_username"`

	AuthoredDate  time.Time `ch:"authored_date"`
	CommittedDate time.Time `ch:"committed_date"`

	AuthorName     string `ch:"author_name"`
	AuthorEmail    string `ch:"author_email"`
	CommitterName  string `ch:"committer_name"`
	CommitterEmail string `ch:"committer_email"`
}

type MergeRequestNoteEvent struct {
	Id                    int64 `ch:"id"`
	MergeRequestId        int64 `ch:"mergerequest_id"`
	MergeRequestIid       int64 `ch:"mergerequest_iid"`
	MergeRequestProjectId int64 `ch:"mergerequest_project_id"`

	CreatedAt  float64 `ch:"created_at"`
	UpdatedAt  float64 `ch:"updated_at"`
	ResolvedAt float64 `ch:"resolved_at"`

	Type     string `ch:"type"`
	System   bool   `ch:"system"`
	Internal bool   `ch:"internal"`

	AuthorId       int64  `ch:"author_id"`
	AuthorUsername string `ch:"author_username"`
	AuthorName     string `ch:"author_name"`

	Resolvable       bool   `ch:"resolvable"`
	Resolved         bool   `ch:"resolved"`
	ResolverId       int64  `ch:"resolver_id"`
	ResolverUsername string `ch:"resolver_username"`
	ResolverName     string `ch:"resolver_name"`
}

type CoverageReport struct {
	Id         string `ch:"id"`
	JobId      int64  `ch:"job_id"`
	PipelineId int64  `ch:"pipeline_id"`
	ProjectId  int64  `ch:"project_id"`

	LineRate     float32 `ch:"line_rate"`
	LinesCovered int32   `ch:"lines_covered"`
	LinesValid   int32   `ch:"lines_valid"`

	BranchRate      float32 `ch:"branch_rate"`
	BranchesCovered int32   `ch:"branches_covered"`
	BranchesValid   int32   `ch:"branches_valid"`

	Complexity float32 `ch:"complexity"`

	Version   string `ch:"version"`
	Timestamp int64  `ch:"timestamp"`

	SourcePaths []string `ch:"source_paths"`
}

type CoveragePackage struct {
	Id         string `ch:"id"`
	ReportId   string `ch:"report_id"`
	JobId      int64  `ch:"job_id"`
	PipelineId int64  `ch:"pipeline_id"`
	ProjectId  int64  `ch:"project_id"`

	Name string `ch:"name"`

	LineRate   float32 `ch:"line_rate"`
	BranchRate float32 `ch:"branch_rate"`
	Complexity float32 `ch:"complexity"`
}

type CoverageClass struct {
	Id         string `ch:"id"`
	PackageId  string `ch:"package_id"`
	ReportId   string `ch:"report_id"`
	JobId      int64  `ch:"job_id"`
	PipelineId int64  `ch:"pipeline_id"`
	ProjectId  int64  `ch:"project_id"`

	PackageName string `ch:"package_name"`
	Name        string `ch:"name"`
	Filename    string `ch:"filename"`

	LineRate   float32 `ch:"line_rate"`
	BranchRate float32 `ch:"branch_rate"`
	Complexity float32 `ch:"complexity"`
}

type CoverageMethod struct {
	Id         string `ch:"id"`
	ClassId    string `ch:"class_id"`
	PackageId  string `ch:"package_id"`
	ReportId   string `ch:"report_id"`
	JobId      int64  `ch:"job_id"`
	PipelineId int64  `ch:"pipeline_id"`
	ProjectId  int64  `ch:"project_id"`

	PackageName string `ch:"package_name"`
	ClassName   string `ch:"class_name"`
	Name        string `ch:"name"`
	Signature   string `ch:"signature"`

	LineRate   float32 `ch:"line_rate"`
	BranchRate float32 `ch:"branch_rate"`
	Complexity float32 `ch:"complexity"`
}

type Deployment struct {
	Id  int64 `ch:"id"`
	Iid int64 `ch:"iid"`

	EnvironmentId   int64  `ch:"environment_id"`
	EnvironmentName string `ch:"environment_name"`
	EnvironmentTier string `ch:"environment_tier"`

	ProjectId int64 `ch:"project_id"`

	JobId      int64 `ch:"job_id"`
	PipelineId int64 `ch:"pipeline_id"`

	TriggererId       int64  `ch:"triggerer_id"`
	TriggererUsername string `ch:"triggerer_username"`
	TriggererName     string `ch:"triggerer_name"`

	CreatedAt  float64 `ch:"created_at"`
	FinishedAt float64 `ch:"finished_at"`
	UpdatedAt  float64 `ch:"updated_at"`

	Status string `ch:"status"`
	Ref    string `ch:"ref"`
	Sha    string `ch:"sha"`
}

type Runner struct {
	Id          int64  `ch:"id"`
	ShortSha    string `ch:"short_sha"`
	Description string `ch:"description"`

	RunnerType string   `ch:"runner_type"`
	TagList    []string `ch:"tag_list"`
	Status     string   `ch:"status"`

	Locked bool `ch:"locked"`
	Paused bool `ch:"paused"`

	RunProtected bool `ch:"run_protected"`
	RunUntagged  bool `ch:"run_untagged"`

	CreatedAt   float64 `ch:"created_at"`
	ContactedAt float64 `ch:"contacted_at"`

	CreatedById       int64  `ch:"created_by_id"`
	CreatedByUsername string `ch:"created_by_username"`
	CreatedByName     string `ch:"created_by_name"`
}
