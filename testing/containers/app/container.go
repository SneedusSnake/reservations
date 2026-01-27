package app

import (
	"context"
	"time"

	"github.com/SneedusSnake/Reservations/testing/utils"
	"github.com/testcontainers/testcontainers-go"
)

const CLOCK_CACHE_PATH = "/tmp/clock_go"

func Start(ctx context.Context, network string, mysqlConnection string, logs ...testcontainers.LogConsumer) (testcontainers.Container, error) {
	persistenceDriver := "memory"
	if mysqlConnection != "" {
		persistenceDriver = "mysql"
	}
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context: utils.TestsRootDir() + "/..",
			Dockerfile: "./build/Docker/Dockerfile",
			PrintBuildLog: true,
		},
		Env: map[string]string{
			"TELEGRAM_API_HOST": "http://telegram-api:8080",
			"TELEGRAM_API_TOKEN": "1234567",
			"CLOCK_DRIVER": "cache",
			"CACHE_CLOCK_PATH": CLOCK_CACHE_PATH,
			"MYSQL_CONNECTION": mysqlConnection,
			"PERSISTENCE_DRIVER": persistenceDriver,
			},
		Networks: []string{network},
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Opts: []testcontainers.LogProductionOption{testcontainers.WithLogProductionTimeout(10*time.Second)},
			Consumers: logs,
		},
	}
	 return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started: true,
	})
}
