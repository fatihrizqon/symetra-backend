package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (BillItem) TableName() string { return "bill_items" }

type BillItem struct {
	Id            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	BillId        uuid.UUID  `gorm:"type:uuid;not null;index;" json:"bill_id"`
	Description   string     `gorm:"type:varchar;not null;" json:"description"`
	Qty           float64    `gorm:"type:numeric(20,4);not null;default:1;" json:"qty"`
	Price         float64    `gorm:"type:numeric(20,4);not null;default:0;" json:"price"`
	Discount      float64    `gorm:"type:numeric(20,4);not null;default:0;" json:"discount"`
	TaxApplicable bool       `gorm:"type:boolean;not null;default:true;" json:"tax_applicable"`
	Amount        float64    `gorm:"type:numeric(20,4);not null;default:0;" json:"amount"`
	AccountId     *uuid.UUID `gorm:"type:uuid;" json:"account_id"`
	Account       *COA       `gorm:"foreignKey:AccountId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"account,omitempty"`
	CreatedAt     time.Time  `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime;" json:"updated_at"`
}

func (BillItem) SearchableFields() []string {
	return []string{"description"}
}

func (BillItem) ApplyFilters(db *gorm.DB, filters map[string][]string) *gorm.DB {
	return db
}
