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
		} else if index := bytes.Index(line, []byte(sectionMarkerStart)); index >= 0 {
			ts, name, err := parseSection(sectionMarkerStart, line)
			if err != nil {
				// TODO: what?
			} else {
				stack.Start(ts, name)
			}
		} else if index := bytes.Index(line, []byte(sectionMarkerEnd)); index >= 0 {
			ts, name, err := parseSection(sectionMarkerEnd, line)
			if err != nil {
				// TODO: what?
			} else {
				data.Sections = append(data.Sections, stack.End(ts, name)...)
			}
		}
	}

	return &data, nil
}
