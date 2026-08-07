package entity

import (
	"time"

	"github.com/google/uuid"
)

type BillStatus string

const (
	BillStatusDraft     BillStatus = "draft"
	BillStatusConfirmed BillStatus = "confirmed"
	BillStatusCancelled BillStatus = "cancelled"
)

type PaymentStatus string

const (
	PaymentStatusUnpaid  PaymentStatus = "unpaid"
	PaymentStatusPartial PaymentStatus = "partial"
	PaymentStatusPaid    PaymentStatus = "paid"
)

type Bill struct {
	Id              uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	CompanyId       uuid.UUID     `gorm:"type:uuid;not null;index;" json:"company_id"`
	BillNumber      string        `gorm:"type:varchar;not null;uniqueIndex:idx_bill_number_company;" json:"bill_number"`
	PurchaseOrderId *uuid.UUID    `gorm:"type:uuid;" json:"purchase_order_id"`
	VendorId        uuid.UUID     `gorm:"type:uuid;not null;" json:"vendor_id"`
	Vendor          *Vendor       `gorm:"foreignKey:VendorId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"vendor,omitempty"`
	BillDate        time.Time     `gorm:"type:date;not null;" json:"bill_date"`
	DueDate         time.Time     `gorm:"type:date;not null;" json:"due_date"`
	Subtotal        float64       `gorm:"type:numeric(20,4);not null;default:0;" json:"subtotal"`
	DiscountTotal   float64       `gorm:"type:numeric(20,4);not null;default:0;" json:"discount_total"`
	Dpp             float64       `gorm:"type:numeric(20,4);not null;default:0;" json:"dpp"`
	TaxRate         float64       `gorm:"type:numeric(5,4);not null;default:0;" json:"tax_rate"`
	TaxAmount       float64       `gorm:"type:numeric(20,4);not null;default:0;" json:"tax_amount"`
	GrandTotal      float64       `gorm:"type:numeric(20,4);not null;default:0;" json:"grand_total"`
	AmountPaid      float64       `gorm:"type:numeric(20,4);not null;default:0;" json:"amount_paid"`
	AmountDue       float64       `gorm:"type:numeric(20,4);not null;default:0;" json:"amount_due"`
	BillStatus      BillStatus    `gorm:"type:varchar;not null;default:'draft';" json:"bill_status"`
	PaymentStatus   PaymentStatus `gorm:"type:varchar;not null;default:'unpaid';" json:"payment_status"`
	Notes           string        `gorm:"type:text;" json:"notes"`
	JournalEntryId  *uuid.UUID    `gorm:"type:uuid;" json:"journal_entry_id"`
	CreatedBy       uuid.UUID     `gorm:"type:uuid;not null;" json:"created_by"`
	CreatedAt       time.Time     `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt       time.Time     `gorm:"autoUpdateTime;" json:"updated_at"`

	Items    []BillItem    `gorm:"foreignKey:BillId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items"`
	Payments []BillPayment `gorm:"foreignKey:BillId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"payments"`
}

func (Bill) TableName() string { return "bills" }
