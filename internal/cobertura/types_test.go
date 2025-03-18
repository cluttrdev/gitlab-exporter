package cobertura_test

import (
	"strings"
	"testing"

	"github.com/cluttrdev/gitlab-exporter/internal/cobertura"
	"github.com/cluttrdev/gitlab-exporter/internal/types"
	"github.com/google/go-cmp/cmp"
)

const data string = `
    <?xml version="1.0" encoding="UTF-8"?>
    <!DOCTYPE coverage SYSTEM "http://cobertura.sourceforge.net/xml/coverage-04.dtd">
    <coverage line-rate="0.013707013" branch-rate="0" version="" timestamp="1742305575745" lines-covered="224" lines-valid="16342" branches-covered="0" branches-valid="0" complexity="0">
      <sources>
        <source>/builds/akun73/gitlab-exporter</source>
      </sources>
      <packages>
        <package name="github.com/cluttrdev/gitlab-exporter/internal/cobertura" line-rate="0.04255319" branch-rate="0" complexity="0">
          <classes>
            <class name="-" filename="internal/cobertura/parse.go" line-rate="0.5714286" branch-rate="0" complexity="0">
              <methods>
                <method name="Parse" signature="" line-rate="0.5714286" branch-rate="0" complexity="0">
                </method>
              </methods>
            </class>
            <class name="-" filename="internal/cobertura/types.go" line-rate="0" branch-rate="0" complexity="0">
              <methods>
                <method name="ConvertCoverageReport" signature="" line-rate="0" branch-rate="0" complexity="0">
                </method>
              </methods>
            </class>
          </classes>
        </package>
        <package name="github.com/cluttrdev/gitlab-exporter/internal/junitxml" line-rate="0" branch-rate="0" complexity="0">
          <classes>
            <class name="-" filename="internal/junitxml/parse.go" line-rate="0" branch-rate="0" complexity="0">
              <methods>
                <method name="Parse" signature="" line-rate="0" branch-rate="0" complexity="0">
                </method>
                <method name="ParseTextAttachments" signature="" line-rate="0" branch-rate="0" complexity="0">
                </method>
                <method name="ParseTextProperties" signature="" line-rate="0" branch-rate="0" complexity="0">
                </method>
              </methods>
            </class>
            <class name="Failure" filename="internal/junitxml/types.go" line-rate="0" branch-rate="0" complexity="0">
              <methods>
                <method name="Output" signature="" line-rate="0" branch-rate="0" complexity="0">
                </method>
              </methods>
            </class>
            <class name="Error" filename="internal/junitxml/types.go" line-rate="0" branch-rate="0" complexity="0">
              <methods>
                <method name="Output" signature="" line-rate="0" branch-rate="0" complexity="0">
                </method>
              </methods>
            </class>
            <class name="-" filename="internal/junitxml/types.go" line-rate="0" branch-rate="0" complexity="0">
              <methods>
                <method name="ConvertTestReport" signature="" line-rate="0" branch-rate="0" complexity="0">
                </method>
                <method name="convertTestSuites" signature="" line-rate="0" branch-rate="0" complexity="0">
                </method>
                <method name="convertTestCases" signature="" line-rate="0" branch-rate="0" complexity="0">
                </method>
                <method name="convertTestProperties" signature="" line-rate="0" branch-rate="0" complexity="0">
                </method>
                <method name="formatTestOutput" signature="" line-rate="0" branch-rate="0" complexity="0">
                </method>
              </methods>
            </class>
          </classes>
        </package>
      </packages>
    </coverage>
`

