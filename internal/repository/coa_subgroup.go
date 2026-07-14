package repository

import (
	"errors"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var coaSubGroupSortColumns = map[string]string{
	"code":       "coa_subgroups.code",
	"name":       "coa_subgroups.name",
	"status":     "coa_subgroups.status",
	"created_at": "coa_subgroups.created_at",
	"updated_at": "coa_subgroups.updated_at",
}

type ICOASubGroupRepository interface {
	Create(entity.COASubGroup) (entity.COASubGroup, error)
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.COASubGroup, int, error)
	FindById(companyID, entityId uuid.UUID) (entity.COASubGroup, error)
	Update(entity.COASubGroup) error
	Delete(companyID, entityId uuid.UUID) error
	SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]entity.COASubGroup, int, error)
}

type COASubGroupRepository struct {
	Db *gorm.DB
}

func NewCOASubGroupRepository(Db *gorm.DB) ICOASubGroupRepository {
	return &COASubGroupRepository{Db: Db}
}

func (r *COASubGroupRepository) Create(ent entity.COASubGroup) (entity.COASubGroup, error) {
	err := r.Db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&ent).Error
	})
	if err != nil {
		return ent, err
	}
	// Preload AFTER commit — safe to use main connection now.
	if err := r.Db.Preload("Group").First(&ent, "id = ?", ent.Id).Error; err != nil {
		return ent, err
	}
	return ent, nil
}

func (r *COASubGroupRepository) FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.COASubGroup, int, error) {
	var entities []entity.COASubGroup
	var totalCount int64

	query := r.Db.Preload("Group").Model(&entity.COASubGroup{}).Where("coa_subgroups.company_id = ?", companyID)
	query = util.ApplySearch(query, qp)
	query = entity.COASubGroup{}.ApplyFilters(query, qp.Filters)

	query = util.ApplySort(query, qp, coaSortColumns, "coa_subgroups.code")
	query = util.ApplyPagination(query, qp)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}

	query = util.ApplySort(query, qp, coaSubGroupSortColumns, "coa_subgroups.created_at")
	query = util.ApplyPagination(query, qp)

	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}
	return entities, int(totalCount), nil
}

func (r *COASubGroupRepository) FindById(companyID, entityId uuid.UUID) (entity.COASubGroup, error) {
	var ent entity.COASubGroup
	if err := r.Db.Preload("Group").Where("id = ? AND company_id = ?", entityId, companyID).First(&ent).Error; err != nil {
		return ent, errors.New("coa subgroup not found")
	}
	return ent, nil
}

func (r *COASubGroupRepository) Update(ent entity.COASubGroup) error {
	return r.Db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&ent).Updates(ent).Error
	})
}

func (r *COASubGroupRepository) Delete(companyID, entityId uuid.UUID) error {
	return r.Db.Transaction(func(tx *gorm.DB) error {
		return tx.Where("id = ? AND company_id = ?", entityId, companyID).Delete(&entity.COASubGroup{}).Error
	})
}

func (r *COASubGroupRepository) SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]entity.COASubGroup, int, error) {
	var entities []entity.COASubGroup
	var totalCount int64

	query := r.Db.Model(&entity.COASubGroup{}).Where("coa_subgroups.company_id = ? AND coa_subgroups.status = 1", companyID)
	query = util.ApplySearch(query, qp)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}

	query = util.ApplySort(query, qp, coaSubGroupSortColumns, "coa_subgroups.code")
	query = util.ApplyPagination(query, qp)

	if err := query.Preload("Group").Find(&entities).Error; err != nil {
		return nil, 0, err
	}
	return entities, int(totalCount), nil
}
