package config

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type JWTService struct {
	Secret        string
	RefreshSecret string
	Expiration    time.Duration
}

func NewJWT(config *viper.Viper, log *logrus.Logger) *JWTService {
	return &JWTService{
		Secret:        config.GetString("jwt.secret"),
		RefreshSecret: config.GetString("jwt.refresh_secret"),
		Expiration:    config.GetDuration("jwt.expiration"),
	}
}
