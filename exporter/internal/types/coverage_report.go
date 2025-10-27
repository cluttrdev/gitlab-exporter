package types

type CoverageReport struct {
	Id  string
	Job JobReference

	LineRate     float32
	LinesCovered int32
	LinesValid   int32

	BranchRate      float32
	BranchesCovered int32
	BranchesValid   int32

	Complexity float32

	Version   string
	Timestamp int64

	SourcePaths []string
}

type CoverageReportReference struct {
	Id  string
	Job JobReference
}

type CoverageSource struct {
	Report CoverageReportReference

	Path string
}

type CoveragePackage struct {
	Id     string
	Report CoverageReportReference

	Name       string
	LineRate   float32
	BranchRate float32
	Complexity float32
}

type CoveragePackageReference struct {
	Id   string
	Name string

	Report CoverageReportReference
}

type CoverageClass struct {
	Id      string
	Package CoveragePackageReference

	Name       string
	Filename   string
	LineRate   float32
	BranchRate float32
	Complexity float32
}

type CoverageClassReference struct {
	Id   string
	Name string

	Package CoveragePackageReference
}

type CoverageMethod struct {
	Id    string
	Class CoverageClassReference

	Name       string
	Signature  string
	LineRate   float32
	BranchRate float32
	Complexity float32
}

type CoverageLine struct {
	ClassId   string
	MethodId  string
	PackageId string
	Report    CoverageReportReference

	Number int32
	Hits   int32
	Branch bool

	ConditionCoverage float32
	ConditionsCovered int32
	ConditionsValid   int32
}
