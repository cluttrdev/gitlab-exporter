package cobertura

import (
	"encoding/xml"
	"io"
)

func Parse(r io.Reader) (CoverageReport, error) {
	var report CoverageReport

	if err := xml.NewDecoder(r).Decode(&report); err != nil {
		return CoverageReport{}, err
	}

	return report, nil
}
