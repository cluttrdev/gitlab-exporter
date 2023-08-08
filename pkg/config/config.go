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
	DefaultGitLabApiUrl string = "https://gitlab.com/api/v4"

	DefaultClickHouseHost     string = "localhost"
	DefaultClickHousePort     int    = 9000
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

	gl_url := getEnv("GITLAB_API_URL", "https://gitlab.com/api/v4")
	gl_token := getEnv("GITLAB_API_TOKEN", "")

	gl := GitLab{
		URL:   gl_url,
		Token: gl_token,
	}

	ch_host := getEnv("CLICKHOUSE_HOST", DefaultClickHouseHost)
	ch_port, err := strconv.ParseInt(getEnv("CLICKHOUSE_PORT", fmt.Sprintf("%d", DefaultClickHousePort)), 10, 16)
	if err != nil {
		return nil, fmt.Errorf("Failed to load clickhouse port: %w", err)
	}
	ch_database := getEnv("CLICKHOUSE_DATABASE", DefaultClickHouseDatabase)
	ch_user := getEnv("CLICKHOUSE_USER", DefaultClickHouseUser)
	ch_password := getEnv("CLICKHOUSE_PASSWORD", DefaultClickHousePassword)

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
