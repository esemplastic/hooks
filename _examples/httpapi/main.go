package main

import (
	"log"
	"net/http"

	"github.com/esemplastic/hooks/_examples/httpapi/hub"
	_ "github.com/esemplastic/hooks/_examples/httpapi/routes"
)

func init() {
	hub.RegisterHook(hub.AddRoute, addRoute)
}

func main() {
	host := "localhost:8080"
	log.Printf("HTTP Server listening on  http://%s", host)
	if err := http.ListenAndServe(host, mux); err != nil {
		log.Fatal(err)
	}
}

var mux = http.NewServeMux()

func addRoute(method string, path string, handler http.HandlerFunc) {
	routeHandler := wrapHandler(method, path, handler)
	mux.Handle(path, routeHandler)
}

func wrapHandler(method string, path string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		handler(w, r)
	}
}
