package main

import (
	"net/http"

	"github.com/esemplastic/hooks/_examples/http/hub"
	_ "github.com/esemplastic/hooks/_examples/http/server"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Hello from index: /"))
	})

	hub.Notifier.Notify(hub.Serve, "localhost:8080", mux)
}
