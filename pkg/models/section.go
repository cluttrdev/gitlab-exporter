package models

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"time"
)

type Section struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Job  struct {
		ID     int64  `json:"id"`
		Name   string `json:"name"`
		Status string `json:"status"`
	} `json:"job"`
	Pipeline struct {
		ID        int64  `json:"id"`
		ProjectID int64  `json:"project_id"`
		Ref       string `json:"ref"`
		Sha       string `json:"sha"`
		Status    string `json:"status"`
	} `json:"pipeline"`
	StartedAt  *time.Time `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	Duration   float64    `json:"duration"`
}

type SectionStack struct {
	Sections []*Section
}

func (s *SectionStack) Len() int {
	return len(s.Sections)
}

func (s *SectionStack) Push(section *Section) {
	s.Sections = append(s.Sections, section)
}

func (s *SectionStack) Pop() *Section {
	size := len(s.Sections)
	if size == 0 {
		return nil
	}
	section := s.Sections[size-1]
	s.Sections = s.Sections[:(size - 1)]
	return section
}

func (s *SectionStack) Start(name string, start *time.Time) {
	s.Push(&Section{
		Name:      name,
		StartedAt: start,
	})
}

func (s *SectionStack) End(name string, end *time.Time) []*Section {
	endedSections := []*Section{}

	section := s.Pop()
	if section == nil {
		return endedSections
	}

	if section.Name != name {
		endedSections = append(endedSections, s.End(section.Name, end)...)
	}

	section.FinishedAt = end
	section.Duration = section.FinishedAt.Sub(*section.StartedAt).Seconds()

	endedSections = append(endedSections, section)

	return endedSections
}

func ParseSections(trace *bytes.Reader) ([]*Section, error) {
	pattern := regexp.MustCompile(`(?P<marker>section_(?:start|end)):(?P<ts>\d+):(?P<name>[\w_]+)`)

	text, err := io.ReadAll(trace)
	if err != nil {
		return nil, fmt.Errorf("[ParseSections] %w", err)
	}

	sections := []*Section{}
	stack := SectionStack{}
	for _, match := range pattern.FindAllSubmatch(text, -1) {
		marker := string(match[pattern.SubexpIndex("marker")])
		time_s := string(match[pattern.SubexpIndex("ts")])
		name := string(match[pattern.SubexpIndex("name")])

		time_i, err := strconv.ParseInt(time_s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("[ParseSections] %w", err)
		}
		time_t := time.Unix(time_i, 0)

		if marker == "section_start" {
			stack.Start(name, &time_t)
		} else if marker == "section_end" {
			sections = append(sections, stack.End(name, &time_t)...)
		} else {
			return nil, fmt.Errorf("Invalid section marker: %s", marker)
		}
	}

	sort.SliceStable(sections, func(i, j int) bool {
		return sections[i].StartedAt.Before(*sections[j].StartedAt)
	})

	return sections, nil
}
