package server

type PriceMonth struct {
	AveragePrice  float64 `json:"avg_usdt_dcr_price"`
	CompleteMonth bool    `json:"complete_month"`
}

type Months map[int]PriceMonth

type PriceYear struct {
	Months `json:"months"`
}
type Years map[int]PriceYear

type PriceHistory struct {
	Years `json:"years"`
}
