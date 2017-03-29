package routing

import "rest/api/handlers"

// GetRoads roads on GET request
var GetRoads = Routing{
	Route{
		"Index",
		"GET",
		"/*",
		false,
		false,
		false,
		handlers.Index,
	},
	Route{
		"Users",
		"GET",
		"/users",
		false,
		true,
		false,
		handlers.GetUsers,
	},
	Route{
		"User",
		"GET",
		"/user/{id}",
		false,
		true,
		false,
		handlers.GetUserByID,
	},
	Route{
		"Search",
		"GET",
		"/search/{q}",
		false,
		true,
		false,
		handlers.Search,
	},
	Route{
		"Search",
		"GET",
		"/search/{q}/{count}",
		false,
		true,
		false,
		handlers.Search,
	},
}
