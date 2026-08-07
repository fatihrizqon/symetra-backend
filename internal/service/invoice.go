package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/request"
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/response"
	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/google/uuid"
)

type IInvoiceService interface {
	Create(companyId, userId uuid.UUID, req request.InvoiceCreateRequest) (response.InvoiceResponse, error)
	FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]response.InvoiceResponse, int, error)
	FindById(companyId, id uuid.UUID) (response.InvoiceResponse, error)
	Update(companyId, id uuid.UUID, req request.InvoiceUpdateRequest) (response.InvoiceResponse, error)
	Delete(companyId, id uuid.UUID) error
	Destroy(companyId uuid.UUID, ids []uuid.UUID) error
	Confirm(companyId, id, userId uuid.UUID) error
	RecordPayment(companyId, id, userId uuid.UUID, req request.InvoicePaymentRequest) error
}

type InvoiceService struct {
	repo           repository.IInvoiceRepository
	customerRepo   repository.ICustomerRepository
	configRepo     repository.ICompanyConfigurationRepository
	journalService IJournalEntryService
}

func NewInvoiceService(
	repo repository.IInvoiceRepository,
	customerRepo repository.ICustomerRepository,
	configRepo repository.ICompanyConfigurationRepository,
	journalService IJournalEntryService,
) IInvoiceService {
	return &InvoiceService{repo: repo, customerRepo: customerRepo, configRepo: configRepo, journalService: journalService}
}

func generateInvoiceNumber(companyId uuid.UUID, date time.Time) string {
	return fmt.Sprintf("INV-%s-%d", date.Format("20060102"), time.Now().UnixMilli())
}

func (s *InvoiceService) calculateTotals(items []entity.InvoiceItem, taxRate float64) (subtotal, discountTotal, dpp, taxAmount, grandTotal float64) {
	for i := range items {
		item := &items[i]
		item.Amount = (item.Qty * item.Price) - item.Discount
		if item.Amount < 0 {
			item.Amount = 0
		}
		subtotal += (item.Qty * item.Price)
		discountTotal += item.Discount
		if item.TaxApplicable {
			dpp += item.Amount
		}
	}
	taxAmount = dpp * taxRate
	grandTotal = (subtotal - discountTotal) + taxAmount
	return
}

func (s *InvoiceService) Create(companyId, userId uuid.UUID, req request.InvoiceCreateRequest) (response.InvoiceResponse, error) {
	_, err := s.customerRepo.FindById(companyId, req.CustomerId)
	if err != nil {
		return response.InvoiceResponse{}, errors.New("customer not found")
	}

	var items []entity.InvoiceItem
	for _, ir := range req.Items {
		items = append(items, entity.InvoiceItem{
			Description:   ir.Description,
			Qty:           ir.Qty,
			Price:         ir.Price,
			Discount:      ir.Discount,
			TaxApplicable: ir.TaxApplicable,
		})
	}

	subtotal, discountTotal, dpp, taxAmount, grandTotal := s.calculateTotals(items, req.TaxRate)

	invoice := entity.Invoice{
		CompanyId:     companyId,
		InvoiceNumber: generateInvoiceNumber(companyId, req.InvoiceDate),
		QuotationId:   req.QuotationId,
		CustomerId:    req.CustomerId,
		InvoiceDate:   req.InvoiceDate,
		DueDate:       req.DueDate,
		Subtotal:      subtotal,
		DiscountTotal: discountTotal,
		Dpp:           dpp,
		TaxRate:       req.TaxRate,
		TaxAmount:     taxAmount,
		GrandTotal:    grandTotal,
		InvoiceStatus: entity.InvoiceStatusDraft,
		TaxStatus:     entity.TaxStatusDraft,
		Notes:         req.Notes,
		CreatedBy:     userId,
		Items:         items,
	}

	if err := s.repo.Create(&invoice); err != nil {
		return response.InvoiceResponse{}, err
	}
	createdInvoice, _ := s.repo.FindById(companyId, invoice.Id)
	return response.FromInvoiceEntity(createdInvoice), nil
}

