package filter

import "supa_go_ltp_updater/model"

func FilterStocks(stocksData []model.Stock) []model.Stock {
	var filteredStocks []model.Stock

	for i := 0; i < len(stocksData); i++ {
		highCloseDiff := stocksData[i].GetPercentageDifferenceBetweenHighAndClose()
		price52WeekHighDiff := stocksData[i].GetPercentagePriceAwayFrom52WeekHigh()
		openCloseDiff := stocksData[i].GetPercentageDifferenceBetweenOpenAndClose()

		crossMatch := highCloseDiff == price52WeekHighDiff

		volumeTimes := stocksData[i].GetVolumeTimes()

		if crossMatch && volumeTimes > 1.5 && openCloseDiff > 2 && highCloseDiff < 2.5 {
			filteredStocks = append(filteredStocks, stocksData[i])
		}

	}

	return filteredStocks
}
