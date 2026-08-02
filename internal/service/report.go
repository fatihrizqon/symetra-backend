package service

import (
	"errors"

	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/response"
	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/google/uuid"
)

type IReportService interface {
	GetGeneralLedger(companyId uuid.UUID, coaId uuid.UUID, dateFrom, dateTo string) (response.GeneralLedgerResponse, error)
	GetTrialBalance(companyId uuid.UUID, dateFrom, dateTo string) (response.TrialBalanceResponse, error)
	GetProfitLoss(companyId uuid.UUID, dateFrom, dateTo string) (response.ProfitLossResponse, error)
	GetBalanceSheet(companyId uuid.UUID, dateTo string) (response.BalanceSheetResponse, error)
	GetCashFlow(companyId uuid.UUID, dateFrom, dateTo string) (response.CashFlowResponse, error)
	GetEquityStatement(companyId uuid.UUID, dateFrom, dateTo string) (response.EquityStatementResponse, error)
	GetJournalBook(companyId uuid.UUID, dateFrom, dateTo string) (response.JournalBookResponse, error)
}

type ReportService struct {
	repo       repository.IReportRepository
	configRepo repository.ICompanyConfigurationRepository
}

func NewReportService(repo repository.IReportRepository, configRepo repository.ICompanyConfigurationRepository) IReportService {
	return &ReportService{
		repo:       repo,
		configRepo: configRepo,
	}
}

func (s *ReportService) GetGeneralLedger(companyId uuid.UUID, coaId uuid.UUID, dateFrom, dateTo string) (response.GeneralLedgerResponse, error) {
	var res response.GeneralLedgerResponse
	if dateFrom == "" || dateTo == "" {
		return res, errors.New("date_from and date_to are required")
	}

	rows, err := s.repo.GetGeneralLedger(companyId, coaId, dateFrom, dateTo)
	if err != nil {
		return res, err
	}

	res.Account.Id = coaId
	if len(rows) > 0 {
		res.Account.Code = rows[0].AccountCode
		res.Account.Name = rows[0].AccountName
	}
	res.OpeningBalance = 0 // In real system, this would be computed from previous periods

	var totalDebit, totalCredit float64
	runningBalance := res.OpeningBalance

	for _, r := range rows {
		runningBalance += (r.Debit - r.Credit)
		entry := response.GeneralLedgerEntry{
			Date:           r.Date,
			JournalNumber:  r.JournalNumber,
			Description:    r.JeDescription,
			Debit:          r.Debit,
			Credit:         r.Credit,
			RunningBalance: runningBalance,
		}
		res.Entries = append(res.Entries, entry)
		totalDebit += r.Debit
		totalCredit += r.Credit
	}

	res.Totals.TotalDebit = totalDebit
	res.Totals.TotalCredit = totalCredit
	res.Totals.ClosingBalance = runningBalance

	return res, nil
}

func (s *ReportService) GetTrialBalance(companyId uuid.UUID, dateFrom, dateTo string) (response.TrialBalanceResponse, error) {
	var res response.TrialBalanceResponse
	res.Period.From = dateFrom
	res.Period.To = dateTo

	if dateFrom == "" || dateTo == "" {
		return res, errors.New("date_from and date_to are required")
	}

	rows, err := s.repo.GetTrialBalance(companyId, dateFrom, dateTo)
	if err != nil {
		return res, err
	}

	var grandTotalDebit, grandTotalCredit float64

	for _, r := range rows {
		res.Accounts = append(res.Accounts, response.TrialBalanceAccount{
			AccountType: r.AccountType,
			Code:        r.Code,
			Name:        r.Name,
			TotalDebit:  r.TotalDebit,
			TotalCredit: r.TotalCredit,
			Balance:     r.Balance,
		})
		grandTotalDebit += r.TotalDebit
		grandTotalCredit += r.TotalCredit
	}

	res.Summary.GrandTotalDebit = grandTotalDebit
	res.Summary.GrandTotalCredit = grandTotalCredit
	// Using a small epsilon to account for floating point errors
	res.Summary.IsBalanced = false
	if (grandTotalDebit - grandTotalCredit) < 0.01 && (grandTotalDebit - grandTotalCredit) > -0.01 {
		res.Summary.IsBalanced = true
	}

	return res, nil
}

