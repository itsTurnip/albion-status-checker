package main

// TODO: Create a better package to work with webhooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"albion-status-checker/checker"
)

type WebhookMessage struct {
	Content   string   `json:"content"`
	Username  string   `json:"username"`
	AvatarURL string   `json:"avatar_url"`
	Embeds    []*Embed `json:"embeds"`
}

type Embed struct {
	Title       string        `json:"title"`
	Type        string        `json:"type"`
	Description string        `json:"description"`
	URL         string        `json:"url"`
	Timestamp   string        `json:"timestamp"`
	Color       int           `json:"color"`
	Fields      []*EmbedField `json:"fields"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

func SendStatusChangeWebhook(message checker.StatusMessage) {
	embed := &Embed{
		Title:       "Статус сервера",
		Description: "Изменился статус сервера",
		Type:        "rich",
		URL:         "https://www.albionstatus.com/",
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
	field := &EmbedField{}
	switch message.Status {
	case "online":
		embed.Color = 0x00FF00
		field.Name = "Онлайн"
		field.Value = message.Message
	case "offline":
		embed.Color = 0xFF0000
		field.Name = "Оффлайн"
		field.Value = message.Message
	}
	embed.Fields = []*EmbedField{
		field,
	}
	webhookMessage := &WebhookMessage{
		// Content:   "@here",
		AvatarURL: "http://www.fau.edu/oit/labs/labimages/Status-dialog-information-icon.png",
		Embeds: []*Embed{
			embed,
		},
	}
	err := SendWebhookMessage(webhook_url, webhookMessage)
	if err != nil {
		fmt.Println(err)
	}
}
func SendWebhookMessage(url string, message *WebhookMessage) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, err = http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	return nil
}
