package telegram_api

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func Start(ctx context.Context, network string, logs ...testcontainers.LogConsumer) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context: testsRootDir() + "/containers/telegram_api/server",
			Dockerfile: "Dockerfile",
			PrintBuildLog: true,
		},
		Env: map[string]string{"BOT_TOKEN": "1234567"},
		Networks: []string{network},
		NetworkAliases: map[string][]string{network: {"telegram-api"}},
		ExposedPorts: []string{"8080"},
		WaitingFor: wait.ForHTTP("/").WithPort("8080"),
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

func testsRootDir() string {
	wd, _ := os.Getwd()

	return strings.SplitAfter(wd, "testing")[0]
}
