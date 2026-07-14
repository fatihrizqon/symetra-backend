package response

import (
	"github.com/google/uuid"
)

type CompanyConfigurationResponse struct {
	Id                     uuid.UUID    `json:"id"`
	CompanyId              uuid.UUID    `json:"company_id"`
	EnableTax              bool         `json:"enable_tax"`
	TaxRate                float64      `json:"tax_rate"`
	ARAccountId            *uuid.UUID   `json:"ar_account_id"`
	ARAccount              *COAResponse `json:"ar_account,omitempty"`
	APAccountId            *uuid.UUID   `json:"ap_account_id"`
	APAccount              *COAResponse `json:"ap_account,omitempty"`
	SalesRevenueAccountId  *uuid.UUID   `json:"sales_revenue_account_id"`
	TaxPayableAccountId    *uuid.UUID   `json:"tax_payable_account_id"`
	TaxReceivableAccountId *uuid.UUID   `json:"tax_receivable_account_id"`
	BankAccountId          *uuid.UUID   `json:"bank_account_id"`
	CashAccountId          *uuid.UUID   `json:"cash_account_id"`
	RetainedEarningsCOAId  *uuid.UUID   `json:"retained_earnings_coa_id"`
	InvoicePrefix          string       `json:"invoice_prefix"`
	QuotationPrefix        string       `json:"quotation_prefix"`
	InvoiceDueDays         int          `json:"invoice_due_days"`
}
