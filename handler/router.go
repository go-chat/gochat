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
	Protected   bool
}

// NewRouter creates a new mux router
func NewRouter() *mux.Router {
	logrus.Info("Serving...")
	var routes []Route

	routes = append(routes, Route{"POST", "/register", "Register", Register, false})
	routes = append(routes, Route{"POST", "/login", "Login", Login, false})
	routes = append(routes, Route{"POST", "/users/groups", "GetGroupsByUserID", GetGroupsByUserID, true})
	routes = append(routes, Route{"POST", "/groups/add_members", "AddMembersToGroupChat", AddMembersToGroupChat, true})
	routes = append(routes, Route{"POST", "/groups/remove_members", "RemoveMembersToGroupChat", RemoveMembersToGroupChat, true})
	routes = append(routes, Route{"DELETE", "/groups/:id", "DeleteGroup", DeleteGroupChat, true})
	routes = append(routes, Route{"GET", "/groups/:id/messages", "ListMessagesOfGroup", ListMessagesOfGroup, true})

	// websocket
	routes = append(routes, Route{"GET", "/ws/connect", "WebSocketConnect", WebSocketConnect, true})

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		handler = CORS(handler)

		if route.Protected {
			handler = TokenAuthMiddleware(handler)
		}

		logrus.Infof("%s %s %s", route.Method, route.Pattern, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
