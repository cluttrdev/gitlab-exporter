package gitlab

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"time"

	_gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-exporter/pkg/models"
)

func (c *Client) GetJobSections(ctx context.Context, projectID int64, jobID int64) ([]*models.Section, error) {
	sections := []*models.Section{}
	for r := range c.ListJobSections(ctx, projectID, jobID) {
		if r.Error != nil {
			return nil, fmt.Errorf("[gitlab.Client.GetSections] %w", r.Error)
		}
		sections = append(sections, r.Section)
	}
	return sections, nil
}

type ListJobSectionsResult struct {
	Section *models.Section
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

		unixTime := func(ts int64) *time.Time {
			const nsec int64 = 0
			t := time.Unix(ts, nsec)
			return &t
		}

		for secnum, secdat := range data {
			section := &models.Section{
				Name:       secdat.Name,
				StartedAt:  unixTime(secdat.Start),
				FinishedAt: unixTime(secdat.End),
				Duration:   float64(secdat.End - secdat.Start),
			}

			section.ID = int64(job.ID*1000 + secnum)
			section.Job.ID = int64(job.ID)
			section.Job.Name = job.Name
			section.Job.Status = job.Status
			section.Pipeline.ID = int64(job.Pipeline.ID)
			section.Pipeline.ProjectID = int64(job.Pipeline.ProjectID)
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

type sectionData struct {
	Name  string
	Start int64
	End   int64
}

func parseSections(trace *bytes.Reader) ([]sectionData, error) {
	sections := []sectionData{}
	stack := sectionStack{}

	scanner := bufio.NewScanner(trace)
	for scanner.Scan() {
		line := scanner.Bytes()
		if index := bytes.Index(line, []byte(sectionMarkerStart)); index >= 0 {
			ts, name, err := parseSection(sectionMarkerStart, line)
			if err != nil {
				// TODO: what?
			} else {
				stack.Start(ts, name)
			}
		} else if index := bytes.Index(line, []byte(sectionMarkerEnd)); index >= 0 {
			ts, name, err := parseSection(sectionMarkerEnd, line)
			if err != nil {
				// TODO: what?
			} else {
				sections = append(sections, stack.End(ts, name)...)
			}
		}
	}

	sort.SliceStable(sections, func(i, j int) bool {
		return sections[i].Start < sections[j].Start
	})

	return sections, nil
}

type sectionStack struct {
	Sections []sectionData
}

func (s *sectionStack) Len() int {
	return len(s.Sections)
}

func (s *sectionStack) Push(section sectionData) {
	s.Sections = append(s.Sections, section)
}

func (s *sectionStack) Pop() sectionData {
	size := len(s.Sections)
	if size == 0 {
		return sectionData{}
	}
	section := s.Sections[size-1]
	s.Sections = s.Sections[:(size - 1)]
	return section
}

func (s *sectionStack) Start(timestamp int64, name string) {
	s.Push(sectionData{
		Name:  name,
		Start: timestamp,
	})
}

func (s *sectionStack) End(timestamp int64, name string) []sectionData {
	endedSections := []sectionData{}

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

func parseSection(marker sectionMarker, line []byte) (timestamp int64, name string, err error) {
	pattern := regexp.MustCompile(fmt.Sprintf(`%s:(?P<ts>\d+):(?P<name>[\w_]+)`, marker))
	match := pattern.FindSubmatch(line)
	if len(match) != 3 {
		err = errors.New("no match found")
		return
	}

	time_s := string(match[pattern.SubexpIndex("ts")])
	name = string(match[pattern.SubexpIndex("name")])

	timestamp, err = strconv.ParseInt(time_s, 10, 64)

	return
}
