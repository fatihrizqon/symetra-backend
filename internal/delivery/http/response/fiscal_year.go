package response

import (
	"time"

	"github.com/google/uuid"
)

type FiscalPeriodResponse struct {
	Id           uuid.UUID `json:"id"`
	FiscalYearId uuid.UUID `json:"fiscal_year_id"`
	Name         string    `json:"name"`
	PeriodNumber int       `json:"period_number"`
	StartDate    string    `json:"start_date"`
	EndDate      string    `json:"end_date"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type FiscalYearResponse struct {
	Id          uuid.UUID              `json:"id"`
	CompanyId   uuid.UUID              `json:"company_id"`
	Name        string                 `json:"name"`
	StartDate   string                 `json:"start_date"`
	EndDate     string                 `json:"end_date"`
	Status      string                 `json:"status"`
	PeriodType  string                 `json:"period_type"`
	ClosingJEId *uuid.UUID             `json:"closing_je_id"`
	ClosedAt    *time.Time             `json:"closed_at"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Periods     []FiscalPeriodResponse `json:"periods,omitempty"`
}

type FiscalYearReadinessCheck struct {
	Name    string `json:"name"`
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
}

type FiscalYearReadinessResponse struct {
	IsReady bool                       `json:"is_ready"`
	Checks  []FiscalYearReadinessCheck `json:"checks"`
}
