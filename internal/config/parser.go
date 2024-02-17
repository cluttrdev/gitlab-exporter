package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config implements the Unmarshaler interface
func (c *Config) UnmarshalYAML(v *yaml.Node) error {
	type _Config struct {
		GitLab GitLab `yaml:"gitlab"`

		Endpoints []Endpoint `yaml:"endpoints"`

		Projects []yaml.Node `yaml:"projects"`

		Server Server `yaml:"server"`

		Log Log `yaml:"log"`
	}

	var _cfg _Config
	_cfg.GitLab = c.GitLab
	_cfg.Endpoints = c.Endpoints
	_cfg.Server = c.Server
	_cfg.Log = c.Log

	if err := v.Decode(&_cfg); err != nil {
		return err
	}

	c.GitLab = _cfg.GitLab
	c.Endpoints = _cfg.Endpoints
	c.Server = _cfg.Server
	c.Log = _cfg.Log

	for _, n := range _cfg.Projects {
		p := Project{
			ProjectSettings: DefaultProjectSettings(),
		}

		if err := n.Decode(&p); err != nil {
			return nil
		}

		c.Projects = append(c.Projects, p)
	}

	return nil
}

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