func (s *ReportService) GetProfitLoss(companyId uuid.UUID, dateFrom, dateTo string) (response.ProfitLossResponse, error) {
	var res response.ProfitLossResponse
	res.Period.From = dateFrom
	res.Period.To = dateTo

	if dateFrom == "" || dateTo == "" {
		return res, errors.New("date_from and date_to are required")
	}

	rows, err := s.repo.GetProfitLoss(companyId, dateFrom, dateTo)
	if err != nil {
		return res, err
	}

	for _, r := range rows {
		acc := response.ProfitLossAccount{
			Code:    r.Code,
			Name:    r.Name,
			Balance: r.Balance, // logic is credit - debit from repo
		}
		if r.Type == "revenue" {
			res.Revenue.Accounts = append(res.Revenue.Accounts, acc)
			res.Revenue.Total += r.Balance
		} else if r.Type == "expense" {
			// For expenses, we flip the balance so it shows as a positive expense number
			acc.Balance = -r.Balance
			res.Expense.Accounts = append(res.Expense.Accounts, acc)
			res.Expense.Total += acc.Balance
		}
	}

	res.NetIncome = res.Revenue.Total - res.Expense.Total
	return res, nil
}

func (s *ReportService) GetBalanceSheet(companyId uuid.UUID, dateTo string) (response.BalanceSheetResponse, error) {
	var res response.BalanceSheetResponse
	res.AsOf = dateTo

	if dateTo == "" {
		return res, errors.New("date_to is required")
	}

	rows, err := s.repo.GetBalanceSheet(companyId, dateTo)
	if err != nil {
		return res, err
	}

	assetGroups := make(map[string]*response.BalanceSheetGroup)
	liabilityGroups := make(map[string]*response.BalanceSheetGroup)
	equityGroups := make(map[string]*response.BalanceSheetGroup)

	for _, r := range rows {
		acc := response.BalanceSheetAccount{
			Code:    r.Code,
			Name:    r.Name,
			Balance: r.Balance,
		}

		if r.AccountType == "asset" {
			if _, ok := assetGroups[r.GroupName]; !ok {
				assetGroups[r.GroupName] = &response.BalanceSheetGroup{Name: r.GroupName}
			}
			assetGroups[r.GroupName].Accounts = append(assetGroups[r.GroupName].Accounts, acc)
			assetGroups[r.GroupName].Subtotal += acc.Balance
			res.Assets.Total += acc.Balance
		} else if r.AccountType == "liability" {
			acc.Balance = -acc.Balance // flip sign for liabilities
			if _, ok := liabilityGroups[r.GroupName]; !ok {
				liabilityGroups[r.GroupName] = &response.BalanceSheetGroup{Name: r.GroupName}
			}
			liabilityGroups[r.GroupName].Accounts = append(liabilityGroups[r.GroupName].Accounts, acc)
			liabilityGroups[r.GroupName].Subtotal += acc.Balance
			res.Liabilities.Total += acc.Balance
		} else if r.AccountType == "equity" {
			acc.Balance = -acc.Balance // flip sign for equity
			if _, ok := equityGroups[r.GroupName]; !ok {
				equityGroups[r.GroupName] = &response.BalanceSheetGroup{Name: r.GroupName}
			}
			equityGroups[r.GroupName].Accounts = append(equityGroups[r.GroupName].Accounts, acc)
			equityGroups[r.GroupName].Subtotal += acc.Balance
			res.Equity.Total += acc.Balance
		}
	}

	for _, v := range assetGroups {
		res.Assets.Groups = append(res.Assets.Groups, *v)
	}
	for _, v := range liabilityGroups {
		res.Liabilities.Groups = append(res.Liabilities.Groups, *v)
	}
	for _, v := range equityGroups {
		res.Equity.Groups = append(res.Equity.Groups, *v)
	}

	res.TotalLiabilitiesAndEquity = res.Liabilities.Total + res.Equity.Total
	
	diff := res.Assets.Total - res.TotalLiabilitiesAndEquity
	if diff < 0.01 && diff > -0.01 {
		res.IsBalanced = true
	}

	return res, nil
}

