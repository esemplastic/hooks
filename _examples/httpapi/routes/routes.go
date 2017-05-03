package routes

import "github.com/esemplastic/hooks/_examples/httpapi/hooks"

// important step 2 of 2:
func init() {
	hooks.AddRoute("GET", "/", indexHandler)
	hooks.AddRoute("GET", "/about", aboutHandler)
}
