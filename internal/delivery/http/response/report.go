package response

import (
	"time"

	"github.com/google/uuid"
)

type GeneralLedgerEntry struct {
	Date           time.Time `json:"date"`
	JournalNumber  string    `json:"journal_number"`
	Description    string    `json:"description"`
	Debit          float64   `json:"debit"`
	Credit         float64   `json:"credit"`
	RunningBalance float64   `json:"running_balance"`
}

type GeneralLedgerResponse struct {
	Account struct {
		Id   uuid.UUID `json:"id"`
		Code string    `json:"code"`
		Name string    `json:"name"`
	} `json:"account"`
	OpeningBalance float64              `json:"opening_balance"`
	Entries        []GeneralLedgerEntry `json:"entries"`
	Totals         struct {
		TotalDebit     float64 `json:"total_debit"`
		TotalCredit    float64 `json:"total_credit"`
		ClosingBalance float64 `json:"closing_balance"`
	} `json:"totals"`
}

type TrialBalanceAccount struct {
	AccountType string  `json:"account_type"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	TotalDebit  float64 `json:"total_debit"`
	TotalCredit float64 `json:"total_credit"`
	Balance     float64 `json:"balance"`
}

type TrialBalanceResponse struct {
	Period struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"period"`
	Accounts []TrialBalanceAccount `json:"accounts"`
	Summary  struct {
		GrandTotalDebit  float64 `json:"grand_total_debit"`
		GrandTotalCredit float64 `json:"grand_total_credit"`
		IsBalanced       bool    `json:"is_balanced"`
	} `json:"summary"`
}

type ProfitLossAccount struct {
	Code    string  `json:"code"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

type ProfitLossSection struct {
	Accounts []ProfitLossAccount `json:"accounts"`
	Total    float64             `json:"total"`
}

type ProfitLossResponse struct {
	Period struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"period"`
	Revenue   ProfitLossSection `json:"revenue"`
	Expense   ProfitLossSection `json:"expense"`
	NetIncome float64           `json:"net_income"`
}

type BalanceSheetAccount struct {
	Code    string  `json:"code"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

type BalanceSheetGroup struct {
	Name     string                `json:"name"`
	Accounts []BalanceSheetAccount `json:"accounts"`
	Subtotal float64               `json:"subtotal"`
}

type BalanceSheetSection struct {
	Groups []BalanceSheetGroup `json:"groups"`
	Total  float64             `json:"total"`
}

type BalanceSheetResponse struct {
	AsOf   string              `json:"as_of"`
	Assets BalanceSheetSection `json:"assets"`
	Liabilities BalanceSheetSection `json:"liabilities"`
	Equity      BalanceSheetSection `json:"equity"`
	TotalLiabilitiesAndEquity float64 `json:"total_liabilities_and_equity"`
	IsBalanced                bool    `json:"is_balanced"`
}

type CashFlowMovement struct {
	Date           time.Time `json:"date"`
	JournalNumber  string    `json:"journal_number"`
	Description    string    `json:"description"`
	Inflow         float64   `json:"inflow"`
	Outflow        float64   `json:"outflow"`
	RunningBalance float64   `json:"running_balance"`
}

type CashFlowResponse struct {
	Period struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"period"`
	OpeningBalance float64            `json:"opening_balance"`
	Movements      []CashFlowMovement `json:"movements"`
	Summary        struct {
		TotalInflow    float64 `json:"total_inflow"`
		TotalOutflow   float64 `json:"total_outflow"`
		NetChange      float64 `json:"net_change"`
		ClosingBalance float64 `json:"closing_balance"`
	} `json:"summary"`
}

type EquityStatementAccount struct {
	Code           string  `json:"code"`
	Name           string  `json:"name"`
	OpeningBalance float64 `json:"opening_balance"`
	NetIncome      float64 `json:"net_income"`
	ClosingBalance float64 `json:"closing_balance"`
}

type EquityStatementResponse struct {
	Period struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"period"`
	EquityAccounts []EquityStatementAccount `json:"equity_accounts"`
}

type JournalBookLine struct {
	CoaCode string  `json:"coa_code"`
	CoaName string  `json:"coa_name"`
	Debit   float64 `json:"debit"`
	Credit  float64 `json:"credit"`
}

type JournalBookEntry struct {
	JournalNumber string            `json:"journal_number"`
	Date          time.Time         `json:"date"`
	Type          string            `json:"type"`
	Description   string            `json:"description"`
	Lines         []JournalBookLine `json:"lines"`
	TotalDebit    float64           `json:"total_debit"`
	TotalCredit   float64           `json:"total_credit"`
}

type JournalBookResponse []JournalBookEntry
