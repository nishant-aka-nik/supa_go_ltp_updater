package supabase

import (
	"log"
	"supa_go_ltp_updater/config"
	"supa_go_ltp_updater/model"
)

func previousDayDataSupabaseUpdater(stocksData []model.PreviousDayData, tableName string) {
	client := GetSupabaseClient()

	var recordsInTable []model.PreviousDayData
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
			"close":  record.Close,
			"symbol": record.Symbol,
			"date":   record.FormatDate(record.Date),
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

func GetPreviousDayData() []model.PreviousDayData {
	client := GetSupabaseClient()

	var recordsInTable []model.PreviousDayData
	//get all table records
	selectErr := client.DB.From(config.AppConfig.TableNames.PreviousDayData).Select("*").Execute(&recordsInTable)
	if selectErr != nil {
		log.Fatalf("Error selecting table: %v", selectErr)
	}

	return recordsInTable
}

func PreviousDayDataUpdater(stocksData []model.Stock, tableName string) {
	//update todays data to previous day data
	allPreviousDayData := []model.PreviousDayData{}

	for _, stock := range stocksData {
		previousDayData := model.PreviousDayData{
			Symbol: stock.Symbol,
			Date:   stock.Date,
			Close:  stock.Close,
		}

		allPreviousDayData = append(allPreviousDayData, previousDayData)
	}

	previousDayDataSupabaseUpdater(allPreviousDayData, tableName)
}

func InsertGapFilteredStocks(stocksData []model.GapFilter, tableName string) {
	client := GetSupabaseClient()

	var recordsInTable []model.GapFilter
	//get all table records
	selectErr := client.DB.From(tableName).Select("symbol").Execute(&recordsInTable)
	if selectErr != nil {
		log.Fatalf("Error selecting table: %v", selectErr)
	}

	var symbolsMap = make(map[string]struct{})
	for _, record := range recordsInTable {
		symbolsMap[record.Symbol] = struct{}{}
	}

	for _, record := range stocksData {
		payload := map[string]interface{}{
			"date":         record.FormatDate(record.Date),
			"symbol":       record.Symbol,
			"gap_pivot":    record.GapPivot,
			"volume_times": record.VolumeTimes,
			"entry":        record.Entry,
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

func GetGapFilteredStocks() []model.GapFilter {
	client := GetSupabaseClient()

	var recordsInTable []model.GapFilter
	//get all table records
	selectErr := client.DB.From(config.AppConfig.TableNames.GapFilter).Select("*").Execute(&recordsInTable)
	if selectErr != nil {
		log.Fatalf("Error selecting table: %v", selectErr)
	}

	return recordsInTable
}
