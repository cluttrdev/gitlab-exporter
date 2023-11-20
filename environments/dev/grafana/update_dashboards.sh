#!/bin/sh

set -eu

git_dir=$(git rev-parse --show-toplevel)
src_dir="${git_dir}/assets/grafana/dashboards"
dst_dir="${git_dir}/environments/dev/grafana/provisioning/dashboards/gitlab-ci-analytics"
cp ${src_dir}/*.json "${dst_dir}/"
