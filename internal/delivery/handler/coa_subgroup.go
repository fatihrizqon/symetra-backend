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

type COASubGroupHandler struct {
	ICOASubGroupService service.ICOASubGroupService
}

func NewCOASubGroupHandler(serv service.ICOASubGroupService) *COASubGroupHandler {
	return &COASubGroupHandler{ICOASubGroupService: serv}
}

// Create a New COA SubGroup
// @Summary Create coa subgroup
// @Description Store a new coa subgroup record
// @Tags COASubGroup
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.COASubGroupCreateRequest true "COA SubGroup Create Request"
// @Success 201 {object} response.JSON "A new record has been stored."
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/coa-subgroups [post]
func (h *COASubGroupHandler) Create(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	req := request.COASubGroupCreateRequest{}
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.ICOASubGroupService.Create(companyID, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.JSON{
		Status:  fiber.StatusCreated,
		Message: "A new record has been stored.",
		Data:    response.NewCOASubGroupResponse(result),
	})
}

// Find All COASubGroup
// @Summary Get all coa subgroups
// @Description Retrieve all coa subgroup records with pagination
// @Tags COASubGroup
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
// @Router /api/v1/coa-subgroups [get]
func (h *COASubGroupHandler) FindAll(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	var entity = entity.COASubGroup{}
	qp := util.ParseQueryParams(ctx, entity.SearchableFields())
	results, totalCount, err := h.ICOASubGroupService.FindAll(companyID, qp)
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
		Data:    response.NewCOASubGroupResponses(results),
		Meta:    &meta,
	})
}

// Find COA SubGroup by Id
// @Summary Get coa subgroup by ID
// @Description Retrieve a single coa subgroup by its ID
// @Tags COASubGroup
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "COA SubGroup ID"
// @Success 200 {object} response.JSON "Successfully retrieved selected record."
// @Failure 404 {object} response.JSON "COA SubGroup not found"
// @Router /api/v1/coa-subgroups/{id} [get]
func (h *COASubGroupHandler) FindById(ctx fiber.Ctx) error {
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

	result, err := h.ICOASubGroupService.FindById(companyID, id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status:  fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Successfully retrieved selected record.",
		Data:    response.NewCOASubGroupResponse(result),
	})
}

// Update COA SubGroup by Id
// @Summary Update coa subgroup
// @Description Update coa subgroup data by ID
// @Tags COASubGroup
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "COA SubGroup ID"
// @Param request body request.COASubGroupUpdateRequest true "COA SubGroup Update Request"
// @Success 200 {object} response.JSON "Selected record has been updated."
// @Failure 404 {object} response.JSON "COA SubGroup not found"
// @Router /api/v1/coa-subgroups/{id} [put]
func (h *COASubGroupHandler) Update(ctx fiber.Ctx) error {
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

	req := request.COASubGroupUpdateRequest{}
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	req.Id = id
	result, err := h.ICOASubGroupService.Update(companyID, req)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status:  fiber.StatusNotFound,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Selected record has been updated.",
		Data:    response.NewCOASubGroupResponse(result),
	})
}

// Delete COA SubGroup by Id
// @Summary Delete coa subgroup
// @Description Remove a coa subgroup record by ID
// @Tags COASubGroup
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "COA SubGroup ID"
// @Success 200 {object} response.JSON "Selected record has been deleted."
// @Failure 404 {object} response.JSON "COA SubGroup not found"
// @Router /api/v1/coa-subgroups/{id} [delete]
func (h *COASubGroupHandler) Delete(ctx fiber.Ctx) error {
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

	err = h.ICOASubGroupService.Delete(companyID, id)
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
// @Summary Get coa subgroup dropdown list
// @Description Retrieve a list of coa subgroups for dropdown selection
// @Tags COASubGroup
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Success 200 {object} response.JSON "Success"
// @Failure 500 {object} response.JSON "Internal Server Error"
// @Router /api/v1/dropdown/coa-subgroups [get]
func (h *COASubGroupHandler) SelectDropdownList(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	qp := util.ParseQueryParams(ctx, entity.COASubGroup{}.SearchableFields())
	items, totalCount, err := h.ICOASubGroupService.SelectDropdownList(companyID, qp)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{Status: fiber.StatusInternalServerError, Message: "Failed to retrieve records", Errors: err.Error()})
	}
	if totalCount == 0 {
		return ctx.Status(fiber.StatusOK).JSON(response.SelectJSON{Data: []response.SelectDropdownListResponse{}})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.SelectJSON{Data: items})
}

// Bulk Destroy COASubGroup
// @Summary Bulk delete coa subgroups
// @Description Delete multiple coa subgroups by their IDs
// @Tags COASubGroup
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.BulkActionRequest true "Bulk Action Request"
// @Success 200 {object} response.JSON "COASubGroup destroyed"
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/coa-subgroups/destroy [delete]
func (h *COASubGroupHandler) Destroy(ctx fiber.Ctx) error {
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

	if err := h.ICOASubGroupService.Destroy(companyID, req.Ids); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Records have been deleted."})
}
