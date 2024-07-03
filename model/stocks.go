package model

type Stock struct {
	Symbol           string  `json:"symbol"`
	LTP              float64 `json:"last_traded_price"`
	ChangePercentage float64 `json:"change_pct"`
	DataDelay        float64 `json:"data_delay"`
}
