package telegram

import (
	"context"
	"log"
	"strings"

	"github.com/SneedusSnake/Reservations/domain/reservations"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type telegramAdapter struct {
	subjectsStore reservations.SubjectsStore
	log *log.Logger
}

func NewAdapter(subStore reservations.SubjectsStore, log *log.Logger) *telegramAdapter {
	return &telegramAdapter{subjectsStore: subStore, log: log}
}

func (ta *telegramAdapter) AddSubjectHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	ta.log.Println("Handling add subject command")
	name := strings.SplitN(update.Message.Text, " ", 2)[1]
	subject := reservations.Subject{Id: ta.subjectsStore.NextIdentity(), Name: name}
	ta.subjectsStore.Add(subject)
}

func (ta *telegramAdapter) AddSubjectTagsHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	ta.log.Println("Handling add subject tags command")
	args := strings.SplitN(update.Message.Text, " ", 3)
	subject, err := ta.subjectsStore.List().Find(args[1])
	if err != nil {
		ta.log.Println(err)
	}
	for tag := range strings.SplitSeq(args[2], " ") {
		ta.subjectsStore.AddTag(subject.Id, string(tag))
	}
}

func (ta *telegramAdapter) ListSubjectsHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	ta.log.Println("Handling list command")
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   ta.subjectsStore.List().Names(),
	})
}

func (ta *telegramAdapter) ListSubjectTagsHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	ta.log.Println("Handling list subject tags command")
	args := strings.SplitN(update.Message.Text, " ", 2)
	subject, err := ta.subjectsStore.List().Find(args[1])
	if err != nil {
		ta.log.Println(err)
	}
	tags, err := ta.subjectsStore.GetTags(subject.Id)
	if err != nil {
		ta.log.Println(err)
	}
	ta.log.Printf("Returning tags for subject %s: %v", subject.Name, tags)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   strings.Join(tags, "\n"),
	})
}
