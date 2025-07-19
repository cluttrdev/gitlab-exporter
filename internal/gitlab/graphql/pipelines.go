package graphql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.cluttr.dev/gitlab-exporter/internal/types"
)

type PipelineFields struct {
	PipelineReferenceFields
	Project ProjectReferenceFields

	PipelineFieldsCore
	PipelineFieldsRelations
}

func ConvertPipeline(pf PipelineFields) (types.Pipeline, error) {
	var (
		id, iid, projectId int64
		err                error
	)
	if id, err = ParseId(pf.Id, GlobalIdPipelinePrefix); err != nil {
		return types.Pipeline{}, fmt.Errorf("parse pipeline id: %w", err)
	}
	if iid, err = ParseId(pf.Iid, ""); err != nil {
		return types.Pipeline{}, fmt.Errorf("parse pipeline iid: %w", err)
	}
	if projectId, err = ParseId(pf.Project.Id, GlobalIdProjectPrefix); err != nil {
		return types.Pipeline{}, fmt.Errorf("parse project id: %w", err)
	}

	p := types.Pipeline{
		Id:  id,
		Iid: iid,
		Project: types.ProjectReference{
			Id:       projectId,
			FullPath: pf.Project.FullPath,
		},

		Name:    valOrZero(pf.Name),
		Ref:     valOrZero(pf.Ref),
		RefPath: valOrZero(pf.RefPath),
		Sha:     valOrZero(pf.Sha),
		Source:  valOrZero(pf.Source),
		Status:  strings.ToLower(string(pf.Status)),

		CommittedAt: pf.CommittedAt,
		CreatedAt:   &pf.CreatedAt,
		UpdatedAt:   &pf.UpdatedAt,
		StartedAt:   pf.StartedAt,
		FinishedAt:  pf.FinishedAt,

		QueuedDuration: time.Duration(valOrZero(pf.QueuedDuration) * float64(time.Second)),
		Duration:       time.Duration(valOrZero(pf.Duration) * int(time.Second)),
		Coverage:       valOrZero(pf.Coverage),
		FailureReason:  valOrZero(pf.FailureReason),

		Warnings:   pf.Warnings,
		YamlErrors: pf.YamlErrors,

		Child: pf.Child,
		// UpstreamPipeline: nil,
		// DownstreamPipelines: nil,

		// MergeRequest: nil,

		// UserId: 0,
	}

	if pf.Upstream != nil {
		upstreamId, _ := ParseId(pf.Upstream.Id, GlobalIdPipelinePrefix)
		upstreamIid, _ := ParseId(pf.Upstream.Iid, "")
		upstreamProjectId, _ := ParseId(pf.Upstream.Project.Id, GlobalIdProjectPrefix)
		p.UpstreamPipeline = &types.PipelineReference{
			Id:  upstreamId,
			Iid: upstreamIid,
			Project: types.ProjectReference{
				Id:       upstreamProjectId,
				FullPath: pf.Upstream.Project.FullPath,
			},
		}
	}
	if pf.Downstream != nil {
		for _, dpf := range pf.Downstream.Nodes {
			downstreamId, _ := ParseId(dpf.Id, GlobalIdPipelinePrefix)
			downstreamIid, _ := ParseId(dpf.Iid, "")
			downstreamProjectId, _ := ParseId(dpf.Project.Id, GlobalIdProjectPrefix)
			p.DownstreamPipelines = append(p.DownstreamPipelines, &types.PipelineReference{
				Id:  downstreamId,
				Iid: downstreamIid,
				Project: types.ProjectReference{
					Id:       downstreamProjectId,
					FullPath: dpf.Project.FullPath,
				},
			})
		}
	}
	if pf.MergeRequest != nil {
		mrId, _ := ParseId(pf.MergeRequest.Id, GlobalIdMergeRequestPrefix)
		mrIid, _ := ParseId(pf.MergeRequest.Iid, "")
		mrProjectId, _ := ParseId(pf.MergeRequest.Project.Id, GlobalIdProjectPrefix)
		p.MergeRequest = &types.MergeRequestReference{
			Id:  mrId,
			Iid: mrIid,
			Project: types.ProjectReference{
				Id:       mrProjectId,
				FullPath: pf.Project.FullPath,
			},
		}
	}
	if pf.User != nil {
		user, err := convertUserReference(pf.User)
		if err != nil {
			return types.Pipeline{}, fmt.Errorf("convert user reference: %w", err)
		}
		p.User = user
	}

	return p, nil
}

