#!/usr/bin/env bash
set -euo pipefail

GITLAB_URL="${GITLAB_URL:-http://localhost:8080}"
ROOT_PASSWORD="${ROOT_PASSWORD:-password1234!}"
ROOT_TOKEN="${ROOT_TOKEN:-glpat-setup-root-token}"
DEBUG="${DEBUG:-false}"

if [ "$DEBUG" = "true" ]; then
  set -x
fi

# --- Helpers ---

log() { echo "==> $*" >&2; }

api() {
  local method=$1 path=$2
  shift 2
  local url="$GITLAB_URL/api/v4$path"
  url="${url// /%20}"
  local response http_code curl_exit
  response=$(curl -s -w "\n%{http_code}" -X "$method" "$url" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    "$@") || curl_exit=$?
  if [ "${curl_exit:-0}" -ne 0 ]; then
    echo "DEBUG: curl failed with exit code $curl_exit for $method /api/v4$path" >&2
    echo "DEBUG: response: $response" >&2
    return "$curl_exit"
  fi
  http_code=$(echo "$response" | tail -n1)
  response=$(echo "$response" | sed '$d')
  if [ "${http_code:0:1}" != "2" ]; then
    echo "DEBUG: API $method /api/v4$path returned HTTP $http_code" >&2
    echo "DEBUG: response body: $response" >&2
    return 1
  fi
  echo "$response"
}

wait_for_gitlab() {
  log "Waiting for GitLab to be ready..."
  local attempts=0
  until curl -sf "$GITLAB_URL/api/v4/version" -H "Private-Token: $ROOT_TOKEN" > /dev/null 2>&1; do
    attempts=$((attempts + 1))
    if [ "$attempts" -ge 120 ]; then
      echo "ERROR: GitLab not ready after 10 minutes" >&2
      exit 1
    fi
    sleep 5
  done
  log "GitLab is healthy"
}

TOKEN=""
get_token() {
  log "Obtaining OAuth token for root..."
  local response curl_exit
  response=$(curl -s -w "\n%{http_code}" "$GITLAB_URL/oauth/token" \
    -d "grant_type=password&username=root&password=$ROOT_PASSWORD") || curl_exit=$?
  if [ "${curl_exit:-0}" -ne 0 ]; then
    echo "ERROR: curl failed with exit code $curl_exit" >&2
    echo "DEBUG: URL: $GITLAB_URL/oauth/token" >&2
    echo "DEBUG: response: $response" >&2
    exit 1
  fi
  local http_code
  http_code=$(echo "$response" | tail -n1)
  response=$(echo "$response" | sed '$d')
  if [ "${http_code:0:1}" != "2" ]; then
    echo "ERROR: OAuth token request returned HTTP $http_code" >&2
    echo "DEBUG: response body: $response" >&2
    exit 1
  fi
  log "DEBUG: OAuth response: $response"
  TOKEN=$(echo "$response" | jq -r '.access_token')
  if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
    echo "ERROR: Failed to extract access_token from response" >&2
    echo "DEBUG: parsed token value: '$TOKEN'" >&2
    echo "DEBUG: full response: $response" >&2
    exit 1
  fi
  log "Token obtained"
}

ensure_user() {
  local username=$1 name=$2 email=$3
  local id
  id=$(api GET "/users?username=$username" | jq '.[0].id // empty')
  if [ -n "$id" ]; then
    log "User '$username' already exists (id=$id)"
    echo "$id"
    return
  fi
  id=$(api POST "/users" -d "$(jq -n \
    --arg u "$username" --arg n "$name" --arg e "$email" \
    '{username:$u, name:$n, email:$e, password:"Xk9#mP2vL7qR4!w", skip_confirmation:true}')" \
    | jq '.id')
  log "Created user '$username' (id=$id)"
  echo "$id"
}

ensure_group() {
  local name=$1 path=$2
  local id
  id=$(api GET "/groups?search=$path" | jq --arg p "$path" '[.[] | select(.path == $p)][0].id // empty')
  if [ -n "$id" ]; then
    log "Group '$path' already exists (id=$id)"
    echo "$id"
    return
  fi
  id=$(api POST "/groups" -d "$(jq -n \
    --arg n "$name" --arg p "$path" \
    '{name:$n, path:$p, visibility:"public"}')" \
    | jq '.id')
  log "Created group '$path' (id=$id)"
  echo "$id"
}

