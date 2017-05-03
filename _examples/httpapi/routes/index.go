package routes

import "net/http"

func indexHandler(w http.ResponseWriter, r *http.Request) {
	contents := []byte("Hello from index!")
	w.Write(contents)
}
