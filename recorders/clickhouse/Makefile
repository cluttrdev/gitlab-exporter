
PKG ?= .
APP := gitlab-exporter-clickhouse-recorder

ifneq "${VERBOSE}" "1"
.SILENT:
endif

.DEFAULT_GOAL := help

.PHONY: help
help: ## Display this help page
	grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[33m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build:  ## Create application binary
	if [ -z "${os}" ]; then goos=$$(go env GOOS); else goos="${os}"; fi; \
	if [ -z "${arch}" ]; then goarch=$$(go env GOARCH); else goarch="${arch}"; fi; \
	if [ -z "${output}" ]; then output="bin/${APP}"; else output="${output}"; fi; \
	export version=$$(make --no-print-directory version); \
	GOOS="$${goos}" GOARCH="$${goarch}" \
	go build \
		-ldflags "-s -w -X 'main.version=$${version}'" \
		-o "$${output}" \
		${PKG}

.PHONY: build-image
build-image: ## Build container image
	export version=$$(make --no-print-directory version); \
	if [ -n "${tag}" ]; then tag="${tag}"; else tag="$${version/+/-}" ; fi; \
	docker build \
		--file Dockerfile \
		--build-arg VERSION=$${version} \
		--tag "${APP}:$${tag}" \
		.
