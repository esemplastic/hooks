package hub

import (
	"github.com/esemplastic/hooks"
)

// practise for you:
// move them to a different package named "hub/hooks".
const (
	Serve_Hook    = "SERVE"
	AddRoute_Hook = "ADD_ROUTE"
)

// practise for you:
// move them to different packages named "hub/registry" and "hub/notifier",
// respectfully.
var (
	Registry hooks.Registry // restrict access to register hooks only, listeners.
	Notifier hooks.Notifier // restrict acceess to notify only, callers.
)

func init() {
	hub := hooks.NewHub()
	Registry = hub
	Notifier = hub
}
