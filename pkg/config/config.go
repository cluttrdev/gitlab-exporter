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
	URL   string `yaml:"url"`
	Token string `yaml:"token"`
}

type ClickHouse struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

const (
	DEFAULT_GITLAB_API_URL string = "https://gitlab.com/api/v4"

	DEFAULT_CLICKHOUSE_HOST     string = "localhost"
	DEFAULT_CLICKHOUSE_PORT     string = "9000"
	DEFAULT_CLICKHOUSE_DATABASE string = "default"
	DEFAULT_CLICKHOUSE_USER     string = "default"
	DEFAULT_CLICKHOUSE_PASSWORD string = ""
)

func LoadEnv() (*Config, error) {
	getEnv := func(key string, defaultVal string) string {
		val, ok := os.LookupEnv(key)
		if !ok {
			val = defaultVal
		}
		return val
	}

	gl_url := getEnv("GITLAB_API_URL", "https://gitlab.com/api/v4")
	gl_token := getEnv("GITLAB_API_TOKEN", "")

	gl := GitLab{
		URL:   gl_url,
		Token: gl_token,
	}

	ch_host := getEnv("CLICKHOUSE_HOST", DEFAULT_CLICKHOUSE_HOST)
	ch_port, err := strconv.ParseInt(getEnv("CLICKHOUSE_PORT", DEFAULT_CLICKHOUSE_PORT), 10, 16)
	if err != nil {
		return nil, fmt.Errorf("Failed to load clickhouse port: %w", err)
	}
	ch_database := getEnv("CLICKHOUSE_DATABASE", DEFAULT_CLICKHOUSE_DATABASE)
	ch_user := getEnv("CLICKHOUSE_USER", DEFAULT_CLICKHOUSE_USER)
	ch_password := getEnv("CLICKHOUSE_PASSWORD", DEFAULT_CLICKHOUSE_PASSWORD)

	ch := ClickHouse{
		Host:     ch_host,
		Port:     int(ch_port),
		Database: ch_database,
		User:     ch_user,
		Password: ch_password,
	}

	return &Config{
		GitLab:     gl,
		ClickHouse: ch,
	}, nil
}
