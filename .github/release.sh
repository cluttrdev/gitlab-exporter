#!/usr/bin/env bash
set -euo pipefail

owner=cluttrdev
repo=gitlab-clickhouse-exporter

tag=${RELEASE_TAG}
assets=${RELEASE_ASSETS}

echo "creating release for ${tag}"
response=$(curl -L -s \
    -H "Accept: application/vnd.github+json" \
    -H "Authorization: Bearer ${GITHUB_TOKEN}" \
    -H "X-GitHub-Api-Version 2022-11-28" \
    -d "{\"tag_name\": \"${tag}\", \"name\": \"${tag}\", \"body\": \"\"}" \
    https://api.github.com/repos/${owner}/${repo}/releases \
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

upload_url=https://uploads.github.com/repos/${owner}/${repo}/releases/${release_id}/assets
for asset in ${assets}; do
    name=$(basename ${asset})
    echo "uploading asset: ${name}"
    response=$(curl -L -s \
        -H "Accept: application/vnd.github+json" \
        -H "Authorization: Bearer ${GITHUB_TOKEN}" \
        -H "X-GitHub-Api-Version 2022-11-28" \
        -H "Content-Type: application/octet-stream" \
        ${upload_url}?name=${name} \
        --data-binary "@${asset}"
    )
done
