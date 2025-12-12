package clickhouse

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

type DeduplicateTableOptions struct {
	Database    string
	Table       string
	Final       *bool
	By          []string
	Except      []string
	ThrowIfNoop *bool
}

func DeduplicateTable(ctx context.Context, opt DeduplicateTableOptions, ch *Client) error {
	if err := validateDeduplicateTableOptions(ctx, opt, ch); err != nil {
		return fmt.Errorf("error validating deduplication options: %w", err)
	}

	query, params := PrepareDeduplicateQuery(opt)

	ctx = WithParameters(ctx, params)
	return ch.Exec(ctx, query)
}

func PrepareDeduplicateQuery(opt DeduplicateTableOptions) (string, map[string]string) {
	var (
		query  string
		params = map[string]string{}
	)

	// OPTIMIZE
	query = "OPTIMIZE TABLE {database:Identifier}.{table:Identifier}"

	var dbName string = opt.Database
	if dbName == "" {
		dbName = "gitlab_ci"
	}
	params["database"] = dbName
	params["table"] = opt.Table

	// FINAL
	if opt.Final == nil || *opt.Final {
		query += " FINAL"
	}

	// DEDUPLICATE
	query += " DEDUPLICATE"

	// BY
	if len(opt.By) > 0 {
		query += " BY " + strings.Join(opt.By, ",")
	} else if len(opt.Except) > 0 {
		query += " BY *"
	}

	// EXCEPT
	if len(opt.Except) == 1 {
		query += " EXCEPT " + opt.Except[0]
	} else if len(opt.Except) > 1 {
		query += " EXCEPT (" + strings.Join(opt.Except, ",") + ")"
	}

	// SETTINGS
	if opt.ThrowIfNoop != nil {
		if *opt.ThrowIfNoop {
			query += " SETTINGS optimize_throw_if_noop=1"
		} else {
			query += " SETTINGS optimize_throw_if_noop=0"
		}
	}

	return query, params
}

func validateDeduplicateTableOptions(ctx context.Context, opt DeduplicateTableOptions, ch *Client) error {
	var dbName string = opt.Database
	if dbName == "" {
		dbName = "gitlab_ci"
	}

	// validate database identifier
	if err := matchIdentifier(dbName); err != nil {
		return err
	}

	// validate table identifier
	if err := matchIdentifier(opt.Table); err != nil {
		return err
	}

	// validate column identifiers
	cols, err := getColumnNames(ctx, dbName, opt.Table, ch)
	if err != nil {
		return fmt.Errorf("error getting table column names: %w", err)
	}

	for _, c := range append(append([]string{}, opt.By...), opt.Except...) {
		if !slices.Contains(cols, c) {
			return fmt.Errorf("Table `%s` has no column `%s`", opt.Table, c)
		}
	}

	return nil
}

func matchIdentifier(s string) error {
	pattern := `^[a-zA-Z_][0-9a-zA-Z_]*$`
	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		return err
	} else if !matched {
		return fmt.Errorf("invalid identifier: `%s`", s)
	}
	return nil
}

func getColumnNames(ctx context.Context, database string, table string, ch *Client) ([]string, error) {
	var columnNames []string

	query_tpl := `
    SELECT DISTINCT COLUMN_NAME FROM information_schema.COLUMNS
      WHERE (TABLE_SCHEMA = '%s') AND (TABLE_NAME = '%s')
    `

	query := fmt.Sprintf(query_tpl, database, table)

	var results []struct {
		ColumnName string `ch:"COLUMN_NAME"`
	}

	if err := ch.Select(ctx, &results, query); err != nil {
		return nil, err
	}

	for _, r := range results {
		columnNames = append(columnNames, r.ColumnName)
	}

	return columnNames, nil
}
