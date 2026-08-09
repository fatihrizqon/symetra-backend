package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportGeneralLedgerRow struct {
	Date            time.Time `json:"date"`
	JournalNumber   string    `json:"journal_number"`
	JeDescription   string    `json:"je_description"`
	LineDescription string    `json:"line_description"`
	Debit           float64   `json:"debit"`
	Credit          float64   `json:"credit"`
	AccountCode     string    `json:"account_code"`
	AccountName     string    `json:"account_name"`
}

type ReportTrialBalanceRow struct {
	AccountType string  `json:"account_type"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	TotalDebit  float64 `json:"total_debit"`
	TotalCredit float64 `json:"total_credit"`
	Balance     float64 `json:"balance"`
}

type ReportProfitLossRow struct {
	Type         string  `json:"type"`
	GroupName    string  `json:"group_name"`
	SubgroupName string  `json:"subgroup_name"`
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Balance      float64 `json:"balance"`
}

type ReportBalanceSheetRow struct {
	AccountType  string  `json:"account_type"`
	GroupName    string  `json:"group_name"`
	SubgroupName string  `json:"subgroup_name"`
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Balance      float64 `json:"balance"`
}

type ReportCashFlowRow struct {
	Date          time.Time `json:"date"`
	JournalNumber string    `json:"journal_number"`
	Description   string    `json:"description"`
	Debit         float64   `json:"debit"`
	Credit        float64   `json:"credit"`
}

type ReportEquityRow struct {
	Code           string  `json:"code"`
	Name           string  `json:"name"`
	OpeningBalance float64 `json:"opening_balance"`
	NetIncome      float64 `json:"net_income"`
	ClosingBalance float64 `json:"closing_balance"`
}

type ReportJournalBookRow struct {
	Id            uuid.UUID `json:"id"`
	JournalNumber string    `json:"journal_number"`
	Date          time.Time `json:"date"`
	Type          string    `json:"type"`
	Description   string    `json:"description"`
	TotalDebit    float64   `json:"total_debit"`
	TotalCredit   float64   `json:"total_credit"`
	LineId        uuid.UUID `json:"line_id"`
	CoaCode       string    `json:"coa_code"`
	CoaName       string    `json:"coa_name"`
	LineDebit     float64   `json:"line_debit"`
	LineCredit    float64   `json:"line_credit"`
}

type IReportRepository interface {
	GetGeneralLedger(companyID uuid.UUID, coaId uuid.UUID, dateFrom, dateTo string) ([]ReportGeneralLedgerRow, error)
	GetTrialBalance(companyID uuid.UUID, dateFrom, dateTo string) ([]ReportTrialBalanceRow, error)
	GetProfitLoss(companyID uuid.UUID, dateFrom, dateTo string) ([]ReportProfitLossRow, error)
	GetBalanceSheet(companyID uuid.UUID, dateTo string) ([]ReportBalanceSheetRow, error)
	GetCashFlow(companyID uuid.UUID, cashCoaIds []uuid.UUID, dateFrom, dateTo string) ([]ReportCashFlowRow, error)
	GetEquityStatement(companyID uuid.UUID, dateFrom, dateTo string) ([]ReportEquityRow, error)
	GetJournalBook(companyID uuid.UUID, dateFrom, dateTo string) ([]ReportJournalBookRow, error)
}

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) IReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetGeneralLedger(companyID uuid.UUID, coaId uuid.UUID, dateFrom, dateTo string) ([]ReportGeneralLedgerRow, error) {
	var rows []ReportGeneralLedgerRow
	query := `
		SELECT
			je.date,
			je.journal_number,
			je.description AS je_description,
			jl.description AS line_description,
			jl.debit,
			jl.credit,
			c.code AS account_code,
			c.name AS account_name
		FROM journal_lines jl
		JOIN journal_entries je ON jl.journal_entry_id = je.id
		JOIN coa c ON jl.coa_id = c.id
		WHERE je.company_id = ?
		  AND je.status = 'posted'
		  AND jl.coa_id = ?
		  AND je.date BETWEEN ? AND ?
		ORDER BY je.date ASC, je.created_at ASC
	`
	err := r.db.Raw(query, companyID, coaId, dateFrom, dateTo).Scan(&rows).Error
	return rows, err
}

func (r *ReportRepository) GetTrialBalance(companyID uuid.UUID, dateFrom, dateTo string) ([]ReportTrialBalanceRow, error) {
	var rows []ReportTrialBalanceRow
	query := `
		SELECT
			cg.type AS account_type,
			c.code,
			c.name,
			SUM(jl.debit) AS total_debit,
			SUM(jl.credit) AS total_credit,
			SUM(jl.debit) - SUM(jl.credit) AS balance
		FROM journal_lines jl
		JOIN journal_entries je ON jl.journal_entry_id = je.id
		JOIN coa c ON jl.coa_id = c.id
		JOIN coa_subgroups cs ON c.subgroup_id = cs.id
		JOIN coa_groups cg ON cs.group_id = cg.id
		WHERE je.company_id = ?
		  AND je.status = 'posted'
		  AND je.date BETWEEN ? AND ?
		GROUP BY cg.type, c.code, c.name
		ORDER BY c.code ASC
	`
	err := r.db.Raw(query, companyID, dateFrom, dateTo).Scan(&rows).Error
	return rows, err
}

func (r *ReportRepository) GetProfitLoss(companyID uuid.UUID, dateFrom, dateTo string) ([]ReportProfitLossRow, error) {
	var rows []ReportProfitLossRow
	query := `
		SELECT
			cg.type,
			cg.name AS group_name,
			cs.name AS subgroup_name,
			c.code,
			c.name,
			SUM(jl.credit) - SUM(jl.debit) AS balance
		FROM journal_lines jl
		JOIN journal_entries je ON jl.journal_entry_id = je.id
		JOIN coa c ON jl.coa_id = c.id
		JOIN coa_subgroups cs ON c.subgroup_id = cs.id
		JOIN coa_groups cg ON cs.group_id = cg.id
		WHERE je.company_id = ?
		  AND je.status = 'posted'
		  AND cg.type IN ('revenue', 'expense')
		  AND je.date BETWEEN ? AND ?
		GROUP BY cg.type, cg.name, cs.name, c.code, c.name
	`
	err := r.db.Raw(query, companyID, dateFrom, dateTo).Scan(&rows).Error
	return rows, err
}

func (r *ReportRepository) GetBalanceSheet(companyID uuid.UUID, dateTo string) ([]ReportBalanceSheetRow, error) {
	var rows []ReportBalanceSheetRow
	query := `
		SELECT
			cg.type AS account_type,
			cg.name AS group_name,
			cs.name AS subgroup_name,
			c.code,
			c.name,
			SUM(jl.debit) - SUM(jl.credit) AS balance
		FROM journal_lines jl
		JOIN journal_entries je ON jl.journal_entry_id = je.id
		JOIN coa c ON jl.coa_id = c.id
		JOIN coa_subgroups cs ON c.subgroup_id = cs.id
		JOIN coa_groups cg ON cs.group_id = cg.id
		WHERE je.company_id = ?
		  AND je.status = 'posted'
		  AND cg.type IN ('asset', 'liability', 'equity')
		  AND je.date <= ?
		GROUP BY cg.type, cg.name, cs.name, c.code, c.name

		UNION ALL

		SELECT
			'equity' AS account_type,
			'Equities' AS group_name,
			'Current Year Earnings' AS subgroup_name,
			'3030101' AS code,
			'Current Year Net Income' AS name,
			SUM(jl.debit) - SUM(jl.credit) AS balance
		FROM journal_lines jl
		JOIN journal_entries je ON jl.journal_entry_id = je.id
		JOIN coa c ON jl.coa_id = c.id
		JOIN coa_subgroups cs ON c.subgroup_id = cs.id
		JOIN coa_groups cg ON cs.group_id = cg.id
		WHERE je.company_id = ?
		  AND je.status = 'posted'
		  AND cg.type IN ('revenue', 'expense')
		  AND je.date <= ?
		HAVING SUM(jl.debit) - SUM(jl.credit) != 0
	`
	err := r.db.Raw(query, companyID, dateTo, companyID, dateTo).Scan(&rows).Error
	return rows, err
}

func (r *ReportRepository) GetCashFlow(companyID uuid.UUID, cashCoaIds []uuid.UUID, dateFrom, dateTo string) ([]ReportCashFlowRow, error) {
	var rows []ReportCashFlowRow
	if len(cashCoaIds) == 0 {
		return rows, nil
	}
	query := `
		SELECT
			je.date,
			je.journal_number,
			je.description,
			jl.debit,
			jl.credit
		FROM journal_lines jl
		JOIN journal_entries je ON jl.journal_entry_id = je.id
		WHERE je.company_id = ?
		  AND je.status = 'posted'
		  AND jl.coa_id IN (?)
		  AND je.date BETWEEN ? AND ?
		ORDER BY je.date ASC
	`
	err := r.db.Raw(query, companyID, cashCoaIds, dateFrom, dateTo).Scan(&rows).Error
	return rows, err
}

func (r *ReportRepository) GetEquityStatement(companyID uuid.UUID, dateFrom, dateTo string) ([]ReportEquityRow, error) {
	// Not full logic yet, as equity needs more complex handling
	var rows []ReportEquityRow
	return rows, nil
}

func (r *ReportRepository) GetJournalBook(companyID uuid.UUID, dateFrom, dateTo string) ([]ReportJournalBookRow, error) {
	var rows []ReportJournalBookRow
	query := `
		SELECT
			je.id,
			je.journal_number,
			je.date,
			je.type,
			je.description,
			je.total_debit,
			je.total_credit,
			jl.id AS line_id,
			c.code AS coa_code,
			c.name AS coa_name,
			jl.debit AS line_debit,
			jl.credit AS line_credit
		FROM journal_entries je
		JOIN journal_lines jl ON jl.journal_entry_id = je.id
		JOIN coa c ON jl.coa_id = c.id
		WHERE je.company_id = ?
		  AND je.status = 'posted'
		  AND je.date BETWEEN ? AND ?
		ORDER BY je.date ASC, je.created_at ASC, jl.created_at ASC
	`
	err := r.db.Raw(query, companyID, dateFrom, dateTo).Scan(&rows).Error
	return rows, err
}
