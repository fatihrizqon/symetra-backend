package request

import (
	"time"

	"github.com/google/uuid"
)

type InvoiceItemRequest struct {
	Description   string  `json:"description" validate:"required"`
	Qty           float64 `json:"qty" validate:"required,min=1"`
	Price         float64 `json:"price" validate:"required,min=0"`
	Discount      float64 `json:"discount" validate:"min=0"`
	TaxApplicable bool    `json:"tax_applicable"`
}

type InvoiceCreateRequest struct {
	CustomerId  uuid.UUID            `json:"customer_id" validate:"required"`
	QuotationId *uuid.UUID           `json:"quotation_id"`
	InvoiceDate time.Time            `json:"invoice_date" validate:"required"`
	DueDate     time.Time            `json:"due_date" validate:"required"`
	TaxRate     float64              `json:"tax_rate" validate:"min=0"`
	Notes       string               `json:"notes"`
	Items       []InvoiceItemRequest `json:"items" validate:"required,min=1,dive"`
}

type InvoiceUpdateRequest struct {
	Id          uuid.UUID            `json:"-"`
	CustomerId  uuid.UUID            `json:"customer_id" validate:"required"`
	QuotationId *uuid.UUID           `json:"quotation_id"`
	InvoiceDate time.Time            `json:"invoice_date" validate:"required"`
	DueDate     time.Time            `json:"due_date" validate:"required"`
	TaxRate     float64              `json:"tax_rate" validate:"min=0"`
	Notes       string               `json:"notes"`
	Items       []InvoiceItemRequest `json:"items" validate:"required,min=1,dive"`
}

type InvoicePaymentRequest struct {
	Amount           float64   `json:"amount" validate:"required,min=0.01"`
	PaymentDate      time.Time `json:"payment_date" validate:"required"`
	PaymentAccountId uuid.UUID `json:"payment_account_id" validate:"required"`
	Notes            string    `json:"notes"`
}
