GIT_DIR := `git rev-parse --show-toplevel`
GIT_SHA := `git rev-parse --short HEAD`
GIT_TAG := `git describe --exact-match 2>/dev/null || true`

GOHOSTOS := `go env GOHOSTOS`
GOHOSTARCH := `go env GOHOSTARCH`

BIN_NAME := "gitlab-clickhouse-exporter"
BIN_DIR := "bin"
DIST_DIR := "dist"

VERSION := if GIT_TAG != "" { GIT_TAG } else { "devel" }

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
    go build -o {{BIN_DIR}}/

# create binary distribution
dist:
    #!/usr/bin/env bash
    set -euo pipefail

    version={{VERSION}}

    declare -A OSARCHMAP=(
        [linux]="amd64,arm,arm64"
        [darwin]="amd64,arm64"
    )
    for os in ${!OSARCHMAP[@]}; do
        for arch in ${OSARCHMAP[$os]//,/ }; do
            tmp_dir={{DIST_DIR}}/{{BIN_NAME}}_${version}_${os}_${arch}

            GOOS=${os} GOARCH=${arch} go build -o ${tmp_dir}/{{BIN_NAME}}
        done
    done

    for dir in $(find {{DIST_DIR}}/ -mindepth 1 -maxdepth 1 -type d); do 
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

    echo "creating release for ${tag}"
    response=$(curl -L -s \
        -H "Accept: application/vnd.github+json" \
        -H "Authorization: Bearer ${GITHUB_TOKEN}" \
        -H "X-GitHub-Api-Version 2022-11-28" \
        -d "{\"tag_name\": \"${tag}\", \"name\": \"${tag}\", \"body\": \"\"}" \
        https://api.github.com/repos/cluttrdev/gitlab-clickhouse-exporter/releases \
        2>/dev/null
    )

    error=$(jq -r '.errors[0].code // empty' <<< $response)
    [ -n "$error" ] && {
        echo $error
        exit 1
    }

    release_id=$(jq -r '.id // empty' <<< $response)
    [ -n "$release_id" ] || {
        echo "No release with tag: ${tag}"
        exit 2
    }

    echo "creating binary distributions"
    just dist

    archives=$(find ${dist_dir}/ -type f -name ${bin_name}_*.tar.gz)
    upload_url=https://uploads.github.com/repos/cluttrdev/gitlab-clickhouse-exporter/releases/${release_id}/assets
    for archive in ${archives}; do
        echo "uploading asset: ${archive}"
        name=$(basename ${archive})
        response=$(curl -L -s \
            -H "Accept: application/vnd.github+json" \
            -H "Authorization: Bearer ${GITHUB_TOKEN}" \
            -H "X-GitHub-Api-Version 2022-11-28" \
            -H "Content-Type: application/octet-stream" \
            ${upload_url}?name=${name} \
            --data-binary "@${archive}"
        )
    done

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

_system-info:
    @echo "{{os()}}-{{arch()}}"
