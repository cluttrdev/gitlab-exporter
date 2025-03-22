# Changelog

## [0.16.0](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.15.1..v0.16.0)

- b467b7c ci: Add chart release job
- 3b0f6ae chore: Fix typo in example config
- 531f82a fix: Fix linting errors
- 4ea996b chore: Fix clickhouse-exporter link in README
- d9ff257 refactor: Move module url to go.cluttr.dev and repo references to gitlab
- d1a15ff test: Move test files to source file dirs
- f2f3a53 patch!: Add pipeline and job fields ref_path
- 34575eb patch: Add project field default_branch
- 6681e4c refactor: Shorten grpc client error messages
- 35e147f feat: Export cobertura coverage reports
- 961eda2 ci: Upload cobertura coverage report

## [0.15.1](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.15.0..v0.15.1)

- 5bc5ec5 chore(release): v0.15.1
- 266c78f fix: Map failure messages to reasons
- 85c96e9 ci: Update golangci-lint

## [0.15.0](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.14.2..v0.15.0)

- bac980e chore(release): v0.15.0
- 70dea2a feat: Download junit report artifacts via api
- 3034d40 ci: Update go version
- f178815 build: Update default go version in Dockerfile
- 2119961 fix: Update dependencies
- 2c71f86 refactor!: Use job reference in testreports

## [0.14.2](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.14.1..v0.14.2)

- f964506 chore(release): v0.14.2
- d154c45 fix: Add missing testcase status when converting to protobuf message

## [0.14.1](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.14.0..v0.14.1)

- afc419c chore(release): v0.14.1
- 33b66d8 chore: Update gitlab remote references
- ed4a170 feat: Use project feature access levels
- 24c0106 fix: Prevent rest.GitLabDeployment nil pointer dereference for user and environment

## [0.14.0](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.13.2..v0.14.0)

- 0675d78 chore(release): v0.14.0
- d01da14 chore: Update example config file
- 4b767cf feat: Export deployments

## [0.13.2](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.13.1..v0.13.2)

- 8e7aa82 chore(release): v0.13.2
- b66e1b4 refactor!: Use user references in merge request note events
- b890609 patch: Add username and name to user references
- 311d75e docs: Fix changelog

## [0.13.1](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.13.0..v0.13.1)

- 690f199 chore(release): v0.13.1
- 81fc2ff chore: Remove unused rest api functions
- 1910756 chore: Remove junitmetrics tool
- e01a804 fix: Remove oauth2 support
- 2629a43 patch: Move to gitlab.com/gitlab-org/api/client-go
- 35cad6b patch: Update dependencies
- 9ff135b fix: Possible gitlab client blocking when fetching many resources
- 99a9333 fix: Session authed http client locking

## [0.13.0](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.12.0..v0.13.0)

- b501326 chore(release): v0.13.0
- 002f430 refactor(cmd): Consolidate gitlab client creation
- 5d6be44 feat: Add session authed gitlab http client (experimental)
- 4b17dfe fix: OAuth redirect handler path

## [0.12.0](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.11.0..v0.12.0)

- 3081af6 chore(release): v0.12.0
- 116e486 docs: Update readme
- ef359ee docs: Add unreleased changes when generating changelog
- 6f6ff42 feat: Export properties in junit test reports
- 17353b9 build: Set default make goal
- 275c4b9 chore: Update example config
- 92f1b8f feat(cmd): Add fetch artifacts command
- 8a58245 feat(cmd): Add oauth request and refresh commands
- 79bac33 feat: Fetch projects pipelines junit reports
- 044abf2 feat: Fetch graphql job artifacts download paths
- bf03891 feat: Add oauthed gitlab http client
- 44dd4a3 feat: Add gitlab oauth2 package
- 50b64a3 refactor: Add internal http client
- 1887445 feat: Add junitxml test report types
- 536201a refactor!: Adjust protobuf messages to internal types
- 88299f0 fix: Add interrupted job sections
- 9515877 refactor: Adjust test report types
- d45bb8a ci(github): Fix docker image tags and repo urls

## [0.11.0](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.10.3..v0.11.0)

