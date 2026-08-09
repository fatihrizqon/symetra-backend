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

type InvoiceHandler struct {
	IInvoiceService service.IInvoiceService
}

func NewInvoiceHandler(serv service.IInvoiceService) *InvoiceHandler {
	return &InvoiceHandler{IInvoiceService: serv}
}

// Create a New Invoice
// @Summary Create invoice
// @Description Store a new invoice record
// @Tags Invoices
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.InvoiceCreateRequest true "Invoice Create Request"
// @Success 201 {object} response.JSON "A new record has been stored."
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/invoices [post]
func (h *InvoiceHandler) Create(ctx fiber.Ctx) error {
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

	req := request.InvoiceCreateRequest{}
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.IInvoiceService.Create(companyID, authorID, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.JSON{
		Status:  fiber.StatusCreated,
		Message: "A new record has been stored.",
		Data:    response.NewInvoiceResponse(result),
	})
}

// Find All Invoices
// @Summary Get all invoices
// @Description Retrieve all invoice records with pagination
// @Tags Invoices
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
// @Router /api/v1/invoices [get]
func (h *InvoiceHandler) FindAll(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	var entity = entity.Invoice{}
	qp := util.ParseQueryParams(ctx, entity.SearchableFields())
	results, totalCount, err := h.IInvoiceService.FindAll(companyID, qp)
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
		Data:    response.NewInvoiceResponses(results),
		Meta:    &meta,
	})
}

// Find Invoice by Id
// @Summary Get invoice by ID
// @Description Retrieve a single invoice by its ID
// @Tags Invoices
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 200 {object} response.JSON "Successfully retrieved selected record."
// @Failure 404 {object} response.JSON "Invoice not found"
// @Router /api/v1/invoices/{id} [get]
func (h *InvoiceHandler) FindById(ctx fiber.Ctx) error {
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

	result, err := h.IInvoiceService.FindById(companyID, id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status:  fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Successfully retrieved selected record.",
		Data:    response.NewInvoiceResponse(result),
	})
}

// Update Invoice by Id
// @Summary Update invoice
// @Description Update invoice data by ID
// @Tags Invoices
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Invoice ID"
// @Param request body request.InvoiceUpdateRequest true "Invoice Update Request"
// @Success 200 {object} response.JSON "Selected record has been updated."
// @Failure 404 {object} response.JSON "Invoice not found"
// @Router /api/v1/invoices/{id} [put]
func (h *InvoiceHandler) Update(ctx fiber.Ctx) error {
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

	req := request.InvoiceUpdateRequest{}
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	req.Id = id
	result, err := h.IInvoiceService.Update(companyID, req)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status:  fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Selected record has been updated.",
		Data:    response.NewInvoiceResponse(result),
	})
}

// Delete Invoice by Id
// @Summary Delete invoice
// @Description Remove a invoice record by ID
// @Tags Invoices
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 200 {object} response.JSON "Selected record has been deleted."
// @Failure 404 {object} response.JSON "Invoice not found"
// @Router /api/v1/invoices/{id} [delete]
func (h *InvoiceHandler) Delete(ctx fiber.Ctx) error {
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

	err = h.IInvoiceService.Delete(companyID, id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Selected record has been deleted.",
	})
}

// Select Dropdown List
// @Summary Get invoice dropdown list
// @Description Retrieve a list of invoices for dropdown selection
// @Tags Invoices
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Success 200 {object} response.JSON "Success"
// @Failure 500 {object} response.JSON "Internal Server Error"
// @Router /api/v1/dropdown/invoices [get]
func (h *InvoiceHandler) SelectDropdownList(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	qp := util.ParseQueryParams(ctx, []string{"invoice_number"})
	qp.PageSize = 100 // default higher limit for dropdown
	entities, _, err := h.IInvoiceService.FindAll(companyID, qp)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{Status: fiber.StatusInternalServerError, Message: "Failed to retrieve records", Errors: err.Error()})
	}
	var data []map[string]interface{}
	for _, v := range entities {
		data = append(data, map[string]interface{}{
			"id":             v.Id,
			"invoice_number": v.InvoiceNumber,
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Success", Data: data})
}

// Bulk Destroy Invoices
// @Summary Bulk delete invoices
// @Description Delete multiple invoices by their IDs
// @Tags Invoices
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.BulkActionRequest true "Bulk Action Request"
// @Success 200 {object} response.JSON "Invoices destroyed"
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/invoices/destroy [delete]
func (h *InvoiceHandler) Destroy(ctx fiber.Ctx) error {
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
	err = h.IInvoiceService.Destroy(companyID, req.Ids)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Invoices destroyed"})
}

// Confirm Invoice
// @Summary Confirm a invoice
// @Description Mark a invoice as confirmed
// @Tags Invoices
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 200 {object} response.JSON "Invoice confirmed"
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/invoices/{id}/confirm [put]
func (h *InvoiceHandler) Confirm(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	claims, ok := ctx.Locals("auth").(*util.Claims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.JSON{Status: fiber.StatusUnauthorized, Message: "unauthorized"})
	}
	userID := claims.UserID
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	err = h.IInvoiceService.Confirm(companyID, id, userID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Invoice confirmed"})
}

// Pay Invoice
// @Summary Record invoice payment
// @Description Record a payment against a confirmed invoice
// @Tags Invoices
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Invoice ID"
// @Param request body request.InvoicePaymentRequest true "Invoice Payment Request"
// @Success 200 {object} response.JSON "Invoice payment recorded"
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/invoices/{id}/pay [post]
func (h *InvoiceHandler) Pay(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	claims, ok := ctx.Locals("auth").(*util.Claims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(response.JSON{Status: fiber.StatusUnauthorized, Message: "unauthorized"})
	}
	userID := claims.UserID
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	var req request.InvoicePaymentRequest
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	err = h.IInvoiceService.RecordPayment(companyID, id, userID, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Invoice payment recorded"})
}
