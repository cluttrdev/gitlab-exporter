package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "strconv"

    clickhouse "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/clickhouse"
    config "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/config"
    gitlab "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/gitlab"
)

func main() {
    flag.Parse()
    if flag.NArg() != 2 {
        fmt.Println("usage: gitlab-clickhouse-exporter PROJECT_ID PIPELINE_ID")
        return
    }

    projectID, err := strconv.ParseInt(flag.Arg(0), 10, 64)
    if err != nil {
        log.Fatal(err)
    }
    pipelineID, err := strconv.ParseInt(flag.Arg(1), 10, 64)
    if err != nil {
        log.Fatal(err)
    }

    cfg, err := config.LoadEnv()
    if err != nil {
        log.Fatal(err)
    }

    gl, err := gitlab.NewGitLabClient(gitlab.ClientConfig{
        URL: cfg.GitLab.URL,
        Token: cfg.GitLab.Token,
    })
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    ph, err := gl.GetPipelineHierarchy(ctx, projectID, pipelineID)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf(
        "Pipeline %d has %d jobs, %d sections, and %d bridges",
        ph.Pipeline.ID, len(ph.Jobs), len(ph.Sections), len(ph.Bridges))

    ch, err := clickhouse.NewClickHouseClient(clickhouse.ClientConfig{
        Host: cfg.ClickHouse.Host,
        Port: cfg.ClickHouse.Port,
        Database: cfg.ClickHouse.Database,
        User: cfg.ClickHouse.User,
        Password: cfg.ClickHouse.Password,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    if err = ch.CreateDatabase(ctx); err != nil {
        log.Fatal(err)
    }
    if err = ch.CreateTables(ctx); err != nil {
        log.Fatal(err)
    }

    if err = clickhouse.InsertPipelineHierarchy(ctx, ph, ch); err != nil {
        log.Fatal(err)
    }

    traces := ph.GetAllTraces()
    if err = clickhouse.InsertTraces(ctx, traces, ch); err != nil {
        log.Fatal(err)
    }
}
