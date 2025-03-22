#!/usr/bin/env bash
set -euo pipefail

echo "Fetching GitLab release"
project_id=50817395 # akun73/gitlab-exporter
response=$(
    curl -sSL --fail-with-body \
        "https://gitlab.com/api/v4/projects/${project_id}/releases/${TAG_NAME}"
)
if [ $? -ne 0 ]; then
    echo "${response}"
    exit 1
fi

# extract release notes
release_notes=$(jq -r '.description' <<< "${response}" | sed 's/$/\\n/' | tr -d '\n')

# adjust commit link url
release_notes=$(sed 's#gitlab.com/akun73/gitlab-exporter/-/commit/#github.com/cluttrdev/gitlab-exporter/commit/#g' <<< "${release_notes}")

release_assets=$(jq -r '.assets.links[]|[.name, .url] | @tsv' <<< "${response}")

echo "Creating GitHub release"
response=$(
    curl "https://api.github.com/repos/cluttrdev/gitlab-exporter/releases" \
        -sSL --fail-with-body \
        --header "Accept: application/vnd.github+json" \
        --header "Authorization: Bearer ${GITHUB_TOKEN}" \
        --header "X-GitHub-Api-Version 2022-11-28" \
        --data "{\"tag_name\":\"${TAG_NAME}\", \"name\":\"${TAG_NAME}\", \"body\":\"${release_notes}\"}"
)
if [ $? -ne 0 ]; then
    echo "${response}"
    exit 1
fi

upload_url=$(jq -r '.upload_url // empty' <<<"${response}")
[ -n "${upload_url}" ] || {
    echo "Missing upload url"
    echo "${response}"
    exit 1
}
upload_url="${upload_url%\{*\}}"

echo "Synching release assets"
echo "${release_assets}" | while IFS=$'\t' read -r name url; do
    printf "\t%s" "${name}"

    printf " donwload..."
    curl \
        -sSL --fail-with-body \
        --output "${name}" \
        "${url}"

    printf " upload..."
    curl "${upload_url}?name=${name}" \
        -sSL --fail-with-body \
        --header "Accept: application/vnd.github+json" \
        --header "Authorization: Bearer ${GITHUB_TOKEN}" \
        --header "X-GitHub-Api-Version 2022-11-28" \
        --header "Content-Type: application/octet-stream" \
        --data-binary "@${name}"

    printf " done\n"
done
