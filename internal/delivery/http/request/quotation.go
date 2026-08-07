package request

import (
	"time"

	"github.com/google/uuid"
)

type QuotationItemRequest struct {
	Description   string  `json:"description" validate:"required"`
	Qty           float64 `json:"qty" validate:"required,min=1"`
	Price         float64 `json:"price" validate:"required,min=0"`
	Discount      float64 `json:"discount" validate:"min=0"`
	TaxApplicable bool    `json:"tax_applicable"`
}

type QuotationCreateRequest struct {
	CustomerId    uuid.UUID              `json:"customer_id" validate:"required"`
	QuotationDate time.Time              `json:"quotation_date" validate:"required"`
	ExpiryDate    *time.Time             `json:"expiry_date"`
	TaxRate       float64                `json:"tax_rate" validate:"min=0"`
	Notes         string                 `json:"notes"`
	Items         []QuotationItemRequest `json:"items" validate:"required,min=1,dive"`
}

type QuotationUpdateRequest struct {
	Id            uuid.UUID              `json:"-"`
	CustomerId    uuid.UUID              `json:"customer_id" validate:"required"`
	QuotationDate time.Time              `json:"quotation_date" validate:"required"`
	ExpiryDate    *time.Time             `json:"expiry_date"`
	TaxRate       float64                `json:"tax_rate" validate:"min=0"`
	Notes         string                 `json:"notes"`
	Items         []QuotationItemRequest `json:"items" validate:"required,min=1,dive"`
}
