package models

import (
	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

func ConvertLabels(labels map[string]string) []*typespb.Metric_Label {
	list := make([]*typespb.Metric_Label, 0, len(labels))
	for name, value := range labels {
		list = append(list, &typespb.Metric_Label{
			Name:  name,
			Value: value,
		})
	}
	return list
}
