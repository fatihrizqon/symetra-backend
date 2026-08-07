package repository

import (
	"fmt"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ICustomerRepository interface {
	Create(customer *entity.Customer) error
	FindById(companyId, id uuid.UUID) (entity.Customer, error)
	FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]entity.Customer, int64, error)
	Update(customer *entity.Customer) error
	Delete(companyId, id uuid.UUID) error
	BulkDestroy(companyId uuid.UUID, ids []uuid.UUID) error
}

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) ICustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) Create(customer *entity.Customer) error {
	return r.db.Create(customer).Error
}

func (r *CustomerRepository) FindById(companyId, id uuid.UUID) (entity.Customer, error) {
	var customer entity.Customer
	err := r.db.Preload("Coa").
		Where("id = ? AND company_id = ?", id, companyId).
		First(&customer).Error
	return customer, err
}

var customerSortColumns = map[string]string{
	"code":       "customers.code",
	"name":       "customers.name",
	"created_at": "customers.created_at",
	"updated_at": "customers.updated_at",
}

func (r *CustomerRepository) FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]entity.Customer, int64, error) {
	var customers []entity.Customer
	var total int64
	query := r.db.Model(&entity.Customer{}).Where("company_id = ?", companyId)

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
		if col, valid := customerSortColumns[qp.SortBy]; valid {
			sort = col
		}
	}
	if qp.SortDir == "asc" || qp.SortDir == "desc" {
		order = qp.SortDir
	}
	query = query.Order(fmt.Sprintf("%s %s", sort, order))

	query = util.ApplyPagination(query, qp)

	err := query.Preload("Coa").Find(&customers).Error
	return customers, total, err
}

func (r *CustomerRepository) Update(customer *entity.Customer) error {
	return r.db.Save(customer).Error
}

func (r *CustomerRepository) Delete(companyId, id uuid.UUID) error {
	return r.db.Where("id = ? AND company_id = ?", id, companyId).Delete(&entity.Customer{}).Error
}

func (r *CustomerRepository) BulkDestroy(companyId uuid.UUID, ids []uuid.UUID) error {
	return r.db.Where("company_id = ? AND id IN ?", companyId, ids).Delete(&entity.Customer{}).Error
}
