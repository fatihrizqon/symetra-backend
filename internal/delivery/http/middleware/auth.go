package middleware

import (
	"fmt"
	"strings"

	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func NewAuth(tokenRepo repository.ITokenRepository) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		authHeader := ctx.Get("Authorization", "")
		if authHeader == "" {
			util.HandleError(ctx, fiber.StatusUnauthorized, fmt.Errorf("missing authorization header"))
			return nil
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			util.HandleError(ctx, fiber.StatusUnauthorized, fmt.Errorf("invalid authorization format"))
			return nil
		}

		token := strings.TrimPrefix(authHeader, bearerPrefix)

		claims, err := util.ParseAccessToken(token)
		if err != nil {
			util.HandleError(ctx, fiber.StatusUnauthorized, err)
			return nil
		}

		sessionID := claims.SessionID
		if sessionID == uuid.Nil {
			util.HandleError(ctx, fiber.StatusUnauthorized, fmt.Errorf("invalid token claims"))
			return nil
		}

		if _, err := tokenRepo.FindSessionByID(sessionID); err != nil {
			util.HandleError(ctx, fiber.StatusUnauthorized, fmt.Errorf("session revoked"))
			return nil
		}

		ctx.Locals("auth", claims)
		ctx.Locals("session_id", sessionID)

		return ctx.Next()
	}
}
