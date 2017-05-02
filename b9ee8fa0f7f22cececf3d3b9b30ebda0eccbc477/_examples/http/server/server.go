package server

import (
	"log"
	"net/http"

	"github.com/esemplastic/hooks"
	"github.com/esemplastic/hooks/_examples/http/hub"
)

func init() {
	hub.Registrar.RegisterHook(hub.Serve, func(p hooks.Payloads) {
		// if p.Len() > 2 || p.Len() == 0 {
		// 	panic("expecting len to be between 1 and 2")
		// }

		addr := p.First().String()
		handler := p.Second().Handler()

		log.Printf("HTTP Server started at host: %s", addr)
		http.ListenAndServe(addr, handler)
	})
}
