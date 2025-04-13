package messages

import (
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.cluttr.dev/gitlab-exporter/internal/types"
	"go.cluttr.dev/gitlab-exporter/protobuf/typespb"
)

func NewSection(section types.Section) *typespb.Section {
	return &typespb.Section{
		Id:  section.Id,
		Job: NewJobReference(section.Job),

		Name:       section.Name,
		StartedAt:  timestamppb.New(valOrZero(section.StartedAt)),
		FinishedAt: timestamppb.New(valOrZero(section.FinishedAt)),
		Duration:   durationpb.New(section.Duration),
	}
}
