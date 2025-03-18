package cobertura

import (
	"encoding/xml"
	"fmt"

	"github.com/cluttrdev/gitlab-exporter/internal/types"
)

// https://github.com/gcovr/gcovr/blob/main/tests/cobertura.coverage-04.dtd

type CoverageReport struct {
	XMLName xml.Name `xml:"coverage"`

	LineRate        float32 `xml:"line-rate,attr"`
	BranchRate      float32 `xml:"branch-rate,attr"`
	LinesCovered    int32   `xml:"lines-covered,attr"`
	LinesValid      int32   `xml:"lines-valid,attr"`
	BranchesCovered int32   `xml:"branches-covered,attr"`
	BranchesValid   int32   `xml:"branches-valid,attr"`
	Complexity      float32 `xml:"complexity,attr"`
	Version         string  `xml:"version,attr"`
	Timestamp       int64   `xml:"timestamp,attr"`

	Sources  []Source  `xml:"sources>source"`
	Packages []Package `xml:"packages>package"`
}

type Source struct {
	Path string `xml:",chardata"`
}

type Package struct {
	Name       string  `xml:"name,attr"`
	LineRate   float32 `xml:"line-rate,attr"`
	BranchRate float32 `xml:"branch-rate,attr"`
	Complexity float32 `xml:"complexity,attr"`

	Classes []Class `xml:"classes>class"`
}

type Class struct {
	Name       string  `xml:"name,attr"`
	Filename   string  `xml:"filename,attr"`
	LineRate   float32 `xml:"line-rate,attr"`
	BranchRate float32 `xml:"branch-rate,attr"`
	Complexity float32 `xml:"complexity,attr"`

	Methods []Method `xml:"methods>method"`
	Lines   []Line   `xml:"lines>line"`
}

type Method struct {
	Name       string  `xml:"name,attr"`
	Signature  string  `xml:"signature,attr"`
	LineRate   float32 `xml:"line-rate,attr"`
	BranchRate float32 `xml:"branch-rate,attr"`
	Complexity float32 `xml:"complexity,attr"`

	Lines []Line `xml:"lines>line"`
}

type Line struct {
	Number            int32  `xml:"number,attr"`
	Hits              int32  `xml:"hits,attr"`
	Branch            bool   `xml:"branch,attr"`
	ConditionCoverage string `xml:"condition-coverage,attr"`

	Conditions []Condition `xml:"conditions>condition"`
}

type Condition struct {
	Number   int32  `xml:"number,attr"`
	Type     string `xml:"type,attr"`
	Coverage string `xml:"coverage,attr"`
}

func ConvertCoverageReport(index int, report CoverageReport, job types.JobReference) (types.CoverageReport, []types.CoveragePackage, []types.CoverageClass, []types.CoverageMethod) {
	covReportId := fmt.Sprintf("%d-%d", job.Id, index)
	covReport := types.CoverageReport{
		Id:  covReportId,
		Job: job,

		LineRate:     report.LineRate,
		LinesCovered: report.LinesCovered,
		LinesValid:   report.LinesValid,

		BranchRate:      report.BranchRate,
		BranchesCovered: report.BranchesCovered,
		BranchesValid:   report.BranchesValid,

		Complexity: report.Complexity,

		Version: report.Version,
		// Timestamp: report.Timestamp,
	}
	for _, source := range report.Sources {
		covReport.SourcePaths = append(covReport.SourcePaths, source.Path)
	}

	// The cobertura spec does not say whether the timestamp should be in seconds
	// or something else. We assume the timestamp represents a time in the range
	// between 2000-01-01T00:00:00Z and 3000-01-01T00:00:00Z
	ts := report.Timestamp
	switch {
	case 946684800 < ts && ts < 32503680000: // seconds
		covReport.Timestamp = ts
	case 946684800000 < ts && ts < 32503680000000: // milliseconds
		covReport.Timestamp = ts / 1_000
	case 946684800000000 < ts && ts < 32503680000000000: // microseconds
		covReport.Timestamp = ts / 1_000_000
	case 946684800000000000 < ts: // nanoseconds
		covReport.Timestamp = ts / 1_000_000_000
	}

	covReportRef := types.CoverageReportReference{
		Id:  covReportId,
		Job: job,
	}
	covPackages, covClasses, covMethods := convertCoveragePackages(report.Packages, covReportRef)

	return covReport, covPackages, covClasses, covMethods
}

func convertCoveragePackages(pkgs []Package, report types.CoverageReportReference) ([]types.CoveragePackage, []types.CoverageClass, []types.CoverageMethod) {
	var (
		covPackages []types.CoveragePackage
		covClasses  []types.CoverageClass
		covMethods  []types.CoverageMethod
	)

	for i, pkg := range pkgs {
		covPackageId := fmt.Sprintf("%s-%d", report.Id, i)
		covPackages = append(covPackages, types.CoveragePackage{
			Id:     covPackageId,
			Report: report,

			Name:       pkg.Name,
			LineRate:   pkg.LineRate,
			BranchRate: pkg.BranchRate,
			Complexity: pkg.Complexity,
		})

		covPackageRef := types.CoveragePackageReference{
			Id:     covPackageId,
			Name:   pkg.Name,
			Report: report,
		}
		ccs, cms := convertCoverageClasses(pkg.Classes, covPackageRef)
		covClasses = append(covClasses, ccs...)
		covMethods = append(covMethods, cms...)
	}

	return covPackages, covClasses, covMethods
}

func convertCoverageClasses(clss []Class, pkg types.CoveragePackageReference) ([]types.CoverageClass, []types.CoverageMethod) {
	var (
		covClasses []types.CoverageClass
		covMethods []types.CoverageMethod
	)

	for j, cls := range clss {
		covClassId := fmt.Sprintf("%s-%d", pkg.Id, j)
		covClasses = append(covClasses, types.CoverageClass{
			Id:      covClassId,
			Package: pkg,

			Name:       cls.Name,
			Filename:   cls.Filename,
			LineRate:   cls.LineRate,
			BranchRate: cls.BranchRate,
			Complexity: cls.Complexity,
		})

		covClassRef := types.CoverageClassReference{
			Id:      covClassId,
			Name:    cls.Name,
			Package: pkg,
		}
		cms := convertCoverageMethods(cls.Methods, covClassRef)
		covMethods = append(covMethods, cms...)
	}

	return covClasses, covMethods
}

func convertCoverageMethods(mtds []Method, cls types.CoverageClassReference) []types.CoverageMethod {
	var covMethods []types.CoverageMethod

	for k, mtd := range mtds {
		covMethodId := fmt.Sprintf("%s-%d", cls.Id, k)

		covMethods = append(covMethods, types.CoverageMethod{
			Id:    covMethodId,
			Class: cls,

			Name:       mtd.Name,
			Signature:  mtd.Signature,
			LineRate:   mtd.LineRate,
			BranchRate: mtd.BranchRate,
			Complexity: mtd.Complexity,
		})
	}

	return covMethods
}
