package exporter

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/cluttrdev/gitlab-exporter/internal/config"
)

func CreateEndpointConfigs(cfg []config.Endpoint) []EndpointConfig {
	endpoints := make([]EndpointConfig, 0, len(cfg))
	for _, cc := range cfg {
		endpoints = append(endpoints, EndpointConfig{
			Address: cc.Address,
			Options: []grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			},
		})
	}
	return endpoints
}
