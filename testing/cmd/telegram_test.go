package cmd

import (
	"testing"
	"github.com/SneedusSnake/Reservations/testing/specifications"
)

type TelegramViewer struct{}

func (v TelegramViewer) List() []string {
	return []string{"Test Subject #1", "Test Subject #2", "Test Subject #3"}
}

func TestList(t *testing.T) {
	specifications.ListSpecification(t, TelegramViewer{})
}
