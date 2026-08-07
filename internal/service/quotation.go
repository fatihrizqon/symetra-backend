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

type IQuotationService interface {
	Create(companyId, userId uuid.UUID, req request.QuotationCreateRequest) (response.QuotationResponse, error)
	FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]response.QuotationResponse, int, error)
	FindById(companyId, id uuid.UUID) (response.QuotationResponse, error)
	Update(companyId, id uuid.UUID, req request.QuotationUpdateRequest) (response.QuotationResponse, error)
	Delete(companyId, id uuid.UUID) error
	Destroy(companyId uuid.UUID, ids []uuid.UUID) error
	Approve(companyId, id uuid.UUID) error
}

type QuotationService struct {
	repo         repository.IQuotationRepository
	customerRepo repository.ICustomerRepository
}

func NewQuotationService(repo repository.IQuotationRepository, customerRepo repository.ICustomerRepository) IQuotationService {
	return &QuotationService{repo: repo, customerRepo: customerRepo}
}

func generateQuotationNumber(companyId uuid.UUID, date time.Time) string {
	return fmt.Sprintf("QUO-%s-%d", date.Format("20060102"), time.Now().UnixMilli())
}

func (s *QuotationService) calculateTotals(items []entity.QuotationItem, taxRate float64) (subtotal, discountTotal, dpp, taxAmount, grandTotal float64) {
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

func mapQuotationToResponse(q entity.Quotation) response.QuotationResponse {
	respItems := make([]response.QuotationItemResponse, 0, len(q.Items))
	for _, item := range q.Items {
		respItems = append(respItems, response.QuotationItemResponse{
			Id:            item.Id,
			Description:   item.Description,
			Qty:           item.Qty,
			Price:         item.Price,
			Discount:      item.Discount,
			TaxApplicable: item.TaxApplicable,
			Amount:        item.Amount,
		})
	}

	resp := response.QuotationResponse{
		Id:                 q.Id,
		CompanyId:          q.CompanyId,
		QuotationNumber:    q.QuotationNumber,
		CustomerId:         q.CustomerId,
		QuotationDate:      q.QuotationDate,
		ExpiryDate:         q.ExpiryDate,
		Subtotal:           q.Subtotal,
		DiscountTotal:      q.DiscountTotal,
		Dpp:                q.Dpp,
		TaxRate:            q.TaxRate,
		TaxAmount:          q.TaxAmount,
		GrandTotal:         q.GrandTotal,
		Status:             string(q.Status),
		Notes:              q.Notes,
		ConvertedInvoiceId: q.ConvertedInvoiceId,
		CreatedBy:          q.CreatedBy,
		CreatedAt:          q.CreatedAt,
		UpdatedAt:          q.UpdatedAt,
		Items:              respItems,
	}

	if q.Customer != nil {
		custResp := response.CustomerResponse{
			Id:        q.Customer.Id,
			CompanyId: q.Customer.CompanyId,
			Code:      q.Customer.Code,
			Name:      q.Customer.Name,
			Email:     q.Customer.Email,
			Phone:     q.Customer.Phone,
			Address:   q.Customer.Address,
			CoaId:     q.Customer.CoaId,
			Status:    q.Customer.Status,
			CreatedAt: q.Customer.CreatedAt,
			UpdatedAt: q.Customer.UpdatedAt,
		}
		if q.Customer.Coa != nil {
			custResp.Coa = &response.COAResponse{
				Id:            q.Customer.Coa.Id,
				Code:          q.Customer.Coa.Code,
				Name:          q.Customer.Coa.Name,
				IsContra:      q.Customer.Coa.IsContra,
				NormalBalance: q.Customer.Coa.GetAbsoluteNormalBalance(),
			}
		}
		resp.Customer = &custResp
	}

	return resp
}

func (s *QuotationService) Create(companyId, userId uuid.UUID, req request.QuotationCreateRequest) (response.QuotationResponse, error) {
	_, err := s.customerRepo.FindById(companyId, req.CustomerId)
	if err != nil {
		return response.QuotationResponse{}, errors.New("customer not found")
	}

	var items []entity.QuotationItem
	for _, ir := range req.Items {
		items = append(items, entity.QuotationItem{
			Description:   ir.Description,
			Qty:           ir.Qty,
			Price:         ir.Price,
			Discount:      ir.Discount,
			TaxApplicable: ir.TaxApplicable,
		})
	}

	subtotal, discountTotal, dpp, taxAmount, grandTotal := s.calculateTotals(items, req.TaxRate)

	quo := entity.Quotation{
		CompanyId:       companyId,
		QuotationNumber: generateQuotationNumber(companyId, req.QuotationDate),
		CustomerId:      req.CustomerId,
		QuotationDate:   req.QuotationDate,
		ExpiryDate:      req.ExpiryDate,
		Subtotal:        subtotal,
		DiscountTotal:   discountTotal,
		Dpp:             dpp,
		TaxRate:         req.TaxRate,
		TaxAmount:       taxAmount,
		GrandTotal:      grandTotal,
		Status:          entity.QuotationStatusDraft,
		Notes:           req.Notes,
		CreatedBy:       userId,
		Items:           items,
	}

	if err := s.repo.Create(&quo); err != nil {
		return response.QuotationResponse{}, err
	}
	createdQuo, _ := s.repo.FindById(companyId, quo.Id)
	return mapQuotationToResponse(createdQuo), nil
}

func (s *QuotationService) FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]response.QuotationResponse, int, error) {
	quos, total, err := s.repo.FindAll(companyId, qp)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []response.QuotationResponse{}, 0, nil
	}

	totalPages := (int(total) + qp.PageSize - 1) / qp.PageSize
	if qp.Page > totalPages {
		return nil, int(total), nil
	}

	resps := make([]response.QuotationResponse, 0, len(quos))
	for _, q := range quos {
		resps = append(resps, mapQuotationToResponse(q))
	}

	return resps, int(total), nil
}

