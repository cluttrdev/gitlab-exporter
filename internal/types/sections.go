package types

import (
	"time"

	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

type Section struct {
	Id  int64
	Job JobReference

	Name       string
	StartedAt  *time.Time
	FinishedAt *time.Time
	Duration   time.Duration
}

func ConvertSection(section Section) *typespb.Section {
	return &typespb.Section{
		Id:  section.Id,
		Job: ConvertJobReference(section.Job),

		Name:       section.Name,
		StartedAt:  ConvertTime(section.StartedAt),
		FinishedAt: ConvertTime(section.FinishedAt),
		Duration:   ConvertDuration(section.Duration.Seconds()),
	}
}
