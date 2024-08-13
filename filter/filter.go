package filter

import (
	"log"
	"supa_go_ltp_updater/model"
	"supa_go_ltp_updater/notification"
	"supa_go_ltp_updater/supabase"
	"supa_go_ltp_updater/utils"
)

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

func Alert(latestStocksData []model.Stock, crossMatchedStocks []model.Stock) {
	var crossMatchedStockMap = make(map[string]model.Stock)
	for _, record := range crossMatchedStocks {
		crossMatchedStockMap[record.Symbol] = record
	}

	var filteredStockSlice []model.Stock

	for _, stock := range latestStocksData {
		if _, ok := crossMatchedStockMap[stock.Symbol]; !ok {
			continue
		}

		entryStart := crossMatchedStockMap[stock.Symbol].GetEntry()
		id := crossMatchedStockMap[stock.Symbol].Id
		highCloseDiff := stock.GetPercentageDifferenceBetweenHighAndClose()
		openCloseDiff := stock.GetPercentageDifferenceBetweenOpenAndClose()

		entry := stock.Close > entryStart
		volume := stock.GetVolumeTimes() > 1.5
		candleGreaterThanWick := openCloseDiff > highCloseDiff
		greenCandle := openCloseDiff > 0

		if entry && volume && candleGreaterThanWick && greenCandle {
			filteredStockSlice = append(filteredStockSlice, stock)

			//---------------
			if !crossMatchedStockMap[stock.Symbol].Entry {
				stoploss := stock.Close - (stock.Close * 0.1)
				target := stock.Close + (stock.Close * 0.05)
				supabase.MarkStoplossTargetEntry(stoploss, target, id, "filter_history")
			}

			//---------------
			if len(filteredStockSlice) > 0 {
				today := utils.GetISTTime()
				log.Printf("--------------------------Top Picks for %v:--------------------------", today.Format("02 January 2006"))

				for index, record := range filteredStockSlice {
					log.Printf("--------------------------%v--------------------------", index+1)
					log.Println("Symbol: ", record.Symbol)
					log.Println("Volume Times: ", record.GetVolumeTimes())
					log.Println("Today's Price change percentage: ", record.GetPercentageDifferenceBetweenOpenAndClose())
					log.Println("Percentage difference between high and close: ", record.GetPercentageDifferenceBetweenHighAndClose())
				}

				log.Println("--------------------------xxx--------------------------")
			}
		}
	}

	if len(filteredStockSlice) > 0 {
		notification.SendMails(notification.GetEntryEmailList(filteredStockSlice))
	}
}

func FilterCrossMatchStocks(latestStocksData []model.Stock) []model.Stock {
	var filteredStocks []model.Stock

	for i := 0; i < len(latestStocksData); i++ {
		highCloseDiff := latestStocksData[i].GetPercentageDifferenceBetweenHighAndClose()
		price52WeekHighDiff := latestStocksData[i].GetPercentagePriceAwayFrom52WeekHigh()
		openCloseDiff := latestStocksData[i].GetPercentageDifferenceBetweenOpenAndClose()

		crossMatch := highCloseDiff == price52WeekHighDiff

		volumeTimes := latestStocksData[i].GetVolumeTimes()

		if crossMatch && volumeTimes > 1.5 && openCloseDiff > 2 && highCloseDiff < 2.5 {
			latestStocksData[i].CrossMatch = true
			latestStocksData[i].CrossMatchPivot = latestStocksData[i].High
			filteredStocks = append(filteredStocks, latestStocksData[i])
		}

	}

	return filteredStocks
}

func Reset(latestStocksData []model.Stock, crossMatchedStocks []model.Stock) []model.Stock {
	// make entry and active both false on stoploss hit
	var latestStocksDataMap = make(map[string]model.Stock)
	for _, record := range latestStocksData {
		latestStocksDataMap[record.Symbol] = record
	}

	var resetStocks []model.Stock

	for _, stock := range crossMatchedStocks {
		LTP := latestStocksDataMap[stock.Symbol].Close

		if LTP < stock.Stoploss {
			stock.Entry = false
			stock.CrossMatch = false
			resetStocks = append(resetStocks, stock)
		}
	}

	return resetStocks
}
