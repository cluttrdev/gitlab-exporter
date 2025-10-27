package rest

import (
	"bytes"
	"context"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func (c *Client) GetProjectJobArtifact(ctx context.Context, projectPath string, jobId int64, artifactPath string) (*bytes.Reader, error) {
	r, _, err := c.client.Jobs.DownloadSingleArtifactsFile(projectPath, int(jobId), artifactPath, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return r, nil
}
