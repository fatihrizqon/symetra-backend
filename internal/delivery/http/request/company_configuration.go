package request

import "github.com/google/uuid"

type CompanyConfigurationUpdateRequest struct {
	EnableTax              bool       `json:"enable_tax"`
	TaxRate                float64    `validate:"min=0,max=1" json:"tax_rate"`
	ARAccountId            *uuid.UUID `json:"ar_account_id"`
	APAccountId            *uuid.UUID `json:"ap_account_id"`
	SalesRevenueAccountId  *uuid.UUID `json:"sales_revenue_account_id"`
	TaxPayableAccountId    *uuid.UUID `json:"tax_payable_account_id"`
	TaxReceivableAccountId *uuid.UUID `json:"tax_receivable_account_id"`
	BankAccountId          *uuid.UUID `json:"bank_account_id"`
	CashAccountId          *uuid.UUID `json:"cash_account_id"`
	RetainedEarningsCOAId  *uuid.UUID `json:"retained_earnings_coa_id"`
	InvoicePrefix          string     `validate:"omitempty,min=1" json:"invoice_prefix"`
	QuotationPrefix        string     `validate:"omitempty,min=1" json:"quotation_prefix"`
	InvoiceDueDays         int        `validate:"min=0" json:"invoice_due_days"`
}
