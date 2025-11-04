package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/SneedusSnake/Reservations/domain/reservations"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TelegramApi struct{
		Host string `envconfig:"TELEGRAM_API_HOST"`
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

func main() {
	log.Print("Starting main")
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Print(err)
		panic(err)
	}


	for {
		updates, err := updates()
		if err != nil {
			log.Print(err)
		}

		if len(updates) != 0 {
			handleUpdate(updates[0])
		}

		time.Sleep(time.Microsecond*500)
	}
}

func updates() ([]Update, error) {
	var updates []Update
	r, err := http.Get(fmt.Sprintf("%s/getUpdates", cfg.TelegramApi.Host))
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	log.Print("recieved updates:", string(body))

	err = json.Unmarshal(body, &updates)

	if err != nil {
		return nil, err
	}

	return updates, nil
}

func subjects() string {
	subjects := []reservations.Subject{
		{Id: 1, Name: "Subject #1"},
		{Id: 2, Name: "Subject #2"},
		{Id: 2, Name: "Subject #3"},
	}

	subjectNames := []string{}
	for _, subject := range subjects {
		subjectNames = append(subjectNames, subject.Name)
	}
	return strings.Join(subjectNames, "\n")
}

func handleUpdate(u Update) error {
	if u.Message.Text == "/list" {
		msg := Message{ChatId: u.Message.Chat.Id, Text: subjects()}
		data, err := json.Marshal(msg)

		r, err := http.Post(fmt.Sprintf("%s/sendMessage", cfg.TelegramApi.Host), "application/json", strings.NewReader(string(data)))

		if err != nil {
			log.Print(err)
			panic(err)
		}
		if r.StatusCode != 200 {
			log.Print(err)
		}
	}

	return nil
}
