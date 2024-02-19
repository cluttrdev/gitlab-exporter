package models

import (
	pb "github.com/cluttrdev/gitlab-exporter/grpc/exporterpb"
)

func ConvertLabels(labels map[string]string) []*pb.Metric_Label {
	list := make([]*pb.Metric_Label, 0, len(labels))
	for name, value := range labels {
		list = append(list, &pb.Metric_Label{
			Name:  name,
			Value: value,
		})
	}
	return list
}
