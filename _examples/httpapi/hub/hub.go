package hub

import (
	"net/http"
	"reflect"
	"runtime"

	"github.com/esemplastic/hooks"
)

// practise for you:
// move them to a different package named "hub/hooks".
// const (
// 	AddRoute_Hook = "ADD_ROUTE"
// )

///TODO: ALL THE BELOW SHOULD BE TRANSFERED TO THE LIB ITSELF.
// With that code, the listeners(hookers(!)) can listen to any function that is being called
// with just one line of code, without refactoring. The most important stuff of my idea is done.

// we need that if we want to remain the statical typing,
// without the need to remove the existing "more dynamic" behavior
// (although we already make checks so the program knows if fails from the beginning)
func nameOfFunc(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

// practise for you:
// move them to different packages named "hub/registry" and "hub/notifier",
// respectfully.
//
// Can be used to restrict access if you wanna transfer the right parts of a hub.
// var (
// 	Registry hooks.Registry // restrict access to register hooks only, listeners.
// 	Notifier hooks.Notifier // restrict acceess to notify only, callers.
// )

var hub = hooks.NewHub()

func RegisterHook(fn interface{}, callback interface{}) {
	hub.RegisterHook(nameOfFunc(fn), callback)
}

func AddRoute(method string, path string, handler http.HandlerFunc) {
	hub.Notify(nameOfFunc(AddRoute), method, path, handler)
}
