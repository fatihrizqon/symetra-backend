package handler

import (
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/response"
	"github.com/fatihrizqon/gofiber-microservice/internal/service"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/gofiber/fiber/v3"
)

type FileHandler struct {
	IFileService service.IFileService
}

func NewFileHandler(serv service.IFileService) *FileHandler {
	return &FileHandler{IFileService: serv}
}

// Upload File
// @Summary Upload a file
// @Description Upload a file (image, pdf, xlsx). Max 5MB.
// @Tags Files
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 201 {object} response.JSON "File uploaded successfully."
// @Failure 400 {object} response.JSON "Bad request"
// @Failure 401 {object} response.JSON "Unauthorized"
// @Router /api/v1/files/upload [post]
func (h *FileHandler) Upload(ctx fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  fiber.StatusBadRequest,
			Message: "No file provided in request",
		})
	}

	claims, ok := ctx.Locals("auth").(*util.Claims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.JSON{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	result, err := h.IFileService.Upload(ctx.RequestCtx(), file, claims.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.JSON{
		Status:  fiber.StatusCreated,
		Message: "File uploaded successfully.",
		Data:    result,
	})
}

// fiber:context-methods migrated
