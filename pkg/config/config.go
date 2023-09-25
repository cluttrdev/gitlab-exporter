package config

import (
	"github.com/creasty/defaults"
)

type Config struct {
	GitLab     GitLab     `yaml:"gitlab"`
	ClickHouse ClickHouse `yaml:"clickhouse"`
}

type GitLab struct {
	Api struct {
		URL   string `default:"https://gitlab.com/api/v4" yaml:"url"`
		Token string `default:"" yaml:"token"`
	} `yaml:"api"`

	Client struct {
		Rate struct {
			Limit float64 `default:"0.0" yaml:"limit"`
		} `yaml:"rate"`
	} `yaml:"client"`
}

type ClickHouse struct {
	Host     string `default:"localhost" yaml:"host"`
	Port     string `default:"9000" yaml:"port"`
	Database string `default:"default" yaml:"database"`
	User     string `default:"default" yaml:"user"`
	Password string `default:"" yaml:"password"`
}

func Default() *Config {
	var cfg Config

	defaults.MustSet(&cfg)

	return &cfg
}