- 7ba86e7 chore(release): v0.11.0
- 8e76142 build: Add make target to generate changelog
- 4120a9f build: Replace justfile and scripts with Makefile
- 97f9eec ci: Push release image to both GitLab and GitHub registries
- f8d1a64 ci(github): Add release sync workflow
- 5829e71 ci: Update gitlab-ci.yml, add release jobs
- 242b14f revert: Add ci image build job
- fd19805 build: Add Makefile
- a40d159 ci: Add ci image build job
- 70428f3 ci: Fix go version
- 4bd3a9e feat: Periodically resolve projects
- e0117ce refactor: Make task controller projects settings map thread safe
- 7453403 refactor: Call export methods concurrently for endpoint clients and data batches

## [0.10.3](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.10.2..v0.10.3)

- ebdc3ca chore(release): v0.10.3
- 6e9045d refactor: Remove custom protobuf time conversion functions
- 6519dd8 fix: Try fetching job log data if runner id is empty
- c00e7fd fix: Job protobuf queued duration value
- 6dabc7c fix: Catchup command project resolution

## [0.10.2](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.10.1..v0.10.2)

- 97c6bcc chore(release): v0.10.2
- dd8cc4c fix: Regenerate protobuf code

## [0.10.1](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.10.0..v0.10.1)

- 6841b39 chore(release): v0.10.1
- a01cb80 refactor: Use pipeline name for trace span name if available
- 58b5a9e fix: Do not export empty trace data
- 1bcaa05 fix: Parsing namespace gids
- 9f9b0e0 chore: Add graphql generation recipe to justfile
- 0b3fae6 fix: Remove gitlab-api-(url|token) config usage

## [0.10.0](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.9.1..v0.10.0)

- d02d72b chore(release): v0.10.0
- b326797 refactor!: Adjust protobuf and grpc stuff
- bf91a56 chore: Update dependencies
- 3f0123d chore: Adjust commands
- d958d22 refactor: Replace jobs with tasks controller
- 8fc1968 feat: Compose rest and graphql gitlab api clients
- 881e0f2 feat: Add gitlab graphql client
- a7824b5 feat: Add internal types

## [0.9.1](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.9.0..v0.9.1)

- a663eae chore(release): v0.9.1
- dcecdb2 fix(cmd): Move namespace project resolution into jobs run group
- 366e18a fix: Listing (user) namespace projects and options
- 0b43bed test: Add config test for export defaults
- ef4ef33 chore: Update helm chart

## [0.9.0](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.8.1..v0.9.0)

- 1b09bd7 chore(release): v0.9.0
- 819e027 feat: Remove custom worker pool package
- 3880953 feat: Serve worker pool prometheus metrics
- a6239ea patch: Improve listing GitLab API resources
- 9480e4d feat: Add merge request note events
- 053d870 refactor: Move proto types conversion to separate package

## [0.8.1](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.8.0..v0.8.1)

- b977639 chore(release): v0.8.1
- f9c9fc2 patch: Add ids to metric proto message and remove job reference message

## [0.8.0](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.7.1..v0.8.0)

- 44bf374 chore(release): v0.8.0
- d50f297 patch: Ensure data is exported in chunks with max size
- bf7c7b2 test: Fix config tests
- 1302c71 patch: Add merge requests block to project export config
- 58f3b84 feat: Exclude namespace projects from being exported
- 94f85e6 patch: Always export project on first iteration
- 30df05e patch: Expose underlying gitlab.Client for internal use
- 4e1582c feat: Export merge requests
- aa052a4 feat: Provide grpc server implementation
- fc64fd8 feat: Add new protos
- 2a9a6b6 Update go version and dependencies
- 06ef8a5 Replace golang.org/x/exp/slices with stdlib package
- e38d47c Add support for project namespaces configuration
- 44fff86 Add support for exporting project data
- 476a55e Remove internal/models package
- 17a4698 Export testcases in chunks
- ff46391 Refactor exporting pipeline hierarchy testreports
- 712c92f Start adding integration tests
- 9da9735 Add gitlab-exporter recorder mock
- 294b529 Add gitlab server mock to serve json testdata
- 72b378e Update go version and dependencies
- 44eb406 Update README.md
- 94749c3 Add custom sections to gitlab-ci.yml
- 7b5d43c Remove obsolete config.catch_up.forced option
- 35d27d3 Add project defaults config options
- 55b1857 Fix helm chart ports config again

