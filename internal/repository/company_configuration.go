package repository

import (
	"errors"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ICompanyConfigurationRepository interface {
	GetByCompanyId(companyID uuid.UUID) (entity.CompanyConfiguration, error)
	Upsert(companyID uuid.UUID, config entity.CompanyConfiguration) (entity.CompanyConfiguration, error)
}

type CompanyConfigurationRepository struct {
	Db *gorm.DB
}

func NewCompanyConfigurationRepository(Db *gorm.DB) ICompanyConfigurationRepository {
	return &CompanyConfigurationRepository{Db: Db}
}

func (r *CompanyConfigurationRepository) GetByCompanyId(companyID uuid.UUID) (entity.CompanyConfiguration, error) {
	var config entity.CompanyConfiguration
	err := r.Db.Where("company_id = ?", companyID).First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Auto-create empty config if not found
			config = entity.CompanyConfiguration{
				CompanyId:      companyID,
				EnableTax:      false,
				TaxRate:        0.11,
				InvoicePrefix:  "INV",
				QuotationPrefix: "QUO",
				InvoiceDueDays: 30,
			}
			err = r.Db.Create(&config).Error
			return config, err
		}
		return config, err
	}
	return config, nil
}

func (r *CompanyConfigurationRepository) Upsert(companyID uuid.UUID, config entity.CompanyConfiguration) (entity.CompanyConfiguration, error) {
	var existing entity.CompanyConfiguration
	err := r.Db.Where("company_id = ?", companyID).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create
			config.CompanyId = companyID
			err = r.Db.Create(&config).Error
			return config, err
		}
		return existing, err
	}
	// Update existing
	config.Id = existing.Id
	config.CompanyId = companyID
	config.CreatedAt = existing.CreatedAt
	err = r.Db.Save(&config).Error
	return config, err
}
