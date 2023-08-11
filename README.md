# gitlab-clickhouse-exporter

`gitlab-clickhouse-exporter` can help you build an observability and analytics
solution to gain insights into your CI pipelines by exporting data retrieved
from the [GitLab API][gitlab-api] to [ClickHouse][clickhouse].

---

**Note:** This project is in an early development stage, so functionality and
configuration options are limited.

---

<p>
    <img src="./assets/project-overview.webp" />
    <img src="./assets/pipeline-trace.webp" />
</p>

## Installation

Make sure you have [Go][go-install] installed on your machine and run

```shell
go install https://github.com/cluttrdev/gitlab-clickhouse-exporter@latest
```

## Usage

`gitlab-clickhouse-exporter` can either run in daemon mode or execute one-off
commands.

### Daemon Mode

To run `gitlab-clickhouse-exporter` in daemon mode use:

```shell
gitlab-clickhouse-exporter run
```

This will periodically export data for updated pipelines of the configured projects,
see [Configuration](#configuration) for configuration options.

### Command Mode

`gitlab-clickhouse-exporter` supports a number of commands that can be excuted
individually. Use the following to get an overview of available command:

```shell
gitlab-clickhouse-exporter -h
```

## Configuration

| Environment Variable        | Default Value                 |
| ---                         | ---                           |
| `GLCHE_GITLAB_API_URL`      | `"https://gitlab.com/api/v4"` |
| `GLCHE_GITLAB_API_TOKEN`    | **required**                  |
| `GLCHE_CLICKHOUSE_HOST`     | `"localhost"`                 |
| `GLCHE_CLICKHOUSE_PORT`     | `9000`                        |
| `GLCHE_CLICKHOUSE_DATABASE` | `"default"`                   |
| `GLCHE_CLICKHOUSE_USER`     | `"default"`                   |
| `GLCHE_CLICKHOUSE_PASSWORD` | `""`                          |
| `GLCHE_PROJECTS`            | `[]`                          |


## License

This project is licensed under the [MIT License](./LICENSE).

[gitlab-api]: https://docs.gitlab.com/ee/api/rest/
[clickhouse]: https://clickhouse.com/
[go-install]: https://go.dev/doc/install
