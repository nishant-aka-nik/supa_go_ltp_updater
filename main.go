package main

import (
	"fmt"
	"supa_go_ltp_updater/config"
	"supa_go_ltp_updater/service"

	"github.com/robfig/cron"
)

func main() {
	//init microservice
	initService()
	//-----------------

	fmt.Printf("app config :%#v", config.AppConfig)

	// Initialize the cron scheduler
	c := InitCronScheduler()

	// Defer the stop of the cron scheduler to ensure it stops when main function exits
	defer c.Stop()
	select {}
}

func initService() {
	err := config.LoadConfig()
	if err != nil {
		fmt.Println("Failed to load config:", err)
		return
	}
}

// InitCronScheduler initializes and starts the cron scheduler
func InitCronScheduler() *cron.Cron {
	// Create a new cron instance
	c := cron.New()

	// Add a cron job that runs every 10 seconds
	c.AddFunc("@every 00h00m10s", service.CronLtpUpdater)

	// Start the cron scheduler
	c.Start()
	fmt.Println("Cron scheduler initialized")
	return c
}
