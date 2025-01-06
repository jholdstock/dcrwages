package poloniex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// endPoloDate is the last date to use hard-coded Poliniex price data (1/4/2019).
var endPoloDate = time.Date(2019, 4, 1, 0, 0, 0, 0, time.UTC)

// endBnDualDate is the last date to use hard-coded Binance dual market price
// data (1/12/2024).
var endBnDualDate = time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)

var bnDualPrices = map[int]map[int]float64{
	2024: {
		11: 13.7122090566241415,
		10: 12.3452206510161044,
		9:  11.9598676192569346,
		8:  10.7436660914080520,
		7:  13.1012799749362561,
		6:  17.4775393256896692,
		5:  20.3536546157506741,
		4:  22.5338887631813378,
		3:  26.2812852716170866,
		2:  17.5807608204906849,
		1:  16.1007062786345543,
	},
	2023: {
		12: 15.8902050235771846,
		11: 14.3495915381172026,
		10: 12.7347690405402698,
		9:  13.2293647933044287,
		8:  13.8877126520067051,
		7:  15.4010666387651050,
		6:  14.2512103966331569,
		5:  17.1296217018691372,
		4:  20.2190440490048644,
		3:  20.6897983432358856,
		2:  24.0321691053528959,
		1:  21.1977158893422626,
	},
	2022: {
		12: 20.2237844266543689,
		11: 22.0023382389147102,
		10: 26.2743116883658097,
		9:  27.1921232482801898,
		8:  31.9309107579060623,
		7:  23.9273217386308694,
		6:  28.0567443092857260,
		5:  41.4586762240704729,
		4:  60.6211010870180189,
		3:  57.2239892345201113,
		2:  62.2365498055088580,
		1:  62.6322576366711985,
	},
	2021: {
		12: 79.5948970759228160,
		11: 108.1269316941815788,
		10: 121.5740041806176066,
		9:  139.5594719252921720,
		8:  161.2414497256680477,
		7:  127.4849892393791606,
		6:  131.5153570075625566,
		5:  173.4682224073590078,
		4:  198.5968727470704778,
		3:  161.0123188961222525,
		2:  113.7606320032663092,
		1:  54.2538104507785235,
	},
	2020: {
		12: 31.0702874840878387,
		11: 18.1877500562256920,
		10: 12.0138116782416091,
		9:  13.2622457262482403,
		8:  17.0247038340469921,
		7:  15.1290513375872422,
		6:  16.0516571275069531,
		5:  14.1097535627650945,
		4:  12.3417514417872081,
		3:  13.4009651019791640,
		2:  20.4783306393921798,
		1:  18.0045197493624407,
	},
	2019: {
		4:  24.2556019881414677,
		5:  27.7628186277074640,
		6:  28.8963648391817252,
		7:  28.9711538570402780,
		8:  26.2279185280189786,
		9:  22.0171573523647623,
		10: 15.5914981306509972,
		11: 19.9686529369072403,
		12: 18.3240334397952935,
	},
}

var poloPrices = map[int]map[int]float64{
	2016: {
		6:  1.7324316392854615,
		7:  1.787538249917163,
		8:  1.5856028639980422,
		9:  1.3566497939719884,
		10: 0.9353982669667259,
		11: 0.6621181838638783,
		12: 0.507987091420884,
	},
	2017: {
		1:  0.967719806974505,
		2:  2.3901023080493693,
		3:  6.01272447309967,
		4:  13.426474042737595,
		5:  18.401926466740733,
		6:  35.25252598351194,
		7:  27.511244258640346,
		8:  29.55360321606234,
		9:  32.35313864742267,
		10: 28.681525162593655,
		11: 36.42929582578028,
		12: 72.20130915329341,
	},
	2018: {
		1:  101.27282575174003,
		2:  78.25946287855886,
		3:  57.99056601467777,
		4:  60.945655873985416,
		5:  93.06751886052632,
		6:  88.34326846692787,
		7:  65.41494432920159,
		8:  42.73673228545482,
		9:  38.16362216703373,
		10: 41.03641945403267,
		11: 32.51391916462286,
		12: 17.50659503883639,
	},
	2019: {
		1: 17.05543787850539,
		2: 16.51164586583946,
		3: 18.14373460080878,
	},
}

const binanceURL = "https://api.binance.com"
const httpTimeout = time.Second * 10

// GetMonthAverage returns the average DCR/USDT price for a given month.
func GetMonthAverage(month time.Month, year int) (float64, error) {
	startTime := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endTime := startTime.AddDate(0, 1, 0)

	unixStart := startTime.Unix()
	unixEnd := endTime.Unix()

	// Use hard-coded Poloniex data if start date is before 31/3/2019.
	if startTime.Before(endPoloDate) {
		return poloPrices[year][int(month)], nil
	}

	// Use hard-coded Binance dual market data if start date is before 31/11/2024.
	if startTime.Before(endBnDualDate) {
		return bnDualPrices[year][int(month)], nil
	}

	// Download prices from Binance.
	prices, err := getPrices(unixStart, unixEnd)
	if err != nil {
		return 0, fmt.Errorf("getPrices: %w", err)
	}

	// Calculate and return the average of all prices.
	var average float64
	for _, price := range prices {
		average += price
	}
	average = average / float64(len(prices))

	return average, nil
}

// getPrices contacts the Binance API to download the hourly DCR/USDT price data
// for a given datetime range. Returns a map of unix timestamp => average price.
func getPrices(startDate int64, endDate int64) (map[uint64]float64, error) {
	// Construct HTTP request and set parameters.
	req, err := http.NewRequest(http.MethodGet, binanceURL+"/api/v1/klines", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	const dcrUsdtSymbol = "DCRUSDT"
	q.Set("symbol", dcrUsdtSymbol)
	q.Set("startTime", strconv.FormatInt(startDate*1000, 10))
	q.Set("endTime", strconv.FormatInt(endDate*1000, 10))

	// Request 1 hour intervals since there is a 1000 point limit on requests.
	// 31 Days * 24 Hours = 720 data points.
	q.Set("interval", "1h")
	q.Set("limit", "1000")

	req.URL.RawQuery = q.Encode()

	// Create HTTP client.
	httpClient := http.Client{
		Timeout: httpTimeout,
	}

	// Send HTTP request.
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var chartData [][]json.RawMessage

	// Read response and deserialise JSON into [][]json.RawMessage.
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&chartData)
	if err != nil {
		return nil, err
	}

	prices := make(map[uint64]float64, len(chartData))
	for _, v := range chartData {
		var openTime uint64 // v[0]
		var highStr string  // v[2]
		var lowStr string   // v[3]

		err := json.Unmarshal(v[0], &openTime)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(v[2], &highStr)
		if err != nil {
			return nil, err
		}
		high, err := strconv.ParseFloat(highStr, 64)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(v[3], &lowStr)
		if err != nil {
			return nil, err
		}
		low, err := strconv.ParseFloat(lowStr, 64)
		if err != nil {
			return nil, err
		}
		// Create a map of unix timestamps => average price.
		prices[openTime/1000] = (high + low) / 2
	}
	return prices, nil
}
