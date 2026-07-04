package entity

import (
	"time"

	"github.com/google/uuid"
)

func (File) TableName() string {
	return "files"
}

type File struct {
	Id           uuid.UUID `gorm:"type:uuid; primaryKey; default:gen_random_uuid();" json:"id"`
	OriginalName string    `gorm:"type:character varying; not null;" json:"original_name"`
	Path         string    `gorm:"type:character varying; not null;" json:"path"`
	MimeType     string    `gorm:"type:character varying; not null;" json:"mime_type"`
	Size         int64     `gorm:"type:bigint; not null;" json:"size"`
	UploadedBy   uuid.UUID `gorm:"type:uuid; not null;" json:"uploaded_by"`
	CreatedAt    time.Time `gorm:"autoCreateTime;" json:"created_at"`
}
