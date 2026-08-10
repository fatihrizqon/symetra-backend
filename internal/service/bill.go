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

type IBillService interface {
	Create(companyID uuid.UUID, authorID uuid.UUID, req request.BillCreateRequest) (entity.Bill, error)
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.Bill, int, error)
	FindById(companyID uuid.UUID, id uuid.UUID) (entity.Bill, error)
	Update(companyID uuid.UUID, req request.BillUpdateRequest) (entity.Bill, error)
	Delete(companyID uuid.UUID, id uuid.UUID) error
	Destroy(companyID uuid.UUID, ids []uuid.UUID) error
	Confirm(companyID uuid.UUID, id uuid.UUID, authorID uuid.UUID) error
	RecordPayment(companyID uuid.UUID, id uuid.UUID, authorID uuid.UUID, req request.BillPaymentRequest) error
}

type BillService struct {
	validate                        *validator.Validate
	IBillRepository                 repository.IBillRepository
	IVendorRepository               repository.IVendorRepository
	ICompanyConfigurationRepository repository.ICompanyConfigurationRepository
	IJournalEntryService            IJournalEntryService
}

func NewBillService(validate *validator.Validate, IBillRepository repository.IBillRepository, IVendorRepository repository.IVendorRepository, ICompanyConfigurationRepository repository.ICompanyConfigurationRepository, IJournalEntryService IJournalEntryService) IBillService {
	return &BillService{
		validate:                        validate,
		IBillRepository:                 IBillRepository,
		IVendorRepository:               IVendorRepository,
		ICompanyConfigurationRepository: ICompanyConfigurationRepository,
		IJournalEntryService:            IJournalEntryService,
	}
}

func (s *BillService) Create(companyID uuid.UUID, authorID uuid.UUID, req request.BillCreateRequest) (entity.Bill, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.Bill{}, err
	}

	billItems := make([]entity.BillItem, 0, len(req.Items))
	for _, item := range req.Items {
		billItems = append(billItems, entity.BillItem{
			Description:   item.Description,
			Qty:           item.Qty,
			Price:         item.Price,
			Discount:      item.Discount,
			TaxApplicable: item.TaxApplicable,
			AccountId:     item.AccountId,
		})
	}

	subtotal, discountTotal, dpp, taxAmount, grandTotal := s.calculateTotals(billItems, req.TaxRate)

	bill := entity.Bill{
		CompanyId:       companyID,
		BillNumber:      generateBillNumber(req.BillDate),
		PurchaseOrderId: req.PurchaseOrderId,
		VendorId:        req.VendorId,
		BillDate:        req.BillDate,
		DueDate:         req.DueDate,
		Subtotal:        subtotal,
		DiscountTotal:   discountTotal,
		Dpp:             dpp,
		TaxRate:         req.TaxRate,
		TaxAmount:       taxAmount,
		GrandTotal:      grandTotal,
		AmountPaid:      0,
		AmountDue:       grandTotal,
		BillStatus:      entity.BillStatusDraft,
		PaymentStatus:   entity.PaymentStatusUnpaid,
		Notes:           req.Notes,
		CreatedBy:       authorID,
		Items:           billItems,
	}

	err := s.IBillRepository.WithTransaction(func(txRepo repository.IBillRepository) error {
		return txRepo.Create(&bill)
	})
	if err != nil {
		return entity.Bill{}, err
	}

	return bill, nil
}

func (s *BillService) FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.Bill, int, error) {
	bills, totalCount, err := s.IBillRepository.FindAll(companyID, qp)
	if err != nil {
		return nil, 0, err
	}

	return bills, totalCount, nil
}

func (s *BillService) FindById(companyID uuid.UUID, id uuid.UUID) (entity.Bill, error) {
	bill, err := s.IBillRepository.FindById(companyID, id)
	if err != nil {
		return entity.Bill{}, err
	}

	return bill, nil
}

