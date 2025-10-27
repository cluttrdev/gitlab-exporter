# GitLab Exporter Protocol Buffers

This module contains the Protocol Buffer definitions that define the API
contract between the exporter and recorder implementations.

## Overview

The protobuf module defines:
- gRPC service interfaces for recording GitLab data
- Message types for all exported entities (pipelines, jobs, merge requests, etc.)
- Common types and references

## Module Structure

```
protobuf/
├── protos/                              # Proto definition files
│   └── gitlabexporter/protobuf/
│       ├── service/service.proto        # gRPC service definition
│       ├── pipeline.proto               # Pipeline messages
│       ├── job.proto                    # Job messages
│       └── ...                          # Other entity types
├── servicepb/                           # Generated gRPC service code
│   ├── service.pb.go
│   └── service_grpc.pb.go
└── typespb/                             # Generated message types
    ├── pipeline.pb.go
    ├── job.pb.go
    └── ...
```

## Usage

### For Exporter (Client)

```go
import (
    "go.cluttr.dev/gitlab-exporter/protobuf/servicepb"
    "go.cluttr.dev/gitlab-exporter/protobuf/typespb"
)

// Create gRPC client and send data
client := servicepb.NewGitLabExporterClient(conn)
pipelines := []*typespb.Pipeline{...}
client.RecordPipelines(ctx, &servicepb.RecordPipelinesRequest{Data: pipelines})
```

### For Recorders (Server)

```go
import (
    "go.cluttr.dev/gitlab-exporter/protobuf/servicepb"
    "go.cluttr.dev/gitlab-exporter/protobuf/typespb"
)

// Implement the server interface
type MyRecorder struct {
    servicepb.UnimplementedGitLabExporterServer
}

func (r *MyRecorder) RecordPipelines(ctx context.Context, req *servicepb.RecordPipelinesRequest) (*servicepb.RecordSummary, error) {
    // Store req.Data in your backend
    return &servicepb.RecordSummary{RecordedCount: int32(len(req.Data))}, nil
}
```

## Regenerating Code

After modifying `.proto` files:

```bash
make protobuf
```

This requires `protoc` with the `protoc-gen-go` and `protoc-gen-go-grpc` plugins installed.
