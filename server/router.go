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
}

var apiRoutes = []route{
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

	router.PathPrefix("/css/").Handler(
		http.StripPrefix("/css", http.FileServer(http.Dir("server/public/css"))))
	router.PathPrefix("/js/").Handler(
		http.StripPrefix("/js", http.FileServer(http.Dir("server/public/js"))))
	router.PathPrefix("/images/").Handler(
		http.StripPrefix("/images", http.FileServer(http.Dir("server/public/images"))))

	// API router
	for _, route := range apiRoutes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = logger(handler, route.Name)
		handler = initChecker(handler)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	// HTML routes
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

// initChecker is a HTTP handler for the JSON API. It returns a 503
// HTTP status if the data model is not yet loaded.
func initChecker(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(FullHistory.Years) == 0 {
			writeJSONResponse(w,
				http.StatusServiceUnavailable,
				errmsg{
					"dcrwages is initialising",
				},
			)
		} else {
			inner.ServeHTTP(w, r)
		}
	})
}
