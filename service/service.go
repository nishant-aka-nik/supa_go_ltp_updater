package service

import (
	"fmt"
	"supa_go_ltp_updater/stocks"
	"supa_go_ltp_updater/supabase"
)

func CronLtpUpdater() {

	stocksData := stocks.GetStocks()
	fmt.Print(stocksData)

	// LtpUpdater(stocksData)
	supabase.LtpUpdater(stocksData)
}
