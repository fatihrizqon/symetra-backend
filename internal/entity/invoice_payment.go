package entity

import (
	"time"

	"github.com/google/uuid"
)

type InvoicePayment struct {
	Id               uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	InvoiceId        uuid.UUID  `gorm:"type:uuid;not null;index;" json:"invoice_id"`
	Amount           float64    `gorm:"type:numeric(20,4);not null;" json:"amount"`
	PaymentDate      time.Time  `gorm:"type:date;not null;" json:"payment_date"`
	PaymentAccountId uuid.UUID  `gorm:"type:uuid;not null;" json:"payment_account_id"`
	PaymentAccount   *COA       `gorm:"foreignKey:PaymentAccountId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"payment_account,omitempty"`
	JournalEntryId   *uuid.UUID `gorm:"type:uuid;" json:"journal_entry_id"`
	Notes            string     `gorm:"type:text;" json:"notes"`
	CreatedBy        uuid.UUID  `gorm:"type:uuid;not null;" json:"created_by"`
	CreatedAt        time.Time  `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime;" json:"updated_at"`
}

func (InvoicePayment) TableName() string { return "invoice_payments" }
