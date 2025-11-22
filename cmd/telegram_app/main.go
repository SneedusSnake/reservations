package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/SneedusSnake/Reservations/adapters/driven/persistence/inmemory"
	"github.com/SneedusSnake/Reservations/domain/reservations"
	"github.com/SneedusSnake/Reservations/adapters/driving/telegram"
	"github.com/go-telegram/bot"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TelegramApi struct{
		Host string `envconfig:"TELEGRAM_API_HOST"`
		Token string `envconfig:"TELEGRAM_API_TOKEN"`
	}
}

var cfg Config;
var subjectsStore reservations.SubjectsStore

func main() {
	log.Print("Starting main")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	loadConfig()

	subjectsStore = inmemory.NewSubjectsStore()
	b := tgBot()
	adapter := telegram.NewAdapter(subjectsStore, log.Default())

	b.RegisterHandler(bot.HandlerTypeMessageText, "/add_subject", bot.MatchTypePrefix, adapter.AddSubjectHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/add_tags", bot.MatchTypePrefix, adapter.AddSubjectTagsHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/list", bot.MatchTypeExact, adapter.ListSubjectsHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/tags", bot.MatchTypePrefix, adapter.ListSubjectTagsHandler)
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
