package config

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
		Rate struct {
			Limit float64 `yaml:"limit"`
		} `yaml:"rate"`
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

	DefaultGitLabClientRateLimit float64 = 0

	DefaultClickHouseHost     string = "localhost"
	DefaultClickHousePort     string = "9000"
	DefaultClickHouseDatabase string = "default"
	DefaultClickHouseUser     string = "default"
	DefaultClickHousePassword string = ""
)
