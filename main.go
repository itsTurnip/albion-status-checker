package main

import (
	"github.com/itsTurnip/albion-status-checker/checker"
	"os"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

var webhookURL string

func main() {
	ParseEnv()
	if webhookURL == "" {
		panic("Webhook url environment variable is not set")
	}
	check := checker.NewChecker()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-sc
		check.Stop()
	}()
	go check.CheckStatus()
	log.Info("Getting current status...")
	status := <-check.Changes
	log.Info("Current server status: ", status.Status)
	check.Start()
	log.Info("Started checking")
	for message := range check.Changes {
		err := SendStatusChangeWebhook(message)
		log.Errorf("Error occured sending status change: %s", err)
	}
}

// ParseEnv parses environment variables
func ParseEnv() {
	env := syscall.Environ()
	for _, line := range env {
		fields := strings.Split(line, "=")
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])
		switch key {
		case "WEBHOOK_URL":
			webhookURL = value
		case "LOGLEVEL":
			switch value {
			case "INFO":
				log.SetLevel(log.InfoLevel)
			case "DEBUG":
				log.SetLevel(log.DebugLevel)
			case "ERROR":
				log.SetLevel(log.ErrorLevel)
			case "FATAL":
				log.SetLevel(log.FatalLevel)
			case "WARN":
				log.SetLevel(log.WarnLevel)
			default:
				log.SetLevel(log.InfoLevel)
			}
		}
	}
}
