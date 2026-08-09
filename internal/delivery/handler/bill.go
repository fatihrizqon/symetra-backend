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

type BillHandler struct {
	IBillService service.IBillService
}

func NewBillHandler(serv service.IBillService) *BillHandler {
	return &BillHandler{IBillService: serv}
}

// Create a New Bill
// @Summary Create bill
// @Description Store a new bill record
// @Tags Bills
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.BillCreateRequest true "Bill Create Request"
// @Success 201 {object} response.JSON "A new record has been stored."
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/bills [post]
func (h *BillHandler) Create(ctx fiber.Ctx) error {
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

	req := request.BillCreateRequest{}
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.IBillService.Create(companyID, authorID, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.JSON{
		Status:  fiber.StatusCreated,
		Message: "A new record has been stored.",
		Data:    response.NewBillResponse(result),
	})
}

// Find All Bills
// @Summary Get all bills
// @Description Retrieve all bill records with pagination
// @Tags Bills
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
// @Router /api/v1/bills [get]
func (h *BillHandler) FindAll(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	var entity = entity.Bill{}
	qp := util.ParseQueryParams(ctx, entity.SearchableFields())
	results, totalCount, err := h.IBillService.FindAll(companyID, qp)
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
		Data:    response.NewBillResponses(results),
		Meta:    &meta,
	})
}

// Find Bill by Id
// @Summary Get bill by ID
// @Description Retrieve a single bill by its ID
// @Tags Bills
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Bill ID"
// @Success 200 {object} response.JSON "Successfully retrieved selected record."
// @Failure 404 {object} response.JSON "Bill not found"
// @Router /api/v1/bills/{id} [get]
func (h *BillHandler) FindById(ctx fiber.Ctx) error {
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

	result, err := h.IBillService.FindById(companyID, id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status:  fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Successfully retrieved selected record.",
		Data:    response.NewBillResponse(result),
	})
}

// Update Bill by Id
// @Summary Update bill
// @Description Update bill data by ID
// @Tags Bills
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Bill ID"
// @Param request body request.BillUpdateRequest true "Bill Update Request"
// @Success 200 {object} response.JSON "Selected record has been updated."
// @Failure 404 {object} response.JSON "Bill not found"
// @Router /api/v1/bills/{id} [put]
func (h *BillHandler) Update(ctx fiber.Ctx) error {
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

	req := request.BillUpdateRequest{}
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	req.Id = id
	result, err := h.IBillService.Update(companyID, req)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status:  fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Selected record has been updated.",
		Data:    response.NewBillResponse(result),
	})
}

// Delete Bill by Id
// @Summary Delete bill
// @Description Remove a bill record by ID
// @Tags Bills
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Bill ID"
// @Success 200 {object} response.JSON "Selected record has been deleted."
// @Failure 404 {object} response.JSON "Bill not found"
// @Router /api/v1/bills/{id} [delete]
func (h *BillHandler) Delete(ctx fiber.Ctx) error {
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

	err = h.IBillService.Delete(companyID, id)
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
// @Summary Get bill dropdown list
// @Description Retrieve a list of bills for dropdown selection
// @Tags Bills
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Success 200 {object} response.JSON "Success"
// @Failure 500 {object} response.JSON "Internal Server Error"
// @Router /api/v1/dropdown/bills [get]
func (h *BillHandler) SelectDropdownList(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	qp := util.ParseQueryParams(ctx, []string{"bill_number"})
	qp.PageSize = 100 // default higher limit for dropdown
	entities, _, err := h.IBillService.FindAll(companyID, qp)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{Status: fiber.StatusInternalServerError, Message: "Failed to retrieve records", Errors: err.Error()})
	}
	var data []map[string]interface{}
	for _, v := range entities {
		data = append(data, map[string]interface{}{
			"id":          v.Id,
			"bill_number": v.BillNumber,
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Success", Data: data})
}

// Bulk Destroy Bills
// @Summary Bulk delete bills
// @Description Delete multiple bills by their IDs
// @Tags Bills
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.BulkActionRequest true "Bulk Action Request"
// @Success 200 {object} response.JSON "Bills destroyed"
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/bills/destroy [delete]
func (h *BillHandler) Destroy(ctx fiber.Ctx) error {
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
	err = h.IBillService.Destroy(companyID, req.Ids)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Bills destroyed"})
}

// Confirm Bill
// @Summary Confirm a bill
// @Description Mark a bill as confirmed and generate related journal entries
// @Tags Bills
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Bill ID"
// @Success 200 {object} response.JSON "Bill confirmed"
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/bills/{id}/confirm [put]
func (h *BillHandler) Confirm(ctx fiber.Ctx) error {
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
	err = h.IBillService.Confirm(companyID, id, userID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Bill confirmed"})
}

// Pay Bill
// @Summary Record bill payment
// @Description Record a payment against a confirmed bill
// @Tags Bills
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Bill ID"
// @Param request body request.BillPaymentRequest true "Bill Payment Request"
// @Success 200 {object} response.JSON "Bill payment recorded"
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/bills/{id}/pay [post]
func (h *BillHandler) Pay(ctx fiber.Ctx) error {
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
	var req request.BillPaymentRequest
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	err = h.IBillService.RecordPayment(companyID, id, userID, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Bill payment recorded"})
}
