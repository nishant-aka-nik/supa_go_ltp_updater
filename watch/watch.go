package watch

import (
	"fmt"
	"log"
	"supa_go_ltp_updater/model"
	"supa_go_ltp_updater/notification"
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
				Subject: "Stoploss Hit",
				Body:    fmt.Sprintf("Stoploss hit for %s", swingLog.Symbol),
			}
			emailist = emailist.PushEmail(email)

			if swingLog.Account.SecondaryEmail != "" {
				email := notification.Email{
					To:      swingLog.Account.UserEmail,
					Subject: "Stoploss Hit",
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

	for _, swingLog := range swingLogs {
		ltp, ok := symbolToLtpMap[swingLog.Symbol]
		if !ok {
			log.Printf("LTP not found for %s", swingLog.Symbol)
			continue
		}
		if ltp > swingLog.Target {
			email := notification.Email{
				To:      swingLog.Account.UserEmail,
				Subject: "Target Hit",
				Body:    fmt.Sprintf("Target hit for %s", swingLog.Symbol),
			}
			emailist = emailist.PushEmail(email)
		}
	}

	notification.SendMails(emailist.GetEmails())
}

func GetSymbolToLtpMap(stocksData []model.Stock) map[string]float64 {
	var symbolToLtpMap = make(map[string]float64)
	for _, stock := range stocksData {
		symbolToLtpMap[stock.Symbol] = stock.Close
	}
	return symbolToLtpMap
}
