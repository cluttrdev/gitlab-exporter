package migrations_test

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	ch "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/golang-migrate/migrate/v4"

	"go.cluttr.dev/gitlab-exporter/recorders/clickhouse/internal/clickhouse"
	"go.cluttr.dev/gitlab-exporter/recorders/clickhouse/tests/integration/conftest"
)

const testSet string = "migrations"

func TestMain(m *testing.M) {
	// Check if Docker is available before attempting to create test environment
	if !conftest.IsDockerAvailable() {
		fmt.Fprintln(os.Stderr, "Skipping integration tests: Docker not available")
		os.Exit(0)
	}

	env, err := conftest.CreateClickHouseTestEnvironment(testSet)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = env.Container.Terminate(context.Background())
	}()

	conftest.SetTestEnvironment(testSet, env)

	if err := conftest.CreateDatabase(testSet); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func testClient(t *testing.T) *client {
	client, err := newClient(testSet, "../migrations")
	if err != nil {
		t.Fatalf("failed to create test client: %v", err)
	}
	return client
}

type client struct {
	Conn      ch.Conn
	Migration *migrate.Migrate
}

func newClient(testSet string, migrationsPath string) (*client, error) {
	opts, err := conftest.GetTestClientOptions(testSet)
	if err != nil {
		return nil, err
	}
	opts.MaxOpenConns = 1

	conn, err := ch.Open(&opts)
	if err != nil {
		return nil, err
	}

	migration, err := newMigration(testSet, migrationsPath)
	if err != nil {
		return nil, err
	}

	return &client{Conn: conn, Migration: migration}, nil
}

func newMigration(testSet string, migrationsPath string) (*migrate.Migrate, error) {
	env, err := conftest.GetTestEnvironment(testSet)
	if err != nil {
		return nil, err
	}

	root, err := filepath.Abs(migrationsPath)
	if err != nil {
		return nil, err
	}

	fsys := OSPathFS{
		Path: filepath.Clean(root),
	}

	opts := clickhouse.MigrationOptions{
		ClientConfig: clickhouse.ClientConfig{
			Host:     env.Host,
			Port:     fmt.Sprint(env.Port),
			Database: env.Database,
			User:     env.Username,
			Password: env.Password,
		},

		FileSystem: &fsys,
		Path:       "",
	}

	return clickhouse.NewMigration(opts)
}

type OSPathFS struct {
	Path string
}

func (fsys *OSPathFS) Open(name string) (fs.File, error) {
	return os.Open(filepath.Join(fsys.Path, name))
}
