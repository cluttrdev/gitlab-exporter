package clickhouse

import (
	"context"
	"fmt"

	"golang.org/x/sync/semaphore"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type Client struct {
	conn   driver.Conn
	dbName string

	sem *semaphore.Weighted
}

type ClientConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

func NewClient(conn driver.Conn, database string) *Client {
	return &Client{
		conn:   conn,
		dbName: database,
	}
}

func (c *Client) SetMaxConcurrentQueries(n int64) {
	if n > 0 {
		c.sem = semaphore.NewWeighted(n)
	} else {
		c.sem = nil
	}
}

func (c *Client) acquire(ctx context.Context, n int64) error {
	if c.sem == nil {
		return nil
	}
	return c.sem.Acquire(ctx, n)
}

func (c *Client) release(n int64) {
	if c.sem == nil {
		return
	}
	c.sem.Release(n)
}

func ClientOptions(cfg ClientConfig) clickhouse.Options {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	return clickhouse.Options{
		Protocol: clickhouse.Native,
		Addr:     []string{addr},
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
				{Name: "gitlab-exporter-clickhouse-recorder", Version: "v0.0.0+unknown"},
			},
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	}
}

func Connect(options *clickhouse.Options) (driver.Conn, error) {
	if options.Settings == nil {
		options.Settings = clickhouse.Settings{
			"connect_timeout": 30,
		}
	}

	return clickhouse.Open(options)
}

func (c *Client) Ping(ctx context.Context) error {
	return c.conn.Ping(ctx)
}

func WithParameters(ctx context.Context, params map[string]string) context.Context {
	return clickhouse.Context(ctx, clickhouse.WithParameters(params))
}

func (c *Client) Exec(ctx context.Context, query string, args ...any) error {
	if err := c.acquire(ctx, 1); err != nil {
		return err
	}
	defer c.release(1)
	return c.conn.Exec(ctx, query, args...)
}

func (c *Client) Select(ctx context.Context, dest any, query string, args ...any) error {
	if err := c.acquire(ctx, 1); err != nil {
		return err
	}
	defer c.release(1)
	return c.conn.Select(ctx, dest, query, args...)
}

func (c *Client) PrepareBatch(ctx context.Context, query string) (driver.Batch, error) {
	if err := c.acquire(ctx, 1); err != nil {
		return nil, err
	}
	defer c.release(1)
	return c.conn.PrepareBatch(ctx, query)
}
