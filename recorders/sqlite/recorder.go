package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"

	"go.cluttr.dev/gitlab-exporter/protobuf/servicepb"
)

// Recorder implements the recorder.Recorder interface for SQLite storage
type Recorder struct {
	servicepb.UnimplementedGitLabExporterServer

	db       *sql.DB
	address  string
	settings Settings
}

// Settings holds SQLite-specific configuration
type Settings struct {
	// Database file path
	Path string `yaml:"path"`

	// Enable WAL mode for better concurrency
	WALMode bool `yaml:"wal_mode"`

	// Batch size for inserts
	BatchSize int `yaml:"batch_size"`
}

// New creates a new SQLite recorder instance
func New(address string) *Recorder {
	return &Recorder{
		address: address,
	}
}

// Name returns the recorder type name
func (r *Recorder) Name() string {
	return "sqlite"
}

// Initialize prepares the SQLite recorder with configuration
func (r *Recorder) Initialize(ctx context.Context, settings Settings) error {
	// Set defaults
	r.settings = Settings{
		Path:      "gitlab-exporter-sqlite.db",
		WALMode:   true,
		BatchSize: 1000,
	}

	// Override with options if provided
	if settings.Path != "" {
		r.settings.Path = settings.Path
	}
	if !settings.WALMode {
		r.settings.WALMode = settings.WALMode
	}
	if settings.BatchSize != 0 {
		r.settings.BatchSize = settings.BatchSize
	}

	// Validate settings
	if r.settings.Path == "" {
		return fmt.Errorf("database path is required")
	}

	return nil
}

// Start opens the database connection and runs migrations
func (r *Recorder) Start(ctx context.Context) error {
	// Open database
	db, err := sql.Open("sqlite", r.settings.Path)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	r.db = db

	// Configure connection
	if err := r.configure(ctx); err != nil {
		return fmt.Errorf("configure database: %w", err)
	}

	// Run migrations
	if err := RunMigrations(ctx, r.db, "gitlab_ci"); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	return nil
}

// Stop closes the database connection
func (r *Recorder) Stop(ctx context.Context) error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

// CheckHealth checks if the database connection is alive
func (r *Recorder) CheckHealth(ctx context.Context) error {
	if r.db == nil {
		return fmt.Errorf("database not initialized")
	}

	return r.db.PingContext(ctx)
}
