package server

import (
	"html/template"
	"net/http"
)

type homepageContext struct {
	FullHistory PriceHistory
}

func homePage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("server/views/index.html")
	t.Execute(w, homepageContext{
		FullHistory: FullHistory,
	})
}
