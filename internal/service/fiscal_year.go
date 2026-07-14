package service

import (
	"fmt"
	"time"

	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/request"
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/response"
	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type IFiscalYearService interface {
	Create(companyID uuid.UUID, req request.FiscalYearCreateRequest) (entity.FiscalYear, error)
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]response.FiscalYearResponse, int, error)
	FindById(companyID, id uuid.UUID) (response.FiscalYearResponse, error)
	Activate(companyID, id uuid.UUID) (entity.FiscalYear, error)
	ClosePeriod(companyID, periodId uuid.UUID, userID *uuid.UUID, req request.FiscalPeriodCloseRequest) (entity.FiscalPeriod, error)
}

type FiscalYearService struct {
	Repo     repository.IFiscalYearRepository
	validate *validator.Validate
}

func NewFiscalYearService(repo repository.IFiscalYearRepository, validate *validator.Validate) IFiscalYearService {
	return &FiscalYearService{Repo: repo, validate: validate}
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

	fy := entity.FiscalYear{
		CompanyId:  companyID,
		Name:       req.Name,
		StartDate:  startDate,
		EndDate:    endDate,
		Status:     entity.FiscalYearStatusDraft,
		PeriodType: req.PeriodType,
	}

	// Generate periods if monthly
	if fy.PeriodType == "monthly" {
		periods := generateMonthlyPeriods(startDate, endDate)
		fy.Periods = periods
	}

	return s.Repo.Create(fy)
}

func (s *FiscalYearService) FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]response.FiscalYearResponse, int, error) {
	entities, totalCount, err := s.Repo.FindAll(companyID, qp)
	if err != nil {
		return nil, 0, err
	}
	resps := make([]response.FiscalYearResponse, 0, len(entities))
	for _, v := range entities {
		resps = append(resps, mapFiscalYear(v))
	}
	return resps, totalCount, nil
}

func (s *FiscalYearService) FindById(companyID, id uuid.UUID) (response.FiscalYearResponse, error) {
	ent, err := s.Repo.FindById(companyID, id)
	if err != nil {
		return response.FiscalYearResponse{}, err
	}
	return mapFiscalYear(ent), nil
}

func (s *FiscalYearService) Activate(companyID, id uuid.UUID) (entity.FiscalYear, error) {
	// Check if any other is active
	if activeFY, err := s.Repo.GetActive(companyID); err == nil {
		return activeFY, fmt.Errorf("FISCAL_YEAR_ALREADY_ACTIVE")
	}

	fy, err := s.Repo.FindById(companyID, id)
	if err != nil {
		return fy, err
	}
	
	if fy.Status != entity.FiscalYearStatusDraft {
		return fy, fmt.Errorf("INVALID_STATE_TRANSITION")
	}

	fy.Status = entity.FiscalYearStatusActive
	err = s.Repo.Update(fy)
	return fy, err
}

func (s *FiscalYearService) ClosePeriod(companyID, periodId uuid.UUID, userID *uuid.UUID, req request.FiscalPeriodCloseRequest) (entity.FiscalPeriod, error) {
	period, err := s.Repo.FindPeriodById(companyID, periodId)
	if err != nil {
		return period, err
	}

	if period.Status != entity.FiscalPeriodStatusOpen {
		return period, fmt.Errorf("INVALID_STATE_TRANSITION")
	}

	// Move to closed
	oldStatus := string(period.Status)
	period.Status = entity.FiscalPeriodStatusClosed
	if err := s.Repo.UpdatePeriod(period); err != nil {
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
	_ = s.Repo.LogPeriodAction(log)

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

func mapFiscalYear(v entity.FiscalYear) response.FiscalYearResponse {
	resp := response.FiscalYearResponse{
		Id:          v.Id,
		CompanyId:   v.CompanyId,
		Name:        v.Name,
		StartDate:   v.StartDate.Format("2006-01-02"),
		EndDate:     v.EndDate.Format("2006-01-02"),
		Status:      string(v.Status),
		PeriodType:  v.PeriodType,
		ClosingJEId: v.ClosingJEId,
		ClosedAt:    v.ClosedAt,
		CreatedAt:   v.CreatedAt,
		UpdatedAt:   v.UpdatedAt,
	}
	
	if len(v.Periods) > 0 {
		var pResps []response.FiscalPeriodResponse
		for _, p := range v.Periods {
			pResps = append(pResps, response.FiscalPeriodResponse{
				Id:           p.Id,
				FiscalYearId: p.FiscalYearId,
				Name:         p.Name,
				PeriodNumber: p.PeriodNumber,
				StartDate:    p.StartDate.Format("2006-01-02"),
				EndDate:      p.EndDate.Format("2006-01-02"),
				Status:       string(p.Status),
				CreatedAt:    p.CreatedAt,
				UpdatedAt:    p.UpdatedAt,
			})
		}
		resp.Periods = pResps
	}
	
	return resp
}
