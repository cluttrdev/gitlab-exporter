package types

import (
	"github.com/xanzy/go-gitlab"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

func ConvertJob(job *gitlab.Job) *typespb.Job {
	artifacts := make([]*typespb.JobArtifacts, 0, len(job.Artifacts))
	for _, a := range job.Artifacts {
		artifacts = append(artifacts, &typespb.JobArtifacts{
			Filename:   a.Filename,
			FileType:   a.FileType,
			FileFormat: a.FileFormat,
			Size:       int64(a.Size),
		})
	}

	return &typespb.Job{
		Id:   int64(job.ID),
		Name: job.Name,
		Pipeline: &typespb.PipelineReference{
			Id:        int64(job.Pipeline.ID),
			ProjectId: int64(job.Pipeline.ProjectID),
			Ref:       job.Pipeline.Ref,
			Sha:       job.Pipeline.Sha,
			Status:    job.Pipeline.Status,
		},
		Ref:            job.Ref,
		CreatedAt:      ConvertTime(job.CreatedAt),
		StartedAt:      ConvertTime(job.StartedAt),
		FinishedAt:     ConvertTime(job.FinishedAt),
		ErasedAt:       ConvertTime(job.ErasedAt),
		Duration:       ConvertDuration(job.Duration),
		QueuedDuration: ConvertDuration(job.QueuedDuration),
		Coverage:       job.Coverage,
		Stage:          job.Stage,
		Status:         job.Status,
		AllowFailure:   job.AllowFailure,
		FailureReason:  job.FailureReason,
		Tag:            job.Tag,
		WebUrl:         job.WebURL,
		TagList:        job.TagList,

		Commit:  convertCommit(job.Commit),
		Project: ConvertProject(job.Project),
		User:    convertUser(job.User),

		Runner: &typespb.JobRunner{
			Id:          int64(job.Runner.ID),
			Name:        job.Runner.Name,
			Description: job.Runner.Description,
			Active:      job.Runner.Active,
			IsShared:    job.Runner.IsShared,
		},

		Artifacts: artifacts,
		ArtifactsFile: &typespb.JobArtifactsFile{
			Filename: job.ArtifactsFile.Filename,
			Size:     int64(job.ArtifactsFile.Size),
		},
		ArtifactsExpireAt: ConvertTime(job.ArtifactsExpireAt),
	}
}

func convertCommit(commit *gitlab.Commit) *typespb.Commit {
	var status string
	if commit.Status != nil {
		status = string(*commit.Status)
	}
	return &typespb.Commit{
		Id:             commit.ID,
		ShortId:        commit.ShortID,
		ParentIds:      commit.ParentIDs,
		ProjectId:      int64(commit.ProjectID),
		AuthorName:     commit.AuthorName,
		AuthorEmail:    commit.AuthorEmail,
		AuthoredDate:   ConvertTime(commit.AuthoredDate),
		CommitterName:  commit.CommitterName,
		CommitterEmail: commit.CommitterEmail,
		CommittedDate:  ConvertTime(commit.CommittedDate),
		CreatedAt:      ConvertTime(commit.CreatedAt),
		Title:          commit.Title,
		Message:        commit.Message,
		Trailers:       commit.Trailers,
		Stats:          convertCommitStats(commit.Stats),
		Status:         status,
		WebUrl:         commit.WebURL,
	}
}

func convertCommitStats(stats *gitlab.CommitStats) *typespb.CommitStats {
	if stats == nil {
		return nil
	}
	return &typespb.CommitStats{
		Additions: int64(stats.Additions),
		Deletions: int64(stats.Deletions),
		Total:     int64(stats.Total),
	}
}

func ConvertBridge(bridge *gitlab.Bridge) *typespb.Bridge {
	// account for downstream pipeline creation failures
	downstreamPipeline := &typespb.PipelineInfo{
		CreatedAt: &timestamppb.Timestamp{},
		UpdatedAt: &timestamppb.Timestamp{},
	}
	if bridge.DownstreamPipeline != nil {
		downstreamPipeline = ConvertPipelineInfo(bridge.DownstreamPipeline)
	}
	return &typespb.Bridge{
		// Commit: ConvertCommit(bridge.Commit),
		Id:             int64(bridge.ID),
		Name:           bridge.Name,
		Pipeline:       ConvertPipelineInfo(&bridge.Pipeline),
		Ref:            bridge.Ref,
		CreatedAt:      ConvertTime(bridge.CreatedAt),
		StartedAt:      ConvertTime(bridge.StartedAt),
		FinishedAt:     ConvertTime(bridge.FinishedAt),
		ErasedAt:       ConvertTime(bridge.ErasedAt),
		Duration:       ConvertDuration(bridge.Duration),
		QueuedDuration: ConvertDuration(bridge.QueuedDuration),
		Coverage:       bridge.Coverage,
		Stage:          bridge.Stage,
		Status:         bridge.Status,
		AllowFailure:   bridge.AllowFailure,
		FailureReason:  bridge.FailureReason,
		Tag:            bridge.Tag,
		WebUrl:         bridge.WebURL,
		// User: ConvertUser(bridge.User),
		DownstreamPipeline: downstreamPipeline,
	}
}
