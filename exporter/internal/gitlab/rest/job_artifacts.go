package rest

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func (c *Client) DownloadProjectJobArtifactsFile(ctx context.Context, projectPath string, jobId int64, artifactPath string) (*bytes.Reader, error) {
	r, _, err := c.client.Jobs.DownloadSingleArtifactsFile(projectPath, int(jobId), artifactPath, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (c *Client) GetProjectJobArtifacts(ctx context.Context, pid any, jobId int, fileType string) (*bytes.Reader, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/jobs/%d/artifacts", gitlab.PathEscape(project), jobId)

	opt := struct {
		FileType string `url:"file_type"`
	}{
		FileType: fileType,
	}

	options := []gitlab.RequestOptionFunc{
		gitlab.WithContext(ctx),
	}

	req, err := c.client.NewRequest(http.MethodGet, u, &opt, options)
	if err != nil {
		return nil, err
	}

	artifactsBuf := new(bytes.Buffer)
	_, err = c.client.Do(req, artifactsBuf)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(artifactsBuf.Bytes()), nil
}
