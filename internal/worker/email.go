package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

const (
	TypeEmailDelivery = "email:deliver"
)

type EmailTaskPayload struct {
	To      string
	Subject string
	Body    string
}

func NewEmailDeliveryTask(to, subject, body string) (*asynq.Task, error) {
	payload := EmailTaskPayload{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeEmailDelivery, payloadBytes), nil
}

type EmailProcessorConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	SenderName  string
	SenderEmail string
}

type EmailProcessor struct {
	config EmailProcessorConfig
	log    *logrus.Logger
}

func NewEmailProcessor(config EmailProcessorConfig, log *logrus.Logger) *EmailProcessor {
	return &EmailProcessor{
		config: config,
		log:    log,
	}
}

func (p *EmailProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload EmailTaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	m := gomail.NewMessage()
	sender := fmt.Sprintf("%s <%s>", p.config.SenderName, p.config.SenderEmail)
	m.SetHeader("From", sender)
	m.SetHeader("To", payload.To)
	m.SetHeader("Subject", payload.Subject)
	m.SetBody("text/html", payload.Body)

	d := gomail.NewDialer(p.config.Host, p.config.Port, p.config.Username, p.config.Password)

	if err := d.DialAndSend(m); err != nil {
		p.log.Errorf("Failed to send email to %s: %v", payload.To, err)
		return err // Return err to trigger Asynq retry
	}

	p.log.Infof("Email successfully sent to %s via Asynq worker", payload.To)
	return nil
}
