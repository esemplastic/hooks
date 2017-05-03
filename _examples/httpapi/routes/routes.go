package routes

import "github.com/esemplastic/hooks/_examples/httpapi/hub"

func init() {
	hub.AddRoute("GET", "/", indexHandler)
	hub.AddRoute("GET", "/about", aboutHandler)
}
