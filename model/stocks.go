package model

import "time"

type Stock struct {
	Symbol              string    `json:"symbol"`
	ChangePercentage    float64   `json:"change_pct"`
	DataDelay           float64   `json:"data_delay"`
	High52              float64   `json:"high52"`
	Date                time.Time `json:"date"`
	Open                float64   `json:"open"`
	High                float64   `json:"high"`
	Low                 float64   `json:"low"`
	Close               float64   `json:"close"`
	VolumeTimes         float64   `json:"volume_times"`
	Volume              float64   `json:"volume"`
	DailyAvgVolume      float64   `json:"daily_avg_volume"`
	PriceAwayFrom52High float64   `json:"price_away_from_52high"`
}
