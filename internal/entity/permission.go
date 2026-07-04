package entity

import "github.com/google/uuid"

type Permission struct {
	Id          uuid.UUID `gorm:"type:uuid; primaryKey; default:gen_random_uuid();" json:"id"`
	Name        string    `gorm:"type:character varying; not null; unique;" json:"name"`
	Description string    `gorm:"type:character varying;" json:"description"`
}

func (Permission) TableName() string {
	return "permissions"
}
