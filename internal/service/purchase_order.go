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

type IPurchaseOrderService interface {
	Create(companyId, userId uuid.UUID, req request.PurchaseOrderCreateRequest) (response.PurchaseOrderResponse, error)
	FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]response.PurchaseOrderResponse, int, error)
	FindById(companyId, id uuid.UUID) (response.PurchaseOrderResponse, error)
	Update(companyId, id uuid.UUID, req request.PurchaseOrderUpdateRequest) (response.PurchaseOrderResponse, error)
	Delete(companyId, id uuid.UUID) error
	Destroy(companyId uuid.UUID, ids []uuid.UUID) error
	Approve(companyId, id uuid.UUID) error
}

type PurchaseOrderService struct {
	repo       repository.IPurchaseOrderRepository
	vendorRepo repository.IVendorRepository
}

func NewPurchaseOrderService(repo repository.IPurchaseOrderRepository, vendorRepo repository.IVendorRepository) IPurchaseOrderService {
	return &PurchaseOrderService{repo: repo, vendorRepo: vendorRepo}
}

func generatePoNumber(companyId uuid.UUID, date time.Time) string {
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

func mapPurchaseOrderToResponse(po entity.PurchaseOrder) response.PurchaseOrderResponse {
	respItems := make([]response.PurchaseOrderItemResponse, 0, len(po.Items))
	for _, item := range po.Items {
		respItems = append(respItems, response.PurchaseOrderItemResponse{
			Id:            item.Id,
			Description:   item.Description,
			Qty:           item.Qty,
			Price:         item.Price,
			Discount:      item.Discount,
			TaxApplicable: item.TaxApplicable,
			Amount:        item.Amount,
		})
	}

	resp := response.PurchaseOrderResponse{
		Id:              po.Id,
		CompanyId:       po.CompanyId,
		PoNumber:        po.PoNumber,
		VendorId:        po.VendorId,
		PoDate:          po.PoDate,
		ExpiryDate:      po.ExpiryDate,
		Subtotal:        po.Subtotal,
		DiscountTotal:   po.DiscountTotal,
		Dpp:             po.Dpp,
		TaxRate:         po.TaxRate,
		TaxAmount:       po.TaxAmount,
		GrandTotal:      po.GrandTotal,
		Status:          string(po.Status),
		Notes:           po.Notes,
		ConvertedBillId: po.ConvertedBillId,
		CreatedBy:       po.CreatedBy,
		CreatedAt:       po.CreatedAt,
		UpdatedAt:       po.UpdatedAt,
		Items:           respItems,
	}

	if po.Vendor != nil {
		vendorResp := response.VendorResponse{
			Id:        po.Vendor.Id,
			CompanyId: po.Vendor.CompanyId,
			Code:      po.Vendor.Code,
			Name:      po.Vendor.Name,
			Email:     po.Vendor.Email,
			Phone:     po.Vendor.Phone,
			Address:   po.Vendor.Address,
			CoaId:     po.Vendor.CoaId,
			Status:    po.Vendor.Status,
			CreatedAt: po.Vendor.CreatedAt,
			UpdatedAt: po.Vendor.UpdatedAt,
		}
		if po.Vendor.Coa != nil {
			vendorResp.Coa = &response.COAResponse{
				Id:            po.Vendor.Coa.Id,
				Code:          po.Vendor.Coa.Code,
				Name:          po.Vendor.Coa.Name,
				IsContra:      po.Vendor.Coa.IsContra,
				NormalBalance: po.Vendor.Coa.GetAbsoluteNormalBalance(),
			}
		}
		resp.Vendor = &vendorResp
	}

	return resp
}

func (s *PurchaseOrderService) Create(companyId, userId uuid.UUID, req request.PurchaseOrderCreateRequest) (response.PurchaseOrderResponse, error) {
	_, err := s.vendorRepo.FindById(companyId, req.VendorId)
	if err != nil {
		return response.PurchaseOrderResponse{}, errors.New("vendor not found")
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
		CompanyId:     companyId,
		PoNumber:      generatePoNumber(companyId, req.PoDate),
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

	if err := s.repo.Create(&po); err != nil {
		return response.PurchaseOrderResponse{}, err
	}
	createdPo, _ := s.repo.FindById(companyId, po.Id)
	return mapPurchaseOrderToResponse(createdPo), nil
}

func (s *PurchaseOrderService) FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]response.PurchaseOrderResponse, int, error) {
	pos, total, err := s.repo.FindAll(companyId, qp)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []response.PurchaseOrderResponse{}, 0, nil
	}

	totalPages := (int(total) + qp.PageSize - 1) / qp.PageSize
	if qp.Page > totalPages {
		return nil, int(total), nil
	}

	resps := make([]response.PurchaseOrderResponse, 0, len(pos))
	for _, p := range pos {
		resps = append(resps, mapPurchaseOrderToResponse(p))
	}

	return resps, int(total), nil
}

func (s *PurchaseOrderService) FindById(companyId, id uuid.UUID) (response.PurchaseOrderResponse, error) {
	po, err := s.repo.FindById(companyId, id)
	if err != nil {
		return response.PurchaseOrderResponse{}, errors.New("purchase order not found")
	}
	return mapPurchaseOrderToResponse(po), nil
}

func (s *PurchaseOrderService) Update(companyId, id uuid.UUID, req request.PurchaseOrderUpdateRequest) (response.PurchaseOrderResponse, error) {
	po, err := s.repo.FindById(companyId, id)
	if err != nil {
		return response.PurchaseOrderResponse{}, errors.New("purchase order not found")
	}

	if po.Status != entity.POStatusDraft {
		return response.PurchaseOrderResponse{}, errors.New("only draft purchase orders can be updated")
	}

	_, err = s.vendorRepo.FindById(companyId, req.VendorId)
	if err != nil {
		return response.PurchaseOrderResponse{}, errors.New("vendor not found")
	}

	var items []entity.PurchaseOrderItem
	for _, ir := range req.Items {
		items = append(items, entity.PurchaseOrderItem{
			PurchaseOrderId: id,
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

	if err := s.repo.Update(&po); err != nil {
		return response.PurchaseOrderResponse{}, err
	}
	updatedPo, _ := s.repo.FindById(companyId, po.Id)
	return mapPurchaseOrderToResponse(updatedPo), nil
}

func (s *PurchaseOrderService) Delete(companyId, id uuid.UUID) error {
	po, err := s.repo.FindById(companyId, id)
	if err != nil {
		return errors.New("purchase order not found")
	}
	if po.Status != entity.POStatusDraft {
		return errors.New("only draft purchase orders can be deleted")
	}
	return s.repo.Delete(companyId, id)
}

func (s *PurchaseOrderService) Destroy(companyId uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return errors.New("no ids provided")
	}
	return s.repo.BulkDestroy(companyId, ids)
}

func (s *PurchaseOrderService) Approve(companyId, id uuid.UUID) error {
	po, err := s.repo.FindById(companyId, id)
	if err != nil {
		return errors.New("purchase order not found")
	}
	if po.Status != entity.POStatusDraft {
		return errors.New("only draft purchase orders can be approved")
	}
	return s.repo.UpdateStatus(companyId, id, entity.POStatusApproved)
}
