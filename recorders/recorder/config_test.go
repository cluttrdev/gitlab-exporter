package recorder

import (
	"context"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestConfigLoaderLoadFromConfigs(t *testing.T) {
	// Create new registry
	registry := NewRegistry()

	registry.Register("test", func() Recorder {
		return &mockRecorder{name: "test-recorder"}
	})

	manager := NewManager()

	configs := []Config{
		{
			Type:    "test",
			Enabled: true,
			Address: "localhost:9100",
		},
	}

	ctx := context.Background()
	recorders, err := LoadFromConfigs(ctx, registry, configs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, r := range recorders {
		manager.Add(r)
	}

	if manager.Count() != 1 {
		t.Fatalf("expected 1 recorder, got %d", manager.Count())
	}

	recordersList := manager.List()
	if recordersList[0].Name() != "test-recorder" {
		t.Errorf("expected name %q, got %q", "test-recorder", recorders[0].Name())
	}
}

func TestConfigLoaderSkipsDisabled(t *testing.T) {
	// Create new registry
	registry := NewRegistry()

	registry.Register("test", func() Recorder {
		return &mockRecorder{name: "test-recorder"}
	})

	manager := NewManager()

	configs := []Config{
		{
			Type:    "test",
			Enabled: false, // disabled
			Address: "localhost:9100",
		},
	}

	ctx := context.Background()
	recorders, err := LoadFromConfigs(ctx, registry, configs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, r := range recorders {
		manager.Add(r)
	}

	if manager.Count() != 0 {
		t.Errorf("expected 0 recorders (disabled), got %d", manager.Count())
	}
}

func TestConfigLoaderErrorsOnMissingType(t *testing.T) {
	registry := NewRegistry()

	configs := []Config{
		{
			Type:    "", // missing type
			Enabled: true,
			Address: "localhost:9100",
		},
	}

	ctx := context.Background()
	_, err := LoadFromConfigs(ctx, registry, configs)
	if err == nil {
		t.Error("expected error for missing type")
	}
}

func TestConfigLoaderErrorsOnUnknownType(t *testing.T) {
	registry := NewRegistry()

	configs := []Config{
		{
			Type:    "unknown",
			Enabled: true,
			Address: "localhost:9100",
		},
	}

	ctx := context.Background()
	_, err := LoadFromConfigs(ctx, registry, configs)
	if err == nil {
		t.Error("expected error for unknown type")
	}
}

func TestConfigUnmarshal(t *testing.T) {
	yamlData := `
type: clickhouse
enabled: true
address: localhost:9000
settings:
  database: gitlab
  username: default
`

	var cfg Config
	err := yaml.Unmarshal([]byte(yamlData), &cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Type != "clickhouse" {
		t.Errorf("expected type %q, got %q", "clickhouse", cfg.Type)
	}

	if !cfg.Enabled {
		t.Error("expected enabled to be true")
	}

	if cfg.Address != "localhost:9000" {
		t.Errorf("expected address %q, got %q", "localhost:9000", cfg.Address)
	}

	if cfg.Settings == nil {
		t.Fatal("expected options to be populated")
	}

	if cfg.Settings["database"] != "gitlab" {
		t.Errorf("expected database %q, got %v", "gitlab", cfg.Settings["database"])
	}

	if cfg.Settings["username"] != "default" {
		t.Errorf("expected username %q, got %v", "default", cfg.Settings["username"])
	}
}
