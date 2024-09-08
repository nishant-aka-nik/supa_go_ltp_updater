package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	GsheetUrl          string     `json:"gsheetUrl"`
	Supabase           Supabase   `json:"supabase"`
	CronSpec           CronSpec   `json:"cronSpec"`
	Email              SMTP       `json:"email"`
	TargetPercentage   float64    `json:"targetPercentage"`
	StoplossPercentage float64    `json:"stoplossPercentage"`
	TableNames         TableNames `json:"tableNames"`
}

type TableNames struct {
	PreviousDayData string `json:"previousDayData"`
	GapFilter       string `josn:"gapFilter"`
}

type CronSpec struct {
	LtpUpdaterCronSpec   string `json:"ltpUpdateCronSpec"`
	FilterStocksCronSpec string `json:"filterStocksCronSpec"`
	TargetHitCronSpec    string `json:"targetHitCronSpec"`
}

type SMTP struct {
	SMTPServer string `json:"smtp_server"`
	SMTPPort   int    `json:"smtp_port"`
	SMTPEmail  string `json:"smtp_email"`
	SMTPPass   string `json:"smtp_pass"`
}

type Supabase struct {
	SupabaseUrl string `json:"supabaseUrl"`
	SupabaseKey string `json:"supabaseKey"`
}

var AppConfig *Config

func LoadConfig() error {
	CONFIG_PATH, filepathErr := filepath.Abs("./config/config.json")
	if filepathErr != nil {
		return fmt.Errorf("could not find config file: %v", filepathErr)
	}

	file, err := os.Open(CONFIG_PATH)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	AppConfig = &Config{}
	err = json.Unmarshal(bytes, AppConfig)
	if err != nil {
		return err
	}

	return nil
}
