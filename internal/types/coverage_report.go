package types

import (
	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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

func ConvertCoverageReport(report CoverageReport) *typespb.CoverageReport {
	return &typespb.CoverageReport{
		Id:  report.Id,
		Job: ConvertJobReference(report.Job),

		LineRate:     report.LineRate,
		LinesCovered: report.LinesCovered,
		LinesValid:   report.LinesValid,

		BranchRate:      report.BranchRate,
		BranchesCovered: report.BranchesCovered,
		BranchesValid:   report.BranchesValid,

		Complexity: report.Complexity,

		Version:   report.Version,
		Timestamp: &timestamppb.Timestamp{Seconds: report.Timestamp},

		SourcePaths: report.SourcePaths,
	}
}

func ConvertCoverageReportReference(report CoverageReportReference) *typespb.CoverageReportReference {
	return &typespb.CoverageReportReference{
		Id:  report.Id,
		Job: ConvertJobReference(report.Job),
	}
}

func ConvertCoveragePackage(pkg CoveragePackage) *typespb.CoveragePackage {
	return &typespb.CoveragePackage{
		Id:     pkg.Id,
		Report: ConvertCoverageReportReference(pkg.Report),

		Name: pkg.Name,

		LineRate:   pkg.LineRate,
		BranchRate: pkg.BranchRate,
		Complexity: pkg.Complexity,
	}
}

func ConvertCoveragePackageReference(pkg CoveragePackageReference) *typespb.CoveragePackageReference {
	return &typespb.CoveragePackageReference{
		Id:   pkg.Id,
		Name: pkg.Name,

		Report: ConvertCoverageReportReference(pkg.Report),
	}
}

func ConvertCoverageClass(cls CoverageClass) *typespb.CoverageClass {
	return &typespb.CoverageClass{
		Id:      cls.Id,
		Package: ConvertCoveragePackageReference(cls.Package),

		Name:     cls.Name,
		Filename: cls.Filename,

		LineRate:   cls.LineRate,
		BranchRate: cls.BranchRate,
		Complexity: cls.Complexity,
	}
}

func ConvertCoverageClassReference(cls CoverageClassReference) *typespb.CoverageClassReference {
	return &typespb.CoverageClassReference{
		Id:   cls.Id,
		Name: cls.Name,

		Package: ConvertCoveragePackageReference(cls.Package),
	}
}

func ConvertCoverageMethod(mtd CoverageMethod) *typespb.CoverageMethod {
	return &typespb.CoverageMethod{
		Id:    mtd.Id,
		Class: ConvertCoverageClassReference(mtd.Class),

		Name:      mtd.Name,
		Signature: mtd.Signature,

		LineRate:   mtd.LineRate,
		BranchRate: mtd.BranchRate,
		Complexity: mtd.Complexity,
	}
}
