package clickhouseclient

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type Client struct {
	sync.RWMutex
	conn driver.Conn
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

	if err := client.configure(cfg); err != nil {
		return nil, err
	}

	return &client, nil
}

func (c *Client) configure(cfg ClientConfig) error {
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
				{Name: "gitlab-clickhouse-exporter", Version: "v0.0.0+unknown"},
			},
		},
	})
	if err != nil {
		return err
	}

	ctx := context.Background()
	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Printf("Exception: [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return err
	}

	c.Lock()
	c.conn = conn
	c.Unlock()
	return nil
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

func (c *Client) CreateDatabase(ctx context.Context) error {
	return c.Exec(ctx, `CREATE DATABASE IF NOT EXISTS gitlab_ci`)
}

func (c *Client) CreateTables(ctx context.Context) error {
	return createTables(ctx, c)
}
