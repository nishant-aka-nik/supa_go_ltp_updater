package watch

import (
	"fmt"
	"log"
	"supa_go_ltp_updater/config"
	"supa_go_ltp_updater/model"
	"supa_go_ltp_updater/notification"
	"supa_go_ltp_updater/supabase"
)

func StoplossHit(stocksData []model.Stock, symbolToLtpMap map[string]float64, swingLogs []model.SwingLog) {
	emailist := notification.EmailList{}

	for _, swingLog := range swingLogs {
		ltp, ok := symbolToLtpMap[swingLog.Symbol]
		if !ok {
			log.Printf("LTP not found for %s", swingLog.Symbol)
			continue
		}

		if ltp < swingLog.Stoploss {
			email := notification.Email{
				To:      swingLog.Account.UserEmail,
				Subject: "Stoploss Hit - Argus Alerts",
				Body:    fmt.Sprintf("Stoploss hit for %s", swingLog.Symbol),
			}
			emailist = emailist.PushEmail(email)

			if swingLog.Account.SecondaryEmail != "" {
				email := notification.Email{
					To:      swingLog.Account.UserEmail,
					Subject: "Stoploss Hit - Argus Alerts",
					Body:    fmt.Sprintf("Stoploss hit for %s", swingLog.Symbol),
				}
				emailist = emailist.PushEmail(email)
			}
		}
	}

	notification.SendMails(emailist.GetEmails())
}

func TargetHit(stocksData []model.Stock, symbolToLtpMap map[string]float64, swingLogs []model.SwingLog) {
	emailist := notification.EmailList{}

	swingLogsToBeUpdated := []model.SwingLog{}

	for _, swingLog := range swingLogs {
		ltp, ok := symbolToLtpMap[swingLog.Symbol]
		if !ok {
			log.Printf("LTP not found for %s", swingLog.Symbol)
			continue
		}

		// check if target hit
		if ltp > swingLog.Target {
			newStoploss := ltp - (ltp * (config.AppConfig.StoplossPercentage / 100))
			swingLog.Stoploss = newStoploss
			newTarget := ltp + (ltp * (config.AppConfig.TargetPercentage / 100))
			swingLog.Target = newTarget
			//FIXME: pivot is not used by ui yet
			swingLog.Pivot = ltp

			emails := []string{swingLog.Account.UserEmail, swingLog.Account.SecondaryEmail}

			for _, email := range emails {
				if email == "" {
					continue
				}
				email := notification.Email{
					To:      email,
					Subject: "Target Hit - Argus Alerts",
					Body:    fmt.Sprintf("Target hit for %v at %v \nNew Stoploss: %v (10 percent) \nNew Target: %v (5 percent)", swingLog.Symbol, ltp, swingLog.Stoploss, swingLog.Target),
				}
				emailist = emailist.PushEmail(email)
			}

			swingLogsToBeUpdated = append(swingLogsToBeUpdated, swingLog)

		}
	}

	notification.SendMails(emailist.GetEmails())
	supabase.TargetUpdater(swingLogsToBeUpdated, "swinglog")

}

func GetSymbolToLtpMap(stocksData []model.Stock) map[string]float64 {
	var symbolToLtpMap = make(map[string]float64)
	for _, stock := range stocksData {
		symbolToLtpMap[stock.Symbol] = stock.Close
	}
	return symbolToLtpMap
}
