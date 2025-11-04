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

type TelegramDriver struct {
	client *http.Client
	host string
	t *testing.T
}

func NewDriver(client *http.Client, host string, t *testing.T) *TelegramDriver {
	return &TelegramDriver{
		client,
		host,
		t,
	}
}

func (d *TelegramDriver) UserRequestsSubjectsList() {
	clientMessage := `{
		"chat": {"id": 1234},
		"text": "/list"
	}`
	_, err := d.client.Post(fmt.Sprintf("%s/sendClientMessage", d.host), "application/json", bytes.NewBuffer([]byte(clientMessage)))
	assert.NoError(d.t, err)
}

func (d *TelegramDriver) UserSeesSubjects(subject ...string) {
	var responseData []struct{
		Text string `json:"text"`
	}

	for i := 0; i < 10 && len(responseData) == 0; i++ {
		r, err := d.client.Get(fmt.Sprintf("%s/getBotMessages", d.host))
		assert.NoError(d.t, err)

		body, err := io.ReadAll(r.Body)
		assert.NoError(d.t, err)
		err = json.Unmarshal(body, &responseData)
		assert.NoError(d.t, err)
	}

	assert.NotEqual(d.t, 0, len(responseData))

	subjects := strings.Split(responseData[len(responseData)-1].Text, "\n")
	assert.Equal(d.t, subject, subjects)
}
