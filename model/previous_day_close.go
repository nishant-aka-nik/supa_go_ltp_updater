package model

import (
	"log"
	"time"
)

type PreviousDayData struct {
	ID     int     `json:"id"`
	Date   string  `json:"date"`
	Close  float64 `json:"close"`
	Symbol string  `json:"symbol"`
}

func (s PreviousDayData) FormatDate(dateStr string) string {
	// Parse the input date string
	t, err := time.Parse("02/01/2006", dateStr)
	if err != nil {
		log.Fatalf("PreviousDayData model Error parsing date: %v", err)
	}
	// Format it in YYYY-MM-DD
	return t.Format("2006-01-02")
}
