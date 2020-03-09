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
	if config.Webhook == nil {
		panic("Webhook url environment variable is not set")
	}
	log.SetLevel(config.LogLevel)
	check := checker.NewChecker()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-sc
		err := check.Stop()
		if err != nil {
			log.Errorf("Error stopping checker %s", err)
		}
	}()
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
	log.Info("Started checking")
	for message := range check.Changes {
		err := SendStatusChangeWebhook(config.Webhook, message)
		if err != nil {
			log.Errorf("Error occurred sending status change: %s", err)
		}
	}
}
