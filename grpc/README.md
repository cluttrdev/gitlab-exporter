# GitLab Exporter gRPC Infrastructure

This module provides shared gRPC client and server infrastructure used by both
the exporter and recorder implementations.

It contains:
- **Server**: gRPC server with health checks, metrics, and Unix socket support
- **Client**: gRPC client with connection management and helper methods

## Module Structure

```
grpc/
├── server/
│   └── server.go      # gRPC server implementation
└── client/
    └── client.go      # gRPC client implementation
```

## Server Usage

The server supports both TCP and Unix domain sockets:

```go
import (
    "go.cluttr.dev/gitlab-exporter/grpc/server"
    "go.cluttr.dev/gitlab-exporter/protobuf/servicepb"
)

// Create recorder implementation
type MyRecorder struct {
    servicepb.UnimplementedGitLabExporterServer
    // ... your implementation
}

recorder := &MyRecorder{}

// Create and start server
srv := server.New(recorder)

// Listen on TCP
srv.ListenAndServe(ctx, "0.0.0.0:9100")

// Or listen on Unix socket
srv.ListenAndServe(ctx, "unix:///tmp/recorder.sock")

// Or with bare path
srv.ListenAndServe(ctx, "/tmp/recorder.sock")
```

### Unix Socket Support

Unix domain sockets provide:
- **Better performance**: Lower latency than TCP for local communication
- **No port conflicts**: No need to manage port allocations
- **Simpler security**: File system permissions control access

The server automatically detects socket addresses:

| Format | Network | Example |
|--------|---------|---------|
| `unix:///path` | Unix | `unix:///tmp/recorder.sock` |
| `/path` | Unix | `/tmp/recorder.sock` |
| `host:port` | TCP | `localhost:9100` |

## Client Usage

```go
import (
    "go.cluttr.dev/gitlab-exporter/grpc/client"
    "google.golang.org/grpc/credentials/insecure"
)

// Connect to recorder
c, err := client.NewClient(
    "unix:///tmp/recorder.sock",
    grpc.WithTransportCredentials(insecure.NewCredentials()),
)
if err != nil {
    log.Fatal(err)
}

// Send data
pipelines := []*typespb.Pipeline{...}
err = client.RecordPipelines(c, ctx, pipelines)
```
