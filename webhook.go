package main

import (
	"github.com/itsTurnip/albion-status-checker/checker"
	"github.com/itsTurnip/dishooks"
)

func SendStatusChangeWebhook(webhook *dishooks.Webhook, message *checker.StatusMessage) error {
	embed := &dishooks.Embed{
		Title:       "Статус сервера",
		Description: "Изменился статус сервера",
		Type:        "rich",
		URL:         "https://www.albionstatus.com/",
	}
	field := &dishooks.EmbedField{
		Value: message.Message,
	}
	switch message.Status {
	case checker.OnlineStatus:
		embed.Color = 0x00FF00
		field.Name = "Онлайн"
	case checker.OfflineStatus:
		embed.Color = 0xFF0000
		field.Name = "Оффлайн"
	case checker.TimeoutStatus:
		embed.Color = 0x735184
		field.Name = "Таймаут"
	case checker.StartingStatus:
		embed.Color = 0xFFFF00
		field.Name = "Запускается"
	default:
		field.Name = message.Status
	}
	embed.Fields = []*dishooks.EmbedField{
		field,
	}
	embed.Timestamp = message.Timestamp
	webhookMessage := &dishooks.WebhookMessage{
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
