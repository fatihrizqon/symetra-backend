package request

import (
	"time"

	"github.com/google/uuid"
)

type BillItemRequest struct {
	Description   string     `json:"description" validate:"required"`
	Qty           float64    `json:"qty" validate:"required,min=1"`
	Price         float64    `json:"price" validate:"required,min=0"`
	Discount      float64    `json:"discount" validate:"min=0"`
	TaxApplicable bool       `json:"tax_applicable"`
	AccountId     *uuid.UUID `json:"account_id"`
}

type BillCreateRequest struct {
	VendorId        uuid.UUID         `json:"vendor_id" validate:"required"`
	PurchaseOrderId *uuid.UUID        `json:"purchase_order_id"`
	BillDate        time.Time         `json:"bill_date" validate:"required"`
	DueDate         time.Time         `json:"due_date" validate:"required"`
	TaxRate         float64           `json:"tax_rate" validate:"min=0"`
	Notes           string            `json:"notes"`
	Items           []BillItemRequest `json:"items" validate:"required,min=1,dive"`
}

type BillUpdateRequest struct {
	Id              uuid.UUID         `json:"-"`
	VendorId        uuid.UUID         `json:"vendor_id" validate:"required"`
	PurchaseOrderId *uuid.UUID        `json:"purchase_order_id"`
	BillDate        time.Time         `json:"bill_date" validate:"required"`
	DueDate         time.Time         `json:"due_date" validate:"required"`
	TaxRate         float64           `json:"tax_rate" validate:"min=0"`
	Notes           string            `json:"notes"`
	Items           []BillItemRequest `json:"items" validate:"required,min=1,dive"`
}

type BillPaymentRequest struct {
	Amount           float64   `json:"amount" validate:"required,min=0.01"`
	PaymentDate      time.Time `json:"payment_date" validate:"required"`
	PaymentAccountId uuid.UUID `json:"payment_account_id" validate:"required"`
	Notes            string    `json:"notes"`
}