ensure_group_member() {
  local group_id=$1 user_id=$2 access_level=${3:-30}
  local existing
  existing=$(api GET "/groups/$group_id/members" | jq --argjson uid "$user_id" '[.[] | select(.id == $uid)][0].id // empty')
  if [ -n "$existing" ]; then
    return
  fi
  api POST "/groups/$group_id/members" -d "$(jq -n \
    --argjson uid "$user_id" --argjson al "$access_level" \
    '{user_id:$uid, access_level:$al}')" > /dev/null
  log "Added user $user_id to group $group_id"
}

ensure_project() {
  local name=$1 namespace_id=$2
  local id
  if [ -z "$namespace_id" ]; then
    id=$(api GET "/projects?search=$name&owned=true" | jq --arg n "$name" '[.[] | select(.name == $n)][0].id // empty')
  else
    id=$(api GET "/groups/$namespace_id/projects?search=$name" | jq --arg n "$name" '[.[] | select(.name == $n)][0].id // empty')
  fi
  if [ -n "$id" ]; then
    log "Project '$name' already exists (id=$id)"
    echo "$id"
    return
  fi
  local payload
  payload=$(jq -n --arg n "$name" '{name:$n, visibility:"public", initialize_with_readme:true}')
  if [ -n "$namespace_id" ]; then
    payload=$(echo "$payload" | jq --argjson ns "$namespace_id" '. + {namespace_id:$ns}')
  fi
  id=$(api POST "/projects" -d "$payload" | jq '.id')
  log "Created project '$name' (id=$id)"
  echo "$id"
}

ensure_label() {
  local project_id=$1 name=$2 color=$3
  local existing
  existing=$(api GET "/projects/$project_id/labels?search=$name" | jq --arg n "$name" '[.[] | select(.name == $n)][0].id // empty')
  if [ -n "$existing" ]; then return; fi
  api POST "/projects/$project_id/labels" -d "$(jq -n \
    --arg n "$name" --arg c "$color" '{name:$n, color:$c}')" > /dev/null
  log "Created label '$name' in project $project_id"
}

ensure_milestone() {
  local project_id=$1 title=$2
  local existing
  existing=$(api GET "/projects/$project_id/milestones" | jq --arg t "$title" '[.[] | select(.title == $t)][0].id // empty')
  if [ -n "$existing" ]; then
    echo "$existing"
    return
  fi
  local id
  id=$(api POST "/projects/$project_id/milestones" -d "$(jq -n --arg t "$title" '{title:$t}')" | jq '.id')
  log "Created milestone '$title' in project $project_id (id=$id)"
  echo "$id"
}

ensure_issue() {
  local project_id=$1 title=$2 labels=$3 milestone_id=$4
  local existing
  existing=$(api GET "/projects/$project_id/issues?search=$title" | jq --arg t "$title" '[.[] | select(.title == $t)][0].iid // empty')
  if [ -n "$existing" ]; then return; fi
  local payload
  payload=$(jq -n --arg t "$title" --arg l "$labels" '{title:$t, labels:$l}')
  if [ -n "$milestone_id" ]; then
    payload=$(echo "$payload" | jq --argjson mid "$milestone_id" '. + {milestone_id:$mid}')
  fi
  api POST "/projects/$project_id/issues" -d "$payload" > /dev/null
  log "Created issue '$title' in project $project_id"
}

