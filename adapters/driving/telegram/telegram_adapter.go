package telegram

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/SneedusSnake/Reservations/application"
	"github.com/SneedusSnake/Reservations/domain"
	"github.com/SneedusSnake/Reservations/domain/reservations"
	"github.com/SneedusSnake/Reservations/domain/users"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type telegramAdapter struct {
	subjectsStore reservations.SubjectsStore
	usersStore users.UsersStore
	tgStore users.TelegramUsersStore
	reservationsRegistry reservations.ReservationsRegistry
	createHandler *application.CreateReservationHandler
	clock domain.Clock
	log *log.Logger
}

func NewAdapter(
	subStore reservations.SubjectsStore,
	usersStore users.UsersStore,
	tgStore users.TelegramUsersStore,
	reservations reservations.ReservationsRegistry,
	createHandler *application.CreateReservationHandler,
	clock domain.Clock,
	log *log.Logger,
) *telegramAdapter {
	return &telegramAdapter{
		subjectsStore: subStore,
		usersStore: usersStore,
		tgStore: tgStore,
		reservationsRegistry: reservations,
		createHandler: createHandler,
		clock: clock,
		log: log,
	}
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

func (ta *telegramAdapter) ReservationHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	ta.log.Println("Handling reservation command")
	args := strings.SplitN(update.Message.Text, " ", 3)
	subject, err := ta.subjectsStore.List().Find(args[1])
	if err != nil {
		ta.log.Println(err)
	}
	minutes, err := strconv.Atoi(args[2])

	if err != nil {
		ta.log.Println(err)
	}
	user, err := ta.tgStore.Get(update.Message.From.ID)

	if err != nil {
		user = ta.createNewTelegramUser(update.Message.From)
	}

	cmd := application.CreateReservation{UserId: user.Id, SubjectId: subject.Id, From: ta.clock.Current(), To: ta.clock.Current().Add(time.Duration(minutes)*time.Minute)}
	r, err := ta.createHandler.Handle(cmd)

	if err != nil {
		if reservedErr, ok := err.(application.AlreadyReservedError); ok {
			r, _ := ta.reservationsRegistry.Get(reservedErr.ReservationIds[0])
			u, _ := ta.usersStore.Get(r.Id)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("Already reserved by %s until %s", u.Name, r.End.Format(time.DateTime)),
			})
		}
		ta.log.Println(err)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Reservation for %s acquired by %s until %s", subject.Name, user.Name, r.End.Format(time.DateTime)),
	})
}

func (ta *telegramAdapter) createNewTelegramUser(user *models.User) users.TelegramUser {
	tgUser := users.TelegramUser{
		TelegramId: user.ID,
		User: users.User{
			Id: ta.usersStore.NextIdentity(),
			Name: user.FirstName,
		},
	}
	ta.usersStore.Add(tgUser.User)
	ta.tgStore.Add(tgUser)

	return tgUser
}
