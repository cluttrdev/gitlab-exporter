#!/bin/sh

export GITLAB_URL="http://localhost:8080"
# export GITLAB_TOKEN=""

export GITLAB_PROJECT_ID="gid://gitlab/Project/2"
export GITLAB_PROJECT_PATH="project-alpha"
export GITLAB_PIPELINE_IID="" # default: 1
export GITLAB_JOB_ID="" # default: gid://gitlab/Ci::Build/1
export GITLAB_MERGE_REQUEST_IID="" # default: 1

# Version gating. Unset (the default) measures the full query, so fields that
# do not exist yet on this server show up as "Field 'x' doesn't exist" - which
# is how you find fields that need a @minVersion annotation.
#
#   auto      gate to this server's own version (detected via metadata)
#   <version> gate to an explicit version, e.g. 17.7
export GITLAB_VERSION_GATE="auto"

go test -v complexity_test.go
