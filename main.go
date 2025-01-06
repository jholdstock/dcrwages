package main

import (
	"log"
	"time"

	"github.com/jholdstock/dcrwages/model"
	"github.com/jholdstock/dcrwages/server"
)

// For production these values should be 6 and 2016.
// BTC/DCR data is not available on Polo before this time.
const earliestMonth = 6
const earliestYear = 2016

const refreshRate = 30 * time.Minute
const listen = ":3000"

func main() {
	// Load data model.
	go model.Init(earliestMonth, earliestYear)

	// Start a ticker to update data model.
	ticker := time.NewTicker(refreshRate)
	go func() {
		for {
			<-ticker.C
			model.Refresh()
		}
	}()

	// Start HTTP server.
	log.Printf("Listening on %s", listen)
	log.Fatal(server.NewRouter().Run(listen))
}
