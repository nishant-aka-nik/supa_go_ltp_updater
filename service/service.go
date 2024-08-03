package service

import (
	"log"
	"supa_go_ltp_updater/filter"
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

	//check for target hit and send notification
	watch.TargetHit(stocksData, symbolToLtpMap, swingLogs)

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
	//FIXME: InsertCrossMatchedStocks should only do the insert operation in supabase
	// FilterCrossMatchStocks shoud do all the calculation
	// only refactor the calculation part from InsertCrossMatchedStocks to FilterCrossMatchStocks
	filterStocks := filter.FilterCrossMatchStocks(latestStocksData)
	// update cross match stocks data in supabase
	supabase.InsertCrossMatchedStocks(filterStocks, "filter_history")

	// Reset Stage
	// filter reset stocks
	resetStocks := filter.Reset(latestStocksData, crossMatchedStocks)
	// update reset stocks data in supabase
	supabase.Reset(resetStocks, "filter_history")

	// log execution time
	end := utils.GetISTTime()
	log.Printf("Job ended at: %s\n", end)
	log.Printf("Job execution time: %v\n", end.Sub(start).Seconds())
	log.Println("Finished running FilterStocks")
}
