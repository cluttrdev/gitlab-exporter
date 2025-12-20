REPO_ROOT := $$(git rev-parse --show-toplevel)
BIN_DIR=${REPO_ROOT}/bin
DIST_DIR=${REPO_ROOT}/dist

APPS=$$(ls cmd/)
PLATFORMS=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

.ONESHELL:

ifneq "${VERBOSE}" "1"
.SILENT:
endif

.DEFAULT_GOAL := help

.PHONY: help
help: ## Display this help page
	grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[33m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: tidy
tidy: ## Run go mod tidy on specified module or all modules
ifdef MOD
	go mod tidy -C ${MOD}
else
	find . -type f -name go.mod -exec sh -c 'mod=$$(dirname {}); echo "Tidying $$mod/..."; go mod tidy -C $$mod' \;
endif

.PHONY: fmt
fmt: ## Run go fmt on specified module or all modules
ifdef MOD
	go fmt -C ${MOD} ./...
else
	find . -type f -name go.mod -exec sh -c 'mod=$$(dirname {}); echo "Formatting $$mod/..."; go fmt -C $$mod ./...' \;
endif

.PHONY: vet
vet: ## Run go vet on specified module or all modules
ifdef MOD
	go vet -C ${MOD} ./...
else
	find . -type f -name go.mod -exec sh -c 'mod=$$(dirname {}); echo "Vetting $$mod/..."; go vet -C $$mod ./...' \;
endif

.PHONY: lint
lint: ## Run golangci-lint on specified module or all modules
ifdef MOD
	cd ${MOD} && golangci-lint run ./...
else
	find . -type f -name go.mod -exec sh -c 'mod=$$(dirname {}); echo "Linting $$mod/..."; (cd $$mod && golangci-lint run ./...)' \;
endif

.PHONY: test
test: ## Run tests on specified module or all modules
ifdef MOD
	go test -C ${MOD} ./...
else
	find . -type f -name go.mod -exec sh -c 'mod=$$(dirname {}); echo "Testing $$mod/..."; go test -C $$mod ./...' \;
endif

.PHONY: build
build: ## Build application binary
	./scripts/build.sh -a "${app}" -p "${platform}" binary

.PHONY: build-all
build-all: ## Build all application binaries for all platforms
	./scripts/build.sh --all binary

.PHONY: build-image
build-image: ## Build application container image
	./scripts/build.sh -a "${app}" -p "${platform}" -t "${tag}" image

.PHONY: build-image-all
build-image-all: ## Build container image for each application
	./scripts/build.sh --all -p "${platform}" -t "${tag}" image

.PHONY: build-image-multiplatform
build-image-multiplatform: ## Build multiplatform application container image
	./scripts/build.sh -a "${app}" -p "${platform}" -t "${tag}" --multiplatform image

.PHONY: build-image-multiplatform-all
build-image-multiplatform-all: ## Build multiplatform container image for each application
	./scripts/build.sh --all -t "${tag}" --multiplatform image

.PHONY: dist
dist: ## Build release distribution artifacts
	if [ -z "${app}" ]; then apps="${APPS}"; else apps="${app}"; fi; \
	if [ -z "${platform}" ]; then platforms="${PLATFORMS}"; else platforms="${platform}"; fi; \
	mkdir -p ${DIST_DIR}; \
	version=$$(make --no-print-directory version | tr '+' '+'); \
	for app in $${apps}; do \
		for platform in $${platforms}; do \
			echo "Building $${app} for $${platform}..."; \
			$(MAKE) --no-print-directory build app="$${app}" platform="$${platform}"; \
			os_arch="$$(echo $${platform} | tr '/' '_')"; \
			binary="${BIN_DIR}/$${os_arch}/$${app}"; \
			archive="$${app}_$${version}_$${os_arch}.tar.gz"; \
			tar -czf "${DIST_DIR}/$${archive}" -C "${BIN_DIR}/$${os_arch}" "$${app}"; \
			(cd ${DIST_DIR} && sha256sum "$${archive}" > "$${archive}.sha256"); \
		done; \
	done; \

.PHONY: clean
clean: ## Remove built binaries and distribution artifacts
	rm -rf ${BIN_DIR}/*
	rm -rf ${DIST_DIR}/*

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

.PHONY: release-notes
release-notes: ## Generate release notes
	./scripts/release-notes.sh $(CI_COMMIT_TAG)

.PHONY: version
version: ## Generate version from git tag and commit information
	git describe --exact-match 2>/dev/null || echo $$(git describe --tags --abbrev=0)-dev.$$(git rev-list --count $$(git describe --tags --abbrev=0)..HEAD)+$$(git rev-parse --short=8 HEAD)
