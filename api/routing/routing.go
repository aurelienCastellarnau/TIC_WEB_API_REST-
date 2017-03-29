package routing

import (
	"fmt"
	"net/http"
	"rest/api/logger"
	"rest/api/security"

	"github.com/gorilla/mux"
)

// NewRouter Routing constructor
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/", http.FileServer(http.Dir("./")))
	var roadBox []Routing
	var routing Routing
	roadBox = append(roadBox, GetRoads, PostRoads, PutRoads, DeleteRoads, Preflights)
	routing.MergeRoutes(roadBox)
	for _, route := range routing {

		var handler http.Handler

		handler = route.HandlerFunc
		handler = security.SecureAPIWithToken(handler, route.TokenProtected, route.Reserved)
		handler = security.SecureAPIWithBasic(handler, route.HTTPBasicProtected, route.Reserved)
		handler = security.SetCorsHeaders(handler)
		handler = logger.Logger(handler, route.Name)

		fmt.Printf("\n<Routing> Load route %s on %s", route.Method, route.Pattern)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
