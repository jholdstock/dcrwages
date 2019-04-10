package server

import (
	"github.com/gin-gonic/gin"
)

// NewRouter initialises a router with routes implementing
// a RESTful HTTP service returning JSON encoded price data
func NewRouter() *gin.Engine {
	router := gin.Default()
	router.Static("/public", "./public/")

	router.LoadHTMLGlob("templates/*")

	router.GET("/", homePage)

	api := router.Group("/api")

	api.Use(apiReady())
	{
		api.GET("/prices", getAllData)
		api.GET("/prices/:year", getYear)
		api.GET("/prices/:year/:month", getMonth)
	}

	return router
}
