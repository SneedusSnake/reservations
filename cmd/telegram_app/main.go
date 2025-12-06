package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/SneedusSnake/Reservations"
	"github.com/go-telegram/bot"
)

func main() {
	log.Print("Starting main")
	application := app.Bootstrap()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	b := application.Resolve(app.TELERAM_BOT).(*bot.Bot)

	b.Start(ctx)
}

