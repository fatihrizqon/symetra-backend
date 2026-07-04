package util

import (
	"time"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var accessSecret []byte
var refreshSecret []byte

type Claims struct {
	UserID    uuid.UUID `json:"uid"`
	SessionID uuid.UUID `json:"sid"`
	jwt.RegisteredClaims
}

func NewJWT(config *viper.Viper) {
	jwtSecret := config.GetString("jwt.secret")
	jwtRefreshSecret := config.GetString("jwt.refresh_secret")
	if jwtSecret == "" || jwtRefreshSecret == "" {
		panic("jwt secret missing")
	}
	accessSecret = []byte(jwtSecret)
	refreshSecret = []byte(jwtRefreshSecret)
}

func CreateAccessToken(user entity.User, sessionID uuid.UUID) (string, error) {
	return create(user, sessionID, accessSecret, 15*time.Minute)
}

func CreateRefreshToken(user entity.User, sessionID uuid.UUID) (string, error) {
	return create(user, sessionID, refreshSecret, 7*24*time.Hour)
}

func ParseAccessToken(token string) (*Claims, error) {
	return parse(token, accessSecret)
}

func ParseRefreshToken(token string) (*Claims, error) {
	return parse(token, refreshSecret)
}

func parse(tokenString string, secret []byte) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, fiber.ErrUnauthorized
			}
			return secret, nil
		},
	)
	if err != nil || !token.Valid {
		return nil, fiber.ErrUnauthorized
	}
	return claims, nil
}

func create(user entity.User, sessionID uuid.UUID, secret []byte, duration time.Duration) (string, error) {
	claims := Claims{
		UserID:    user.Id,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
