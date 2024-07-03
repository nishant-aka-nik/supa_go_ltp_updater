package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	GsheetUrl   string   `json:"gsheetUrl"`
	RunEveryMin int      `json:"runEveryMin"`
	Supabase    Supabase `json:"supabase"`
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
