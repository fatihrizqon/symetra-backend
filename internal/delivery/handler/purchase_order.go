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

type PurchaseOrderHandler struct {
	IPurchaseOrderService service.IPurchaseOrderService
}

func NewPurchaseOrderHandler(serv service.IPurchaseOrderService) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{IPurchaseOrderService: serv}
}

// Create a New Purchase Order
// @Summary Create purchase order
// @Description Store a new purchase order record
// @Tags PurchaseOrders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.PurchaseOrderCreateRequest true "Purchase Order Create Request"
// @Success 201 {object} response.JSON "A new record has been stored."
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/purchase-orders [post]
func (h *PurchaseOrderHandler) Create(ctx fiber.Ctx) error {
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

	req := request.PurchaseOrderCreateRequest{}
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.IPurchaseOrderService.Create(companyID, authorID, req)
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

// Find All PurchaseOrders
// @Summary Get all purchase orders
// @Description Retrieve all purchase order records with pagination
// @Tags PurchaseOrders
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
// @Router /api/v1/purchase-orders [get]
func (h *PurchaseOrderHandler) FindAll(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	var entity = entity.PurchaseOrder{}
	qp := util.ParseQueryParams(ctx, entity.SearchableFields())
	results, totalCount, err := h.IPurchaseOrderService.FindAll(companyID, qp)
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

// Find Purchase Order by Id
// @Summary Get purchase order by ID
// @Description Retrieve a single purchase order by its ID
// @Tags PurchaseOrders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Order ID"
// @Success 200 {object} response.JSON "Successfully retrieved selected record."
// @Failure 404 {object} response.JSON "Purchase Order not found"
// @Router /api/v1/purchase-orders/{id} [get]
func (h *PurchaseOrderHandler) FindById(ctx fiber.Ctx) error {
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

	result, err := h.IPurchaseOrderService.FindById(companyID, id)
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

// Update Purchase Order by Id
// @Summary Update purchase order
// @Description Update purchase order data by ID
// @Tags PurchaseOrders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Order ID"
// @Param request body request.PurchaseOrderUpdateRequest true "Purchase Order Update Request"
// @Success 200 {object} response.JSON "Selected record has been updated."
// @Failure 404 {object} response.JSON "Purchase Order not found"
// @Router /api/v1/purchase-orders/{id} [put]
func (h *PurchaseOrderHandler) Update(ctx fiber.Ctx) error {
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

	req := request.PurchaseOrderUpdateRequest{}
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	req.Id = id
	result, err := h.IPurchaseOrderService.Update(companyID, req)
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

// Delete Purchase Order by Id
// @Summary Delete purchase order
// @Description Remove a purchase order record by ID
// @Tags PurchaseOrders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Order ID"
// @Success 200 {object} response.JSON "Selected record has been deleted."
// @Failure 404 {object} response.JSON "Purchase Order not found"
// @Router /api/v1/purchase-orders/{id} [delete]
func (h *PurchaseOrderHandler) Delete(ctx fiber.Ctx) error {
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

	err = h.IPurchaseOrderService.Delete(companyID, id)
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
// @Summary Get purchase order dropdown list
// @Description Retrieve a list of purchase orders for dropdown selection
// @Tags PurchaseOrders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Success 200 {object} response.JSON "Success"
// @Failure 500 {object} response.JSON "Internal Server Error"
// @Router /api/v1/dropdown/purchase-orders [get]
func (h *PurchaseOrderHandler) SelectDropdownList(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	qp := util.ParseQueryParams(ctx, []string{"po_number"})
	qp.PageSize = 100 // default higher limit for dropdown
	entities, _, err := h.IPurchaseOrderService.FindAll(companyID, qp)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{Status: fiber.StatusInternalServerError, Message: "Failed to retrieve records", Errors: err.Error()})
	}
	var data []map[string]interface{}
	for _, v := range entities {
		data = append(data, map[string]interface{}{
			"id":        v.Id,
			"po_number": v.PoNumber,
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Success", Data: data})
}

// Bulk Destroy PurchaseOrders
// @Summary Bulk delete purchase orders
// @Description Delete multiple purchase orders by their IDs
// @Tags PurchaseOrders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.BulkActionRequest true "Bulk Action Request"
// @Success 200 {object} response.JSON "PurchaseOrders destroyed"
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/purchase-orders/destroy [delete]
func (h *PurchaseOrderHandler) Destroy(ctx fiber.Ctx) error {
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
	err = h.IPurchaseOrderService.Destroy(companyID, req.Ids)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Purchase Orders destroyed"})
}

// Approve Purchase Order
// @Summary Approve a purchase order
// @Description Mark a purchase order as approved
// @Tags PurchaseOrders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Purchase Order ID"
// @Success 200 {object} response.JSON "Purchase Order approved"
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/purchase-orders/{id}/approve [put]
func (h *PurchaseOrderHandler) Approve(ctx fiber.Ctx) error {
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
	err = h.IPurchaseOrderService.Approve(companyID, id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Purchase Order approved"})
}
