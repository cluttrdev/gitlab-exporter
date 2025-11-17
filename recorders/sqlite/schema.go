package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"time"

	"go.cluttr.dev/migrate"
	"go.cluttr.dev/migrate/source/iofs"

	sqlite_driver "go.cluttr.dev/gitlab-exporter/recorders/sqlite/db/driver"
)

//go:embed db/migrations/*.sql
var migrationsFS embed.FS

// migrationsPath is the path to the migrations directory in the embedded filesystem
const migrationsPath = "db/migrations"

func RunMigrations(ctx context.Context, db *sql.DB, name string) error {
	srcDrv, err := iofs.New(migrationsFS, migrationsPath)
	if err != nil {
		return fmt.Errorf("create migration source driver: %w", err)
	}
	dbDrv, err := sqlite_driver.WithInstance(db, &sqlite_driver.Config{
		MigrationsTable: sqlite_driver.DefaultMigrationsTable,
		DatabaseName:    name,
		NoTxWrap:        false,
	})
	if err != nil {
		return fmt.Errorf("create migration database driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", srcDrv, "sqlite", dbDrv)
	if err != nil {
		return fmt.Errorf("create migration: %w", err)
	}

	if err := m.Up(); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}

// configure sets up database pragmas and settings
func (r *Recorder) configure(ctx context.Context) error {
	r.db.SetMaxOpenConns(1) // only single writer ever
	r.db.SetConnMaxIdleTime(1 * time.Minute)

	pragmas := []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA synchronous=NORMAL",
		"PRAGMA temp_store=MEMORY",
		"PRAGMA busy_timeout=5000",
		"PRAGMA cache_size=-64000", // 64MB cache
	}

	if !r.settings.WALMode {
		pragmas[0] = "PRAGMA journal_mode=DELETE"
	}

	for _, pragma := range pragmas {
		if _, err := r.db.ExecContext(ctx, pragma); err != nil {
			return fmt.Errorf("execute %s: %w", pragma, err)
		}
	}

	return nil
}
