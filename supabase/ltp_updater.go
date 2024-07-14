package supabase

import (
	"log"
	"supa_go_ltp_updater/config"
	"supa_go_ltp_updater/model"

	"github.com/nedpals/supabase-go"
)

func LtpUpdater(stocksData []model.Stock) {
	// Initialize the Supabase client
	supabaseUrl := config.AppConfig.Supabase.SupabaseUrl
	supabaseKey := config.AppConfig.Supabase.SupabaseKey

	client := supabase.CreateClient(supabaseUrl, supabaseKey)

	// Define the table and the update payload
	table := "last_traded_price"

	// Define the payload
	var recordsInTable []model.Stock
	//get all table records
	selectErr := client.DB.From(table).Select("symbol").Execute(&recordsInTable)
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
			"last_traded_price": record.LTP,
			"change_pct":        record.ChangePercentage,
			"symbol":            record.Symbol,
		}

		if _, ok := symbolsMap[record.Symbol]; !ok {
			log.Printf("Inserting record: %v\n", record)
			var result []map[string]interface{}
			err := client.DB.From(table).Insert(payload).Execute(&result)
			if err != nil {
				log.Fatalf("Error updating table: %v", err)
			}
			continue
		}

		// FIXME: this needs to be optimised it is very slow 24s a lot
		log.Printf("Updating record: %v\n", record)
		var result []map[string]interface{}
		err := client.DB.From(table).Update(payload).Eq("symbol", record.Symbol).Execute(&result)
		if err != nil {
			log.Fatalf("Error updating table: %v for records: %#v", err, record)
		}
	}
}
