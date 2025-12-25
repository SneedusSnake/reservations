package acceptance

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/SneedusSnake/Reservations/internal/adapters/driven/clock/cache"
	"github.com/SneedusSnake/Reservations/testing/acceptance/drivers"
	"github.com/SneedusSnake/Reservations/testing/acceptance/drivers/telegram"
	"github.com/SneedusSnake/Reservations/testing/acceptance/specifications"
	"github.com/SneedusSnake/Reservations/testing/containers/app"
	"github.com/SneedusSnake/Reservations/testing/containers/telegram_api"
	"github.com/alecthomas/assert/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
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
		cache.NewClock(app.CLOCK_CACHE_PATH),
		testApp.App,
		t,
	)

	prepareTestFixtures(driver)
	cleanUp := func () {
		driver.CleanUp()
	}

	t.Run("User can see list of all existing subjects", func(t *testing.T) {
		specifications.ListSpecification(t, driver)
		t.Cleanup(cleanUp)
	})

	t.Run("User can see list of all tags attached to a subject", func(t *testing.T) {
		specifications.SubjectTagsSpecification(t, driver)
		t.Cleanup(cleanUp)
	})

	t.Run("User can make a reservation for a subject", func(t *testing.T) {
		specifications.ReserveSubjectSpecification(t, driver)
		t.Cleanup(cleanUp)
	})

	t.Run("User can remove reservations for a subject", func(t *testing.T) {
		specifications.RemoveReservationSpecification(t, driver)
		t.Cleanup(cleanUp)
	})

	t.Run("User can see list of all reservations", func(t *testing.T) {
		specifications.ListReservedSubjects(t, driver)
		t.Cleanup(cleanUp)
	})
}

func bootApplication(t *testing.T) *TestApplication {
	ctx := t.Context()
	net, err := network.New(ctx)
	assert.NoError(t, err)

	testApp := &TestApplication{Ctx: ctx, Network: net}
	apiContainer, err := telegram_api.Start(ctx, net.Name, &StdoutLogConsumer{Container: "Telegram test server"})
	testcontainers.CleanupContainer(t, apiContainer)
	assert.NoError(t, err)
	appContainer, err := app.Start(ctx, net.Name, &StdoutLogConsumer{Container: "Application"})
	testcontainers.CleanupContainer(t, appContainer)
	assert.NoError(t, err)
	testApp.TelegramApi = apiContainer
	testApp.App = appContainer
	
	return testApp
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
