package tasks_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/tasks"
)

func checkQuery(t *testing.T, want string, got string) {
	if want != got {
		t.Errorf("Expected `%s`, got `%s`", want, got)
	}
}

func checkParams(t *testing.T, want map[string]string, got map[string]string) {
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Config mismatch (-want +got):\n%s", diff)
	}
}

func TestPrepareDeduplicateQuery_Minimal(t *testing.T) {
	opt := tasks.DeduplicateTableOptions{
		Database:    "",
		Table:       "pipelines",
		Final:       &[]bool{false}[0],
		By:          []string{},
		Except:      []string{},
		ThrowIfNoop: &[]bool{false}[0],
	}

	expectedQuery := "OPTIMIZE TABLE {database:Identifier}.{table:Identifier} DEDUPLICATE"
	expectedParams := map[string]string{
		"database": "gitlab_ci",
		"table":    "pipelines",
	}

	query, params := tasks.PrepareDeduplicateQuery(opt)

	checkQuery(t, expectedQuery, query)
	checkParams(t, expectedParams, params)
}

func TestPrepareDeduplicateQuery_Full(t *testing.T) {
	opt := tasks.DeduplicateTableOptions{
		Database:    "gitlab_ci",
		Table:       "pipelines",
		Final:       &[]bool{true}[0],
		By:          []string{"id", "project_id"},
		Except:      []string{"finished_at", "status"},
		ThrowIfNoop: &[]bool{true}[0],
	}

	expectedQuery := "" +
		"OPTIMIZE TABLE {database:Identifier}.{table:Identifier} FINAL DEDUPLICATE" +
		" BY {by_1:Identifier},{by_2:Identifier}" +
		" EXCEPT ({except_1:Identifier},{except_2:Identifier})" +
		" SETTINGS optimize_throw_if_noop=1"
	expectedParams := map[string]string{
		"database": "gitlab_ci",
		"table":    "pipelines",
		"by_1":     "id",
		"by_2":     "project_id",
		"except_1": "finished_at",
		"except_2": "status",
	}

	query, params := tasks.PrepareDeduplicateQuery(opt)

	checkQuery(t, expectedQuery, query)
	checkParams(t, expectedParams, params)
}

func TestPrepareDeduplicateQuery_WithFinal(t *testing.T) {
	opt := tasks.DeduplicateTableOptions{
		Database:    "",
		Table:       "pipelines",
		Final:       &[]bool{true}[0],
		By:          []string{},
		Except:      []string{},
		ThrowIfNoop: &[]bool{false}[0],
	}

	expectedQuery := "OPTIMIZE TABLE {database:Identifier}.{table:Identifier} FINAL DEDUPLICATE"
	expectedParams := map[string]string{
		"database": "gitlab_ci",
		"table":    "pipelines",
	}

	query, params := tasks.PrepareDeduplicateQuery(opt)

	checkQuery(t, expectedQuery, query)
	checkParams(t, expectedParams, params)
}

func TestPrepareDeduplicateQuery_WithThrowIfNoop(t *testing.T) {
	opt := tasks.DeduplicateTableOptions{
		Database:    "",
		Table:       "pipelines",
		Final:       &[]bool{false}[0],
		By:          []string{},
		Except:      []string{},
		ThrowIfNoop: &[]bool{true}[0],
	}

	expectedQuery := "" +
		"OPTIMIZE TABLE {database:Identifier}.{table:Identifier} DEDUPLICATE" +
		" SETTINGS optimize_throw_if_noop=1"
	expectedParams := map[string]string{
		"database": "gitlab_ci",
		"table":    "pipelines",
	}

	query, params := tasks.PrepareDeduplicateQuery(opt)

	checkQuery(t, expectedQuery, query)
	checkParams(t, expectedParams, params)
}

func TestPrepareDeduplicateQuery_WithBy(t *testing.T) {
	opt := tasks.DeduplicateTableOptions{
		Database:    "",
		Table:       "pipelines",
		Final:       &[]bool{false}[0],
		By:          []string{"id", "project_id"},
		Except:      []string{},
		ThrowIfNoop: &[]bool{false}[0],
	}

	expectedQuery := "" +
		"OPTIMIZE TABLE {database:Identifier}.{table:Identifier} DEDUPLICATE" +
		" BY {by_1:Identifier},{by_2:Identifier}"
	expectedParams := map[string]string{
		"database": "gitlab_ci",
		"table":    "pipelines",
		"by_1":     "id",
		"by_2":     "project_id",
	}

	query, params := tasks.PrepareDeduplicateQuery(opt)

	checkQuery(t, expectedQuery, query)
	checkParams(t, expectedParams, params)
}

func TestPrepareDeduplicateQuery_WithSingleExcept(t *testing.T) {
	opt := tasks.DeduplicateTableOptions{
		Database:    "",
		Table:       "pipelines",
		Final:       &[]bool{false}[0],
		By:          []string{},
		Except:      []string{"project_id"},
		ThrowIfNoop: &[]bool{false}[0],
	}

	expectedQuery := "" +
		"OPTIMIZE TABLE {database:Identifier}.{table:Identifier} DEDUPLICATE" +
		" BY * EXCEPT {except_1:Identifier}"
	expectedParams := map[string]string{
		"database": "gitlab_ci",
		"table":    "pipelines",
		"except_1": "project_id",
	}

	query, params := tasks.PrepareDeduplicateQuery(opt)

	checkQuery(t, expectedQuery, query)
	checkParams(t, expectedParams, params)
}

func TestPrepareDeduplicateQuery_WithMultipleExcept(t *testing.T) {
	opt := tasks.DeduplicateTableOptions{
		Database:    "",
		Table:       "pipelines",
		Final:       &[]bool{false}[0],
		By:          []string{},
		Except:      []string{"project_id", "status"},
		ThrowIfNoop: &[]bool{false}[0],
	}

	expectedQuery := "" +
		"OPTIMIZE TABLE {database:Identifier}.{table:Identifier} DEDUPLICATE" +
		" BY * EXCEPT ({except_1:Identifier},{except_2:Identifier})"
	expectedParams := map[string]string{
		"database": "gitlab_ci",
		"table":    "pipelines",
		"except_1": "project_id",
		"except_2": "status",
	}

	query, params := tasks.PrepareDeduplicateQuery(opt)

	checkQuery(t, expectedQuery, query)
	checkParams(t, expectedParams, params)
}

func TestPrepareDeduplicateQuery_WithByAndExcept(t *testing.T) {
	opt := tasks.DeduplicateTableOptions{
		Database:    "",
		Table:       "pipelines",
		Final:       &[]bool{false}[0],
		By:          []string{"id"},
		Except:      []string{"project_id", "status"},
		ThrowIfNoop: &[]bool{false}[0],
	}

	expectedQuery := "" +
		"OPTIMIZE TABLE {database:Identifier}.{table:Identifier} DEDUPLICATE" +
		" BY {by_1:Identifier} EXCEPT ({except_1:Identifier},{except_2:Identifier})"
	expectedParams := map[string]string{
		"database": "gitlab_ci",
		"table":    "pipelines",
		"by_1":     "id",
		"except_1": "project_id",
		"except_2": "status",
	}

	query, params := tasks.PrepareDeduplicateQuery(opt)

	checkQuery(t, expectedQuery, query)
	checkParams(t, expectedParams, params)
}
