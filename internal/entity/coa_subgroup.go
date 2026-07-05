package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (COASubGroup) TableName() string { return "coa_subgroups" }

type COASubGroup struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	CompanyId uuid.UUID `gorm:"type:uuid;not null;index;" json:"company_id"` // ← NEW
	GroupId   uuid.UUID `gorm:"type:uuid;not null;index;" json:"group_id"`
	Group     *COAGroup `gorm:"foreignKey:GroupId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"group,omitempty"`
	Code      string    `gorm:"type:character varying;not null;" json:"code"`
	Name      string    `gorm:"type:character varying;not null;" json:"name"`
	Status    int       `gorm:"type:int;not null;default:1;" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;" json:"updated_at"`
}

func (COASubGroup) SearchableFields() []string { return []string{"code", "name"} }

func (COASubGroup) ApplyFilters(db *gorm.DB, filters map[string][]string) *gorm.DB {
	if values, ok := filters["group_id"]; ok {
		db = db.Where("coa_subgroups.group_id IN ?", values)
	}
	if values, ok := filters["status"]; ok {
		db = db.Where("coa_subgroups.status IN ?", values)
	}
	return db
}
