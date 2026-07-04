package handler

import (
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/request"
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/response"
	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/service"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type CompanyHandler struct {
	ICompanyService service.ICompanyService
}

func NewCompanyHandler(serv service.ICompanyService) *CompanyHandler {
	return &CompanyHandler{ICompanyService: serv}
}

// Create a New Company
// @Summary Create company
// @Description Store a new company record
// @Tags Companies
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.CompanyCreateRequest true "Company Create Request"
// @Success 201 {object} response.JSON "A new record has been stored."
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/companies [post]
func (h *CompanyHandler) Create(ctx fiber.Ctx) error {
	userId, err := util.GetAuthor(ctx)
	req := request.CompanyCreateRequest{}
	if err := ctx.Bind().Body(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.ICompanyService.Create(req, userId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.JSON{
		Status:  fiber.StatusCreated,
		Message: "A new record has been stored.",
		Data:    result,
	})
}

// Find All Companies
// @Summary Get all companies
// @Description Retrieve all company records with pagination
// @Tags Companies
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param sort query string false "Sort column (name, legal_name, email, status, created_at, updated_at)"
// @Param order query string false "Sort direction (asc, desc)"
// @Param status query string false "Filter by status"
// @Param verified query string false "Filter by email verified (true, false)"
// @Success 200 {object} response.JSON "Successfully retrieved all records."
// @Failure 500 {object} response.JSON "Internal Server Error"
// @Router /api/v1/companies [get]
func (h *CompanyHandler) FindAll(ctx fiber.Ctx) error {
	qp := util.ParseQueryParams(ctx, entity.Company{}.SearchableFields())

	entities, totalCount, err := h.ICompanyService.FindAll(qp)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to retrieve records",
			Errors:  err.Error(),
		})
	}

	if totalCount == 0 || (qp.Page-1)*qp.PageSize >= totalCount {
		return ctx.Status(fiber.StatusOK).JSON(response.JSON{
			Status:  fiber.StatusOK,
			Message: "No records found.",
			Data:    []response.CompanyResponse{},
			Meta:    nil,
		})
	}

	baseURL := ctx.Protocol() + "://" + ctx.Hostname() + ctx.Path()
	meta := util.GenerateMeta(baseURL, qp, totalCount)

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Successfully retrieved all records.",
		Data:    entities,
		Meta:    &meta,
	})
}

// Find Company by Id
// @Summary Get company by ID
// @Description Retrieve a single company by its ID
// @Tags Companies
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Success 200 {object} response.JSON "Successfully retrieved selected record."
// @Failure 404 {object} response.JSON "Company not found"
// @Router /api/v1/companies/{id} [get]
func (h *CompanyHandler) FindById(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.ICompanyService.FindById(parsedId)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status:  fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Successfully retrieved selected record.",
		Data:    result,
	})
}

// Update Company by Id
// @Summary Update company
// @Description Update company data by ID
// @Tags Companies
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param request body request.CompanyUpdateRequest true "Company Update Request"
// @Success 200 {object} response.JSON "Selected record has been updated."
// @Failure 404 {object} response.JSON "Company not found"
// @Router /api/v1/companies/{id} [put]
func (h *CompanyHandler) Update(ctx fiber.Ctx) error {
	req := request.CompanyUpdateRequest{}
	if err := ctx.Bind().Body(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	parsedId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	req.Id = parsedId

	result, err := h.ICompanyService.Update(req)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status:  fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Selected record has been updated.",
		Data:    result,
	})
}

// Delete Company by Id
// @Summary Delete company
// @Description Remove a company record by ID
// @Tags Companies
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Success 200 {object} response.JSON "Selected record has been deleted."
// @Failure 404 {object} response.JSON "Company not found"
// @Router /api/v1/companies/{id} [delete]
func (h *CompanyHandler) Delete(ctx fiber.Ctx) error {
	parsedId, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.ICompanyService.Delete(parsedId)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status:  fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Selected record has been deleted.",
		Data:    result,
	})
}