func (s *QuotationService) FindById(companyId, id uuid.UUID) (response.QuotationResponse, error) {
	quo, err := s.repo.FindById(companyId, id)
	if err != nil {
		return response.QuotationResponse{}, errors.New("quotation not found")
	}
	return mapQuotationToResponse(quo), nil
}

func (s *QuotationService) Update(companyId, id uuid.UUID, req request.QuotationUpdateRequest) (response.QuotationResponse, error) {
	quo, err := s.repo.FindById(companyId, id)
	if err != nil {
		return response.QuotationResponse{}, errors.New("quotation not found")
	}

	if quo.Status != entity.QuotationStatusDraft {
		return response.QuotationResponse{}, errors.New("only draft quotations can be updated")
	}

	_, err = s.customerRepo.FindById(companyId, req.CustomerId)
	if err != nil {
		return response.QuotationResponse{}, errors.New("customer not found")
	}

	var items []entity.QuotationItem
	for _, ir := range req.Items {
		items = append(items, entity.QuotationItem{
			QuotationId:   id,
			Description:   ir.Description,
			Qty:           ir.Qty,
			Price:         ir.Price,
			Discount:      ir.Discount,
			TaxApplicable: ir.TaxApplicable,
		})
	}

	subtotal, discountTotal, dpp, taxAmount, grandTotal := s.calculateTotals(items, req.TaxRate)

	quo.CustomerId = req.CustomerId
	quo.QuotationDate = req.QuotationDate
	quo.ExpiryDate = req.ExpiryDate
	quo.TaxRate = req.TaxRate
	quo.Notes = req.Notes
	quo.Items = items
	quo.Subtotal = subtotal
	quo.DiscountTotal = discountTotal
	quo.Dpp = dpp
	quo.TaxAmount = taxAmount
	quo.GrandTotal = grandTotal

	if err := s.repo.Update(&quo); err != nil {
		return response.QuotationResponse{}, err
	}
	updatedQuo, _ := s.repo.FindById(companyId, quo.Id)
	return mapQuotationToResponse(updatedQuo), nil
}

func (s *QuotationService) Delete(companyId, id uuid.UUID) error {
	quo, err := s.repo.FindById(companyId, id)
	if err != nil {
		return errors.New("quotation not found")
	}
	if quo.Status != entity.QuotationStatusDraft {
		return errors.New("only draft quotations can be deleted")
	}
	return s.repo.Delete(companyId, id)
}

func (s *QuotationService) Destroy(companyId uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return errors.New("no ids provided")
	}
	return s.repo.BulkDestroy(companyId, ids)
}

func (s *QuotationService) Approve(companyId, id uuid.UUID) error {
	quo, err := s.repo.FindById(companyId, id)
	if err != nil {
		return errors.New("quotation not found")
	}
	if quo.Status != entity.QuotationStatusDraft {
		return errors.New("only draft quotations can be approved")
	}
	return s.repo.UpdateStatus(companyId, id, entity.QuotationStatusAccepted)
}
