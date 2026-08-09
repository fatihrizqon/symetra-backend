package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (Customer) TableName() string { return "customers" }

type Customer struct {
	Id        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	CompanyId uuid.UUID  `gorm:"type:uuid;not null;index;" json:"company_id"`
	Code      string     `gorm:"type:varchar;not null;" json:"code"`
	Name      string     `gorm:"type:varchar;not null;" json:"name"`
	Email     string     `gorm:"type:varchar;" json:"email"`
	Phone     string     `gorm:"type:varchar;" json:"phone"`
	Address   string     `gorm:"type:text;" json:"address"`
	CoaId     *uuid.UUID `gorm:"type:uuid;" json:"coa_id"`
	Coa       *COA       `gorm:"foreignKey:CoaId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"coa,omitempty"`
	Status    int        `gorm:"type:int;not null;default:1;" json:"status"`
	CreatedAt time.Time  `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime;" json:"updated_at"`
}

func (Customer) SearchableFields() []string {
	return []string{"code", "name", "email", "phone"}
}

func (Customer) ApplyFilters(db *gorm.DB, filters map[string][]string) *gorm.DB {
	if values, ok := filters["status"]; ok {
		db = db.Where("status IN ?", values)
	}

	return db
}
