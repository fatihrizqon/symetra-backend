package entity

import (
	"time"

	"github.com/google/uuid"
)

func (Credential) TableName() string {
	return "credentials"
}

type Credential struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey"`
	SessionID    uuid.UUID  `gorm:"type:uuid;index;not null"`
	Type         string     `gorm:"type:varchar(50);not null"`
	RefreshToken string     `gorm:"type:text;uniqueIndex;not null"`
	ExpiresAt    time.Time  `gorm:"index"`
	RevokedAt    *time.Time `gorm:"default:null"`
	CreatedAt    time.Time
}
