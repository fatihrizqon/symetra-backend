package config

import (
	"github.com/fatihrizqon/gofiber-microservice/internal/notifier"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewNotifier(v *viper.Viper, log *logrus.Logger) {
	provider := v.GetString("notifier.provider")
	webhook := v.GetString("notifier.webhook")

	var notify notifier.Notifier

	switch provider {
	case "discord":
		notify = notifier.NewDiscord(webhook)
	case "slack":
		notify = notifier.NewSlack(webhook)
	default:
		return
	}

	log.AddHook(&notifier.Hook{
		Notifier: notify,
		LogLevels: []logrus.Level{
			logrus.ErrorLevel,
			logrus.FatalLevel,
			logrus.PanicLevel,
		},
	})
}
