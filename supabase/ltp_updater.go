package supabase

import (
	"fmt"
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

	// Perform the update operation
	for _, record := range stocksData {
		payload := map[string]interface{}{
			"last_traded_price": record.LTP,
			"change_pct":        record.ChangePercentage,
		}

		var result []map[string]interface{}
		err := client.DB.From(table).Update(payload).Eq("symbol", record.Symbol).Execute(&result)
		if err != nil {
			log.Fatalf("Error updating table: %v", err)
		}

		// Print the result
		for _, row := range result {
			fmt.Printf("%v\n", row)
		}
	}
}
