package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type errmsg struct {
	Err string `json:"error"`
}

func writeJSONResponse(w http.ResponseWriter, httpStatus int, i interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(i)
}

func getIntParam(key string, w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	value, err := strconv.Atoi(vars[key])
	if err != nil {
		writeJSONResponse(w,
			http.StatusBadRequest,
			errmsg{
				fmt.Sprintf("Could not decode param {%s}. Should be an integer.", key),
			},
		)
		return -1, err
	}
	return value, nil
}

// Return all available data, json encoded
func getFullHistory(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w,
		http.StatusOK,
		FullHistory)
}

// Return a single year, json encoded
func getYear(w http.ResponseWriter, r *http.Request) {
	year, err := getIntParam("year", w, r)
	if err != nil {
		return
	}

	if _, found := FullHistory.Years[year]; !found {
		writeJSONResponse(w,
			http.StatusNotFound,
			errmsg{
				fmt.Sprintf("No data for year %d", year),
			})

		return
	}

	writeJSONResponse(w,
		http.StatusOK,
		FullHistory.Years[year])
}

// Return a single month, json encoded
func getMonth(w http.ResponseWriter, r *http.Request) {
	year, err := getIntParam("year", w, r)
	if err != nil {
		return
	}

	if _, found := FullHistory.Years[year]; !found {
		writeJSONResponse(w,
			http.StatusNotFound,
			errmsg{
				fmt.Sprintf("No data for year %d", year),
			})

		return
	}

	month, err := getIntParam("month", w, r)
	if err != nil {
		return
	}

	if _, found := FullHistory.Years[year].Months[month]; !found {
		writeJSONResponse(w,
			http.StatusNotFound,
			errmsg{
				fmt.Sprintf("No data for month %d", month),
			})

		return
	}

	writeJSONResponse(w,
		http.StatusOK,
		FullHistory.Years[year].Months[month])
}
