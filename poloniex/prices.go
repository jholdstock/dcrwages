package poloniex

import (
	"encoding/json"
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

// GetPrices contacts the Poloniex API to download
// price data for a given CC pairing. Returns a map
// of unix timestamp => average price
func GetPrices(pairing string, startDate int64, endDate int64) (map[uint64]float64, error) {
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
