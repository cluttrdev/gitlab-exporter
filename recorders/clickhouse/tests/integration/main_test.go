package integration_tests

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	ch "github.com/ClickHouse/clickhouse-go/v2"

	"go.cluttr.dev/gitlab-exporter-clickhouse-recorder/internal/clickhouse"
)

const testSet string = "native"

func TestMain(m *testing.M) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		fmt.Fprintln(os.Stderr, "Skipping integration tests")
		os.Exit(0)
	}

	env, err := CreateClickHouseTestEnvironment(testSet)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = env.Container.Terminate(context.Background())
	}()

	SetTestEnvironment(testSet, env)

	if err := CreateDatabase(testSet); err != nil {
		panic(err)
	}

	if err := RunSchemaMigrations(testSet); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func GetTestClient(testSet string) (*clickhouse.Client, error) {
	te, err := GetTestEnvironment(testSet)
	if err != nil {
		return nil, err
	}

	opts := ClientOptionsFromEnv(te, ch.Settings{})
	opts.MaxOpenConns = 1

	conn, err := clickhouse.Connect(&opts)
	if err != nil {
		return nil, err
	}

	return clickhouse.NewClient(conn, te.Database), nil
}

type OSPathFS struct {
	Path string
}

func (fsys *OSPathFS) Open(name string) (fs.File, error) {
	return os.Open(filepath.Join(fsys.Path, name))
}

func RunSchemaMigrations(testSet string) error {
	env, err := GetTestEnvironment(testSet)
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
