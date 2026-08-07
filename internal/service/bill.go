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

type IBillService interface {
	Create(companyId, userId uuid.UUID, req request.BillCreateRequest) (response.BillResponse, error)
	FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]response.BillResponse, int, error)
	FindById(companyId, id uuid.UUID) (response.BillResponse, error)
	Update(companyId, id uuid.UUID, req request.BillUpdateRequest) (response.BillResponse, error)
	Delete(companyId, id uuid.UUID) error
	Destroy(companyId uuid.UUID, ids []uuid.UUID) error
	Confirm(companyId, id, userId uuid.UUID) error
	RecordPayment(companyId, id, userId uuid.UUID, req request.BillPaymentRequest) error
}

type BillService struct {
	repo           repository.IBillRepository
	vendorRepo     repository.IVendorRepository
	configRepo     repository.ICompanyConfigurationRepository
	journalService IJournalEntryService
}

func NewBillService(
	repo repository.IBillRepository,
	vendorRepo repository.IVendorRepository,
	configRepo repository.ICompanyConfigurationRepository,
	journalService IJournalEntryService,
) IBillService {
	return &BillService{repo: repo, vendorRepo: vendorRepo, configRepo: configRepo, journalService: journalService}
}

func generateBillNumber(companyId uuid.UUID, date time.Time) string {
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

func mapBillToResponse(b entity.Bill) response.BillResponse {
	respItems := make([]response.BillItemResponse, 0, len(b.Items))
	for _, item := range b.Items {
		respItems = append(respItems, response.BillItemResponse{
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

	resp := response.BillResponse{
		Id:              b.Id,
		CompanyId:       b.CompanyId,
		BillNumber:      b.BillNumber,
		VendorId:        b.VendorId,
		PurchaseOrderId: b.PurchaseOrderId,
		BillDate:        b.BillDate,
		DueDate:         b.DueDate,
		Subtotal:        b.Subtotal,
		DiscountTotal:   b.DiscountTotal,
		Dpp:             b.Dpp,
		TaxRate:         b.TaxRate,
		TaxAmount:       b.TaxAmount,
		GrandTotal:      b.GrandTotal,
		AmountDue:       b.AmountDue,
		AmountPaid:      b.AmountPaid,
		BillStatus:      string(b.BillStatus),
		PaymentStatus:   string(b.PaymentStatus),
		Notes:           b.Notes,
		JournalEntryId:  b.JournalEntryId,
		CreatedBy:       b.CreatedBy,
		CreatedAt:       b.CreatedAt,
		UpdatedAt:       b.UpdatedAt,
		Items:           respItems,
	}

	if b.Vendor != nil {
		vendorResp := response.VendorResponse{
			Id:        b.Vendor.Id,
			CompanyId: b.Vendor.CompanyId,
			Code:      b.Vendor.Code,
			Name:      b.Vendor.Name,
			Email:     b.Vendor.Email,
			Phone:     b.Vendor.Phone,
			Address:   b.Vendor.Address,
			CoaId:     b.Vendor.CoaId,
			Status:    b.Vendor.Status,
			CreatedAt: b.Vendor.CreatedAt,
			UpdatedAt: b.Vendor.UpdatedAt,
		}
		if b.Vendor.Coa != nil {
			vendorResp.Coa = &response.COAResponse{
				Id:            b.Vendor.Coa.Id,
				Code:          b.Vendor.Coa.Code,
				Name:          b.Vendor.Coa.Name,
				IsContra:      b.Vendor.Coa.IsContra,
				NormalBalance: b.Vendor.Coa.GetAbsoluteNormalBalance(),
			}
		}
		resp.Vendor = &vendorResp
	}

	return resp
}

func (s *BillService) Create(companyId, userId uuid.UUID, req request.BillCreateRequest) (response.BillResponse, error) {
	_, err := s.vendorRepo.FindById(companyId, req.VendorId)
	if err != nil {
		return response.BillResponse{}, errors.New("vendor not found")
	}

	var items []entity.BillItem
	for _, ir := range req.Items {
		items = append(items, entity.BillItem{
			Description:   ir.Description,
			Qty:           ir.Qty,
			Price:         ir.Price,
			Discount:      ir.Discount,
			TaxApplicable: ir.TaxApplicable,
			AccountId:     ir.AccountId,
		})
	}

	subtotal, discountTotal, dpp, taxAmount, grandTotal := s.calculateTotals(items, req.TaxRate)

	bill := entity.Bill{
		CompanyId:       companyId,
		BillNumber:      generateBillNumber(companyId, req.BillDate),
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
		CreatedBy:       userId,
		Items:           items,
	}

	if err := s.repo.Create(&bill); err != nil {
		return response.BillResponse{}, err
	}
	createdBill, _ := s.repo.FindById(companyId, bill.Id)
	return mapBillToResponse(createdBill), nil
}

func (s *BillService) FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]response.BillResponse, int, error) {
	bills, total, err := s.repo.FindAll(companyId, qp)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []response.BillResponse{}, 0, nil
	}

	totalPages := (int(total) + qp.PageSize - 1) / qp.PageSize
	if qp.Page > totalPages {
		return nil, int(total), nil
	}

	resps := make([]response.BillResponse, 0, len(bills))
	for _, b := range bills {
		resps = append(resps, mapBillToResponse(b))
	}

	return resps, int(total), nil
}

func (s *BillService) FindById(companyId, id uuid.UUID) (response.BillResponse, error) {
	bill, err := s.repo.FindById(companyId, id)
	if err != nil {
		return response.BillResponse{}, errors.New("bill not found")
	}
	return mapBillToResponse(bill), nil
}

func (s *BillService) Update(companyId, id uuid.UUID, req request.BillUpdateRequest) (response.BillResponse, error) {
	bill, err := s.repo.FindById(companyId, id)
	if err != nil {
		return response.BillResponse{}, errors.New("bill not found")
	}

	if bill.BillStatus != entity.BillStatusDraft {
		return response.BillResponse{}, errors.New("only draft bills can be updated")
	}

	_, err = s.vendorRepo.FindById(companyId, req.VendorId)
	if err != nil {
		return response.BillResponse{}, errors.New("vendor not found")
	}

	var items []entity.BillItem
	for _, ir := range req.Items {
		items = append(items, entity.BillItem{
			BillId:        id,
			Description:   ir.Description,
			Qty:           ir.Qty,
			Price:         ir.Price,
			Discount:      ir.Discount,
			TaxApplicable: ir.TaxApplicable,
			AccountId:     ir.AccountId,
		})
	}

	subtotal, discountTotal, dpp, taxAmount, grandTotal := s.calculateTotals(items, req.TaxRate)

	bill.VendorId = req.VendorId
	bill.PurchaseOrderId = req.PurchaseOrderId
	bill.BillDate = req.BillDate
	bill.DueDate = req.DueDate
	bill.TaxRate = req.TaxRate
	bill.Notes = req.Notes
	bill.Items = items
	bill.Subtotal = subtotal
	bill.DiscountTotal = discountTotal
	bill.Dpp = dpp
	bill.TaxAmount = taxAmount
	bill.GrandTotal = grandTotal
	bill.AmountDue = grandTotal

	if err := s.repo.Update(&bill); err != nil {
		return response.BillResponse{}, err
	}
	updatedBill, _ := s.repo.FindById(companyId, bill.Id)
	return mapBillToResponse(updatedBill), nil
}

func (s *BillService) Delete(companyId, id uuid.UUID) error {
	bill, err := s.repo.FindById(companyId, id)
	if err != nil {
		return errors.New("bill not found")
	}
	if bill.BillStatus != entity.BillStatusDraft {
		return errors.New("only draft bills can be deleted")
	}
	return s.repo.Delete(companyId, id)
}

func (s *BillService) Destroy(companyId uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return errors.New("no ids provided")
	}
	return s.repo.BulkDestroy(companyId, ids)
}

