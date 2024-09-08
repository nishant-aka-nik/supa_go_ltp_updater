package model

import (
	"log"
	"time"
)

type GapFilter struct {
	ID          int     `json:"id"`
	Date        string  `json:"date"`
	Symbol      string  `json:"symbol"`
	GapPivot    float64 `json:"gap_pivot"`
	VolumeTimes float64 `json:"volume_times"`
	Entry       bool    `json:"entry"`
}

func (s GapFilter) FormatDate(dateStr string) string {
	// Parse the input date string
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Fatalf("Error parsing date: %v", err)
	}
	// Format it in YYYY-MM-DD
	return t.Format("2006-01-02")
}
