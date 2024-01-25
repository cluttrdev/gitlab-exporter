GIT_DIR := `git rev-parse --show-toplevel`

MAIN := "."
BIN_NAME := `basename {{GIT_DIR}}`
BIN_DIR := "bin"
DIST_DIR := "dist"

GITHUB_OWNER := "cluttrdev"
GITHUB_REPO := "gitlab-exporter"

# list available recipes
default:
    @just --list

# format code
fmt:
    go fmt ./...

# lint code
lint:
    golangci-lint run ./...

# vet code
vet:
    go vet ./...

# build application
build *args="":
    {{GIT_DIR}}/scripts/build.sh -p {{MAIN}} {{args}}

# create binary distribution
dist *args="":
    {{GIT_DIR}}/scripts/dist.sh -p {{MAIN}} {{args}}

# create a new release
release *args="":
    #!/bin/sh
    export GITHUB_OWNER={{GITHUB_OWNER}}
    export GITHUB_REPO={{GITHUB_REPO}}
    {{GIT_DIR}}/scripts/release.sh -p {{MAIN}} {{args}}

changes from="" to="":
    #!/bin/sh
    source {{GIT_DIR}}/scripts/functions.sh
    get_changes {{from}} {{to}}

clean:
    @# build artifacts
    @echo "rm {{BIN_DIR}}/{{BIN_NAME}}"
    @-[ -f {{BIN_DIR}}/{{BIN_NAME}} ] && rm {{BIN_DIR}}/{{BIN_NAME}}
    @-[ -d {{BIN_DIR}} ] && rmdir {{BIN_DIR}}

    @# distribution binaries
    @echo "rm {{DIST_DIR}}/{{BIN_NAME}}_*"
    @rm {{DIST_DIR}}/{{BIN_NAME}}_* 2>/dev/null || true
    @-[ -d {{DIST_DIR}} ] && rmdir {{DIST_DIR}}

###############################################################################

proto-gen:
    #!/bin/sh
    protoc \
        -I protos/ \
        -I protos/vendor/opentelemetry-proto \
        --go_out=. --go_opt=module=github.com/cluttrdev/gitlab-exporter \
        --go-grpc_out=. --go-grpc_opt=module=github.com/cluttrdev/gitlab-exporter \
        protos/gitlabexporter/proto/models/* protos/gitlabexporter/proto/service/*

docker-build:
    #!/bin/sh

    github_owner={{GITHUB_OWNER}}
    github_repo={{GITHUB_REPO}}
    git_dir={{GIT_DIR}}

    set -eu

    source ${git_dir}/scripts/functions.sh

    image=${github_owner}/${github_repo}
    version=$(get_version)
    if is_dirty; then
        version="${version}-modified"
    fi

    docker build \
        -f Dockerfile \
        --build-arg VERSION=${version} \
        -t ${image}:${version} \
        .

    if is_tagged; then
        docker tag ${image}:${version} ${image}:latest
    fi

docker-push: docker-build
    #!/bin/sh

    github_owner={{GITHUB_OWNER}}
    github_owner={{GITHUB_REPO}}
    git_dir={{GIT_DIR}}

    set -eu

    source ${git_dir}/scripts/functions.sh

    if is_dirty; then
        echo "Working diretory is dirty"
        exit 1
    fi

    image=${github_owner}/${github_repo}
    version=$(get_version)

    docker push ${image}:${version}

    if is_tagged; then
        docker push ${image}:latest
    fi
