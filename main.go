package main

import (
	"context"
	"fmt"
	"log"
	"supa_go_ltp_updater/config"
	"supa_go_ltp_updater/constants"
	contextkeys "supa_go_ltp_updater/context"
	"supa_go_ltp_updater/service"
	"supa_go_ltp_updater/utils"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	ctx := context.Background()

	//init microservice
	initService()
	//-----------------

	fmt.Printf("app config :%#v", config.AppConfig)

	// Initialize the cron scheduler
	c := InitCronScheduler(ctx)

	RunServiceOnStartup(ctx)

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

func RunServiceOnStartup(ctx context.Context) {
	start := utils.GetISTTime()
	// Check if today is Saturday (6) or Sunday (0)
	if start.Weekday() == time.Saturday || start.Weekday() == time.Sunday {
		log.Println("It's the weekend! Skipping execution.")
		return // Skip further execution
	}

	if !utils.IsSafeTimeToRun() {
		log.Printf("It's not safe time to run skipping execution. time: %v", utils.GetISTTime())
		return
	}

	ctx = contextkeys.SetCaller(ctx, constants.RunServiceOnStartupCaller)

	service.CronLtpUpdater()
	service.FilterStocks(ctx)
	service.TargetHitCheckerCron(ctx)
	service.Gaptor()
}

// InitCronScheduler initializes and starts the cron scheduler
func InitCronScheduler(ctx context.Context) *cron.Cron {
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
	cronEntryID, cronErr = addCronJobWithContext(c, config.AppConfig.CronSpec.FilterStocksCronSpec, "FilterStocks", service.FilterStocks)
	if cronErr != nil {
		log.Fatalf("Failed to add cron job: %v", cronErr)
	}
	log.Printf("Cron job added with ID: %d\n", cronEntryID)

	// Target hit cron job
	log.Printf("Adding target hit cron job with spec: %s\n", config.AppConfig.CronSpec.TargetHitCronSpec)
	cronEntryID, cronErr = addCronJobWithContext(c, config.AppConfig.CronSpec.TargetHitCronSpec, "TargetHitCheckerCron", service.TargetHitCheckerCron)
	if cronErr != nil {
		log.Fatalf("Failed to add cron job: %v", cronErr)
	}
	log.Printf("Cron job added with ID: %d\n", cronEntryID)

	// Target hit cron job
	log.Printf("Adding gaptor cron job with spec: %s\n", config.AppConfig.CronSpec.GaptorCronSpec)
	cronEntryID, cronErr = c.AddFunc(config.AppConfig.CronSpec.GaptorCronSpec, service.Gaptor)
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

func addCronJobWithContext(c *cron.Cron, spec string, jobName string, jobFunc func(ctx context.Context)) (cron.EntryID, error) {
	return c.AddFunc(spec, func() {
		ctx := context.Background()
		ctx = contextkeys.SetCaller(ctx, constants.CronCaller)

		log.Printf("Starting cron job: %s", jobName)
		jobFunc(ctx)
		log.Printf("Completed cron job: %s", jobName)
	})
}
