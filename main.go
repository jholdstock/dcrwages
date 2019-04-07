package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jholdstock/dcrwages/poloniex"
	"github.com/jholdstock/dcrwages/server"
)

// Scan from the present month until the month specified below.
// For production these values should be 6 and 2016.
// BTC/DCR data is not available on Polo before this time.
const earliestMonth = 6
const earliestYear = 2016

// API settings
const listen = ":3000"

func populateDataModel() {
	now := time.Now()
	month, year := now.Month(), now.Year()

	// Initialise API data model
	server.FullHistory = server.PriceHistory{
		Years: map[int]server.Year{
			year: {
				Months: map[int]server.Month{},
			},
		},
	}

	// Starting with the current month, calculate monthly average
	// prices until the end date specified in config
	completeMonth := false
	for {
		// Get the month's average USDT/DCR price
		average, err := poloniex.GetMonthAverage(month, year)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%.4f USDT/DCR (%d-%d)", average, month, year)
		fmt.Println()

		// Store the price in the API data model
		server.FullHistory.Years[year].Months[int(month)] = server.Month{
			AveragePrice:  average,
			CompleteMonth: completeMonth,
		}

		// Stop if month/year specified in config
		if month == earliestMonth && year == earliestYear {
			break
		}

		// Proceed to the next month
		completeMonth = true
		month--
		// If required, roll over to a new year
		if month == 0 {
			month = 12
			year--
			server.FullHistory.Years[year] = server.Year{
				Months: map[int]server.Month{},
			}
		}
	}
}

func main() {
	populateDataModel()

	// Start API server
	var router = server.NewRouter()

	fmt.Printf("Starting API server on \"%s\"", listen)
	fmt.Println()

	log.Fatal(http.ListenAndServe(listen, router))
}
