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

type UpdateHandler func(ctx context.Context, b *bot.Bot, update *models.Update) (string, error)

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

func (ta *telegramAdapter) AddSubjectHandler(ctx context.Context, b *bot.Bot, update *models.Update) (string, error) {
	ta.log.Println("Handling add subject command")
	name := strings.SplitN(update.Message.Text, " ", 2)[1]
	subject := reservations.Subject{Id: ta.subjectsStore.NextIdentity(), Name: name}
	ta.subjectsStore.Add(subject)

	return "", nil
}

func (ta *telegramAdapter) AddSubjectTagsHandler(ctx context.Context, b *bot.Bot, update *models.Update) (string, error) {
	ta.log.Println("Handling add subject tags command")
	args := strings.SplitN(update.Message.Text, " ", 3)
	subject, err := ta.subjectsStore.List().Find(args[1])
	if err != nil {
		return "", err
	}

	for tag := range strings.SplitSeq(args[2], " ") {
		ta.subjectsStore.AddTag(subject.Id, string(tag))
	}

	return "", nil
}

func (ta *telegramAdapter) ListSubjectsHandler(ctx context.Context, b *bot.Bot, update *models.Update) (string, error) {
	ta.log.Println("Handling list command")

	return ta.subjectsStore.List().Names(), nil;
}

func (ta *telegramAdapter) ListSubjectTagsHandler(ctx context.Context, b *bot.Bot, update *models.Update) (string, error) {
	ta.log.Println("Handling list subject tags command")
	args := strings.SplitN(update.Message.Text, " ", 2)
	subject, err := ta.subjectsStore.List().Find(args[1])
	if err != nil {
		return "", err
	}
	tags, err := ta.subjectsStore.GetTags(subject.Id)
	if err != nil {
		return "", err
	}

	return strings.Join(tags, "\n"), nil
}

func (ta *telegramAdapter) ReservationHandler(ctx context.Context, b *bot.Bot, update *models.Update) (string, error) {
	ta.log.Println("Handling reservation command")
	args := strings.SplitN(update.Message.Text, " ", 3)
	subject, err := ta.subjectsStore.List().Find(args[1])
	if err != nil {
		return "", err
	}
	minutes, err := strconv.Atoi(args[2])

	if err != nil {
		return "", err
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
			return fmt.Sprintf("Already reserved by %s until %s", u.Name, r.End.Format(time.DateTime)), nil
		}
		return "", err
	}

	return fmt.Sprintf("Reservation for %s acquired by %s until %s", subject.Name, user.Name, r.End.Format(time.DateTime)), nil
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
