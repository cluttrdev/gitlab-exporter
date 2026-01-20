#!/usr/bin/env bash

DIST_DIR="${DIST_DIR:-dist}"

RELEASE_TAG=${CI_COMMIT_TAG:-$(make version)}

set -euo pipefail

printf "Uploading release assets...\n"
package_url="${CI_API_V4_URL}/projects/${CI_PROJECT_ID}/packages/generic"
assets=$(find dist/ -type f -regextype posix-extended -regex ".*/gitlab-exporter(-.*)?_${RELEASE_TAG/+/\\+}_.*(\.tar\.gz|\.zip)(\.sha256)?" | sort)
for asset in ${assets}; do
    asset_name=$(basename "${asset}")
    bin_name=${asset_name%%_*}
    asset_url="${package_url}/${bin_name}/${RELEASE_TAG}/${asset_name}"
    printf "\tname=%s" "${asset_name}"

    printf " upload..."
    curl \
        -sSL --fail-with-body \
        --header "JOB-TOKEN: ${CI_JOB_TOKEN}" \
        --upload-file "${asset}" \
        "${asset_url}"

    # printf " link..."
    # curl \
    #     -sSL --fail-with-body \
    #     --request POST \
    #     --header "JOB-TOKEN: ${CI_JOB_TOKEN}" \
    #     --data name="${asset_name}" \
    #     --data direct_asset_path="/${asset_name}" \
    #     --data link_type="package" \
    #     --data url="${asset_url}" \
    #     "${CI_API_V4_URL}/projects/${CI_PROJECT_ID}/releases/${RELEASE_TAG}/assets/links"

    printf " done\n"
done

printf "Creating release..."
# Create changelog
git fetch --quiet --tags
changes=$(make changes from="$(git describe --tags --abbrev=0 ${RELEASE_TAG}^)" to="${RELEASE_TAG}")
changelog=$(awk '{ printf "%s\\n", $0}' <<-EOF
## Assets

Download links for the release assets can be found here:
[gitlab-exporter.gitlab.io/packages](https://gitlab-exporter.gitlab.io/packages).

## Changelog
$(awk -v commit_url="${CI_PROJECT_URL}/-/commit/" '{ print "  - " substr($0, index($0, $2)) "([" $1 "](" commit_url $1 "))" }' <<< "${changes}")
EOF
)

curl \
    -sSL --fail-with-body \
    --request POST \
    --header "JOB-TOKEN: ${CI_JOB_TOKEN}" \
    --header "Content-Type: application/json" \
    --data "{\"tag_name\":\"${RELEASE_TAG}\", \"description\":\"${changelog}\"}" \
    "${CI_API_V4_URL}/projects/${CI_PROJECT_ID}/releases"
