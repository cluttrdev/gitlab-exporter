package rest

import (
	"errors"
	"regexp"
	"strconv"
)

type SectionData struct {
	Name  string `json:"name"`
	Start int64  `json:"start"`
	End   int64  `json:"end"`
}

type sectionStack struct {
	Sections []SectionData
}

func (s *sectionStack) Len() int {
	return len(s.Sections)
}

func (s *sectionStack) Push(section SectionData) {
	s.Sections = append(s.Sections, section)
}

func (s *sectionStack) Pop() SectionData {
	size := len(s.Sections)
	if size == 0 {
		return SectionData{}
	}
	section := s.Sections[size-1]
	s.Sections = s.Sections[:(size - 1)]
	return section
}

func (s *sectionStack) Start(timestamp int64, name string) {
	s.Push(SectionData{
		Name:  name,
		Start: timestamp,
	})
}

func (s *sectionStack) End(timestamp int64, name string) []SectionData {
	endedSections := []SectionData{}

	section := s.Pop()
	if section.Name == "" {
		return endedSections
	}

	if section.Name != name {
		endedSections = append(endedSections, s.End(timestamp, section.Name)...)
	}

	section.End = timestamp

	endedSections = append(endedSections, section)

	return endedSections
}

type sectionMarker string

const (
	sectionMarkerStart sectionMarker = "section_start"
	sectionMarkerEnd   sectionMarker = "section_end"
)

func parseSection(line []byte) (marker string, timestamp int64, name string, err error) {
	pattern := regexp.MustCompile(`(?P<marker>section_(?:start|end)):(?P<ts>\d+):(?P<name>[\w_]+)`)
	match := pattern.FindSubmatch(line)
	if len(match) != 4 {
		err = errors.New("no match found")
		return
	}

	marker = string(match[pattern.SubexpIndex("marker")])
	time_s := string(match[pattern.SubexpIndex("ts")])
	name = string(match[pattern.SubexpIndex("name")])

	timestamp, err = strconv.ParseInt(time_s, 10, 64)

	return
}
