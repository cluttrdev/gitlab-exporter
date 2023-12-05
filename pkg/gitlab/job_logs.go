package gitlab

import (
	"bufio"
	"bytes"
	"context"

	_gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-exporter/internal/expfmt"
)

func (c *Client) GetJobLog(ctx context.Context, projectID int64, jobID int64) (*bytes.Reader, error) {
	c.RLock()
	trace, _, err := c.client.Jobs.GetTraceFile(int(projectID), int(jobID), _gitlab.WithContext(ctx))
	c.RUnlock()
	return trace, err
}

type jobLogData struct {
	Sections []sectionData
	Metrics  []*expfmt.Metric
}

func parseJobLog(trace *bytes.Reader) (*jobLogData, error) {
	var (
		data   jobLogData
		stack  sectionStack
		parser expfmt.TextParser
	)

	scanner := bufio.NewScanner(trace)
	for scanner.Scan() {
		line := scanner.Bytes()
		if bytes.HasPrefix(line, []byte(`METRIC_`)) {
			metric, err := parser.LineToMetric(line[7:])
			if err != nil {
				// TODO: what?
				continue
			}
			data.Metrics = append(data.Metrics, metric)
		}

		var i, j int
		sep := []byte(`section_`)
		for {
			j = bytes.Index(line[i:], sep)
			if j < 0 {
				break
			}

			marker, ts, name, err := parseSection(line[i:])
			if err != nil {
				// TODO: what?
			} else if marker == string(sectionMarkerStart) {
				stack.Start(ts, name)
			} else if marker == string(sectionMarkerEnd) {
				data.Sections = append(data.Sections, stack.End(ts, name)...)
			}

			i = i + j + 1
		}
	}

	return &data, nil
}
