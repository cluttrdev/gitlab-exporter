package clickhouse

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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
	opt := DeduplicateTableOptions{
		Database:    "",
		Table:       "pipelines",
		Final:       &[]bool{false}[0],
		By:          []string{},
		Except:      []string{},
		ThrowIfNoop: nil,
	}

	expectedQuery := "OPTIMIZE TABLE {database:Identifier}.{table:Identifier} DEDUPLICATE"
	expectedParams := map[string]string{
		"database": "gitlab_ci",
		"table":    "pipelines",
	}

	query, params := PrepareDeduplicateQuery(opt)

	checkQuery(t, expectedQuery, query)
	checkParams(t, expectedParams, params)
}

func TestPrepareDeduplicateQuery_Full(t *testing.T) {
	opt := DeduplicateTableOptions{
		Database:    "gitlab_ci",
		Table:       "pipelines",
		Final:       &[]bool{true}[0],
		By:          []string{"id", "project_id"},
		Except:      []string{"finished_at", "status"},
		ThrowIfNoop: &[]bool{true}[0],
	}

	expectedQuery := "" +
		"OPTIMIZE TABLE {database:Identifier}.{table:Identifier} FINAL DEDUPLICATE" +
		" BY id,project_id" +
		" EXCEPT (finished_at,status)" +
		" SETTINGS optimize_throw_if_noop=1"
	expectedParams := map[string]string{
		"database": "gitlab_ci",
		"table":    "pipelines",
	}

	query, params := PrepareDeduplicateQuery(opt)

	checkQuery(t, expectedQuery, query)
	checkParams(t, expectedParams, params)
}

func TestPrepareDeduplicateQuery_WithFinal(t *testing.T) {
	opt := DeduplicateTableOptions{
		Database:    "",
		Table:       "pipelines",
		Final:       &[]bool{true}[0],
		By:          []string{},
		Except:      []string{},
		ThrowIfNoop: nil,
	}

	expectedQuery := "OPTIMIZE TABLE {database:Identifier}.{table:Identifier} FINAL DEDUPLICATE"
	expectedParams := map[string]string{
		"database": "gitlab_ci",
		"table":    "pipelines",
	}

	query, params := PrepareDeduplicateQuery(opt)

	checkQuery(t, expectedQuery, query)
	checkParams(t, expectedParams, params)
}

func TestPrepareDeduplicateQuery_WithThrowIfNoopTrue(t *testing.T) {
	opt := DeduplicateTableOptions{
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

	query, params := PrepareDeduplicateQuery(opt)

	checkQuery(t, expectedQuery, query)
	checkParams(t, expectedParams, params)
}

func TestPrepareDeduplicateQuery_WithThrowIfNoopFalse(t *testing.T) {
	opt := DeduplicateTableOptions{
		Database:    "",
		Table:       "pipelines",
		Final:       &[]bool{false}[0],
		By:          []string{},
		Except:      []string{},
		ThrowIfNoop: &[]bool{false}[0],
	}

	expectedQuery := "" +
		"OPTIMIZE TABLE {database:Identifier}.{table:Identifier} DEDUPLICATE" +
		" SETTINGS optimize_throw_if_noop=0"
	expectedParams := map[string]string{
		"database": "gitlab_ci",
		"table":    "pipelines",
	}

	query, params := PrepareDeduplicateQuery(opt)

	checkQuery(t, expectedQuery, query)
	checkParams(t, expectedParams, params)
}

func TestPrepareDeduplicateQuery_WithBy(t *testing.T) {
	opt := DeduplicateTableOptions{
		Database:    "",
		Table:       "pipelines",
		Final:       &[]bool{false}[0],
		By:          []string{"id", "project_id"},
		Except:      []string{},
		ThrowIfNoop: nil,
	}

	expectedQuery := "" +
		"OPTIMIZE TABLE {database:Identifier}.{table:Identifier} DEDUPLICATE" +
		" BY id,project_id"
	expectedParams := map[string]string{
		"database": "gitlab_ci",
		"table":    "pipelines",
	}

	query, params := PrepareDeduplicateQuery(opt)

	checkQuery(t, expectedQuery, query)
	checkParams(t, expectedParams, params)
}

func TestPrepareDeduplicateQuery_WithSingleExcept(t *testing.T) {
	opt := DeduplicateTableOptions{
		Database:    "",
		Table:       "pipelines",
		Final:       &[]bool{false}[0],
		By:          []string{},
		Except:      []string{"project_id"},
		ThrowIfNoop: nil,
	}

	expectedQuery := "" +
		"OPTIMIZE TABLE {database:Identifier}.{table:Identifier} DEDUPLICATE" +
		" BY * EXCEPT project_id"
	expectedParams := map[string]string{
		"database": "gitlab_ci",
		"table":    "pipelines",
	}

	query, params := PrepareDeduplicateQuery(opt)

	checkQuery(t, expectedQuery, query)
	checkParams(t, expectedParams, params)
}

func TestPrepareDeduplicateQuery_WithMultipleExcept(t *testing.T) {
	opt := DeduplicateTableOptions{
		Database:    "",
		Table:       "pipelines",
		Final:       &[]bool{false}[0],
		By:          []string{},
		Except:      []string{"project_id", "status"},
		ThrowIfNoop: nil,
	}

	expectedQuery := "" +
		"OPTIMIZE TABLE {database:Identifier}.{table:Identifier} DEDUPLICATE" +
		" BY * EXCEPT (project_id,status)"
	expectedParams := map[string]string{
		"database": "gitlab_ci",
		"table":    "pipelines",
	}

	query, params := PrepareDeduplicateQuery(opt)

	checkQuery(t, expectedQuery, query)
	checkParams(t, expectedParams, params)
}

func TestPrepareDeduplicateQuery_WithByAndExcept(t *testing.T) {
	opt := DeduplicateTableOptions{
		Database:    "",
		Table:       "pipelines",
		Final:       &[]bool{false}[0],
		By:          []string{"id"},
		Except:      []string{"project_id", "status"},
		ThrowIfNoop: nil,
	}

	expectedQuery := "" +
		"OPTIMIZE TABLE {database:Identifier}.{table:Identifier} DEDUPLICATE" +
		" BY id EXCEPT (project_id,status)"
	expectedParams := map[string]string{
		"database": "gitlab_ci",
		"table":    "pipelines",
	}

	query, params := PrepareDeduplicateQuery(opt)

	checkQuery(t, expectedQuery, query)
	checkParams(t, expectedParams, params)
}
