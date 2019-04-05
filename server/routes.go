package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// FullHistory contains all of the historical price data
var FullHistory PriceHistory

func writeJSONObject(w http.ResponseWriter, i interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(i)
}

// Return all available data, json encoded
func getFullHistory(w http.ResponseWriter, r *http.Request) {
	writeJSONObject(w, FullHistory)
}

// Return a single year, json encoded
func getYear(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	year, err := strconv.Atoi(vars["year"])
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	writeJSONObject(w, FullHistory.Years[year])
}

// Return a single month, json encoded
func getMonth(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	year, err := strconv.Atoi(vars["year"])
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	month, err := strconv.Atoi(vars["month"])
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	writeJSONObject(w, FullHistory.Years[year].Months[month])
}
