# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.10.1] - 2024-10-30

### Changed

- Use pipeline name for trace span name if available

### Fixed

- Remove unused gitlab-api-(url|token) config option
- Parsing graphql namespace gids
- Skip exporting empty trace data

## [0.10.0] - 2024-10-30

### Added

- Internal data types to separate upstream data from protobufs
- GitLab GraphQL API client
- Task controller to handle multitude of projects

### Changed

- Adjust protobuf and grpc types to new graphql data model and internal types

### Removed

- Job types and worker pools

## [0.9.1] - 2024-10-02

### Fixed

- Move project resolution to jobs run group

## [0.9.0] - 2024-10-02

### Added

- Export custom merge request note events

### Changed

- Replace worker pool implementation with alitto/pond

## [0.8.1] - 2024-07-18

### Changed

- Add ids to metric proto message and remove nested job reference

## [0.8.0] - 2024-07-16

### Added

- Provide gRPC server implementation
- Export merge requests data
- Export projects data
- Configure project namespaces for export
- Configure default project export options

### Changed

- Export data in chunks respecting max size limit

### Removed

- Obsolete `catch_up.forced` config option

## [0.7.1] - 2024-02-27

### Added

- Started adding functional tests
- Helm chart podmonitor template

### Changed

- Make grpc client creation more flexible
- Add metadata to pipeline hierarchy export errors

### Fixed

- Helm chart selector labels helper

## [0.7.0] - 2024-02-25

### Added

- gRPC client metrics

### Changed

- Switch to unary grpc calls

## [0.6.3] - 2024-02-20

### Fixed

- Fix testreport exports
- Fix list project pipelines error handling

## [0.6.2] - 2024-02-19

### Fixed

- Command flags env var prefix

## [0.6.1] - 2024-02-19

### Fixed

- Command config flags propagation

## [0.6.0] - 2024-02-19

### Added

- New dedicated `catchup` subcommand
- Helm template to run catchup job

### Changed

- Default to not catching on historical project data when using `run` command
- Use structured logging with `log/slog`
- Restructure protobuf and grpc files and packages

### Fixed

- Support for log embedded metrics

### Removed

- Need for a `controller` structure

## [0.5.1] - 2024-02-04

### Added

- Fetching pipeline test report summaries

### Changed

- Use pipeline test report summary build ids to set test suite/case ids

### Removed

- Unused projects model

## [0.5.0] - 2024-01-26

This is quite a big release that introduces a lot of breaking changes.

### Added

- gRPC service and protocol buffer schemas for client-server model
- Support export of job log embedded metrics (experimental)
- HTTP server implementation with health, debug and metrics endpoints

### Changed

- Project and binary name

### Removed

- Storage backend specific implementations

### Fixed

- Version info subcommand

## [0.4.1] - 2023-10-24

### Changed

- No retry loop in project catch-up worker

### Fixed

- Checking for closed channel when producing in catch-up worker

## [0.4.0] - 2023-10-11

### Added

- Loading configuration from yaml file
- More detailed project configuration options
- Fetching project data from the GitLab API
- Option to force data export during catch-up from pipelines that were not updated

### Changed

- Use workers to handle controller tasks

## [0.3.1] - 2023-09-14

### Added

- Pipeline trace span attributes

### Changed

- Build docker runtime image from scratch

### Fixed

- Add version info to docker build binary

## [0.3.0] - 2023-09-11

### Added

- Default configuration file
- CHANGELOG

### Changed

- Renamed GitLab client `requests-per-second` to `rate-limit` flag

### Removed

- Unused `config.LoadEnv()` function

### Fixed

- Support rate limiting flag in cli and config file
- Printing custom usage function

## [0.2.0] - 2023-09-04

### Added

- Worker pool for task management
- GitLab client rate limiting
- `version` command
- justfile release recipes

### Changed

- Custom usage function
- Init controller and clients only for commands that use them

## [0.1.0] - 2023-08-24

Initial release.

### Added

- Exporting full pipeline hierarchies
- Exporting pipeline test reports
- CLI with daemon mode (`run`) and one-off commands (`fetch`, `export`)
- Dockerfile
- justfile
- LICENSE
- README


<!-- Links -->

[Unreleased]: https://github.com/cluttrdev/gitlab-exporter/compare/v0.8.0...HEAD
[0.8.0]: https://github.com/cluttrdev/gitlab-exporter/compare/v0.7.1...v0.8.0
[0.7.1]: https://github.com/cluttrdev/gitlab-exporter/compare/v0.7.0...v0.7.1
[0.7.0]: https://github.com/cluttrdev/gitlab-exporter/compare/v0.6.3...v0.7.0
[0.6.3]: https://github.com/cluttrdev/gitlab-exporter/compare/v0.6.2...v0.6.3
[0.6.2]: https://github.com/cluttrdev/gitlab-exporter/compare/v0.6.1...v0.6.2
[0.6.1]: https://github.com/cluttrdev/gitlab-exporter/compare/v0.6.0...v0.6.1
[0.6.0]: https://github.com/cluttrdev/gitlab-exporter/compare/v0.5.1...v0.6.0
[0.5.1]: https://github.com/cluttrdev/gitlab-exporter/compare/v0.5.0...v0.5.1
[0.5.0]: https://github.com/cluttrdev/gitlab-exporter/compare/v0.4.1...v0.5.0
[0.4.1]: https://github.com/cluttrdev/gitlab-exporter/compare/v0.4.0...v0.4.1
[0.4.0]: https://github.com/cluttrdev/gitlab-exporter/compare/v0.3.1...v0.4.0
[0.3.1]: https://github.com/cluttrdev/gitlab-exporter/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/cluttrdev/gitlab-exporter/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/cluttrdev/gitlab-exporter/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/cluttrdev/gitlab-exporter/releases/tag/v0.1.0
