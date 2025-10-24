package integration_tests

// DISCLAIMER:
//
// The code in this file was taken from
//
//     https://github.com/ClickHouse/clickhouse-go/blob/main/tests/utils.go
//
// and adjusted to fit the needs of this project.
//
// The original code is released under the Apache License, Version 2.0;
// you may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/docker/go-units"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	ch "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

var testUUID = uuid.NewString()[0:12]
var testTimestamp = time.Now().UnixMilli()

const defaultClickHouseVersion = "latest"

func GetClickHouseTestVersion() string {
	return GetEnv("CLICKHOUSE_VERSION", defaultClickHouseVersion)
}

type ClickHouseTestEnvironment struct {
	Port        int
	Host        string
	Username    string
	Password    string
	Database    string
	ContainerIP string
	Container   testcontainers.Container `json:"-"`
}

func CreateClickHouseTestEnvironment(testSet string) (ClickHouseTestEnvironment, error) {
	// create a ClickHouse Container
	ctx := context.Background()
	// attempt use docker for CI
	provider, err := testcontainers.ProviderDocker.GetProvider()
	if err != nil {
		fmt.Printf("Docker is not running and no clickhouse connections details were provided. Skipping tests: %s\n", err)
		os.Exit(0)
	}
	err = provider.Health(ctx)
	if err != nil {
		fmt.Printf("Docker provider health check failed, skipping integration tests: %v\n", err)
		os.Exit(0)
	}
	fmt.Printf("Using Docker for IT tests\n")

	expected := []*units.Ulimit{
		{
			Name: "nofile",
			Hard: 262144,
			Soft: 262144,
		},
	}
	req := testcontainers.ContainerRequest{
		Image: fmt.Sprintf("clickhouse/clickhouse-server:%s", GetClickHouseTestVersion()),
		Name:  fmt.Sprintf("clickhouse-go-%s-%d", strings.ToLower(testSet), time.Now().UnixNano()),
		Env: map[string]string{
			"CLICKHOUSE_USER":     "default",
			"CLICKHOUSE_PASSWORD": "ClickHouse",
		},
		ExposedPorts: []string{"9000/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForSQL("9000/tcp", "clickhouse", func(host string, port nat.Port) string {
				return fmt.Sprintf("clickhouse://default:ClickHouse@%s:%s", host, port.Port())
			}),
		).WithDeadline(time.Second * time.Duration(120)),
		Resources: container.Resources{
			Ulimits: expected,
		},
	}
	clickhouseContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return ClickHouseTestEnvironment{}, err
	}
	p, _ := clickhouseContainer.MappedPort(ctx, "9000")
	ip, _ := clickhouseContainer.ContainerIP(ctx)
	testEnv := ClickHouseTestEnvironment{
		Port:        p.Int(),
		Host:        "127.0.0.1",
		Username:    "default",
		Password:    "ClickHouse",
		Container:   clickhouseContainer,
		ContainerIP: ip,
		Database:    GetEnv("CLICKHOUSE_DATABASE", getDatabaseName(testSet)),
	}
	return testEnv, nil
}

func SetTestEnvironment(testSet string, environment ClickHouseTestEnvironment) {
	bytes, err := json.Marshal(environment)
	if err != nil {
		panic(err)
	}
	os.Setenv(fmt.Sprintf("CLICKHOUSE_%s_ENV", strings.ToUpper(testSet)), string(bytes))
}

func GetTestEnvironment(testSet string) (ClickHouseTestEnvironment, error) {
	sEnv := os.Getenv(fmt.Sprintf("CLICKHOUSE_%s_ENV", strings.ToUpper(testSet)))
	if sEnv == "" {
		return ClickHouseTestEnvironment{}, errors.New("unable to find environment")
	}
	var env ClickHouseTestEnvironment
	if err := json.Unmarshal([]byte(sEnv), &env); err != nil {
		return ClickHouseTestEnvironment{}, err
	}
	return env, nil
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func ClientOptionsFromEnv(env ClickHouseTestEnvironment, settings ch.Settings) ch.Options {
	timeout, err := strconv.Atoi(GetEnv("CLICKHOUSE_DIAL_TIMEOUT", "10"))
	if err != nil {
		timeout = 10
	}

	return ch.Options{
		Addr:     []string{fmt.Sprintf("%s:%d", env.Host, env.Port)},
		Settings: settings,
		Auth: ch.Auth{
			Database: env.Database,
			Username: env.Username,
			Password: env.Password,
		},
		DialTimeout: time.Duration(timeout) * time.Second,
		Compression: &ch.Compression{
			Method: ch.CompressionLZ4,
		},
	}
}

func getDatabaseName(testSet string) string {
	return fmt.Sprintf("clickhouse-go-%s-%s-%d", testSet, testUUID, testTimestamp)
}

func getConnection(env ClickHouseTestEnvironment, database string, settings ch.Settings) (driver.Conn, error) {
	var err error

	port := env.Port
	if settings == nil {
		settings = ch.Settings{}
	}
	settings["insert_quorum"], err = strconv.Atoi(GetEnv("CLICKHOUSE_QUORUM_INSERT", "1"))
	settings["insert_quorum_parallel"] = 0
	settings["select_sequential_consistency"] = 1
	if err != nil {
		return nil, err
	}

	timeout, err := strconv.Atoi(GetEnv("CLICKHOUSE_DIAL_TIMEOUT", "10"))
	if err != nil {
		return nil, err
	}
	conn, err := ch.Open(&ch.Options{
		Addr:     []string{fmt.Sprintf("%s:%d", env.Host, port)},
		Settings: settings,
		Auth: ch.Auth{
			Database: database,
			Username: env.Username,
			Password: env.Password,
		},
		DialTimeout: time.Duration(timeout) * time.Second,
	})
	return conn, err
}

func CreateDatabase(testSet string) error {
	env, err := GetTestEnvironment(testSet)
	if err != nil {
		return err
	}
	conn, err := getConnection(env, "default", nil)
	if err != nil {
		return err
	}
	return conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE `%s`", env.Database))
}
