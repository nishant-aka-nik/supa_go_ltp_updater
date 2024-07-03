package supabase

import (
	"fmt"
	"log"

	"github.com/nedpals/supabase-go"
)

func LtpUpdater() {
	// Initialize the Supabase client
	supabaseUrl := "https://yutkfbluvhuoshqgbtxe.supabase.co"
	supabaseKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Inl1dGtmYmx1dmh1b3NocWdidHhlIiwicm9sZSI6ImFub24iLCJpYXQiOjE3MTc2OTI5NDEsImV4cCI6MjAzMzI2ODk0MX0.eopTb_b240QiYOvDfsFa0UcXh4c4xHXS9NvU4TdOECc"

	client := supabase.CreateClient(supabaseUrl, supabaseKey)

	// Define the table and the update payload
	table := "last_traded_price"
	payload := map[string]interface{}{
		"ltp": 100.0,
	}

	// Perform the update operation
	var result []map[string]interface{}
	err := client.DB.From(table).Update(payload).Eq("name", "UNOMINDA").Execute(&result)
	if err != nil {
		log.Fatalf("Error updating table: %v", err)
	}

	// Print the result
	for _, row := range result {
		fmt.Printf("%v\n", row)
	}
}
