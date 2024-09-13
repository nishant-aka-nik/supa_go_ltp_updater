package filter

import (
	"log"
	"supa_go_ltp_updater/utils"
	"time"
)

func IsSafeToGapFilter(currentDateString string, previousDateStr string) bool {
	currentTime := utils.GetISTTime()

	// Parse the date strings
	currentDate, err := time.Parse("02/01/2006", currentDateString)
	if err != nil {
		log.Fatalf("Error parsing current date: %v", err)
	}

	previousDate, err := time.Parse("2006-01-02", previousDateStr)
	if err != nil {
		log.Fatalf("Error parsing previous date: %v", err)
	}

	thresholdTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 15, 59, 0, 0, currentTime.Location())

	safeDate := previousDate.Before(currentDate)
	safeTime := currentTime.After(thresholdTime)
	if safeDate && safeTime {
		return true
	}
	return false
}