## [0.7.1](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.7.0..v0.7.1)

- 6d6364e Release v0.7.1
- eaf813a Fix helm selector labels helper
- 5cc58ea Fix podmonitor helm template
- a694cc2 Add metadata to pipeline hierarchy export errors
- d7eb91b Fix helm http monitoring and service config
- cf65803 Add podmonitor template to helm chart
- 4f91e98 Start adding tests
- fbd0bf9 Make grpc client creation more flexible

## [0.7.0](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.6.3..v0.7.0)

- e5e8943 Release v0.7.0
- 687536e Fix execution when no projects are configured
- 2c04805 Add http healthchecks again
- 246d0b2 Add log output for signal handler
- 51028d5 Fix example config file
- 8b7f29e Remove unused internal/server package
- 10bad08 Refactor run and catchup commands
- f563081 Add grpc client metrics suppport
- 4e20117 Ammend http config and fix unmarshalling
- 68778f6 Rename config 'server' block to 'http'
- ed71393 Switch to unary grpc calls

## [0.6.3](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.6.2..v0.6.3)

- 4f6df05 Release v0.6.3
- 7d00d04 Fix testreport exports
- f853c10 Fix list project pipelines error handling

## [0.6.2](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.6.1..v0.6.2)

- 33bf746 Release v0.6.2
- c18c9a4 Fix command flags env var prefix

## [0.6.1](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.6.0..v0.6.1)

- 14dc288 Release v0.6.1
- 878d557 Fix command config flags propagation
- 4edf8c0 chore: Support helm chart podLabels value
- 070bcac Fix .dockerignore

## [0.6.0](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.5.1..v0.6.0)

- 725ecc3 Release v0.6.0
- a23dc7a Update dependencies
- 418de01 Add helm template to run catchup job
- 9ae88fb Update CHANGELOG.md and example config
- 9138e1f Restructure protobuf and grpc files and packages
- 7f50ade Rename LogEmbeddedMetric to just Metric
- 14441bb Skip exporting empty data slices
- a195a88 Add --catchup flag to run command
- 3f89585 Add dedicated 'catchup' command
- 5435973 Remove the need for a 'controller'
- ac50a8e Update README.md
- 1e11979 Use structured logging
- 2e6498c Fix support for log embedded metrics export
- 13e79b0 Refactor data recording methods
- 71d3e92 Add template support for helm env and config values
- f10fc73 Update helm chart versions

## [0.5.1](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.5.0..v0.5.1)

- 4c1550d Release v0.5.1
- 9edd1b5 Remove unused projects model
- cf7b06a go mod tidy
- f864a08 Use test summary build ids to set test suite and case ids
- 3ab95af Add pipeline testreport summary method
- eafaec2 Update helm chart appVersion and default registry

## [0.5.0](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.4.1..v0.5.0)

