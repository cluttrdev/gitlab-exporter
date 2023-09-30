package config_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/config"
)

func defaultConfig() *config.Config {
	var cfg config.Config

	cfg.GitLab.Api.URL = "https://gitlab.com/api/v4"
	cfg.GitLab.Api.Token = ""
	cfg.GitLab.Client.Rate.Limit = 0.0
	cfg.ClickHouse.Host = "localhost"
	cfg.ClickHouse.Port = "9000"
	cfg.ClickHouse.Database = "default"
	cfg.ClickHouse.User = "default"
	cfg.ClickHouse.Password = ""

	cfg.Projects = []config.Project{}

	return &cfg
}

func checkConfig(t *testing.T, want *config.Config, got *config.Config) {
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Config mismatch (-want +got):\n%s", diff)
	}
}

func Test_NewDefault(t *testing.T) {
	expected := defaultConfig()

	cfg := config.Default()

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

    clickhouse:
      host: clickhouse.example.com
    `)

	expected := defaultConfig()
	expected.GitLab.Api.Token = "glpat-xxxxxxxxxxxxxxxxxxxx"
	expected.GitLab.Client.Rate.Limit = 20
	expected.ClickHouse.Host = "clickhouse.example.com"

	cfg := defaultConfig()
	if err := config.Load(data, cfg); err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	checkConfig(t, expected, cfg)
}

func TestLoad_DataWithProjects(t *testing.T) {
	data := []byte(`
    projects:
      - id: 1337  # foo/bar
        sections:
          enabled: true
        testreports:
          enabled: false
        traces:
          enabled: true
        catch_up:
          enabled: true
      - id: 42
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
			Id: 1337, // "foo/bar",
			Sections: config.ProjectSections{
				Enabled: true,
			},
            TestReports: config.ProjectTestReports{
                Enabled: false,
            },
            Traces: config.ProjectTraces{
                Enabled: true,
            },
			CatchUp: config.ProjectCatchUp{
				Enabled:       true,
				UpdatedAfter:  "",
				UpdatedBefore: "",
			},
		},
		config.Project{
			Id: 42, // "42",
			Sections: config.ProjectSections{
				Enabled: false,
			},
            TestReports: config.ProjectTestReports{
                Enabled: true,
            },
            Traces: config.ProjectTraces{
                Enabled: false,
            },
			CatchUp: config.ProjectCatchUp{
				Enabled:       true,
				UpdatedAfter:  "2019-03-15T08:00:00Z",
				UpdatedBefore: "",
			},
		},
	)

	cfg := defaultConfig()
	if err := config.Load(data, cfg); err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	checkConfig(t, expected, cfg)
}
