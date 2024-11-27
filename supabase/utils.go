package supabase

import (
	"log"
	"sync"

	"github.com/nedpals/supabase-go"
)

func InsertSupabase(wg *sync.WaitGroup, client *supabase.Client, payload map[string]interface{}, tableName string) {
	defer wg.Done()
	log.Printf("Inserting record: %v\n", payload)
	var result []map[string]interface{}
	err := client.DB.From(tableName).Insert(payload).Execute(&result)
	if err != nil {
		log.Fatalf("Error updating table: %v", err)
	}
}

func UpdateSupabase(wg *sync.WaitGroup, client *supabase.Client, payload map[string]interface{}, tableName string, symbol string) {
	defer wg.Done()
	log.Printf("Updating record: %v\n", payload)
	var result []map[string]interface{}
	err := client.DB.From(tableName).Update(payload).Eq("symbol", symbol).Execute(&result)
	if err != nil {
		log.Fatalf("Error updating table: %v for records: %#v", err, symbol)
	}
}
