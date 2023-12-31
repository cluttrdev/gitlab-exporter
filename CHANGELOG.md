# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Support export of job log embedded metrics (experimental)
- HTTP server implementation with health, debug and metrics endpoints
- ClickHouse table deduplication task and subcommand

### Changed

- Provide data store interface to make exporter more flexible

### Removed

- ClickHouse database creation (instead use the one configured to connect with)

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


[Unreleased]: https://github.com/cluttrdev/gitlab-exporter/compare/v0.4.1...HEAD
[0.4.1]: https://github.com/cluttrdev/gitlab-exporter/compare/v0.4.0...v0.4.1
[0.4.0]: https://github.com/cluttrdev/gitlab-exporter/compare/v0.3.1...v0.4.0
[0.3.1]: https://github.com/cluttrdev/gitlab-exporter/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/cluttrdev/gitlab-exporter/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/cluttrdev/gitlab-exporter/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/cluttrdev/gitlab-exporter/releases/tag/v0.1.0

