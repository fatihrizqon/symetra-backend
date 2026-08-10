package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/request"
	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type IInvoiceService interface {
	Create(companyID uuid.UUID, userId uuid.UUID, req request.InvoiceCreateRequest) (entity.Invoice, error)
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.Invoice, int, error)
	FindById(companyID uuid.UUID, id uuid.UUID) (entity.Invoice, error)
	Update(companyID uuid.UUID, req request.InvoiceUpdateRequest) (entity.Invoice, error)
	Delete(companyID uuid.UUID, id uuid.UUID) error
	Destroy(companyID uuid.UUID, ids []uuid.UUID) error
	Confirm(companyID uuid.UUID, id uuid.UUID, userId uuid.UUID) error
	RecordPayment(companyID uuid.UUID, id uuid.UUID, userId uuid.UUID, req request.InvoicePaymentRequest) error
}

type InvoiceService struct {
	validate                        *validator.Validate
	IInvoiceRepository              repository.IInvoiceRepository
	ICustomerRepository             repository.ICustomerRepository
	ICompanyConfigurationRepository repository.ICompanyConfigurationRepository
	IJournalEntryService            IJournalEntryService
}

func NewInvoiceService(validate *validator.Validate, IInvoiceRepository repository.IInvoiceRepository, ICustomerRepository repository.ICustomerRepository, ICompanyConfigurationRepository repository.ICompanyConfigurationRepository, IJournalEntryService IJournalEntryService) IInvoiceService {
	return &InvoiceService{
		validate:                        validate,
		IInvoiceRepository:              IInvoiceRepository,
		ICustomerRepository:             ICustomerRepository,
		ICompanyConfigurationRepository: ICompanyConfigurationRepository,
		IJournalEntryService:            IJournalEntryService,
	}
}

func (s *InvoiceService) Create(companyID, userId uuid.UUID, req request.InvoiceCreateRequest) (entity.Invoice, error) {
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
		CompanyId:     companyID,
		InvoiceNumber: generateInvoiceNumber(req.InvoiceDate),
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

	if err := s.IInvoiceRepository.Create(&invoice); err != nil {
		return entity.Invoice{}, err
	}
	createdInvoice, _ := s.IInvoiceRepository.FindById(companyID, invoice.Id)
	return createdInvoice, nil
}

func (s *InvoiceService) FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.Invoice, int, error) {
	invoices, totalCount, err := s.IInvoiceRepository.FindAll(companyID, qp)
	if err != nil {
		return nil, 0, err
	}
	return invoices, totalCount, nil
}

func (s *InvoiceService) FindById(companyID uuid.UUID, id uuid.UUID) (entity.Invoice, error) {
	invoice, err := s.IInvoiceRepository.FindById(companyID, id)
	if err != nil {
		return entity.Invoice{}, err
	}

	return invoice, nil
}

func (s *InvoiceService) Update(companyID uuid.UUID, req request.InvoiceUpdateRequest) (entity.Invoice, error) {
	invoice, err := s.IInvoiceRepository.FindById(companyID, req.Id)
	if err != nil {
		return entity.Invoice{}, errors.New("invoice not found")
	}

	if invoice.InvoiceStatus != entity.InvoiceStatusDraft {
		return entity.Invoice{}, errors.New("only draft invoices can be updated")
	}

	_, err = s.ICustomerRepository.FindById(companyID, req.CustomerId)
	if err != nil {
		return entity.Invoice{}, errors.New("customer not found")
	}

	var items []entity.InvoiceItem
	for _, ir := range req.Items {
		items = append(items, entity.InvoiceItem{
			InvoiceId:     req.Id,
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

	if err := s.IInvoiceRepository.Update(&invoice); err != nil {
		return entity.Invoice{}, err
	}
	updatedInvoice, _ := s.IInvoiceRepository.FindById(companyID, invoice.Id)
	return updatedInvoice, nil
}

func (s *InvoiceService) Delete(companyID uuid.UUID, id uuid.UUID) error {
	return s.IInvoiceRepository.Delete(companyID, id)
}

func (s *InvoiceService) Destroy(companyID uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return errors.New("no ids provided")
	}
	return s.IInvoiceRepository.BulkDestroy(companyID, ids)
}

func (s *InvoiceService) Confirm(companyID, id, userId uuid.UUID) error {
	invoice, err := s.IInvoiceRepository.FindById(companyID, id)
	if err != nil {
		return errors.New("invoice not found")
	}
	if invoice.InvoiceStatus != entity.InvoiceStatusDraft {
		return errors.New("only draft invoices can be confirmed")
	}

	config, err := s.ICompanyConfigurationRepository.GetByCompanyId(companyID)
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

	jeRes, err := s.IJournalEntryService.Create(companyID, userId, jeReq)
	if err != nil {
		return fmt.Errorf("failed to generate journal entry: %v", err)
	}

	// Automatically post the journal entry
	if err := s.IJournalEntryService.Post(companyID, jeRes.Id); err != nil {
		return fmt.Errorf("failed to post journal entry: %v", err)
	}

	return s.IInvoiceRepository.UpdateStatus(companyID, id, entity.InvoiceStatusConfirmed, &jeRes.Id)
}

func (s *InvoiceService) RecordPayment(companyID, id, userId uuid.UUID, req request.InvoicePaymentRequest) error {
	invoice, err := s.IInvoiceRepository.FindById(companyID, id)
	if err != nil {
		return errors.New("invoice not found")
	}
	if invoice.InvoiceStatus != entity.InvoiceStatusConfirmed {
		return errors.New("invoice must be confirmed before recording payment")
	}

	config, err := s.ICompanyConfigurationRepository.GetByCompanyId(companyID)
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

	jeRes, err := s.IJournalEntryService.Create(companyID, userId, jeReq)
	if err != nil {
		return fmt.Errorf("failed to generate payment journal entry: %v", err)
	}

	if err := s.IJournalEntryService.Post(companyID, jeRes.Id); err != nil {
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

	if err := s.IInvoiceRepository.AddPayment(&payment); err != nil {
		return err
	}

	newAmountPaid := invoice.AmountPaid + req.Amount
	var pStatus entity.PaymentStatus = entity.PaymentStatusPartial
	if newAmountPaid >= invoice.GrandTotal {
		pStatus = entity.PaymentStatusPaid
	}

	return s.IInvoiceRepository.UpdatePaymentStatus(companyID, id, newAmountPaid, pStatus)
}

func generateInvoiceNumber(date time.Time) string {
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
