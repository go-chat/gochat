package handler

import (
	"net/http"

	"github.com/go-chat/gochat/middleware"
	"github.com/gorilla/mux"
)

// Route defines indeed information for single route
type Route struct {
	Method      string
	Pattern     string
	Name        string
	HandlerFunc http.HandlerFunc
}

// NewRouter creates a new mux router
func NewRouter() *mux.Router {
	var routes []Route

	routes = append(routes, Route{"GET", "/index", "Index", Index})

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = middleware.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
