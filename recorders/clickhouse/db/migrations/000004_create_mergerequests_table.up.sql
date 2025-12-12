-- mergerequests
CREATE TABLE IF NOT EXISTS mergerequests (
    id Int64,
    iid Int64,
    project_id Int64,

    created_at Float64,
    updated_at Float64,
    merged_at Float64,
    closed_at Float64,

    source_project_id Int64,
    target_project_id Int64,
    source_branch String,
    target_branch String,

    title String,
    state String,
    merge_status String,
    draft Bool,
    has_conflicts Bool,
    merge_error String,

    diff_ref_base_sha String,
    diff_ref_head_sha String,
    diff_ref_start_sha String,

    author_id Int64,
    assignee_id Int64,
    assignees_id Array(Int64),
    reviewers_id Array(Int64),
    merge_user_id Int64,
    close_user_id Int64,

    labels Array(String),

    sha String,
    merge_commit_sha String,
    squash_commit_sha String,

    changes_count String,
    user_notes_count Int64,
    upvotes Int64,
    downvotes Int64,

    head_pipeline_id Int64,
    milestone_id Int64,

    web_url String
)
ENGINE ReplacingMergeTree(updated_at)
ORDER BY id
;
