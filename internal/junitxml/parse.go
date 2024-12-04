package junitxml

import (
	"encoding/xml"
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
