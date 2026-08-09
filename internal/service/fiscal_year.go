package service

import (
	"fmt"
	"time"

	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/request"
	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type IFiscalYearService interface {
	Create(companyID uuid.UUID, req request.FiscalYearCreateRequest) (entity.FiscalYear, error)
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.FiscalYear, int, error)
	FindById(companyID, id uuid.UUID) (entity.FiscalYear, error)
	Activate(companyID, id uuid.UUID) (entity.FiscalYear, error)
	ClosePeriod(companyID, periodId uuid.UUID, userID *uuid.UUID, req request.FiscalPeriodCloseRequest) (entity.FiscalPeriod, error)
}

type FiscalYearService struct {
	validate              *validator.Validate
	IFiscalYearRepository repository.IFiscalYearRepository
}

func NewFiscalYearService(validate *validator.Validate, repo repository.IFiscalYearRepository) IFiscalYearService {
	return &FiscalYearService{
		validate:              validate,
		IFiscalYearRepository: repo,
	}
}

func (s *FiscalYearService) Create(companyID uuid.UUID, req request.FiscalYearCreateRequest) (entity.FiscalYear, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.FiscalYear{}, err
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	if startDate.After(endDate) {
		return entity.FiscalYear{}, fmt.Errorf("start_date must be before end_date")
	}

	fiscal_year := entity.FiscalYear{
		CompanyId:  companyID,
		Name:       req.Name,
		StartDate:  startDate,
		EndDate:    endDate,
		Status:     entity.FiscalYearStatusDraft,
		PeriodType: req.PeriodType,
	}

	// Generate periods if monthly
	if fiscal_year.PeriodType == "monthly" {
		periods := generateMonthlyPeriods(startDate, endDate)
		fiscal_year.Periods = periods
	}

	return s.IFiscalYearRepository.Create(fiscal_year)
}

func (s *FiscalYearService) FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.FiscalYear, int, error) {
	fiscal_years, totalCount, err := s.IFiscalYearRepository.FindAll(companyID, qp)
	if err != nil {
		return nil, 0, err
	}
	results := make([]entity.FiscalYear, 0, len(fiscal_years))
	for _, fiscal_year := range fiscal_years {
		result := entity.FiscalYear{
			Id:          fiscal_year.Id,
			CompanyId:   fiscal_year.CompanyId,
			Name:        fiscal_year.Name,
			StartDate:   fiscal_year.StartDate,
			EndDate:     fiscal_year.EndDate,
			Status:      fiscal_year.Status,
			PeriodType:  fiscal_year.PeriodType,
			ClosingJEId: fiscal_year.ClosingJEId,
			ClosedAt:    fiscal_year.ClosedAt,
			CreatedAt:   fiscal_year.CreatedAt,
			UpdatedAt:   fiscal_year.UpdatedAt,
		}
		if len(fiscal_year.Periods) > 0 {
			var pResps []entity.FiscalPeriod
			for _, period := range fiscal_year.Periods {
				pResps = append(pResps, entity.FiscalPeriod{
					Id:           period.Id,
					FiscalYearId: period.FiscalYearId,
					Name:         period.Name,
					PeriodNumber: period.PeriodNumber,
					StartDate:    period.StartDate,
					EndDate:      period.EndDate,
					Status:       period.Status,
					CreatedAt:    period.CreatedAt,
					UpdatedAt:    period.UpdatedAt,
				})
			}
			result.Periods = pResps
		}
		results = append(results, result)
	}
	return results, totalCount, nil
}

