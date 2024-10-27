package types

import "github.com/cluttrdev/gitlab-exporter/protobuf/typespb"

type Metric struct {
	Id  string
	Iid int64
	Job JobReference

	Name      string
	Labels    map[string]string
	Value     float64
	Timestamp int64
}

func ConvertMetric(metric Metric) *typespb.Metric {
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
		Job: ConvertJobReference(metric.Job),

		Name:      metric.Name,
		Labels:    labels,
		Value:     metric.Value,
		Timestamp: ConvertUnixMilli(metric.Timestamp),
	}
}