ensure_branch_and_mr() {
  local project_id=$1 branch=$2 title=$3 file_path=$4 file_content=$5
  # Check if MR exists
  local existing
  existing=$(api GET "/projects/$project_id/merge_requests?source_branch=$branch&state=opened" | jq '.[0].iid // empty')
  if [ -n "$existing" ]; then
    log "MR '$title' already exists (iid=$existing)"
    return
  fi
  # Create commit on new branch
  api POST "/projects/$project_id/repository/commits" -d "$(jq -n \
    --arg b "$branch" --arg msg "Add $file_path" --arg fp "$file_path" --arg fc "$file_content" \
    '{branch:$b, start_branch:"main", commit_message:$msg, actions:[{action:"create", file_path:$fp, content:$fc}]}')" > /dev/null 2>&1 || true
  # Create MR
  api POST "/projects/$project_id/merge_requests" -d "$(jq -n \
    --arg t "$title" --arg sb "$branch" \
    '{title:$t, source_branch:$sb, target_branch:"main"}')" > /dev/null
  log "Created MR '$title' in project $project_id"
}

ensure_ci_config() {
  local project_id=$1
  # Check if .gitlab-ci.yml exists
  if api GET "/projects/$project_id/repository/files/.gitlab-ci.yml?ref=main" > /dev/null 2>&1; then
    log "CI config already exists in project $project_id"
    return
  fi
  local ci_content
  ci_content='stages:
  - build
  - test
  - deploy

build-job:
  stage: build
  script:
    - echo "Building..."

test-job:
  stage: test
  script:
    - echo "Testing..."

deploy-job:
  stage: deploy
  script:
    - echo "Deploying..."
'
  api POST "/projects/$project_id/repository/commits" -d "$(jq -n \
    --arg content "$ci_content" \
    '{branch:"main", commit_message:"Add CI configuration", actions:[{action:"create", file_path:".gitlab-ci.yml", content:$content}]}')" > /dev/null
  log "Added .gitlab-ci.yml to project $project_id"
}

ensure_exporter_token() {
  local user_id=$1 token_name="gitlab-exporter-token"
  # Revoke existing token if present
  local existing_id
  existing_id=$(api GET "/personal_access_tokens?user_id=$user_id" | jq --arg n "$token_name" '[.[] | select(.name == $n and .revoked == false and .active == true)][0].id // empty')
  if [ -n "$existing_id" ]; then
    log "Revoking existing token '$token_name' (id=$existing_id)"
    api DELETE "/personal_access_tokens/$existing_id" > /dev/null
  fi
  local result
  result=$(api POST "/users/$user_id/personal_access_tokens" -d "$(jq -n \
    --arg n "$token_name" \
    '{name:$n, scopes:["read_api"]}')")
  EXPORTER_TOKEN=$(echo "$result" | jq -r '.token')
  if [ -z "$EXPORTER_TOKEN" ] || [ "$EXPORTER_TOKEN" = "null" ]; then
    log "WARNING: Failed to extract token from response"
    log "DEBUG: response: $result"
    return 1
  fi
  log "Created PAT '$token_name' for gitlab-exporter user"
}

ensure_runner() {
  local existing
  existing=$(api GET "/runners/all" | jq '[.[] | select(.description == "dummy-runner")][0].id // empty')
  if [ -n "$existing" ]; then
    log "Dummy runner already exists (id=$existing)"
    return
  fi
  # Get registration token from admin settings
  local reg_token
  reg_token=$(api GET "/application/settings" | jq -r '.runners_registration_token // empty')
  if [ -z "$reg_token" ]; then
    log "WARNING: Could not get runner registration token, skipping runner setup"
    return
  fi
  # Register runner
  curl -sf -X POST "$GITLAB_URL/api/v4/runners" \
    -d "token=$reg_token&description=dummy-runner&tag_list=dummy&run_untagged=true" > /dev/null
  log "Registered dummy runner"
}

trigger_pipeline() {
  local project_id=$1 ref=${2:-main}
  local existing
  existing=$(api GET "/projects/$project_id/pipelines?ref=$ref&per_page=1" | jq '.[0].id // empty')
  if [ -n "$existing" ]; then
    log "Pipeline already exists for project $project_id ref=$ref (id=$existing)"
    return
  fi
  api POST "/projects/$project_id/pipeline" -d "$(jq -n --arg r "$ref" '{ref:$r}')" > /dev/null
  log "Triggered pipeline for project $project_id ref=$ref"
}

