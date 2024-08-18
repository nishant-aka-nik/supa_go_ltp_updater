package main

import (
	"fmt"
	"log"
	"supa_go_ltp_updater/config"
	"supa_go_ltp_updater/service"
	"supa_go_ltp_updater/utils"
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

	RunServiceOnStartup()

	// Defer the stop of the cron scheduler to ensure it stops when main function exits
	defer c.Stop()
	select {}
}

func initService() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}
}

func RunServiceOnStartup() {
	service.CronLtpUpdater()
	service.FilterStocks()
}

// InitCronScheduler initializes and starts the cron scheduler
func InitCronScheduler() *cron.Cron {
	currentTime := utils.GetISTTime()
	log.Printf("Initializing cron job at %v\n", currentTime)

	c := cron.New(cron.WithLocation(time.FixedZone("IST", 5*60*60+30*60)),
		cron.WithParser(
			cron.NewParser(
				cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow)),
	)

	// LTP updater cron job
	log.Printf("Adding LTP updater cron job with spec: %s\n", config.AppConfig.CronSpec.LtpUpdaterCronSpec)
	cronEntryID, cronErr := c.AddFunc(config.AppConfig.CronSpec.LtpUpdaterCronSpec, service.CronLtpUpdater)
	if cronErr != nil {
		log.Fatalf("Failed to add cron job: %v", cronErr)
	}
	log.Printf("Cron job added with ID: %d\n", cronEntryID)

	// Filter stocks cron job
	log.Printf("Adding filter stocks cron job with spec: %s\n", config.AppConfig.CronSpec.FilterStocksCronSpec)
	cronEntryID, cronErr = c.AddFunc(config.AppConfig.CronSpec.FilterStocksCronSpec, service.FilterStocks)
	if cronErr != nil {
		log.Fatalf("Failed to add cron job: %v", cronErr)
	}
	log.Printf("Cron job added with ID: %d\n", cronEntryID)

	// Target hit cron job
	log.Printf("Adding target hit cron job with spec: %s\n", config.AppConfig.CronSpec.TargetHitCronSpec)
	cronEntryID, cronErr = c.AddFunc(config.AppConfig.CronSpec.FilterStocksCronSpec, service.TargetHitCheckerCron)
	if cronErr != nil {
		log.Fatalf("Failed to add cron job: %v", cronErr)
	}
	log.Printf("Cron job added with ID: %d\n", cronEntryID)

	// Start the cron scheduler
	log.Println("Starting cron scheduler")
	c.Start()
	log.Println("Cron scheduler initialized")
	return c
}
