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

type QuotationHandler struct {
	IQuotationService service.IQuotationService
}

func NewQuotationHandler(serv service.IQuotationService) *QuotationHandler {
	return &QuotationHandler{IQuotationService: serv}
}

// Create a New Quotation
// @Summary Create quotation
// @Description Store a new quotation record
// @Tags Quotations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.QuotationCreateRequest true "Quotation Create Request"
// @Success 201 {object} response.JSON "A new record has been stored."
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/quotations [post]
func (h *QuotationHandler) Create(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	authorID, err := util.GetAuthorID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	req := request.QuotationCreateRequest{}
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.IQuotationService.Create(companyID, authorID, req)
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

// Find All Quotations
// @Summary Get all quotations
// @Description Retrieve all quotation records with pagination
// @Tags Quotations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param sort query string false "Sort column"
// @Param order query string false "Sort direction (asc, desc)"
// @Success 200 {object} response.JSON "Successfully retrieved all records."
// @Failure 500 {object} response.JSON "Internal Server Error"
// @Router /api/v1/quotations [get]
func (h *QuotationHandler) FindAll(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	var entity = entity.Quotation{}
	qp := util.ParseQueryParams(ctx, entity.SearchableFields())
	results, totalCount, err := h.IQuotationService.FindAll(companyID, qp)
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
		})
	}

	baseURL := ctx.Protocol() + "://" + ctx.Hostname() + ctx.Path()
	meta := util.GenerateMeta(baseURL, qp, totalCount)
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Successfully retrieved all records.",
		Data:    results,
		Meta:    &meta,
	})
}

// Select Dropdown List
// @Summary Get quotation dropdown list
// @Description Retrieve a list of quotations for dropdown selection
// @Tags Quotations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Success 200 {object} response.JSON "Success"
// @Failure 500 {object} response.JSON "Internal Server Error"
// @Router /api/v1/dropdown/quotations [get]
func (h *QuotationHandler) SelectDropdownList(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	qp := util.ParseQueryParams(ctx, []string{"quotation_number"})
	qp.PageSize = 100 // default higher limit for dropdown
	entities, _, err := h.IQuotationService.FindAll(companyID, qp)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{Status: fiber.StatusInternalServerError, Message: "Failed to retrieve records", Errors: err.Error()})
	}
	var data []map[string]interface{}
	for _, v := range entities {
		data = append(data, map[string]interface{}{
			"id":               v.Id,
			"quotation_number": v.QuotationNumber,
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Success", Data: data})
}

// Find Quotation by Id
// @Summary Get quotation by ID
// @Description Retrieve a single quotation by its ID
// @Tags Quotations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Quotation ID"
// @Success 200 {object} response.JSON "Successfully retrieved selected record."
// @Failure 404 {object} response.JSON "Quotation not found"
// @Router /api/v1/quotations/{id} [get]
func (h *QuotationHandler) FindById(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.IQuotationService.FindById(companyID, id)
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

// Update Quotation by Id
// @Summary Update quotation
// @Description Update quotation data by ID
// @Tags Quotations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Quotation ID"
// @Param request body request.QuotationUpdateRequest true "Quotation Update Request"
// @Success 200 {object} response.JSON "Selected record has been updated."
// @Failure 404 {object} response.JSON "Quotation not found"
// @Router /api/v1/quotations/{id} [put]
func (h *QuotationHandler) Update(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	req := request.QuotationUpdateRequest{}
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	req.Id = id
	result, err := h.IQuotationService.Update(companyID, req)
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

// Delete Quotation by Id
// @Summary Delete quotation
// @Description Remove a quotation record by ID
// @Tags Quotations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Quotation ID"
// @Success 200 {object} response.JSON "Selected record has been deleted."
// @Failure 404 {object} response.JSON "Quotation not found"
// @Router /api/v1/quotations/{id} [delete]
func (h *QuotationHandler) Delete(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	err = h.IQuotationService.Delete(companyID, id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Quotation deleted"})
}

// Bulk Destroy Quotations
// @Summary Bulk delete quotations
// @Description Delete multiple quotations by their IDs
// @Tags Quotations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.BulkActionRequest true "Bulk Action Request"
// @Success 200 {object} response.JSON "Quotations destroyed"
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/quotations/destroy [delete]
func (h *QuotationHandler) Destroy(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	var req request.BulkActionRequest
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	err = h.IQuotationService.Destroy(companyID, req.Ids)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Quotations destroyed"})
}

// Approve Quotation
// @Summary Approve a quotation
// @Description Mark a quotation as approved
// @Tags Quotations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Quotation ID"
// @Success 200 {object} response.JSON "Quotation approved"
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/quotations/{id}/approve [put]
func (h *QuotationHandler) Approve(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	err = h.IQuotationService.Approve(companyID, id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Quotation approved"})
}
