#!/bin/bash

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

export API_V4_URL=gitlab.com/api/v4
export PROJECT_ID=50817395

echo "${API_V4_URL}/projects/${PROJECT_ID}"
glab api projects/${PROJECT_ID}?statistics=true \
    | jq -r \
    > ${SCRIPT_DIR}/${API_V4_URL}/projects/${PROJECT_ID}.json

export PROJECT_DIR=${SCRIPT_DIR}/${API_V4_URL}/projects/${PROJECT_ID}
mkdir -p ${PROJECT_DIR}

echo "${API_V4_URL}/projects/${PROJECT_ID}/pipelines"
pipelines=$(glab api --paginate projects/${PROJECT_ID}/pipelines)

echo "${pipelines}" | jq -r > ${PROJECT_DIR}/pipelines.json
mkdir -p ${PROJECT_DIR}/pipelines

fetch_pipeline_hierarchy() {
    local pipeline_id=$1

    echo "${API_V4_URL}/projects/${PROJECT_ID}/pipelines/${pipeline_id}"
    glab api projects/${PROJECT_ID}/pipelines/${pipeline_id} \
        | jq -r \
        > ${PROJECT_DIR}/pipelines/${pipeline_id}.json

    mkdir -p ${PROJECT_DIR}/pipelines/${pipeline_id}
    for resource in jobs bridges test_report test_report_summary; do
        echo "${API_V4_URL}/projects/${PROJECT_ID}/pipelines/${pipeline_id}/${resource}"
        glab api --paginate projects/${PROJECT_ID}/pipelines/${pipeline_id}/${resource} \
            | jq -r \
            > ${PROJECT_DIR}/pipelines/${pipeline_id}/${resource}.json
    done
}
export -f fetch_pipeline_hierarchy

echo "${pipelines}" | jq -r '.[].id' | xargs -P 12 -I {} sh -c 'fetch_pipeline_hierarchy "$@"' _ {}