type GetPipelinesOptions struct {
	UpdatedAfter  *time.Time
	UpdatedBefore *time.Time
}

func (c *Client) GetProjectsPipelines(ctx context.Context, ids []string, opts GetPipelinesOptions) ([]PipelineFields, error) {
	var pipelines []PipelineFields

	pipelines, err := c.getProjectsPipelines(ctx, ids, getPipelinesOptions{
		GetPipelinesOptions: opts,
	})
	if err != nil {
		return nil, err
	}

	childPipelines, err := c.getProjectsPipelines(ctx, ids, getPipelinesOptions{
		GetPipelinesOptions: opts,
		Source:              ptr("parent_pipeline"),
	})
	if err != nil {
		return nil, err
	}

	downstreamPipelines, err := c.getProjectsPipelines(ctx, ids, getPipelinesOptions{
		GetPipelinesOptions: opts,
		Source:              ptr("pipeline"),
	})
	if err != nil {
		return nil, err
	}

	pipelines = append(pipelines, append(childPipelines, downstreamPipelines...)...)

	return pipelines, nil
}

func (c *Client) GetProjectPipeline(ctx context.Context, projectId string, pipelineId string) (PipelineFields, error) {
	return c.getProjectIdPipeline(ctx, projectId, pipelineId)
}

type getPipelinesOptions struct {
	GetPipelinesOptions
	Source *string

	endCursor *string
}

func (c *Client) getProjectsPipelines(ctx context.Context, ids []string, opts getPipelinesOptions) ([]PipelineFields, error) {
	var pfs []PipelineFields

	for {
		resp, err := getProjectsPipelines(
			ctx,
			c.client,
			ids,
			opts.Source,
			opts.UpdatedAfter,
			opts.UpdatedBefore,
			opts.endCursor,
		)
		if err != nil {
			return nil, err
		}

		for _, project_ := range resp.Projects.Nodes {
			for _, pipeline_ := range project_.Pipelines.Nodes {
				pf := PipelineFields{
					PipelineReferenceFields: pipeline_.PipelineReferenceFields,
					Project: ProjectReferenceFields{
						Id:       project_.Id,
						FullPath: project_.FullPath,
					},

					PipelineFieldsCore:      pipeline_.PipelineFieldsCore,
					PipelineFieldsRelations: pipeline_.PipelineFieldsRelations,
				}

				if pipeline_.Downstream != nil && pipeline_.Downstream.PageInfo.HasNextPage {
					dpconn, err := c.getProjectPipelineDownstreamConnection(ctx, project_.FullPath, pipeline_.Iid, pipeline_.Downstream.PageInfo.EndCursor)
					if err != nil {
						return nil, err
					}
					pf.Downstream.Nodes = append(pf.Downstream.Nodes, dpconn.Nodes...)
					pf.Downstream.PageInfo = dpconn.PageInfo
				}

				pfs = append(pfs, pf)
			}

			if project_.Pipelines.PageInfo.HasNextPage {
				opts_ := opts
				opts_.endCursor = project_.Pipelines.PageInfo.EndCursor
				pfs_, err := c.getProjectPipelines(ctx, project_.FullPath, opts_)
				if err != nil {
					return nil, err
				}

				pfs = append(pfs, pfs_...)
			}
		}

		if !resp.Projects.PageInfo.HasNextPage {
			break
		}

		opts.endCursor = resp.Projects.PageInfo.EndCursor
	}

	return pfs, nil
}

