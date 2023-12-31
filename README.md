# gitlab-exporter

`gitlab-exporter` can help you build an observability and analytics solution to
gain insights into your CI pipelines by exporting data retrieved from the
[GitLab API][gitlab-api] to a data store backend such as
[ClickHouse][clickhouse].

---

**Note:** This project is in an early development stage, so functionality and
configuration options are limited.

---

<p>
    <img src="./assets/project-overview.webp" />
    <img src="./assets/pipeline-trace.webp" />
</p>

## Installation

To install `gitlab-exporter` you can download a 
[prebuilt binary][prebuilt-binaries] that matches your system, e.g.

```shell
# Download
OS=linux
ARCH=amd64
RELEASE_TAG=$(curl -sSfL https://api.github.com/repos/cluttrdev/gitlab-exporter/releases/latest | jq -r '.tag_name')
curl -sSfL https://github.com/cluttrdev/gitlab-exporter/releases/download/${RELEASE_TAG}/gitlab-exporter_${RELEASE_TAG}_${OS}_${ARCH}.tar.gz
# Install
tar -xf gitlab-exporter_*.tar.gz gitlab-exporter
install gitlab-exporter ~/.local/bin/gitlab-exporter
```

Alternatively, if you have the [Go][go-install] tools installed on your
machine, you can use

```shell
go install github.com/cluttrdev/gitlab-exporter@latest
```

## Usage

`gitlab-exporter` can either run in daemon mode or execute one-off
commands.

### Daemon Mode

To run `gitlab-exporter` in daemon mode use:

```shell
gitlab-exporter --config CONFIG_FILE run
```

This will periodically export data for updated pipelines of the configured projects,
see [Configuration](#configuration) for configuration options.

### Command Mode

`gitlab-exporter` supports a number of commands that can be executed
individually. Use the following to get an overview of available commands:

```shell
gitlab-exporter -h
```

## Configuration

Configuration options can be specified in a config file that is passed to the
application using the `--config` command-line flag.

For an overview of available configuration options and their default values,
see [configs/gitlab-exporter.yaml](./configs/gitlab-exporter.yaml).

Some common options can also be passed as command-line flags and/or environment
variables, with flags taking precedence.

| Flag                  | Environment Variable      | Default Value                 |
| ---                   | ---                       | ---                           |
| --gitlab-api-url      | `GLE_GITLAB_API_URL`      | `"https://gitlab.com/api/v4"` |
| --gitlab-api-token    | `GLE_GITLAB_API_TOKEN`    | **required**                  |
| --clickhouse-host     | `GLE_CLICKHOUSE_HOST`     | `"127.0.0.1"`                 |
| --clickhouse-port     | `GLE_CLICKHOUSE_PORT`     | `"9000"`                      |
| --clickhouse-database | `GLE_CLICKHOUSE_DATABASE` | `"default"`                   |
| --clickhouse-user     | `GLE_CLICKHOUSE_USER`     | `"default"`                   |
| --clickhouse-password | `GLE_CLICKHOUSE_PASSWORD` | `""`                          |

## Development Environment

To test the application during development or to just see what it has to offer,
a [docker-compose](./environments/dev/docker-compose.yaml) setup is provided
that can be used to create a simple environment consisting of a ClickHouse server
and a Grafana instance that includes some predefined dashboards.

To use this, change directory to `environments/dev/` and run:

```shell
docker compose up -d
```

The ClickHouse server will listen on `127.0.0.1:9000` and have the following
database and user created:
```
database: gitlab_ci
user:     gitlab_exporter
password: gitlab_exporter
```
See the
[config.xml](./environments/dev/clickhouse/config.xml),
[users.xml](./environments/dev/clickhouse/users.xml) and
[init-db.sh](./environments/dev/clickhouse/initdb.d/init-db.sh)
files for more details.

Then, create a configuration file or set the necessary environment variables
and run `gitlab-exporter` (either in daemon mode or using one-off commands):
```shell
# create simple config file with the database settings
cat <<EOF > config.yaml
clickhouse:
  database: "gitlab_ci"
  user: "gitlab_exporter"
  password: "gitlab_exporter"
EOF

# set required api token as environment variable (or put it in the config file)
export GLE_GITLAB_API_TOKEN=<your-gitlab-token>

# run the application
gitlab-exporter run --config ./config.yaml --projects <project-ids>
```

You should then be able to login to Grafana on <http://localhost:3000> (using
the default `admin:admin` credentials) and explore the data.

## Acknowledgements

This project was inspired by Maxime Visonneau's
[gitlab-ci-pipeline-exporter][github-gcpe].

The project logo is based on the original [artwork][gopher-artwork] created by
Ashley McNamara.

## License

This project is licensed under the [MIT License](./LICENSE).

[gitlab-api]: https://docs.gitlab.com/ee/api/rest/
[clickhouse]: https://clickhouse.com/
[go-install]: https://go.dev/doc/install
[prebuilt-binaries]: https://github.com/cluttrdev/gitlab-exporter/releases/latest
[github-gcpe]: https://github.com/mvisonneau/gitlab-ci-pipelines-exporter
[gopher-artwork]: https://github.com/ashleymcnamara/gophers
