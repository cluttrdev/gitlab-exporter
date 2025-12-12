package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func LoadFile(filename string, cfg *Config) error {
	file, err := os.Open(filepath.Clean(filename))
	if err != nil {
		return fmt.Errorf("error opening configuration file: %w", err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading configuration file: %w", err)
	}

	return Load(data, cfg)
}

func Load(data []byte, cfg *Config) error {
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("error parsing configuration: %w", err)
	}
	return nil
}
