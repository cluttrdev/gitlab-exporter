.PHONY: protobuf clean

protobuf:  ## Generate Protocol Buffer and gRPC code
	protoc \
		-I protos/ \
		-I protos/vendor/opentelemetry-proto \
		--go_out=. --go_opt=module=go.cluttr.dev/gitlab-exporter/protobuf \
		--go-grpc_out=. --go-grpc_opt=module=go.cluttr.dev/gitlab-exporter/protobuf \
		protos/gitlabexporter/protobuf/*.proto protos/gitlabexporter/protobuf/service/*.proto

clean:  ## Remove generated files
	rm -f servicepb/*.pb.go
	rm -f typespb/*.pb.go

.PHONY: help
help:  ## Display this help page
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[33m%-30s\033[0m %s\n", $$1, $$2}'

ifneq "${VERBOSE}" "1"
.SILENT:
endif