func (s *InvoiceService) FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]response.InvoiceResponse, int, error) {
	invoices, total, err := s.repo.FindAll(companyId, qp)
	if err != nil {
		return nil, 0, err
	}
		resps := make([]response.InvoiceResponse, 0, len(invoices))
	for _, i := range invoices {
		resp := response.FromInvoiceEntity(i)
		resps = append(resps, resp)
	}

	return resps, int(total), nil
}

func (s *InvoiceService) FindById(companyId, id uuid.UUID) (response.InvoiceResponse, error) {
	invoice, err := s.repo.FindById(companyId, id)
	if err != nil {
		return response.InvoiceResponse{}, errors.New("invoice not found")
	}
	resp := response.FromInvoiceEntity(invoice)
	return resp, nil
}

func (s *InvoiceService) Update(companyId, id uuid.UUID, req request.InvoiceUpdateRequest) (response.InvoiceResponse, error) {
	invoice, err := s.repo.FindById(companyId, id)
	if err != nil {
		return response.InvoiceResponse{}, errors.New("invoice not found")
	}

	if invoice.InvoiceStatus != entity.InvoiceStatusDraft {
		return response.InvoiceResponse{}, errors.New("only draft invoices can be updated")
	}

	_, err = s.customerRepo.FindById(companyId, req.CustomerId)
	if err != nil {
		return response.InvoiceResponse{}, errors.New("customer not found")
	}

	var items []entity.InvoiceItem
	for _, ir := range req.Items {
		items = append(items, entity.InvoiceItem{
			InvoiceId:     id,
			Description:   ir.Description,
			Qty:           ir.Qty,
			Price:         ir.Price,
			Discount:      ir.Discount,
			TaxApplicable: ir.TaxApplicable,
		})
	}

	subtotal, discountTotal, dpp, taxAmount, grandTotal := s.calculateTotals(items, req.TaxRate)

	invoice.CustomerId = req.CustomerId
	invoice.QuotationId = req.QuotationId
	invoice.InvoiceDate = req.InvoiceDate
	invoice.DueDate = req.DueDate
	invoice.TaxRate = req.TaxRate
	invoice.Notes = req.Notes
	invoice.Items = items
	invoice.Subtotal = subtotal
	invoice.DiscountTotal = discountTotal
	invoice.Dpp = dpp
	invoice.TaxAmount = taxAmount
	invoice.GrandTotal = grandTotal

	if err := s.repo.Update(&invoice); err != nil {
		return response.InvoiceResponse{}, err
	}
	updatedInvoice, _ := s.repo.FindById(companyId, invoice.Id)
	return response.FromInvoiceEntity(updatedInvoice), nil
}

func (s *InvoiceService) Delete(companyId, id uuid.UUID) error {
	invoice, err := s.repo.FindById(companyId, id)
	if err != nil {
		return errors.New("invoice not found")
	}
	if invoice.InvoiceStatus != entity.InvoiceStatusDraft {
		return errors.New("only draft invoices can be deleted")
	}
	return s.repo.Delete(companyId, id)
}

func (s *InvoiceService) Destroy(companyId uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return errors.New("no ids provided")
	}
	return s.repo.BulkDestroy(companyId, ids)
}

