package entity

import (
	"time"
)

type RedisJob struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)"`
	Type      string    `gorm:"type:varchar(255);not null"`
	Payload   []byte    `gorm:"type:jsonb;not null"`
	Status    string    `gorm:"type:varchar(50);not null;default:'PENDING'"`
	Error     string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
