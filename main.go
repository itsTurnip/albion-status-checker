package main

import (
	"github.com/itsTurnip/albion-status-checker/checker"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func main() {
	config := parseEnv()
	if len(config.Webhooks) == 0 {
		log.Fatal("Webhook url environment variable is not set")
	}
	log.SetLevel(config.LogLevel)
	check := checker.NewChecker()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	log.Info("Getting current status...")
	err := check.CheckStatus()
	if err != nil {
		log.Fatal(err)
	}
	status := <-check.Changes
	log.Info("Current server status: ", status.Status)
	err = check.Start()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		<-sc
		err := check.Stop()
		if err != nil {
			log.Errorf("Error stopping checker %s", err)
		}
	}()
	log.Info("Started checking")
	for message := range check.Changes {
		for _, webhook := range config.Webhooks {
			err := SendStatusChangeWebhook(webhook, message)
			if err != nil {
				log.Errorf("Error occurred sending status change: %s", err)
			}
		}
	}
}
