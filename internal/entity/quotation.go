package entity

import (
	"time"

	"github.com/google/uuid"
)

type QuotationStatus string

const (
	QuotationStatusDraft    QuotationStatus = "draft"
	QuotationStatusSent     QuotationStatus = "sent"
	QuotationStatusAccepted QuotationStatus = "accepted"
	QuotationStatusDeclined QuotationStatus = "declined"
)

type Quotation struct {
	Id                 uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	CompanyId          uuid.UUID       `gorm:"type:uuid;not null;index;" json:"company_id"`
	QuotationNumber    string          `gorm:"type:varchar;not null;uniqueIndex:idx_quotation_number_company;" json:"quotation_number"`
	CustomerId         uuid.UUID       `gorm:"type:uuid;not null;" json:"customer_id"`
	Customer           *Customer       `gorm:"foreignKey:CustomerId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"customer,omitempty"`
	QuotationDate      time.Time       `gorm:"type:date;not null;" json:"quotation_date"`
	ExpiryDate         *time.Time      `gorm:"type:date;" json:"expiry_date"`
	Subtotal           float64         `gorm:"type:numeric(20,4);not null;default:0;" json:"subtotal"`
	DiscountTotal      float64         `gorm:"type:numeric(20,4);not null;default:0;" json:"discount_total"`
	Dpp                float64         `gorm:"type:numeric(20,4);not null;default:0;" json:"dpp"`
	TaxRate            float64         `gorm:"type:numeric(5,4);not null;default:0;" json:"tax_rate"`
	TaxAmount          float64         `gorm:"type:numeric(20,4);not null;default:0;" json:"tax_amount"`
	GrandTotal         float64         `gorm:"type:numeric(20,4);not null;default:0;" json:"grand_total"`
	Status             QuotationStatus `gorm:"type:varchar;not null;default:'draft';" json:"status"`
	Notes              string          `gorm:"type:text;" json:"notes"`
	ConvertedInvoiceId *uuid.UUID      `gorm:"type:uuid;" json:"converted_invoice_id"`
	CreatedBy          uuid.UUID       `gorm:"type:uuid;not null;" json:"created_by"`
	CreatedAt          time.Time       `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt          time.Time       `gorm:"autoUpdateTime;" json:"updated_at"`

	Items []QuotationItem `gorm:"foreignKey:QuotationId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items"`
}

func (Quotation) TableName() string { return "quotations" }
