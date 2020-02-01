package main

import (
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// Config doesn't need to be explained
type Config struct {
	// WebhookURL of discord webhook
	WebhookURL string
	// Logging level
	LogLevel   log.Level
}

func parseEnv() (c *Config) {
	c = &Config{}
	env := syscall.Environ()
	for _, line := range env {
		fields := strings.Split(line, "=")
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])
		switch key {
		case "WEBHOOK_URL":
			c.WebhookURL = value
		case "LOGLEVEL":
			switch value {
			case "INFO":
				c.LogLevel = log.InfoLevel
			case "DEBUG":
				c.LogLevel = log.DebugLevel
			case "ERROR":
				c.LogLevel = log.ErrorLevel
			case "FATAL":
				c.LogLevel = log.FatalLevel
			case "WARN":
				c.LogLevel = log.WarnLevel
			default:
				c.LogLevel = log.InfoLevel
			}
		}
	}
	return
}
