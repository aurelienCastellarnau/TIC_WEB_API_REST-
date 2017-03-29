package logger

import (
	"fmt"
	"net/http"
	"time"
)

// Logger kind of dependancy injection into the routing/server process
func Logger(process http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\n\n\n[Handlers stack trace]NEW REQUEST <Handler.Logger()> :")
		start := time.Now()
		fmt.Printf(
			"\n\nMethod: %s\tURI: %s\tName: %s\tSince: %s\n",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
		process.ServeHTTP(w, r)
	})
}
