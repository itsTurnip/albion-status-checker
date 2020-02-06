package main

import (
	"time"

	"github.com/itsTurnip/albion-status-checker/checker"
	"github.com/itsTurnip/dishooks"
)

func SendStatusChangeWebhook(webhook *dishooks.Webhook, message checker.StatusMessage) error {
	embed := &dishooks.Embed{
		Title:       "Статус сервера",
		Description: "Изменился статус сервера",
		Type:        "rich",
		URL:         "https://www.albionstatus.com/",
		Timestamp:   dishooks.FormatTime(time.Now()),
	}
	field := &dishooks.EmbedField{
		Value: message.Message,
	}
	switch message.Status {
	case "online":
		embed.Color = 0x00FF00
		field.Name = "Онлайн"
	case "offline":
		embed.Color = 0xFF0000
		field.Name = "Оффлайн"
	case "timeout":
		embed.Color = 0x735184
		field.Name = "Таймаут"
	default:
		field.Name = "?"
	}
	embed.Fields = []*dishooks.EmbedField{
		field,
	}
	webhookMessage := &dishooks.WebhookMessage{
		// Content:   "@here",
		AvatarURL: "http://www.fau.edu/oit/labs/labimages/Status-dialog-information-icon.png",
		Embeds: []*dishooks.Embed{
			embed,
		},
	}
	_, err := webhook.SendMessage(webhookMessage)
	if err != nil {
		return err
	}
	return nil
}
