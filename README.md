# gitlab-clickhouse-exporter

`gitlab-clickhouse-exporter` can help you build an observability and analytics
solution to gain insights into your CI pipelines by exporting data retrieved
from the [GitLab API][gitlab-api] to [ClickHouse][clickhouse].

---

**Note:** This project is in an early development stage, so functionality and
configuration options are limited.

---


## Installation

Make sure you have [Go][go-install] installed on your machine and run

```shell
go install https://github.com/cluttrdev/gitlab-clickhouse-exporter@latest
```

## Usage

After [installation](#installation) simply run

```shell
gitlab-clickhouse-exporter PROJECT_ID
```

This will periodically export data for updated pipelines of the specified project,
see [Configuration](#configuration) for configuration options.

## Configuration

| Environment Variable  | Default Value               |
| ---                   | ---                         |
| `GITLAB_API_URL`      | `https://gitlab.com/api/v4` |
| `GITLAB_API_TOKEN`    | **required**                |
| `CLICKHOUSE_HOST`     | `localhost`                 |
| `CLICKHOUSE_PORT`     | `9000`                      |
| `CLICKHOUSE_DATABASE` | `default`                   |
| `CLICKHOUSE_USER`     | `default`                   |
| `CLICKHOUSE_PASSWORD` | `""`                        |


## License

This project is licensed under the [MIT License](./LICENSE)

[gitlab-api]: https://docs.gitlab.com/ee/api/rest/
[clickhouse]: https://clickhouse.com/
[go-install]: https://go.dev/doc/install
