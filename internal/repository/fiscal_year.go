package repository

import (
	"errors"
	"fmt"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IFiscalYearRepository interface {
	Create(entity.FiscalYear) (entity.FiscalYear, error)
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.FiscalYear, int, error)
	FindById(companyID, id uuid.UUID) (entity.FiscalYear, error)
	Update(entity.FiscalYear) error
	Delete(companyID, id uuid.UUID) error
	GetActive(companyID uuid.UUID) (entity.FiscalYear, error)
	FindPeriods(companyID, fiscalYearId uuid.UUID, status string) ([]entity.FiscalPeriod, error)
	FindPeriodById(companyID, id uuid.UUID) (entity.FiscalPeriod, error)
	UpdatePeriod(entity.FiscalPeriod) error
	LogPeriodAction(entity.FiscalPeriodLog) error
}

type FiscalYearRepository struct {
	Db *gorm.DB
}

func NewFiscalYearRepository(Db *gorm.DB) IFiscalYearRepository {
	return &FiscalYearRepository{Db: Db}
}

func (r *FiscalYearRepository) Create(ent entity.FiscalYear) (entity.FiscalYear, error) {
	err := r.Db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&ent).Error
	})
	return ent, err
}

func (r *FiscalYearRepository) FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.FiscalYear, int, error) {
	var entities []entity.FiscalYear
	var totalCount int64

	query := r.Db.Model(&entity.FiscalYear{}).Preload("Periods").Where("company_id = ?", companyID)
	// TODO: Apply filters if necessary

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}

	query = query.Order("start_date DESC")
	query = util.ApplyPagination(query, qp)

	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}
	return entities, int(totalCount), nil
}

func (r *FiscalYearRepository) FindById(companyID, id uuid.UUID) (entity.FiscalYear, error) {
	var ent entity.FiscalYear
	if err := r.Db.Preload("Periods", func(db *gorm.DB) *gorm.DB {
		return db.Order("period_number ASC")
	}).Where("id = ? AND company_id = ?", id, companyID).First(&ent).Error; err != nil {
		return ent, errors.New("fiscal year not found")
	}
	return ent, nil
}

func (r *FiscalYearRepository) Update(ent entity.FiscalYear) error {
	return r.Db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&ent).Updates(ent).Error
	})
}

func (r *FiscalYearRepository) Delete(companyID, id uuid.UUID) error {
	return r.Db.Transaction(func(tx *gorm.DB) error {
		return tx.Where("id = ? AND company_id = ?", id, companyID).Delete(&entity.FiscalYear{}).Error
	})
}

func (r *FiscalYearRepository) GetActive(companyID uuid.UUID) (entity.FiscalYear, error) {
	var ent entity.FiscalYear
	if err := r.Db.Where("company_id = ? AND status = ?", companyID, entity.FiscalYearStatusActive).First(&ent).Error; err != nil {
		return ent, fmt.Errorf("active fiscal year not found")
	}
	return ent, nil
}

func (r *FiscalYearRepository) FindPeriods(companyID, fiscalYearId uuid.UUID, status string) ([]entity.FiscalPeriod, error) {
	var entities []entity.FiscalPeriod
	query := r.Db.Joins("JOIN fiscal_years ON fiscal_years.id = fiscal_periods.fiscal_year_id").
		Where("fiscal_years.company_id = ? AND fiscal_periods.fiscal_year_id = ?", companyID, fiscalYearId)
	
	if status != "" {
		query = query.Where("fiscal_periods.status = ?", status)
	}

	if err := query.Order("fiscal_periods.period_number ASC").Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *FiscalYearRepository) FindPeriodById(companyID, id uuid.UUID) (entity.FiscalPeriod, error) {
	var ent entity.FiscalPeriod
	if err := r.Db.Joins("JOIN fiscal_years ON fiscal_years.id = fiscal_periods.fiscal_year_id").
		Where("fiscal_periods.id = ? AND fiscal_years.company_id = ?", id, companyID).First(&ent).Error; err != nil {
		return ent, errors.New("fiscal period not found")
	}
	return ent, nil
}

func (r *FiscalYearRepository) UpdatePeriod(ent entity.FiscalPeriod) error {
	return r.Db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&ent).Updates(ent).Error
	})
}

func (r *FiscalYearRepository) LogPeriodAction(log entity.FiscalPeriodLog) error {
	return r.Db.Create(&log).Error
}
