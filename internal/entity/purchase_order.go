package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (PurchaseOrder) TableName() string { return "purchase_orders" }

type PurchaseOrderStatus string

const (
	POStatusDraft    PurchaseOrderStatus = "draft"
	POStatusSent     PurchaseOrderStatus = "sent"
	POStatusApproved PurchaseOrderStatus = "approved"
	POStatusDeclined PurchaseOrderStatus = "declined"
)

type PurchaseOrder struct {
	Id              uuid.UUID           `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	CompanyId       uuid.UUID           `gorm:"type:uuid;not null;index;" json:"company_id"`
	PoNumber        string              `gorm:"type:varchar;not null;uniqueIndex:idx_po_number_company;" json:"po_number"`
	VendorId        uuid.UUID           `gorm:"type:uuid;not null;" json:"vendor_id"`
	Vendor          *Vendor             `gorm:"foreignKey:VendorId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"vendor,omitempty"`
	PoDate          time.Time           `gorm:"type:date;not null;" json:"po_date"`
	ExpiryDate      *time.Time          `gorm:"type:date;" json:"expiry_date"`
	Subtotal        float64             `gorm:"type:numeric(20,4);not null;default:0;" json:"subtotal"`
	DiscountTotal   float64             `gorm:"type:numeric(20,4);not null;default:0;" json:"discount_total"`
	Dpp             float64             `gorm:"type:numeric(20,4);not null;default:0;" json:"dpp"`
	TaxRate         float64             `gorm:"type:numeric(5,4);not null;default:0;" json:"tax_rate"`
	TaxAmount       float64             `gorm:"type:numeric(20,4);not null;default:0;" json:"tax_amount"`
	GrandTotal      float64             `gorm:"type:numeric(20,4);not null;default:0;" json:"grand_total"`
	Status          PurchaseOrderStatus `gorm:"type:varchar;not null;default:'draft';" json:"status"`
	Notes           string              `gorm:"type:text;" json:"notes"`
	ConvertedBillId *uuid.UUID          `gorm:"type:uuid;" json:"converted_bill_id"`
	CreatedBy       uuid.UUID           `gorm:"type:uuid;not null;" json:"created_by"`
	CreatedAt       time.Time           `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt       time.Time           `gorm:"autoUpdateTime;" json:"updated_at"`

	Items []PurchaseOrderItem `gorm:"foreignKey:PurchaseOrderId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items"`
}

func (PurchaseOrder) SearchableFields() []string {
	return []string{"po_number", "notes"}
}

func (PurchaseOrder) ApplyFilters(db *gorm.DB, filters map[string][]string) *gorm.DB {
	return db
}
