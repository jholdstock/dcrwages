package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jholdstock/dcrwages/model"
	"github.com/jholdstock/dcrwages/server"
)

// Scan from the present month until the month specified below.
// For production these values should be 6 and 2016.
// BTC/DCR data is not available on Polo before this time.
const earliestMonth = 6
const earliestYear = 2016

// API settings
const listen = ":3000"

func main() {
	// Load data model
	go model.Init(earliestMonth, earliestYear)

	// Start a ticker to update data model
	ticker := time.NewTicker(30 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				model.Refresh()
			}
		}
	}()

	// Start HTTP server
	log.Printf("Listening on \"%s\"", listen)
	log.Fatal(http.ListenAndServe(listen, server.NewRouter()))
}
