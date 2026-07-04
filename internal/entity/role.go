package entity

import "github.com/google/uuid"

type Role struct {
	Id          uuid.UUID    `gorm:"type:uuid; primaryKey; default:gen_random_uuid();" json:"id"`
	Name        string       `gorm:"type:character varying; not null; unique;" json:"name"`
	Description string       `gorm:"type:character varying;" json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

func (Role) TableName() string {
	return "roles"
}
