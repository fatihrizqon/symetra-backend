package repository

import (
	"fmt"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IPurchaseOrderRepository interface {
	Create(po *entity.PurchaseOrder) error
	FindById(companyId, id uuid.UUID) (entity.PurchaseOrder, error)
	FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]entity.PurchaseOrder, int64, error)
	Update(po *entity.PurchaseOrder) error
	Delete(companyId, id uuid.UUID) error
	BulkDestroy(companyId uuid.UUID, ids []uuid.UUID) error
	UpdateStatus(companyId, id uuid.UUID, status entity.PurchaseOrderStatus) error
}

type PurchaseOrderRepository struct {
	db *gorm.DB
}

func NewPurchaseOrderRepository(db *gorm.DB) IPurchaseOrderRepository {
	return &PurchaseOrderRepository{db: db}
}

func (r *PurchaseOrderRepository) Create(po *entity.PurchaseOrder) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(po).Error
	})
}

func (r *PurchaseOrderRepository) FindById(companyId, id uuid.UUID) (entity.PurchaseOrder, error) {
	var po entity.PurchaseOrder
	err := r.db.Preload("Vendor").Preload("Items").
		Where("id = ? AND company_id = ?", id, companyId).
		First(&po).Error
	return po, err
}

var poSortColumns = map[string]string{
	"po_number":  "purchase_orders.po_number",
	"po_date":    "purchase_orders.po_date",
	"status":     "purchase_orders.status",
	"created_at": "purchase_orders.created_at",
	"updated_at": "purchase_orders.updated_at",
}

func (r *PurchaseOrderRepository) FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]entity.PurchaseOrder, int64, error) {
	var pos []entity.PurchaseOrder
	var total int64
	query := r.db.Model(&entity.PurchaseOrder{}).Where("company_id = ?", companyId)

	if qp.Search != "" {
		searchLike := "%" + qp.Search + "%"
		query = query.Where("po_number ILIKE ? OR notes ILIKE ?", searchLike, searchLike)
	}

	if statusValues, ok := qp.Filters["status"]; ok && len(statusValues) > 0 && statusValues[0] != "" {
		query = query.Where("status = ?", statusValues[0])
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sort := "created_at"
	order := "desc"
	if qp.SortBy != "" {
		if col, valid := poSortColumns[qp.SortBy]; valid {
			sort = col
		}
	}
	if qp.SortDir == "asc" || qp.SortDir == "desc" {
		order = qp.SortDir
	}
	query = query.Order(fmt.Sprintf("%s %s", sort, order))

	query = util.ApplyPagination(query, qp)

	err := query.Preload("Vendor").Find(&pos).Error
	return pos, total, err
}

func (r *PurchaseOrderRepository) Update(po *entity.PurchaseOrder) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("purchase_order_id = ?", po.Id).Delete(&entity.PurchaseOrderItem{}).Error; err != nil {
			return err
		}
		return tx.Save(po).Error
	})
}

func (r *PurchaseOrderRepository) Delete(companyId, id uuid.UUID) error {
	return r.db.Where("id = ? AND company_id = ? AND status = ?", id, companyId, entity.POStatusDraft).Delete(&entity.PurchaseOrder{}).Error
}

func (r *PurchaseOrderRepository) BulkDestroy(companyId uuid.UUID, ids []uuid.UUID) error {
	return r.db.Where("company_id = ? AND id IN ? AND status = ?", companyId, ids, entity.POStatusDraft).Delete(&entity.PurchaseOrder{}).Error
}

func (r *PurchaseOrderRepository) UpdateStatus(companyId, id uuid.UUID, status entity.PurchaseOrderStatus) error {
	return r.db.Model(&entity.PurchaseOrder{}).
		Where("id = ? AND company_id = ?", id, companyId).
		Update("status", status).Error
}
