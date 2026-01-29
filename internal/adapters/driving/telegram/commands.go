package telegram

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot/models"
)

type AddSubject struct {
	Name string
}

type AddTags struct {
	SubjectName string
	Tags []string
}

type ListTags struct {
	SubjectName string
}

type CreateReservation struct {
	SubjectName string
	Duration int
}

type RemoveReservation struct {
	SubjectName string
}

type ActiveReservations struct {
	Tags []string
}

func ParseAddSubject(update *models.Update) (AddSubject, error) {
	parts := strings.SplitN(update.Message.Text, " ", 2)
	if len(parts) < 2 {
		return AddSubject{}, fmt.Errorf("Invalid format for add subject command. Expected: /add_subject <name>")
	}
	name := parts[1]

	return AddSubject{Name: name}, nil
}

func ParseAddTags(update *models.Update) (AddTags, error) {
	args := strings.SplitN(update.Message.Text, " ", 3)
	if len(args) < 3 {
		return AddTags{}, fmt.Errorf("Invalid format for add tags command. Expected: /add_tags <subject_name> <tag1> [tag2] [tag3]...")
	}
	tags := strings.Split(args[2], " ")

	return AddTags{
		SubjectName: args[1],
		Tags: tags,
	}, nil
}

func ParseListTags(update *models.Update) (ListTags, error) {
	parts := strings.SplitN(update.Message.Text, " ", 2)
	if len(parts) < 2 {
		return ListTags{}, fmt.Errorf("Invalid format for list tags command. Expected: /tags <subject_name>")
	}
	name := parts[1]

	return ListTags{SubjectName: name}, nil
}

func ParseCreateReservation(update *models.Update) (CreateReservation, error) {
	error := func () (CreateReservation, error) {
		return CreateReservation{}, fmt.Errorf("Invalid format for reserve command. Expected: /reserve <subject_name> <duration_in_minutes>")
	}

	parts := strings.SplitN(update.Message.Text, " ", 3)

	if len(parts) < 3 {
		return error()
	}
	subjectName := parts[1]
	minutes, err := strconv.Atoi(parts[2])

	if err != nil {
		return error()
	}
	return CreateReservation{
		SubjectName: subjectName,
		Duration: minutes,
	}, nil
}

func ParseRemoveReservation(update *models.Update) (RemoveReservation, error) {
	parts := strings.SplitN(update.Message.Text, " ", 2)
	if len(parts) < 2 {
		return RemoveReservation{}, fmt.Errorf("Invalid format for remove reservation command. Expected: /remove <subject_name>")
	}
	name := parts[1]

	return RemoveReservation{SubjectName: name}, nil
}

func ParseActiveReservations(update *models.Update) (ActiveReservations, error) {
	tags := make([]string, 0)
	parts := strings.SplitN(update.Message.Text, " ", 2)
	if len(parts) == 2 {
		tags = strings.Split(parts[1], " ")
	}
	return ActiveReservations{
		Tags: tags,
	}, nil
}
