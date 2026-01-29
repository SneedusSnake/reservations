package telegram_test

import (
	"testing"

	"github.com/SneedusSnake/Reservations/internal/adapters/driving/telegram"
	"github.com/alecthomas/assert/v2"
	"github.com/go-telegram/bot/models"
)

func TestInputParsers(t *testing.T) {
	t.Run("it parses AddSubject command", func(t *testing.T) {
		update := telegramUpdate("/add_subject Test")
		cmd, err := telegram.ParseAddSubject(update)
		assert.NoError(t, err)
		assert.Equal(t, cmd.Name, "Test")
	})

	t.Run("it returns error given no arguments provided to AddSubject", func(t *testing.T) {
		update := telegramUpdate("/add_subject")
		_, err := telegram.ParseAddSubject(update)
		assert.Error(t, err)
	})

	t.Run("it parses AddTags command", func(t *testing.T) {
		update := telegramUpdate("/add_tags Test tag1 tag2 tag3")
		cmd, err := telegram.ParseAddTags(update)
		assert.NoError(t, err)
		assert.Equal(t, cmd.SubjectName, "Test")
		assert.SliceContains(t, cmd.Tags, "tag1")
		assert.SliceContains(t, cmd.Tags, "tag2")
		assert.SliceContains(t, cmd.Tags, "tag3")
	})

	t.Run("it returns error given wrong format provided to AddSbujectTags", func(t *testing.T) {
		update := telegramUpdate("/add_tags Test")
		_, err := telegram.ParseAddTags(update)
		assert.Error(t, err)

		update = telegramUpdate("/add_tags")
		_, err = telegram.ParseAddTags(update)
		assert.Error(t, err)
	})

	t.Run("it parses ListTags command", func(t *testing.T) {
		update := telegramUpdate("/list_tags Test")
		cmd, err := telegram.ParseListTags(update)
		assert.NoError(t, err)
		assert.Equal(t, cmd.SubjectName, "Test")
	})

	t.Run("it returns error given no arguments provided to ListTags", func(t *testing.T) {
		update := telegramUpdate("/list_tags")
		_, err := telegram.ParseListTags(update)
		assert.Error(t, err)
	})

	t.Run("it parses CreateReservation command", func(t *testing.T) {
		update := telegramUpdate("/reserve Test 10")
		cmd, err := telegram.ParseCreateReservation(update)
		assert.NoError(t, err)
		assert.Equal(t, cmd.SubjectName, "Test")
		assert.Equal(t, cmd.Duration, 10)
	})
	
	t.Run("it returns error given wrong arguments provided to CreateReservation", func(t *testing.T) {
		update := telegramUpdate("/reserve")
		_, err := telegram.ParseCreateReservation(update)
		assert.Error(t, err)

		update = telegramUpdate("/reserve Test")
		_, err = telegram.ParseCreateReservation(update)
		assert.Error(t, err)

		update = telegramUpdate("/reserve Test invalid_duration")
		_, err = telegram.ParseCreateReservation(update)
		assert.Error(t, err)
	})

	t.Run("it parses RemoveReservation command", func(t *testing.T) {
		update := telegramUpdate("/remove Test")
		cmd, err := telegram.ParseRemoveReservation(update)
		assert.NoError(t, err)
		assert.Equal(t, cmd.SubjectName, "Test")
	})

	t.Run("it returns error given no arguments provided to RemoveReservation", func(t *testing.T) {
		update := telegramUpdate("/remove")
		_, err := telegram.ParseRemoveReservation(update)
		assert.Error(t, err)
	})

	t.Run("it parses ActiveReservations command", func(t *testing.T) {
		update := telegramUpdate("/reserved")
		cmd, err := telegram.ParseActiveReservations(update)
		assert.NoError(t, err)
		assert.Equal(t, cmd.Tags, []string{})

		update = telegramUpdate("/reserved tag1 tag2 tag3")
		cmd, err = telegram.ParseActiveReservations(update)
		assert.NoError(t, err)
		assert.Equal(t, cmd.Tags, []string{"tag1", "tag2", "tag3"})
	})
}

func telegramUpdate(text string) *models.Update {
	return &models.Update{
		Message: &models.Message{
			Text: text,
		},
	}
}
