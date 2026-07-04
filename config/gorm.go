package config

import (
	"fmt"
	"time"

	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(viper *viper.Viper, log *logrus.Logger) *gorm.DB {
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	database := viper.GetString("database.name")

	maxIdle := viper.GetInt("database.pool.idle")
	maxOpen := viper.GetInt("database.pool.max")
	lifetime := viper.GetInt("database.pool.lifetime")

	credentials := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, database,
	)

	db, err := gorm.Open(postgres.Open(credentials), &gorm.Config{
		Logger: logger.New(&logrusWriter{Logger: log}, logger.Config{
			SlowThreshold:             5 * time.Second,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  logger.Info,
		}),
	})

	util.PanicIfError(err)

	sqlDB, err := db.DB()

	util.PanicIfError(err)

	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(lifetime) * time.Second)

	log.Infof(
		"Database connected (pool: idle=%d max=%d lifetime=%ds)",
		maxIdle, maxOpen, lifetime,
	)

	return db
}

type logrusWriter struct {
	Logger *logrus.Logger
}

func (l *logrusWriter) Printf(message string, args ...interface{}) {
	l.Logger.Tracef(message, args...)
}
