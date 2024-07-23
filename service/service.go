package service

import (
	"log"
	"supa_go_ltp_updater/stocks"
	"supa_go_ltp_updater/supabase"
	"supa_go_ltp_updater/watch"
	"time"
)

func CronLtpUpdater() {
	start := time.Now()
	log.Printf("Job started at: %s\n", start)
	log.Println("Running CronLtpUpdater")

	// fetch stocks data from google sheets
	stocksData := stocks.GetStocks()
	log.Printf("--------------------------xxx--------------------------")
	log.Printf("Fetched stocks data: %v", stocksData)
	log.Printf("--------------------------xxx--------------------------")

	// update last traded price in supabase
	supabase.LtpUpdater(stocksData)

	// get symbol to ltp map from stocks data
	symbolToLtpMap := watch.GetSymbolToLtpMap(stocksData)

	// get swing logs from supabase
	swingLogs := supabase.GetLogsFromSupbase()

	//check for stoploss hit and send notification
	watch.StoplossHit(stocksData, symbolToLtpMap, swingLogs)

	//check for target hit and send notification
	watch.TargetHit(stocksData, symbolToLtpMap, swingLogs)

	// log execution time
	end := time.Now()
	log.Printf("Job ended at: %s\n", end)
	log.Printf("Job execution time: %v\n", end.Sub(start).Seconds())
	log.Println("Finished running CronLtpUpdater")
}
