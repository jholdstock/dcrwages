package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jholdstock/dcrwages/model"
)

var monthNames = []string{
	"",
	"January",
	"February",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"October",
	"November",
	"December",
}

func homePage(c *gin.Context) {
	duration := time.Since(model.LastUpdated)
	lastUpdated := int(duration.Minutes())

	c.HTML(http.StatusOK, "homepage.html", gin.H{
		"PriceData":    model.FullHistory,
		"Initialised":  model.Initialised,
		"MonthNames":   monthNames,
		"CurrentYear":  time.Now().Year(),
		"CurrentMonth": int(time.Now().Month()),
		"LastUpdated":  fmt.Sprintf("%d minutes ago", lastUpdated),
	})
}
