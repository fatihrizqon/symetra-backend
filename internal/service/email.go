package service

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/fatihrizqon/gofiber-microservice/internal/worker"
	"github.com/sirupsen/logrus"
)

// EmailConfig is now removed.

type IEmailService interface {
	CreateVerificationEmailJob(to, name, verificationLink string) (*entity.RedisJob, error)
	CreateResetPasswordEmailJob(to, name, resetLink string) (*entity.RedisJob, error)
	CreateNotificationEmailJob(to, subject, message string) (*entity.RedisJob, error)
}

type EmailService struct {
	log *logrus.Logger
}

func NewEmailService(log *logrus.Logger) IEmailService {
	return &EmailService{
		log: log,
	}
}

func (s *EmailService) parseTemplate(templateName string, data interface{}) (string, error) {
	tmplPath := filepath.Join("templates", "email", templateName)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return buf.String(), nil
}

func (s *EmailService) createEmailJob(to, subject, body string) (*entity.RedisJob, error) {
	task, err := worker.NewEmailDeliveryTask(to, subject, body)
	if err != nil {
		s.log.Errorf("Failed to create email delivery task: %v", err)
		return nil, err
	}

	job := &entity.RedisJob{
		ID:      util.GenerateUUID().String(),
		Type:    task.Type(),
		Payload: task.Payload(),
		Status:  "PENDING",
	}

	return job, nil
}

func (s *EmailService) CreateVerificationEmailJob(to, name, verificationLink string) (*entity.RedisJob, error) {
	data := map[string]string{
		"Name":             name,
		"VerificationLink": verificationLink,
	}

	body, err := s.parseTemplate("verification.html", data)
	if err != nil {
		return nil, err
	}

	return s.createEmailJob(to, "Verifikasi Alamat Email Anda", body)
}

func (s *EmailService) CreateResetPasswordEmailJob(to, name, resetLink string) (*entity.RedisJob, error) {
	data := map[string]string{
		"Name":      name,
		"ResetLink": resetLink,
	}

	body, err := s.parseTemplate("reset_password.html", data)
	if err != nil {
		return nil, err
	}

	return s.createEmailJob(to, "Reset Password Akun Anda", body)
}

func (s *EmailService) CreateNotificationEmailJob(to, subject, message string) (*entity.RedisJob, error) {
	data := map[string]string{
		"Message": message,
	}

	body, err := s.parseTemplate("notification.html", data)
	if err != nil {
		return nil, err
	}

	return s.createEmailJob(to, subject, body)
}