cancel_pipeline() {
  local project_id=$1 ref=${2:-main}
  local pipeline_id
  pipeline_id=$(api GET "/projects/$project_id/pipelines?ref=$ref&per_page=1" | jq '.[0].id // empty')
  if [ -z "$pipeline_id" ]; then
    log "No pipeline found for project $project_id ref=$ref to cancel"
    return
  fi
  local status
  status=$(api GET "/projects/$project_id/pipelines/$pipeline_id" | jq -r '.status')
  if [ "$status" = "canceled" ] || [ "$status" = "success" ] || [ "$status" = "failed" ]; then
    log "Pipeline $pipeline_id already finished (status=$status)"
    return
  fi
  api POST "/projects/$project_id/pipelines/$pipeline_id/cancel" > /dev/null
  log "Canceled pipeline $pipeline_id for project $project_id"
}

# --- Main ---

wait_for_gitlab
get_token

# Users
USER1_ID=$(ensure_user "testuser1" "Test User One" "testuser1@example.com")
USER2_ID=$(ensure_user "testuser2" "Test User Two" "testuser2@example.com")
EXPORTER_USER_ID=$(ensure_user "gitlab-exporter" "GitLab Exporter" "exporter@example.com")
EXPORTER_TOKEN=""
ensure_exporter_token "$EXPORTER_USER_ID"

# Group
GROUP_ID=$(ensure_group "Test Group" "test-group")
ensure_group_member "$GROUP_ID" "$USER1_ID" 30  # Developer
ensure_group_member "$GROUP_ID" "$USER2_ID" 20  # Reporter

# Projects
PROJECT_ALPHA_ID=$(ensure_project "project-alpha" "$GROUP_ID")
PROJECT_BETA_ID=$(ensure_project "project-beta" "")

# Labels
for pid in "$PROJECT_ALPHA_ID" "$PROJECT_BETA_ID"; do
  ensure_label "$pid" "bug" "#d9534f"
  ensure_label "$pid" "feature" "#428bca"
  ensure_label "$pid" "urgent" "#f0ad4e"
done

# Milestones & Issues
for pid in "$PROJECT_ALPHA_ID" "$PROJECT_BETA_ID"; do
  MS1_ID=$(ensure_milestone "$pid" "v1.0")
  MS2_ID=$(ensure_milestone "$pid" "v2.0")

  ensure_issue "$pid" "Fix login page crash" "bug,urgent" "$MS1_ID"
  ensure_issue "$pid" "Add dark mode support" "feature" "$MS2_ID"
  ensure_issue "$pid" "Update dependencies" "bug" ""
done

# Merge Requests
ensure_branch_and_mr "$PROJECT_ALPHA_ID" "feature/add-docs" "Add documentation" \
  "docs/README.md" "# Documentation\n\nThis is the project documentation.\n"
ensure_branch_and_mr "$PROJECT_ALPHA_ID" "fix/typo" "Fix typo in config" \
  "config/fix.txt" "Fixed typo\n"
ensure_branch_and_mr "$PROJECT_BETA_ID" "feature/new-api" "Implement new API endpoint" \
  "src/api.py" "def handler():\n    return {'status': 'ok'}\n"

# CI / Pipelines
ensure_ci_config "$PROJECT_ALPHA_ID"
ensure_ci_config "$PROJECT_BETA_ID"

# Runner
ensure_runner

# Trigger pipelines
trigger_pipeline "$PROJECT_ALPHA_ID" "main"
trigger_pipeline "$PROJECT_BETA_ID" "main"

# Cancel pipelines so they finish without a runner
cancel_pipeline "$PROJECT_ALPHA_ID" "main"
cancel_pipeline "$PROJECT_BETA_ID" "main"

log "Done! Test data setup complete."
log "  GitLab: $GITLAB_URL"
log "  Login:  root / $ROOT_PASSWORD"
log "  Users:  testuser1, testuser2 (password: Xk9#mP2vL7qR4!w)"
log "  Group:  test-group"
log "  Projects: project-alpha (group), project-beta (personal)"
if [ -n "${EXPORTER_TOKEN:-}" ]; then
  log "  Exporter token (read_api): $EXPORTER_TOKEN"
fi
