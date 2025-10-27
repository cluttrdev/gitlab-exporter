# Recorder Framework

This package provides the core framework for implementing built-in recorders
in the GitLab Exporter.

It enables:
- **Pluggable backends**: Multiple storage implementations (ClickHouse, SQLite, etc.)
- **Clean separation**: Recorders implement a standard interface
- **Lifecycle management**: Coordinated start/stop and health checking
- **Easy configuration**: YAML-based configuration with type-specific options

## Architecture

```
┌─────────────────────┐
│  Recorder Interface │  <- All recorders implement this
└──────────┬──────────┘
           │
           ├─── Registry (type -> factory mapping)
           │
           └─── Manager (lifecycle orchestration)
```

## Core Interface

```go
type Recorder interface {
    // Embed gRPC server interface for all Record* methods
    servicepb.GitLabExporterServer

    // Lifecycle methods
    Name() string
    Initialize(ctx context.Context, config []byte) error
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    CheckHealth(ctx context.Context) error
}
```

## Implementing a Recorder

### 1. Define Your Recorder Struct

```go
package myrecorder

import (
    "context"
    "go.cluttr.dev/gitlab-exporter/protobuf/servicepb"
    "go.cluttr.dev/gitlab-exporter/recorders/recorder"
)

type MyRecorder struct {
    servicepb.UnimplementedGitLabExporterServer
    // ... your fields
}
```

### 2. Implement Required Methods

```go
func (r *MyRecorder) Name() string {
    return "myrecorder"
}

func (r *MyRecorder) Initialize(ctx context.Context, config []byte) error {
    // Parse config
    var cfg recorder.Config
    if err := yaml.Unmarshal(config, &cfg); err != nil {
        return err
    }

    // Extract recorder-specific options
    // cfg.Settings is map[string]any - handle as needed

    return nil
}

func (r *MyRecorder) Start(ctx context.Context) error {
    // Open connections, start workers, etc.
    return nil
}

func (r *MyRecorder) Stop(ctx context.Context) error {
    // Clean shutdown
    return nil
}

func (r *MyRecorder) Health(ctx context.Context) error {
    // Check if recorder is healthy
    return nil
}
```

### 3. Implement Record Methods

Implement whichever `Record*` methods your recorder needs:

```go
func (r *MyRecorder) RecordPipelines(ctx context.Context, req *servicepb.RecordPipelinesRequest) (*servicepb.RecordSummary, error) {
    // Store pipelines in your backend
    count := len(req.Data)
    return &servicepb.RecordSummary{RecordedCount: int32(count)}, nil
}
```

Methods not implemented will use the embedded
`UnimplementedGitLabExporterServer` which returns "not implemented" errors.

### 4. Register Your Recorder

```go
func init() {
    recorder.Register("myrecorder", func() recorder.Recorder {
        return &MyRecorder{}
    })
}
```

## Configuration

Recorders are configured via YAML:

```yaml
recorders:
  - type: myrecorder
    enabled: true
    address: localhost:9000  # optional, use as needed
    settings:
      # Recorder-specific settings
      # Your recorder should parse these in Initialize()
      database: mydb
      custom_option: value
```

### Config Structure

```go
type Config struct {
    Type     string         `yaml:"type"`      // Recorder type name
    Enabled  bool           `yaml:"enabled"`   // Enable/disable
    Address  string         `yaml:"address"`   // Listen address
    Settings map[string]any `yaml:"settings"`  // Recorder-specific config
}
```

The `Settings` field is intentionally generic - each recorder can parse it
according to its needs.

## Usage

### Loading Recorders

```go
import "go.cluttr.dev/gitlab-exporter/recorders/recorder"

// Create registry
registry := recorder.NewRegistry()

registry.Register("clickhouse", ... )
registry.Register("sqlite", ... )

// Load from config
configs := []recorder.Config{
    {Type: "clickhouse", Enabled: true, Address: "localhost:9000"},
    {Type: "sqlite", Enabled: true, Address: "unix:///tmp/gitlab-exporter-sqlite.sock"},
}

recorders, err := recorder.LoadFromConfigs(ctx, registry, configs)
if err != nil {
    log.Fatal(err)
}

// Create manager
manager := recorder.NewManager()

for _, rec := range recorders {
    manager.Add(rec)
}

// Start all recorders
if err := manager.StartAll(ctx); err != nil {
    log.Fatal(err)
}
defer manager.StopAll(context.Background())
```

### Health Checks

```go
statuses := manager.HealthCheckAll(ctx)
for _, status := range statuses {
    if !status.Healthy {
        log.Printf("Recorder %s unhealthy: %s", status.Name, status.Message)
    }
}
```

## Registry

The registry maps recorder type names to factory functions:

```go
// Create new registry
registry := recorder.NewRegistry()

// Register a recorder type
registry.Register("mytype", func() recorder.Recorder {
    return &MyRecorder{}
})

// List registered types
types := registry.List()

// Create instance
rec, err := registry.Create("mytype")
```
