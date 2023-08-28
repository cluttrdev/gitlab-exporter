package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	GitLab     GitLab     `yaml:"gitlab"`
	ClickHouse ClickHouse `yaml:"clickhouse"`
}

type GitLab struct {
	Api struct {
		URL   string `yaml:"url"`
		Token string `yaml:"token"`
	} `yaml:"api"`

	Client struct {
		RequestsPerSecond int64 `yaml:"requests_per_second"`
	} `yaml:"client"`
}

type ClickHouse struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

const (
	DefaultGitLabApiUrl   string = "https://gitlab.com/api/v4"
	DefaultGitLabApiToken string = ""

	DefaultGitLabClientRequestsPerSecond int64 = 0

	DefaultClickHouseHost     string = "localhost"
	DefaultClickHousePort     string = "9000"
	DefaultClickHouseDatabase string = "default"
	DefaultClickHouseUser     string = "default"
	DefaultClickHousePassword string = ""
)

func LoadEnv() (*Config, error) {
	getEnv := func(key string, defaultVal string) string {
		val, ok := os.LookupEnv(key)
		if !ok {
			val = defaultVal
		}
		return val
	}

	gl := GitLab{}
	gl.Api.URL = getEnv("GLCHE_GITLAB_API_URL", DefaultGitLabApiUrl)
	gl.Api.Token = getEnv("GLCHE_GITLAB_API_TOKEN", DefaultGitLabApiToken)

	gl_rps := getEnv("GLCHE_GITLAB_CLIENT_REQUESTSPERSECOND", "0")
	val, err := strconv.ParseInt(gl_rps, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing environment variables: %w", err)
	}
	gl.Client.RequestsPerSecond = val

	ch_host := getEnv("GLCHE_CLICKHOUSE_HOST", DefaultClickHouseHost)
	ch_port := getEnv("GLCHE_CLICKHOUSE_PORT", DefaultClickHousePort)
	ch_database := getEnv("GLCHE_CLICKHOUSE_DATABASE", DefaultClickHouseDatabase)
	ch_user := getEnv("GLCHE_CLICKHOUSE_USER", DefaultClickHouseUser)
	ch_password := getEnv("GLCHE_CLICKHOUSE_PASSWORD", DefaultClickHousePassword)

	ch := ClickHouse{
		Host:     ch_host,
		Port:     ch_port,
		Database: ch_database,
		User:     ch_user,
		Password: ch_password,
	}

	return &Config{
		GitLab:     gl,
		ClickHouse: ch,
	}, nil
}
