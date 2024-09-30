package gitlab

import (
	"bufio"
	"bytes"
	"context"

	_gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-exporter/internal/expfmt"
)

func (c *Client) GetJobLog(ctx context.Context, projectID int64, jobID int64) (*bytes.Reader, error) {
	trace, _, err := c.client.Jobs.GetTraceFile(int(projectID), int(jobID), _gitlab.WithContext(ctx))
	return trace, err
}

type JobLogData struct {
	Sections []SectionData `json:"sections"`
	Metrics  []*MetricData `json:"metrics"`
}

type MetricData struct {
	Name      string            `json:"name"`
	Labels    map[string]string `json:"labels"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp"`
}

func ParseJobLog(trace *bytes.Reader) (*JobLogData, error) {
	var (
		data   JobLogData
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
			data.Metrics = append(data.Metrics, &MetricData{
				Name:      metric.Name,
				Labels:    convertMetricLabels(metric.Labels),
				Value:     metric.Value,
				Timestamp: metric.TimestampMs,
			})
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

func convertMetricLabels(pairs []expfmt.MetricLabelPair) map[string]string {
	if len(pairs) == 0 {
		return nil
	}
	labels := map[string]string{}
	for _, p := range pairs {
		labels[p.Name] = p.Value
	}
	return labels
}
