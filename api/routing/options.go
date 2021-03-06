package routing

import (
	handlers "rest/api/handlers"
)

// Preflights roads for OPTIONS requests
var Preflights = Routing{
	Route{
		"PreflightHome",
		"OPTIONS",
		"/",
		false,
		false,
		false,
		handlers.Index,
	},
	Route{
		"PreflightGetUsers",
		"OPTIONS",
		"/users",
		false,
		false,
		false,
		handlers.Index,
	},
	Route{
		"PreflightGetId",
		"OPTIONS",
		"/user/{id}",
		false,
		false,
		false,
		handlers.Index,
	},
	Route{
		"PreflightHTTPBasicAuth",
		"OPTIONS",
		"/auth",
		false,
		false,
		false,
		handlers.Index,
	},
	Route{
		"PreflightHTTPBasicUnset",
		"OPTIONS",
		"/logout",
		false,
		false,
		false,
		handlers.Index,
	},
	Route{
		"PreflightAddUser",
		"OPTIONS",
		"/user",
		false,
		false,
		false,
		handlers.Index,
	},
	Route{
		"PreflightSearch",
		"OPTIONS",
		"/search/{q}",
		false,
		false,
		false,
		handlers.Index,
	},
	Route{
		"PreflightSearch",
		"OPTIONS",
		"/search/{q}/{count}",
		false,
		false,
		false,
		handlers.Index,
	},
}
