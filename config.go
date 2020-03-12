package main

import (
	"github.com/itsTurnip/dishooks"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// Config doesn't need to be explained
type Config struct {
	// Webhooks structs of discord webhooks
	Webhooks []*dishooks.Webhook
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
			links := strings.Split(value, ",")
			for _, link := range links {
				webhook, err := dishooks.WebhookFromURL(link)
				if err != nil {
					log.Error("Error occurred while getting webhook: ", err)
					continue
				}
				c.Webhooks = append(c.Webhooks, webhook)
			}
		case "LOGLEVEL":
			switch value {
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
