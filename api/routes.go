package main

import "net/http"

// Route structure defining routing process
type Route struct {
	Name               string
	Method             string
	Pattern            string
	TokenProtected     bool
	HTTPBasicProtected bool
	HandlerFunc        http.HandlerFunc
}

// Routing type to serve slice of routes
type Routing []Route

var routing = Routing{
	Route{
		"Index",
		"GET",
		"/",
		false,
		false,
		Index,
	},
	Route{
		"Users",
		"GET",
		"/user",
		false,
		false,
		GetUsers,
	},
	Route{
		"User",
		"GET",
		"/user/{id}",
		false,
		false,
		GetUserByID,
	},
	Route{
		"addUser",
		"OPTIONS",
		"/user",
		false,
		false,
		AddUser,
	},
	Route{
		"addUser",
		"POST",
		"/user",
		false,
		true,
		AddUser,
	},
	Route{
		"token",
		"OPTIONS",
		"/token",
		false,
		false,
		SetToken,
	},
	Route{
		"token",
		"POST",
		"/token",
		false,
		false,
		SetToken,
	},
}
