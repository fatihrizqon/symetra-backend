package handler

import (
	"time"

	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/request"
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/response"
	"github.com/fatihrizqon/gofiber-microservice/internal/helper"
	"github.com/fatihrizqon/gofiber-microservice/internal/service"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	IAuthService service.IAuthService
	Cookie       helper.RefreshCookie
}

func NewAuthHandler(serv service.IAuthService, production bool) *AuthHandler {
	return &AuthHandler{
		IAuthService: serv,
		Cookie: helper.RefreshCookie{
			TTL:        7 * 24 * time.Hour,
			Production: production,
		},
	}
}

// Register godoc
// @Summary User registration
// @Description Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.RegisterRequest true "Register request"
// @Success 201 {object} response.JSON "User registered successfully"
// @Failure 400 {object} response.JSON "Invalid request format"
// @Failure 500 {object} response.JSON "Internal server error"
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(ctx fiber.Ctx) error {
	var req request.RegisterRequest
	if err := ctx.Bind().Body(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.IAuthService.Register(req)
	if err != nil {
		util.HandleError(ctx, fiber.StatusInternalServerError, err)
		return nil
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "user registered successfully",
		Data: response.UserInfo{
			Id:              result.Id,
			Username:        result.Username,
			Name:            result.Name,
			Email:           result.Email,
			Status:          result.Status,
			EmailVerifiedAt: result.EmailVerifiedAt.Format(time.RFC3339),
		},
	})
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return a JWT token in a cookie
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "Login request"
// @Success 200 {object} response.AuthJSON
// @Failure 400 {object} response.JSON "Invalid request format"
// @Failure 401 {object} response.JSON "Authentication failed"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(ctx fiber.Ctx) error {
	var req request.LoginRequest
	if err := ctx.Bind().Body(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.IAuthService.Login(req)
	if err != nil {
		util.HandleError(ctx, fiber.StatusUnauthorized, err)
		return nil
	}

	h.Cookie.Set(ctx, "refresh_token", result.RefreshToken)

	return ctx.Status(fiber.StatusOK).JSON(response.AuthJSON{
		Message: "you are authenticated",
		Status:  fiber.StatusOK,
		User: response.UserInfo{
			Id:              result.User.Id,
			Username:        result.User.Username,
			Name:            result.User.Name,
			Email:           result.User.Email,
			Status:          result.User.Status,
			EmailVerifiedAt: result.User.EmailVerifiedAt.Format(time.RFC3339),
		},
		AccessToken: result.AccessToken,
	})
}

// Refresh godoc
// @Summary Refresh access token
// @Description Generate new access token and rotate refresh token
// @Tags Auth
// @Produce json
// @Success 200 {object} response.AuthJSON
// @Failure 401 {object} response.JSON "Unauthorized"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) Refresh(ctx fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.JSON{
			Status:  fiber.StatusUnauthorized,
			Message: "missing refresh token",
		})
	}

	result, err := h.IAuthService.RefreshToken(refreshToken)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.JSON{
			Status:  fiber.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	h.Cookie.Clear(ctx, "refresh_token")
	h.Cookie.Set(ctx, "refresh_token", result.RefreshToken)

	return ctx.Status(fiber.StatusOK).JSON(response.AuthJSON{
		Status:  fiber.StatusOK,
		Message: "token refreshed",
		User: response.UserInfo{
			Id:              result.User.Id,
			Username:        result.User.Username,
			Name:            result.User.Name,
			Email:           result.User.Email,
			Status:          result.User.Status,
			EmailVerifiedAt: result.User.EmailVerifiedAt.Format(time.RFC3339),
		},
		AccessToken: result.AccessToken,
	})
}

// Logout godoc
// @Summary Logout user
// @Description Logout user, blacklist refresh token, clear cookie
// @Tags Auth
// @Success 200 {object} response.JSON "Successfully logged out"
// @Failure 401 {object} response.JSON "Unauthorized"
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(ctx fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.JSON{
			Status:  fiber.StatusUnauthorized,
			Message: "unauthorized",
		})
	}

	if err := h.IAuthService.Logout(refreshToken); err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.JSON{
			Status:  fiber.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	h.Cookie.Clear(ctx, "refresh_token")

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "successfully logged out",
	})
}
