package exporter

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpc_client "github.com/cluttrdev/gitlab-exporter/grpc/client"

	"github.com/cluttrdev/gitlab-exporter/internal/config"
)

func CreateEndpointConfigs(cfg []config.Endpoint) []grpc_client.EndpointConfig {
	endpoints := make([]grpc_client.EndpointConfig, 0, len(cfg))
	for _, cc := range cfg {
		endpoints = append(endpoints, grpc_client.EndpointConfig{
			Address: cc.Address,
			Options: []grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			},
		})
	}
	return endpoints
}
