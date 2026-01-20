package exporter

import (
	"testing"
	"time"

	tracepb_v1 "go.opentelemetry.io/proto/otlp/trace/v1"

	"go.cluttr.dev/gitlab-exporter/exporter/internal/types"
)

func Test_convert_and_filter(t *testing.T) {
	var cfun convertFunc[types.Job, *tracepb_v1.Span] = func(_ types.Job) *tracepb_v1.Span {
		return nil
	}

	data := []types.Job{
		{Id: 1, StartedAt: &[]time.Time{time.Unix(1, 0)}[0], FinishedAt: &[]time.Time{time.Unix(2, 0)}[0]},
		{Id: 2, StartedAt: &[]time.Time{time.Unix(3, 0)}[0], FinishedAt: nil},
		{Id: 3, StartedAt: nil, FinishedAt: nil},
	}

	results := convert(data, cfun)
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	results = filterNil(results)
	if len(results) != 0 {
		t.Errorf("Expected 0 results after filtering, got %d", len(results))
	}
}
