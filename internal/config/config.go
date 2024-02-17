package config

import (
	"github.com/creasty/defaults"
)

type Config struct {
	GitLab    GitLab     `default:"{}" yaml:"gitlab"`
	Endpoints []Endpoint `default:"[]" yaml:"endpoints"`
	Projects  []Project  `default:"[]" yaml:"projects"`
	Server    Server     `default:"{}" yaml:"server"`
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

type Endpoint struct {
	Address string `default:"" yaml:"address"`
}

type Project struct {
	ProjectSettings `default:"{}" yaml:",inline"`

	Id int64 `yaml:"id"`
}

type ProjectSettings struct {
	Export  ProjectExport  `default:"{}" yaml:"export"`
	CatchUp ProjectCatchUp `default:"{}" yaml:"catch_up"`
}

type ProjectExport struct {
	Sections           ProjectExportSections           `default:"{}" yaml:"sections"`
	TestReports        ProjectExportTestReports        `default:"{}" yaml:"testreports"`
	Traces             ProjectExportTraces             `default:"{}" yaml:"traces"`
	LogEmbeddedMetrics ProjectExportLogEmbeddedMetrics `default:"{}" yaml:"log_embedded_metrics"`
}

type ProjectExportSections struct {
	Enabled bool `default:"true" yaml:"enabled"`
}

type ProjectExportTestReports struct {
	Enabled bool `default:"true" yaml:"enabled"`
}

type ProjectExportTraces struct {
	Enabled bool `default:"true" yaml:"enabled"`
}

type ProjectExportLogEmbeddedMetrics struct {
	Enabled bool `default:"true" yaml:"enabled"`
}

type ProjectCatchUp struct {
	Enabled       bool   `default:"false" yaml:"enabled"`
	Forced        bool   `default:"false" yaml:"forced"`
	UpdatedAfter  string `default:"" yaml:"updated_after"`
	UpdatedBefore string `default:"" yaml:"updated_before"`
}

type Server struct {
	Host string `default:"127.0.0.1" yaml:"host"`
	Port string `default:"8080" yaml:"port"`
}

func Default() Config {
	var cfg Config

	defaults.MustSet(&cfg)

	return cfg
}

func DefaultProjectSettings() ProjectSettings {
	var cfg ProjectSettings

	defaults.MustSet(&cfg)

	return cfg
}
