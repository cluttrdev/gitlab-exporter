package types

import (
	"strconv"

	"github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

func ConvertPipelineInfo(pipeline *gitlab.PipelineInfo) *typespb.PipelineInfo {
	if pipeline == nil {
		return nil
	}

	return &typespb.PipelineInfo{
		Id:        int64(pipeline.ID),
		Iid:       int64(pipeline.IID),
		ProjectId: int64(pipeline.ProjectID),
		Status:    pipeline.Status,
		Source:    pipeline.Source,
		Ref:       pipeline.Ref,
		Sha:       pipeline.SHA,
		WebUrl:    pipeline.WebURL,
		CreatedAt: ConvertTime(pipeline.CreatedAt),
		UpdatedAt: ConvertTime(pipeline.UpdatedAt),
	}
}

func ConvertPipeline(pipeline *gitlab.Pipeline) *typespb.Pipeline {
	return &typespb.Pipeline{
		Id:             int64(pipeline.ID),
		Iid:            int64(pipeline.IID),
		ProjectId:      int64(pipeline.ProjectID),
		Status:         pipeline.Status,
		Source:         pipeline.Source,
		Ref:            pipeline.Ref,
		Sha:            pipeline.SHA,
		BeforeSha:      pipeline.BeforeSHA,
		Tag:            pipeline.Tag,
		YamlErrors:     pipeline.YamlErrors,
		CreatedAt:      ConvertTime(pipeline.CreatedAt),
		UpdatedAt:      ConvertTime(pipeline.UpdatedAt),
		StartedAt:      ConvertTime(pipeline.StartedAt),
		FinishedAt:     ConvertTime(pipeline.FinishedAt),
		CommittedAt:    ConvertTime(pipeline.CommittedAt),
		Duration:       ConvertDuration(float64(pipeline.Duration)),
		QueuedDuration: ConvertDuration(float64(pipeline.QueuedDuration)),
		Coverage:       convertCoverage(pipeline.Coverage),
		WebUrl:         pipeline.WebURL,
		User:           convertBasicUser(pipeline.User),
	}
}

func convertCoverage(coverage string) float64 {
	cov, err := strconv.ParseFloat(coverage, 64)
	if err != nil {
		cov = 0.0
	}
	return cov
}
