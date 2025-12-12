package main

import (
	"embed"
	"fmt"
	"os"

	"go.cluttr.dev/gitlab-exporter-clickhouse-recorder/internal/cmd"
)

var version string

//go:embed db/migrations/*.sql
var migrationsFS embed.FS

func main() {
	cmd.Version = version

	cmd.MigrationsFileSystem = migrationsFS
	cmd.MigrationsPath = "db/migrations"

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
