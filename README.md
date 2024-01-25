# gitlab-exporter

`gitlab-exporter` can help you build an observability and analytics solution to
gain insights into your CI pipelines. 
It provides the [protocol buffer][protobuf] schemas and client implementation 
for a [gRPC] service to record data fetched from the [GitLab API][gitlab-api].
This allows the data to be exported and stored in various storage backends by
implementing the correspondong gRPC server interface.

The following server implementations are available:

  - [gitlab-clickhouse-exporter][gh-glche]

<p>
    <img src="./assets/project-overview.webp" />
    <img src="./assets/pipeline-trace.webp" />
</p>

## Installation

To install `gitlab-exporter` you can download a 
[prebuilt binary][prebuilt-binaries] that matches your system, e.g.

```shell
# download latest release archive
RELEASE_TAG=$(curl -sSfL https://api.github.com/repos/cluttrdev/gitlab-exporter/releases/latest | jq -r '.tag_name')
curl -sSfL -o /tmp/gitlab-exporter.tar.gz \
    https://github.com/cluttrdev/gitlab-exporter/releases/download/${RELEASE_TAG}/gitlab-exporter_${RELEASE_TAG}_linux_amd64.tar.gz
# extract executable binary into install dir (must exist)
INSTALL_DIR=$HOME/.local/bin
tar -C ${INSTALL_DIR} -zxof /tmp/gitlab-exporter.tar.gz gitlab-exporter
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
gitlab-exporter run --config CONFIG_FILE 
```

This will periodically fetch data of the configured projects and send it to
configured gRPC server endpoints.
See [Configuration](#configuration) for configuration options.

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

Some options can also be overriden with command-line flags and/or environment
variables, where flags take precedence.

| Flag               | Environment Variable   | Default Value                 |
| ---                | ---                    | ---                           |
| --gitlab-api-url   | `GLE_GITLAB_API_URL`   | `"https://gitlab.com/api/v4"` |
| --gitlab-api-token | `GLE_GITLAB_API_TOKEN` | **required**                  |

## Acknowledgements

This project was inspired by Maxime Visonneau's
[gitlab-ci-pipeline-exporter][github-gcpe].

The project logo is based on the original [artwork][gopher-artwork] created by
Ashley McNamara.

## License

This project is licensed under the [MIT License](./LICENSE).

[protobuf]: https://protobuf.dev/
[grpc]: https://grpc.io/
[gitlab-api]: https://docs.gitlab.com/ee/api/rest/
[gh-glche]: https://github.com/cluttrdev/gitlab-clickhouse-exporter
[go-install]: https://go.dev/doc/install
[prebuilt-binaries]: https://github.com/cluttrdev/gitlab-exporter/releases/latest
[github-gcpe]: https://github.com/mvisonneau/gitlab-ci-pipelines-exporter
[gopher-artwork]: https://github.com/ashleymcnamara/gophers
