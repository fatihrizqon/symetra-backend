package main

import (
	"context"
	"log"

	"github.com/fatihrizqon/gofiber-microservice/config"
)

func main() {
	viper := config.NewViper()
	logger := config.NewLogger(viper)
	db := config.NewDatabase(viper, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv, mux, asynqClient := config.BootstrapWorker(ctx, &config.BootstrapWorkerConfig{
		DB:     db,
		Log:    logger,
		Config: viper,
	})
	defer asynqClient.Close()

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
