package server

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/jholdstock/dcrwages/model"
)

type homepageContext struct {
	PriceData    model.PriceHistory
	Initialised  bool
	MonthNames   []string
	CurrentYear  int
	CurrentMonth int
	LastUpdated  string
}

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

func homePage(w http.ResponseWriter, r *http.Request) {
	duration := time.Since(model.LastUpdated)
	lastUpdated := int(duration.Minutes())

	t, _ := template.ParseFiles("server/views/homepage.html")
	t.Execute(w, homepageContext{
		PriceData:    model.FullHistory,
		Initialised:  model.Initialised,
		MonthNames:   monthNames,
		CurrentYear:  time.Now().Year(),
		CurrentMonth: int(time.Now().Month()),
		LastUpdated:  fmt.Sprintf("%d minutes ago", lastUpdated),
	})
}
