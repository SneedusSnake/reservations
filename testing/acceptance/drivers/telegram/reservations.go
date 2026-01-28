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

	"github.com/SneedusSnake/Reservations/internal/adapters/driven/clock/cache"
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

type Reservation struct {
	User string
	Subject string
	Time time.Time
}

type Response struct {
	ChatId int `json:"chat_id"`
	Text string `json:"text"`
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
	responses []Response
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
		[]Response{},
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
	d.waitForBotResponse()
}

func (d *TelegramDriver) UserRequestsSubjectTags(subject string) {
	msg := Message{
		Id: d.messageId,
		Text: "/tags " + subject,
	}

	d.sendClientMessage(msg)
	d.waitForBotResponse()
}

func (d *TelegramDriver) UserRequestsReservationForSubject(user string, subject string, minutes int) {
	msg := Message{
		Id: d.messageId,
		Text: "/reserve " + subject + " " + strconv.Itoa(minutes),
		From: User{Id: d.getUserId(user), FirstName: user},
	}

	d.sendClientMessage(msg)
	d.waitForBotResponse()
}

func (d *TelegramDriver) UserRequestsReservationRemoval(user string, subject string) {
	msg := Message{
		Id: d.messageId,
		Text: "/remove " + subject,
		From: User{Id: d.getUserId(user), FirstName: user},
	}

	d.sendClientMessage(msg)
	d.waitForBotResponse()
}

func (d *TelegramDriver) UserRequestsReservationsList(tags ...string) {
	msg := Message{
		Id: d.messageId,
		Text: "/reserved",
	}

	if len(tags) > 0 {
		msg.Text += " " + strings.Join(tags, " ")
	}

	d.sendClientMessage(msg)
	d.waitForBotResponse()
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

func (d *TelegramDriver) UserSeesReservations(reservations ...string) {
	msg := d.getLastBotResponse()

	listReservations := d.reservationsFromList(msg)
	assert.Equal(d.t, len(reservations), len(listReservations))
	for _, r := range reservations {
		assert.SliceContains(d.t, listReservations, d.reservationFromSpec(r))
	}
}

func (d *TelegramDriver) UserDoesNotSeeReservations(subject string) {
	msg := d.getLastBotResponse()
	seen := false

	listReservations := d.reservationsFromList(msg)
	for _, r := range listReservations {
		if r.Subject == subject {
			seen = true
		}
	}

	assert.False(d.t, seen)
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

func (d *TelegramDriver) CleanUp() {
	d.responses = []Response{}
	d.messageId = 0
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

func (d *TelegramDriver) waitForBotResponse() {
	var responseData []Response
	ticker := time.NewTicker(time.Millisecond*500)
	done := make(chan bool, 1)
	timer := time.NewTimer(time.Second*20)

	go func() {
		<- timer.C
		done <- true
	}()

	for {
		select {
		case <- done:
			ticker.Stop()
			return
		case <- ticker.C:
			r, err := d.client.Get(fmt.Sprintf("%s/testing/getBotMessages", d.host))
			assert.NoError(d.t, err)

			body, err := io.ReadAll(r.Body)
			assert.NoError(d.t, err)
			err = json.Unmarshal(body, &responseData)
			assert.NoError(d.t, err)
			for _, response := range responseData {
				d.responses = append(d.responses, response)
			}

			if len(responseData) > 0 {
				done <- true
			}
		}
	}
}

func (d *TelegramDriver) getLastBotResponse() string {
	assert.NotEqual(d.t, 0, len(d.responses))
	botMessage := d.responses[len(d.responses) - 1]
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

func (d *TelegramDriver) reservationsFromList(list string) []Reservation {
	var reservations []Reservation
	lines := strings.Split(strings.Trim(list, "\n"), "\n")
	
	for _, line := range lines[1:] {
		reservations = append(reservations, d.reservationFromList(line))
	}

	return reservations
}

func (d *TelegramDriver) reservationFromList(r string) Reservation {
	data := strings.Split(r, "\t")
	assert.Equal(d.t, 4, len(data))
	t, err := time.Parse(time.DateTime, data[1])
	assert.NoError(d.t, err)
	
	return Reservation{
		Subject: data[0],
		User: data[3],
		Time: t,
	}
}

func (d *TelegramDriver) reservationFromSpec(r string) Reservation {
	data := strings.Split(r, " ")
	assert.Equal(d.t, 3, len(data))
	specTime := strings.Split(data[2], ":")
	assert.Equal(d.t, 2, len(specTime))
	hours, err := strconv.Atoi(specTime[0])
	assert.NoError(d.t, err)
	minutes, err := strconv.Atoi(specTime[1])
	assert.NoError(d.t, err)
	t := time.Date(
		d.clock.Current().Year(),
		d.clock.Current().Month(),
		d.clock.Current().Day(),
		hours,
		minutes,
		d.clock.Current().Second(),
		d.clock.Current().Nanosecond(),
		d.clock.Current().Location(),
	)
	
	return Reservation{
		Subject: data[1],
		User: data[0],
		Time: t,
	}
}
