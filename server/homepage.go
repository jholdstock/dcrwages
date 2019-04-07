package server

import (
	"html/template"
	"net/http"
)

type homepageContext struct {
	FullHistory PriceHistory
	MonthNames  []string
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
		FullHistory: FullHistory,
		MonthNames:  monthNames,
	})
}
