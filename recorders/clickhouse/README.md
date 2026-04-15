# gitlab-exporter-clickhouse-recorder

`gitlab-exporter-clickhouse-recorder` serves a gRPC endpoint that can receive data
from a [gitlab-exporter](https://gitlab.com/gitlab-exporter/gitlab-exporter) and records it
in a [ClickHouse](https://clickhouse.com) database.

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
