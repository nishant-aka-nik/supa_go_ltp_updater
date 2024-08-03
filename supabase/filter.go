package supabase

import (
	"log"
	"strconv"
	"supa_go_ltp_updater/model"
	"supa_go_ltp_updater/notification"
	"time"
)

func InsertCrossMatchedStocks(stocksData []model.Stock, tableName string) {
	client := GetSupabaseClient()

	var recordsInTable []model.Stock
	//get all table records
	selectErr := client.DB.From(tableName).Select("symbol").Eq("cross_match", "TRUE").Execute(&recordsInTable)
	if selectErr != nil {
		log.Fatalf("Error selecting table: %v", selectErr)
	}

	var symbolsMap = make(map[string]struct{})
	for _, record := range recordsInTable {
		symbolsMap[record.Symbol] = struct{}{}
	}

	//TODO: filteredStockString to be removed when we will use the better email html template
	var filteredStockString []string
	var filteredStockSlice []model.Stock

	// Perform the update operation
	for _, record := range stocksData {
		filteredStockString = append(filteredStockString, record.Symbol)
		filteredStockSlice = append(filteredStockSlice, record)

		if _, ok := symbolsMap[record.Symbol]; !ok {
			payload := map[string]interface{}{
				"close":             record.Close,
				"change_pct":        record.ChangePercentage,
				"symbol":            record.Symbol,
				"high52":            record.High52,
				"high":              record.High,
				"low":               record.Low,
				"open":              record.Open,
				"date":              record.Date,
				"daily_avg_volume":  record.DailyAvgVolume,
				"volume":            record.Volume,
				"cross_match":       record.CrossMatch,
				"cross_match_pivot": record.CrossMatchPivot,
			}
			log.Printf("Inserting record: %v\n", record)
			var result []map[string]interface{}
			err := client.DB.From(tableName).Insert(payload).Execute(&result)
			if err != nil {
				log.Fatalf("Error updating table: %v", err)
			}
		}
	}

	if len(filteredStockSlice) > 0 {
		today := time.Now()
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

	if len(filteredStockString) > 0 {
		notification.SendMails(notification.GetTopsPicksEmailList(filteredStockString))
	}

	//TODO: add the exit updater also with tax calculation
}

func GetCrossMatchedStocks(tableName string) []model.Stock {
	client := GetSupabaseClient()

	var recordsInTable []model.Stock
	//get all table records
	selectErr := client.DB.From(tableName).Select("*").Eq("cross_match", "TRUE").Execute(&recordsInTable)
	if selectErr != nil {
		log.Fatalf("Error selecting table: %v", selectErr)
	}

	return recordsInTable
}

func MarkStoplossTargetEntry(stoploss float64, target float64, id int, tableName string) {
	client := GetSupabaseClient()

	payload := map[string]interface{}{
		"target":   target,
		"stoploss": stoploss,
		"entry":    "TRUE",
	}

	var result []map[string]interface{}
	err := client.DB.From(tableName).Update(payload).Eq("id", strconv.FormatInt(int64(id), 10)).Execute(&result)
	if err != nil {
		log.Fatalf("Error updating table: %v", err)
	}
}

func Reset(records []model.Stock, tableName string) {
	client := GetSupabaseClient()

	for _, record := range records {
		payload := map[string]interface{}{
			"entry":             record.Entry,
			"cross_match":       record.CrossMatch,
		}

		var result []map[string]interface{}

		err := client.DB.From(tableName).Update(payload).Eq("id", strconv.FormatInt(int64(record.Id), 10)).Execute(&result)
		if err != nil {
			log.Fatalf("Error updating table: %v", err)
		}
	}
}
