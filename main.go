package main

import (
	"fmt"
	"log"
	"supa_go_ltp_updater/config"
	"supa_go_ltp_updater/service"
	"time"

	"github.com/robfig/cron/v3"
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
	istLocation, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Fatalf("Failed to load IST location: %v", err)
	}
	currentTime := time.Now().In(istLocation)
	fmt.Printf("Cron job executed at %v\n", currentTime)

	c := cron.New(cron.WithLocation(time.FixedZone("IST", 5*60*60+30*60)))

	// Add a cron job that runs every 10 seconds
	c.AddFunc(config.AppConfig.CronSpec, service.CronLtpUpdater)

	// Start the cron scheduler
	c.Start()
	fmt.Println("Cron scheduler initialized")
	return c
}
