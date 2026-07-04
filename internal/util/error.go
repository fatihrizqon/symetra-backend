package util

import (
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/response"
	"github.com/gofiber/fiber/v3"
)

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func HandleError(ctx fiber.Ctx, status int, err error) {
	if err != nil {
		ctx.Status(status).JSON(response.JSON{
			Status:  status,
			Message: err.Error(),
		})
	}
}