func (s *FiscalYearService) FindById(companyID, id uuid.UUID) (entity.FiscalYear, error) {
	fiscal_year, err := s.IFiscalYearRepository.FindById(companyID, id)
	if err != nil {
		return entity.FiscalYear{}, err
	}

	result := entity.FiscalYear{
		Id:          fiscal_year.Id,
		CompanyId:   fiscal_year.CompanyId,
		Name:        fiscal_year.Name,
		StartDate:   fiscal_year.StartDate,
		EndDate:     fiscal_year.EndDate,
		Status:      fiscal_year.Status,
		PeriodType:  fiscal_year.PeriodType,
		ClosingJEId: fiscal_year.ClosingJEId,
		ClosedAt:    fiscal_year.ClosedAt,
		CreatedAt:   fiscal_year.CreatedAt,
		UpdatedAt:   fiscal_year.UpdatedAt,
	}
	if len(fiscal_year.Periods) > 0 {
		var pResps []entity.FiscalPeriod
		for _, period := range fiscal_year.Periods {
			pResps = append(pResps, entity.FiscalPeriod{
				Id:           period.Id,
				FiscalYearId: period.FiscalYearId,
				Name:         period.Name,
				PeriodNumber: period.PeriodNumber,
				StartDate:    period.StartDate,
				EndDate:      period.EndDate,
				Status:       period.Status,
				CreatedAt:    period.CreatedAt,
				UpdatedAt:    period.UpdatedAt,
			})
		}
		result.Periods = pResps
	}
	return result, nil
}

func (s *FiscalYearService) Activate(companyID uuid.UUID, id uuid.UUID) (entity.FiscalYear, error) {
	// Check if any other is active
	if activeFY, err := s.IFiscalYearRepository.GetActive(companyID); err == nil {
		return activeFY, fmt.Errorf("FISCAL_YEAR_ALREADY_ACTIVE")
	}

	fiscal_year, err := s.IFiscalYearRepository.FindById(companyID, id)
	if err != nil {
		return fiscal_year, err
	}

	if fiscal_year.Status != entity.FiscalYearStatusDraft {
		return fiscal_year, fmt.Errorf("INVALID_STATE_TRANSITION")
	}

	fiscal_year.Status = entity.FiscalYearStatusActive
	err = s.IFiscalYearRepository.Update(fiscal_year)
	return fiscal_year, err
}

func (s *FiscalYearService) ClosePeriod(companyID uuid.UUID, periodId uuid.UUID, userID *uuid.UUID, req request.FiscalPeriodCloseRequest) (entity.FiscalPeriod, error) {
	period, err := s.IFiscalYearRepository.FindPeriodById(companyID, periodId)
	if err != nil {
		return period, err
	}

	if period.Status != entity.FiscalPeriodStatusOpen {
		return period, fmt.Errorf("INVALID_STATE_TRANSITION")
	}

	// Move to closed
	oldStatus := string(period.Status)
	period.Status = entity.FiscalPeriodStatusClosed
	if err := s.IFiscalYearRepository.UpdatePeriod(period); err != nil {
		return period, err
	}

	// Log action
	log := entity.FiscalPeriodLog{
		FiscalPeriodId: period.Id,
		Action:         "close",
		FromStatus:     oldStatus,
		ToStatus:       string(period.Status),
		Reason:         req.Reason,
		PerformedBy:    userID,
	}
	_ = s.IFiscalYearRepository.LogPeriodAction(log)

	return period, nil
}

func generateMonthlyPeriods(start, end time.Time) []entity.FiscalPeriod {
	var periods []entity.FiscalPeriod
	current := start
	periodNum := 1

	for current.Before(end) || current.Equal(end) {
		// Calculate the end of the current month
		nextMonth := current.AddDate(0, 1, 0)
		endOfMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, current.Location()).Add(-time.Second)

		if endOfMonth.After(end) {
			endOfMonth = end
		}

		periods = append(periods, entity.FiscalPeriod{
			Name:         current.Format("January 2006"),
			PeriodNumber: periodNum,
			StartDate:    current,
			EndDate:      endOfMonth,
			Status:       entity.FiscalPeriodStatusOpen,
		})

		// Move to the first day of next month
		current = endOfMonth.Add(time.Second)
		periodNum++
	}
	return periods
}
