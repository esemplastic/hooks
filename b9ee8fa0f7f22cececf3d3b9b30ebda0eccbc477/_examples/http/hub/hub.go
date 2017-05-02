package hub

import (
	"github.com/esemplastic/hooks"
)

const (
	// server
	Serve hooks.State = iota
	Shutdown
	// routing
	Handle
)

var (
	Hub       *hooks.Hub
	Registrar hooks.Registrar
	Notifier  hooks.Notifier
)

func init() {
	Hub = hooks.NewHub()        // full access
	Registrar = Hub.Registrar() // restrict access for registrations only
	Notifier = Hub.Notifier()   // restrict access for notify only
}
