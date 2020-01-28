package main

import (
	"albion-status-checker/checker"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var webhook_url string = ""

func main() {
	ParseEnv()
	if webhook_url == "" {
		panic("Webhook url environment is not set")
	}
	check := checker.NewChecker()
	check.Start()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-sc
		check.Stop()
	}()
	err := check.CheckStatus()
	if err == nil {
		status := <-check.Changes
		fmt.Printf("Current server status: %s\n", status.Status)
	}
	fmt.Println("Started checking")
	for message := range check.Changes {
		SendStatusChangeWebhook(message)
	}
}

func ParseEnv() {
	env := syscall.Environ()
	for _, line := range env {
		fields := strings.Split(line, "=")
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])
		if key == "WEBHOOK_URL" {
			webhook_url = value
		}
	}
}
