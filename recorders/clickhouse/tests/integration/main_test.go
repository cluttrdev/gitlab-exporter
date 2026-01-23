package integration_tests

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"go.cluttr.dev/gitlab-exporter/recorders/clickhouse/internal/clickhouse"
	"go.cluttr.dev/gitlab-exporter/recorders/clickhouse/tests/integration/conftest"
)

const testSet string = "native"

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

	if err := RunSchemaMigrations(testSet); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func GetTestClient(testSet string) (*clickhouse.Client, error) {
	opts, err := conftest.GetTestClientOptions(testSet)
	if err != nil {
		return nil, err
	}
	opts.MaxOpenConns = 1

	conn, err := clickhouse.Connect(&opts)
	if err != nil {
		return nil, err
	}

	return clickhouse.NewClient(conn, opts.Auth.Database), nil
}

type OSPathFS struct {
	Path string
}

func (fsys *OSPathFS) Open(name string) (fs.File, error) {
	return os.Open(filepath.Join(fsys.Path, name))
}

func RunSchemaMigrations(testSet string) error {
	env, err := conftest.GetTestEnvironment(testSet)
	if err != nil {
		return err
	}

	root, err := filepath.Abs("../../db/migrations")
	if err != nil {
		return err
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

	if err := clickhouse.MigrateUp(opts); err != nil {
		return err
	}
	if err := clickhouse.MigrateDown(opts); err != nil {
		return err
	}
	return clickhouse.MigrateUp(opts)
}
