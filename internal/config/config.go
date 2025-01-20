package config

import (
	"os"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
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

	OAuth GitLabOAuth `default:"{}" yaml:"oauth"`

	Client struct {
		Rate struct {
			Limit float64 `default:"0.0" yaml:"limit"`
		} `yaml:"rate"`
	} `yaml:"client"`
}

type GitLabOAuth struct {
	GitLabOAuthSecrets `default:"{}" yaml:",inline"`

	SecretsFile string `default:"" yaml:"secrets_file"`
	FlowType    string `default:"" yaml:"flow_type"`
}

type GitLabOAuthSecrets struct {
	ClientId     string `default:"" yaml:"client_id"`
	ClientSecret string `default:"" yaml:"client_secret"`

	Username string `default:"" yaml:"username"`
	Password string `default:"" yaml:"password"`
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
	Sections      ProjectExportSections      `default:"{}" yaml:"sections"`
	TestReports   ProjectExportTestReports   `default:"{}" yaml:"testreports"`
	Reports       ProjectExportReports       `default:"{}" yaml:"reports"`
	Traces        ProjectExportTraces        `default:"{}" yaml:"traces"`
	Metrics       ProjectExportMetrics       `default:"{}" yaml:"metrics"`
	MergeRequests ProjectExportMergeRequests `default:"{}" yaml:"mergerequests"`
}

type ProjectExportSections struct {
	Enabled bool `default:"true" yaml:"enabled"`
}

type ProjectExportTestReports struct {
	Enabled bool `default:"true" yaml:"enabled"`
}

type ProjectExportReports struct {
	Enabled bool `default:"false" yaml:"enabled"`

	Junit ProjectExportReportsJunit `default:"{}" yaml:"junit"`
}

type ProjectExportReportsJunit struct {
	Enabled bool `default:"false" yaml:"enabled"`
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
	if cfg.ProjectDefaults.Export.Reports.Enabled {
		return true
	}

	for _, p := range cfg.Projects {
		if p.Export.Reports.Enabled {
			return true
		}
	}

	for _, n := range cfg.Namespaces {
		if n.Export.Reports.Enabled {
			return true
		}
	}

	return false
}

func LoadOAuthSecretsFile(path string) (GitLabOAuthSecrets, error) {
	file, err := os.Open(path)
	if err != nil {
		return GitLabOAuthSecrets{}, err
	}

	var secrets GitLabOAuthSecrets
	if err := yaml.NewDecoder(file).Decode(&secrets); err != nil {
		return GitLabOAuthSecrets{}, err
	}

	return secrets, nil
}
