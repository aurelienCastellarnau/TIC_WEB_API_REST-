package routing

import "rest/api/handlers"

// DeleteRoads roads on DELETE requests
var DeleteRoads = Routing{
	Route{
		"deleteUser",
		"DELETE",
		"/user/{id}",
		false,
		true,
		true,
		handlers.DeleteUser,
	},
}
