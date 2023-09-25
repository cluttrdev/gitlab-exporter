package config_test

import (
	// "reflect"
	"testing"

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

	return &cfg
}

func checkConfig(t *testing.T, want *config.Config, got *config.Config) {
	if want.GitLab.Api.URL != got.GitLab.Api.URL {
		t.Errorf("Expected GitLab.Api.URL to be %s, got %s", got.GitLab.Api.URL, want.GitLab.Api.URL)
	}

	if want.GitLab.Api.Token != got.GitLab.Api.Token {
		t.Errorf("Expected GitLab.Api.Token to be %s, got %s", got.GitLab.Api.Token, want.GitLab.Api.Token)
	}

	if want.GitLab.Client.Rate.Limit != got.GitLab.Client.Rate.Limit {
		t.Errorf("Expected GitLab.Client.Rate.Limit to be %f, got %f", got.GitLab.Client.Rate.Limit, want.GitLab.Client.Rate.Limit)
	}

	if want.ClickHouse.Host != got.ClickHouse.Host {
		t.Errorf("Expected ClickHouse.Host to be %s, got %s", got.ClickHouse.Host, want.ClickHouse.Host)
	}

	if want.ClickHouse.Port != got.ClickHouse.Port {
		t.Errorf("Expected ClickHouse.Port to be %s, got %s", got.ClickHouse.Port, want.ClickHouse.Port)
	}

	if want.ClickHouse.Database != got.ClickHouse.Database {
		t.Errorf("Expected ClickHouse.Database to be %s, got %s", got.ClickHouse.Database, want.ClickHouse.Database)
	}

	if want.ClickHouse.User != got.ClickHouse.User {
		t.Errorf("Expected ClickHouse.User to be %s, got %s", got.ClickHouse.User, want.ClickHouse.User)
	}

	if want.ClickHouse.Password != got.ClickHouse.Password {
		t.Errorf("Expected ClickHouse.Password to be %s, got %s", got.ClickHouse.Password, want.ClickHouse.Password)
	}

	// gotValue := reflect.ValueOf(*cfg)
	// wantValue := reflect.ValueOf(expected)
	// if gotValue.Interface() != wantValue.Interface() {
	//     t.Errorf("Expected %+v, got %+v", wantValue,wantValue)
	// }
	//
	// for i := 0; i < gotValue.NumField(); i++ {
	//     name := gotValue.Type().Field(i).Name
	//     got := gotValue.FieldByName(name).Interface()
	//     want := wantValue.FieldByName(name).Interface()
	//     if got != want {
	//         t.Errorf("Expected %s to be %v, got %v", name, want, got)
	//     }
	// }
}

func Test_NewDefault(t *testing.T) {
	cfg := config.Default()

	expected := defaultConfig()

	checkConfig(t, expected, cfg)
}

func TestLoad_EmptyData(t *testing.T) {
	data := []byte{}

	var cfg config.Config
	if err := config.Load(data, &cfg); err != nil {
		t.Errorf("Expected no error when loading empty data, got: %v", err)
	}

	var expected config.Config

	checkConfig(t, &expected, &cfg)
}

func TestLoad_PartialData(t *testing.T) {
	data := []byte(`
    gitlab:
      api:
        url: https://git.example.com/api/v4
    `)

	var cfg config.Config
	if err := config.Load(data, &cfg); err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	var expected config.Config
	expected.GitLab.Api.URL = "https://git.example.com/api/v4"

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

	var cfg config.Config
	if err := config.Load(data, &cfg); err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	var expected config.Config

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

	cfg := defaultConfig()
	if err := config.Load(data, cfg); err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expected := defaultConfig()
	expected.GitLab.Api.Token = "glpat-xxxxxxxxxxxxxxxxxxxx"
	expected.GitLab.Client.Rate.Limit = 20
	expected.ClickHouse.Host = "clickhouse.example.com"

	checkConfig(t, expected, cfg)
}
