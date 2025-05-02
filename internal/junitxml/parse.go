package junitxml

import (
	"encoding/xml"
	"errors"
	"io"
	"regexp"
)

func Parse(r io.Reader) (TestReport, error) {
	var report TestReport

	if err := xml.NewDecoder(r).Decode(&report); err != nil {
		return TestReport{}, err
	}
	return report, nil
}

func ParseMany(r io.Reader) ([]TestReport, error) {
	var reports []TestReport

	decoder := xml.NewDecoder(r)
	for {
		tok, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}

		if elem, ok := tok.(xml.StartElement); ok {
			var report TestReport
			if err := decoder.DecodeElement(&report, &elem); err != nil {
				return nil, err
			}
			reports = append(reports, report)
		}
	}

	return reports, nil
}

const (
	AttachmentRegex = `\[\[ATTACHMENT\|(?<path>[^\[\]\|]+?)\]\]`
	PropertyRegex   = `\[\[PROPERTY\|(?<name>[^\[\]\|=]+)=(?<value>[^\[\]\|=]+)\]\]`

	// PropertyTextValueRegex = `\[\[PROPERTY\|(?<name>[^\[\]\|=]+)\]\]\n(?<text>(?:.*\n)+)\[\[/PROPERTY\]\]`
)

func ParseTextAttachments(s string) []string {
	var paths []string

	re := regexp.MustCompile(AttachmentRegex)
	for _, match := range re.FindAllStringSubmatch(s, -1) {
		paths = append(paths, match[1])
	}

	return paths
}

func ParseTextProperties(s string) []Property {
	var properties []Property

	re := regexp.MustCompile(PropertyRegex)
	for _, match := range re.FindAllStringSubmatch(s, -1) {
		properties = append(properties, Property{
			Name:  match[1],
			Value: match[2],
		})
	}

	return properties
}
