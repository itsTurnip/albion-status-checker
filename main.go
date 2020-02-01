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
	if config.WebhookURL == "" {
		panic("Webhook url environment variable is not set")
	}
	check := checker.NewChecker()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-sc
		check.Stop()
	}()
	log.Info("Getting current status...")
	err := check.CheckStatus()
	if err != nil {
		panic(err)
	}
	status := <-check.Changes
	log.Info("Current server status: ", status.Status)
	check.Start()
	log.Info("Started checking")
	for message := range check.Changes {
		err := SendStatusChangeWebhook(config.WebhookURL, message)
		log.Errorf("Error occured sending status change: %s", err)
	}
}
