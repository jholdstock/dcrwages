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
const earliestMonth = 11
const earliestYear = 2018

// API settings
const listen = ":3000"

// getMonthAverage returns the average USD/DCR price for a given month
func getMonthAverage(month time.Month, year int) (float64, error) {
	startTime := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endTime := startTime.AddDate(0, 1, 0)

	unixStart := startTime.Unix()
	unixEnd := endTime.Unix()

	// Download BTC/DCR and USDT/BTC prices from Poloniex
	dcrPrices, err := poloniex.GetPrices("BTC_DCR", unixStart, unixEnd)
	if err != nil {
		return 0, err
	}
	btcPrices, err := poloniex.GetPrices("USDT_BTC", unixStart, unixEnd)
	if err != nil {
		return 0, err
	}

	// Create a map of unix timestamps => average price
	usdtDcrPrices := make(map[uint64]float64)

	// Select only timestamps which appear in both charts to
	// populate the result set. Multiply BTC/DCR rate by
	// USDT/BTC rate to get USDT/DCR rate.
	for timestamp, dcr := range dcrPrices {
		if btc, ok := btcPrices[timestamp]; ok {
			usdtDcrPrices[timestamp] = dcr * btc
		}
	}

	// Calculate and return the average of all USDT/DCR prices
	var average float64
	for _, price := range usdtDcrPrices {
		average += price
	}
	average = average / float64(len(usdtDcrPrices))

	// Polo doesn't like >6 requests per second
	time.Sleep(300 * time.Millisecond)

	return average, nil
}

func main() {
	now := time.Now()
	month, year := now.Month(), now.Year()

	// Initialise API data model
	server.FullHistory = server.PriceHistory{
		Years: map[int]server.PriceYear{
			year: {
				Months: map[int]server.PriceMonth{},
			},
		},
	}

	// Starting with the current month, calculate monthly average
	// prices until the end date specified in config
	completeMonth := false
	for {
		// Get the month's average USDT/DCR price
		average, err := getMonthAverage(month, year)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%.4f USDT/DCR (%d-%d)", average, month, year)
		fmt.Println()

		// Store the price in the API data model
		server.FullHistory.Years[year].Months[int(month)] = server.PriceMonth{
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
			server.FullHistory.Years[year] = server.PriceYear{
				Months: map[int]server.PriceMonth{},
			}
		}
	}

	// Start API server
	var router = server.NewRouter()

	fmt.Printf("Starting API server on \"%s\"", listen)
	fmt.Println()

	log.Fatal(http.ListenAndServe(listen, router))
}
