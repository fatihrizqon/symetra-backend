package repository

import (
	"fmt"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IVendorRepository interface {
	Create(vendor *entity.Vendor) error
	FindById(companyId, id uuid.UUID) (entity.Vendor, error)
	FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]entity.Vendor, int64, error)
	Update(vendor *entity.Vendor) error
	Delete(companyId, id uuid.UUID) error
	BulkDestroy(companyId uuid.UUID, ids []uuid.UUID) error
}

type VendorRepository struct {
	db *gorm.DB
}

func NewVendorRepository(db *gorm.DB) IVendorRepository {
	return &VendorRepository{db: db}
}

func (r *VendorRepository) Create(vendor *entity.Vendor) error {
	return r.db.Create(vendor).Error
}

func (r *VendorRepository) FindById(companyId, id uuid.UUID) (entity.Vendor, error) {
	var vendor entity.Vendor
	err := r.db.Preload("Coa").
		Where("id = ? AND company_id = ?", id, companyId).
		First(&vendor).Error
	return vendor, err
}

var vendorSortColumns = map[string]string{
	"code":       "vendors.code",
	"name":       "vendors.name",
	"created_at": "vendors.created_at",
	"updated_at": "vendors.updated_at",
}

func (r *VendorRepository) FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]entity.Vendor, int64, error) {
	var vendors []entity.Vendor
	var total int64
	query := r.db.Model(&entity.Vendor{}).Where("company_id = ?", companyId)

	if qp.Search != "" {
		searchLike := "%" + qp.Search + "%"
		query = query.Where("code ILIKE ? OR name ILIKE ? OR email ILIKE ?", searchLike, searchLike, searchLike)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sort := "name"
	order := "asc"
	if qp.SortBy != "" {
		if col, valid := vendorSortColumns[qp.SortBy]; valid {
			sort = col
		}
	}
	if qp.SortDir == "asc" || qp.SortDir == "desc" {
		order = qp.SortDir
	}
	query = query.Order(fmt.Sprintf("%s %s", sort, order))

	query = util.ApplyPagination(query, qp)

	err := query.Preload("Coa").Find(&vendors).Error
	return vendors, total, err
}

func (r *VendorRepository) Update(vendor *entity.Vendor) error {
	return r.db.Save(vendor).Error
}

func (r *VendorRepository) Delete(companyId, id uuid.UUID) error {
	return r.db.Where("id = ? AND company_id = ?", id, companyId).Delete(&entity.Vendor{}).Error
}

func (r *VendorRepository) BulkDestroy(companyId uuid.UUID, ids []uuid.UUID) error {
	return r.db.Where("company_id = ? AND id IN ?", companyId, ids).Delete(&entity.Vendor{}).Error
}
