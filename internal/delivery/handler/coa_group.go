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

type COAGroupHandler struct {
	ICOAGroupService service.ICOAGroupService
}

func NewCOAGroupHandler(serv service.ICOAGroupService) *COAGroupHandler {
	return &COAGroupHandler{ICOAGroupService: serv}
}

// Create a New COA Group
// @Summary Create coa group
// @Description Store a new coa group record
// @Tags COAGroup
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.COAGroupCreateRequest true "COA Group Create Request"
// @Success 201 {object} response.JSON "A new record has been stored."
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/coa-groups [post]
func (h *COAGroupHandler) Create(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	req := request.COAGroupCreateRequest{}
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.ICOAGroupService.Create(companyID, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.JSON{
		Status:  fiber.StatusCreated,
		Message: "A new record has been stored.",
		Data:    response.NewCOAGroupResponse(result),
	})
}

// Find All COAGroup
// @Summary Get all coa groups
// @Description Retrieve all coa group records with pagination
// @Tags COAGroup
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
// @Router /api/v1/coa-groups [get]
func (h *COAGroupHandler) FindAll(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	var entity = entity.COAGroup{}
	qp := util.ParseQueryParams(ctx, entity.SearchableFields())
	results, totalCount, err := h.ICOAGroupService.FindAll(companyID, qp)
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
		Data:    response.NewCOAGroupResponses(results),
		Meta:    &meta,
	})
}

// Find COA Group by Id
// @Summary Get coa group by ID
// @Description Retrieve a single coa group by its ID
// @Tags COAGroup
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "COA Group ID"
// @Success 200 {object} response.JSON "Successfully retrieved selected record."
// @Failure 404 {object} response.JSON "COA Group not found"
// @Router /api/v1/coa-groups/{id} [get]
func (h *COAGroupHandler) FindById(ctx fiber.Ctx) error {
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

	result, err := h.ICOAGroupService.FindById(companyID, id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status:  fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Successfully retrieved selected record.",
		Data:    response.NewCOAGroupResponse(result),
	})
}

// Update COA Group by Id
// @Summary Update coa group
// @Description Update coa group data by ID
// @Tags COAGroup
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "COA Group ID"
// @Param request body request.COAGroupUpdateRequest true "COA Group Update Request"
// @Success 200 {object} response.JSON "Selected record has been updated."
// @Failure 404 {object} response.JSON "COA Group not found"
// @Router /api/v1/coa-groups/{id} [put]
func (h *COAGroupHandler) Update(ctx fiber.Ctx) error {
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

	req := request.COAGroupUpdateRequest{}
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	req.Id = id
	result, err := h.ICOAGroupService.Update(companyID, req)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status:  fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Selected record has been updated.",
		Data:    response.NewCOAGroupResponse(result),
	})
}

// Delete COA Group by Id
// @Summary Delete coa group
// @Description Remove a coa group record by ID
// @Tags COAGroup
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "COA Group ID"
// @Success 200 {object} response.JSON "Selected record has been deleted."
// @Failure 404 {object} response.JSON "COA Group not found"
// @Router /api/v1/coa-groups/{id} [delete]
func (h *COAGroupHandler) Delete(ctx fiber.Ctx) error {
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

	err = h.ICOAGroupService.Delete(companyID, id)
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
// @Summary Get coa group dropdown list
// @Description Retrieve a list of coa groups for dropdown selection
// @Tags COAGroup
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Success 200 {object} response.JSON "Success"
// @Failure 500 {object} response.JSON "Internal Server Error"
// @Router /api/v1/dropdown/coa-groups [get]
func (h *COAGroupHandler) SelectDropdownList(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	qp := util.ParseQueryParams(ctx, entity.COAGroup{}.SearchableFields())
	items, totalCount, err := h.ICOAGroupService.SelectDropdownList(companyID, qp)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{Status: fiber.StatusInternalServerError, Message: "Failed to retrieve records", Errors: err.Error()})
	}
	if totalCount == 0 {
		return ctx.Status(fiber.StatusOK).JSON(response.SelectJSON{Data: []response.SelectDropdownListResponse{}})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.SelectJSON{Data: items})
}

// Bulk Destroy COAGroup
// @Summary Bulk delete coa groups
// @Description Delete multiple coa groups by their IDs
// @Tags COAGroup
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.BulkActionRequest true "Bulk Action Request"
// @Success 200 {object} response.JSON "COAGroup destroyed"
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/coa-groups/destroy [delete]
func (h *COAGroupHandler) Destroy(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	req := request.BulkActionRequest{}
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	if err := h.ICOAGroupService.Destroy(companyID, req.Ids); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Records have been deleted."})
}
