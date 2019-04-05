package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

const poloURL = "https://poloniex.com/public"
const httpTimeout = time.Second * 3
const pricePeriod = 900

type poloChartData struct {
	Date            uint64  `json:"date"`
	WeightedAverage float64 `json:"weightedAverage"`
}

// downloadPrices returns a map of unix timestamps => average price
func downloadPrices(pairing string, startDate int64, endDate int64) (map[uint64]float64, error) {
	// Construct HTTP request and set parameters
	req, err := http.NewRequest(http.MethodGet, poloURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Set("command", "returnChartData")
	q.Set("currencyPair", pairing)
	q.Set("start", strconv.FormatInt(startDate, 10))
	q.Set("end", strconv.FormatInt(endDate, 10))
	q.Set("period", strconv.Itoa(pricePeriod))
	req.URL.RawQuery = q.Encode()

	// Create HTTP client,
	httpClient := http.Client{
		Timeout: httpTimeout,
	}

	// Send HTTP request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Read response and deserialise JSON
	decoder := json.NewDecoder(resp.Body)
	var chartData []poloChartData
	err = decoder.Decode(&chartData)
	if err != nil {
		return nil, err
	}

	// Create a map of unix timestamps => average price
	prices := make(map[uint64]float64, len(chartData))
	for _, data := range chartData {
		prices[data.Date] = data.WeightedAverage
	}

	return prices, nil
}

// getMonthAverage returns the average USD/DCR price for a given month
func getMonthAverage(month time.Month, year int) (float64, error) {
	startTime := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endTime := startTime.AddDate(0, 1, 0)

	unixStart := startTime.Unix()
	unixEnd := endTime.Unix()

	// Download BTC/DCR and USDT/BTC prices from Poloniex
	dcrPrices, err := downloadPrices("BTC_DCR", unixStart, unixEnd)
	if err != nil {
		return -1, err
	}
	btcPrices, err := downloadPrices("USDT_BTC", unixStart, unixEnd)
	if err != nil {
		return -1, err
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

	return average, nil
}

// writeMonthAverage retrieves the average USD/DCR price for a
// given month and then prints it to the terminal
func writeMonthAverage(month time.Month, year int) {
	avg, err := getMonthAverage(month, year)
	if err != nil {
		log.Fatal(err)
	}

	// Polo doesn't like >6 requests per second
	time.Sleep(300 * time.Millisecond)

	fmt.Printf("%.4f USDT/DCR (%d-%d)", avg, month, year)
	fmt.Println()
}

func main() {
	writeMonthAverage(6, 2016)
	writeMonthAverage(7, 2016)
	writeMonthAverage(8, 2016)
	writeMonthAverage(9, 2016)
	writeMonthAverage(10, 2016)
	writeMonthAverage(11, 2016)
	writeMonthAverage(12, 2016)

	writeMonthAverage(1, 2017)
	writeMonthAverage(2, 2017)
	writeMonthAverage(3, 2017)
	writeMonthAverage(4, 2017)
	writeMonthAverage(5, 2017)
	writeMonthAverage(6, 2017)
	writeMonthAverage(7, 2017)
	writeMonthAverage(8, 2017)
	writeMonthAverage(9, 2017)
	writeMonthAverage(10, 2017)
	writeMonthAverage(11, 2017)
	writeMonthAverage(12, 2017)

	writeMonthAverage(1, 2018)
	writeMonthAverage(2, 2018)
	writeMonthAverage(3, 2018)
	writeMonthAverage(4, 2018)
	writeMonthAverage(5, 2018)
	writeMonthAverage(6, 2018)
	writeMonthAverage(7, 2018)
	writeMonthAverage(8, 2018)
	writeMonthAverage(9, 2018)
	writeMonthAverage(10, 2018)
	writeMonthAverage(11, 2018)
	writeMonthAverage(12, 2018)

	writeMonthAverage(1, 2019)
	writeMonthAverage(2, 2019)
	writeMonthAverage(3, 2019)
}
