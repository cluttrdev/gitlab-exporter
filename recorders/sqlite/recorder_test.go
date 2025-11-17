package sqlite

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestNew(t *testing.T) {
	r := New("unix:///tmp/sqlite.sock")
	if r == nil {
		t.Fatal("New() returned nil")
	}

	if r.Name() != "sqlite" {
		t.Errorf("Name() = %s, want sqlite", r.Name())
	}
}

func TestRecorder_Initialize(t *testing.T) {
	tests := []struct {
		name      string
		config    string
		wantErr   bool
		checkPath string
	}{
		{
			name: "valid config with path in settings",
			config: `
path: /tmp/override.db
wal_mode: false
batch_size: 500
`,
			wantErr:   false,
			checkPath: "/tmp/override.db",
		},
		{
			name: "empty path",
			config: `
path: ""
`,
			wantErr:   false,
			checkPath: "gitlab-exporter-sqlite.db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New("unix:///tmp/sqlite.sock")

			var settings Settings
			_ = yaml.Unmarshal([]byte(tt.config), &settings)
			err := r.Initialize(context.Background(), settings)

			if (err != nil) != tt.wantErr {
				t.Errorf("Initialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.checkPath != "" {
				if r.settings.Path != tt.checkPath {
					t.Errorf("Path = %s, want %s", r.settings.Path, tt.checkPath)
				}
			}
		})
	}
}

func TestRecorder_Lifecycle(t *testing.T) {
	// Create a temporary directory for test database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create recorder
	r := New("unix:///tmp/sqlite.sock")

	// Initialize
	config := map[string]interface{}{
		"path":     dbPath,
		"wal_mode": false, // Disable WAL for simpler testing
	}
	configBytes, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}
	var settings Settings
	_ = yaml.Unmarshal(configBytes, &settings)

	ctx := context.Background()
	err = r.Initialize(ctx, settings)
	if err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	// Start (this creates the database and runs migrations)
	err = r.Start(ctx)
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Verify database file was created
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Errorf("Database file was not created at %s", dbPath)
	}

	// Health check
	err = r.CheckHealth(ctx)
	if err != nil {
		t.Errorf("CheckHealth() error = %v", err)
	}

	// Stop
	err = r.Stop(ctx)
	if err != nil {
		t.Errorf("Stop() error = %v", err)
	}

	// Health check after stop should fail
	err = r.CheckHealth(ctx)
	if err == nil {
		t.Error("CheckHealth() after Stop() should return error")
	}
}

func TestRecorder_Health_NotInitialized(t *testing.T) {
	r := New("unix:///tmp/sqlite.sock")
	ctx := context.Background()

	err := r.CheckHealth(ctx)
	if err == nil {
		t.Error("CheckHealth() on uninitialized recorder should return error")
	}
}

func TestRecorder_Stop_NotStarted(t *testing.T) {
	r := New("unix:///tmp/sqlite.sock")
	ctx := context.Background()

	err := r.Stop(ctx)
	if err != nil {
		t.Errorf("Stop() on non-started recorder should not error, got: %v", err)
	}
}
