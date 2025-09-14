# gitlab-exporter-clickhouse-recorder

`gitlab-exporter-clickhouse-recorder` serves a gRPC endpoint that can receive data
from a [gitlab-exporter](https://gitlab.com/akun73/gitlab-exporter) and records it
in a [ClickHouse](https://clickhouse.com) database.

## Dashboards

For some Grafana dashboards and screenshots see [here](https://gitlab.com/gitlab-exporter/grafana-dashboards).

## Installation

To install `gitlab-exporter-clickhouse-recorder` you can download a 
[prebuilt binary](https://gitlab.com/gitlab-exporter/clickhouse-recorder/-/releases)
that matches your system, e.g.

```shell
# download latest release archive
RELEASES_URL=https://gitlab.com/api/v4/projects/gitlab-exporter%2Fclickhouse-recorder/releases
RELEASE_TAG=$(curl -sSfL ${RELEASES_URL} | jq -r '.[0].tag_name')
curl -sSfL ${RELEASES_URL}/${RELEASE_TAG}/downloads/gitlab-exporter-clickhouse-recorder_${RELEASE_TAG}_linux_amd64.tar.gz \
    -o /tmp/gitlab-exporter-clickhouse-recorder.tar.gz

# extract executable binary into install dir (must exist)
INSTALL_DIR=$HOME/.local/bin
tar -C ${INSTALL_DIR} -zxof /tmp/gitlab-exporter-clickhouse-recorder.tar.gz gitlab-exporter-clickhouse-recorder

# check
${INSTALL_DIR}/gitlab-exporter-clickhouse-recorder version
```

## Usage

`gitlab-exporter-clickhouse-recorder` can either run in server mode or execute one-off
commands.

### Server Mode

To run `gitlab-exporter-clickhouse-recorder` in server mode use:

```shell
gitlab-exporter-clickhouse-recorder run --config CONFIG_FILE
```

This will start a gRPC server that exports recorded data to the configured 
ClickHouse database. See [Configuration](#configuration) for configuration options.

### Command Mode

`gitlab-exporter-clickhouse-recorder` supports commands that can be executed
individually. Use the following to get an overview of available commands:

```shell
gitlab-exporter-clickhouse-recorder -h
```

## Configuration

Configuration options can be specified in a config file that is passed to the
application using the `--config` command-line flag.

For an overview of available configuration options and their default values,
see [configs/gitlab-exporter-clickhouse-recorder.yaml](./configs/gitlab-exporter-clickhouse-recorder.yaml).

Common options can also be overridden with command-line flags and/or environment
variables, where flags take precedence.

| Flag                  | Environment Variable        | Default Value |
| ---                   | ---                         | ---           |
| # global options      |                             |               |
| --clickhouse-host     | `GLCHR_CLICKHOUSE_HOST`     | `"127.0.0.1"` |
| --clickhouse-port     | `GLCHR_CLICKHOUSE_PORT`     | `"9000"`      |
| --clickhouse-database | `GLCHR_CLICKHOUSE_DATABASE` | `"default"`   |
| --clickhouse-user     | `GLCHR_CLICKHOUSE_USER`     | `"default"`   |
| --clickhouse-password | `GLCHR_CLICKHOUSE_PASSWORD` | `""`          |
| # run options         |                             |               |
| --server-host         | `GLCHR_SERVER_HOST`         | `"0.0.0.0"`   |
| --server-port         | `GLCHR_SERVER_PORT`         | `"0"`         |
| --log-level           | `GLCHR_LOG_LEVEL`           | `"info"`      |
| --log-format          | `GLCHR_LOG_FORMAT`          | `"text"`      |

## Getting Started

To get up and running, have a look at the [demo](./examples/demo/README.md)
example which contains a `docker compose` setup to provision a ClickHouse server
and a Grafana instance that includes predefined dashboards.

## License

This project is licensed under the [MIT License](./LICENSE).
