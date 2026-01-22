package tasks

import (
	"testing"
	"time"
)

func Test_convertSectionTimestamps(t *testing.T) {
	expectedTime := time.Date(2006, time.January, 02, 15, 4, 5, 0, time.UTC)

	tests := []struct {
		testName string

		timestamp    int64
		expectedTime time.Time
	}{
		{"seconds", expectedTime.Unix(), expectedTime},
		{"milliseconds", expectedTime.UnixMilli(), expectedTime},
		{"microseconds", expectedTime.UnixMicro(), expectedTime},
		{"nanoseconds", expectedTime.UnixNano(), expectedTime},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got := convertSectionTimestamp(tt.timestamp)
			if !tt.expectedTime.Equal(got) {
				t.Errorf("want %v, got %v", tt.expectedTime.UTC(), got.UTC())
			}
		})
	}
}
