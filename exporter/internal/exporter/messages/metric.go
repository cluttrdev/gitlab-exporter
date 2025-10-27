package messages

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"go.cluttr.dev/gitlab-exporter/internal/types"
	"go.cluttr.dev/gitlab-exporter/protobuf/typespb"
)

func NewMetric(metric types.Metric) *typespb.Metric {
	var labels []*typespb.Metric_Label
	for name, value := range metric.Labels {
		labels = append(labels, &typespb.Metric_Label{
			Name:  name,
			Value: value,
		})
	}

	return &typespb.Metric{
		Id:  []byte(metric.Id),
		Iid: metric.Iid,
		Job: NewJobReference(metric.Job),

		Name:      metric.Name,
		Labels:    labels,
		Value:     metric.Value,
		Timestamp: timestamppb.New(time.UnixMilli(metric.Timestamp)),
	}
}
