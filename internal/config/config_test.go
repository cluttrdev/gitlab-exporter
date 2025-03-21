package config_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"go.cluttr.dev/gitlab-exporter/internal/config"
)

func defaultConfig() config.Config {
	var cfg config.Config

	cfg.GitLab.Url = "https://gitlab.com"
	cfg.GitLab.Token = ""
	cfg.GitLab.Client.Rate.Limit = 0.0

	cfg.Endpoints = []config.Endpoint{}

	cfg.ProjectDefaults = defaultProjectSettings()

	cfg.Projects = []config.Project{}
	cfg.Namespaces = []config.Namespace{}

	cfg.HTTP.Enabled = true
	cfg.HTTP.Host = "127.0.0.1"
	cfg.HTTP.Port = "9100"
	cfg.HTTP.Debug = false

	cfg.Log.Level = "info"
	cfg.Log.Format = "text"

	return cfg
}

func defaultProjectSettings() config.ProjectSettings {
	var cfg config.ProjectSettings

	cfg.Export.Deployments.Enabled = true
	cfg.Export.MergeRequests.Enabled = true
	cfg.Export.MergeRequests.NoteEvents = true
	cfg.Export.Metrics.Enabled = true
	cfg.Export.Sections.Enabled = true
	cfg.Export.TestReports.Enabled = true
	cfg.Export.Traces.Enabled = true

	cfg.CatchUp.Enabled = false
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
      url: https://git.example.com
    `)

	var expected config.Config
	expected.GitLab.Url = "https://git.example.com"

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
      token: glpat-xxxxxxxxxxxxxxxxxxxx
      client:
        rate:
          limit: 20

    log:
      format: json
    `)

	expected := defaultConfig()
	expected.GitLab.Token = "glpat-xxxxxxxxxxxxxxxxxxxx"
	expected.GitLab.Client.Rate.Limit = 20
	expected.Log.Format = "json"

	cfg := defaultConfig()
	if err := config.Load(data, &cfg); err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	checkConfig(t, expected, cfg)
}

