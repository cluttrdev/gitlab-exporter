module go.cluttr.dev/gitlab-exporter/recorders/recorder

go 1.24.3

replace (
	go.cluttr.dev/gitlab-exporter/grpc v0.0.0 => ../../grpc
	go.cluttr.dev/gitlab-exporter/protobuf v0.0.0 => ../../protobuf
)

require (
	go.cluttr.dev/gitlab-exporter/grpc v0.0.0
	go.cluttr.dev/gitlab-exporter/protobuf v0.0.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/kr/pretty v0.3.1 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	go.opentelemetry.io/proto/otlp v1.8.0 // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8 // indirect
	google.golang.org/grpc v1.76.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)
