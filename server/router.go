package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var routes = []route{
	{
		"getFullHistory",
		"GET",
		"/prices",
		getFullHistory,
	},
	{
		"getYear",
		"GET",
		"/prices/{year}",
		getYear,
	},
	{
		"getMonth",
		"GET",
		"/prices/{year}/{month}",
		getMonth,
	},
}

// NewRouter initialises a router with routes implementing
// a RESTful HTTP service returning JSON encoded price data
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
