package routing

import "rest/api/handlers"

// PutRoads roads on DELETE request
var PutRoads = Routing{
	Route{
		"putUser",
		"PUT",
		"/user/{id}",
		false,
		true,
		true,
		handlers.EditUser,
	},
}
