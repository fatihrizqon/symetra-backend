package response

import (
	"time"
	"github.com/google/uuid"
)

type PurchaseOrderItemResponse struct {
	Id            uuid.UUID `json:"id"`
	Description   string    `json:"description"`
	Qty           float64   `json:"qty"`
	Price         float64   `json:"price"`
	Discount      float64   `json:"discount"`
	TaxApplicable bool      `json:"tax_applicable"`
	Amount        float64   `json:"amount"`
}

type PurchaseOrderResponse struct {
	Id              uuid.UUID                   `json:"id"`
	CompanyId       uuid.UUID                   `json:"company_id"`
	PoNumber        string                      `json:"po_number"`
	VendorId        uuid.UUID                   `json:"vendor_id"`
	Vendor          *VendorResponse             `json:"vendor,omitempty"`
	PoDate          time.Time                   `json:"po_date"`
	ExpiryDate      *time.Time                  `json:"expiry_date"`
	Subtotal        float64                     `json:"subtotal"`
	DiscountTotal   float64                     `json:"discount_total"`
	Dpp             float64                     `json:"dpp"`
	TaxRate         float64                     `json:"tax_rate"`
	TaxAmount       float64                     `json:"tax_amount"`
	GrandTotal      float64                     `json:"grand_total"`
	Status          string                      `json:"status"`
	Notes           string                      `json:"notes"`
	ConvertedBillId *uuid.UUID                  `json:"converted_bill_id"`
	CreatedBy       uuid.UUID                   `json:"created_by"`
	CreatedAt       time.Time                   `json:"created_at"`
	UpdatedAt       time.Time                   `json:"updated_at"`
	Items           []PurchaseOrderItemResponse `json:"items,omitempty"`
}


