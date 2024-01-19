package controller

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpc_client "github.com/cluttrdev/gitlab-exporter/grpc/client"

	"github.com/cluttrdev/gitlab-exporter/pkg/clickhouse"
	"github.com/cluttrdev/gitlab-exporter/pkg/config"
	"github.com/cluttrdev/gitlab-exporter/pkg/gitlab"

	"github.com/cluttrdev/gitlab-exporter/internal/datastore"
	"github.com/cluttrdev/gitlab-exporter/internal/worker"
)

type Controller struct {
	config    config.Config
	GitLab    gitlab.Client
	Exporter  *Exporter
	DataStore datastore.DataStore

	workers []worker.Worker
}

func NewController(cfg config.Config) (*Controller, error) {
	c := &Controller{
		config: cfg,
	}

	if err := c.configure(cfg); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Controller) configure(cfg config.Config) error {
	if err := c.configureGitLabClient(cfg.GitLab); err != nil {
		return err
	}

	if err := c.configureClickHouseDataStore(cfg.ClickHouse); err != nil {
		return err
	}

	if err := c.configureExporter(cfg.Endpoints); err != nil {
		return err
	}

	if err := c.configureWorkers(cfg); err != nil {
		return err
	}

	return nil
}

func (c *Controller) configureGitLabClient(cfg config.GitLab) error {
	return c.GitLab.Configure(gitlab.ClientConfig{
		URL:   cfg.Api.URL,
		Token: cfg.Api.Token,

		RateLimit: cfg.Client.Rate.Limit,
	})
}

func (c *Controller) configureClickHouseDataStore(cfg config.ClickHouse) error {
	var client clickhouse.Client

	conf := clickhouse.ClientConfig{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Database: cfg.Database,
		User:     cfg.User,
		Password: cfg.Password,
	}
	if err := client.Configure(conf); err != nil {
		return err
	}

	c.DataStore = clickhouse.NewClickHouseDataStore(&client)
	return nil
}

func (c *Controller) configureExporter(cfg []config.Endpoint) error {
	endpoints := make([]grpc_client.EndpointConfig, 0, len(cfg))
	for _, cc := range cfg {
		endpoints = append(endpoints, grpc_client.EndpointConfig{
			Address: cc.Address,
			Options: []grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			},
		})
	}
	exp, err := NewExporter(endpoints)
	if err != nil {
		return err
	}
	c.Exporter = exp
	return nil
}

func (c *Controller) configureWorkers(cfg config.Config) error {
	workers := []worker.Worker{}

	for _, prj := range cfg.Projects {
		if prj.CatchUp.Enabled {
			workers = append(workers, worker.NewCatchUpProjectWorker(prj, &c.GitLab, c.DataStore))
		}
		workers = append(workers, worker.NewExportProjectWorker(prj, &c.GitLab, c.DataStore))
	}

	c.workers = workers

	return nil
}

func (c *Controller) init(ctx context.Context) error {
	if err := c.DataStore.Initialize(ctx); err != nil {
		return err
	}

	return nil
}

func (c *Controller) CheckReadiness(ctx context.Context) error {
	if err := c.GitLab.CheckReadiness(ctx); err != nil {
		return err
	}

	if err := c.DataStore.CheckReadiness(ctx); err != nil {
		return err
	}

	return nil
}

func (c *Controller) Run(ctx context.Context) error {
	if err := c.waitForReady(ctx); err != nil {
		return fmt.Errorf("error getting ready: %w", err)
	}

	if err := c.init(ctx); err != nil {
		return err
	}

	log.Println("Starting workers...")
	for _, w := range c.workers {
		go w.Start(ctx)
	}

	<-ctx.Done()

	log.Println("Stopping workers...")
	for _, w := range c.workers {
		w.Stop()
		<-w.Done()
	}

	return nil
}

func (c *Controller) waitForReady(ctx context.Context) error {
	var (
		maxTries         int     = 5
		backoffBaseSec   float64 = 1.0
		backoffJitterSec float64 = 1.0
	)

	ticker := time.NewTicker(time.Second)

	var err error
	for i := 0; i < maxTries; i++ {
		if err = c.CheckReadiness(ctx); err == nil {
			return nil
		}

		log.Println(fmt.Errorf("Readiness check failed: %w", err))
		delaySec := backoffBaseSec*math.Pow(2, float64(i)) + backoffJitterSec*rand.Float64()
		delay := time.Duration(delaySec) * time.Second
		ticker.Reset(delay)

		select {
		case <-ticker.C:
		case <-ctx.Done():
			ticker.Stop()
			return context.Canceled
		}
	}

	return errors.New("Failed to get ready")
}
