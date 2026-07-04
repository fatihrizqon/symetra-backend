package repository

import (
	"time"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// companySortColumns is the whitelist of sortable columns for the companies table.
// Key = public-facing ?sort= value, Value = safe SQL column expression.
var companySortColumns = map[string]string{
	"name":       "companies.name",
	"legal_name": "companies.legal_name",
	"email":      "companies.email",
	"status":     "companies.status",
	"created_at": "companies.created_at",
	"updated_at": "companies.updated_at",
}

type ICompanyRepository interface {
	WithTransaction(fn func(txRepo ICompanyRepository) error) error
	Create(entity.Company) (entity.Company, error)
	FindAll(qp *util.QueryParams) ([]entity.Company, int, error)
	FindById(id uuid.UUID) (entity.Company, error)
	Update(entity.Company) error
	Delete(id uuid.UUID) error
}

type CompanyRepository struct {
	Db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) ICompanyRepository {
	return &CompanyRepository{Db: db}
}

func (r *CompanyRepository) WithTransaction(fn func(txRepo ICompanyRepository) error) error {
	tx := r.Db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	txRepo := &CompanyRepository{Db: tx}
	if err := fn(txRepo); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *CompanyRepository) Create(u entity.Company) (entity.Company, error) {
	if err := r.Db.Create(&u).Error; err != nil {
		return u, err
	}
	return u, nil
}

func (r *CompanyRepository) FindAll(qp *util.QueryParams) ([]entity.Company, int, error) {
	var entities []entity.Company
	var totalCount int64

	query := r.Db.Model(&entity.Company{}).Where("deleted_at IS NULL")
	query = util.ApplySearch(query, qp)
	query = entity.Company{}.ApplyFilters(query, qp.Filters)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}

	query = util.ApplySort(query, qp, companySortColumns, "companies.created_at")
	query = util.ApplyPagination(query, qp)

	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, int(totalCount), nil
}

func (r *CompanyRepository) FindById(id uuid.UUID) (entity.Company, error) {
	var u entity.Company
	if err := r.Db.Where("id = ?", id).First(&u).Error; err != nil {
		return u, err
	}
	return u, nil
}

func (r *CompanyRepository) Update(u entity.Company) error {
	if err := r.Db.Model(&u).Updates(u).Error; err != nil {
		return err
	}
	return nil
}

func (r *CompanyRepository) Delete(id uuid.UUID) error {
	if err := r.Db.Model(&entity.Company{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}
	return nil
}