func (s *InvoiceService) Confirm(companyId, id, userId uuid.UUID) error {
	invoice, err := s.repo.FindById(companyId, id)
	if err != nil {
		return errors.New("invoice not found")
	}
	if invoice.InvoiceStatus != entity.InvoiceStatusDraft {
		return errors.New("only draft invoices can be confirmed")
	}

	config, err := s.configRepo.GetByCompanyId(companyId)
	if err != nil {
		return errors.New("company configuration not found")
	}

	if config.ARAccountId == nil {
		return errors.New("accounts receivable (A/R) account is not configured in company settings")
	}
	if config.SalesRevenueAccountId == nil {
		return errors.New("sales revenue account is not configured in company settings")
	}

	var jeLines []request.JournalLineRequest

	// Debit A/R
	jeLines = append(jeLines, request.JournalLineRequest{
		CoaId:       *config.ARAccountId,
		Description: util.StringPtr(fmt.Sprintf("A/R for Invoice %s", invoice.InvoiceNumber)),
		Debit:       invoice.GrandTotal,
	})

	// Credit Sales Revenue
	revenueAmount := invoice.Subtotal - invoice.DiscountTotal
	jeLines = append(jeLines, request.JournalLineRequest{
		CoaId:       *config.SalesRevenueAccountId,
		Description: util.StringPtr(fmt.Sprintf("Sales Revenue for Invoice %s", invoice.InvoiceNumber)),
		Credit:      revenueAmount,
	})

	// Credit Tax Payable if tax amount > 0
	if invoice.TaxAmount > 0 {
		if config.TaxPayableAccountId == nil {
			return errors.New("tax payable account is not configured in company settings")
		}
		jeLines = append(jeLines, request.JournalLineRequest{
			CoaId:       *config.TaxPayableAccountId,
			Description: util.StringPtr(fmt.Sprintf("Tax for Invoice %s", invoice.InvoiceNumber)),
			Credit:      invoice.TaxAmount,
		})
	}

	jeReq := request.JournalEntryCreateRequest{
		Date:        invoice.InvoiceDate,
		Description: fmt.Sprintf("Auto-generated Journal Entry for Invoice %s", invoice.InvoiceNumber),
		Lines:       jeLines,
	}

	jeRes, err := s.journalService.Create(companyId, userId, jeReq)
	if err != nil {
		return fmt.Errorf("failed to generate journal entry: %v", err)
	}

	// Automatically post the journal entry
	if err := s.journalService.Post(companyId, jeRes.Id); err != nil {
		return fmt.Errorf("failed to post journal entry: %v", err)
	}

	return s.repo.UpdateStatus(companyId, id, entity.InvoiceStatusConfirmed, &jeRes.Id)
}

func (s *InvoiceService) RecordPayment(companyId, id, userId uuid.UUID, req request.InvoicePaymentRequest) error {
	invoice, err := s.repo.FindById(companyId, id)
	if err != nil {
		return errors.New("invoice not found")
	}
	if invoice.InvoiceStatus != entity.InvoiceStatusConfirmed {
		return errors.New("invoice must be confirmed before recording payment")
	}

	config, err := s.configRepo.GetByCompanyId(companyId)
	if err != nil {
		return errors.New("company configuration not found")
	}

	if config.ARAccountId == nil {
		return errors.New("accounts receivable (A/R) account is not configured in company settings")
	}

	// Journal Entry for Payment: Debit Payment Account (Bank/Cash), Credit A/R
	var jeLines []request.JournalLineRequest
	jeLines = append(jeLines, request.JournalLineRequest{
		CoaId:       req.PaymentAccountId,
		Description: util.StringPtr(fmt.Sprintf("Payment for Invoice %s", invoice.InvoiceNumber)),
		Debit:       req.Amount,
	})
	jeLines = append(jeLines, request.JournalLineRequest{
		CoaId:       *config.ARAccountId,
		Description: util.StringPtr(fmt.Sprintf("Payment for Invoice %s", invoice.InvoiceNumber)),
		Credit:      req.Amount,
	})

	jeReq := request.JournalEntryCreateRequest{
		Date:        req.PaymentDate,
		Description: fmt.Sprintf("Payment for Invoice %s", invoice.InvoiceNumber),
		Lines:       jeLines,
	}

	jeRes, err := s.journalService.Create(companyId, userId, jeReq)
	if err != nil {
		return fmt.Errorf("failed to generate payment journal entry: %v", err)
	}

	if err := s.journalService.Post(companyId, jeRes.Id); err != nil {
		return fmt.Errorf("failed to post payment journal entry: %v", err)
	}

	payment := entity.InvoicePayment{
		InvoiceId:        id,
		Amount:           req.Amount,
		PaymentDate:      req.PaymentDate,
		PaymentAccountId: req.PaymentAccountId,
		JournalEntryId:   &jeRes.Id,
		Notes:            req.Notes,
		CreatedBy:        userId,
	}

	if err := s.repo.AddPayment(&payment); err != nil {
		return err
	}

	newAmountPaid := invoice.AmountPaid + req.Amount
	var pStatus entity.PaymentStatus = entity.PaymentStatusPartial
	if newAmountPaid >= invoice.GrandTotal {
		pStatus = entity.PaymentStatusPaid
	}

	return s.repo.UpdatePaymentStatus(companyId, id, newAmountPaid, pStatus)
}
