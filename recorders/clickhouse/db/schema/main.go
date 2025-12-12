package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/docker/go-units"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/testcontainers/testcontainers-go"
	tcexec "github.com/testcontainers/testcontainers-go/exec"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
)

const (
	clickhouseServerVersion string = "latest"
	clickhouseUsername      string = "default"
	clickhousePassword      string = "ClickHouse"
	clickhouseDatabase      string = "gitlab_ci"
)

var tables = []string{
	"bridges",
	"coverage_reports",
	"coverage_packages",
	"coverage_classes",
	"coverage_methods",
	"deployments",
	"issues",
	"jobs",
	"mergerequest_noteevents",
	"mergerequests",
	"metrics",
	"pipelines",
	"projects",
	"sections",
	"testcases",
	"testreports",
	"testsuites",
	"traces",
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	cont, err := createContainer(ctx)
	if err != nil {
		return fmt.Errorf("create container: %w", err)
	}
	defer func() {
		_ = cont.Terminate(context.Background())
	}()

	p, _ := cont.MappedPort(ctx, "9000")
	conn, err := getConnection("127.0.0.1", p.Port(), "default")
	if err != nil {
		return fmt.Errorf("get connection: %w", err)
	}

	if err := conn.Exec(ctx, fmt.Sprintf("CREATE DATABASE `%s`", clickhouseDatabase)); err != nil {
		return fmt.Errorf("create database: %w", err)
	}

	if err := runMigrations(ctx, "127.0.0.1", p.Port()); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	var buf bytes.Buffer
	for _, table := range tables {
		buf.Reset()

		// -- table
		buf.WriteString(fmt.Sprintf("-- %s\n", table))
		query, err := getCreateTableQuery(ctx, conn, table)
		if err != nil {
			return err
		}
		fquery, err := formatQuery(ctx, cont, query)
		if err != nil {
			return fmt.Errorf("format query: %w", err)
		}
		_, _ = buf.Write(fquery)

		// -- table_in
		buf.WriteString(fmt.Sprintf("-- %s_in\n", table))
		query_in := fmt.Sprintf(`
            CREATE TABLE IF NOT EXISTS %s.%s_in AS %s.%s ENGINE = Null;
        `, clickhouseDatabase, table, clickhouseDatabase, table)
		fquery_in, err := formatQuery(ctx, cont, string(query_in))
		if err != nil {
			return fmt.Errorf("format query: %w", err)
		}
		_, _ = buf.Write(fquery_in)

		// -- table_mv
		buf.WriteString(fmt.Sprintf("-- %s_mv\n", table))
		querySelect_mv, err := getAsSelectQuery(ctx, conn, table+"_mv")
		if err != nil {
			return err
		}
		query_mv := fmt.Sprintf(`
            CREATE MATERIALIZED VIEW IF NOT EXISTS %s.%s_mv TO %s.%s
            AS %s
        `, clickhouseDatabase, table, clickhouseDatabase, table, querySelect_mv)
		fquery_mv, err := formatQuery(ctx, cont, query_mv)
		if err != nil {
			return fmt.Errorf("format query: %w", err)
		}
		_, _ = buf.Write(fquery_mv)

		fmt.Println(buf.String())
	}

	// additional stuff for traces
	buf.Reset()

	// -- traces_trace_id_ts
	buf.WriteString("-- traces_trace_id_ts\n")
	query, err := getCreateTableQuery(ctx, conn, "traces_trace_id_ts")
	if err != nil {
		return err
	}
	fquery, err := formatQuery(ctx, cont, query)
	if err != nil {
		return fmt.Errorf("format query: %w", err)
	}
	_, _ = buf.Write(fquery)

	// -- traces_trace_id_ts_mv
	buf.WriteString("-- traces_trace_id_ts_mv\n")
	query_mv, err := getCreateTableQuery(ctx, conn, "traces_trace_id_ts_mv")
	if err != nil {
		return err
	}
	fquery_mv, err := formatQuery(ctx, cont, query_mv)
	if err != nil {
		return fmt.Errorf("format query: %w", err)
	}
	_, _ = buf.Write(fquery_mv)

	// --trace_view
	buf.WriteString("-- trace_view\n")
	query_v, err := getCreateTableQuery(ctx, conn, "trace_view")
	if err != nil {
		return err
	}
	fquery_v, err := formatQuery(ctx, cont, query_v)
	if err != nil {
		return fmt.Errorf("format query: %w", err)
	}
	_, _ = buf.Write(fquery_v)

	fmt.Println(buf.String())

	return nil
}

