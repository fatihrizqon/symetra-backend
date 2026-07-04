package notifier

import "github.com/sirupsen/logrus"

type Hook struct {
	Notifier  Notifier
	LogLevels []logrus.Level
}

func (h *Hook) Levels() []logrus.Level {
	return h.LogLevels
}

func (h *Hook) Fire(entry *logrus.Entry) error {
	msg := Message{
		Title:   entry.Level.String(),
		Content: entry.Message,
		Level:   entry.Level.String(),
		Fields:  entry.Data,
	}

	go h.Notifier.Send(msg)
	return nil
}
