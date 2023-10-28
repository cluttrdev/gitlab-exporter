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

To install `gitlab-clickhouse-exporter` you can download a 
[prebuilt binary][prebuilt-binaries] that matches your system, e.g.

```shell
# Download
OS=linux
ARCH=amd64
RELEASE_TAG=$(curl -sSfL https://api.github.com/repos/cluttrdev/gitlab-clickhouse-exporter/releases/latest | jq -r '.tag_name')
curl -sSfL https://github.com/cluttrdev/gitlab-clickhouse-exporter/releases/download/${RELEASE_TAG}/gitlab-clickhouse-exporter_${RELEASE_TAG}_${OS}_${ARCH}.tar.gz
# Install
tar -xf gitlab-clickhouse-exporter_*.tar.gz gitlab-clickhouse-exporter
install gitlab-clickhouse-exporter ~/.local/bin/gitlab-clickhouse-exporter
```

Alternatively, if you have the [Go][go-install] tools installed on your
machine, you can use

```shell
go install github.com/cluttrdev/gitlab-clickhouse-exporter@latest
```

## Usage

`gitlab-clickhouse-exporter` can either run in daemon mode or execute one-off
commands.

### Daemon Mode

To run `gitlab-clickhouse-exporter` in daemon mode use:

```shell
gitlab-clickhouse-exporter --config CONFIG_FILE run
```

This will periodically export data for updated pipelines of the configured projects,
see [Configuration](#configuration) for configuration options.

### Command Mode

`gitlab-clickhouse-exporter` supports a number of commands that can be executed
individually. Use the following to get an overview of available commands:

```shell
gitlab-clickhouse-exporter -h
```

## Configuration

Configuration options can be specified in a config file that is passed to the
application using the `--config` command-line flag.

For an overview of available configuration options and their default values,
see [configs/gitlab-clickhouse-exporter.yaml](./configs/gitlab-clickhouse-exporter.yaml).

Some common options can also be passed as command-line flags and/or environment
variables, with flags taking precedence.

| Flag                  | Environment Variable        | Default Value                 |
| ---                   | ---                         | ---                           |
| --gitlab-api-url      | `GLCHE_GITLAB_API_URL`      | `"https://gitlab.com/api/v4"` |
| --gitlab-api-token    | `GLCHE_GITLAB_API_TOKEN`    | **required**                  |
| --clickhouse-host     | `GLCHE_CLICKHOUSE_HOST`     | `"localhost"`                 |
| --clickhouse-port     | `GLCHE_CLICKHOUSE_PORT`     | `9000`                        |
| --clickhouse-database | `GLCHE_CLICKHOUSE_DATABASE` | `"default"`                   |
| --clickhouse-user     | `GLCHE_CLICKHOUSE_USER`     | `"default"`                   |
| --clickhouse-password | `GLCHE_CLICKHOUSE_PASSWORD` | `""`                          |

## Development Environment

To test the application during development or to just see what it has to offer,
a [docker-compose.yaml](./environments/dev/docker-compose.yaml) file is provided
that can be used to set up a simple environment consisting of a ClickHouse server
and a Grafana instance that includes some predefined dashboards.

To use this, simply change directory to `environments/dev/` and run:

```shell
docker compose up -d
```

Then, set the necessary environment variables and run `gitlab-clickhouse-exporter`
(either in daemon mode or using one-off commands):
```shell
export GLCHE_GITLAB_API_TOKEN=<your-gitlab-token>

gitlab-clickhouse-exporter run --projects <project-ids>
```

You can then login to Grafana on <http://localhost:3000> to explore the data.

## Acknowledgements

This project was inspired by [Maxime Visonneau's][github-mvisonneau] 
[gitlab-ci-pipeline-exporter][github-gcpe].

## License

This project is licensed under the [MIT License](./LICENSE).

[gitlab-api]: https://docs.gitlab.com/ee/api/rest/
[clickhouse]: https://clickhouse.com/
[go-install]: https://go.dev/doc/install
[prebuilt-binaries]: https://github.com/cluttrdev/gitlab-clickhouse-exporter/releases/latest
[github-mvisonneau]: https://github.com/mvisonneau
[github-gcpe]: https://github.com/mvisonneau/gitlab-ci-pipelines-exporter
