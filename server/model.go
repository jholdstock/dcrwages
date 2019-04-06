package server

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

// FullHistory contains all of the historical price data
var FullHistory PriceHistory
