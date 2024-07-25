package watch

import (
	"log"
	"supa_go_ltp_updater/model"
)

func StoplossHit(stocksData []model.Stock, symbolToLtpMap map[string]float64, swingLogs []model.SwingLog) {
	for _, swingLog := range swingLogs {
		ltp, ok := symbolToLtpMap[swingLog.Symbol]
		if !ok {
			log.Printf("LTP not found for %s", swingLog.Symbol)
			continue
		}

		if ltp < swingLog.Stoploss {
			// TODO: email notification
			log.Printf("Stoploss hit for %s", swingLog.Symbol)
		}
	}
}

func TargetHit(stocksData []model.Stock, symbolToLtpMap map[string]float64, swingLogs []model.SwingLog) {
	for _, swingLog := range swingLogs {
		ltp, ok := symbolToLtpMap[swingLog.Symbol]
		if !ok {
			log.Printf("LTP not found for %s", swingLog.Symbol)
			continue
		}
		if ltp > swingLog.Target {
			// TODO: email notification
			log.Printf("Target hit for %s", swingLog.Symbol)
		}
	}
}

func GetSymbolToLtpMap(stocksData []model.Stock) map[string]float64 {
	var symbolToLtpMap = make(map[string]float64)
	for _, stock := range stocksData {
		symbolToLtpMap[stock.Symbol] = stock.Close
	}
	return symbolToLtpMap
}
