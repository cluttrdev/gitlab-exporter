package gitlab

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"regexp"
	"sort"
	"strconv"

	_gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-exporter/internal/models"
	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

type ListJobSectionsResult struct {
	Section *typespb.Section
	Error   error
}

func (c *Client) ListJobSections(ctx context.Context, projectID int64, jobID int64) <-chan ListJobSectionsResult {
	ch := make(chan ListJobSectionsResult)

	go func() {
		defer close(ch)

		c.RLock()
		job, _, err := c.client.Jobs.GetJob(int(projectID), int(jobID), _gitlab.WithContext(ctx))
		c.RUnlock()
		if err != nil {
			ch <- ListJobSectionsResult{
				Error: err,
			}
			return
		}

		c.RLock()
		trace, _, err := c.client.Jobs.GetTraceFile(int(projectID), int(jobID), _gitlab.WithContext(ctx))
		c.RUnlock()
		if err != nil {
			ch <- ListJobSectionsResult{
				Error: err,
			}
			return
		}

		data, err := parseSections(trace)
		if err != nil {
			ch <- ListJobSectionsResult{
				Error: err,
			}
			return
		}

		for secnum, secdat := range data {
			section := &typespb.Section{
				Name:       secdat.Name,
				StartedAt:  models.ConvertUnixSeconds(secdat.Start),
				FinishedAt: models.ConvertUnixSeconds(secdat.End),
				Duration:   models.ConvertDuration(float64(secdat.End - secdat.Start)),
			}

			section.Id = int64(job.ID*1000 + secnum)
			section.Job.Id = int64(job.ID)
			section.Job.Name = job.Name
			section.Job.Status = job.Status
			section.Pipeline.Id = int64(job.Pipeline.ID)
			section.Pipeline.ProjectId = int64(job.Pipeline.ProjectID)
			section.Pipeline.Ref = job.Pipeline.Ref
			section.Pipeline.Sha = job.Pipeline.Sha
			section.Pipeline.Status = job.Pipeline.Status

			ch <- ListJobSectionsResult{
				Section: section,
			}
		}
	}()

	return ch
}

type SectionData struct {
	Name  string `json:"name"`
	Start int64  `json:"start"`
	End   int64  `json:"end"`
}

func parseSections(trace *bytes.Reader) ([]SectionData, error) {
	sections := []SectionData{}
	stack := sectionStack{}

	scanner := bufio.NewScanner(trace)
	for scanner.Scan() {
		line := scanner.Bytes()
		var i, j int
		sep := []byte(`section_`)
		for {
			j = bytes.Index(line[i:], sep)
			if j < 0 {
				break
			}

			marker, ts, name, err := parseSection(line[i:])
			if err != nil {
				// TODO: what?
			} else if marker == string(sectionMarkerStart) {
				stack.Start(ts, name)
			} else if marker == string(sectionMarkerEnd) {
				sections = append(sections, stack.End(ts, name)...)
			}

			i = i + j + 1
		}
	}

	sort.SliceStable(sections, func(i, j int) bool {
		return sections[i].Start < sections[j].Start
	})

	return sections, nil
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
