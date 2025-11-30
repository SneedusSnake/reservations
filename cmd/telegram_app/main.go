package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/SneedusSnake/Reservations/adapters/driven/clock/cache"
	"github.com/SneedusSnake/Reservations/adapters/driven/clock/system"
	"github.com/SneedusSnake/Reservations/adapters/driven/persistence/inmemory"
	"github.com/SneedusSnake/Reservations/adapters/driving/telegram"
	"github.com/SneedusSnake/Reservations/application"
	"github.com/SneedusSnake/Reservations/domain"
	"github.com/SneedusSnake/Reservations/domain/reservations"
	"github.com/SneedusSnake/Reservations/domain/users"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TelegramApi struct{
		Host string `envconfig:"TELEGRAM_API_HOST"`
		Token string `envconfig:"TELEGRAM_API_TOKEN"`
	}
	Clock string `envconfig:"CLOCK_DRIVER"`
	CacheClockPath string `envconfig:"CACHE_CLOCK_PATH"`
}

var cfg Config;
var subjectsStore reservations.SubjectsStore
var usersStore users.UsersStore
var tgUsersStore users.TelegramUsersStore
var reservationsRegistry reservations.ReservationsRegistry
var clock domain.Clock

func main() {
	log.Print("Starting main")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	loadConfig()

	subjectsStore = inmemory.NewSubjectsStore()
	usersStore = inmemory.NewUsersStore()
	tgUsersStore = inmemory.NewTelegramUsersStore(usersStore)
	reservationsRegistry = inmemory.NewReservationStore()
	clock = system.SystemClock{}
	if cfg.Clock == "cache" {
		clock = cache.NewClock(cfg.CacheClockPath)
	} 
	createHandler := application.NewCreateReservationHandler(subjectsStore, reservationsRegistry, usersStore, clock)

	b := tgBot()
	adapter := telegram.NewAdapter(subjectsStore, usersStore, tgUsersStore, reservationsRegistry, createHandler, clock, log.Default())

	b.RegisterHandler(bot.HandlerTypeMessageText, "/add_subject", bot.MatchTypePrefix, botHandlerFunc(adapter.AddSubjectHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/add_tags", bot.MatchTypePrefix, botHandlerFunc(adapter.AddSubjectTagsHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/list", bot.MatchTypeExact, botHandlerFunc(adapter.ListSubjectsHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/tags", bot.MatchTypePrefix, botHandlerFunc(adapter.ListSubjectTagsHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/reserve", bot.MatchTypePrefix, botHandlerFunc(adapter.ReservationHandler))
	b.Start(ctx)
}

func loadConfig() {
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Print(err)
		panic(err)
	}
}

func tgBot() *bot.Bot {
	opts := []bot.Option{
		bot.WithServerURL(cfg.TelegramApi.Host),
	}
	b, err := bot.New(cfg.TelegramApi.Token, opts...)

	if err != nil {
		log.Print(err)
		panic(err)
	}

	return b
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
