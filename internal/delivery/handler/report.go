package handler

import (
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/response"
	"github.com/fatihrizqon/gofiber-microservice/internal/service"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ReportHandler struct {
	Service service.IReportService
}

func NewReportHandler(serv service.IReportService) *ReportHandler {
	return &ReportHandler{Service: serv}
}

func (h *ReportHandler) GeneralLedger(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	dateFrom := ctx.Query("date_from")
	dateTo := ctx.Query("date_to")
	coaIdParam := ctx.Query("coa_id")
	coaId, err := uuid.Parse(coaIdParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: "invalid coa_id"})
	}

	result, err := h.Service.GetGeneralLedger(companyID, coaId, dateFrom, dateTo)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "General Ledger fetched successfully", Data: result})
}

func (h *ReportHandler) TrialBalance(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	dateFrom := ctx.Query("date_from")
	dateTo := ctx.Query("date_to")

	result, err := h.Service.GetTrialBalance(companyID, dateFrom, dateTo)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Trial Balance fetched successfully", Data: result})
}

func (h *ReportHandler) ProfitLoss(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	dateFrom := ctx.Query("date_from")
	dateTo := ctx.Query("date_to")

	result, err := h.Service.GetProfitLoss(companyID, dateFrom, dateTo)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Profit & Loss fetched successfully", Data: result})
}

func (h *ReportHandler) BalanceSheet(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	dateTo := ctx.Query("date_to")

	result, err := h.Service.GetBalanceSheet(companyID, dateTo)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Balance Sheet fetched successfully", Data: result})
}

func (h *ReportHandler) CashFlow(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	dateFrom := ctx.Query("date_from")
	dateTo := ctx.Query("date_to")

	result, err := h.Service.GetCashFlow(companyID, dateFrom, dateTo)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Cash Flow fetched successfully", Data: result})
}

func (h *ReportHandler) EquityStatement(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	dateFrom := ctx.Query("date_from")
	dateTo := ctx.Query("date_to")

	result, err := h.Service.GetEquityStatement(companyID, dateFrom, dateTo)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Equity Statement fetched successfully", Data: result})
}

func (h *ReportHandler) JournalBook(ctx fiber.Ctx) error {
	companyID, err := util.GetCompanyID(ctx)
	if err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	dateFrom := ctx.Query("date_from")
	dateTo := ctx.Query("date_to")

	result, err := h.Service.GetJournalBook(companyID, dateFrom, dateTo)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{Status: fiber.StatusBadRequest, Message: err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{Status: fiber.StatusOK, Message: "Journal Book fetched successfully", Data: result})
}
