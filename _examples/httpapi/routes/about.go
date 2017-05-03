package routes

import "net/http"

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	contents := []byte("About page!")
	w.Write(contents)
}
