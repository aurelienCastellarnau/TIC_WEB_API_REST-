package routing

import "rest/api/handlers"

// PostRoads roads for POST requests
var PostRoads = Routing{
	Route{
		"HTTPBasicAuth",
		"POST",
		"/auth",
		false,
		false,
		false,
		handlers.Auth,
	},
	Route{
		"HTTPBasicUnset",
		"POST",
		"/logout",
		false,
		true,
		false,
		handlers.RefreshAuth,
	},
	Route{
		"addUser",
		"POST",
		"/user",
		false,
		true,
		true,
		handlers.AddUser,
	},
}
