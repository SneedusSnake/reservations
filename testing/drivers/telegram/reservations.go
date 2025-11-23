package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/SneedusSnake/Reservations/adapters/driven/clock/cache"
	"github.com/alecthomas/assert/v2"
	"github.com/testcontainers/testcontainers-go"
)

type Message struct {
	Id int `json:"message_id"`
	Text string `json:"text"`
	Chat Chat `json:"chat"`
	From User `json:"from"`
}
type Chat struct {
	Id int `json:"id"`
}

type User struct {
	Id int `json:"id"`
	FirstName string `json:"first_name"`
}

type TelegramDriver struct {
	client *http.Client
	host string
	messageId int
	userId int
	chatId int
	clock cache.CacheClock
	users map[string]int
	appContainer testcontainers.Container
	t *testing.T
}

func NewDriver(client *http.Client, host string,clock cache.CacheClock, app testcontainers.Container, t *testing.T) *TelegramDriver {
	return &TelegramDriver{
		client,
		host,
		0,
		0,
		1234,
		clock,
		make(map[string]int),
		app,
		t,
	}
}

func (d *TelegramDriver) AdminAddsSubject(subject string) {
	msg := Message{
		Text: "/add_subject " + subject,
	}

	d.sendClientMessage(msg)
}

func (d *TelegramDriver) AdminAddsTagsToSubject(subject string, tags ...string) {
	msg := Message{
		Text: fmt.Sprintf("/add_tags %s %s", subject, strings.Join(tags, " ")),
	}

	d.sendClientMessage(msg)
}

func (d *TelegramDriver) UserRequestsSubjectsList() {
	msg := Message{
		Id: d.messageId,
		Text: "/list",
	}

	d.sendClientMessage(msg)
}

func (d *TelegramDriver) UserRequestsSubjectTags(subject string) {
	msg := Message{
		Id: d.messageId,
		Text: "/tags " + subject,
	}

	d.sendClientMessage(msg)
}

func (d *TelegramDriver) UserRequestsReservationForSubject(user string, subject string, minutes int) {
	msg := Message{
		Id: d.messageId,
		Text: "/reserve " + subject + " " + strconv.Itoa(minutes),
		From: User{Id: d.getUserId(user), FirstName: user},
	}

	d.sendClientMessage(msg)
}

func (d *TelegramDriver) UserSeesSubjects(subject ...string) {
	msg := d.getLastBotResponse()

	subjects := strings.Split(msg, "\n")
	for _, s := range subject {
		assert.SliceContains(d.t, subjects, s)
	}
}

func (d *TelegramDriver) UserSeesSubjectTags(tags ...string) {
	msg := d.getLastBotResponse()

	recievedTags := strings.Split(msg, "\n")
	for _, tag := range tags {
		assert.SliceContains(d.t, recievedTags, tag)
	}
}

func (d *TelegramDriver) UserAcquiredReservationForSubject(user string, subject string, until string) {
	msg := d.getLastBotResponse()

	assert.Contains(d.t, msg, subject)
	assert.Contains(d.t, msg, user)
	assert.Contains(d.t, msg, until)
}

func (d *TelegramDriver) SubjectHasAlreadyBeenReservedBy(user string, until string) {
	msg := d.getLastBotResponse()

	assert.Contains(d.t, msg, "Already reserved by")
	assert.Contains(d.t, msg, user)
	assert.Contains(d.t, msg, until)
}

func (d *TelegramDriver) ClockSet(t string) {
	now := time.Now()
	parsed, err := time.Parse(time.TimeOnly, t + ":00")
	if err != nil {
		d.t.Fatal(err)
	}
	year, month, day := now.Date()
	hour, minute, second := parsed.Clock()
	result := time.Date(year, month, day, hour, minute, second, 0, time.Local)

	d.clock.Set(result)
	d.appContainer.CopyFileToContainer(d.t.Context(), d.clock.Path(), d.clock.Path(), 0o666)
}

func (d *TelegramDriver) sendClientMessage(msg Message) {
	d.messageId++
	msg.Id = d.messageId
	msg.Chat = Chat{Id: d.chatId}
	
	encoded, err := json.Marshal(msg)	
	assert.NoError(d.t, err)
	_, err = d.client.Post(fmt.Sprintf("%s/testing/sendClientMessage", d.host), "application/json", bytes.NewBuffer(encoded))
	assert.NoError(d.t, err)
}

func (d *TelegramDriver) getLastBotResponse() string {
	var responseData []struct{
		ChatId int `json:"chat_id"`
		Text string `json:"text"`
	}

	for range 5 {
		r, err := d.client.Get(fmt.Sprintf("%s/testing/getBotMessages", d.host))
		assert.NoError(d.t, err)

		body, err := io.ReadAll(r.Body)
		assert.NoError(d.t, err)
		err = json.Unmarshal(body, &responseData)
		assert.NoError(d.t, err)
	}

	assert.NotEqual(d.t, 0, len(responseData))
	botMessage := responseData[len(responseData) - 1]
	assert.Equal(d.t, d.chatId, botMessage.ChatId)

	return botMessage.Text
}

func (d *TelegramDriver) getUserId(name string) int {
	id, ok := d.users[name]

	if !ok {
		d.userId++
		id = d.userId
		d.users[name] = id
	}

	return id
}
