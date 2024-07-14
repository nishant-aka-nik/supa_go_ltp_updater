package service

import (
	"log"
	"supa_go_ltp_updater/stocks"
	"supa_go_ltp_updater/supabase"
	"time"
)

func CronLtpUpdater() {
	start := time.Now()
	log.Printf("Job started at: %s\n", start)
	log.Println("Running CronLtpUpdater")

	stocksData := stocks.GetStocks()
	log.Printf("--------------------------xxx--------------------------")
	log.Printf("Fetched stocks data: %v", stocksData)
	log.Printf("--------------------------xxx--------------------------")

	supabase.LtpUpdater(stocksData)
	end := time.Now()
	log.Printf("Job ended at: %s\n", end)
	log.Printf("Job execution time: %v\n", end.Sub(start).Seconds())
	log.Println("Finished running CronLtpUpdater")
}
