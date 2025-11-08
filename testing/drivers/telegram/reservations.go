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
	ChatId int `json:"chat.id"`
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
		ChatId: 1234,
		
	}
	encoded, err := json.Marshal(clientMessage)	
	assert.NoError(d.t, err)
	_, err = d.client.Post(fmt.Sprintf("%s/testing/sendClientMessage", d.host), "application/json", bytes.NewBuffer(encoded))
	assert.NoError(d.t, err)
}

func (d *TelegramDriver) UserRequestsSubjectsList() {
	clientMessage := `{
		"chat": {"id": 1234},
		"text": "/list"
	}`
	_, err := d.client.Post(fmt.Sprintf("%s/testing/sendClientMessage", d.host), "application/json", bytes.NewBuffer([]byte(clientMessage)))
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
