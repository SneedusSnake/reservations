package cmd

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/SneedusSnake/Reservations/testing/specifications"
)

type TelegramViewer struct{
	client http.Client
}

func (v TelegramViewer) List() ([]string, error) {
	r, err := v.client.Get("http://localhost:8080")
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	return strings.Split(string(body), "\n"), nil
}

func TestList(t *testing.T) {
	kill, err := runServer()
	if err != nil {
		t.Fatal(err)
	}

	defer kill()

	specifications.ListSpecification(t, TelegramViewer{
		client: *http.DefaultClient,
	})
}

func runServer() (func (), error) {
	bin, err := buildBinary()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(bin)

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	kill := func () {
		_ = cmd.Process.Kill()
	}

	return kill, waitForServerListening()
}

func waitForServerListening() error {
	port := "8080"
	for i := 0; i < 30; i++ {
		conn, _ := net.Dial("tcp", net.JoinHostPort("localhost", port))
		if conn != nil {
			conn.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("nothing seems to be listening on localhost:%s", port)
}


func buildBinary() (string, error) {
	binName := "../temp/bin/telegram_server"
	build := exec.Command("go", "build", "-o", binName, "./server/telegram")
	if err := build.Run(); err != nil {
			return "", fmt.Errorf("cannot build tool %s: %s", binName, err)
	}
	return binName, nil
}
