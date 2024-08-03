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
