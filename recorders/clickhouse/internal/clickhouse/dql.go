package clickhouse

import "context"

func SelectPipelineMaxUpdatedAt(c *Client, ctx context.Context) (map[int64]float64, error) {
	const query string = `
        SELECT id, max(updated_at) AS updated_at
        FROM {db:Identifier}.{table:Identifier}
        GROUP BY id
        `
	var params = map[string]string{
		"db":    c.dbName,
		"table": PipelinesTable,
	}

	ctx = WithParameters(ctx, params)

	var results []struct {
		ID        int64   `ch:"id"`
		UpdatedAt float64 `ch:"updated_at"`
	}

	if err := c.Select(ctx, &results, query); err != nil {
		return nil, err
	}

	m := make(map[int64]float64, len(results))
	for _, res := range results {
		m[res.ID] = res.UpdatedAt
	}

	return m, nil
}

func SelectTableIDLastestUpdates(c *Client, ctx context.Context, table string, idColumn string, updatedAtColumn string) (map[int64]float64, error) {
	const query string = `
        SELECT {id:Identifier} AS id, max({updated_at:Identifier}) AS updated_at
        FROM {db:Identifier}.{table:Identifier}
        GROUP BY id
        `
	var params = map[string]string{
		"db":         c.dbName,
		"table":      table,
		"id":         idColumn,
		"updated_at": updatedAtColumn,
	}

	ctx = WithParameters(ctx, params)

	var results []struct {
		ID        int64   `ch:"id"`
		UpdatedAt float64 `ch:"updated_at"`
	}

	if err := c.Select(ctx, &results, query); err != nil {
		return nil, err
	}

	m := make(map[int64]float64, len(results))
	for _, res := range results {
		m[res.ID] = res.UpdatedAt
	}

	return m, nil
}

func SelectTableIDs[T int64 | string](c *Client, ctx context.Context, table string, column string) (map[T]struct{}, error) {
	const query string = `
        SELECT DISTINCT {column:Identifier} AS id FROM {db:Identifier}.{table:Identifier}
        `
	var params = map[string]string{
		"db":     c.dbName,
		"table":  table,
		"column": column,
	}

	ctx = WithParameters(ctx, params)

	var results []struct {
		ID T `ch:"id"`
	}

	if err := c.Select(ctx, &results, query); err != nil {
		return nil, err
	}

	m := make(map[T]struct{}, len(results))
	for _, res := range results {
		m[res.ID] = struct{}{}
	}

	return m, nil
}

func SelectTraceSpanIDs(c *Client, ctx context.Context) (map[string]struct{}, error) {
	const query string = `
        SELECT TraceId, SpanId FROM {db:Identifier}.{table:Identifier}
        `
	var params = map[string]string{
		"db":    c.dbName,
		"table": TraceSpansTable,
	}

	ctx = WithParameters(ctx, params)

	var results []struct {
		TraceId string `ch:"TraceId"`
		SpanId  string `ch:"SpanId"`
	}

	if err := c.Select(ctx, &results, query); err != nil {
		return nil, err
	}

	m := make(map[string]struct{}, len(results))
	var key string
	for _, r := range results {
		key = r.TraceId + "-" + r.SpanId
		m[key] = struct{}{}
	}

	return m, nil
}
