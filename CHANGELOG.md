# Changelog

## [Unreleased](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.15.0..HEAD)

- 287b3fd ci: Remove github release sync
- 4f631b8 chore: Update default image repository in helm chart
- 728ecde docs: Update releases url in README, Makefile and CHANGELOG.md

## [0.15.0](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.14.0..v0.15.0)

- 5e87bf4 chore(release): v0.15.0
- 7225aec ci: Use default.tags
- 3d8d9c8 ci: Use hook to embed job properties
- 0a49fd3 ci: Switch to buildkit for container image builds
- aa8f535 test: Fix clickhouse testcontainer version
- b6b1da8 feat: Record job exit_code
- b361c03 fix(deps): Update dependencies

## [0.14.0](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.13.2..v0.14.0)

- 7cbbdad chore(release): v0.14.0
- ab14a1f feat: Record issues
- e85bb88 build: Set tag for build-image target
- 4af0a72 fix: Update dependencies

## [0.13.2](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.13.1..v0.13.2)

- 9ad70bf chore(release): v0.13.2
- 1c28a81 fix: Version number of add_pipelines_downstream migrations

## [0.13.1](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.13.0..v0.13.1)

- 19c7523 chore(release): v0.13.1
- 3a08e83 ci(github): Link to gitlab release assets instead of re-uploading to github
- f700c76 patch: Record pipeline downstream references
- 624dd4c fix: Update dependencies
- 89ecf9f ci: Fix release artifacts builds
- 6222dca build: Fix build target os/arch usage
- 7e48599 ci: Print some job properties
- 1df8d70 fix(db): Add missing pipelines name column

## [0.13.0](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.12.0..v0.13.0)

- 9bb44ac chore(release): v0.13.0
- b4f7cad Merge branch 'job-properties' into 'main'
- f335909 chore: Update demo example
- 7e93bd3 chore: Update README
- 3af22cd fix: Update dependencies
- dae0cfa style: Format code
- 48554e4 feat: Record job log properties
- f238fa7 ci: Fix release-chart job
- 0a72230 ci: Specify chart release job image tag

## [0.12.0](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.11.2..v0.12.0)

- c3c1ebf chore(release): v0.12.0
- 457abd2 chore: Fix linting errors
- b00b196 ci: Set release name in github release script
- 00aa6a9 ci: Add chart release job
- d39df91 chore: Fix README
- 8321fd4 path: Add pipelines and jbos column ref_path
- 33617f7 path: Add projects column default_branch
- fc2fb84 feat: Record coverage report data
- 1733925 fix: Update dependencies
- 2954ed9 ci: Push to docker hub
- c551d86 docs: Move repo references to gitlab
- cc990a3 refactor: Move to module vanity url
- eb51d61 chore: Fix db schema cli
- 69ace73 chore: Fix failure_reason mutation snippet
- e02f5de chore: Add sql snippet to map job failure reasons
- 55d1eb4 chore: Add cli to generate schema from migrations

## [0.11.2](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.11.1..v0.11.2)

- c1b07f6 chore(release): v0.11.2
- 73f925c fix: InsertDeployments log message
- 447ebc5 fix: Nullify latest error in health check

## [0.11.1](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.11.0..v0.11.1)

- d115c2f chore(release): v0.11.1
- 562f2e9 patch: Add job_id to test entity tables
- 231b708 ci: Update golangci-lint
- 531852b ci: Update go version
- dec0038 build: Update Dockerfile go version
- 401363a fix: Update dependencies

## [0.11.0](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.10.1..v0.11.0)

- 819b9fa chore(release): v0.11.0
- a24833a ci: Use go version 1.23
- dc815f5 chore: Update gitlab remote references
- 5b5bd8d build: Bump container image go version
- 3fb4a91 feat: Record deployments
- b622558 fix: Update dependencies

## [0.10.1](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.10.0..v0.10.1)

- 912f9c8 chore(release): v0.10.1
- b2e8e35 patch: Add user reference names in mr tables
- 0418208 fix: Update dependencies
- dfc1ae0 chore: Allow disabling migrations job in helm chart

## [0.10.0](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.9.0..v0.10.0)

- 175e5f1 chore(release): v0.10.0
- daccd6a ci: Skip integration tests
- 0aaeb93 chore: Update dependencies
- 4698250 refactor: Use table model structs for batch insertion
- 85ec729 fix: Fix some down migrations
- df75a58 refactor: Alter table mergerequests
- b3dbbb3 refactor: Alter table metrics
- ff09f89 refactor: Alter tables testreports, -suites, -cases
- a057b8d refactor: Alter table sections
- 09dda28 refactor: Alter table jobs
- 68d7c1d refactor: Alter table pipelines
- d19d524 refactor: Alter table projects
- e296702 chore: Remove unused cache struct

