package handler

import (
	"net/http"

	"github.com/Sirupsen/logrus"
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
	logrus.Info("Serving...")
	var routes []Route

	routes = append(routes, Route{"GET", "/index", "Index", Index})
	routes = append(routes, Route{"POST", "/register", "Register", Register})

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		logrus.Infof("%s %s %s", route.Method, route.Pattern, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
