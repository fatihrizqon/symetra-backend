package response

import (
	"time"
	"github.com/google/uuid"
)

type QuotationItemResponse struct {
	Id            uuid.UUID `json:"id"`
	Description   string    `json:"description"`
	Qty           float64   `json:"qty"`
	Price         float64   `json:"price"`
	Discount      float64   `json:"discount"`
	TaxApplicable bool      `json:"tax_applicable"`
	Amount        float64   `json:"amount"`
}

type QuotationResponse struct {
	Id                 uuid.UUID               `json:"id"`
	CompanyId          uuid.UUID               `json:"company_id"`
	QuotationNumber    string                  `json:"quotation_number"`
	CustomerId         uuid.UUID               `json:"customer_id"`
	Customer           *CustomerResponse       `json:"customer,omitempty"`
	QuotationDate      time.Time               `json:"quotation_date"`
	ExpiryDate         *time.Time              `json:"expiry_date"`
	Subtotal           float64                 `json:"subtotal"`
	DiscountTotal      float64                 `json:"discount_total"`
	Dpp                float64                 `json:"dpp"`
	TaxRate            float64                 `json:"tax_rate"`
	TaxAmount          float64                 `json:"tax_amount"`
	GrandTotal         float64                 `json:"grand_total"`
	Status             string                  `json:"status"`
	Notes              string                  `json:"notes"`
	ConvertedInvoiceId *uuid.UUID              `json:"converted_invoice_id"`
	CreatedBy          uuid.UUID               `json:"created_by"`
	CreatedAt          time.Time               `json:"created_at"`
	UpdatedAt          time.Time               `json:"updated_at"`
	Items              []QuotationItemResponse `json:"items,omitempty"`
}


