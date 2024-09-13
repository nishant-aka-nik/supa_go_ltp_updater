package filter

import (
	"log"
	"supa_go_ltp_updater/config"
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

func Alert(latestStocksData []model.Stock, crossMatchedStocks []model.BreakoutFilter) {
	var crossMatchedStockMap = make(map[string]model.BreakoutFilter)
	for _, record := range crossMatchedStocks {
		crossMatchedStockMap[record.Symbol] = record
	}

	var filteredStockSlice []model.Stock

	for _, stock := range latestStocksData {
		if _, ok := crossMatchedStockMap[stock.Symbol]; !ok {
			continue
		}

		entryStart := crossMatchedStockMap[stock.Symbol].GetEntry()
		id := crossMatchedStockMap[stock.Symbol].ID
		highCloseDiff := stock.GetPercentageDifferenceBetweenHighAndClose()
		openCloseDiff := stock.GetPercentageDifferenceBetweenOpenAndClose()

		entry := stock.Close > entryStart
		volume := stock.GetVolumeTimes() > 0.2
		candleGreaterThanWick := openCloseDiff > highCloseDiff
		greenCandle := openCloseDiff > 0

		if entry && volume && candleGreaterThanWick && greenCandle {
			filteredStockSlice = append(filteredStockSlice, stock)

			//---------------
			entryPrice := stock.Close
			entryDate := stock.FormatDate(stock.Date)
			supabase.MarkEntry(entryPrice, entryDate, id, config.AppConfig.TableNames.BreakoutFilter)

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
		notification.SendMails(notification.GetEntryEmailList(filteredStockSlice, "Breakout Entry"))
	}
}

func FilterCrossMatchStocks(latestStocksData []model.Stock) []model.BreakoutFilter {
	var filteredStocks []model.BreakoutFilter

	for i := 0; i < len(latestStocksData); i++ {
		highCloseDiff := latestStocksData[i].GetPercentageDifferenceBetweenHighAndClose()
		price52WeekHighDiff := latestStocksData[i].GetPercentagePriceAwayFrom52WeekHigh()
		openCloseDiff := latestStocksData[i].GetPercentageDifferenceBetweenOpenAndClose()

		crossMatch := highCloseDiff == price52WeekHighDiff

		volumeTimes := latestStocksData[i].GetVolumeTimes()

		candleGreaterThanWick := openCloseDiff > highCloseDiff

		if crossMatch && volumeTimes > 1.5 && openCloseDiff > 2 && candleGreaterThanWick {
			crossMatchedStock := model.BreakoutFilter{
				Symbol:          latestStocksData[i].Symbol,
				CrossMatch:      true,
				CrossMatchPivot: latestStocksData[i].High,
				CrossMatchDate:  latestStocksData[i].FormatDate(latestStocksData[i].Date),
				VolumeTimes:     latestStocksData[i].GetVolumeTimes(),
			}
			filteredStocks = append(filteredStocks, crossMatchedStock)
		}

	}

	return filteredStocks
}

func Reset(latestStocksData []model.Stock, enteredStocks []model.BreakoutFilter) []model.BreakoutFilter {
	// make entry and active both false on stoploss hit
	var latestStocksDataMap = make(map[string]model.Stock)
	for _, record := range latestStocksData {
		latestStocksDataMap[record.Symbol] = record
	}

	var resetStocks []model.BreakoutFilter

	for _, stock := range enteredStocks {
		LTP := latestStocksDataMap[stock.Symbol].Close

		price4PercentAboveEntry := LTP > (stock.EntryPrice * 1.04)

		if price4PercentAboveEntry {
			stock.Entry = false
			resetStocks = append(resetStocks, stock)
		}
	}

	return resetStocks
}
