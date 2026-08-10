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

type IPurchaseOrderService interface {
	Create(companyID uuid.UUID, userId uuid.UUID, req request.PurchaseOrderCreateRequest) (entity.PurchaseOrder, error)
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.PurchaseOrder, int, error)
	FindById(companyID uuid.UUID, id uuid.UUID) (entity.PurchaseOrder, error)
	Update(companyID uuid.UUID, req request.PurchaseOrderUpdateRequest) (entity.PurchaseOrder, error)
	Delete(companyID uuid.UUID, id uuid.UUID) error
	Destroy(companyID uuid.UUID, ids []uuid.UUID) error
	Approve(companyID uuid.UUID, id uuid.UUID) error
}

type PurchaseOrderService struct {
	validate                 *validator.Validate
	IPurchaseOrderRepository repository.IPurchaseOrderRepository
	IVendorRepository        repository.IVendorRepository
}

func NewPurchaseOrderService(validate *validator.Validate, IPurchaseOrderRepository repository.IPurchaseOrderRepository, IVendorRepository repository.IVendorRepository) IPurchaseOrderService {
	return &PurchaseOrderService{
		validate:                 validate,
		IPurchaseOrderRepository: IPurchaseOrderRepository,
		IVendorRepository:        IVendorRepository,
	}
}

func generatePoNumber(companyID uuid.UUID, date time.Time) string {
	// simple po number generator, in real app you might want this in repo to prevent race conditions
	return fmt.Sprintf("PO-%s-%d", date.Format("20060102"), time.Now().UnixMilli())
}

func (s *PurchaseOrderService) calculateTotals(items []entity.PurchaseOrderItem, taxRate float64) (subtotal, discountTotal, dpp, taxAmount, grandTotal float64) {
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

func (s *PurchaseOrderService) Create(companyID uuid.UUID, userId uuid.UUID, req request.PurchaseOrderCreateRequest) (entity.PurchaseOrder, error) {
	_, err := s.IVendorRepository.FindById(companyID, req.VendorId)
	if err != nil {
		return entity.PurchaseOrder{}, errors.New("vendor not found")
	}

	var items []entity.PurchaseOrderItem
	for _, ir := range req.Items {
		items = append(items, entity.PurchaseOrderItem{
			Description:   ir.Description,
			Qty:           ir.Qty,
			Price:         ir.Price,
			Discount:      ir.Discount,
			TaxApplicable: ir.TaxApplicable,
		})
	}

	subtotal, discountTotal, dpp, taxAmount, grandTotal := s.calculateTotals(items, req.TaxRate)

	po := entity.PurchaseOrder{
		CompanyId:     companyID,
		PoNumber:      generatePoNumber(companyID, req.PoDate),
		VendorId:      req.VendorId,
		PoDate:        req.PoDate,
		ExpiryDate:    req.ExpiryDate,
		Subtotal:      subtotal,
		DiscountTotal: discountTotal,
		Dpp:           dpp,
		TaxRate:       req.TaxRate,
		TaxAmount:     taxAmount,
		GrandTotal:    grandTotal,
		Status:        entity.POStatusDraft,
		Notes:         req.Notes,
		CreatedBy:     userId,
		Items:         items,
	}

	if err := s.IPurchaseOrderRepository.Create(&po); err != nil {
		return entity.PurchaseOrder{}, err
	}

	createdPo, err := s.IPurchaseOrderRepository.FindById(companyID, po.Id)
	if err != nil {
		return entity.PurchaseOrder{}, err
	}

	poItems := make([]entity.PurchaseOrderItem, 0, len(createdPo.Items))
	for _, item := range createdPo.Items {
		poItems = append(poItems, entity.PurchaseOrderItem{
			Id:            item.Id,
			Description:   item.Description,
			Qty:           item.Qty,
			Price:         item.Price,
			Discount:      item.Discount,
			TaxApplicable: item.TaxApplicable,
			Amount:        item.Amount,
		})
	}

	result := entity.PurchaseOrder{
		Id:              createdPo.Id,
		CompanyId:       createdPo.CompanyId,
		PoNumber:        createdPo.PoNumber,
		VendorId:        createdPo.VendorId,
		PoDate:          createdPo.PoDate,
		ExpiryDate:      createdPo.ExpiryDate,
		Subtotal:        createdPo.Subtotal,
		DiscountTotal:   createdPo.DiscountTotal,
		Dpp:             createdPo.Dpp,
		TaxRate:         createdPo.TaxRate,
		TaxAmount:       createdPo.TaxAmount,
		GrandTotal:      createdPo.GrandTotal,
		Status:          createdPo.Status,
		Notes:           createdPo.Notes,
		ConvertedBillId: createdPo.ConvertedBillId,
		CreatedBy:       createdPo.CreatedBy,
		CreatedAt:       createdPo.CreatedAt,
		UpdatedAt:       createdPo.UpdatedAt,
		Items:           poItems,
	}

	if createdPo.Vendor != nil {
		vendor := entity.Vendor{
			Id:        createdPo.Vendor.Id,
			CompanyId: createdPo.Vendor.CompanyId,
			Code:      createdPo.Vendor.Code,
			Name:      createdPo.Vendor.Name,
			Email:     createdPo.Vendor.Email,
			Phone:     createdPo.Vendor.Phone,
			Address:   createdPo.Vendor.Address,
			CoaId:     createdPo.Vendor.CoaId,
			Status:    createdPo.Vendor.Status,
			CreatedAt: createdPo.Vendor.CreatedAt,
			UpdatedAt: createdPo.Vendor.UpdatedAt,
		}
		if createdPo.Vendor.Coa != nil {
			vendor.Coa = &entity.COA{
				Id:       createdPo.Vendor.Coa.Id,
				Code:     createdPo.Vendor.Coa.Code,
				Name:     createdPo.Vendor.Coa.Name,
				IsContra: createdPo.Vendor.Coa.IsContra,
			}
		}
		result.Vendor = &vendor
	}

	return result, nil
}

func (s *PurchaseOrderService) FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.PurchaseOrder, int, error) {
	purchase_orders, totalCount, err := s.IPurchaseOrderRepository.FindAll(companyID, qp)
	if err != nil {
		return nil, 0, err
	}
	return purchase_orders, totalCount, nil
}

