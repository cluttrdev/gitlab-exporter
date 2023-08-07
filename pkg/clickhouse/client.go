package clickhouseclient

import (
	"context"
	"fmt"
	"log"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type Client struct {
	Conn driver.Conn
}

type ClientConfig struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

func NewClickHouseClient(cfg ClientConfig) (*Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
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
					{Name: "gitlab-clickhouse-exporter", Version: "0.1.0"},
				},
			},
		})
	)

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Printf("Exception: [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}

	return &Client{
		Conn: conn,
	}, nil
}

func (c *Client) CreateDatabase(ctx context.Context) error {
	return c.Conn.Exec(ctx, `CREATE DATABASE IF NOT EXISTS gitlab_ci`)
}

func (c *Client) CreateTables(ctx context.Context) error {
	return createTables(ctx, c)
}

func (c *Client) PrepareBatch(ctx context.Context, query string) (driver.Batch, error) {
	return c.Conn.PrepareBatch(ctx, query)
}
