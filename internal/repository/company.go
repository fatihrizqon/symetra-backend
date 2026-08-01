package repository

import (
	"context"
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
	FindMyCompanies(userID uuid.UUID) ([]entity.CompanyMember, error)
	FindMember(ctx context.Context, companyId uuid.UUID, userId uuid.UUID) (entity.CompanyMember, error)
	FindMembersByCompany(companyID uuid.UUID) ([]entity.CompanyMember, error)
	AssignMember(member entity.CompanyMember) (entity.CompanyMember, error)
	UpdateMemberRole(companyId uuid.UUID, userId uuid.UUID, role string) error
	RemoveMember(companyId uuid.UUID, userId uuid.UUID) error
	Update(entity.Company) error
	Delete(id uuid.UUID) error
	BulkDestroy(ids []uuid.UUID) error
}

type CompanyRepository struct {
	Db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) ICompanyRepository {
	return &CompanyRepository{Db: db}
}

func (r *CompanyRepository) WithTransaction(fn func(txRepo ICompanyRepository) error) error {
	return r.Db.Transaction(func(tx *gorm.DB) error {
		txRepo := &CompanyRepository{Db: tx}
		return fn(txRepo)
	})
}

func (r *CompanyRepository) Create(entity entity.Company) (entity.Company, error) {
	if err := r.Db.Create(&entity).Error; err != nil {
		return entity, err
	}
	return entity, nil
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
	var entity entity.Company
	if err := r.Db.Where("id = ?", id).First(&entity).Error; err != nil {
		return entity, err
	}
	return entity, nil
}

func (r *CompanyRepository) FindMyCompanies(userId uuid.UUID) ([]entity.CompanyMember, error) {
	var members []entity.CompanyMember
	err := r.Db.Model(&entity.CompanyMember{}).
		Preload("Company").
		Where("user_id = ?", userId).
		Find(&members).Error
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (r *CompanyRepository) FindMember(ctx context.Context, companyId uuid.UUID, userId uuid.UUID) (entity.CompanyMember, error) {
	var member entity.CompanyMember
	err := r.Db.WithContext(ctx).Model(&entity.CompanyMember{}).
		Where("company_id = ? AND user_id = ?", companyId, userId).
		First(&member).Error
	return member, err
}

func (r *CompanyRepository) FindMembersByCompany(companyId uuid.UUID) ([]entity.CompanyMember, error) {
	var members []entity.CompanyMember
	err := r.Db.Model(&entity.CompanyMember{}).
		Preload("Company").
		Preload("User").
		Where("company_id = ?", companyId).
		Find(&members).Error
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (r *CompanyRepository) AssignMember(member entity.CompanyMember) (entity.CompanyMember, error) {
	if err := r.Db.Create(&member).Error; err != nil {
		return member, err
	}
	err := r.Db.Preload("Company").Preload("User").First(&member, "id = ?", member.Id).Error
	return member, err
}

func (r *CompanyRepository) UpdateMemberRole(companyId uuid.UUID, userId uuid.UUID, role string) error {
	if err := r.Db.Model(&entity.CompanyMember{}).
		Where("company_id = ?", companyId).
		Where("user_id = ?", userId).
		Update("role", role).Error; err != nil {
		return err
	}
	return nil
}

func (r *CompanyRepository) RemoveMember(companyId uuid.UUID, userId uuid.UUID) error {
	if err := r.Db.Model(&entity.CompanyMember{}).
		Where("company_id = ?", companyId).
		Where("user_id = ?", userId).
		Delete(&entity.CompanyMember{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *CompanyRepository) Update(entity entity.Company) error {
	if err := r.Db.Model(&entity).Updates(entity).Error; err != nil {
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

func (r *CompanyRepository) BulkDestroy(ids []uuid.UUID) error {
	if err := r.Db.Model(&entity.Company{}).Where("id IN ?", ids).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}
	return nil
}

