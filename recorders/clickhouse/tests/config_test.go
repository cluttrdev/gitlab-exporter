package config_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"go.cluttr.dev/gitlab-exporter-clickhouse-recorder/internal/config"
)

func defaultConfig() config.Config {
	var cfg config.Config

	cfg.ClickHouse.Host = "127.0.0.1"
	cfg.ClickHouse.Port = "9000"
	cfg.ClickHouse.Database = "default"
	cfg.ClickHouse.User = "default"
	cfg.ClickHouse.Password = ""

	cfg.Server.Host = "0.0.0.0"
	cfg.Server.Port = "0"

	cfg.HTTP.Enabled = true
	cfg.HTTP.Host = "127.0.0.1"
	cfg.HTTP.Port = "9100"
	cfg.HTTP.Debug = false

	cfg.Log.Level = "info"
	cfg.Log.Format = "text"

	return cfg
}

func checkConfig(t *testing.T, want interface{}, got interface{}) {
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Config mismatch (-want +got):\n%s", diff)
	}
}

func Test_NewDefault(t *testing.T) {
	expected := defaultConfig()

	var cfg config.Config
	config.SetDefaults(&cfg)

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
    clickhouse:
      database: "gitlab_ci"
    `)

	var expected config.Config
	expected.ClickHouse.Database = "gitlab_ci"

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

    server:
      question:
        answer: "42"
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
      non_unique_key: 1
      non_unique_key: 2
    `)

	var cfg config.Config
	if err := config.Load(data, &cfg); err != nil {
		t.Error("Expected error when loading invalid data, got `nil`")
	}
}

func TestLoad_DataWithDefaults(t *testing.T) {
	data := []byte(`
    clickhouse:
      user: gitlab-exporter
      password: supersecret

    server:
      port: 36275

    http:
      host: 0.0.0.0
      port: 9443

    log:
      format: json
    `)

	expected := defaultConfig()
	expected.ClickHouse.User = "gitlab-exporter"
	expected.ClickHouse.Password = "supersecret"
	expected.Server.Port = "36275"
	expected.HTTP.Host = "0.0.0.0"
	expected.HTTP.Port = "9443"
	expected.Log.Format = "json"

	cfg := defaultConfig()
	if err := config.Load(data, &cfg); err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	checkConfig(t, expected, cfg)
}
