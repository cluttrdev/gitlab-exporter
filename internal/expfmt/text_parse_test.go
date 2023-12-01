package expfmt_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/cluttrdev/gitlab-exporter/internal/expfmt"
)

func TestTextParser_LineToMetric(t *testing.T) {
	var scenarios = []struct {
		in  string
		out *expfmt.Metric
	}{
		{
			in:  "",
			out: nil,
		},
		{
			in: "minimal_metric 42.1337",
			out: &expfmt.Metric{
				Name:  "minimal_metric",
				Value: 42.1337,
			},
		},
		{
			in: "no_labels{} 3",
			out: &expfmt.Metric{
				Name:  "no_labels",
				Value: 3,
			},
		},
		{
			in: `metric_name{label_name="label_value",} 1 1234567890`,
			out: &expfmt.Metric{
				Name:        "metric_name",
				Labels:      []expfmt.MetricLabelPair{{Name: "label_name", Value: "label_value"}},
				Value:       1,
				TimestampMs: 1234567890,
			},
		},
		{
			in: `metric_name{label_name1="label_value1",label_name2="label_value2"} 1 1234567890`,
			out: &expfmt.Metric{
				Name: "metric_name",
				Labels: []expfmt.MetricLabelPair{
					{Name: "label_name1", Value: "label_value1"},
					{Name: "label_name2", Value: "label_value2"},
				},
				Value:       1,
				TimestampMs: 1234567890,
			},
		},
	}

	var parser expfmt.TextParser

	for i, scenario := range scenarios {
		out, err := parser.LineToMetric([]byte(scenario.in))
		if err != nil {
			t.Errorf("%d. error: %s", i, err)
			continue
		}

		if diff := cmp.Diff(scenario.out, out); diff != "" {
			t.Errorf("Config mismatch (-want +got):\n%s", diff)
		}
	}
}