func (s *PurchaseOrderService) FindById(companyID uuid.UUID, id uuid.UUID) (entity.PurchaseOrder, error) {
	purchase_order, err := s.IPurchaseOrderRepository.FindById(companyID, id)
	if err != nil {
		return entity.PurchaseOrder{}, err
	}

	return purchase_order, nil
}

func (s *PurchaseOrderService) Update(companyID uuid.UUID, req request.PurchaseOrderUpdateRequest) (entity.PurchaseOrder, error) {
	po, err := s.IPurchaseOrderRepository.FindById(companyID, req.Id)
	if err != nil {
		return entity.PurchaseOrder{}, errors.New("purchase order not found")
	}

	if po.Status != entity.POStatusDraft {
		return entity.PurchaseOrder{}, errors.New("only draft purchase orders can be updated")
	}

	_, err = s.IVendorRepository.FindById(companyID, req.VendorId)
	if err != nil {
		return entity.PurchaseOrder{}, errors.New("vendor not found")
	}

	var items []entity.PurchaseOrderItem
	for _, ir := range req.Items {
		items = append(items, entity.PurchaseOrderItem{
			PurchaseOrderId: req.Id,
			Description:     ir.Description,
			Qty:             ir.Qty,
			Price:           ir.Price,
			Discount:        ir.Discount,
			TaxApplicable:   ir.TaxApplicable,
		})
	}

	subtotal, discountTotal, dpp, taxAmount, grandTotal := s.calculateTotals(items, req.TaxRate)

	po.VendorId = req.VendorId
	po.PoDate = req.PoDate
	po.ExpiryDate = req.ExpiryDate
	po.TaxRate = req.TaxRate
	po.Notes = req.Notes
	po.Items = items
	po.Subtotal = subtotal
	po.DiscountTotal = discountTotal
	po.Dpp = dpp
	po.TaxAmount = taxAmount
	po.GrandTotal = grandTotal

	if err := s.IPurchaseOrderRepository.Update(&po); err != nil {
		return entity.PurchaseOrder{}, err
	}

	updatedPo, err := s.IPurchaseOrderRepository.FindById(companyID, po.Id)
	if err != nil {
		return entity.PurchaseOrder{}, err
	}

	respItems := make([]entity.PurchaseOrderItem, 0, len(updatedPo.Items))
	for _, item := range updatedPo.Items {
		respItems = append(respItems, entity.PurchaseOrderItem{
			Id:            item.Id,
			Description:   item.Description,
			Qty:           item.Qty,
			Price:         item.Price,
			Discount:      item.Discount,
			TaxApplicable: item.TaxApplicable,
			Amount:        item.Amount,
		})
	}

	resp := entity.PurchaseOrder{
		Id:              updatedPo.Id,
		CompanyId:       updatedPo.CompanyId,
		PoNumber:        updatedPo.PoNumber,
		VendorId:        updatedPo.VendorId,
		PoDate:          updatedPo.PoDate,
		ExpiryDate:      updatedPo.ExpiryDate,
		Subtotal:        updatedPo.Subtotal,
		DiscountTotal:   updatedPo.DiscountTotal,
		Dpp:             updatedPo.Dpp,
		TaxRate:         updatedPo.TaxRate,
		TaxAmount:       updatedPo.TaxAmount,
		GrandTotal:      updatedPo.GrandTotal,
		Status:          updatedPo.Status,
		Notes:           updatedPo.Notes,
		ConvertedBillId: updatedPo.ConvertedBillId,
		CreatedBy:       updatedPo.CreatedBy,
		CreatedAt:       updatedPo.CreatedAt,
		UpdatedAt:       updatedPo.UpdatedAt,
		Items:           respItems,
	}

	if updatedPo.Vendor != nil {
		vendor := entity.Vendor{
			Id:        updatedPo.Vendor.Id,
			CompanyId: updatedPo.Vendor.CompanyId,
			Code:      updatedPo.Vendor.Code,
			Name:      updatedPo.Vendor.Name,
			Email:     updatedPo.Vendor.Email,
			Phone:     updatedPo.Vendor.Phone,
			Address:   updatedPo.Vendor.Address,
			CoaId:     updatedPo.Vendor.CoaId,
			Status:    updatedPo.Vendor.Status,
			CreatedAt: updatedPo.Vendor.CreatedAt,
			UpdatedAt: updatedPo.Vendor.UpdatedAt,
		}
		if updatedPo.Vendor.Coa != nil {
			vendor.Coa = &entity.COA{
				Id:       updatedPo.Vendor.Coa.Id,
				Code:     updatedPo.Vendor.Coa.Code,
				Name:     updatedPo.Vendor.Coa.Name,
				IsContra: updatedPo.Vendor.Coa.IsContra,
			}
		}
		resp.Vendor = &vendor
	}

	return resp, nil
}

func (s *PurchaseOrderService) Delete(companyID uuid.UUID, id uuid.UUID) error {
	po, err := s.IPurchaseOrderRepository.FindById(companyID, id)
	if err != nil {
		return errors.New("purchase order not found")
	}
	if po.Status != entity.POStatusDraft {
		return errors.New("only draft purchase orders can be deleted")
	}
	return s.IPurchaseOrderRepository.Delete(companyID, id)
}

func (s *PurchaseOrderService) Destroy(companyID uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return errors.New("no ids provided")
	}
	return s.IPurchaseOrderRepository.BulkDestroy(companyID, ids)
}

func (s *PurchaseOrderService) Approve(companyID uuid.UUID, id uuid.UUID) error {
	po, err := s.IPurchaseOrderRepository.FindById(companyID, id)
	if err != nil {
		return errors.New("purchase order not found")
	}
	if po.Status != entity.POStatusDraft {
		return errors.New("only draft purchase orders can be approved")
	}
	return s.IPurchaseOrderRepository.UpdateStatus(companyID, id, entity.POStatusApproved)
}
