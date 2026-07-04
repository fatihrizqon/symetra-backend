package helper

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

type RefreshCookie struct {
	Name       string
	TTL        time.Duration
	Production bool
}

func (c RefreshCookie) Set(ctx fiber.Ctx, name string, token string) {
	ctx.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   c.Production,
		SameSite: fiber.CookieSameSiteLaxMode,
		MaxAge:   int(c.TTL.Seconds()),
	})
}

func (c RefreshCookie) Clear(ctx fiber.Ctx, name string) {
	ctx.Cookie(&fiber.Cookie{
		Name:     name,
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: true,
	})
}
