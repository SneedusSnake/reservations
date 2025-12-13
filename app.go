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
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/kelseyhightower/envconfig"
)

const (
	CLOCK = "clock"

	STORE_SUBJECTS     = "subjects_store"
	STORE_USERS        = "users_store"
	STORE_TG_USERS     = "tg_users_store"
	STORE_RESERVATIONS = "reservations_store"

	SERVICE_SUBJECT = "subject_service"
	SERVICE_USER = "user_service"
	SERVICE_TELEGRAM_USER = "telegram_user_service"
	SERVICE_RESERVATION = "reservation_service"

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

	reservationService := application.NewReservationService(
		subjectsStore,
		reservationsStore,
		usersStore,
		clock,
	)
	subjectService := application.NewSubjectService(subjectsStore)
	userService := application.NewUserService(usersStore)
	tgUserService := telegram.NewTelegramUserService(tgUsersStore, *userService)

	tgBot := app.telegramBot()

	app.container[STORE_SUBJECTS] = subjectsStore
	app.container[STORE_USERS] = usersStore
	app.container[STORE_TG_USERS] = tgUsersStore
	app.container[STORE_RESERVATIONS] = reservationsStore
	app.container[CLOCK] = clock

	app.container[SERVICE_RESERVATION] = reservationService
	app.container[SERVICE_SUBJECT] = subjectService
	app.container[SERVICE_USER] = userService
	app.container[SERVICE_TELEGRAM_USER] = tgUserService

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
		app.Resolve(SERVICE_SUBJECT).(*application.SubjectService),
		app.Resolve(SERVICE_RESERVATION).(*application.ReservationService),
		app.Resolve(SERVICE_USER).(*application.UserService),
		app.Resolve(SERVICE_TELEGRAM_USER).(*telegram.TelegramUserService),
		app.Resolve(CLOCK).(ports.Clock),
		app.Log,
	)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/add_subject", bot.MatchTypePrefix, botHandlerFunc(adapter.AddSubjectHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/add_tags", bot.MatchTypePrefix, botHandlerFunc(adapter.AddSubjectTagsHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/list", bot.MatchTypeExact, botHandlerFunc(adapter.ListSubjectsHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/tags", bot.MatchTypePrefix, botHandlerFunc(adapter.ListSubjectTagsHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/reserve", bot.MatchTypePrefix, botHandlerFunc(adapter.CreateReservationHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/remove", bot.MatchTypePrefix, botHandlerFunc(adapter.RemoveReservationHandler))
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
				MessageThreadID: update.Message.MessageThreadID,
				Text:   text,
			})
		}
	}
}
