package supabase

import (
	"log"
	"strconv"
	"supa_go_ltp_updater/config"
	"supa_go_ltp_updater/model"
)

func InsertCrossMatchedStocks(filteredStocks []model.BreakoutFilter, tableName string) {
	client := GetSupabaseClient()

	crossMatchedStocks := GetCrossMatchedStocks(config.AppConfig.TableNames.BreakoutFilter)

	var symbolsMap = make(map[string]struct{})
	for _, record := range crossMatchedStocks {
		symbolsMap[record.Symbol] = struct{}{}
	}

	// Perform the update operation
	for _, record := range filteredStocks {

		if _, ok := symbolsMap[record.Symbol]; !ok {
			payload := map[string]interface{}{
				"symbol":            record.Symbol,
				"cross_match_date":  record.CrossMatchDate,
				"cross_match":       record.CrossMatch,
				"cross_match_pivot": record.CrossMatchPivot,
				"volume_times":      record.VolumeTimes,
			}
			log.Printf("Inserting record: %v\n", record)
			var result []map[string]interface{}
			err := client.DB.From(tableName).Insert(payload).Execute(&result)
			if err != nil {
				log.Fatalf("Error updating table: %v", err)
			}
		}
	}
}

func GetCrossMatchedStocks(tableName string) []model.BreakoutFilter {
	client := GetSupabaseClient()

	var recordsInTable []model.BreakoutFilter
	//get all table records
	selectErr := client.DB.From(tableName).Select("*").Eq("cross_match", "TRUE").Execute(&recordsInTable)
	if selectErr != nil {
		log.Fatalf("Error selecting table: %v", selectErr)
	}

	return recordsInTable
}

func MarkEntry(entryPrice float64, entryDate string, id int, tableName string) {
	client := GetSupabaseClient()

	payload := map[string]interface{}{
		"entry":       "TRUE",
		"cross_match": "FALSE",
		"entry_date":  entryDate,
		"entry_price": entryPrice,
	}

	var result []map[string]interface{}
	err := client.DB.From(tableName).Update(payload).Eq("id", strconv.FormatInt(int64(id), 10)).Execute(&result)
	if err != nil {
		log.Fatalf("Error updating table: %v", err)
	}
}

func Reset(records []model.BreakoutFilter, tableName string) {
	client := GetSupabaseClient()

	for _, record := range records {
		payload := map[string]interface{}{
			"entry": record.Entry,
		}

		var result []map[string]interface{}

		err := client.DB.From(tableName).Update(payload).Eq("id", strconv.FormatInt(int64(record.ID), 10)).Execute(&result)
		if err != nil {
			log.Fatalf("Error updating table: %v", err)
		}
	}
}

func GetEnteredStocks(tableName string) []model.BreakoutFilter {
	client := GetSupabaseClient()

	var recordsInTable []model.BreakoutFilter
	//get all table records
	selectErr := client.DB.From(tableName).Select("*").Eq("entry", "TRUE").Execute(&recordsInTable)
	if selectErr != nil {
		log.Fatalf("Error selecting table: %v", selectErr)
	}

	return recordsInTable
}
