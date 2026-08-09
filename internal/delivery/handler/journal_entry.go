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

type JournalEntryHandler struct {
	IJournalEntryService service.IJournalEntryService
}

func NewJournalEntryHandler(serv service.IJournalEntryService) *JournalEntryHandler {
	return &JournalEntryHandler{IJournalEntryService: serv}
}

// Create a New Journal Entry
// @Summary Create journal entry
// @Description Store a new journal entry record
// @Tags JournalEntries
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.JournalEntryCreateRequest true "Journal Entry Create Request"
// @Success 201 {object} response.JSON "A new record has been stored."
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/journal_entries [post]
func (h *JournalEntryHandler) Create(ctx fiber.Ctx) error {
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

	req := request.JournalEntryCreateRequest{}
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	result, err := h.IJournalEntryService.Create(companyID, authorID, req)
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

// Find All JournalEntries
// @Summary Get all journal entrys
// @Description Retrieve all journal entry records with pagination
// @Tags JournalEntries
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
// @Router /api/v1/journal_entries [get]
func (h *JournalEntryHandler) FindAll(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	var entity = entity.JournalEntry{}
	qp := util.ParseQueryParams(ctx, entity.SearchableFields())
	results, totalCount, err := h.IJournalEntryService.FindAll(companyID, qp)
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

// Find Journal Entry by Id
// @Summary Get journal entry by ID
// @Description Retrieve a single journal entry by its ID
// @Tags JournalEntries
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Journal Entry ID"
// @Success 200 {object} response.JSON "Successfully retrieved selected record."
// @Failure 404 {object} response.JSON "Journal Entry not found"
// @Router /api/v1/journal_entries/{id} [get]
func (h *JournalEntryHandler) FindById(ctx fiber.Ctx) error {
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

	result, err := h.IJournalEntryService.FindById(companyID, id)
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

// Update Journal Entry by Id
// @Summary Update journal entry
// @Description Update journal entry data by ID
// @Tags JournalEntries
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Journal Entry ID"
// @Param request body request.JournalEntryUpdateRequest true "Journal Entry Update Request"
// @Success 200 {object} response.JSON "Selected record has been updated."
// @Failure 404 {object} response.JSON "Journal Entry not found"
// @Router /api/v1/journal_entries/{id} [put]
func (h *JournalEntryHandler) Update(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	req := request.JournalEntryUpdateRequest{}
	if err := util.Parse(ctx, &req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	req.Id = id
	result, err := h.IJournalEntryService.Update(companyID, req)
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

// Delete Journal Entry by Id
// @Summary Delete journal entry
// @Description Remove a journal entry record by ID
// @Tags JournalEntries
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Journal Entry ID"
// @Success 200 {object} response.JSON "Selected record has been deleted."
// @Failure 404 {object} response.JSON "Journal Entry not found"
// @Router /api/v1/journal_entries/{id} [delete]
func (h *JournalEntryHandler) Delete(ctx fiber.Ctx) error {
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

	err = h.IJournalEntryService.Delete(companyID, id)
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

// Post Journal Entry
// @Summary Post a journal entry
// @Description Post a journal entry to the ledger
// @Tags JournalEntries
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Journal Entry ID"
// @Success 200 {object} response.JSON "Journal Entry posted"
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/journal_entries/{id}/post [put]
func (h *JournalEntryHandler) Post(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	if err := h.IJournalEntryService.Post(companyID, id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}

	// reload for response
	ent, _ := h.IJournalEntryService.FindById(companyID, id)
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Journal entry posted successfully.", Data: ent})
}

// Void Journal Entry
// @Summary Void a journal entry
// @Description Void a journal entry
// @Tags JournalEntries
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Journal Entry ID"
// @Success 200 {object} response.JSON "Journal Entry voided"
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/journal_entries/{id}/void [put]
func (h *JournalEntryHandler) Void(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	if err := h.IJournalEntryService.Void(companyID, id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}

	ent, _ := h.IJournalEntryService.FindById(companyID, id)
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Journal entry voided successfully.", Data: ent})
}

// Bulk Destroy JournalEntries
// @Summary Bulk delete journal entrys
// @Description Delete multiple journal entrys by their IDs
// @Tags JournalEntries
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.BulkActionRequest true "Bulk Action Request"
// @Success 200 {object} response.JSON "JournalEntries destroyed"
// @Failure 400 {object} response.JSON "Bad request"
// @Router /api/v1/journal_entries/destroy [delete]
func (h *JournalEntryHandler) Destroy(ctx fiber.Ctx) error {
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

	if err := h.IJournalEntryService.Destroy(companyID, req.Ids); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Records have been deleted."})
}
