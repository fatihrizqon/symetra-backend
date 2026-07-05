package repository

import (
	"errors"
	"fmt"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var coaGroupSortColumns = map[string]string{
	"code":           "coa_groups.code",
	"name":           "coa_groups.name",
	"normal_balance": "coa_groups.normal_balance",
	"status":         "coa_groups.status",
	"created_at":     "coa_groups.created_at",
	"updated_at":     "coa_groups.updated_at",
}

type ICOAGroupRepository interface {
	Create(entity.COAGroup) (entity.COAGroup, error)
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.COAGroup, int, error)
	FindById(companyID, entityId uuid.UUID) (entity.COAGroup, error)
	FindByName(companyID uuid.UUID, name string) (entity.COAGroup, error)
	Update(entity.COAGroup) error
	Delete(companyID, entityId uuid.UUID) error
	SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]entity.COAGroup, int, error)
}

type COAGroupRepository struct {
	Db *gorm.DB
}

func NewCOAGroupRepository(Db *gorm.DB) ICOAGroupRepository {
	return &COAGroupRepository{Db: Db}
}

func (r *COAGroupRepository) Create(g entity.COAGroup) (entity.COAGroup, error) {
	tx := r.Db.Begin()
	if err := tx.Create(&g).Error; err != nil {
		tx.Rollback()
		return g, err
	}
	tx.Commit()
	return g, nil
}

func (r *COAGroupRepository) FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.COAGroup, int, error) {
	var entities []entity.COAGroup
	var totalCount int64

	query := r.Db.Model(&entity.COAGroup{}).Where("company_id = ?", companyID)
	query = util.ApplySearch(query, qp)
	query = entity.COAGroup{}.ApplyFilters(query, qp.Filters)

	query = util.ApplySort(query, qp, coaSortColumns, "coa_groups.code")
	query = util.ApplyPagination(query, qp)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}

	query = util.ApplySort(query, qp, coaGroupSortColumns, "coa_groups.created_at")
	query = util.ApplyPagination(query, qp)

	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}
	return entities, int(totalCount), nil
}

func (r *COAGroupRepository) FindById(companyID, entityId uuid.UUID) (entity.COAGroup, error) {
	var g entity.COAGroup
	if err := r.Db.Where("id = ? AND company_id = ?", entityId, companyID).First(&g).Error; err != nil {
		return g, errors.New("coa group not found")
	}
	return g, nil
}

func (r *COAGroupRepository) Update(g entity.COAGroup) error {
	tx := r.Db.Begin()
	if err := tx.Model(&g).Updates(g).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *COAGroupRepository) Delete(companyID, entityId uuid.UUID) error {
	tx := r.Db.Begin()
	if err := tx.Where("id = ? AND company_id = ?", entityId, companyID).Delete(&entity.COAGroup{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *COAGroupRepository) SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]entity.COAGroup, int, error) {
	var entities []entity.COAGroup
	var totalCount int64

	query := r.Db.Model(&entity.COAGroup{}).Where("company_id = ? AND status = 1", companyID)
	query = util.ApplySearch(query, qp)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}

	query = util.ApplySort(query, qp, coaGroupSortColumns, "coa_groups.code")
	query = util.ApplyPagination(query, qp)

	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}
	return entities, int(totalCount), nil
}

func (r *COAGroupRepository) FindByName(companyID uuid.UUID, name string) (entity.COAGroup, error) {
	var g entity.COAGroup
	// Case-insensitive contains match, consistent with report service behavior
	if err := r.Db.Where("company_id = ? AND LOWER(name) = LOWER(?)", companyID, name).First(&g).Error; err != nil {
		return g, fmt.Errorf("coa group '%s' not found", name)
	}
	return g, nil
}
