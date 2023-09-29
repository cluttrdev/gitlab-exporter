package config

import (
	"github.com/creasty/defaults"
)

type Config struct {
	GitLab     GitLab     `yaml:"gitlab"`
	ClickHouse ClickHouse `yaml:"clickhouse"`
	Projects   []Project  `yaml:"projects" default:"[]"`
}

type GitLab struct {
	Api struct {
		URL   string `yaml:"url" default:"https://gitlab.com/api/v4"`
		Token string `yaml:"token" default:""`
	} `yaml:"api"`

	Client struct {
		Rate struct {
			Limit float64 `yaml:"limit" default:"0.0"`
		} `yaml:"rate"`
	} `yaml:"client"`
}

type ClickHouse struct {
	Host     string `yaml:"host" default:"localhost"`
	Port     string `yaml:"port" default:"9000"`
	Database string `yaml:"database" default:"default"`
	User     string `yaml:"user" default:"default"`
	Password string `yaml:"password" default:""`
}

type Project struct {
	Id       int64           `yaml:"id"`
	Sections ProjectSections `yaml:"sections"`
	CatchUp  ProjectCatchUp  `yaml:"catch_up"`
}

type ProjectSections struct {
	Enabled bool `yaml:"enabled" default:"false"`
}

type ProjectCatchUp struct {
	Enabled       bool   `yaml:"enabled" default:"false"`
	UpdatedAfter  string `yaml:"updated_after" default:""`
	UpdatedBefore string `yaml:"updated_before" default:""`
}

func Default() *Config {
	var cfg Config

	defaults.MustSet(&cfg)

	return &cfg
}
