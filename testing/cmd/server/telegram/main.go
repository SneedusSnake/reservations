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

type Message struct {
	Id int `json:"message_id"`
	Text string `json:"text"`
	From User `json:"from"`
	Chat Chat `json:"chat"`
}

type Update struct {
	Id int `json:"update_id"`
	Message Message `json:"message"`
}

type BotMessage struct {
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
	log.Print("Recieved message from bot: " + string(body))
	var message BotMessage
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = json.Unmarshal(body, &message)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	botMessages = append(botMessages, message)
	log.Printf("Parsed message: %v, \n Total messages: %v\n", message, botMessages)
	fmt.Fprint(w, "OK")
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	body, err := io.ReadAll(r.Body)
	var message Message
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


	log.Print("writing client message", message)
	messages = append(messages, message)
	fmt.Fprint(w, "OK")
}

func getUpdates(w http.ResponseWriter, r *http.Request) {
	var updates []Update
	for i := 0; i < len(messages); i++ {
		updates = append(updates, Update{Id: i, Message: messages[i]})
	}

	data, err := json.Marshal(updates)
	
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Printf("updates requested, returning %v", updates)
	fmt.Fprint(w, string(data))
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(botMessages)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Printf("bot messages requested, returning %v", botMessages)

	fmt.Fprint(w, string(data))
}

var botMessages []BotMessage
var messages []Message


func main () {
	handler := http.NewServeMux()
	handler.Handle("/", http.HandlerFunc(index))
	handler.Handle("/sendMessage", http.HandlerFunc(sendBotMessage))
	handler.Handle("/sendClientMessage", http.HandlerFunc(sendMessage))
	handler.Handle("/getBotMessages", http.HandlerFunc(getMessages))
	handler.Handle("/getUpdates", http.HandlerFunc(getUpdates))

	s := &http.Server{
		Addr:           ":8080",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
