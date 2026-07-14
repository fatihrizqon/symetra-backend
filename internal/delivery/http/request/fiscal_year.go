package request

import "github.com/google/uuid"

type FiscalYearCreateRequest struct {
	Name       string `validate:"required,min=1" json:"name"`
	StartDate  string `validate:"required,datetime=2006-01-02" json:"start_date"`
	EndDate    string `validate:"required,datetime=2006-01-02" json:"end_date"`
	PeriodType string `validate:"required,oneof=monthly" json:"period_type"`
}

type FiscalYearUpdateRequest struct {
	Id         uuid.UUID
	Name       string `validate:"required,min=1" json:"name"`
	StartDate  string `validate:"required,datetime=2006-01-02" json:"start_date"`
	EndDate    string `validate:"required,datetime=2006-01-02" json:"end_date"`
	PeriodType string `validate:"required,oneof=monthly" json:"period_type"`
}

type FiscalPeriodCloseRequest struct {
	Reason string `json:"reason"`
}
