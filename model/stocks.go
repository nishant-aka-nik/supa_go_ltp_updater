package model

import (
	"log"
	"supa_go_ltp_updater/utils"
	"time"
)

type Stock struct {
	Id               int     `json:"id"`
	Symbol           string  `json:"symbol"`
	ChangePercentage float64 `json:"change_pct"`
	High52           float64 `json:"high52"`
	Date             string  `json:"date"`
	Open             float64 `json:"open"`
	High             float64 `json:"high"`
	Low              float64 `json:"low"`
	Close            float64 `json:"close"`
	Volume           float64 `json:"volume"`
	DailyAvgVolume   float64 `json:"daily_avg_volume"`
	CrossMatchPivot  float64 `json:"cross_match_pivot"`
	CrossMatch       bool    `json:"cross_match"`
	Entry            bool    `json:"entry"`
	Target           float64 `json:"target"`
	Stoploss         float64 `json:"stoploss"`
}

func (s Stock) GetVolumeTimes() float64 {
	// this can happen sometime when the symbol change in google finance it will not have a DailyAvgVolume
	// to prevent divide by zero issue 
	if s.DailyAvgVolume == 0 {
		return 0
	}
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

func (s Stock) GetPercentageDifferenceBetweenOpenAndClose() float64 {
	percentage := ((s.Close - s.Open) / s.Open) * 100
	return utils.FormatToTwoDecimalPlaces(percentage)
}

func (s Stock) GetEntry() float64 {
	return s.CrossMatchPivot + (s.CrossMatchPivot * 0.01)
}

func (s Stock) StoplossHit() bool {
	return s.Close < s.Stoploss
}

func (s Stock) TargetHit() bool {
	return s.Close > s.Target
}

func (s Stock) FormatDate(dateStr string) string {
	// Parse the input date string
	t, err := time.Parse("02/01/2006", dateStr)
	if err != nil {
		log.Fatalf("Stock model Error parsing date: %v", err)
	}
	// Format it in YYYY-MM-DD
	return t.Format("2006-01-02")
}