func (c *Client) getProjectPipelines(ctx context.Context, path string, opts getPipelinesOptions) ([]PipelineFields, error) {
	var pfs []PipelineFields

	for {
		resp, err := getProjectPipelines(
			ctx,
			c.client,
			path,
			opts.Source,
			opts.UpdatedAfter,
			opts.UpdatedBefore,
			opts.endCursor,
		)
		if err != nil {
			return nil, err
		}

		project_ := resp.Project
		if project_ == nil {
			return nil, fmt.Errorf("project not found: %v", path)
		}

		for _, pipeline_ := range project_.Pipelines.Nodes {
			pf := PipelineFields{
				PipelineReferenceFields: pipeline_.PipelineReferenceFields,
				Project: ProjectReferenceFields{
					Id:       project_.Id,
					FullPath: project_.FullPath,
				},

				PipelineFieldsCore:      pipeline_.PipelineFieldsCore,
				PipelineFieldsRelations: pipeline_.PipelineFieldsRelations,
			}

			if pipeline_.Downstream != nil && pipeline_.Downstream.PageInfo.HasNextPage {
				dpconn, err := c.getProjectPipelineDownstreamConnection(ctx, project_.FullPath, pipeline_.Iid, pipeline_.Downstream.PageInfo.EndCursor)
				if err != nil {
					return nil, err
				}
				pf.Downstream.Nodes = append(pf.Downstream.Nodes, dpconn.Nodes...)
				pf.Downstream.PageInfo = dpconn.PageInfo
			}

			pfs = append(pfs, pf)
		}

		if !project_.Pipelines.PageInfo.HasNextPage {
			break
		}

		opts.endCursor = project_.Pipelines.PageInfo.EndCursor
	}

	return pfs, nil
}

func (c *Client) getProjectPipelineDownstreamConnection(ctx context.Context, projectPath string, pipelineIid string, endCursor *string) (*PipelineFieldsRelationsDownstreamPipelineConnection, error) {
	var dpconn PipelineFieldsRelationsDownstreamPipelineConnection

	for {
		resp, err := getProjectPipelineDownstream(ctx, c.client, projectPath, pipelineIid, endCursor)
		if err != nil {
			return nil, err
		}

		project_ := resp.Project
		if project_ == nil {
			return nil, fmt.Errorf("project not found: %v", projectPath)
		}
		pipeline_ := project_.Pipeline
		if pipeline_ == nil {
			return nil, fmt.Errorf("project pipeline not found: %v (%v)", pipelineIid, projectPath)
		}
		downstream_ := pipeline_.Downstream
		if downstream_ == nil {
			return nil, nil
		}

		for _, node_ := range pipeline_.Downstream.Nodes {
			dpconn.Nodes = append(dpconn.Nodes, &PipelineFieldsRelationsDownstreamPipelineConnectionNodesPipeline{
				PipelineReferenceFields: node_.PipelineReferenceFields,
				Project: &PipelineFieldsRelationsDownstreamPipelineConnectionNodesPipelineProject{
					ProjectReferenceFields: node_.Project.ProjectReferenceFields,
				},
			})
		}

		if !downstream_.PageInfo.HasNextPage {
			break
		}

		endCursor = downstream_.PageInfo.EndCursor
	}

	return &dpconn, nil
}

func (c *Client) getProjectIdPipeline(ctx context.Context, projectId string, pipelineId string) (PipelineFields, error) {
	resp, err := getProjectIdPipeline(ctx, c.client, projectId, pipelineId)
	if err != nil {
		return PipelineFields{}, err
	}

	if len(resp.Projects.Nodes) == 0 {
		return PipelineFields{}, fmt.Errorf("project not found: %v", projectId)
	}
	project_ := resp.Projects.Nodes[0]
	pipeline_ := project_.Pipeline
	if pipeline_ == nil {
		return PipelineFields{}, fmt.Errorf("project pipeline not found: %v/%v", projectId, pipelineId)
	}

	return PipelineFields{
		PipelineReferenceFields: pipeline_.PipelineReferenceFields,
		Project: ProjectReferenceFields{
			Id:       project_.Id,
			FullPath: project_.FullPath,
		},

		PipelineFieldsCore: pipeline_.PipelineFieldsCore,
	}, nil
}
