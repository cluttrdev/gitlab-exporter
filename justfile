GIT_DIR := `git rev-parse --show-toplevel`
GIT_SHA := `git rev-parse --short HEAD`
GIT_TAG := `git describe --exact-match 2>/dev/null || true`

GOHOSTOS := `go env GOHOSTOS`
GOHOSTARCH := `go env GOHOSTARCH`

BIN_NAME := "gitlab-clickhouse-exporter"
BIN_DIR := "bin"
DIST_DIR := "dist"

# list available recipes
default:
    @just --list --unsorted

# format code
fmt:
    go fmt ./...

# vet code
vet:
    go vet ./...

# build application
build: vet
    go build -o {{BIN_DIR}}/{{BIN_NAME}}

# create binary distribution
dist:
    #!/usr/bin/env bash
    set -euo pipefail

    dist_dir={{DIST_DIR}}
    bin_name={{BIN_NAME}}
    version=$(just _version)

    declare -A OSARCHMAP=(
        [linux]="amd64,arm,arm64"
        [darwin]="amd64,arm64"
    )
    for os in ${!OSARCHMAP[@]}; do
        for arch in ${OSARCHMAP[$os]//,/ }; do
            tmp_dir=${dist_dir}/${bin_name}_${version}_${os}_${arch}

            GOOS=${os} GOARCH=${arch} go build -o ${tmp_dir}/${bin_name}
        done
    done

    for dir in $(find ${dist_dir}/ -mindepth 1 -maxdepth 1 -type d -name ${bin_name}_${version}_*); do 
        find $dir -printf "%P\n" \
        | tar -czf ${dir}.tar.gz --no-recursion -C ${dir} -T -

        rm -r ${dir}
    done

release: _check-tag
    #!/usr/bin/env bash
    set -euo pipefail

    tag={{GIT_TAG}}
    dist_dir={{DIST_DIR}}
    bin_name={{BIN_NAME}}

    echo "creating binary distributions"
    just dist
    assets=$(find ${dist_dir}/ -type f -name ${bin_name}_${tag}_*.tar.gz)

    export RELEASE_TAG=${tag}
    export RELEASE_ASSETS=${assets}
    ./.github/release.sh

clean:
    @echo "rm {{BIN_DIR}}/{{BIN_NAME}}"
    @rm {{BIN_DIR}}/{{BIN_NAME}} 2>/dev/null || true
    @echo "rmdir {{BIN_DIR}}"
    @rmdir {{BIN_DIR}} 2>/dev/null || true

    @echo "rm -r {{DIST_DIR}}"
    @rm -r {{DIST_DIR}} 2>/dev/null || true

docker-build:
    #!/bin/sh
    set -euo pipefail

    image=cluttrdev/{{BIN_NAME}}
    version=$(just _version)

    docker build . -f build/Dockerfile -t ${image}:${version}

    if just _check-tag > /dev/null 2>&1; then
        docker tag ${image}:${version} ${image}:latest
    fi

docker-push: _check-dirty && docker-build
    #!/bin/sh
    set -euo pipefail

    image=cluttrdev/{{BIN_NAME}}
    version=$(just _version)

    docker push ${image}:${version}

    if just _check-tag > /dev/null 2>&1; then
        docker push ${image}:latest
    fi

# ---

# fail if working directory is dirty
[no-exit-message]
_check-dirty:
    @git diff --quiet || (echo "Working directory is dirty" && exit 1)

# fail if current commit is not tagged
[no-exit-message]
_check-tag:
    #!/bin/sh
    set -euo pipefail

    just _check-dirty

    [ -n "{{GIT_TAG}}" ] || {
        echo "No tag exactly matches current commit"
        exit 2
    }

_version:
    @if [ -n "{{GIT_TAG}}" ]; then echo "{{GIT_TAG}}"; else just _pseudo-version; fi

_pseudo-version prefix="" object="HEAD":
    #!/bin/sh

    latest_tag=$(git describe --tags --abbrev=0 || true)

    if [ -n "{{prefix}}" ]; then
        prefix={{prefix}}
    elif [ -n "${latest_tag}" ]; then
        prefix=${latest_tag}-$(git rev-list ${latest_tag}..{{object}} --count)
    else
        prefix=v0.0.0
    fi

    # UTC time the revision was created (yymmddhhmmss).
    timestamp=$(date -u +%y%m%d%H%M%S -d @$(git log -n 1 --format=%ct {{object}}))

    # 12-character prefix of the commit hash
    revision=$(git rev-parse --short=12 --verify {{object}})

    echo "${prefix}-${timestamp}-${revision}"

_system-info:
    @echo "{{os()}}_{{arch()}}"