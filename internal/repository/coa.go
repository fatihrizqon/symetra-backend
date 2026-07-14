package repository

import (
	"errors"
	"fmt"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var coaSortColumns = map[string]string{
	"code":       "coa.code",
	"name":       "coa.name",
	"status":     "coa.status",
	"created_at": "coa.created_at",
	"updated_at": "coa.updated_at",
}

type ICOARepository interface {
	Create(entity.COA) (entity.COA, error)
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.COA, int, error)
	FindById(companyID, entityId uuid.UUID) (entity.COA, error)
	FindByCode(companyID uuid.UUID, code string) (entity.COA, error)
	Update(entity.COA) error
	Delete(companyID, entityId uuid.UUID) error
	SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]entity.COA, int, error)
}

type COARepository struct {
	Db *gorm.DB
}

func NewCOARepository(Db *gorm.DB) ICOARepository {
	return &COARepository{Db: Db}
}

func (r *COARepository) Create(entity entity.COA) (entity.COA, error) {
	err := r.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&entity).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return entity, err
	}
	if err := r.Db.Preload("SubGroup").Preload("SubGroup.Group").First(&entity, "id = ?", entity.Id).Error; err != nil {
		return entity, err
	}
	return entity, nil
}

func (r *COARepository) FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.COA, int, error) {
	var entities []entity.COA
	var totalCount int64

	query := r.Db.Model(&entity.COA{}).Where("coa.company_id = ?", companyID)
	query = util.ApplySearch(query, qp)
	query = entity.COA{}.ApplyFilters(query, qp.Filters)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}

	query = util.ApplySort(query, qp, coaSortColumns, "coa.code")
	query = util.ApplyPagination(query, qp)

	if err := query.Preload("SubGroup").Preload("SubGroup.Group").Find(&entities).Error; err != nil {
		return nil, 0, err
	}
	return entities, int(totalCount), nil
}

func (r *COARepository) FindById(companyID, entityId uuid.UUID) (entity.COA, error) {
	var entity entity.COA
	if err := r.Db.Preload("SubGroup").Preload("SubGroup.Group").Where("id = ? AND company_id = ?", entityId, companyID).First(&entity).Error; err != nil {
		return entity, errors.New("coa not found")
	}
	return entity, nil
}

func (r *COARepository) Update(entity entity.COA) error {
	return r.Db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&entity).Updates(entity).Error
	})
}

func (r *COARepository) Delete(companyID, entityId uuid.UUID) error {
	return r.Db.Transaction(func(tx *gorm.DB) error {
		return tx.Where("id = ? AND company_id = ?", entityId, companyID).Delete(&entity.COA{}).Error
	})
}

func (r *COARepository) SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]entity.COA, int, error) {
	var entities []entity.COA
	var totalCount int64

	query := r.Db.Model(&entity.COA{}).Where("coa.company_id = ? AND coa.status = 1 AND coa.active = true", companyID)
	query = util.ApplySearch(query, qp)
	query = entity.COA{}.ApplyFilters(query, qp.Filters)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}

	query = util.ApplySort(query, qp, coaSortColumns, "coa.code")
	query = util.ApplyPagination(query, qp)

	if err := query.Preload("SubGroup").Find(&entities).Error; err != nil {
		return nil, 0, err
	}
	return entities, int(totalCount), nil
}

func (r *COARepository) FindByCode(companyID uuid.UUID, code string) (entity.COA, error) {
	var entity entity.COA
	if err := r.Db.Where("company_id = ? AND code = ?", companyID, code).First(&entity).Error; err != nil {
		return entity, fmt.Errorf("coa with code %s not found", code)
	}
	return entity, nil
}
