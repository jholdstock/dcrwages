package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var routes = []route{
	{
		"homePage",
		"GET",
		"/",
		homePage,
	},
	{
		"getFullHistory",
		"GET",
		"/api/prices",
		getFullHistory,
	},
	{
		"getYear",
		"GET",
		"/api/prices/{year}",
		getYear,
	},
	{
		"getMonth",
		"GET",
		"/api/prices/{year}/{month}",
		getMonth,
	},
}

// NewRouter initialises a router with routes implementing
// a RESTful HTTP service returning JSON encoded price data
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.PathPrefix("/css/").Handler(http.StripPrefix("/css", http.FileServer(http.Dir("server/public/css"))))
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js", http.FileServer(http.Dir("server/public/js"))))
	router.PathPrefix("/images/").Handler(http.StripPrefix("/images", http.FileServer(http.Dir("server/public/images"))))

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
