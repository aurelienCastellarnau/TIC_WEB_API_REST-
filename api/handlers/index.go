package handlers

import "net/http"

// Index handler nothing to do with it
func Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", 200)
}