## [0.9.0](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.8.3..v0.9.0)

- 08bdc6d chore(release): v0.9.0
- 2cd8e8f ci: Disable/skip integrations tests
- d18adc6 ci: Fix dind
- 755d0a0 ci: Fix tests job
- da7ca54 fix: Update dependencies
- c463374 ci: Set up CI workflows
- c8d8667 build: Replace justfile/scripts with Makefile
- adf5339 feat: Adjust to gitlab-exporter v0.12.0 proto changes
- f886c1d chore: Update dependencies
- 8b9199c feat: Allow setting clickhouse client max concurrent queries

## [0.8.3](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.8.2..v0.8.3)

- 1218609 chore(release): v0.8.3
- ec29c05 fix: Fix trace view sql
- 8c72134 chore(deps): Update gitlab-exporter dependency

## [0.8.2](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.8.1..v0.8.2)

- 00f2812 chore(release): v0.8.2
- a053aeb fix: Update gitlab-exporter to v0.10.2

## [0.8.1](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.8.0..v0.8.1)

- e5ac371 chore(release): v0.8.1
- 4fac122 chore: Update demo
- 97e44ce fix: Do not try to insert empty data
- 09eebaf chore(db): Add links to trace view
- 88b512b fix: Adjust to gitlab-exporter v0.10.0 changes
- ea64c57 chore(deps): Update dependencies

## [0.8.0](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.7.1..v0.8.0)

- 46f447a chore(release): v0.8.0
- b6cabb7 chore: Update grafana dashboards
- 2002a83 chore: Update demo config
- 6c04441 feat: Record merge request note events
- dc073d1 chore: Update dependencies

## [0.7.1](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.7.0..v0.7.1)

- 86d5437 chore(release): v0.7.1
- 8820b5c patch: Optimize deduplicating insert queries
- 59bd1df chore: Add rbac role and binding to allow init container query job
- 7e82c81 chore: Fix job template
- 174b003 chore: Apply security context to init container
- d357d38 chore: Add init container to helm chart deployment template
- 7397491 chore: Add migrate job template to helm chart

## [0.7.0](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.6.2..v0.7.0)

- 0d50af5 chore(release): v0.7.0
- efcf7f3 patch: Use async inserts
- ed2aa7b patch: Use deduplicating insert queries instead of client-side caching
- 77da431 patch: Adjust to changes in metrics proto message
- 95e9ec5 chore: Update gitlab-exporter dependency to v0.8.1
- 8419e55 patch: Remove entity id cache
- 22816e4 chore: Update dependencies
- d54149c chore: Update grafana dashboards
- 82e6f60 chore: Add migration container to demo example
- d4a5395 feat: Check database schema version on run
- 156f43b feat: Record projects
- e043fdc feat: Record merge requests
- d150029 feat: Use migrations instead of ddl in code
- 0b05aab refactor: Remove whitespace when using query parameters
- 613f9fd refactor: Use gitlab-exporter/grpc/server
- 2a90177 chore: Update go version and dependencies

## [0.6.2](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.6.1..v0.6.2)

- 8366b57 Release v0.6.2
- 190aa11 chore: Update demo with readonly user
- 009d1e5 chore: Update dependencies
- 5d4a88f test: Localize integration test setup
- 960281e refactor: Reduce memory required to cache entity ids
- 87f7f64 chore: Fix helm chart ports config again
- 3e86c7b chore: Update helm chart version
- 3dfbddb chore: Fix helm chart http service and monitoring config
- 35b8a05 chore: Add helm chart support for podLabels values
- 3bbaf9f chore: Fix helm chart selector labels helper template
- 84d3779 chore: Fix docker compose glchr image

## [0.6.1](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.6.0..v0.6.1)

- 189358d Release v0.6.1
- aa71003 chore: Extend quickstart example to full demo
- 02b8508 test: Started adding integration tests
- 8a5a4cf refactor: Make clickhouse client creation more flexible
- 128b588 fix: Catch jobs without pipeline
- 474465d refactor: Improve table creation functions

## [0.6.0](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.5.3..v0.6.0)

