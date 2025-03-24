package rest

import (
	"bufio"
	"bytes"
	"context"
	"fmt"

	_gitlab "gitlab.com/gitlab-org/api/client-go"

	"go.cluttr.dev/gitlab-exporter/internal/expfmt"
)

type JobLogData struct {
	Sections []SectionData `json:"sections"`
	Metrics  []MetricData  `json:"metrics"`
}

type MetricData struct {
	Name      string            `json:"name"`
	Labels    map[string]string `json:"labels"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp"`
}

func (c *Client) GetJobLogData(ctx context.Context, projectId int64, jobId int64) (JobLogData, error) {
	log, _, err := c.client.Jobs.GetTraceFile(int(projectId), int(jobId), _gitlab.WithContext(ctx))
	if err != nil {
		return JobLogData{}, fmt.Errorf("get job log: %w", err)
	}

	data, err := ParseJobLog(log)
	if err != nil {
		return JobLogData{}, fmt.Errorf("parse job log: %w", err)
	}

	return data, nil
}

func ParseJobLog(trace *bytes.Reader) (JobLogData, error) {
	var (
		data     JobLogData
		sections sectionStack
		parser   expfmt.TextParser
	)

	const (
		METRIC_MARKER  = `METRIC_`
		SECTION_MARKER = `section_`
	)

	scanner := bufio.NewScanner(trace)
	for scanner.Scan() {
		line := scanner.Bytes()
		if bytes.HasPrefix(line, []byte(METRIC_MARKER)) {
			metric, err := parser.LineToMetric(line[7:])
			if err != nil {
				// TODO: what?
				continue
			}
			data.Metrics = append(data.Metrics, MetricData{
				Name:      metric.Name,
				Labels:    convertMetricLabels(metric.Labels),
				Value:     metric.Value,
				Timestamp: metric.TimestampMs,
			})
		}

		var i, j int
		sep := []byte(SECTION_MARKER)
		for {
			j = bytes.Index(line[i:], sep)
			if j < 0 {
				break
			}

			marker, ts, name, err := parseSection(line[i:])
			if err != nil {
				// TODO: what?
			} else if marker == string(sectionMarkerStart) {
				sections.Start(ts, name)
			} else if marker == string(sectionMarkerEnd) {
				data.Sections = append(data.Sections, sections.End(ts, name)...)
			}

			i = i + j + 1
		}
	}

	// add all unfinished sections (e.g. due to job interruption) open-ended
	// to give caller a chance to set end timestamp
	for sections.Size() > 0 {
		data.Sections = append(data.Sections, sections.Pop())
	}

	return data, nil
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
