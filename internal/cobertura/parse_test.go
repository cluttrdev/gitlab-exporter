package cobertura_test

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/cluttrdev/gitlab-exporter/internal/cobertura"
)

func TestParse(t *testing.T) {
	const data string = `
        <?xml version='1.0' encoding='UTF-8'?>
        <!DOCTYPE coverage SYSTEM 'http://cobertura.sourceforge.net/xml/coverage-04.dtd'>
        <coverage line-rate="0.8571428571428571" branch-rate="0.5" lines-covered="6" lines-valid="7" branches-covered="1" branches-valid="2" complexity="0.0" timestamp="1728808696" version="gcovr 8.3">
          <sources>
            <source>doc/examples</source>
          </sources>
          <packages>
            <package name="" line-rate="0.8571428571428571" branch-rate="0.5" complexity="0.0">
              <classes>
                <class name="example_cpp" filename="example.cpp" line-rate="0.8571428571428571" branch-rate="0.5" complexity="0.0">
                  <methods/>
                  <lines>
                    <line number="3" hits="1" branch="false"/>
                    <line number="5" hits="1" branch="true" condition-coverage="50% (1/2)">
                      <conditions>
                        <condition number="0" type="jump" coverage="50%"/>
                      </conditions>
                    </line>
                    <line number="7" hits="0" branch="false"/>
                    <line number="11" hits="1" branch="false"/>
                    <line number="15" hits="1" branch="false"/>
                    <line number="17" hits="1" branch="false"/>
                    <line number="19" hits="1" branch="false"/>
                  </lines>
                </class>
              </classes>
            </package>
          </packages>
        </coverage>
    `

	report, err := cobertura.Parse(strings.NewReader(data))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	want := cobertura.CoverageReport{
		XMLName:  xml.Name{Local: "coverage"},
		LineRate: 0.8571428571428571, BranchRate: 0.5, LinesCovered: 6, LinesValid: 7, BranchesCovered: 1, BranchesValid: 2, Complexity: 0.0, Version: "gcovr 8.3", Timestamp: 1728808696,
		Sources: []cobertura.Source{
			{Path: "doc/examples"},
		},
		Packages: []cobertura.Package{
			{
				Name: "", LineRate: 0.8571428571428571, BranchRate: 0.5, Complexity: 0.0,
				Classes: []cobertura.Class{
					{
						Name: "example_cpp", Filename: "example.cpp", LineRate: 0.8571428571428571, BranchRate: 0.5, Complexity: 0.0,
						Methods: nil,
						Lines: []cobertura.Line{
							{Number: 3, Hits: 1, Branch: false, ConditionCoverage: "", Conditions: nil},
							{
								Number: 5, Hits: 1, Branch: true, ConditionCoverage: "50% (1/2)",
								Conditions: []cobertura.Condition{
									{Number: 0, Type: "jump", Coverage: "50%"},
								},
							},
							{Number: 7, Hits: 0, Branch: false},
							{Number: 11, Hits: 1, Branch: false},
							{Number: 15, Hits: 1, Branch: false},
							{Number: 17, Hits: 1, Branch: false},
							{Number: 19, Hits: 1, Branch: false},
						},
					},
				},
			},
		},
	}

	if diff := cmp.Diff(want, report); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}
