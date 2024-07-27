package utils

import (
	"fmt"
	"strconv"
)

func FormatToTwoDecimalPlaces(value float64) float64 {
	formattedValue, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return formattedValue
}
