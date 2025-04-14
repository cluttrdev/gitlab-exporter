package config

import (
	"github.com/creasty/defaults"
)

// Config holds all the parameter settings for the application.
type Config struct {
	// GitLab client settings
	GitLab GitLab `default:"{}" yaml:"gitlab"`
	// List of recorder endpoints to export to
	Endpoints []Endpoint `default:"[]" yaml:"endpoints"`
	// Default settings for projects
	ProjectDefaults ProjectSettings `default:"{}" yaml:"project_defaults"`
	// List of project to export
	Projects []Project `default:"[]" yaml:"projects"`
	// List of namespaces of which to export projects
	Namespaces []Namespace `default:"[]" yaml:"namespaces"`
	// HTTP server settings
	HTTP HTTP `default:"{}" yaml:"http"`
	// Log configuration settings
	Log Log `default:"{}" yaml:"log"`
}

type GitLab struct {
	Url   string `default:"https://gitlab.com" yaml:"url"`
	Token string `default:"" yaml:"token"`

	Username string `default:"" yaml:"username"`
	Password string `default:"" yaml:"password"`

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
	Deployments   ProjectExportDeployments   `default:"{}" yaml:"deployments"`
	Issues        ProjectExportIssues        `default:"{}" yaml:"issues"`
	Jobs          ProjectExportJobs          `default:"{}" yaml:"jobs"`
	Sections      ProjectExportSections      `default:"{}" yaml:"sections"`
	TestReports   ProjectExportTestReports   `default:"{}" yaml:"testreports"`
	Reports       ProjectExportReports       `default:"{}" yaml:"reports"`
	Traces        ProjectExportTraces        `default:"{}" yaml:"traces"`
	Metrics       ProjectExportMetrics       `default:"{}" yaml:"metrics"`
	MergeRequests ProjectExportMergeRequests `default:"{}" yaml:"mergerequests"`
}

type ProjectExportDeployments struct {
	Enabled bool `default:"true" yaml:"enabled"`
}

type ProjectExportIssues struct {
	Enabled bool `default:"true" yaml:"enabled"`
}

type ProjectExportJobs struct {
	Properties ProjectExportJobsProperties `default:"{}" yaml:"properties"`
}

type ProjectExportJobsProperties struct {
	Enabled bool `default:"true" yaml:"enabled"`
}

type ProjectExportSections struct {
	Enabled bool `default:"true" yaml:"enabled"`
}

type ProjectExportTestReports struct {
	Enabled bool `default:"true" yaml:"enabled"`
}

type ProjectExportReports struct {
	Enabled bool `default:"false" yaml:"enabled"`

	Junit    ProjectExportReportsSettings `default:"{}" yaml:"junit"`
	Coverage ProjectExportReportsSettings `default:"{}" yaml:"coverage"`
}

type ProjectExportReportsSettings struct {
	Enabled bool     `default:"false" yaml:"enabled"`
	Paths   []string `default:"" yaml:"paths"`
}

type ProjectExportTraces struct {
	Enabled bool `default:"true" yaml:"enabled"`
}

type ProjectExportMetrics struct {
	Enabled bool `default:"true" yaml:"enabled"`
}

type ProjectExportMergeRequests struct {
	Enabled bool `default:"true" yaml:"enabled"`

	NoteEvents bool `default:"true" yaml:"note_events"`
}

type ProjectCatchUp struct {
	Enabled       bool   `default:"false" yaml:"enabled"`
	UpdatedAfter  string `default:"" yaml:"updated_after"`
	UpdatedBefore string `default:"" yaml:"updated_before"`
}

type Namespace struct {
	ProjectSettings `default:"{}" yaml:",inline"`

	Id   string `yaml:"id"`
	Kind string `default:"" yaml:"kind"`

	Visibility       string `default:"" yaml:"visibility"`
	WithShared       bool   `default:"false" yaml:"with_shared"`
	IncludeSubgroups bool   `default:"false" yaml:"include_subgroups"`

	ExcludeProjects []string `default:"[]" yaml:"exclude_projects"`
}

type HTTP struct {
	Enabled bool   `default:"true" yaml:"enabled"`
	Host    string `default:"127.0.0.1" yaml:"host"`
	Port    string `default:"9100" yaml:"port"`
	Debug   bool   `default:"false" yaml:"debug"`
}

type Log struct {
	Level  string `default:"info" yaml:"level"`
	Format string `default:"text" yaml:"format"`
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

func IsAuthedHTTPRequired(cfg Config) bool {
	if !cfg.ProjectDefaults.Export.Reports.Enabled {
		return false
	}

	for _, p := range cfg.Projects {
		if !p.Export.Reports.Enabled {
			continue
		}
		if p.Export.Reports.Junit.Enabled && len(p.Export.Reports.Junit.Paths) == 0 {
			return true
		}
	}

	for _, n := range cfg.Namespaces {
		if !n.Export.Reports.Enabled {
			continue
		}
		if n.Export.Reports.Junit.Enabled && len(n.Export.Reports.Junit.Paths) == 0 {
			return true
		}
	}

	return false
}
