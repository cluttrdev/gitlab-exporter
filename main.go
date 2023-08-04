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
    ctx, cancel := context.WithCancel(ctx)

    cfg, err := config.LoadEnv()
    if err != nil {
        log.Fatal(err)
    }

    defer func() {
        cancel()
    }()

    if err = cmd.Run(ctx, *cfg, os.Stdout); err != nil {
        fmt.Fprintf(os.Stderr, "%s\n", err)
        os.Exit(1)
    }
}
