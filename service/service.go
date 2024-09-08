package service

import (
	"log"
	"supa_go_ltp_updater/config"
	"supa_go_ltp_updater/filter"
	"supa_go_ltp_updater/model"
	"supa_go_ltp_updater/notification"
	"supa_go_ltp_updater/stocks"
	"supa_go_ltp_updater/supabase"
	"supa_go_ltp_updater/utils"
	"supa_go_ltp_updater/watch"
)

func CronLtpUpdater() {
	start := utils.GetISTTime()
	log.Printf("Job started at: %s\n", start)
	log.Println("Running CronLtpUpdater")

	// fetch stocks data from google sheets
	stocksData := stocks.GetStocks()
	log.Printf("--------------------------xxx--------------------------")
	log.Printf("Fetched stocks data: %v", stocksData)
	log.Printf("--------------------------xxx--------------------------")

	// get symbol to ltp map from stocks data
	symbolToLtpMap := watch.GetSymbolToLtpMap(stocksData)

	// get swing logs from supabase
	swingLogs := supabase.GetLogsFromSupbase()

	//check for stoploss hit and send notification
	watch.StoplossHit(stocksData, symbolToLtpMap, swingLogs)

	// update todays data in supabase
	supabase.LtpUpdater(stocksData, "todays_data")

	// log execution time
	end := utils.GetISTTime()
	log.Printf("Job ended at: %s\n", end)
	log.Printf("Job execution time: %v\n", end.Sub(start).Seconds())
	log.Println("Finished running CronLtpUpdater")
}

func FilterStocks() {
	start := utils.GetISTTime()
	log.Printf("FilterStocks Job started at: %s\n", start)
	log.Println("Running FilterStocks")

	// fetch stocks data from google sheets
	latestStocksData := stocks.GetStocks()
	log.Printf("--------------------------xxx--------------------------")
	log.Printf("Fetched stocks data: %v", latestStocksData)
	log.Printf("--------------------------xxx--------------------------")

	// get active stocks
	crossMatchedStocks := supabase.GetCrossMatchedStocks("filter_history")

	// Alert Stage
	filter.Alert(latestStocksData, crossMatchedStocks)

	// Insert Stage
	// filter cross match stocks
	filterStocks := filter.FilterCrossMatchStocks(latestStocksData)
	// update cross match stocks data in supabase
	supabase.InsertCrossMatchedStocks(filterStocks, "filter_history")

	// Reset Stage
	// filter reset stocks
	resetStocks := filter.Reset(latestStocksData, crossMatchedStocks)
	// update reset stocks data in supabase
	supabase.Reset(resetStocks, "filter_history")

	notification.SendMails(notification.GetHealthCheckEmailList("FilterStocks"))

	// log execution time
	end := utils.GetISTTime()
	log.Printf("Job ended at: %s\n", end)
	log.Printf("Job execution time: %v\n", end.Sub(start).Seconds())
	log.Println("Finished running FilterStocks")
}

func TargetHitCheckerCron() {
	start := utils.GetISTTime()
	log.Printf("Job started at: %s\n", start)
	log.Println("Running CronLtpUpdater")

	// fetch stocks data from google sheets
	stocksData := stocks.GetStocks()
	log.Printf("--------------------------xxx--------------------------")
	log.Printf("Fetched stocks data: %v", stocksData)
	log.Printf("--------------------------xxx--------------------------")

	// get symbol to ltp map from stocks data
	symbolToLtpMap := watch.GetSymbolToLtpMap(stocksData)

	// get swing logs from supabase
	swingLogs := supabase.GetLogsFromSupbase()

	//check for target hit and send notification
	watch.TargetHit(stocksData, symbolToLtpMap, swingLogs)

	notification.SendMails(notification.GetHealthCheckEmailList("TargetHitCheckerCron"))

	// log execution time
	end := utils.GetISTTime()
	log.Printf("Job ended at: %s\n", end)
	log.Printf("Job execution time: %v\n", end.Sub(start).Seconds())
	log.Println("Finished running CronLtpUpdater")
}

func Gaptor() {
	start := utils.GetISTTime()
	log.Printf("Job started at: %s\n", start)
	log.Println("Running Gaptor")
	log.Printf("start.Weekday(): %#v\n", start.Weekday())

	// fetch stocks data from google sheets
	stocksData := stocks.GetStocks()
	log.Printf("--------------------------xxx--------------------------")
	log.Printf("Fetched stocks data: %v", stocksData)
	log.Printf("--------------------------xxx--------------------------")

	// get previous day close data
	previousDayDataSlice := supabase.GetPreviousDayData()
	log.Printf("previousDayData: %#v\n", previousDayDataSlice)

	//safeguard
	if !filter.IsSafeToGapFilter(stocksData[0].Date, previousDayDataSlice[0].Date) {
		return
	}
	//-----------------

	symbolToPreviousDay3PercentUpCloseMap := make(map[string]model.PreviousDayData)

	for _, previousDayData := range previousDayDataSlice {
		threePercentUpPrice := utils.PercentageIncrease(previousDayData.Close, 3)
		previousDayData.Close = threePercentUpPrice
		symbolToPreviousDay3PercentUpCloseMap[previousDayData.Symbol] = previousDayData
	}

	var gapFilteredData []model.GapFilter

	for _, stockData := range stocksData {
		volumeTimes := stockData.GetVolumeTimes()

		threePercentUpPrice := symbolToPreviousDay3PercentUpCloseMap[stockData.Symbol].Close

		priceGap := stockData.Open > threePercentUpPrice
		volumeGap := volumeTimes > 3

		//GapPivot is the high on the gap day
		if priceGap && volumeGap {
			gapStock := model.GapFilter{
				Date:        stockData.FormatDate(stockData.Date),
				Symbol:      stockData.Symbol,
				GapPivot:    stockData.High,
				VolumeTimes: volumeTimes,
				Entry:       false,
			}

			gapFilteredData = append(gapFilteredData, gapStock)
		}
	}

	//update gapFilteredData into supabase
	supabase.InsertGapFilteredStocks(gapFilteredData, config.AppConfig.TableNames.GapFilter)

	//get all gapFilteredData from supabase and validate stockData.close > gappivot
	// if stockData.close > gappivot trigger email
	gapFilteredStocks := supabase.GetGapFilteredStocks()
	log.Printf("gapFilteredStocks: %#v\n", gapFilteredStocks)

	symbolToGapFilteredMap := make(map[string]model.GapFilter)
	for _, gapFilteredStock := range gapFilteredStocks {
		symbolToGapFilteredMap[gapFilteredStock.Symbol] = gapFilteredStock
	}

	var GapEntryStocks []model.Stock
	for _, stockData := range stocksData {
		gapPivot := symbolToGapFilteredMap[stockData.Symbol].GapPivot

		entry := stockData.Close > (gapPivot + 2)

		if entry {
			GapEntryStocks = append(GapEntryStocks, stockData)
		}
	}

	//alert via email
	if len(GapEntryStocks) > 0 {
		notification.SendMails(notification.GetEntryEmailList(GapEntryStocks, "Gap Entry"))
		notification.SendMails(notification.GetEntryEmailList(GapEntryStocks, "Gap Entry"))
	}

	//update todays data to previous day data
	supabase.PreviousDayDataUpdater(stocksData, config.AppConfig.TableNames.PreviousDayData)

	// log execution time
	end := utils.GetISTTime()
	log.Printf("Job ended at: %s\n", end)
	log.Printf("Job execution time: %v\n", end.Sub(start).Seconds())
	log.Println("Finished running Gaptor")
}
