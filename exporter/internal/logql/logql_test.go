package logql_test

import (
	"bytes"
	"testing"

	"go.cluttr.dev/gitlab-exporter/internal/logql"
)

func TestCount(t *testing.T) {
	tests := []struct {
		testName string // description of this test case
		// Named input parameters for target function.
		log     []byte
		filters []logql.LineFilter
		want    []int
		wantErr bool
	}{
		{
			testName: "logfmt",
			log: []byte(`
                level=info ts=2022-03-23T11:55:29.846163306Z caller=main.go:112 msg="Starting Grafana Enterprise Logs"
                level=debug ts=2022-03-23T11:55:29.846226372Z caller=main.go:113 version=v1.3.0 branch=HEAD Revision=e071a811 LokiVersion=v2.4.2 LokiRevision=525040a3
                level=warn ts=2022-03-23T11:55:45.213901602Z caller=added_modules.go:198 msg="found valid license" cluster=enterprise-logs-test-fixture
                level=info ts=2022-03-23T11:55:45.214611239Z caller=server.go:269 http=[::]:3100 grpc=[::]:9095 msg="server listening on addresses"
                level=debug ts=2022-03-23T11:55:45.219665469Z caller=module_service.go:64 msg=initialising module=license
                level=warm ts=2022-03-23T11:55:45.219678992Z caller=module_service.go:64 msg=initialising module=server
                level=error ts=2022-03-23T11:55:45.221140583Z caller=manager.go:132 msg="license manager up and running"
                level=info ts=2022-03-23T11:55:45.221254326Z caller=loki.go:355 msg="Loki started"
            `),
			filters: []logql.LineFilter{
				{
					{Operator: `|=`, Patterns: []string{"level=info"}},
				},
				{
					{Operator: `|=`, Patterns: []string{"level=warn", "level=warm"}},
				},
				{
					{Operator: `|~`, Patterns: []string{`level=war(n|m)`}},
				},
			},
			want:    []int{3, 2, 2},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got, gotErr := logql.Count(bytes.NewReader(tt.log), tt.filters)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Count() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Count() succeeded unexpectedly")
			}

			if len(got) != len(tt.want) {
				t.Fatalf("Count() returned %d counts, want %d", len(got), len(tt.want))
			}

			for i, count := range got {
				if count != tt.want[i] {
					t.Errorf("Count()[%d] = %v, want %v", i, count, tt.want[i])
				}
			}

		})
	}
}
