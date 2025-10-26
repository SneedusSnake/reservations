package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func apiServer(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	requests = append(requests, string(body))
	fmt.Fprint(w, "Test Subject #1\nTest Subject #2\nTest Subject #3");
}

var requests []string

func main () {
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(apiServer),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
