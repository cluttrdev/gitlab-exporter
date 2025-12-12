#!/usr/bin/env bash

DIST_DIR="${DIST_DIR:-dist}"
BIN_NAME=gitlab-exporter-clickhouse-recorder

set -euo pipefail

printf "Creating binary distribution archives...\n"

declare -A OSARCHMAP=(
    [linux]="amd64 arm64 arm"
    [darwin]="amd64 arm64"
    [windows]="amd64"
)

for os in "${!OSARCHMAP[@]}"; do
    for arch in ${OSARCHMAP[$os]}; do
        printf "\tos=%s arch=%s" "${os}" "${arch}"
        ext=""
        [ "${os}" = "windows" ] && ext=".exe"

        output_dir="${DIST_DIR}/${BIN_NAME}_${CI_COMMIT_TAG}_${os}_${arch}"
        output_file="${BIN_NAME}${ext}"
        output="${output_dir}/${output_file}"

        printf " build..."
        make build os="${os}" arch="${arch}" output="${output}"

        printf " archive..."
        case "${os}" in
            linux | darwin)
                tar -C "${output_dir}" -czf "${output_dir}.tar.gz" "${output_file}"
                ;;
            windows)
                "${GOPATH}"/bin/arc archive "${output_dir}.zip" "${output}"
                ;;
        esac
        rm -r "${output_dir}"

        printf " done\n"
    done
done

printf "Creating release..."
# Create changelog
git fetch --quiet --tags
changes=$(make changes from="$(git describe --tags --abbrev=0 ${CI_COMMIT_TAG}^)" to="${CI_COMMIT_TAG}")
changelog=$(awk '{ printf "%s\\n", $0}' <<-EOF
## Changelog
$(awk '{ print "  - [" $1 "][" $1 "] " substr($0, index($0, $2)) }' <<< "${changes}")

<!-- Links -->
$(awk -v commit_url="${CI_PROJECT_URL}/-/commit/" '{ print "[" $1 "]: " commit_url $1 }' <<< "${changes}")
EOF
)

curl \
    -sSL --fail-with-body \
    --request POST \
    --header "JOB-TOKEN: ${CI_JOB_TOKEN}" \
    --header "Content-Type: application/json" \
    --data "{\"tag_name\":\"${CI_COMMIT_TAG}\", \"description\":\"${changelog}\"}" \
    "${CI_API_V4_URL}/projects/${CI_PROJECT_ID}/releases"

printf "Uploading release assets...\n"
package_url="${CI_API_V4_URL}/projects/${CI_PROJECT_ID}/packages/generic"
assets=$(find dist/ -type f -regextype posix-extended -regex ".*/${BIN_NAME}_${CI_COMMIT_TAG}_.*(\.tar\.gz|\.zip)")
for asset in ${assets}; do
    asset_name=$(basename "${asset}")
    asset_url="${package_url}/${BIN_NAME}/${CI_COMMIT_TAG}/${asset_name}"
    printf "\tname=%s" "${asset_name}"

    printf " upload..."
    curl \
        -sSL --fail-with-body \
        --header "JOB-TOKEN: ${CI_JOB_TOKEN}" \
        --upload-file "${asset}" \
        "${asset_url}"

    printf " link..."
    curl \
        -sSL --fail-with-body \
        --request POST \
        --header "JOB-TOKEN: ${CI_JOB_TOKEN}" \
        --data name="${asset_name}" \
        --data direct_asset_path="/${asset_name}" \
        --data link_type="package" \
        --data url="${asset_url}" \
        "${CI_API_V4_URL}/projects/${CI_PROJECT_ID}/releases/${CI_COMMIT_TAG}/assets/links"

    printf " done\n"
done
