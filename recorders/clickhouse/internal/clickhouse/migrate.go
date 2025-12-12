package clickhouse

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"net/url"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type MigrationOptions struct {
	ClientConfig

	FileSystem fs.FS
	Path       string
}

var (
	ErrMigrateNoChange   = migrate.ErrNoChange
	ErrMigrateNilVersion = migrate.ErrNilVersion
)

const migrationsTable string = "schema_migrations"

func GetSchemaVersion(c *Client, ctx context.Context) (uint, bool, error) {
	if err := c.acquire(ctx, 1); err != nil {
		return 0, false, err
	}
	defer c.release(1)

	var (
		version int64
		dirty   uint8
		query   = "SELECT version, dirty FROM `" + migrationsTable + "` ORDER BY sequence DESC LIMIT 1"
	)
	if err := c.conn.QueryRow(ctx, query).Scan(&version, &dirty); err != nil {
		if err == sql.ErrNoRows {
			return 0, false, ErrMigrateNilVersion
		}
		return 0, false, err
	}

	if version < 0 {
		return 0, false, fmt.Errorf("expected version >= 0, got %d", version)
	}
	return uint(version), dirty == 1, nil
}

func GetLatestMigrationVersion(fsys fs.FS, path string) (uint, error) {
	entries, err := fs.ReadDir(fsys, path)
	if err != nil {
		return 0, fmt.Errorf("error reading migrations: %w", err)
	}

	var version uint = 0
	nilVersion := true
	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		m, err := source.Parse(e.Name())
		if err != nil {
			continue
		}

		nilVersion = false
		if m.Version > version {
			version = m.Version
		}
	}
	if nilVersion {
		return 0, ErrMigrateNilVersion
	}
	return version, nil
}

func NewMigration(opts MigrationOptions) (*migrate.Migrate, error) {
	if opts.FileSystem == nil {
		return nil, errors.New("missing migrations file system")
	}
	drv, err := iofs.New(opts.FileSystem, opts.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to create source driver: %w", err)
	}

	q := url.Values{}
	q.Set("x-multi-statement", "true")
	q.Set("x-migrations-table", migrationsTable)
	q.Set("x-migrations-table-engine", "MergeTree")
	// q.Set("x-cluster-name", "")

	dsn := url.URL{
		Scheme:   "clickhouse",
		Host:     fmt.Sprintf("%s:%s", opts.ClientConfig.Host, opts.ClientConfig.Port),
		Path:     opts.ClientConfig.Database,
		User:     url.UserPassword(opts.ClientConfig.User, opts.ClientConfig.Password),
		RawQuery: q.Encode(),
	}

	return migrate.NewWithSourceInstance("iofs", drv, dsn.String())
}

func MigrateUp(opts MigrationOptions) error {
	m, err := NewMigration(opts)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil {
		return fmt.Errorf("failed to apply up migrations: %w", err)
	}
	return nil
}

func MigrateDown(opts MigrationOptions) error {
	m, err := NewMigration(opts)
	if err != nil {
		return fmt.Errorf("create migration instance: %w", err)
	}

	if err := m.Down(); err != nil {
		return fmt.Errorf("apply down migrations: %w", err)
	}
	return nil
}