- 7c8e421 Release v0.5.0
- 7c655b2 Fix version info subcommand
- 69cc1f0 Fix go version used in .gitlab-ci.yml
- 5e0c3f2 Fix testreport cases pb conversion
- 5099d6c Add justfile recipe to compile protobufs
- eb37e12 Add opentelemetry-proto submodule
- 9a49cdc Update README.md and CHANGELOG.md
- 0b01c7b Re-add docker justfile recipes and move Dockerfile to root dir
- 5560891 Update README.md
- a00463a Move environments/dev setup to gitlab-clickhouse-exporter repo
- 400a91d Update justfile and add scripts
- b8e57cf Replace peterbourgon/ff/v3/ffcli with cluttrdev/cli
- 1937560 Remove unused internal/util package
- a75cd24 Move internal/worker to pkg/
- 16ce430 Move internal/healthz to pkg/
- a23f029 Move pkg/server to internal/
- ce1346a Move pkg/gitlab to internal/
- 5b6e327 Move pkg/controller to internal/
- 7d13403 Move pkg/config to internal/
- 8f8bec9 BREAKING: Replace pkg/models package with protobufs
- 9a073f4 BREAKING: Remove ClickHouse/Datastore functionality in favor of gRPC exporter
- cf31800 Add exporter type and endpoint config
- a7bb87a Fix traces grpc export
- 37a4e20 Add protobuf schemas for grpc exporter service
- 85ff360 Change project logo
- 0b3adf1 Update changelog
- 91360cb Remove obsolete default clickhouse db name constant
- 0e8fc41 Let golangci-lint be more verbose
- 1bf1a3a Merge branch 'log_embedded_metrics' into 'main'
- ba94c3c Add job log parsing to public interface
- 47f61db Fix parsing job log sections
- ca4277c Add build cache dir to ci cache
- aa3f9b3 Set GOMAXPROCS in ci
- 607c81f Add some log embedded metrics to ci jobs
- 1f798d9 Add tool to turn junit test report results into log embedded metrics
- 4f39e5e Set default ci tags and non-root writable cache locations
- 9c1a4c7 Update CHANGELOG.md
- 4a21cba Enable job embedded metrics export
- 8066536 Fix expfmt text parse test constant quotes
- ed39f8c Add datastore job metrics insert interface and clickhouse implementation
- f7b15d6 Add clickhouse job metrics ddl and dml
- cf31730 Add job embedded metrics model
- cb35f3e Add fetching and parsing job logs for sections and metrics
- 8655f94 Refactor job section parsing
- 42fb459 Add embedded metrics text parser
- 3c53000 Fix gitlab-ci config file name
- 9743fab Fix server block in example config
- 9a30cc8 Extend README dev environment section
- a622a7e Add database initialization to dev environment setup
- 3aff75f Remove clickhouse db creation on datastore init
- d8e070a Replace hardcoded clickhouse database with configured one
- e860ecd BREAKING: Rename project
- cf8ae29 Rename gitlabclient package
- 57736e3 Move worker package to internal
- 7488884 Replace clickhouse.Client usage with new DataStore interface
- 6b635de Add ClickHouseDataStore implementing the DataStore interface
- 0f823a0 Add DataStore interface
- 79b1fec Add models.Trace type declaration
- 23c2584 Rename clickhouseclient package
- 84af351 Add Helm chart
- 1225be3 Split server address config into host and port
- 7180ae5 Update dev environment dashboard provisioning
- 2490395 Remove vet job since golangci-lint includes it
- 2224c54 Cache gopath in ci
- 6c38aee Fix deduplication task tests
- bb84d2b Remove unused parseID and pathEsacpe
- 9f0f202 Add just lint recipe and fix some linting errors
- c4109e3 Update required go version to 1.20
- fe73f90 Add .gitlab/gitlab-ci.yaml
- 9d8163d Remove flag types from usage output
- f737313 Update README
- 73b42cf Add logo files
- aa9a562 Add some readiness checks
- 2eb733a Update changelog
- 1bbc380 Change pseudo-version timestamp command
- b4cb56c Fix deduplication queries
- f2ba5b2 Add deduplication task and subcommand
- d1405b2 Remove need for default checkOK in healthz handler
- 81c710c Add debug and preliminary health endpoints
- 68561f6 Add simple http server with minimal prometheus metrics

## [0.4.1](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.4.0..v0.4.1)

- 86e8e26 Release v0.4.1
- db2a5c6 Remove retry loop in catch-up worker and check for closed channel when producing

## [0.4.0](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.3.1..v0.4.0)

