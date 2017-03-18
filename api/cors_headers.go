package main

import (
	"fmt"
	"net/http"
)

// SetCorsHeaders is implemented to authorize
// Cors request from angular client and xhr generally
// see http://stackoverflow.com/questions/12830095/setting-http-headers-in-golang
func SetCorsHeaders(process http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Printf("\n[Handlers stack trace] call of Handler.SetCorsHeaders()")

		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers",
				"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		if r.Method == "OPTIONS" {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		process.ServeHTTP(w, r)
	})
}
