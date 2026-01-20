package containers

import (
	"fmt"

	"github.com/testcontainers/testcontainers-go"
)

type StdoutLogConsumer struct{
	Container string
}

func (lc *StdoutLogConsumer) Accept(l testcontainers.Log) {
	fmt.Printf("%s: %s\n", lc.Container, string(l.Content))
}

func Stdout(name string) *StdoutLogConsumer {
	return &StdoutLogConsumer{Container: name}
}