- 1ade409 Release v0.6.0
- 594d73d refactor: Add some log output in run command
- 628d584 refactor: Switch to unary grpc calls
- ba16e4a chore: Update dependencies
- bd4f4aa refactor: Remove http probes pkg again
- 71eda46 feat: Add grpc metrics
- e0863bb chore: Update quickstart example, adding prometheus
- 4c6a95c chore: Fix helm service monitor template port
- ed9b518 chore: Fix helm chart service monitor typo

## [0.5.3](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.5.2..v0.5.3)

- 5ceea0d Release v0.5.3
- 759c03a refactor: Add some debug output
- 1490691 fix: Not waiting on retry ticker after stopped

## [0.5.2](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.5.1..v0.5.2)

- 418c5a0 Release v0.5.2
- 2b41766 chore: Add helm chart service monitor template
- 42b960a feat: Add debug flag
- 49d658e feat: Add http probes again
- 0dccec1 fix: ClickHouse entity id cache update allocations
- 551fc46 fix: Improve config heap escapes
- 629e207 fix: Improve retry allocations

## [0.5.1](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.5.0..v0.5.1)

- 3218669 Release v0.5.1
- 5af58ba fix: Adjust config env var prefix
- c9faeaf chore: Allow name overrides in helm chart
- 377a04e chore: Fix helm chart version
- 8f5cc22 chore: Rename helm chart directory

## [0.5.0](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.4.2..v0.5.0)

- 25547c8 Release v0.5.0
- 3a8dce9 BREAKING CHANGE: Rename project
- d9a0d5e chore: Update quickstart example

## [0.4.2](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.4.1..v0.4.2)

- 1acff02 Release v0.4.2
- 0c66605 fix: Fix trace spans insertion cache update
- 8ae58b6 fix: Fix log embedded metrics cache update

## [0.4.1](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.4.0..v0.4.1)

- f838ba7 Release v0.4.1
- 2624351 fix: Rename RecordLogEmbeddedMetrics to RecordMetrics
- f2529c4 chore: Update helm chart versions

## [0.4.0](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.3.1..v0.4.0)

- d9d6d91 Release v0.4.0
- d6b7c6f chore: Update gitlab-exporter to v0.6.0
- 4dc0cee chore: Update dependencies
- c71dfaf chore: Update helm chart versions

## [0.3.1](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.3.0..v0.3.1)

- 673fabf Release v0.3.1
- 873f3c5 fix: Convert labels when inserting log embedded metrics
- ed453a6 fix: Use constants for table identifiers in insert methods
- 169441c refactor: Check readiness every 3s
- ebdaa73 fix: Log readiness check failures as errors
- a053d19 fix: Set initial serving status to UNKNOWN

## [0.3.0](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.2.1..v0.3.0)

- 3e123e3 Release v0.3.0
- 297cfe9 refactor: Use gRPC health checks instead of HTTP probes

## [0.2.1](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.2.0..v0.2.1)

- 61742b2 Release v0.2.1
- 3690639 refactor: Improve readiness and retry logic in run command
- aca115b chore(helm): Add template support for env and config values
- f80365d fix: Handle both grpc and http ports in helm chart

## [0.2.0](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/compare/v0.1.0..v0.2.0)

- d7131d9 Release v0.2.0
- fcbf0ee feat: Add http probes server
- 0395b73 chore: Update README and CHANGELOG
- f0f49c8 feat: Add structured logging
- f526a49 chore: Ignore grafana dashboard files in quickstart provisioning dir
- c6b5621 refactor: Add column arg to entity ids query function
- 3a67da6 refactor: Use constants and query params for table names
- f2e2ae7 refactor: Add generic record function
- e4c3cec refactor: Change argument order in insert functions
- 9be5083 refactor: Return number of inserted entities
- 0ddd60d feat: Cache inserted entity ids
- 94f06dc feat: Add dql to select entity ids
- 1698ebb refactor: Remove unused latest pipeline update method
- 6ab6479 chore: Add deployment helm chart

## [0.1.0](https://gitlab.com/gitlab-exporter/gitlab-exporter-clickhouse-recorder/-/commits/v0.1.0)

- 91a8680 Release v0.1.0
- 13931f0 build: Fix gitlab-exporter reference in go.mod
- f0e4c97 feat: Add version info subcommand
- 19c20c5 chore: Fix run command in README.md
- 30a4caa build: Add docker build stuff
- d3b067b chore: Add CHANGELOG.md
- 9cbbb89 chore: Update README.md and add quickstart example
- 5d1a1b7 feat: Add config file option
- e9681e8 build: Add justfile and scripts for workflow tasks
- a092919 fix: Adjust testreports ddl
- 079bd8a chore: Add license
- 630aaab Initial commit, proof-of-concept

