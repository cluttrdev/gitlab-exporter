package types

import (
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.cluttr.dev/gitlab-exporter/protobuf/typespb"
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
		StartedAt:  timestamppb.New(valOrZero(section.StartedAt)),
		FinishedAt: timestamppb.New(valOrZero(section.FinishedAt)),
		Duration:   durationpb.New(section.Duration),
	}
}
