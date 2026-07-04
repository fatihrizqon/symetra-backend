package entity

import (
	"time"

	"github.com/google/uuid"
)

func (Session) TableName() string {
	return "sessions"
}

type Session struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID         uuid.UUID `gorm:"type:uuid;index;not null"`
	UserAgent      string    `gorm:"type:text"`
	IPAddress      string    `gorm:"type:varchar(45)"`
	LastActivityAt *time.Time
	RevokedAt      *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Credentials    []Credential `gorm:"foreignKey:SessionID"`
	User           User         `gorm:"foreignKey:UserID;references:Id"`
}
