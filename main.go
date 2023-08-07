package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/cmd"
	config "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/config"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	if err = cmd.Run(ctx, *cfg, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
