package supabase

import (
	"log"
	"strconv"
	"supa_go_ltp_updater/config"
	"supa_go_ltp_updater/model"

	"github.com/nedpals/supabase-go"
)

func LtpUpdater(stocksData []model.Stock, tableName string) {
	client := GetSupabaseClient()

	var recordsInTable []model.Stock
	//get all table records
	selectErr := client.DB.From(tableName).Select("symbol").Execute(&recordsInTable)
	if selectErr != nil {
		log.Fatalf("Error selecting table: %v", selectErr)
	}

	var symbolsMap = make(map[string]struct{})
	for _, record := range recordsInTable {
		symbolsMap[record.Symbol] = struct{}{}
	}

	// Perform the update operation
	for _, record := range stocksData {
		payload := map[string]interface{}{
			"close":            record.Close,
			"change_pct":       record.ChangePercentage,
			"symbol":           record.Symbol,
			"high52":           record.High52,
			"high":             record.High,
			"low":              record.Low,
			"open":             record.Open,
			"date":             record.Date,
			"daily_avg_volume": record.DailyAvgVolume,
			"volume":           record.Volume,
		}

		if _, ok := symbolsMap[record.Symbol]; !ok {
			log.Printf("Inserting record: %v\n", record)
			var result []map[string]interface{}
			err := client.DB.From(tableName).Insert(payload).Execute(&result)
			if err != nil {
				log.Fatalf("Error updating table: %v", err)
			}
			continue
		}

		// FIXME: this needs to be optimised it is very slow 24s a lot
		log.Printf("Updating record: %v\n", record)
		var result []map[string]interface{}
		err := client.DB.From(tableName).Update(payload).Eq("symbol", record.Symbol).Execute(&result)
		if err != nil {
			log.Fatalf("Error updating table: %v for records: %#v", err, record)
		}

		// TODO: add deletion also supabase should mirror the google sheet
	}
}

func GetSupabaseClient() *supabase.Client {
	//FIXME: how to close supabase client as we are opening but not closing
	supabaseUrl := config.AppConfig.Supabase.SupabaseUrl
	supabaseKey := config.AppConfig.Supabase.SupabaseKey

	client := supabase.CreateClient(supabaseUrl, supabaseKey)
	return client
}

func GetLogsFromSupbase() []model.SwingLog {
	table := "swinglog"
	client := GetSupabaseClient()

	var recordsInTable []model.SwingLog
	//get all table records
	selectErr := client.DB.From(table).Select("*").Execute(&recordsInTable)
	if selectErr != nil {
		log.Fatalf("Error selecting table:%v Err: %v", table, selectErr)
	}

	table = "accounts"
	var accounts []model.Account
	selectErr = client.DB.From(table).Select("*").Execute(&accounts)
	if selectErr != nil {
		log.Fatalf("Error selecting table:%v Err: %v", table, selectErr)
	}

	// Create a map to store accounts by UserID
	accountMap := make(map[string]model.Account)
	for _, account := range accounts {
		accountMap[account.UserID] = account
	}

	// Iterate through recordsInTable and set the Account
	for i, record := range recordsInTable {
		if account, exists := accountMap[record.UserID]; exists {
			recordsInTable[i].Account = account
		} else {
			// Handle the case where the account is not found
			log.Printf("Account not found for UserID: %s\n", record.UserID)
			// Optional: set a default value or take other actions
			recordsInTable[i].Account = model.Account{} // or some default value
		}
	}

	return recordsInTable
}

func TargetUpdater(stocksData []model.SwingLog, tableName string) {
	if len(stocksData) == 0 {
		log.Println("No records to update")
		return
	}

	client := GetSupabaseClient()
	// Perform the update operation
	for _, record := range stocksData {
		payload := map[string]interface{}{
			"stoploss": record.Stoploss,
			"target":   record.Target,
			"pivot":    record.Pivot,
		}

		log.Printf("Updating target for record: %v\n", record)
		var result []map[string]interface{}
		err := client.DB.From(tableName).Update(payload).Eq("id", strconv.FormatInt(int64(record.ID), 10)).Execute(&result)
		if err != nil {
			log.Printf("Error in target updating table: %v for records: %#v", err, record)
		}
	}
}