func (s *ReportService) GetCashFlow(companyId uuid.UUID, dateFrom, dateTo string) (response.CashFlowResponse, error) {
	var res response.CashFlowResponse
	res.Period.From = dateFrom
	res.Period.To = dateTo

	if dateFrom == "" || dateTo == "" {
		return res, errors.New("date_from and date_to are required")
	}

	conf, err := s.configRepo.GetByCompanyId(companyId)
	if err != nil {
		return res, err
	}

	var coaIds []uuid.UUID
	if conf.CashAccountId != nil {
		coaIds = append(coaIds, *conf.CashAccountId)
	}
	if conf.BankAccountId != nil {
		coaIds = append(coaIds, *conf.BankAccountId)
	}

	rows, err := s.repo.GetCashFlow(companyId, coaIds, dateFrom, dateTo)
	if err != nil {
		return res, err
	}

	runningBalance := res.OpeningBalance // 0 for now
	var totalInflow, totalOutflow float64

	for _, r := range rows {
		inflow := r.Debit
		outflow := r.Credit
		runningBalance += (inflow - outflow)

		res.Movements = append(res.Movements, response.CashFlowMovement{
			Date:           r.Date,
			JournalNumber:  r.JournalNumber,
			Description:    r.Description,
			Inflow:         inflow,
			Outflow:        outflow,
			RunningBalance: runningBalance,
		})
		totalInflow += inflow
		totalOutflow += outflow
	}

	res.Summary.TotalInflow = totalInflow
	res.Summary.TotalOutflow = totalOutflow
	res.Summary.NetChange = totalInflow - totalOutflow
	res.Summary.ClosingBalance = runningBalance

	return res, nil
}

func (s *ReportService) GetEquityStatement(companyId uuid.UUID, dateFrom, dateTo string) (response.EquityStatementResponse, error) {
	var res response.EquityStatementResponse
	res.Period.From = dateFrom
	res.Period.To = dateTo

	// Placeholder logic
	return res, nil
}

func (s *ReportService) GetJournalBook(companyId uuid.UUID, dateFrom, dateTo string) (response.JournalBookResponse, error) {
	var res response.JournalBookResponse

	if dateFrom == "" || dateTo == "" {
		return res, errors.New("date_from and date_to are required")
	}

	rows, err := s.repo.GetJournalBook(companyId, dateFrom, dateTo)
	if err != nil {
		return res, err
	}

	entryMap := make(map[string]*response.JournalBookEntry)
	var order []string

	for _, r := range rows {
		key := r.Id.String()
		if _, exists := entryMap[key]; !exists {
			entryMap[key] = &response.JournalBookEntry{
				JournalNumber: r.JournalNumber,
				Date:          r.Date,
				Type:          r.Type,
				Description:   r.Description,
				TotalDebit:    r.TotalDebit,
				TotalCredit:   r.TotalCredit,
			}
			order = append(order, key)
		}
		
		entryMap[key].Lines = append(entryMap[key].Lines, response.JournalBookLine{
			CoaCode: r.CoaCode,
			CoaName: r.CoaName,
			Debit:   r.LineDebit,
			Credit:  r.LineCredit,
		})
	}

	for _, key := range order {
		res = append(res, *entryMap[key])
	}

	return res, nil
}
