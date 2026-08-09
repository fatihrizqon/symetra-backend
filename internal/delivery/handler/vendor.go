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

type VendorHandler struct {
	IVendorService service.IVendorService
}

func NewVendorHandler(serv service.IVendorService) *VendorHandler {
	return &VendorHandler{IVendorService: serv}
}

// Create a New Vendor
// @Summary Create vendor
// @Description Store a new vendor record
// @Tags Vendors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.VendorCreateRequest true "Vendor Create Request"
// @Success 201 {object} response.JSON "A new record has been stored."
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/vendors [post]
func (h *VendorHandler) Create(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	req := request.VendorCreateRequest{}
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.IVendorService.Create(companyID, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.JSON{
		Status:  fiber.StatusCreated,
		Message: "A new record has been stored.",
		Data:    response.NewVendorResponse(result),
	})
}

// Find All Vendors
// @Summary Get all vendors
// @Description Retrieve all vendor records with pagination
// @Tags Vendors
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
// @Router /api/v1/vendors [get]
func (h *VendorHandler) FindAll(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	var entity = entity.Vendor{}
	qp := util.ParseQueryParams(ctx, entity.SearchableFields())
	results, totalCount, err := h.IVendorService.FindAll(companyID, qp)
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
		Data:    response.NewVendorResponses(results),
		Meta:    &meta,
	})
}

// Find Vendor by Id
// @Summary Get vendor by ID
// @Description Retrieve a single vendor by its ID
// @Tags Vendors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Vendor ID"
// @Success 200 {object} response.JSON "Successfully retrieved selected record."
// @Failure 404 {object} response.JSON "Vendor not found"
// @Router /api/v1/vendors/{id} [get]
func (h *VendorHandler) FindById(ctx fiber.Ctx) error {
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

	result, err := h.IVendorService.FindById(companyID, id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status:  fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Successfully retrieved selected record.",
		Data:    response.NewVendorResponse(result),
	})
}

// Update Vendor by Id
// @Summary Update vendor
// @Description Update vendor data by ID
// @Tags Vendors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Vendor ID"
// @Param request body request.VendorUpdateRequest true "Vendor Update Request"
// @Success 200 {object} response.JSON "Selected record has been updated."
// @Failure 404 {object} response.JSON "Vendor not found"
// @Router /api/v1/vendors/{id} [put]
func (h *VendorHandler) Update(ctx fiber.Ctx) error {
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

	req := request.VendorUpdateRequest{}
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	req.Id = id
	result, err := h.IVendorService.Update(companyID, req)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status:  fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Selected record has been updated.",
		Data:    response.NewVendorResponse(result),
	})
}

// Delete Vendor by Id
// @Summary Delete vendor
// @Description Remove a vendor record by ID
// @Tags Vendors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Vendor ID"
// @Success 200 {object} response.JSON "Selected record has been deleted."
// @Failure 404 {object} response.JSON "Vendor not found"
// @Router /api/v1/vendors/{id} [delete]
func (h *VendorHandler) Delete(ctx fiber.Ctx) error {
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

	err = h.IVendorService.Delete(companyID, id)
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
// @Summary Get vendor dropdown list
// @Description Retrieve a list of vendors for dropdown selection
// @Tags Vendors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Success 200 {object} response.JSON "Success"
// @Failure 500 {object} response.JSON "Internal Server Error"
// @Router /api/v1/dropdown/vendors [get]
func (h *VendorHandler) SelectDropdownList(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	qp := util.ParseQueryParams(ctx, []string{"name"})
	qp.PageSize = 100 // default higher limit for dropdown
	entities, _, err := h.IVendorService.FindAll(companyID, qp)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{Status: fiber.StatusInternalServerError, Message: "Failed to retrieve records", Errors: err.Error()})
	}
	var data []map[string]interface{}
	for _, v := range entities {
		data = append(data, map[string]interface{}{
			"id":   v.Id,
			"name": v.Name,
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Success", Data: data})
}

// Bulk Destroy Vendors
// @Summary Bulk delete vendors
// @Description Delete multiple vendors by their IDs
// @Tags Vendors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.BulkActionRequest true "Bulk Action Request"
// @Success 200 {object} response.JSON "Vendors destroyed"
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/vendors/destroy [delete]
func (h *VendorHandler) Destroy(ctx fiber.Ctx) error {
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
	err = h.IVendorService.Destroy(companyID, req.Ids)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Vendors destroyed"})
}
