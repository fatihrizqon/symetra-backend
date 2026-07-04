package config

import (
	"context"
	"fmt"
	"time"

	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/fatihrizqon/gofiber-microservice/internal/worker"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapWorkerConfig struct {
	DB     *gorm.DB
	Log    *logrus.Logger
	Config *viper.Viper
}

func BootstrapWorker(ctx context.Context, deps *BootstrapWorkerConfig) (*asynq.Server, *asynq.ServeMux, *asynq.Client) {
	cfg := deps.Config
	logger := deps.Log
	db := deps.DB

	redisOpt := asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%d", cfg.GetString("redis.host"), cfg.GetInt("redis.port")),
		Password: cfg.GetString("redis.password"),
		DB:       cfg.GetInt("redis.db"),
	}

	emailConfig := worker.EmailProcessorConfig{
		Host:        cfg.GetString("smtp.host"),
		Port:        cfg.GetInt("smtp.port"),
		Username:    cfg.GetString("smtp.username"),
		Password:    cfg.GetString("smtp.password"),
		SenderName:  cfg.GetString("smtp.sender_name"),
		SenderEmail: cfg.GetString("smtp.sender_email"),
	}

	emailProcessor := worker.NewEmailProcessor(emailConfig, logger)

	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"default": 10,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(worker.TypeEmailDelivery, emailProcessor.ProcessTask)

	redisJobRepo := repository.NewRedisJobRepository(db)

	asynqClient := asynq.NewClient(redisOpt)

	jobSweeper := worker.NewJobSweeper(redisJobRepo, asynqClient, logger)

	go jobSweeper.Start(ctx, 15*time.Second)

	return srv, mux, asynqClient
}