func (s *BillService) Update(companyID uuid.UUID, req request.BillUpdateRequest) (entity.Bill, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.Bill{}, err
	}

	var bill entity.Bill
	err := s.IBillRepository.WithTransaction(func(txRepo repository.IBillRepository) error {
		var txErr error
		bill, txErr = txRepo.FindById(companyID, req.Id)
		if txErr != nil {
			return errors.New("bill not found")
		}

		if bill.BillStatus != entity.BillStatusDraft {
			return errors.New("only draft bills can be updated")
		}

		var billItems []entity.BillItem
		for _, ir := range req.Items {
			billItems = append(billItems, entity.BillItem{
				BillId:        req.Id,
				Description:   ir.Description,
				Qty:           ir.Qty,
				Price:         ir.Price,
				Discount:      ir.Discount,
				TaxApplicable: ir.TaxApplicable,
				AccountId:     ir.AccountId,
			})
		}

		subtotal, discountTotal, dpp, taxAmount, grandTotal := s.calculateTotals(billItems, req.TaxRate)

		bill.VendorId = req.VendorId
		bill.PurchaseOrderId = req.PurchaseOrderId
		bill.BillDate = req.BillDate
		bill.DueDate = req.DueDate
		bill.TaxRate = req.TaxRate
		bill.Notes = req.Notes
		bill.Items = billItems
		bill.Subtotal = subtotal
		bill.DiscountTotal = discountTotal
		bill.Dpp = dpp
		bill.TaxAmount = taxAmount
		bill.GrandTotal = grandTotal
		bill.AmountDue = grandTotal

		if txErr := txRepo.Update(&bill); txErr != nil {
			return txErr
		}
		return nil
	})

	if err != nil {
		return entity.Bill{}, err
	}

	updatedBill, err := s.IBillRepository.FindById(companyID, bill.Id)
	if err != nil {
		return entity.Bill{}, err
	}

	billItems := make([]entity.BillItem, 0, len(updatedBill.Items))
	for _, item := range updatedBill.Items {
		billItems = append(billItems, entity.BillItem{
			Id:            item.Id,
			Description:   item.Description,
			Qty:           item.Qty,
			Price:         item.Price,
			Discount:      item.Discount,
			TaxApplicable: item.TaxApplicable,
			Amount:        item.Amount,
			AccountId:     item.AccountId,
		})
	}

	result := entity.Bill{
		Id:              updatedBill.Id,
		CompanyId:       updatedBill.CompanyId,
		BillNumber:      updatedBill.BillNumber,
		VendorId:        updatedBill.VendorId,
		PurchaseOrderId: updatedBill.PurchaseOrderId,
		BillDate:        updatedBill.BillDate,
		DueDate:         updatedBill.DueDate,
		Subtotal:        updatedBill.Subtotal,
		DiscountTotal:   updatedBill.DiscountTotal,
		Dpp:             updatedBill.Dpp,
		TaxRate:         updatedBill.TaxRate,
		TaxAmount:       updatedBill.TaxAmount,
		GrandTotal:      updatedBill.GrandTotal,
		AmountDue:       updatedBill.AmountDue,
		AmountPaid:      updatedBill.AmountPaid,
		BillStatus:      updatedBill.BillStatus,
		PaymentStatus:   updatedBill.PaymentStatus,
		Notes:           updatedBill.Notes,
		JournalEntryId:  updatedBill.JournalEntryId,
		CreatedBy:       updatedBill.CreatedBy,
		CreatedAt:       updatedBill.CreatedAt,
		UpdatedAt:       updatedBill.UpdatedAt,
		Items:           billItems,
	}

	if updatedBill.Vendor != nil {
		vendor := entity.Vendor{
			Id:        updatedBill.Vendor.Id,
			CompanyId: updatedBill.Vendor.CompanyId,
			Code:      updatedBill.Vendor.Code,
			Name:      updatedBill.Vendor.Name,
			Email:     updatedBill.Vendor.Email,
			Phone:     updatedBill.Vendor.Phone,
			Address:   updatedBill.Vendor.Address,
			CoaId:     updatedBill.Vendor.CoaId,
			Status:    updatedBill.Vendor.Status,
			CreatedAt: updatedBill.Vendor.CreatedAt,
			UpdatedAt: updatedBill.Vendor.UpdatedAt,
		}
		if updatedBill.Vendor.Coa != nil {
			vendor.Coa = &entity.COA{
				Id:       updatedBill.Vendor.Coa.Id,
				Code:     updatedBill.Vendor.Coa.Code,
				Name:     updatedBill.Vendor.Coa.Name,
				IsContra: updatedBill.Vendor.Coa.IsContra,
			}
		}
		result.Vendor = &vendor
	}

	return result, nil
}

func (s *BillService) Delete(companyID uuid.UUID, id uuid.UUID) error {
	return s.IBillRepository.Delete(companyID, id)
}

func (s *BillService) Destroy(companyID uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return errors.New("no ids provided")
	}
	return s.IBillRepository.BulkDestroy(companyID, ids)
}

