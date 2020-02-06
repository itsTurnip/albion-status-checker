package main

import (
	"github.com/itsTurnip/dishooks"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// Config doesn't need to be explained
type Config struct {
	// WebhookURL of discord webhook
	Webhook *dishooks.Webhook
	// Logging level
	LogLevel log.Level
}

func parseEnv() (c *Config) {
	c = &Config{
		LogLevel: log.InfoLevel,
	}
	env := syscall.Environ()
	for _, line := range env {
		fields := strings.Split(line, "=")
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])
		switch key {
		case "WEBHOOK_URL":
			webhook, err := dishooks.WebhookFromURL(value)
			if err != nil {
				log.Error("Error occurred while getting webhook: ", err)
				break
			}
			c.Webhook = webhook
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
