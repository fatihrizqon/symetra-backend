package config

import (
	"fmt"

	"github.com/gofiber/fiber/v3/middleware/static"

	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/handler"
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/middleware"
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/route"
	"github.com/fatihrizqon/gofiber-microservice/internal/rbac"
	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/fatihrizqon/gofiber-microservice/internal/service"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/fatihrizqon/gofiber-microservice/internal/util/storage"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	Cors       *CORSConfig
	DB         *gorm.DB
	App        *fiber.App
	Log        *logrus.Logger
	Validate   *validator.Validate
	Config     *viper.Viper
	JWT        *JWTService
	Production bool
}

func Bootstrap(deps *BootstrapConfig) {
	cfg := deps.Config

	deps.App.Use(deps.Cors.Handler())

	// ── Static Files ─────────────────────────────────────────────────────────
	deps.App.Get("/uploads*", static.New("./public/uploads"))

	// ── Storage ──────────────────────────────────────────────────────────────
	baseURL := fmt.Sprintf("%s/uploads", cfg.GetString("web.base_url"))
	localStorage := storage.NewLocalStorage("./public/uploads", baseURL)

	// ── Repositories ──────────────────────────────────────────────────────────
	userRepository := repository.NewUserRepository(deps.DB)
	authRepository := repository.NewAuthRepository(deps.DB)
	tokenRepository := repository.NewTokenRepository(deps.DB)
	rbacRepository := repository.NewRbacRepository(deps.DB)
	fileRepository := repository.NewFileRepository(deps.DB)

	// ── Email Service ─────────────────────────────────────────────────────────
	emailService := service.NewEmailService(deps.Log)

	// ── Services ──────────────────────────────────────────────────────────────
	fileService := service.NewFileService(fileRepository, localStorage)
	userService := service.NewUserService(userRepository, deps.Validate, fileService)
	appURL := cfg.GetString("web.base_url")
	authService := service.NewAuthService(authRepository, tokenRepository, deps.Validate, emailService, appURL)
	rbacService := service.NewRbacService(rbacRepository, userRepository, deps.Validate)

	// ── Handlers ──────────────────────────────────────────────────────────────
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService, deps.Production)
	fileHandler := handler.NewFileHandler(fileService)
	rbacHandler := handler.NewRbacHandler(rbacService)

	// ── Middleware ────────────────────────────────────────────────────────────
	authMiddleware := middleware.NewAuth(tokenRepository)

	rbacEngine := rbac.New(rbac.Config{
		Store: rbacRepository,
		UserLookup: func(c fiber.Ctx) string {
			claims, ok := c.Locals("auth").(*util.Claims)
			if !ok {
				return ""
			}
			return claims.UserID.String()
		},
	})

	// ── Routes ────────────────────────────────────────────────────────────────
	routeConfig := route.RouteConfig{
		App:            deps.App,
		AuthMiddleware: authMiddleware,
		RbacEngine:     rbacEngine,
		UserHandler:    userHandler,
		AuthHandler:    authHandler,
		FileHandler:    fileHandler,
		RbacHandler:    rbacHandler,
		Production:     deps.Production,
	}

	routeConfig.Setup()
}
