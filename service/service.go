package service

import (
	"fmt"
	"supa_go_ltp_updater/stocks"
)

func CronLtpUpdater() {

	stocksData := stocks.GetStocks()
	fmt.Print(stocksData)

}
