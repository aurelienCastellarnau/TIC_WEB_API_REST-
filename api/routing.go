package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter Routing constructor
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/", http.FileServer(http.Dir("./")))
	for _, route := range routing {

		var handler http.Handler

		handler = route.HandlerFunc
		handler = SetCorsHeaders(handler)
		handler = SecureAPIWithToken(handler, route.TokenProtected)
		handler = SecureAPIWithBasic(handler, route.HTTPBasicProtected)
		handler = Logger(handler, route.Name)

		fmt.Printf("\n<Routing> Load route %s on %s", route.Method, route.Pattern)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
