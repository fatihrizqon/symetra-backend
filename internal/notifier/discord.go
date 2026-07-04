package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DiscordNotifier struct {
	WebhookURL string
}

func NewDiscord(webhook string) *DiscordNotifier {
	return &DiscordNotifier{WebhookURL: webhook}
}

func (d *DiscordNotifier) Send(msg Message) error {
	payload := map[string]interface{}{
		"username": "Go Service Logger",
		"embeds": []map[string]interface{}{
			{
				"title":       msg.Title,
				"description": msg.Content,
				"color":       levelColor(msg.Level),
				"fields":      formatFields(msg.Fields),
			},
		},
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(
		"POST",
		d.WebhookURL,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	_, err = http.DefaultClient.Do(req)
	return err
}

func levelColor(level string) int {
	switch level {
	case "error", "fatal", "panic":
		return 0xFF0000
	case "warn":
		return 0xFFA500
	default:
		return 0x00FF00
	}
}

func formatFields(fields map[string]interface{}) []map[string]string {
	result := []map[string]string{}
	for k, v := range fields {
		result = append(result, map[string]string{
			"name":  k,
			"value": fmt.Sprintf("%v", v),
		})
	}
	return result
}
