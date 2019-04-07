package server

import (
	"html/template"
	"net/http"
	"time"
)

type homepageContext struct {
	FullHistory  PriceHistory
	MonthNames   []string
	CurrentYear  int
	CurrentMonth int
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
	t, _ := template.ParseFiles("server/views/index.html")
	t.Execute(w, homepageContext{
		FullHistory:  FullHistory,
		MonthNames:   monthNames,
		CurrentYear:  time.Now().Year(),
		CurrentMonth: int(time.Now().Month()),
	})
}