func TestLoad_WithProjectsEmptyDefaults(t *testing.T) {
	data := []byte(`
    project_defaults: {}
    projects:
      - id: 42
    `)

	expected := defaultConfig()
	expected.Projects = append(expected.Projects,
		config.Project{
			Id: 42,
			ProjectSettings: config.ProjectSettings{
				Export: config.ProjectExport{
					Deployments:   config.ProjectExportDeployments{Enabled: true},
					MergeRequests: config.ProjectExportMergeRequests{Enabled: true, NoteEvents: true},
					Metrics:       config.ProjectExportMetrics{Enabled: true},
					Sections:      config.ProjectExportSections{Enabled: true},
					TestReports:   config.ProjectExportTestReports{Enabled: true},
					Traces:        config.ProjectExportTraces{Enabled: true},
				},
				CatchUp: config.ProjectCatchUp{
					Enabled:       false,
					UpdatedAfter:  "",
					UpdatedBefore: "",
				},
			},
		},
	)

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
					Deployments: config.ProjectExportDeployments{
						Enabled: true,
					},
					Sections: config.ProjectExportSections{
						Enabled: true,
					},
					TestReports: config.ProjectExportTestReports{
						Enabled: false,
					},
					Traces: config.ProjectExportTraces{
						Enabled: true,
					},
					MergeRequests: config.ProjectExportMergeRequests{
						Enabled:    true,
						NoteEvents: true,
					},
					Metrics: config.ProjectExportMetrics{
						Enabled: true,
					},
				},
				CatchUp: config.ProjectCatchUp{
					Enabled:       true,
					UpdatedAfter:  "",
					UpdatedBefore: "",
				},
			},
			Id: 1337, // "foo/bar",
		},
		config.Project{
			ProjectSettings: config.ProjectSettings{
				Export: config.ProjectExport{
					Deployments: config.ProjectExportDeployments{
						Enabled: true,
					},
					Sections: config.ProjectExportSections{
						Enabled: false,
					},
					TestReports: config.ProjectExportTestReports{
						Enabled: true,
					},
					Traces: config.ProjectExportTraces{
						Enabled: false,
					},
					MergeRequests: config.ProjectExportMergeRequests{
						Enabled:    true,
						NoteEvents: true,
					},
					Metrics: config.ProjectExportMetrics{
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

func TestLoad_DataWithProjectDefaults(t *testing.T) {
	data := []byte(`
    project_defaults:
      export:
        traces:
          enabled: false
        metrics:
          enabled: false
      catch_up:
        enabled: true
    projects:
      - id: 314
      - id: 1337
        export:
          metrics:
            enabled: true
        catchup: {}
      - id: 42
        export:
          sections:
            enabled: false
          testreports:
            enabled: true
          traces:
            enabled: true
          metrics:
            enabled: false
        catch_up:
          enabled: true
          updated_after: "2019-03-15T08:00:00Z"
    `)

	expected := defaultConfig()
	expected.ProjectDefaults = config.ProjectSettings{
		Export: config.ProjectExport{
			Deployments:   config.ProjectExportDeployments{Enabled: true},
			Sections:      config.ProjectExportSections{Enabled: true},
			TestReports:   config.ProjectExportTestReports{Enabled: true},
			Traces:        config.ProjectExportTraces{Enabled: false},
			MergeRequests: config.ProjectExportMergeRequests{Enabled: true, NoteEvents: true},
			Metrics:       config.ProjectExportMetrics{Enabled: false},
		},
		CatchUp: config.ProjectCatchUp{
			Enabled:       true,
			UpdatedAfter:  "",
			UpdatedBefore: "",
		},
	}
	expected.Projects = append(expected.Projects,
		config.Project{
			ProjectSettings: config.ProjectSettings{
				Export: config.ProjectExport{
					Deployments:   config.ProjectExportDeployments{Enabled: true},
					Sections:      config.ProjectExportSections{Enabled: true},
					TestReports:   config.ProjectExportTestReports{Enabled: true},
					Traces:        config.ProjectExportTraces{Enabled: false},
					MergeRequests: config.ProjectExportMergeRequests{Enabled: true, NoteEvents: true},
					Metrics:       config.ProjectExportMetrics{Enabled: false},
				},
				CatchUp: config.ProjectCatchUp{
					Enabled:       true,
					UpdatedAfter:  "",
					UpdatedBefore: "",
				},
			},
			Id: 314,
		},
		config.Project{
			ProjectSettings: config.ProjectSettings{
				Export: config.ProjectExport{
					Deployments:   config.ProjectExportDeployments{Enabled: true},
					Sections:      config.ProjectExportSections{Enabled: true},
					TestReports:   config.ProjectExportTestReports{Enabled: true},
					Traces:        config.ProjectExportTraces{Enabled: false},
					MergeRequests: config.ProjectExportMergeRequests{Enabled: true, NoteEvents: true},
					Metrics:       config.ProjectExportMetrics{Enabled: true},
				},
				CatchUp: config.ProjectCatchUp{
					Enabled:       true,
					UpdatedAfter:  "",
					UpdatedBefore: "",
				},
			},
			Id: 1337,
		},
		config.Project{
			ProjectSettings: config.ProjectSettings{
				Export: config.ProjectExport{
					Deployments:   config.ProjectExportDeployments{Enabled: true},
					Sections:      config.ProjectExportSections{Enabled: false},
					TestReports:   config.ProjectExportTestReports{Enabled: true},
					Traces:        config.ProjectExportTraces{Enabled: true},
					MergeRequests: config.ProjectExportMergeRequests{Enabled: true, NoteEvents: true},
					Metrics:       config.ProjectExportMetrics{Enabled: false},
				},
				CatchUp: config.ProjectCatchUp{
					Enabled:       true,
					UpdatedAfter:  "2019-03-15T08:00:00Z",
					UpdatedBefore: "",
				},
			},
			Id: 42,
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
    http:
      host: "0.0.0.0"
      port: "8443"
    `)

	expected := defaultConfig()
	expected.HTTP.Host = "0.0.0.0"
	expected.HTTP.Port = "8443"

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

func TestLoad_WithNamespaces(t *testing.T) {
	data := []byte(`
    project_defaults:
      export:
        sections:
          enabled: false
        metrics:
          enabled: false
    namespaces:
      - id: akun73
        kind: user
        visibility: public
      - id: gitlab-exporter
        kind: group
        include_subgroups: true
        export:
          metrics:
            enabled: true
    `)

	expected := defaultConfig()
	expected.ProjectDefaults.Export.Sections.Enabled = false
	expected.ProjectDefaults.Export.Metrics.Enabled = false
	expected.Namespaces = []config.Namespace{
		{
			Id:              "akun73",
			Kind:            "user",
			Visibility:      "public",
			ProjectSettings: defaultProjectSettings(),
		},
		{
			Id:               "gitlab-exporter",
			Kind:             "group",
			IncludeSubgroups: true,
			ProjectSettings:  defaultProjectSettings(),
		},
	}
	expected.Namespaces[0].ProjectSettings.Export.Sections.Enabled = false
	expected.Namespaces[0].ProjectSettings.Export.Metrics.Enabled = false
	expected.Namespaces[1].ProjectSettings.Export.Sections.Enabled = false
	expected.Namespaces[1].ProjectSettings.Export.Metrics.Enabled = true

	cfg := config.Default()
	if err := config.Load(data, &cfg); err != nil {
		t.Fatal(err)
	}

	checkConfig(t, expected, cfg)
}