func TestConvert(t *testing.T) {
	r, err := cobertura.Parse(strings.NewReader(data))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	report, packages, classes, methods := cobertura.ConvertCoverageReport(0, r, types.JobReference{})

	wantReport := types.CoverageReport{
		Id: "0-0",

		LineRate:     0.013707013,
		LinesCovered: 224,
		LinesValid:   16342,

		BranchRate:      0.0,
		BranchesCovered: 0,
		BranchesValid:   0,

		Complexity: 0.0,

		Version:   "",
		Timestamp: 1742305575,

		SourcePaths: []string{"/builds/akun73/gitlab-exporter"},
	}

	if diff := cmp.Diff(wantReport, report); diff != "" {
		t.Errorf("Report mismatch (-want +got):\n%s", diff)
	}

	wantPackages := []types.CoveragePackage{
		{
			Id: "0-0-0",
			Report: types.CoverageReportReference{
				Id: "0-0",
			},

			Name:       "github.com/cluttrdev/gitlab-exporter/internal/cobertura",
			LineRate:   0.04255319,
			BranchRate: 0.0,
			Complexity: 0.0,
		},
		{
			Id: "0-0-1",
			Report: types.CoverageReportReference{
				Id: "0-0",
			},

			Name:       "github.com/cluttrdev/gitlab-exporter/internal/junitxml",
			LineRate:   0.0,
			BranchRate: 0.0,
			Complexity: 0.0,
		},
	}

	if diff := cmp.Diff(wantPackages, packages); diff != "" {
		t.Errorf("Packages mismatch (-want +got):\n%s", diff)
	}

	wantClasses := []types.CoverageClass{
		{
			Id: "0-0-0-0", Package: types.CoveragePackageReference{Id: "0-0-0", Name: "github.com/cluttrdev/gitlab-exporter/internal/cobertura", Report: types.CoverageReportReference{Id: "0-0"}},
			Name: "-", Filename: "internal/cobertura/parse.go",
			LineRate: 0.5714286, BranchRate: 0.0, Complexity: 0.0,
		},
		{
			Id: "0-0-0-1", Package: types.CoveragePackageReference{Id: "0-0-0", Name: "github.com/cluttrdev/gitlab-exporter/internal/cobertura", Report: types.CoverageReportReference{Id: "0-0"}},
			Name: "-", Filename: "internal/cobertura/types.go",
			LineRate: 0.0, BranchRate: 0.0, Complexity: 0.0,
		},
		{
			Id: "0-0-1-0", Package: types.CoveragePackageReference{Id: "0-0-1", Name: "github.com/cluttrdev/gitlab-exporter/internal/junitxml", Report: types.CoverageReportReference{Id: "0-0"}},
			Name: "-", Filename: "internal/junitxml/parse.go",
			LineRate: 0.0, BranchRate: 0.0, Complexity: 0.0,
		},
		{
			Id: "0-0-1-1", Package: types.CoveragePackageReference{Id: "0-0-1", Name: "github.com/cluttrdev/gitlab-exporter/internal/junitxml", Report: types.CoverageReportReference{Id: "0-0"}},
			Name: "Failure", Filename: "internal/junitxml/types.go",
			LineRate: 0.0, BranchRate: 0.0, Complexity: 0.0,
		},
		{
			Id: "0-0-1-2", Package: types.CoveragePackageReference{Id: "0-0-1", Name: "github.com/cluttrdev/gitlab-exporter/internal/junitxml", Report: types.CoverageReportReference{Id: "0-0"}},
			Name: "Error", Filename: "internal/junitxml/types.go",
			LineRate: 0.0, BranchRate: 0.0, Complexity: 0.0,
		},
		{
			Id: "0-0-1-3", Package: types.CoveragePackageReference{Id: "0-0-1", Name: "github.com/cluttrdev/gitlab-exporter/internal/junitxml", Report: types.CoverageReportReference{Id: "0-0"}},
			Name: "-", Filename: "internal/junitxml/types.go",
			LineRate: 0.0, BranchRate: 0.0, Complexity: 0.0,
		},
	}

	if diff := cmp.Diff(wantClasses, classes); diff != "" {
		t.Errorf("Classes mismatch (-want +got):\n%s", diff)
	}

	wantMethods := []types.CoverageMethod{
		{
			Id: "0-0-0-0-0", Class: types.CoverageClassReference{Id: "0-0-0-0", Name: "-", Package: types.CoveragePackageReference{Id: "0-0-0", Name: "github.com/cluttrdev/gitlab-exporter/internal/cobertura", Report: types.CoverageReportReference{Id: "0-0"}}},
			Name: "Parse", Signature: "",
			LineRate: 0.5714286, BranchRate: 0.0, Complexity: 0.0,
		},
		{
			Id: "0-0-0-1-0", Class: types.CoverageClassReference{Id: "0-0-0-1", Name: "-", Package: types.CoveragePackageReference{Id: "0-0-0", Name: "github.com/cluttrdev/gitlab-exporter/internal/cobertura", Report: types.CoverageReportReference{Id: "0-0"}}},
			Name: "ConvertCoverageReport", Signature: "",
			LineRate: 0.0, BranchRate: 0.0, Complexity: 0.0,
		},
		{
			Id: "0-0-1-0-0", Class: types.CoverageClassReference{Id: "0-0-1-0", Name: "-", Package: types.CoveragePackageReference{Id: "0-0-1", Name: "github.com/cluttrdev/gitlab-exporter/internal/junitxml", Report: types.CoverageReportReference{Id: "0-0"}}},
			Name: "Parse", Signature: "",
			LineRate: 0.0, BranchRate: 0.0, Complexity: 0.0,
		},
		{
			Id: "0-0-1-0-1", Class: types.CoverageClassReference{Id: "0-0-1-0", Name: "-", Package: types.CoveragePackageReference{Id: "0-0-1", Name: "github.com/cluttrdev/gitlab-exporter/internal/junitxml", Report: types.CoverageReportReference{Id: "0-0"}}},
			Name: "ParseTextAttachments", Signature: "",
			LineRate: 0.0, BranchRate: 0.0, Complexity: 0.0,
		},
		{
			Id: "0-0-1-0-2", Class: types.CoverageClassReference{Id: "0-0-1-0", Name: "-", Package: types.CoveragePackageReference{Id: "0-0-1", Name: "github.com/cluttrdev/gitlab-exporter/internal/junitxml", Report: types.CoverageReportReference{Id: "0-0"}}},
			Name: "ParseTextProperties", Signature: "",
			LineRate: 0.0, BranchRate: 0.0, Complexity: 0.0,
		},
		{
			Id: "0-0-1-1-0", Class: types.CoverageClassReference{Id: "0-0-1-1", Name: "Failure", Package: types.CoveragePackageReference{Id: "0-0-1", Name: "github.com/cluttrdev/gitlab-exporter/internal/junitxml", Report: types.CoverageReportReference{Id: "0-0"}}},
			Name: "Output", Signature: "",
			LineRate: 0.0, BranchRate: 0.0, Complexity: 0.0,
		},
		{
			Id: "0-0-1-2-0", Class: types.CoverageClassReference{Id: "0-0-1-2", Name: "Error", Package: types.CoveragePackageReference{Id: "0-0-1", Name: "github.com/cluttrdev/gitlab-exporter/internal/junitxml", Report: types.CoverageReportReference{Id: "0-0"}}},
			Name: "Output", Signature: "",
			LineRate: 0.0, BranchRate: 0.0, Complexity: 0.0,
		},
		{
			Id: "0-0-1-3-0", Class: types.CoverageClassReference{Id: "0-0-1-3", Name: "-", Package: types.CoveragePackageReference{Id: "0-0-1", Name: "github.com/cluttrdev/gitlab-exporter/internal/junitxml", Report: types.CoverageReportReference{Id: "0-0"}}},
			Name: "ConvertTestReport", Signature: "",
			LineRate: 0.0, BranchRate: 0.0, Complexity: 0.0,
		},
		{
			Id: "0-0-1-3-1", Class: types.CoverageClassReference{Id: "0-0-1-3", Name: "-", Package: types.CoveragePackageReference{Id: "0-0-1", Name: "github.com/cluttrdev/gitlab-exporter/internal/junitxml", Report: types.CoverageReportReference{Id: "0-0"}}},
			Name: "convertTestSuites", Signature: "",
			LineRate: 0.0, BranchRate: 0.0, Complexity: 0.0,
		},
		{
			Id: "0-0-1-3-2", Class: types.CoverageClassReference{Id: "0-0-1-3", Name: "-", Package: types.CoveragePackageReference{Id: "0-0-1", Name: "github.com/cluttrdev/gitlab-exporter/internal/junitxml", Report: types.CoverageReportReference{Id: "0-0"}}},
			Name: "convertTestCases", Signature: "",
			LineRate: 0.0, BranchRate: 0.0, Complexity: 0.0,
		},
		{
			Id: "0-0-1-3-3", Class: types.CoverageClassReference{Id: "0-0-1-3", Name: "-", Package: types.CoveragePackageReference{Id: "0-0-1", Name: "github.com/cluttrdev/gitlab-exporter/internal/junitxml", Report: types.CoverageReportReference{Id: "0-0"}}},
			Name: "convertTestProperties", Signature: "",
			LineRate: 0.0, BranchRate: 0.0, Complexity: 0.0,
		},
		{
			Id: "0-0-1-3-4", Class: types.CoverageClassReference{Id: "0-0-1-3", Name: "-", Package: types.CoveragePackageReference{Id: "0-0-1", Name: "github.com/cluttrdev/gitlab-exporter/internal/junitxml", Report: types.CoverageReportReference{Id: "0-0"}}},
			Name: "formatTestOutput", Signature: "",
			LineRate: 0.0, BranchRate: 0.0, Complexity: 0.0,
		},
	}

	if diff := cmp.Diff(wantMethods, methods); diff != "" {
		t.Errorf("Methods mismatch (-want +got):\n%s", diff)
	}

}
