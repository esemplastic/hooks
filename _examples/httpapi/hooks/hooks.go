package hooks

import (
	"net/http"

	"github.com/esemplastic/hooks"
)

var hub = hooks.NewHub()

// Hooks (funcs, so we keep the type safety on the public api) and notifiers

func AddRoute(method string, path string, handler http.HandlerFunc) {
	hub.NotifyFunc(AddRoute, method, path, handler)
}

// Registrar, just conversion in order to import one path for hooks: this.

func Register(hookFunc interface{}, callback interface{}) {
	hub.RegisterHookFunc(hookFunc, callback)
}
