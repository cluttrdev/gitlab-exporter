package recorder

import (
	"context"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Config represents the configuration for a recorder instance
type Config struct {
	// Type identifies which recorder implementation to use (e.g., "clickhouse", "sqlite")
	Type string `yaml:"type"`

	// Enabled determines if this recorder should be active
	Enabled bool `yaml:"enabled"`

	// Address is the connection network address or unix socket file path
	Address string `yaml:"address,omitempty"`

	// Settings contains recorder-specific configuration as raw YAML
	// Each recorder implementation unmarshals this according to its needs
	Settings map[string]any `yaml:"settings,omitempty"`
}

// LoadFromConfigs creates and initializes recorders from configuration
func LoadFromConfigs(ctx context.Context, registry *Registry, configs []Config) ([]Recorder, error) {
	var recorders []Recorder

	for i, cfg := range configs {
		if !cfg.Enabled {
			continue
		}

		if cfg.Type == "" {
			return nil, fmt.Errorf("recorder config %d: type is required", i)
		}

		// Create recorder instance
		rec, err := registry.Create(cfg.Type)
		if err != nil {
			return nil, fmt.Errorf("recorder config %d: %w", i, err)
		}

		// Marshal the entire config to YAML for the recorder to parse
		configBytes, err := yaml.Marshal(cfg)
		if err != nil {
			return nil, fmt.Errorf("recorder config %d: failed to marshal config: %w", i, err)
		}

		// Initialize the recorder with its config
		if err := rec.Initialize(ctx, configBytes); err != nil {
			return nil, fmt.Errorf("recorder %q: initialization failed: %w", rec.Name(), err)
		}

		recorders = append(recorders, rec)
	}

	return recorders, nil
}
