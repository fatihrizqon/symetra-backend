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

type IQuotationService interface {
	Create(companyID uuid.UUID, userId uuid.UUID, req request.QuotationCreateRequest) (entity.Quotation, error)
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.Quotation, int, error)
	FindById(companyID uuid.UUID, id uuid.UUID) (entity.Quotation, error)
	Update(companyID uuid.UUID, req request.QuotationUpdateRequest) (entity.Quotation, error)
	Delete(companyID uuid.UUID, id uuid.UUID) error
	Destroy(companyID uuid.UUID, ids []uuid.UUID) error
	Approve(companyID uuid.UUID, id uuid.UUID) error
}

type QuotationService struct {
	validate             *validator.Validate
	IQuotationRepository repository.IQuotationRepository
	ICustomerRepository  repository.ICustomerRepository
}

func NewQuotationService(validate *validator.Validate, repo repository.IQuotationRepository, customerRepo repository.ICustomerRepository) IQuotationService {
	return &QuotationService{
		validate:             validate,
		IQuotationRepository: repo,
		ICustomerRepository:  customerRepo,
	}
}

func generateQuotationNumber(companyID uuid.UUID, date time.Time) string {
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

func (s *QuotationService) Create(companyID uuid.UUID, userId uuid.UUID, req request.QuotationCreateRequest) (entity.Quotation, error) {
	_, err := s.ICustomerRepository.FindById(companyID, req.CustomerId)
	if err != nil {
		return entity.Quotation{}, errors.New("customer not found")
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
		CompanyId:       companyID,
		QuotationNumber: generateQuotationNumber(companyID, req.QuotationDate),
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

	if err := s.IQuotationRepository.Create(&quo); err != nil {
		return entity.Quotation{}, err
	}

	createdQuo, err := s.IQuotationRepository.FindById(companyID, quo.Id)
	if err != nil {
		return entity.Quotation{}, err
	}

	quotationItems := make([]entity.QuotationItem, 0, len(createdQuo.Items))
	for _, item := range createdQuo.Items {
		quotationItems = append(quotationItems, entity.QuotationItem{
			Id:            item.Id,
			Description:   item.Description,
			Qty:           item.Qty,
			Price:         item.Price,
			Discount:      item.Discount,
			TaxApplicable: item.TaxApplicable,
			Amount:        item.Amount,
		})
	}

	resp := entity.Quotation{
		Id:                 createdQuo.Id,
		CompanyId:          createdQuo.CompanyId,
		QuotationNumber:    createdQuo.QuotationNumber,
		CustomerId:         createdQuo.CustomerId,
		QuotationDate:      createdQuo.QuotationDate,
		ExpiryDate:         createdQuo.ExpiryDate,
		Subtotal:           createdQuo.Subtotal,
		DiscountTotal:      createdQuo.DiscountTotal,
		Dpp:                createdQuo.Dpp,
		TaxRate:            createdQuo.TaxRate,
		TaxAmount:          createdQuo.TaxAmount,
		GrandTotal:         createdQuo.GrandTotal,
		Status:             createdQuo.Status,
		Notes:              createdQuo.Notes,
		ConvertedInvoiceId: createdQuo.ConvertedInvoiceId,
		CreatedBy:          createdQuo.CreatedBy,
		CreatedAt:          createdQuo.CreatedAt,
		UpdatedAt:          createdQuo.UpdatedAt,
		Items:              quotationItems,
	}

	if createdQuo.Customer != nil {
		customer := entity.Customer{
			Id:        createdQuo.Customer.Id,
			CompanyId: createdQuo.Customer.CompanyId,
			Code:      createdQuo.Customer.Code,
			Name:      createdQuo.Customer.Name,
			Email:     createdQuo.Customer.Email,
			Phone:     createdQuo.Customer.Phone,
			Address:   createdQuo.Customer.Address,
			CoaId:     createdQuo.Customer.CoaId,
			Status:    createdQuo.Customer.Status,
			CreatedAt: createdQuo.Customer.CreatedAt,
			UpdatedAt: createdQuo.Customer.UpdatedAt,
		}
		if createdQuo.Customer.Coa != nil {
			customer.Coa = &entity.COA{
				Id:       createdQuo.Customer.Coa.Id,
				Code:     createdQuo.Customer.Coa.Code,
				Name:     createdQuo.Customer.Coa.Name,
				IsContra: createdQuo.Customer.Coa.IsContra,
			}
		}
		resp.Customer = &customer
	}

	return resp, nil
}

func (s *QuotationService) FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.Quotation, int, error) {
	quos, total, err := s.IQuotationRepository.FindAll(companyID, qp)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []entity.Quotation{}, 0, nil
	}

	totalPages := (int(total) + qp.PageSize - 1) / qp.PageSize
	if qp.Page > totalPages {
		return nil, int(total), nil
	}

	resps := make([]entity.Quotation, 0, len(quos))
	for _, q := range quos {
		respItems := make([]entity.QuotationItem, 0, len(q.Items))
		for _, item := range q.Items {
			respItems = append(respItems, entity.QuotationItem{
				Id:            item.Id,
				Description:   item.Description,
				Qty:           item.Qty,
				Price:         item.Price,
				Discount:      item.Discount,
				TaxApplicable: item.TaxApplicable,
				Amount:        item.Amount,
			})
		}

		resp := entity.Quotation{
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
			Status:             q.Status,
			Notes:              q.Notes,
			ConvertedInvoiceId: q.ConvertedInvoiceId,
			CreatedBy:          q.CreatedBy,
			CreatedAt:          q.CreatedAt,
			UpdatedAt:          q.UpdatedAt,
			Items:              respItems,
		}

		if q.Customer != nil {
			customer := entity.Customer{
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
				customer.Coa = &entity.COA{
					Id:       q.Customer.Coa.Id,
					Code:     q.Customer.Coa.Code,
					Name:     q.Customer.Coa.Name,
					IsContra: q.Customer.Coa.IsContra,
				}
			}
			resp.Customer = &customer
		}

		resps = append(resps, resp)
	}

	return resps, int(total), nil
}

