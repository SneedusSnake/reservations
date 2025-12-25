package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"time"
)

type Chat struct {
	Id int `json:"id"`
}

type User struct {
	Id int `json:"id"`
	IsBot bool `json:"is_bot"`
	FirstName string `json:"first_name"`
	UserName string `json:"username"`
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

type getUpdatesResponse struct {
	OK     bool             `json:"ok"`
	Result []Update `json:"result"`
}

func getMe(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"ok":true,"result":{}}`)
}

func sendBotMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	var message Message

	if r.Header.Get("Content-Type") != "application/json" {
		chatId, _ := strconv.Atoi(r.FormValue("chat_id"))
		message.ChatId = chatId
		message.Text = r.FormValue("text")
	} else {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		err = json.Unmarshal(body, &message)
		if err != nil {
			log.Print(err, string(debug.Stack()))
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	log.Print("recieved bot message: ", message)
	botMessages = append(botMessages, message)
	fmt.Fprint(w, "OK")
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Print(err, string(debug.Stack()))
		w.WriteHeader(http.StatusBadRequest)
	}

	log.Print("Recieved client message", string(body))
	var message UpdateMessage
	message.Id = len(messages) + 1

	err = json.Unmarshal(body, &message)

	if err != nil {
		log.Print(err, string(debug.Stack()))
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

	data, err := json.Marshal(getUpdatesResponse{OK: true, Result: updates})
	
	if err != nil {
		log.Print(err, string(debug.Stack()))
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Fprint(w, string(data))

	updatesCount := len(updates)
	if updatesCount > 0 {
		lastReadId = updates[len(updates) - 1].Id
	}
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
var token string

func main () {
	token = os.Getenv("BOT_TOKEN")
	handler := http.NewServeMux()
	handler.Handle("/", http.HandlerFunc(getMe))
	handler.Handle(url("/getMe"), http.HandlerFunc(getMe))
	handler.Handle(url("/sendMessage"), http.HandlerFunc(sendBotMessage))
	handler.Handle(url("/getUpdates"), http.HandlerFunc(getUpdates))
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

func url(url string) string {
	return fmt.Sprintf("/bot%s%s", token, url)
}
