package repository

import (
	"fmt"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IQuotationRepository interface {
	Create(q *entity.Quotation) error
	FindById(companyId, id uuid.UUID) (entity.Quotation, error)
	FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]entity.Quotation, int64, error)
	Update(q *entity.Quotation) error
	Delete(companyId, id uuid.UUID) error
	BulkDestroy(companyId uuid.UUID, ids []uuid.UUID) error
	UpdateStatus(companyId, id uuid.UUID, status entity.QuotationStatus) error
}

type QuotationRepository struct {
	db *gorm.DB
}

func NewQuotationRepository(db *gorm.DB) IQuotationRepository {
	return &QuotationRepository{db: db}
}

func (r *QuotationRepository) Create(q *entity.Quotation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(q).Error
	})
}

func (r *QuotationRepository) FindById(companyId, id uuid.UUID) (entity.Quotation, error) {
	var q entity.Quotation
	err := r.db.Preload("Customer").Preload("Items").
		Where("id = ? AND company_id = ?", id, companyId).
		First(&q).Error
	return q, err
}

var quotationSortColumns = map[string]string{
	"quotation_number": "quotations.quotation_number",
	"quotation_date":   "quotations.quotation_date",
	"status":           "quotations.status",
	"created_at":       "quotations.created_at",
	"updated_at":       "quotations.updated_at",
}

func (r *QuotationRepository) FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]entity.Quotation, int64, error) {
	var qs []entity.Quotation
	var total int64
	query := r.db.Model(&entity.Quotation{}).Where("company_id = ?", companyId)

	if qp.Search != "" {
		searchLike := "%" + qp.Search + "%"
		query = query.Where("quotation_number ILIKE ? OR notes ILIKE ?", searchLike, searchLike)
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
		if col, valid := quotationSortColumns[qp.SortBy]; valid {
			sort = col
		}
	}
	if qp.SortDir == "asc" || qp.SortDir == "desc" {
		order = qp.SortDir
	}
	query = query.Order(fmt.Sprintf("%s %s", sort, order))

	query = util.ApplyPagination(query, qp)

	err := query.Preload("Customer").Find(&qs).Error
	return qs, total, err
}

func (r *QuotationRepository) Update(q *entity.Quotation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("quotation_id = ?", q.Id).Delete(&entity.QuotationItem{}).Error; err != nil {
			return err
		}
		return tx.Save(q).Error
	})
}

func (r *QuotationRepository) Delete(companyId, id uuid.UUID) error {
	return r.db.Where("id = ? AND company_id = ? AND status = ?", id, companyId, entity.QuotationStatusDraft).Delete(&entity.Quotation{}).Error
}

func (r *QuotationRepository) BulkDestroy(companyId uuid.UUID, ids []uuid.UUID) error {
	return r.db.Where("company_id = ? AND id IN ? AND status = ?", companyId, ids, entity.QuotationStatusDraft).Delete(&entity.Quotation{}).Error
}

func (r *QuotationRepository) UpdateStatus(companyId, id uuid.UUID, status entity.QuotationStatus) error {
	return r.db.Model(&entity.Quotation{}).
		Where("id = ? AND company_id = ?", id, companyId).
		Update("status", status).Error
}
