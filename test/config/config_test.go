package config_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/cluttrdev/gitlab-exporter/internal/config"
)

func defaultConfig() config.Config {
	var cfg config.Config

	cfg.GitLab.Api.URL = "https://gitlab.com/api/v4"
	cfg.GitLab.Api.Token = ""
	cfg.GitLab.Client.Rate.Limit = 0.0

	cfg.Endpoints = []config.Endpoint{}

	cfg.Projects = []config.Project{}

	cfg.Server.Host = "127.0.0.1"
	cfg.Server.Port = "8080"

	return cfg
}

func defaultProjectSettings() config.ProjectSettings {
	var cfg config.ProjectSettings

	cfg.Export.Sections.Enabled = true
	cfg.Export.TestReports.Enabled = true
	cfg.Export.Traces.Enabled = true
	cfg.Export.LogEmbeddedMetrics.Enabled = true

	cfg.CatchUp.Enabled = false
	cfg.CatchUp.Forced = false
	cfg.CatchUp.UpdatedAfter = ""
	cfg.CatchUp.UpdatedBefore = ""

	return cfg
}

func checkConfig(t *testing.T, want interface{}, got interface{}) {
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Config mismatch (-want +got):\n%s", diff)
	}
}

func Test_NewDefault(t *testing.T) {
	expected := defaultConfig()

	cfg := config.Default()

	checkConfig(t, expected, cfg)
}

func Test_NewDefaultProjectSettings(t *testing.T) {
	expected := defaultProjectSettings()

	cfg := config.DefaultProjectSettings()

	checkConfig(t, expected, cfg)
}

func TestLoad_EmptyData(t *testing.T) {
	data := []byte{}

	var expected config.Config

	var cfg config.Config
	if err := config.Load(data, &cfg); err != nil {
		t.Errorf("Expected no error when loading empty data, got: %v", err)
	}

	checkConfig(t, &expected, &cfg)
}

func TestLoad_PartialData(t *testing.T) {
	data := []byte(`
    gitlab:
      api:
        url: https://git.example.com/api/v4
    `)

	var expected config.Config
	expected.GitLab.Api.URL = "https://git.example.com/api/v4"

	var cfg config.Config
	if err := config.Load(data, &cfg); err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	checkConfig(t, &expected, &cfg)
}

func TestLoad_UnknownData(t *testing.T) {
	data := []byte(`
    global:
      unknown: true

    gitlab:
      unknown:
        answer: 42
    `)

	var expected config.Config

	var cfg config.Config
	if err := config.Load(data, &cfg); err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	checkConfig(t, &expected, &cfg)
}

func TestLoad_InvalidData(t *testing.T) {
	data := []byte(`
    invalid:
      data: 1
      data: 2
    `)

	var cfg config.Config
	if err := config.Load(data, &cfg); err != nil {
		t.Error("Expected error when loading invalid data, got `nil`")
	}
}

func TestLoad_DataWithDefaults(t *testing.T) {
	data := []byte(`
    gitlab:
      api:
        token: glpat-xxxxxxxxxxxxxxxxxxxx
      client:
        rate:
          limit: 20
    `)

	expected := defaultConfig()
	expected.GitLab.Api.Token = "glpat-xxxxxxxxxxxxxxxxxxxx"
	expected.GitLab.Client.Rate.Limit = 20

	cfg := defaultConfig()
	if err := config.Load(data, &cfg); err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	checkConfig(t, expected, cfg)
}

func TestLoad_DataWithProjects(t *testing.T) {
	data := []byte(`
    projects:
      - id: 314
      - id: 1337  # foo/bar
        export:
          sections:
            enabled: true
          testreports:
            enabled: false
          traces:
            enabled: true
        catch_up:
          enabled: true
          forced: true
      - id: 42
        export:
          sections:
            enabled: false
          testreports:
            enabled: true
          traces:
            enabled: false
        catch_up:
          enabled: true
          updated_after: "2019-03-15T08:00:00Z"
    `)

	expected := defaultConfig()
	expected.Projects = append(expected.Projects,
		config.Project{
			ProjectSettings: defaultProjectSettings(),
			Id:              314,
		},
		config.Project{
			ProjectSettings: config.ProjectSettings{
				Export: config.ProjectExport{
					Sections: config.ProjectExportSections{
						Enabled: true,
					},
					TestReports: config.ProjectExportTestReports{
						Enabled: false,
					},
					Traces: config.ProjectExportTraces{
						Enabled: true,
					},
					LogEmbeddedMetrics: config.ProjectExportLogEmbeddedMetrics{
						Enabled: true,
					},
				},
				CatchUp: config.ProjectCatchUp{
					Enabled:       true,
					Forced:        true,
					UpdatedAfter:  "",
					UpdatedBefore: "",
				},
			},
			Id: 1337, // "foo/bar",
		},
		config.Project{
			ProjectSettings: config.ProjectSettings{
				Export: config.ProjectExport{
					Sections: config.ProjectExportSections{
						Enabled: false,
					},
					TestReports: config.ProjectExportTestReports{
						Enabled: true,
					},
					Traces: config.ProjectExportTraces{
						Enabled: false,
					},
					LogEmbeddedMetrics: config.ProjectExportLogEmbeddedMetrics{
						Enabled: true,
					},
				},
				CatchUp: config.ProjectCatchUp{
					Enabled:       true,
					UpdatedAfter:  "2019-03-15T08:00:00Z",
					UpdatedBefore: "",
				},
			},
			Id: 42, // "42",
		},
	)

	cfg := defaultConfig()
	if err := config.Load(data, &cfg); err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	checkConfig(t, expected, cfg)
}

func TestLoad_DataWithCustomServerAddress(t *testing.T) {
	data := []byte(`
    server:
      host: "0.0.0.0"
      port: "8443"
    `)

	expected := defaultConfig()
	expected.Server.Host = "0.0.0.0"
	expected.Server.Port = "8443"

	cfg := config.Default()
	if err := config.Load(data, &cfg); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	checkConfig(t, expected, cfg)
}

func TestLoad_DataWithEndpoints(t *testing.T) {
	data := []byte(`
    endpoints:
      - address: "127.0.0.1:36275"
    `)

	expected := defaultConfig()
	expected.Endpoints = []config.Endpoint{
		{Address: "127.0.0.1:36275"},
	}

	cfg := config.Default()
	if err := config.Load(data, &cfg); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	checkConfig(t, expected, cfg)
}
