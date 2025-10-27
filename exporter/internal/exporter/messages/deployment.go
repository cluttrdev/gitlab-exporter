package messages

import (
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"

	"go.cluttr.dev/gitlab-exporter/internal/types"
	"go.cluttr.dev/gitlab-exporter/protobuf/typespb"
)

func NewEnvironmentReference(env types.EnvironmentReference) *typespb.EnvironmentReference {
	var tier typespb.DeploymentTier = typespb.DeploymentTier_DEPLOYMENT_TIER_UNSPECIFIED
	switch strings.ToLower(env.Tier) {
	case "production":
		tier = typespb.DeploymentTier_DEPLOYMENT_TIER_PRODUCTION
	case "staging":
		tier = typespb.DeploymentTier_DEPLOYMENT_TIER_STAGING
	case "testing":
		tier = typespb.DeploymentTier_DEPLOYMENT_TIER_TESTING
	case "development":
		tier = typespb.DeploymentTier_DEPLOYMENT_TIER_DEVELOPMENT
	case "other":
		tier = typespb.DeploymentTier_DEPLOYMENT_TIER_OTHER
	}

	return &typespb.EnvironmentReference{
		Id:   env.Id,
		Name: env.Name,
		Tier: tier,

		Project: NewProjectReference(env.Project),
	}
}

func NewDeployment(dep types.Deployment) *typespb.Deployment {
	var status typespb.DeploymentStatus = typespb.DeploymentStatus_DEPLOYMENT_STATUS_UNSPECIFIED
	switch strings.ToLower(dep.Status) {
	case "created":
		status = typespb.DeploymentStatus_DEPLOYMENT_STATUS_CREATED
	case "running":
		status = typespb.DeploymentStatus_DEPLOYMENT_STATUS_RUNNING
	case "success":
		status = typespb.DeploymentStatus_DEPLOYMENT_STATUS_SUCCESS
	case "failed":
		status = typespb.DeploymentStatus_DEPLOYMENT_STATUS_FAILED
	case "canceled":
		status = typespb.DeploymentStatus_DEPLOYMENT_STATUS_CANCELED
	case "skipped":
		status = typespb.DeploymentStatus_DEPLOYMENT_STATUS_SKIPPED
	case "blocked":
		status = typespb.DeploymentStatus_DEPLOYMENT_STATUS_BLOCKED
	}

	return &typespb.Deployment{
		Id:  dep.Id,
		Iid: dep.Iid,

		Job:         NewJobReference(dep.Job),
		Triggerer:   NewUserReference(dep.Triggerer),
		Environment: NewEnvironmentReference(dep.Environment),

		Timestamps: &typespb.DeploymentTimestamps{
			CreatedAt:  timestamppb.New(valOrZero(dep.CreatedAt)),
			FinishedAt: timestamppb.New(valOrZero(dep.FinishedAt)),
			UpdatedAt:  timestamppb.New(valOrZero(dep.UpdatedAt)),
		},

		Status: status,
		Ref:    dep.Ref,
		Sha:    dep.Sha,
	}
}
