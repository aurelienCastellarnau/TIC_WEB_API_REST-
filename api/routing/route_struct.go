package routing

import (
	"net/http"
)

// Route structure defining routing process
type Route struct {
	Name               string
	Method             string
	Pattern            string
	TokenProtected     bool
	HTTPBasicProtected bool
	Reserved           bool
	HandlerFunc        http.HandlerFunc
}

// Routing type to serve slice of routes
type Routing []Route

// MergeRoutes allow to merge slices of Routing into one Routing object
func (r *Routing) MergeRoutes(roadBox []Routing) {
	for _, slice := range roadBox {
		for _, route := range slice {
			*r = append(*r, route)
		}
	}
}
