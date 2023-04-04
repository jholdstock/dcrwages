package model

import (
	"fmt"
	"log"
	"time"

	"github.com/jholdstock/dcrwages/poloniex"
)

type Month struct {
	AveragePrice  float64 `json:"avg_usdt_dcr_price"`
	CompleteMonth bool    `json:"complete_month"`
}

type Months map[int]Month

type Year struct {
	Months `json:"months"`
}
type Years map[int]Year

type PriceHistory struct {
	Years `json:"years"`
}

func (r PriceHistory) FindYear(yearParam int) (*Year, error) {
	year, found := r.Years[yearParam]
	if !found {
		return nil, fmt.Errorf("no data for year %d", yearParam)
	}

	return &year, nil
}

func (r PriceHistory) FindMonth(yearParam int, monthParam int) (*Month, error) {
	year, err := r.FindYear(yearParam)
	if err != nil {
		return nil, err
	}

	month, found := year.Months[monthParam]
	if !found {
		return nil, fmt.Errorf("no data for month %d of year %d", monthParam, yearParam)
	}

	return &month, nil
}

// FullHistory contains all of the historical price data
var FullHistory = PriceHistory{
	Years: make(map[int]Year),
}

// Initialised is set to true when all historic data has
// been retreived and processed
var Initialised = false
var LastUpdated time.Time

// Init will initialise the data model with historic price data
func Init(earliestMonth time.Month, earliestYear int) {
	log.Println("Retrieving historic data...")
	now := time.Now()
	month, year := now.Month(), now.Year()

	// Starting with the current month, calculate monthly average
	// prices until the end date
	completeMonth := false
	for {
		err := storeMonthInModel(month, year, completeMonth)
		if err != nil {
			log.Fatal("ERROR: Failed to retrieve history.")
		}
		// Stop if month/year
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
		}
	}

	Initialised = true
	log.Println("Initialisation completed.")
}

// Refresh will update price data for the current month.
// If a new month has just started, it will update the value for
// the previous month and set complete_month=true
func Refresh() {
	log.Println("Refreshing model.")

	now := time.Now()
	month, year := now.Month(), now.Year()

	// Check if data for the current month already exists
	if _, ok := FullHistory.Years[year].Months[int(month)]; !ok {
		// A new month has just started.
		log.Printf("%2.d-%d has just started", month, year)

		// Update month which just finished
		previousMonth, previousYear := month, year
		previousMonth--
		if previousMonth == 0 {
			previousMonth = 12
			previousYear--
		}

		_ = storeMonthInModel(previousMonth, previousYear, true)
	}

	_ = storeMonthInModel(month, year, false)
}

func storeMonthInModel(month time.Month, year int, completeMonth bool) error {
	// Get the month's average USDT/DCR price
	average, err := poloniex.GetMonthAverage(month, year)
	if err != nil {
		log.Printf("ERROR: Failed to update month %2.d-%d", month, year)
		log.Printf("%v", err)
		return err
	}

	// Create the year if it doesn't already exist
	if _, ok := FullHistory.Years[year]; !ok {
		FullHistory.Years[year] = Year{
			Months: make(map[int]Month),
		}
	}

	m := Month{
		AveragePrice:  average,
		CompleteMonth: completeMonth,
	}

	// Store the price in the data model
	log.Printf("Storing rate for %2.d-%d: %.4f USDT/DCR", month, year, average)
	LastUpdated = time.Now()
	FullHistory.Years[year].Months[int(month)] = m
	return nil
}
