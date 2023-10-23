package tasks

import (
	"context"
	"fmt"
	"strings"

	clickhouse "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/clickhouse"
)

type DeduplicateTableOptions struct {
	Database    string
	Table       string
	Final       *bool
	By          []string
	Except      []string
	ThrowIfNoop *bool
}

func DeduplicateTable(ctx context.Context, opt DeduplicateTableOptions, ch *clickhouse.Client) error {
	query, params := PrepareDeduplicateQuery(opt)

	ctx = clickhouse.WithParameters(ctx, params)
	return ch.Exec(ctx, query)
}

func PrepareDeduplicateQuery(opt DeduplicateTableOptions) (string, map[string]string) {
	var (
		query  = ""
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
	var byPrefix string = "by_"
	var byParams []string = columnParameterList(len(opt.By), byPrefix)

	if len(byParams) > 0 {
		query += " BY " + strings.Join(byParams, ",")
	} else if len(opt.Except) > 0 {
		query += " BY *"
	}

	for i, val := range opt.By {
		name := columnParameterName(byPrefix, i)
		params[name] = val
	}

	// EXCEPT
	var exceptPrefix string = "except_"
	var exceptParams []string = columnParameterList(len(opt.Except), exceptPrefix)

	if len(exceptParams) == 1 {
		query += " EXCEPT " + exceptParams[0]
	} else if len(exceptParams) > 1 {
		query += " EXCEPT (" + strings.Join(exceptParams, ",") + ")"
	}

	for i, val := range opt.Except {
		name := columnParameterName(exceptPrefix, i)
		params[name] = val
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

func columnParameterName(prefix string, i int) string {
	return fmt.Sprintf("%s%d", prefix, i+1)
}

func columnParameterList(n int, prefix string) []string {
	columns := make([]string, n)
	for i := 0; i < n; i++ {
		name := columnParameterName(prefix, i)
		columns[i] = fmt.Sprintf("{%s:Identifier}", name)
	}
	return columns
}
