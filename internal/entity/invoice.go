package entity

import (
	"time"

	"github.com/google/uuid"
)

type InvoiceStatus string

const (
	InvoiceStatusDraft     InvoiceStatus = "draft"
	InvoiceStatusConfirmed InvoiceStatus = "confirmed"
	InvoiceStatusCancelled InvoiceStatus = "cancelled"
)

type TaxStatus string

const (
	TaxStatusDraft    TaxStatus = "draft"
	TaxStatusIssued   TaxStatus = "issued"
	TaxStatusReported TaxStatus = "reported"
)

type Invoice struct {
	Id               uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	CompanyId        uuid.UUID     `gorm:"type:uuid;not null;index;" json:"company_id"`
	InvoiceNumber    string        `gorm:"type:varchar;not null;uniqueIndex:idx_invoice_number_company;" json:"invoice_number"`
	QuotationId      *uuid.UUID    `gorm:"type:uuid;" json:"quotation_id"`
	CustomerId       uuid.UUID     `gorm:"type:uuid;not null;" json:"customer_id"`
	Customer         *Customer     `gorm:"foreignKey:CustomerId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"customer,omitempty"`
	InvoiceDate      time.Time     `gorm:"type:date;not null;" json:"invoice_date"`
	DueDate          time.Time     `gorm:"type:date;not null;" json:"due_date"`
	Subtotal         float64       `gorm:"type:numeric(20,4);not null;default:0;" json:"subtotal"`
	DiscountTotal    float64       `gorm:"type:numeric(20,4);not null;default:0;" json:"discount_total"`
	Dpp              float64       `gorm:"type:numeric(20,4);not null;default:0;" json:"dpp"`
	TaxRate          float64       `gorm:"type:numeric(5,4);not null;default:0;" json:"tax_rate"`
	TaxAmount        float64       `gorm:"type:numeric(20,4);not null;default:0;" json:"tax_amount"`
	GrandTotal       float64       `gorm:"type:numeric(20,4);not null;default:0;" json:"grand_total"`
	AmountPaid       float64       `gorm:"type:numeric(20,4);not null;default:0;" json:"amount_paid"`
	AmountDue        float64       `gorm:"type:numeric(20,4);not null;default:0;" json:"amount_due"`
	InvoiceStatus    InvoiceStatus `gorm:"type:varchar;not null;default:'draft';" json:"invoice_status"`
	PaymentStatus    PaymentStatus `gorm:"type:varchar;not null;default:'unpaid';" json:"payment_status"`
	TaxStatus        TaxStatus     `gorm:"type:varchar;not null;default:'draft';" json:"tax_status"`
	Notes            string        `gorm:"type:text;" json:"notes"`
	JournalEntryId   *uuid.UUID    `gorm:"type:uuid;" json:"journal_entry_id"`
	PaymentJournalId *uuid.UUID    `gorm:"type:uuid;" json:"payment_journal_id"`
	CreatedBy        uuid.UUID     `gorm:"type:uuid;not null;" json:"created_by"`
	CreatedAt        time.Time     `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt        time.Time     `gorm:"autoUpdateTime;" json:"updated_at"`

	Items    []InvoiceItem    `gorm:"foreignKey:InvoiceId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items"`
	Payments []InvoicePayment `gorm:"foreignKey:InvoiceId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"payments"`
}

func (Invoice) TableName() string { return "invoices" }
