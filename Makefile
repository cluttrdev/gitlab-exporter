
PKG ?= .
APP := gitlab-exporter

.PHONY: fmt
fmt: ## Format source code
	go fmt ${PKG}/...

.PHONY: lint
lint: ## Run set of static code analysis tools
	golangci-lint run ${PKG}/...

.PHONY: vet
vet: ## Examine code for suspicious constructs
	go vet ${PKG}/...

.PHONY: graphql
graphql: ## Generate the GitLab GraphQL API client code
	genqlient internal/gitlab/graphql/genqlient.yaml

.PHONY: protobuf
protobuf:  ## Generate Protocol Buffer and gRPC code
	protoc \
		-I protos/ \
		-I protos/vendor/opentelemetry-proto \
		--go_out=. --go_opt=module=github.com/cluttrdev/gitlab-exporter \
		--go-grpc_out=. --go-grpc_opt=module=github.com/cluttrdev/gitlab-exporter \
		protos/gitlabexporter/protobuf/*.proto protos/gitlabexporter/protobuf/service/*.proto

.PHONY: build
build:  ## Create application binary
	export output="bin/${APP}"; \
	if [ -n "${output}" ]; then export output="${output}"; fi; \
	export version=$$(make --no-print-directory version); \
	go build \
		-ldflags "-s -w -X 'main.version=$${version}'" \
		-o "$${output}" \
		${PKG}

.PHONY: build-image
build-image: ## Build container image
	export version=$$(make --no-print-directory version); \
	docker build \
		--file Dockerfile \
		--build-arg VERSION=$${version} \
		--tag "${APP}:$${version/+/-}" \
		.

.PHONY: test
test: ## Run tests
	go test ${PKG}/...

.PHONY: changes
changes: ## Get commits since last release
	from=$$(git describe --tags --abbrev=0 2>/dev/null); \
	if [ -n "${from}" ]; then from="${from}"; fi; \
	to=HEAD; \
	if [ -n "${to}" ]; then to="${to}"; fi; \
	git log --oneline --no-decorate $${from}..$${to}

.PHONY: version
version: ## Generate version from git tag and commit information
	git describe --exact-match 2>/dev/null || echo $$(git describe --tags --abbrev=0)-dev.$$(git rev-list --count $$(git describe --tags --abbrev=0)..HEAD)+$$(git rev-parse --short=8 HEAD)

.PHONY: help
help: ## Display this help page
	grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[33m%-30s\033[0m %s\n", $$1, $$2}'

ifneq "${VERBOSE}" "1"
.SILENT:
endif