- c0a9843 Release v0.4.0
- f840fd5 Add option to force data export during catch-up
- 18c2b1b Add worker management to controller and some smaller improvements
- 725af26 Separate worker from controller and detach some controller methods
- 0167933 Safeguard gitlab and clickhouse clients with rwmutexes
- ec27b0d Update changelog
- 3ec3d73 Group project export config options
- 7c0a207 Fix loading default project settings
- 97937ca Add testreport and traces config options to default/exmaple config file
- c569d3a Use default project settings for run command projects
- b115411 Fix embedded project settings unmarshaling
- e6fd8b8 Support Project.TestReports and Project.Traces config options
- a5a7d3e Add Project.TestReports and Project.Traces config options
- 9e35a9a Support config.Project.Sections.Enabled option
- bf76578 Split catch-up worker run into produce and process
- 0e8bd65 Fix logging config
- 54553a0 Change commands to use workers
- 3a5bb47 Add workers
- 0a91433 Add project model
- f98e92b Restrict project.id config option to integers (for now)
- c3bc990 Add more detailed projects config options
- 89d5e66 Use custom logic to load config from file
- ac13fbd Add loading configuration from yaml and first tests

## [0.3.1](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.3.0..v0.3.1)

- 8dcff37 Release v0.3.1
- 0dd72e1 Add docker run just recipe
- 1df7df8 Fix adding binary version info to docker build
- dadde55 Update dev environment, adding clickhouse config files
- d98a2c9 Add .gitignore
- 135b7d8 Add .dockerignore and update Dockerfile
- 8ad8346 Add some attributes to trace spans

## [0.3.0](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.2.0..v0.3.0)

- c33b90b Release v0.3.0
- 5f81568 Add CHANGELOG.md
- beebab4 Fix some usage function stuff
- ce4ff40 Fix printing flag.ErrHelp after usage func
- 226e5ac Remove unused config.LoadEnv() function
- 4f6ceb4 Add example/default config file
- 9632a12 Change requests-per-second to rate-limit
- 901c09d Add dirty info to version and fix docker-push recipe dependencies
- bfae306 Tidy up
- 191567c Add support for gitlab-client-requests-per-second flag
- b6ce63d Sanitize keys read from config file

## [0.2.0](https://gitlab.com/akun73/gitlab-exporter/-/compare/v0.1.0..v0.2.0)

- c2a02f4 Add version command
- e1d3512 Init controller only in commands that use it
- d64b763 Implement adjusted default usage func
- 00c2c5a Add internal util package
- 9a7a55d Add docker build and push just recipes
- 22b333d Add github release script and update justfile
- ef1e667 Add worker pool and rate limiter

## [0.1.0](https://gitlab.com/akun73/gitlab-exporter/-/commits/v0.1.0)

- 546cfbb Add justfile
- 5cc1a21 Update README and tidy up
- bf44b5f Format code
- 69b46ab Change clickhouse port config option type to string
- aa76804 Enable configuration via config file
- 6928212 Change flag error handling strategy
- b9a8f9d List cmdline flags in README
- e650d01 Update README to include dev environment section
- d5bcf86 Adjust methods for concurrent api calls
- 8821ca9 Add build dockerfile
- 93d1538 Add first screenshots to the README.md
- 494332b Update project-overview dashboard
- b2b6f27 Update project overview dashboard
- e07647f Add export subcommands
- 715a2e4 Change test result struct id fields data type to int64
- 806a06e Change test result table id column data type to Int64
- ac6c983 Remove redundant trace insertion
- 6733dce Tidy up go module
- 73299e0 Update fetch testreport cmd function name
- 2ccd25a Split up nested test report table
- 1f1f423 Add test reports export when exporting pipeline
- 3ff9a04 Add method to get pipeline hierarchy test reports
- faf6759 Add id and reference fields to test report types
- 51dcb6e Format code
- 359155a Extend cli with fetch subcommands
- 781dbbc Enable json marshalling for pipeline hierarchy struct
- ae863f3 Add proper cli
- 027d6d2 Add types and methods to get pipeline test reports
- 0c1aa14 Update README to reflect changes in cli
- 32c3b3e Reformat code
- 7906ed6 Add daemon run
- 5430469 Add gitlab ListProjectPipelineOption, pagination and fix argument data types
- 4739949 Account for null downstream pipelines in bridges
- ec358d0 Explicitly log to stdout
- b80b10a Add controller type and basic console application structure
- ed1871f Fix minor README issues
- 7dcb91c Add license
- ad19b6f Initial commit, proof of concept

