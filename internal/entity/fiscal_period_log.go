package entity

import (
	"time"

	"github.com/google/uuid"
)

func (FiscalPeriodLog) TableName() string { return "fiscal_period_logs" }

type FiscalPeriodLog struct {
	Id             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	FiscalPeriodId uuid.UUID  `gorm:"type:uuid;not null;index;" json:"fiscal_period_id"`
	Action         string     `gorm:"type:character varying;not null;" json:"action"`
	FromStatus     string     `gorm:"type:character varying;not null;" json:"from_status"`
	ToStatus       string     `gorm:"type:character varying;not null;" json:"to_status"`
	Reason         string     `gorm:"type:text;" json:"reason"`
	PerformedBy    *uuid.UUID `gorm:"type:uuid;" json:"performed_by"`
	PerformedAt    time.Time  `gorm:"autoCreateTime;" json:"performed_at"`
}
