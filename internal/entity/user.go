package entity

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (User) TableName() string {
	return "users"
}

type User struct {
	Id              uuid.UUID  `gorm:"type:uuid; primaryKey; default:gen_random_uuid();" json:"id"`
	Username        string     `gorm:"type:character varying; not null; unique;" json:"username"`
	Name            string     `gorm:"type:character varying; not null;" json:"name"`
	Email           string     `gorm:"type:character varying; not null; unique;" json:"email"`
	Status          int        `gorm:"type:int; not null; default:1;" json:"status"`
	EmailVerifiedAt time.Time  `gorm:"default:null;" json:"email_verified_at"`
	Password        string     `json:"password"`
	Roles           []Role     `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	AvatarFileId    *uuid.UUID `gorm:"type:uuid;default:null;" json:"avatar_file_id"`
	Avatar          *File      `gorm:"foreignKey:AvatarFileId" json:"avatar,omitempty"`
	CreatedAt       time.Time  `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime;" json:"updated_at"`
	DeletedAt       time.Time  `gorm:"default:null;" json:"deleted_at"`
}

func (User) SearchableFields() []string {
	return []string{"username", "email"}
}

// Supports multi-value filters via repeated params (e.g. ?status=1&status=2).
func (User) ApplyFilters(db *gorm.DB, filters map[string][]string) *gorm.DB {
	if values, ok := filters["status"]; ok {
		db = db.Where("status IN ?", values)
	}
	if values, ok := filters["verified"]; ok && len(values) == 1 {
		switch strings.ToLower(strings.TrimSpace(values[0])) {
		case "true":
			db = db.Where("email_verified_at IS NOT NULL")
		case "false":
			db = db.Where("email_verified_at IS NULL")
		}
	}

	return db
}