func (s *BillService) Confirm(companyId, id, userId uuid.UUID) error {
	bill, err := s.repo.FindById(companyId, id)
	if err != nil {
		return errors.New("bill not found")
	}
	if bill.BillStatus != entity.BillStatusDraft {
		return errors.New("only draft bills can be confirmed")
	}

	config, err := s.configRepo.GetByCompanyId(companyId)
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

	jeRes, err := s.journalService.Create(companyId, userId, jeReq)
	if err != nil {
		return fmt.Errorf("failed to generate journal entry: %v", err)
	}

	// Automatically post the journal entry
	if err := s.journalService.Post(companyId, jeRes.Id); err != nil {
		return fmt.Errorf("failed to post journal entry: %v", err)
	}

	return s.repo.UpdateStatus(companyId, id, entity.BillStatusConfirmed, &jeRes.Id)
}

func (s *BillService) RecordPayment(companyId, id, userId uuid.UUID, req request.BillPaymentRequest) error {
	bill, err := s.repo.FindById(companyId, id)
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

	config, err := s.configRepo.GetByCompanyId(companyId)
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

	jeRes, err := s.journalService.Create(companyId, userId, jeReq)
	if err != nil {
		return fmt.Errorf("failed to generate payment journal entry: %v", err)
	}

	if err := s.journalService.Post(companyId, jeRes.Id); err != nil {
		return fmt.Errorf("failed to post payment journal entry: %v", err)
	}

	payment := entity.BillPayment{
		BillId:           id,
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

	newAmountPaid := bill.AmountPaid + req.Amount
	var pStatus entity.PaymentStatus = entity.PaymentStatusPartial
	if newAmountPaid >= bill.GrandTotal {
		pStatus = entity.PaymentStatusPaid
	}

	return s.repo.UpdatePaymentStatus(companyId, id, newAmountPaid, pStatus)
}
