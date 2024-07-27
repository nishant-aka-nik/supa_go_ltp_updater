package model

import (
	"supa_go_ltp_updater/utils"
	"time"
)

type Stock struct {
	Symbol           string    `json:"symbol"`
	ChangePercentage float64   `json:"change_pct"`
	High52           float64   `json:"high52"`
	Date             time.Time `json:"date"`
	Open             float64   `json:"open"`
	High             float64   `json:"high"`
	Low              float64   `json:"low"`
	Close            float64   `json:"close"`
	Volume           float64   `json:"volume"`
	DailyAvgVolume   float64   `json:"daily_avg_volume"`
}

func (s Stock) GetVolumeTimes() float64 {
	volumeTimes := s.Volume / s.DailyAvgVolume
	return utils.FormatToTwoDecimalPlaces(volumeTimes)
}

func (s Stock) GetPercentagePriceAwayFrom52WeekHigh() float64 {
	percentage := ((s.High52 - s.Close) / s.High52) * 100
	return utils.FormatToTwoDecimalPlaces(percentage)
}

func (s Stock) GetPercentageDifferenceBetweenHighAndClose() float64 {
	percentage := ((s.High - s.Close) / s.High) * 100
	return utils.FormatToTwoDecimalPlaces(percentage)
}
