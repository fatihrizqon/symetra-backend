package entity

import (
	"time"

	"github.com/google/uuid"
)

func (InvoiceItem) TableName() string { return "invoice_items" }

type InvoiceItem struct {
	Id            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	InvoiceId     uuid.UUID `gorm:"type:uuid;not null;index;" json:"invoice_id"`
	Description   string    `gorm:"type:varchar;not null;" json:"description"`
	Qty           float64   `gorm:"type:numeric(20,4);not null;default:1;" json:"qty"`
	Price         float64   `gorm:"type:numeric(20,4);not null;default:0;" json:"price"`
	Discount      float64   `gorm:"type:numeric(20,4);not null;default:0;" json:"discount"`
	TaxApplicable bool      `gorm:"type:boolean;not null;default:true;" json:"tax_applicable"`
	Amount        float64   `gorm:"type:numeric(20,4);not null;default:0;" json:"amount"`
	CreatedAt     time.Time `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime;" json:"updated_at"`
}
