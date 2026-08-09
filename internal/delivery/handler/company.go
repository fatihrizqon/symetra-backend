package handler

import (
	"errors"

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
	authorID, err := util.GetAuthorID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	req := request.CompanyCreateRequest{}
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.ICompanyService.Create(req, authorID)
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
	var entity = entity.Company{}
	qp := util.ParseQueryParams(ctx, entity.SearchableFields())
	results, totalCount, err := h.ICompanyService.FindAll(qp)
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
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.ICompanyService.FindById(id)
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

// FindMyCompanies godoc
// @Summary Get all companies the authenticated user belongs to
// @Tags Company
// @Security BearerAuth
// @Router /api/v1/companies/mine [get]
func (h *CompanyHandler) FindMyCompanies(ctx fiber.Ctx) error {
	authorID, err := util.GetAuthorID(ctx)

	if err != nil {
		util.HandleError(ctx, fiber.StatusUnauthorized, err)
		return nil
	}

	result, err := h.ICompanyService.FindMyCompanies(authorID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{
			Status: fiber.StatusInternalServerError, Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status: fiber.StatusOK, Message: "Successfully retrieved your companies.", Data: result,
	})
}

// FindMembersByCompany godoc
// @Summary Get all companies the authenticated user belongs to
// @Tags Company
// @Security BearerAuth
// @Router /api/v1/companies/mine [get]
func (h *CompanyHandler) FindMembersByCompany(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	parsedId, err := uuid.Parse(id)

	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.ICompanyService.FindMembersByCompany(parsedId)
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

// AssignMember godoc
// @Summary Assign member to company
// @Tags Company
// @Security BearerAuth
// @Router /api/v1/companies/{id}/members [post]
func (h *CompanyHandler) AssignMember(ctx fiber.Ctx) error {
	authorID, err := util.GetAuthorID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusUnauthorized, err)
		return nil
	}

	companyID := ctx.Params("id")
	parsedCompanyId, err := uuid.Parse(companyID)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	req := request.AssignMemberRequest{}
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	req.CompanyId = parsedCompanyId

	result, err := h.ICompanyService.AssignMember(req, authorID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{
			Status: fiber.StatusInternalServerError, Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status: fiber.StatusOK, Message: "Successfully assigned member to company.", Data: result,
	})
}

// UpdateMemberRole godoc
// @Summary Update member role
// @Tags Company
// @Security BearerAuth
// @Router /api/v1/companies/{id}/members/{userId}/role [put]
func (h *CompanyHandler) UpdateMemberRole(ctx fiber.Ctx) error {
	companyID := ctx.Params("id")
	userId := ctx.Params("user_id")

	parsedCompanyId, err := uuid.Parse(companyID)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	req := request.UpdateMemberRoleRequest{}
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	req.CompanyId = parsedCompanyId
	req.UserID = parsedUserId

	err = h.ICompanyService.UpdateMemberRole(req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{
			Status: fiber.StatusInternalServerError, Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status: fiber.StatusOK, Message: "Successfully updated member role.",
	})
}

// RemoveMember godoc
// @Summary Remove member from company
// @Tags Company
// @Security BearerAuth
// @Router /api/v1/companies/{id}/members/{userId} [delete]
func (h *CompanyHandler) RemoveMember(ctx fiber.Ctx) error {
	companyID := ctx.Params("id")
	userId := ctx.Params("user_id")

	parsedCompanyId, err := uuid.Parse(companyID)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	err = h.ICompanyService.RemoveMember(parsedCompanyId, parsedUserId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{
			Status: fiber.StatusInternalServerError, Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status: fiber.StatusOK, Message: "Successfully removed member from company.",
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
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	req.Id = id

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
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	err = h.ICompanyService.Delete(id)
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

// SelectCompany godoc
// @Summary Select or switch active company
// @Description Set the active company for the current session. Call this again to switch company.
// @Tags Companies
// @Security BearerAuth
// @Produce json
// @Param id path string true "Company ID"
// @Success 200 {object} response.JSON "Active company has been set."
// @Failure 400 {object} response.JSON "Invalid company ID"
// @Failure 401 {object} response.JSON "Unauthorized"
// @Failure 403 {object} response.JSON "Not a member of this company"
// @Router /api/v1/companies/{id}/select [post]
func (h *CompanyHandler) SelectCompany(ctx fiber.Ctx) error {
	userID, err := util.GetAuthorID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusUnauthorized, err)
		return nil
	}

	sessionID, ok := ctx.Locals("session_id").(uuid.UUID)
	if !ok || sessionID == uuid.Nil {
		util.HandleError(ctx, fiber.StatusUnauthorized, errors.New("invalid session"))
		return nil
	}

	companyID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, errors.New("invalid company id"))
		return nil
	}

	result, err := h.ICompanyService.SelectCompany(sessionID, companyID, userID)
	if err != nil {
		util.HandleError(ctx, fiber.StatusForbidden, err)
		return nil
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Active company has been set.",
		Data:    result,
	})
}

// GetActiveCompany godoc
// @Summary Get current active company
// @Description Returns the company that is currently active in this session
// @Tags Companies
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.JSON "Successfully retrieved active company."
// @Failure 401 {object} response.JSON "Unauthorized"
// @Failure 404 {object} response.JSON "No active company selected"
// @Router /api/v1/companies/active [get]
func (h *CompanyHandler) GetActiveCompany(ctx fiber.Ctx) error {
	sessionID, ok := ctx.Locals("session_id").(uuid.UUID)
	if !ok || sessionID == uuid.Nil {
		util.HandleError(ctx, fiber.StatusUnauthorized, errors.New("invalid session"))
		return nil
	}

	result, err := h.ICompanyService.GetActiveCompany(sessionID)
	if err != nil {
		util.HandleError(ctx, fiber.StatusNotFound, err)
		return nil
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Successfully retrieved active company.",
		Data:    result,
	})
}

func (h *CompanyHandler) Destroy(ctx fiber.Ctx) error {
	req := request.BulkActionRequest{}
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	if err := h.ICompanyService.Destroy(req.Ids); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Records have been deleted.",
	})
}
