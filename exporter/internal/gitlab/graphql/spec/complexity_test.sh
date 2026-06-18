#!/bin/sh

export GITLAB_URL="http://localhost:8080"
export GITLAB_TOKEN="glpat-PuryR9JzHKfVx8gvEWJa"

export GITLAB_PROJECT_ID="gid://gitlab/Project/2"
export GITLAB_PROJECT_PATH="project-alpha"
export GITLAB_PIPELINE_IID=""
export GITLAB_MERGE_REQUEST_IID=""

go test -v complexity_test.go
