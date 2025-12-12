-- projects
CREATE TABLE IF NOT EXISTS projects (
    id Int64,
    namespace_id Int64,
    owner_id Int64,
    creator_id Int64,
    name String,
    name_with_namespace String,
    path String,
    path_with_namespace String,
    description String,
    visibility String,
    created_at Float64,
    updated_at Float64,
    last_activity_at Float64,
    topics Array(String),
    default_branch String,
    empty_repo Bool,
    archived Bool,
    forks_count Int64,
    stars_count Int64,
    commit_count Int64,
    storage_size Int64,
    repository_size Int64,
    wiki_size Int64,
    lfs_objects_size Int64,
    job_artifacts_size Int64,
    pipeline_artifacts_size Int64,
    packages_size Int64,
    snippets_size Int64,
    uploads_size Int64,
    open_issues_count Int64,
    web_url String
)
ENGINE ReplacingMergeTree(last_activity_at)
ORDER BY id
;
