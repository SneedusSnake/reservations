package cmd

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/SneedusSnake/Reservations/testing/drivers"
	"github.com/SneedusSnake/Reservations/testing/drivers/telegram"
	"github.com/SneedusSnake/Reservations/testing/specifications"
	"github.com/alecthomas/assert/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

type StdoutLogConsumer struct{
	Container string
}

func (lc *StdoutLogConsumer) Accept(l testcontainers.Log) {
	fmt.Printf("%s: %s\n", lc.Container, string(l.Content))
}

type TestApplication struct {
	BotToken string
	TelegramApi testcontainers.Container
	App testcontainers.Container
	Network *testcontainers.DockerNetwork
	Ctx context.Context
}

func (app *TestApplication) TelegramApiHost() (string, error) {
	host, err := app.TelegramApi.Endpoint(app.Ctx, "")
	
	if err != nil {
		return "", err
	}

	return "http://" + host, nil
}

func TestSuite(t *testing.T) {
	testApp := bootApplication(t)
	telegramApiHost, err := testApp.TelegramApiHost()
	assert.NoError(t, err)

	driver := telegram.NewDriver(
		http.DefaultClient,
		telegramApiHost,
		t,
	)

	prepareTestFixtures(driver)

	t.Run("User can see list of all existing subjects", func(t *testing.T) {
		specifications.ListSpecification(t, driver)
	})

	t.Run("User can see list of all tags attached to a subject", func(t *testing.T) {
		specifications.SubjectTagsSpecification(t, driver)
	})
}

func bootApplication(t *testing.T) *TestApplication {
	ctx := context.Background()
	net, err := network.New(ctx)
	assert.NoError(t, err)

	testApp := &TestApplication{Ctx: ctx, Network: net}
	apiContainer, err := bootTelegramApiContainer(testApp)
	testcontainers.CleanupContainer(t, apiContainer)
	assert.NoError(t, err)
	appContainer, err := bootAppContainer(testApp)
	testcontainers.CleanupContainer(t, appContainer)
	assert.NoError(t, err)
	testApp.TelegramApi = apiContainer
	testApp.App = appContainer
	
	return testApp
}

func bootTelegramApiContainer(app *TestApplication) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context: "./server/telegram",
			Dockerfile: "Dockerfile",
			PrintBuildLog: true,
		},
		Env: map[string]string{"BOT_TOKEN": "1234567"},
		Networks: []string{app.Network.Name},
		NetworkAliases: map[string][]string{app.Network.Name: []string{"telegram-api"}},
		ExposedPorts: []string{"8080"},
		WaitingFor: wait.ForHTTP("/").WithPort("8080"),
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Opts: []testcontainers.LogProductionOption{testcontainers.WithLogProductionTimeout(10*time.Second)},
			Consumers: []testcontainers.LogConsumer{&StdoutLogConsumer{Container: "Telegram test server"}},
		},
	}
	 return testcontainers.GenericContainer(app.Ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started: true,
	})
}

func bootAppContainer(app *TestApplication) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context: "../..",
			Dockerfile: "./build/Docker/Dockerfile",
			PrintBuildLog: true,
		},
		Env: map[string]string{"TELEGRAM_API_HOST": "http://telegram-api:8080", "TELEGRAM_API_TOKEN": "1234567"},
		Networks: []string{app.Network.Name},
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Opts: []testcontainers.LogProductionOption{testcontainers.WithLogProductionTimeout(10*time.Second)},
			Consumers: []testcontainers.LogConsumer{&StdoutLogConsumer{Container: "Application"}},
		},
	}
	 return testcontainers.GenericContainer(app.Ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started: true,
	})
}

func prepareTestFixtures(driver drivers.Reservations) {
	driver.AdminAddsSubject("Subject#1")
	driver.AdminAddsSubject("Subject#2")
	driver.AdminAddsSubject("Subject#3")
	time.Sleep(time.Millisecond*100)
	driver.AdminAddsTagsToSubject("Subject#1", "Subject#1", "This_is_a_first_subject", "test")
	driver.AdminAddsTagsToSubject("Subject#2", "Subject#2", "This_is_a_second_subject", "test")
	time.Sleep(time.Millisecond*500)
}
