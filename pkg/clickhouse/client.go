package clickhouse

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type Client struct {
	sync.RWMutex
	conn driver.Conn

	dbName string
}

type ClientConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

func NewClickHouseClient(cfg ClientConfig) (*Client, error) {
	var client Client

	if err := client.Configure(cfg); err != nil {
		return nil, err
	}

	return &client, nil
}

func (c *Client) Configure(cfg ClientConfig) error {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.User,
			Password: cfg.Password,
		},
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "gitlab-exporter", Version: "v0.0.0+unknown"},
			},
		},
	})
	if err != nil {
		return err
	}

	c.Lock()
	c.conn = conn
	c.dbName = cfg.Database
	c.Unlock()
	return nil
}

func (c *Client) CheckReadiness(ctx context.Context) error {
	if err := c.conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			return fmt.Errorf("clickhouse exception: [%d] %s", exception.Code, exception.Message)
		} else {
			return fmt.Errorf("error pinging clickhouse: %w", err)
		}
	}
	return nil
}

func WithParameters(ctx context.Context, params map[string]string) context.Context {
	return clickhouse.Context(ctx, clickhouse.WithParameters(params))
}

func (c *Client) Exec(ctx context.Context, query string, args ...any) error {
	c.RLock()
	defer c.RUnlock()
	return c.conn.Exec(ctx, query, args...)
}

func (c *Client) Select(ctx context.Context, dest any, query string, args ...any) error {
	c.RLock()
	defer c.RUnlock()
	return c.conn.Select(ctx, dest, query, args...)
}

func (c *Client) PrepareBatch(ctx context.Context, query string) (driver.Batch, error) {
	c.RLock()
	defer c.RUnlock()
	return c.conn.PrepareBatch(ctx, query)
}

func (c *Client) CreateTables(ctx context.Context) error {
	return createTables(ctx, c.dbName, c)
}

func (c *Client) QueryProjectPipelinesLatestUpdate(ctx context.Context, projectID int64) (map[int64]time.Time, error) {
	const (
		msPerSecond float64 = 1000
	)

	var results []struct {
		PipelineID   int64   `ch:"id"`
		LatestUpdate float64 `ch:"latest_update"`
	}

	query := fmt.Sprintf(`
        SELECT id, max(updated_at) AS latest_update
        FROM %s.pipelines
        WHERE project_id = %d
        GROUP BY id
    `, c.dbName, projectID)

	if err := c.Select(ctx, &results, query); err != nil {
		return nil, err
	}

	m := map[int64]time.Time{}
	for _, r := range results {
		m[r.PipelineID] = time.UnixMilli(int64(r.LatestUpdate * msPerSecond)).UTC()
	}

	return m, nil
}
