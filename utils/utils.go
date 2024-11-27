package utils

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

func FormatToTwoDecimalPlaces(value float64) float64 {
	formattedValue, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return formattedValue
}

// GetISTTime returns the current time in IST (Indian Standard Time)
func GetISTTime() time.Time {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Fatal(err)
	}
	return time.Now().In(loc)
}

func PercentageIncrease(price float64, percentage float64) float64 {
	price = price + (price * (percentage / 100))
	return price
}

func IsSafeTimeToRun() bool {
	currentTime := GetISTTime()

	thresholdTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 15, 39, 0, 0, currentTime.Location())

	safeTime := currentTime.After(thresholdTime)

	return safeTime
}
