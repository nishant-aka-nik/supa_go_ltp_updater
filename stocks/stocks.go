package stocks

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"supa_go_ltp_updater/config"
	"supa_go_ltp_updater/model"
	"time"
)

func GetStocks() []model.Stock {
	url := config.AppConfig.GsheetUrl
	csvData, err := fetchCSV(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching CSV: %v\n", err)
		return nil
	}

	stocks, err := csvToJSON(csvData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting CSV to JSON: %v\n", err)
		return nil
	}
	
	return stocks
}

// csvToJSON converts CSV data to a slice of Stock structs.
func csvToJSON(csvData []byte) ([]model.Stock, error) {
	reader := csv.NewReader(strings.NewReader(string(csvData)))

	// Read the header row
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	var stocks []model.Stock
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		row := make(map[string]string)
		for i, value := range record {
			row[headers[i]] = value
		}

		// Convert the map to a Stock struct
		stock, err := mapToStock(row)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

func fetchCSV(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// mapToStock converts a map of strings to a Stock struct dynamically.
func mapToStock(row map[string]string) (model.Stock, error) {
	var stock model.Stock
	stockValue := reflect.ValueOf(&stock).Elem()
	stockType := stockValue.Type()

	for i := 0; i < stockType.NumField(); i++ {
		field := stockType.Field(i)
		jsonTag := field.Tag.Get("json")

		if value, ok := row[jsonTag]; ok {
			fieldValue := stockValue.Field(i)
			x := fieldValue.Kind()
			switch x {
			case reflect.String:
				fieldValue.SetString(value)
			case reflect.Float64:
				if parsedValue, err := strconv.ParseFloat(value, 64); err == nil {
					fieldValue.SetFloat(parsedValue)
				} else {
					return stock, err
				}
			case reflect.Struct:
				if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
					layout := "02/01/2006"
					if parsedDate, err := time.Parse(layout, value); err == nil {
						fieldValue.Set(reflect.ValueOf(parsedDate))
					} else {
						return stock, err
					}
				}
			default:
				return stock, errors.New("unsupported field type")
			}
		}
	}

	return stock, nil
}