func getCreateTableQuery(ctx context.Context, conn driver.Conn, table string) (string, error) {
	var result []struct {
		CreateTableQuery string `ch:"create_table_query"`
	}
	q := fmt.Sprintf(`SELECT create_table_query FROM system.tables WHERE name = '%s'`, table)
	if err := conn.Select(ctx, &result, q); err != nil {
		return "", fmt.Errorf("get create table query: %w", err)
	}

	if len(result) == 0 {
		return "", fmt.Errorf("empty result")
	}

	return result[0].CreateTableQuery, nil
}

func getAsSelectQuery(ctx context.Context, conn driver.Conn, view string) (string, error) {
	var result []struct {
		AsSelect string `ch:"as_select"`
	}
	q := fmt.Sprintf(`SELECT as_select FROM system.tables WHERE name = '%s'`, view)
	if err := conn.Select(ctx, &result, q); err != nil {
		return "", fmt.Errorf("get select query: %w", err)
	}

	if len(result) == 0 {
		return "", fmt.Errorf("empty result")
	}

	return result[0].AsSelect, nil
}

func formatQuery(ctx context.Context, cont testcontainers.Container, query string) ([]byte, error) {
	_, r, err := cont.Exec(ctx, []string{
		"clickhouse-format",
		"--comments",
		"--multiquery",
		"--query", query,
	}, tcexec.Multiplexed())
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func createContainer(ctx context.Context) (testcontainers.Container, error) {
	provider, err := testcontainers.ProviderDocker.GetProvider()
	if err != nil {
		return nil, fmt.Errorf("get provider: %w", err)
	}
	if err := provider.Health(ctx); err != nil {
		return nil, fmt.Errorf("check health: %w", err)
	}

	expected := []*units.Ulimit{
		{
			Name: "nofile",
			Hard: 262144,
			Soft: 262144,
		},
	}
	req := testcontainers.ContainerRequest{
		Image: fmt.Sprintf("clickhouse/clickhouse-server:%s", clickhouseServerVersion),
		Name:  fmt.Sprintf("clickhouse-go-%d", time.Now().UnixNano()),
		Env: map[string]string{
			"CLICKHOUSE_USER":     clickhouseUsername,
			"CLICKHOUSE_PASSWORD": clickhousePassword,
		},
		ExposedPorts: []string{"9000/tcp"},
		WaitingFor: tcwait.ForAll(
			tcwait.ForSQL("9000/tcp", "clickhouse", func(host string, port nat.Port) string {
				return fmt.Sprintf("clickhouse://%s:%s@%s:%s", clickhouseUsername, clickhousePassword, host, port.Port())
			}),
		).WithDeadline(time.Second * time.Duration(30)),
		Resources: container.Resources{
			Ulimits: expected,
		},
	}

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func getConnection(host string, port string, database string) (driver.Conn, error) {
	if database == "" {
		database = "default"
	}
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr:     []string{fmt.Sprintf("%s:%s", host, port)},
		Settings: clickhouse.Settings{},
		Auth: clickhouse.Auth{
			Database: database,
			Username: clickhouseUsername,
			Password: clickhousePassword,
		},
		DialTimeout: time.Duration(10) * time.Second,
	})
	return conn, err
}

type osPathFS struct {
	Path string
}

func (fsys *osPathFS) Open(name string) (fs.File, error) {
	return os.Open(filepath.Join(fsys.Path, name))
}

func createMigration(fsys fs.FS, path string, host string, port string, database string) (*migrate.Migrate, error) {
	drv, err := iofs.New(fsys, path)
	if err != nil {
		return nil, fmt.Errorf("create source driver: %w", err)
	}

	q := url.Values{}
	q.Set("x-multi-statement", "true")
	q.Set("x-migrations-table", "schema_migrations")
	q.Set("x-migrations-table-engine", "MergeTree")
	// q.Set("x-cluster-name", "")

	dsn := url.URL{
		Scheme:   "clickhouse",
		Host:     fmt.Sprintf("%s:%s", host, port),
		Path:     database,
		User:     url.UserPassword(clickhouseUsername, clickhousePassword),
		RawQuery: q.Encode(),
	}

	return migrate.NewWithSourceInstance("iofs", drv, dsn.String())
}

func runMigrations(ctx context.Context, host string, port string) error {
	root, err := filepath.Abs("db/migrations")
	if err != nil {
		return err
	}

	fsys := osPathFS{
		Path: filepath.Clean(root),
	}

	m, err := createMigration(&fsys, "", host, port, clickhouseDatabase)
	if err != nil {
		return fmt.Errorf("create migration: %w", err)
	}

	return m.Up()
}
