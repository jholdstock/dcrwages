package poloniex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Set the last date to use polo as 4/1/2019.  If any start date requested is
// after that date then use Binance instead.
var endPoloDate = time.Date(2019, 4, 1, 0, 0, 0, 0, time.UTC)

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
		1:  17.05543787850539,
		2:  16.51164586583946,
		3:  18.14373460080878,
		4:  24.22247913957305,
		5:  27.709146359675785,
		6:  28.868318646874954,
		7:  28.961653167721696,
		8:  26.179791372204296,
		9:  21.98694310045591,
		10: 15.56276183286695,
		11: 19.93416568376405,
		12: 18.295713556992688,
	},
}

const binanceURL = "https://api.binance.com"
const httpTimeout = time.Second * 3

const dcrSymbolBinance = "DCRBTC"
const usdtSymbolBinance = "BTCUSDT"

// GetMonthAverage returns the average USD/DCR price for a given month
func GetMonthAverage(month time.Month, year int) (float64, error) {
	startTime := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endTime := startTime.AddDate(0, 1, 0)

	unixStart := startTime.Unix()
	unixEnd := endTime.Unix()

	// Use Poloniex data if start date is before 3/31/19.
	if startTime.Before(endPoloDate) {
		return poloPrices[year][int(month)], nil
	}

	// Download BTC/DCR and USDT/BTC prices from Binance
	dcrPrices, err := getPrices(dcrSymbolBinance, unixStart, unixEnd)
	if err != nil {
		return 0, fmt.Errorf("getPrices %v: %v", dcrSymbolBinance, err)
	}
	btcPrices, err := getPrices(usdtSymbolBinance, unixStart, unixEnd)
	if err != nil {
		return 0, fmt.Errorf("getPricesBinance %v: %v", usdtSymbolBinance, err)
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

// GetPrices contacts the Poloniex API to download
// price data for a given CC pairing. Returns a map
// of unix timestamp => average price
func getPrices(pairing string, startDate int64, endDate int64) (map[uint64]float64, error) {
	// Construct HTTP request and set parameters
	req, err := http.NewRequest(http.MethodGet, binanceURL+"/api/v1/klines", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Set("symbol", pairing)
	q.Set("startTime", strconv.FormatInt(startDate*1000, 10))
	q.Set("endTime", strconv.FormatInt(endDate*1000, 10))

	// Request 1 hour intervals since there is a 1000 point limit on requests
	// 31 Days * 24 Hours = 720 data points
	q.Set("interval", "1h")
	q.Set("limit", "1000")

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

	var chartData [][]json.RawMessage

	// Read response and deserialise JSON into [][]json.RawMessage
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
		// Create a map of unix timestamps => average price
		prices[openTime/1000] = (high + low) / 2
	}
	return prices, nil
}
