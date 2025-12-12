# Gitlab Exporter Clickhouse Recorder Demo

This demo provides a `docker compose` setup that provisions the following
services:

- a `gitlab-exporter` that exports data fetched from the GitLab API
- a `gitlab-exporter-clickhouse-recorder` that receives the exported data and stores it in a ClickHouse database
- a `ClickHouse` server that acts as a storage backend
- a `Prometheus` server to scrape metrics from `gitlab-exporter` and `gitlab-exporter-clickhouse-recorder`
- a `Grafana` instance to visualize data stored in the ClickHouse server as well as metrics stored in the Prometheus TSDB

To get started, run:
```shell
docker compose up -d
```

The ClickHouse server will listen on `127.0.0.1:9000` and have the following
database and user credentials created:
```
database: gitlab_ci
user:     glchr
password: glchr
```
See the
[config.xml](./clickhouse/config.xml),
[users.xml](./clickhouse/users.xml) and
[init-db.sh](./clickhouse/initdb.d/init-db.sh)
files for more details.

<!-- Links -->
[gh-gitlab-exporter]: https://github.com/cluttrdev/gitlab-exporter

