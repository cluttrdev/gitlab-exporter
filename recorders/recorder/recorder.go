package recorder

import (
	"context"
	"time"

	"go.cluttr.dev/gitlab-exporter/protobuf/servicepb"
)

// Recorder defines the interface that all recorder implementations must satisfy.
// Recorders are responsible for storing GitLab data in various backends.
type Recorder interface {
	// Embed the gRPC server interface - all recorders must implement these methods
	servicepb.GitLabExporterServer

	// Name returns the unique name/identifier of this recorder type
	Name() string

	// Initialize prepares the recorder with the given configuration.
	// Called once before any recording operations.
	Initialize(ctx context.Context, config []byte) error

	// Start begins the recorder's operation (e.g., opens connections, starts background workers)
	// The implementation must be non-blocking.
	Start(ctx context.Context) error

	// Stop gracefully shuts down the recorder
	Stop(ctx context.Context) error

	// CheckHealth returns the current health status of the recorder
	CheckHealth(ctx context.Context) error
}

// Factory creates a new recorder instance.
// The factory receives raw configuration bytes that it can unmarshal as needed.
type Factory func() Recorder

// HealthStatus represents the health state of a recorder
type HealthStatus struct {
	Name      string
	Healthy   bool
	Message   string
	CheckedAt time.Time
}
