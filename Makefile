REPO_ROOT := $$(git rev-parse --show-toplevel)

.DEFAULT_GOAL := help

.PHONY: help
help: ## Display this help page
	grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[33m%-30s\033[0m %s\n", $$1, $$2}'

ifneq "${VERBOSE}" "1"
.SILENT:
endif

.PHONY: tidy
tidy: ## Run go mod tidy on specified module or all modules
ifdef MOD
	go mod tidy -C ${MOD}
else
	find . -type f -name go.mod -exec sh -c 'go mod tidy -C $$(dirname {})' \;
endif

.PHONY: fmt
fmt: ## Run go fmt on specified module or all modules
ifdef MOD
	go fmt -C ${MOD} ./...
else
	find . -type f -name go.mod -exec sh -c 'go fmt -C $$(dirname {}) ./...' \;
endif

.PHONY: vet
vet: ## Run go vet on specified module or all modules
ifdef MOD
	go vet -C ${MOD} ./...
else
	find . -type f -name go.mod -exec sh -c 'go vet -C $$(dirname {}) ./...' \;
endif

.PHONY: lint
lint: ## Run linter on specified module or all modules
ifdef MOD
	cd ${MOD} && golangci-lint run ./...
else
	find . -type f -name go.mod -exec sh -c 'cd $$(dirname {}) && golangci-lint run ./...' \;
endif

.PHONY: test
test: ## Run tests on specified module or all modules
ifdef MOD
	go test -C ${MOD} ./...
else
	find . -type f -name go.mod -exec sh -c 'go test -C $$(dirname {}) ./...' \;
endif

.PHONY: build
build:  ## Create application binary
	if [ -z "${MOD}" ]; then echo "Please specify module path of which to build package!"; exit 1; fi; \
	if [ -z "${pkg}" ]; then pkg='.'; else pkg="${pkg}"; fi; \
	if [ -z "${os}" ]; then goos=$$(go env GOOS); else goos="${os}"; fi; \
	if [ -z "${arch}" ]; then goarch=$$(go env GOARCH); else goarch="${arch}"; fi; \
	if [ -z "${output}" ]; then output="${REPO_ROOT}/bin/${APP}"; else output="${output}"; fi; \
	export version=$$(make --no-print-directory version); \
	GOOS="$${goos}" GOARCH="$${goarch}" \
	go build \
		-C ${MOD} \
		-ldflags "-s -w -X 'main.version=$${version}'" \
		-o "$${output}" \
		${pkg}

build-exporter:
	$(MAKE) --no-print-directory build MOD="${REPO_ROOT}/cmd/gitlab-exporter"
build-clickhouse-recorder:
	$(MAKE) --no-print-directory build MOD="${REPO_ROOT}/cmd/gitlab-exporter-clickhouse-recorder"

.PHONY: changes
changes: ## Get commits since last release
	to=HEAD; \
	if [ -n "${to}" ]; then to="${to}"; fi; \
	from=$$(git describe --tags --abbrev=0 "$${to}^" 2>/dev/null); \
	if [ -n "${from}" ]; then from="${from}"; fi; \
	if [ -n "$${from}" ]; then \
		git log --oneline --no-decorate $${from}..$${to}; \
	else \
		git log --oneline --no-decorate $${to}; \
	fi

.PHONY: changelog
changelog:
	printf "# Changelog\n\n"; \
	latest=$$(git describe --tags --abbrev=0); \
	changes=$$(make --no-print-directory changes from="$${latest}" | awk '{ print "- " $$0 }'); \
	if [ -n "$${changes}" ]; then \
		url="https://gitlab.com/gitlab-exporter/gitlab-exporter/-/compare/$${latest}..HEAD"; \
		printf "## [Unreleased](%s)\n\n%s\n\n" "$${url}" "$${changes}"; \
	fi; \
	for tag in $$(git tag --list | sort --version-sort --reverse); do \
		previous=$$(git describe --tags --abbrev=0 "$${tag}^" 2>/dev/null); \
		changes=$$(make --no-print-directory changes to=$${tag} | awk '{ print "- " $$0 }'); \
		if [ -n "$${previous}" ]; then \
			url="https://gitlab.com/gitlab-exporter/gitlab-exporter/-/compare/$${previous}..$${tag}"; \
		else \
			url="https://gitlab.com/gitlab-exporter/gitlab-exporter/-/commits/$${tag}"; \
		fi; \
		printf "## [%s](%s)\n\n%s\n\n" "$${tag#v}" "$${url}" "$${changes}"; \
	done

.PHONY: version
version: ## Generate version from git tag and commit information
	git describe --exact-match 2>/dev/null || echo $$(git describe --tags --abbrev=0)-dev.$$(git rev-list --count $$(git describe --tags --abbrev=0)..HEAD)+$$(git rev-parse --short=8 HEAD)
