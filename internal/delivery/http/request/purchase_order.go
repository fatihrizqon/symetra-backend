package request

import (
	"time"

	"github.com/google/uuid"
)

type PurchaseOrderItemRequest struct {
	Description   string  `json:"description" validate:"required"`
	Qty           float64 `json:"qty" validate:"required,min=1"`
	Price         float64 `json:"price" validate:"required,min=0"`
	Discount      float64 `json:"discount" validate:"min=0"`
	TaxApplicable bool    `json:"tax_applicable"`
}

type PurchaseOrderCreateRequest struct {
	VendorId   uuid.UUID                  `json:"vendor_id" validate:"required"`
	PoDate     time.Time                  `json:"po_date" validate:"required"`
	ExpiryDate *time.Time                 `json:"expiry_date"`
	TaxRate    float64                    `json:"tax_rate" validate:"min=0"`
	Notes      string                     `json:"notes"`
	Items      []PurchaseOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type PurchaseOrderUpdateRequest struct {
	Id         uuid.UUID                  `json:"-"`
	VendorId   uuid.UUID                  `json:"vendor_id" validate:"required"`
	PoDate     time.Time                  `json:"po_date" validate:"required"`
	ExpiryDate *time.Time                 `json:"expiry_date"`
	TaxRate    float64                    `json:"tax_rate" validate:"min=0"`
	Notes      string                     `json:"notes"`
	Items      []PurchaseOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}
