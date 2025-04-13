package messages

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.cluttr.dev/gitlab-exporter/internal/types"
	"go.cluttr.dev/gitlab-exporter/protobuf/typespb"
)

func NewCoverageReport(report types.CoverageReport) *typespb.CoverageReport {
	return &typespb.CoverageReport{
		Id:  report.Id,
		Job: NewJobReference(report.Job),

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

func NewCoverageReportReference(report types.CoverageReportReference) *typespb.CoverageReportReference {
	return &typespb.CoverageReportReference{
		Id:  report.Id,
		Job: NewJobReference(report.Job),
	}
}

func NewCoveragePackage(pkg types.CoveragePackage) *typespb.CoveragePackage {
	return &typespb.CoveragePackage{
		Id:     pkg.Id,
		Report: NewCoverageReportReference(pkg.Report),

		Name: pkg.Name,

		LineRate:   pkg.LineRate,
		BranchRate: pkg.BranchRate,
		Complexity: pkg.Complexity,
	}
}

func NewCoveragePackageReference(pkg types.CoveragePackageReference) *typespb.CoveragePackageReference {
	return &typespb.CoveragePackageReference{
		Id:   pkg.Id,
		Name: pkg.Name,

		Report: NewCoverageReportReference(pkg.Report),
	}
}

func NewCoverageClass(cls types.CoverageClass) *typespb.CoverageClass {
	return &typespb.CoverageClass{
		Id:      cls.Id,
		Package: NewCoveragePackageReference(cls.Package),

		Name:     cls.Name,
		Filename: cls.Filename,

		LineRate:   cls.LineRate,
		BranchRate: cls.BranchRate,
		Complexity: cls.Complexity,
	}
}

func NewCoverageClassReference(cls types.CoverageClassReference) *typespb.CoverageClassReference {
	return &typespb.CoverageClassReference{
		Id:   cls.Id,
		Name: cls.Name,

		Package: NewCoveragePackageReference(cls.Package),
	}
}

func NewCoverageMethod(mtd types.CoverageMethod) *typespb.CoverageMethod {
	return &typespb.CoverageMethod{
		Id:    mtd.Id,
		Class: NewCoverageClassReference(mtd.Class),

		Name:      mtd.Name,
		Signature: mtd.Signature,

		LineRate:   mtd.LineRate,
		BranchRate: mtd.BranchRate,
		Complexity: mtd.Complexity,
	}
}
