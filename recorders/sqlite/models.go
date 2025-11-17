package sqlite

type Project struct {
	Id          int
	NamespaceId int

	Data []byte
}

type Pipeline struct {
	Id        int
	Iid       int
	ProjectId int

	Data []byte
}

type Job struct {
	Id         int
	PipelineId int
	ProjectId  int

	Data []byte
}

type Section struct {
	Id         int
	JobId      int
	PipelineId int
	ProjectId  int

	Data []byte
}

type Metric struct {
	Id         string
	Iid        int
	JobId      int
	PipelineId int
	ProjectId  int

	Data []byte
}

type CoverageReport struct {
	Id string

	JobId      int
	PipelineId int
	ProjectId  int

	Data []byte
}

type CoveragePackage struct {
	Id       string
	ReportId string

	JobId      int
	PipelineId int
	ProjectId  int

	Data []byte
}

type CoverageClass struct {
	Id        string
	PackageId string
	ReportId  string

	JobId      int
	PipelineId int
	ProjectId  int

	Data []byte
}

type CoverageMethod struct {
	Id        string
	ClassId   string
	PackageId string
	ReportId  string

	JobId      int
	PipelineId int
	ProjectId  int

	Data []byte
}

type TestReport struct {
	Id string

	JobId      int
	PipelineId int
	ProjectId  int

	Data []byte
}

type TestSuite struct {
	Id           string
	TestReportId string

	JobId      int
	PipelineId int
	ProjectId  int

	Data []byte
}

type TestCase struct {
	Id           string
	TestSuiteId  string
	TestReportId string

	JobId      int
	PipelineId int
	ProjectId  int

	Data []byte
}

type Deployment struct {
	Id            int
	Iid           int
	EnvironmentId int

	JobId      int
	PipelineId int
	ProjectId  int

	Data []byte
}

type Issue struct {
	Id        int
	Iid       int
	ProjectId int

	Data []byte
}

type MergeRequest struct {
	Id        int
	Iid       int
	ProjectId int

	Data []byte
}

type MergeRequestNoteEvent struct {
	Id                    int
	MergeRequestId        int
	MergeRequestIid       int
	MergeRequestProjectId int

	Data []byte
}

type Runner struct {
	Id int

	Data []byte
}

type TraceSpan struct {
	Timestamp          uint64 // Unix Nano
	TraceId            []byte
	SpanId             []byte
	ParentSpanId       []byte
	TraceState         string
	SpanName           string
	SpanKind           string
	ServiceName        string
	ResourceAttributes []byte // JSON object
	ScopeName          string
	ScopeVersion       string
	SpanAttributes     []byte // JSON object
	Duration           int64  // [ns]
	StatusCode         int32
	StatusMessage      string
	Events             []byte // JSON array of {Timestamp, Name, Attributes}
	Links              []byte // JSON array of {TraceId, SpanId, TraceState, Attributes}
}
