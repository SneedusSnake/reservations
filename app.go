package app

import (
	"log"
	"context"

	"github.com/SneedusSnake/Reservations/internal/adapters/driven/clock/cache"
	"github.com/SneedusSnake/Reservations/internal/adapters/driven/clock/system"
	"github.com/SneedusSnake/Reservations/internal/adapters/driven/persistence/inmemory"
	"github.com/SneedusSnake/Reservations/internal/adapters/driving/telegram"
	"github.com/SneedusSnake/Reservations/internal/application"
	"github.com/SneedusSnake/Reservations/internal/ports"
	"github.com/SneedusSnake/Reservations/internal/ports/reservations"
	"github.com/SneedusSnake/Reservations/internal/ports/users"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/kelseyhightower/envconfig"
)

const (
	CLOCK              = "clock"

	STORE_SUBJECTS     = "subjects_store"
	STORE_USERS        = "users_store"
	STORE_TG_USERS     = "tg_users_store"
	STORE_RESERVATIONS = "reservations_store"

	HANDLER_CREATE_RESERVATION = "reservation_create_handler"

	TELERAM_BOT = "telegram_bot"
)

type App struct{
	Config Config
	Log *log.Logger
	container map[string]any
}

type Config struct {
	TelegramApi struct{
		Host string `envconfig:"TELEGRAM_API_HOST"`
		Token string `envconfig:"TELEGRAM_API_TOKEN"`
	}
	Clock string `envconfig:"CLOCK_DRIVER"`
	CacheClockPath string `envconfig:"CACHE_CLOCK_PATH"`
}

func (app *App) Resolve(dependency string) any {
	result, ok := app.container[dependency]
	if !ok {
		app.Log.Fatalf("Could not resolve %s", dependency)
	}

	return result
}

func Bootstrap() *App {
	app := &App{
		Log: log.Default(),
		container: make(map[string]any),
	}
	app.loadConfig()
	app.registerDependencies()

	return app
}

func (app *App) loadConfig() {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		app.Log.Print(err)
		panic(err)
	}
	app.Config = cfg
}

func (app *App) registerDependencies() {
	var clock ports.Clock

	clock = system.SystemClock{}
	if app.Config.Clock == "cache" {
		clock = cache.NewClock(app.Config.CacheClockPath)
	} 

	subjectsStore := inmemory.NewSubjectsStore()
	usersStore := inmemory.NewUsersStore()
	tgUsersStore := inmemory.NewTelegramUsersStore(usersStore)
	reservationsStore := inmemory.NewReservationStore()


	createReservationHandler := application.NewCreateReservationHandler(
		subjectsStore,
		reservationsStore,
		usersStore,
		clock,
	)

	tgBot := app.telegramBot()

	app.container[STORE_SUBJECTS] = subjectsStore
	app.container[STORE_USERS] = usersStore
	app.container[STORE_TG_USERS] = tgUsersStore
	app.container[STORE_RESERVATIONS] = reservationsStore
	app.container[CLOCK] = clock

	app.container[HANDLER_CREATE_RESERVATION] = createReservationHandler

	app.container[TELERAM_BOT] = tgBot
	app.registerTelegramBotHandlers()
}

func (app *App) telegramBot() *bot.Bot {
	opts := []bot.Option{
		bot.WithServerURL(app.Config.TelegramApi.Host),
	}
	b, err := bot.New(app.Config.TelegramApi.Token, opts...)

	if err != nil {
		app.Log.Print(err)
		panic(err)
	}

	return b
}

func (app *App) registerTelegramBotHandlers() {
	b := app.Resolve(TELERAM_BOT).(*bot.Bot)
	adapter := telegram.NewAdapter(
		app.Resolve(STORE_SUBJECTS).(reservations.SubjectsRepository),
		app.Resolve(STORE_USERS).(users.UsersRepository),
		app.Resolve(STORE_TG_USERS).(users.TelegramUsersRepository),
		app.Resolve(STORE_RESERVATIONS).(reservations.ReservationsRepository),
		app.Resolve(HANDLER_CREATE_RESERVATION).(*application.CreateReservationHandler),
		app.Resolve(CLOCK).(ports.Clock),
		app.Log,
	)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/add_subject", bot.MatchTypePrefix, botHandlerFunc(adapter.AddSubjectHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/add_tags", bot.MatchTypePrefix, botHandlerFunc(adapter.AddSubjectTagsHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/list", bot.MatchTypeExact, botHandlerFunc(adapter.ListSubjectsHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/tags", bot.MatchTypePrefix, botHandlerFunc(adapter.ListSubjectTagsHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/reserve", bot.MatchTypePrefix, botHandlerFunc(adapter.ReservationHandler))
}

type UpdateHandler func(ctx context.Context, b *bot.Bot, update *models.Update) (string, error)

func botHandlerFunc(h UpdateHandler) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		text, err := h(ctx, b, update)

		if err != nil {
			log.Print(err)
			return
		}

		if text != "" {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   text,
			})
		}
	}
}
