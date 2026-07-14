package entity

import (
	"time"

	"github.com/google/uuid"
)

type FiscalYearStatus string

const (
	FiscalYearStatusDraft  FiscalYearStatus = "draft"
	FiscalYearStatusActive FiscalYearStatus = "active"
	FiscalYearStatusClosed FiscalYearStatus = "closed"
)

func (FiscalYear) TableName() string { return "fiscal_years" }

type FiscalYear struct {
	Id             uuid.UUID        `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	CompanyId      uuid.UUID        `gorm:"type:uuid;not null;index;" json:"company_id"`
	Name           string           `gorm:"type:character varying;not null;" json:"name"`
	StartDate      time.Time        `gorm:"type:date;not null;" json:"start_date"`
	EndDate        time.Time        `gorm:"type:date;not null;" json:"end_date"`
	Status         FiscalYearStatus `gorm:"type:character varying;not null;default:'draft';" json:"status"`
	PeriodType     string           `gorm:"type:character varying;not null;default:'monthly';" json:"period_type"`
	ClosingJEId    *uuid.UUID       `gorm:"type:uuid;" json:"closing_je_id"`
	ClosedAt       *time.Time       `json:"closed_at"`
	CreatedAt      time.Time        `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt      time.Time        `gorm:"autoUpdateTime;" json:"updated_at"`
	Periods        []FiscalPeriod   `gorm:"foreignKey:FiscalYearId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"periods"`
}