func (s *BillService) Confirm(companyID uuid.UUID, id uuid.UUID, authorID uuid.UUID) error {
	bill, err := s.IBillRepository.FindById(companyID, id)
	if err != nil {
		return errors.New("bill not found")
	}
	if bill.BillStatus != entity.BillStatusDraft {
		return errors.New("only draft bills can be confirmed")
	}

	config, err := s.ICompanyConfigurationRepository.GetByCompanyId(companyID)
	if err != nil {
		return errors.New("company configuration not found")
	}

	if config.APAccountId == nil {
		return errors.New("accounts payable (A/P) account is not configured in company settings")
	}

	var jeLines []request.JournalLineRequest

	// Credit AP
	jeLines = append(jeLines, request.JournalLineRequest{
		CoaId:       *config.APAccountId,
		Description: util.StringPtr(fmt.Sprintf("A/P for Bill %s", bill.BillNumber)),
		Credit:      bill.GrandTotal,
	})

	// Debit Expense/Inventory items
	for _, item := range bill.Items {
		if item.AccountId == nil {
			return fmt.Errorf("bill item '%s' is missing an account assignment", item.Description)
		}
		jeLines = append(jeLines, request.JournalLineRequest{
			CoaId:       *item.AccountId,
			Description: util.StringPtr(fmt.Sprintf("Bill %s: %s", bill.BillNumber, item.Description)),
			Debit:       item.Amount,
		})
	}

	// Debit Tax Receivable if tax amount > 0
	if bill.TaxAmount > 0 {
		if config.TaxReceivableAccountId == nil {
			return errors.New("tax receivable account is not configured in company settings")
		}
		jeLines = append(jeLines, request.JournalLineRequest{
			CoaId:       *config.TaxReceivableAccountId,
			Description: util.StringPtr(fmt.Sprintf("Tax for Bill %s", bill.BillNumber)),
			Debit:       bill.TaxAmount,
		})
	}

	jeReq := request.JournalEntryCreateRequest{
		Date:        bill.BillDate,
		Description: fmt.Sprintf("Auto-generated Journal Entry for Bill %s", bill.BillNumber),
		Lines:       jeLines,
	}

	jeRes, err := s.IJournalEntryService.Create(companyID, authorID, jeReq)
	if err != nil {
		return fmt.Errorf("failed to generate journal entry: %v", err)
	}

	// Automatically post the journal entry
	if err := s.IJournalEntryService.Post(companyID, jeRes.Id); err != nil {
		return fmt.Errorf("failed to post journal entry: %v", err)
	}

	return s.IBillRepository.UpdateStatus(companyID, id, entity.BillStatusConfirmed, &jeRes.Id)
}

func (s *BillService) RecordPayment(companyID uuid.UUID, id uuid.UUID, authorID uuid.UUID, req request.BillPaymentRequest) error {
	bill, err := s.IBillRepository.FindById(companyID, id)
	if err != nil {
		return errors.New("bill not found")
	}
	if bill.BillStatus != entity.BillStatusConfirmed {
		return errors.New("bill must be confirmed before recording payment")
	}
	if bill.AmountDue <= 0 {
		return errors.New("bill is already fully paid")
	}
	if req.Amount > bill.AmountDue {
		return errors.New("payment amount cannot exceed amount due")
	}

	config, err := s.ICompanyConfigurationRepository.GetByCompanyId(companyID)
	if err != nil {
		return errors.New("company configuration not found")
	}

	if config.APAccountId == nil {
		return errors.New("accounts payable (A/P) account is not configured in company settings")
	}

	// Journal Entry for Payment: Debit A/P, Credit Payment Account (Bank/Cash)
	var jeLines []request.JournalLineRequest
	jeLines = append(jeLines, request.JournalLineRequest{
		CoaId:       *config.APAccountId,
		Description: util.StringPtr(fmt.Sprintf("Payment for Bill %s", bill.BillNumber)),
		Debit:       req.Amount,
	})
	jeLines = append(jeLines, request.JournalLineRequest{
		CoaId:       req.PaymentAccountId,
		Description: util.StringPtr(fmt.Sprintf("Payment for Bill %s", bill.BillNumber)),
		Credit:      req.Amount,
	})

	jeReq := request.JournalEntryCreateRequest{
		Date:        req.PaymentDate,
		Description: fmt.Sprintf("Payment for Bill %s", bill.BillNumber),
		Lines:       jeLines,
	}

	jeRes, err := s.IJournalEntryService.Create(companyID, authorID, jeReq)
	if err != nil {
		return fmt.Errorf("failed to generate payment journal entry: %v", err)
	}

	if err := s.IJournalEntryService.Post(companyID, jeRes.Id); err != nil {
		return fmt.Errorf("failed to post payment journal entry: %v", err)
	}

	payment := entity.BillPayment{
		BillId:           id,
		Amount:           req.Amount,
		PaymentDate:      req.PaymentDate,
		PaymentAccountId: req.PaymentAccountId,
		JournalEntryId:   &jeRes.Id,
		Notes:            req.Notes,
		CreatedBy:        authorID,
	}

	if err := s.IBillRepository.AddPayment(&payment); err != nil {
		return err
	}

	newAmountPaid := bill.AmountPaid + req.Amount
	var pStatus entity.PaymentStatus = entity.PaymentStatusPartial
	if newAmountPaid >= bill.GrandTotal {
		pStatus = entity.PaymentStatusPaid
	}

	return s.IBillRepository.UpdatePaymentStatus(companyID, id, newAmountPaid, pStatus)
}

func generateBillNumber(date time.Time) string {
	return fmt.Sprintf("BILL-%s-%d", date.Format("20060102"), time.Now().UnixMilli())
}

func (s *BillService) calculateTotals(items []entity.BillItem, taxRate float64) (subtotal, discountTotal, dpp, taxAmount, grandTotal float64) {
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
