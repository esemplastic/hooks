package routes

import (
	"net/http"

	"github.com/esemplastic/hooks/_examples/httpapi/hub"
)

func init() {
	hub.Notifier.Notify(hub.AddRoute_Hook, "GET", "/", indexHandler)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	contents := []byte("Hello from index!")
	w.Write(contents)
}
