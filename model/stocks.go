package model

import "time"

type Stock struct {
	Symbol              string    `json:"symbol"`
	LTP                 float64   `json:"last_traded_price"`
	ChangePercentage    float64   `json:"change_pct"`
	DataDelay           float64   `json:"data_delay"`
	Low                 float64   `json:"low"`
	High52              float64   `json:"high52"`
	Date                time.Time `json:"date"`
	Open                float64   `json:"open"`
	High                float64   `json:"high"`
	VolumeTimes         float64   `json:"volume_times"`
	PriceAwayFrom52High float64   `json:"price_away_from_52high"`
}
