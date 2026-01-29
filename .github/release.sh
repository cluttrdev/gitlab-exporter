#!/usr/bin/env bash
set -euo pipefail

echo "Fetching GitLab release"
project_id=50817395 # gitlab.com/gitlab-exporter/gitlab-exporter
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