func (s *QuotationService) FindById(companyID uuid.UUID, id uuid.UUID) (entity.Quotation, error) {
	q, err := s.IQuotationRepository.FindById(companyID, id)
	if err != nil {
		return entity.Quotation{}, errors.New("quotation not found")
	}

	respItems := make([]entity.QuotationItem, 0, len(q.Items))
	for _, item := range q.Items {
		respItems = append(respItems, entity.QuotationItem{
			Id:            item.Id,
			Description:   item.Description,
			Qty:           item.Qty,
			Price:         item.Price,
			Discount:      item.Discount,
			TaxApplicable: item.TaxApplicable,
			Amount:        item.Amount,
		})
	}

	resp := entity.Quotation{
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
		Status:             q.Status,
		Notes:              q.Notes,
		ConvertedInvoiceId: q.ConvertedInvoiceId,
		CreatedBy:          q.CreatedBy,
		CreatedAt:          q.CreatedAt,
		UpdatedAt:          q.UpdatedAt,
		Items:              respItems,
	}

	if q.Customer != nil {
		customer := entity.Customer{
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
			customer.Coa = &entity.COA{
				Id:       q.Customer.Coa.Id,
				Code:     q.Customer.Coa.Code,
				Name:     q.Customer.Coa.Name,
				IsContra: q.Customer.Coa.IsContra,
			}
		}
		resp.Customer = &customer
	}

	return resp, nil
}

