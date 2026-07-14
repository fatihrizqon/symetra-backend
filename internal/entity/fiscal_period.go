package entity

import (
	"time"

	"github.com/google/uuid"
)

type FiscalPeriodStatus string

const (
	FiscalPeriodStatusOpen   FiscalPeriodStatus = "open"
	FiscalPeriodStatusClosed FiscalPeriodStatus = "closed"
	FiscalPeriodStatusLocked FiscalPeriodStatus = "locked"
)

func (FiscalPeriod) TableName() string { return "fiscal_periods" }

type FiscalPeriod struct {
	Id           uuid.UUID          `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	FiscalYearId uuid.UUID          `gorm:"type:uuid;not null;index;" json:"fiscal_year_id"`
	Name         string             `gorm:"type:character varying;not null;" json:"name"`
	PeriodNumber int                `gorm:"type:int;not null;" json:"period_number"`
	StartDate    time.Time          `gorm:"type:date;not null;" json:"start_date"`
	EndDate      time.Time          `gorm:"type:date;not null;" json:"end_date"`
	Status       FiscalPeriodStatus `gorm:"type:character varying;not null;default:'open';" json:"status"`
	CreatedAt    time.Time          `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt    time.Time          `gorm:"autoUpdateTime;" json:"updated_at"`
}
