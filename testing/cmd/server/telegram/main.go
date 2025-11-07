package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Chat struct {
	Id int `json:"id"`
}

type User struct {
	Id int `json:"id"`
	IsBot bool `json:"is_bot"`
	FirstName string `json:"first_name"`
	Username string `json:"username"`
}

type UpdateMessage struct {
	Id int `json:"message_id"`
	Text string `json:"text"`
	From User `json:"from"`
	Chat Chat `json:"chat"`
}

type Update struct {
	Id int `json:"update_id"`
	Message UpdateMessage `json:"message"`
}

type Message struct {
	ChatId int `json:"chat_id"`
	Text string `json:"text"`
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

func sendBotMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	body, err := io.ReadAll(r.Body)
	var message Message
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = json.Unmarshal(body, &message)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	botMessages = append(botMessages, message)
	fmt.Fprint(w, "OK")
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	body, err := io.ReadAll(r.Body)
	var message UpdateMessage
	message.Id = len(messages) + 1
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	err = json.Unmarshal(body, &message)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	messages = append(messages, message)
	fmt.Fprint(w, "OK")
}

func getUpdates(w http.ResponseWriter, r *http.Request) {
	var updates []Update
	for i := lastReadId; i < len(messages); i++ {
		updates = append(updates, Update{Id: i+1, Message: messages[i]})
	}

	data, err := json.Marshal(updates)
	
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Fprint(w, string(data))
	lastReadId = updates[len(updates) - 1].Id
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(botMessages)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Fprint(w, string(data))
}

var botMessages []Message
var messages []UpdateMessage
var lastReadId int


func main () {
	handler := http.NewServeMux()
	handler.Handle("/", http.HandlerFunc(index))
	handler.Handle("/sendMessage", http.HandlerFunc(sendBotMessage))
	handler.Handle("/getUpdates", http.HandlerFunc(getUpdates))
	handler.Handle("/testing/sendClientMessage", http.HandlerFunc(sendMessage))
	handler.Handle("/testing/getBotMessages", http.HandlerFunc(getMessages))

	s := &http.Server{
		Addr:           ":8080",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
