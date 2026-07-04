package notifier

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type SlackNotifier struct {
	WebhookURL string
}

func NewSlack(webhook string) *SlackNotifier {
	return &SlackNotifier{WebhookURL: webhook}
}

func (s *SlackNotifier) Send(msg Message) error {
	payload := map[string]interface{}{
		"text": "*" + msg.Title + "*\n" + msg.Content,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(
		"POST",
		s.WebhookURL,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	_, err = http.DefaultClient.Do(req)
	return err
}