func (s *QuotationService) Update(companyID uuid.UUID, req request.QuotationUpdateRequest) (entity.Quotation, error) {
	quo, err := s.IQuotationRepository.FindById(companyID, req.Id)
	if err != nil {
		return entity.Quotation{}, errors.New("quotation not found")
	}

	if quo.Status != entity.QuotationStatusDraft {
		return entity.Quotation{}, errors.New("only draft quotations can be updated")
	}

	_, err = s.ICustomerRepository.FindById(companyID, req.CustomerId)
	if err != nil {
		return entity.Quotation{}, errors.New("customer not found")
	}

	var items []entity.QuotationItem
	for _, ir := range req.Items {
		items = append(items, entity.QuotationItem{
			QuotationId:   req.Id,
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

	if err := s.IQuotationRepository.Update(&quo); err != nil {
		return entity.Quotation{}, err
	}

	updatedQuo, err := s.IQuotationRepository.FindById(companyID, quo.Id)
	if err != nil {
		return entity.Quotation{}, err
	}

	respItems := make([]entity.QuotationItem, 0, len(updatedQuo.Items))
	for _, item := range updatedQuo.Items {
		respItems = append(respItems, entity.QuotationItem{
			Id:            item.Id,
			Description:   item.Description,
			Qty:           item.Qty,
			Price:         item.Price,
			Discount:      item.Discount,
			TaxApplicable: item.TaxApplicable,
			Amount:        item.Amount,
		})
	}

	resp := entity.Quotation{
		Id:                 updatedQuo.Id,
		CompanyId:          updatedQuo.CompanyId,
		QuotationNumber:    updatedQuo.QuotationNumber,
		CustomerId:         updatedQuo.CustomerId,
		QuotationDate:      updatedQuo.QuotationDate,
		ExpiryDate:         updatedQuo.ExpiryDate,
		Subtotal:           updatedQuo.Subtotal,
		DiscountTotal:      updatedQuo.DiscountTotal,
		Dpp:                updatedQuo.Dpp,
		TaxRate:            updatedQuo.TaxRate,
		TaxAmount:          updatedQuo.TaxAmount,
		GrandTotal:         updatedQuo.GrandTotal,
		Status:             updatedQuo.Status,
		Notes:              updatedQuo.Notes,
		ConvertedInvoiceId: updatedQuo.ConvertedInvoiceId,
		CreatedBy:          updatedQuo.CreatedBy,
		CreatedAt:          updatedQuo.CreatedAt,
		UpdatedAt:          updatedQuo.UpdatedAt,
		Items:              respItems,
	}

	if updatedQuo.Customer != nil {
		customer := entity.Customer{
			Id:        updatedQuo.Customer.Id,
			CompanyId: updatedQuo.Customer.CompanyId,
			Code:      updatedQuo.Customer.Code,
			Name:      updatedQuo.Customer.Name,
			Email:     updatedQuo.Customer.Email,
			Phone:     updatedQuo.Customer.Phone,
			Address:   updatedQuo.Customer.Address,
			CoaId:     updatedQuo.Customer.CoaId,
			Status:    updatedQuo.Customer.Status,
			CreatedAt: updatedQuo.Customer.CreatedAt,
			UpdatedAt: updatedQuo.Customer.UpdatedAt,
		}
		if updatedQuo.Customer.Coa != nil {
			customer.Coa = &entity.COA{
				Id:       updatedQuo.Customer.Coa.Id,
				Code:     updatedQuo.Customer.Coa.Code,
				Name:     updatedQuo.Customer.Coa.Name,
				IsContra: updatedQuo.Customer.Coa.IsContra,
			}
		}
		resp.Customer = &customer
	}

	return resp, nil
}

func (s *QuotationService) Delete(companyID uuid.UUID, id uuid.UUID) error {
	quo, err := s.IQuotationRepository.FindById(companyID, id)
	if err != nil {
		return errors.New("quotation not found")
	}
	if quo.Status != entity.QuotationStatusDraft {
		return errors.New("only draft quotations can be deleted")
	}
	return s.IQuotationRepository.Delete(companyID, id)
}

func (s *QuotationService) Destroy(companyID uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return errors.New("no ids provided")
	}
	return s.IQuotationRepository.BulkDestroy(companyID, ids)
}

func (s *QuotationService) Approve(companyID uuid.UUID, id uuid.UUID) error {
	quo, err := s.IQuotationRepository.FindById(companyID, id)
	if err != nil {
		return errors.New("quotation not found")
	}
	if quo.Status != entity.QuotationStatusDraft {
		return errors.New("only draft quotations can be approved")
	}
	return s.IQuotationRepository.UpdateStatus(companyID, id, entity.QuotationStatusAccepted)
}
