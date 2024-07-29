package filter

import "supa_go_ltp_updater/model"

func FilterStocks(stocksData []model.Stock) []model.Stock {
	var filteredStocks []model.Stock

	for i := 0; i < len(stocksData); i++ {
		highCloseDiff := stocksData[i].GetPercentageDifferenceBetweenHighAndClose()
		price52WeekHighDiff := stocksData[i].GetPercentagePriceAwayFrom52WeekHigh()
		openCloseDiff := stocksData[i].GetPercentageDifferenceBetweenOpenAndClose()

		crossMatch := highCloseDiff == price52WeekHighDiff

		if crossMatch &&
			price52WeekHighDiff < 1.6 &&
			openCloseDiff > 1.4 {

			filteredStocks = append(filteredStocks, stocksData[i])
		}

	}

	return filteredStocks
}
