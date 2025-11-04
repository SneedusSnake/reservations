package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

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

type TelegramDriver struct{
	client http.Client
	host string
	t *testing.T
}

func (d TelegramDriver) UserRequestsSubjectsList() {
	clientMessage := `{
		"chat": {"id": 1234},
		"text": "/list"
	}`
	_, err := d.client.Post(fmt.Sprintf("%s/sendClientMessage", d.host), "application/json", bytes.NewBuffer([]byte(clientMessage)))
	assert.NoError(d.t, err)
}

func (d TelegramDriver) UserSeesSubjects(subject ...string) {
	var responseData []struct{
		Text string `json:"text"`
	}

	for i := 0; i < 10 && len(responseData) == 0; i++ {
		r, err := d.client.Get(fmt.Sprintf("%s/getBotMessages", d.host))
		assert.NoError(d.t, err)

		body, err := io.ReadAll(r.Body)
		assert.NoError(d.t, err)
		err = json.Unmarshal(body, &responseData)
		assert.NoError(d.t, err)
	}

	assert.NotEqual(d.t, 0, len(responseData))

	subjects := strings.Split(responseData[len(responseData)-1].Text, "\n")
	assert.Equal(d.t, subject, subjects)
}

func TestList(t *testing.T) {
	ctx := context.Background()
	net, err := network.New(ctx)
	assert.NoError(t, err)
	apiContainer, err := bootTelegramApiContainer(ctx, net)
	testcontainers.CleanupContainer(t, apiContainer)
	assert.NoError(t, err)
	appContainer, err := bootAppContainer(ctx, net)
	testcontainers.CleanupContainer(t, appContainer)
	assert.NoError(t, err)
	apiHost, err := apiContainer.Endpoint(ctx, "")
	assert.NoError(t, err)

	specifications.ListSpecification(t, TelegramDriver{
		client: *http.DefaultClient,
		t: t,
		host: "http://" + apiHost,
	})
}

func bootTelegramApiContainer(ctx context.Context, network *testcontainers.DockerNetwork) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context: "./server/telegram",
			Dockerfile: "Dockerfile",
			PrintBuildLog: true,
		},
		Networks: []string{network.Name},
		NetworkAliases: map[string][]string{network.Name: []string{"telegram-api"}},
		ExposedPorts: []string{"8080"},
		WaitingFor: wait.ForHTTP("/").WithPort("8080"),
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Opts: []testcontainers.LogProductionOption{testcontainers.WithLogProductionTimeout(10*time.Second)},
			Consumers: []testcontainers.LogConsumer{&StdoutLogConsumer{Container: "Telegram test server"}},
		},
	}
	 return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started: true,
	})
}

func bootAppContainer(ctx context.Context, network *testcontainers.DockerNetwork) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context: "../..",
			Dockerfile: "./build/Docker/Dockerfile",
			PrintBuildLog: true,
		},
		Env: map[string]string{"TELEGRAM_API_HOST": "http://telegram-api:8080"},
		Networks: []string{network.Name},
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Opts: []testcontainers.LogProductionOption{testcontainers.WithLogProductionTimeout(10*time.Second)},
			Consumers: []testcontainers.LogConsumer{&StdoutLogConsumer{Container: "Application"}},
		},
	}
	 return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started: true,
	})
}

