package telegram

import (
	"fmt"
	"net/http"
	"testing"
	"io"
	"strings"
	"encoding/json"
	"bytes"

	"github.com/alecthomas/assert/v2"
)

type Message struct {
	Id int `json:"message_id"`
	Text string `json:"text"`
	Chat Chat `json:"chat"`
}
type Chat struct {
	Id int `json:"id"`
}

type TelegramDriver struct {
	client *http.Client
	host string
	t *testing.T
	messageId int
}

func NewDriver(client *http.Client, host string, t *testing.T) *TelegramDriver {
	return &TelegramDriver{
		client,
		host,
		t,
		0,
	}
}

func (d *TelegramDriver) AdminAddsSubject(subject string) {
	d.messageId++
	clientMessage := Message{
		Id: d.messageId,
		Text: "/add_subject " + subject,
		Chat: Chat{Id: 1234},
	}
	encoded, err := json.Marshal(clientMessage)	
	assert.NoError(d.t, err)
	_, err = d.client.Post(fmt.Sprintf("%s/testing/sendClientMessage", d.host), "application/json", bytes.NewBuffer(encoded))
	assert.NoError(d.t, err)
}

func (d *TelegramDriver) AdminAddsTagsToSubject(subject string, tags ...string) {
	d.messageId++
	clientMessage := Message{
		Id: d.messageId,
		Text: fmt.Sprintf("/add_tags %s %s", subject, strings.Join(tags, " ")),
		Chat: Chat{Id: 1234},
	}

	encoded, err := json.Marshal(clientMessage)	
	assert.NoError(d.t, err)
	_, err = d.client.Post(fmt.Sprintf("%s/testing/sendClientMessage", d.host), "application/json", bytes.NewBuffer(encoded))
	assert.NoError(d.t, err)
}

func (d *TelegramDriver) UserRequestsSubjectsList() {
	d.messageId++
	clientMessage := Message{
		Id: d.messageId,
		Text: "/list",
		Chat: Chat{Id: 1234},
	}
	encoded, err := json.Marshal(clientMessage)	
	assert.NoError(d.t, err)
	_, err = d.client.Post(fmt.Sprintf("%s/testing/sendClientMessage", d.host), "application/json", bytes.NewBuffer(encoded))
	assert.NoError(d.t, err)
}

func (d *TelegramDriver) UserRequestsSubjectTags(subject string) {
	d.messageId++
	clientMessage := Message{
		Id: d.messageId,
		Text: "/tags " + subject,
		Chat: Chat{Id: 1234},
	}
	encoded, err := json.Marshal(clientMessage)	
	assert.NoError(d.t, err)
	_, err = d.client.Post(fmt.Sprintf("%s/testing/sendClientMessage", d.host), "application/json", bytes.NewBuffer(encoded))
	assert.NoError(d.t, err)
}

func (d *TelegramDriver) UserSeesSubjects(subject ...string) {
	var responseData []struct{
		ChatId int `json:"chat_id"`
		Text string `json:"text"`
	}

	for i := 0; i < 10 && len(responseData) < len(subject); i++ {
		r, err := d.client.Get(fmt.Sprintf("%s/testing/getBotMessages", d.host))
		assert.NoError(d.t, err)

		body, err := io.ReadAll(r.Body)
		assert.NoError(d.t, err)
		err = json.Unmarshal(body, &responseData)
		assert.NoError(d.t, err)
	}

	assert.NotEqual(d.t, 0, len(responseData))
	botMessage := responseData[len(responseData) - 1]

	subjects := strings.Split(botMessage.Text, "\n")
	for _, s := range subject {
		assert.SliceContains(d.t, subjects, s)
	}
	assert.Equal(d.t, 1234, botMessage.ChatId)
}

func (d *TelegramDriver) UserSeesSubjectTags(tags ...string) {
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

	recievedTags := strings.Split(botMessage.Text, "\n")
	for _, tag := range tags {
		assert.SliceContains(d.t, recievedTags, tag)
	}
	assert.Equal(d.t, 1234, botMessage.ChatId)
}
