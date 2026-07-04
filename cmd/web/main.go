package main

import (
	"fmt"

	"github.com/fatihrizqon/gofiber-microservice/config"
	"github.com/fatihrizqon/gofiber-microservice/database"
	"github.com/gofiber/fiber/v3"
)

// @title Go REST API with Fiber Framework
// @version 1.0
// @description This is an Official Documentation for Go REST API with Fiber Framework
// @termsOfService http://swagger.io/terms/
// @contact.name Fatih Rizqon
// @contact.email fatihrizqon@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 127.0.0.1:3000
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	viper := config.NewViper()
	cors := config.NewCORS(viper)
	log := config.NewLogger(viper)
	db := config.NewDatabase(viper, log)
	validate := config.NewValidator(viper)
	app := config.NewFiber(viper, log)
	jwt := config.NewJWT(viper, log)
	prod := viper.GetBool("production")

	config.Bootstrap(&config.BootstrapConfig{
		Cors:       cors,
		DB:         db,
		App:        app,
		Log:        log,
		Validate:   validate,
		Config:     viper,
		JWT:        jwt,
		Production: prod,
	})

	database.Migrate(db)

	config.NewNotifier(viper, log)

	port := viper.GetInt("web.port")

	err := app.Listen(fmt.Sprintf(":%d", port), fiber.ListenConfig{EnablePrefork: viper.GetBool("web.prefork")})

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
