package entity

import (
	"time"

	"github.com/google/uuid"
)

func (CompanyConfiguration) TableName() string { return "company_configurations" }

type CompanyConfiguration struct {
	Id                      uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	CompanyId               uuid.UUID  `gorm:"type:uuid;uniqueIndex;not null;" json:"company_id"`
	EnableTax               bool       `gorm:"type:boolean;not null;default:false;" json:"enable_tax"`
	TaxRate                 float64    `gorm:"type:numeric(5,4);not null;default:0.11;" json:"tax_rate"` // e.g. 0.11 for 11%
	ARAccountId             *uuid.UUID `gorm:"type:uuid;" json:"ar_account_id"`
	ARAccount               *COA       `gorm:"foreignKey:ARAccountId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"ar_account,omitempty"`
	APAccountId             *uuid.UUID `gorm:"type:uuid;" json:"ap_account_id"`
	APAccount               *COA       `gorm:"foreignKey:APAccountId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"ap_account,omitempty"`
	SalesRevenueAccountId   *uuid.UUID `gorm:"type:uuid;" json:"sales_revenue_account_id"`
	SalesRevenueAccount     *COA       `gorm:"foreignKey:SalesRevenueAccountId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"sales_revenue_account,omitempty"`
	TaxPayableAccountId     *uuid.UUID `gorm:"type:uuid;" json:"tax_payable_account_id"`
	TaxPayableAccount       *COA       `gorm:"foreignKey:TaxPayableAccountId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"tax_payable_account,omitempty"`
	TaxReceivableAccountId  *uuid.UUID `gorm:"type:uuid;" json:"tax_receivable_account_id"`
	TaxReceivableAccount    *COA       `gorm:"foreignKey:TaxReceivableAccountId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"tax_receivable_account,omitempty"`
	BankAccountId           *uuid.UUID `gorm:"type:uuid;" json:"bank_account_id"`
	BankAccount             *COA       `gorm:"foreignKey:BankAccountId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"bank_account,omitempty"`
	CashAccountId           *uuid.UUID `gorm:"type:uuid;" json:"cash_account_id"`
	CashAccount             *COA       `gorm:"foreignKey:CashAccountId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"cash_account,omitempty"`
	RetainedEarningsCOAId   *uuid.UUID `gorm:"type:uuid;" json:"retained_earnings_coa_id"`
	RetainedEarningsCOA     *COA       `gorm:"foreignKey:RetainedEarningsCOAId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"retained_earnings_coa,omitempty"`
	InvoicePrefix           string     `gorm:"type:character varying;not null;default:'INV';" json:"invoice_prefix"`
	QuotationPrefix         string     `gorm:"type:character varying;not null;default:'QUO';" json:"quotation_prefix"`
	InvoiceDueDays          int        `gorm:"type:int;not null;default:30;" json:"invoice_due_days"`
	CreatedAt               time.Time  `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt               time.Time  `gorm:"autoUpdateTime;" json:"updated_at"`
}
