package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/SneedusSnake/Reservations/internal/ports"
	"github.com/SneedusSnake/Reservations/internal/application"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type telegramAdapter struct {
	subjectService *application.SubjectService
	reservationsService *application.ReservationService
	userService *application.UserService
	telegramUserService *TelegramUserService
	clock ports.Clock
	log *log.Logger
}

type UpdateHandler func(ctx context.Context, b *bot.Bot, update *models.Update) (string, error)

func NewAdapter(
	subjectService *application.SubjectService,
	reservationService *application.ReservationService,
	userService *application.UserService,
	telegramUserService *TelegramUserService,
	clock ports.Clock,
	log *log.Logger,
) *telegramAdapter {
	return &telegramAdapter{
		subjectService: subjectService,
		reservationsService: reservationService,
		telegramUserService: telegramUserService,
		userService: userService,
		clock: clock,
		log: log,
	}
}

func (ta *telegramAdapter) AddSubjectHandler(ctx context.Context, b *bot.Bot, update *models.Update) (string, error) {
	ta.log.Println("Handling add subject command")
	input, err := ParseAddSubject(update)
	if err != nil {
		return err.Error(), nil
	}
	_, err = ta.subjectService.Create(input.Name)

	if err != nil {
		return "", err
	}

	return "", nil
}

func (ta *telegramAdapter) AddSubjectTagsHandler(ctx context.Context, b *bot.Bot, update *models.Update) (string, error) {
	ta.log.Println("Handling add subject tags command")
	input, err := ParseAddTags(update)
	if err != nil {
		return err.Error(), nil
	}
	subject, err := ta.subjectService.GetByName(input.SubjectName)
	if err != nil {
		return "", err
	}
	cmd := application.AddTags{SubjectId: subject.Id, Tags: input.Tags}
	err = ta.subjectService.AddTags(cmd)

	if err != nil {
		return "", err
	}

	return "", nil
}

func (ta *telegramAdapter) ListSubjectsHandler(ctx context.Context, b *bot.Bot, update *models.Update) (string, error) {
	ta.log.Println("Handling list command")
	subjects, err := ta.subjectService.List()

	if err != nil {
		return "", nil
	}

	return subjects.Names(), nil;
}

func (ta *telegramAdapter) ListSubjectTagsHandler(ctx context.Context, b *bot.Bot, update *models.Update) (string, error) {
	ta.log.Println("Handling list subject tags command")
	input, err := ParseListTags(update)
	if err != nil {
		return err.Error(), nil
	}
	subject, err := ta.subjectService.GetByName(input.SubjectName)
	if err != nil {
		return "", err
	}
	tags, err := ta.subjectService.ListTags(subject.Id)
	if err != nil {
		return "", err
	}

	return strings.Join(tags, "\n"), nil
}

func (ta *telegramAdapter) CreateReservationHandler(ctx context.Context, b *bot.Bot, update *models.Update) (string, error) {
	ta.log.Println("Handling create reservation command")
	input, err := ParseCreateReservation(update)
	if err != nil {
		return err.Error(), nil
	}

	subject, err := ta.subjectService.GetByName(input.SubjectName)
	if err != nil {
		return "", err
	}

	user, err := ta.telegramUserService.Get(update.Message.From.ID)

	if err != nil {
		user, err = ta.telegramUserService.Create(CreateUser{update.Message.From.ID, update.Message.From.FirstName})
		if err != nil {
			return "", err
		}
	}

	cmd := application.CreateReservation{UserId: user.Id, SubjectId: subject.Id, From: ta.clock.Current(), To: ta.clock.Current().Add(time.Duration(input.Duration)*time.Minute)}
	r, err := ta.reservationsService.Create(cmd)

	if err != nil {
		if reservedErr, ok := err.(application.AlreadyReservedError); ok {
			r, _ := ta.reservationsService.Get(reservedErr.ReservationIds[0])
			u, _ := ta.userService.Get(r.UserId)
			return fmt.Sprintf("Already reserved by %s until %s", u.Name, r.End.Format(time.DateTime)), nil
		}
		return "", err
	}

	return fmt.Sprintf("Reservation for %s acquired by %s until %s", subject.Name, user.Name, r.End.Format(time.DateTime)), nil
}

func (ta *telegramAdapter) RemoveReservationHandler(ctx context.Context, b *bot.Bot, update *models.Update) (string, error) {
	ta.log.Println("Handling remove reservation command")
	input, err := ParseRemoveReservation(update)
	if err != nil {
		return err.Error(), nil
	}
	subject, err := ta.subjectService.GetByName(input.SubjectName)
	if err != nil {
		return "", err
	}

	user, err := ta.telegramUserService.Get(update.Message.From.ID)

	if err != nil {
		return "", err
	}

	err = ta.reservationsService.Remove(application.RemoveReservations{UserId: user.Id, SubjectId: subject.Id})
	
	if err != nil {
		return "", err
	}

	return subject.Name, nil
}

func (ta *telegramAdapter) ActiveReservationsHandler(ctx context.Context, b *bot.Bot, update *models.Update) (string, error) {
	ta.log.Println("Handling active reservations command")
	input, err := ParseActiveReservations(update)
	if err != nil {
		return err.Error(), nil
	}
	list, err := ta.reservationsService.ActiveReservations(ta.clock.Current(), input.Tags...)

	if err != nil {
		return "", err
	}

	text := "Subject\tReserved Until\t\tUser\n"
	for _, reservation := range list {
		text += fmt.Sprintf("%s\t%s\t\t%s\n", reservation.Subject, reservation.End.Format(time.DateTime), reservation.User)
	}

	return text, nil
}
