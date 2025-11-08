package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/SneedusSnake/Reservations/adapters/driven/persistence/inmemory"
	"github.com/SneedusSnake/Reservations/domain/reservations"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TelegramApi struct{
		Host string `envconfig:"TELEGRAM_API_HOST"`
		Token string `envconfig:"TELEGRAM_API_TOKEN"`
	}
}

type Message struct {
	ChatId int `json:"chat_id"`
	Text string `json:"text"`
}

type Update struct {
	Id int `json:"update_id"`
	Message struct{
		Text string `json:"text"`
		Chat struct{
			Id int `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

var cfg Config;
var subjectsStore reservations.SubjectsStore

func main() {
	log.Print("Starting main")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Print(err)
		panic(err)
	}
	subjectsStore = inmemory.NewSubjectsStore()

	opts := []bot.Option{
		bot.WithServerURL(cfg.TelegramApi.Host),
	}
	b, err := bot.New(cfg.TelegramApi.Token, opts...)

	if err != nil {
		log.Print(err)
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/add_subject", bot.MatchTypePrefix, addSubjectHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/list", bot.MatchTypeExact, listSubjectsHandler)
	b.Start(ctx)
}

func subjects() string {
	subjects := subjectsStore.List()
	subjectNames := []string{}
	for _, subject := range subjects {
		subjectNames = append(subjectNames, subject.Name)
	}
	return strings.Join(subjectNames, "\n")
}

func addSubjectHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	log.Println("Handling add subject command")
	name := strings.SplitN(update.Message.Text, " ", 2)[1]
	subject := reservations.Subject{Id: subjectsStore.NextIdentity(), Name: name}
	subjectsStore.Add(subject)
}

func listSubjectsHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	log.Println("Handling list command")
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   subjects(),
	})
}
