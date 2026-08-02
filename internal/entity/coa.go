package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (COA) TableName() string { return "coa" }

type COA struct {
	Id           uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	CompanyId    uuid.UUID    `gorm:"type:uuid;not null;index;" json:"company_id"`
	SubgroupId   uuid.UUID    `gorm:"type:uuid;not null;index;" json:"subgroup_id"`
	SubGroup     *COASubGroup `gorm:"foreignKey:SubgroupId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"subgroup,omitempty"`
	Code         string       `gorm:"type:character varying;not null; uniqueIndex;" json:"code"`
	Name         string       `gorm:"type:character varying;not null;" json:"name"`
	CurrencyCode string       `gorm:"type:character varying;not null;default:'IDR';" json:"currency_code"`
	IsContra     bool         `gorm:"type:boolean;not null;default:false;" json:"is_contra"`
	ControlType  string       `gorm:"type:character varying;not null;default:'';" json:"control_type"`
	Active       bool         `gorm:"type:boolean;not null;default:true;" json:"active"`
	Status       int          `gorm:"type:int;not null;default:1;" json:"status"`
	CreatedAt    time.Time    `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt    time.Time    `gorm:"autoUpdateTime;" json:"updated_at"`
}

func (COA) SearchableFields() []string { return []string{"coa.name", "coa.code"} }

func (COA) ApplyFilters(db *gorm.DB, filters map[string][]string) *gorm.DB {
	if values, ok := filters["status"]; ok {
		db = db.Where("coa.status IN ?", values)
	}
	if values, ok := filters["group_type"]; ok {
		db = db.
			Joins("JOIN coa_subgroups ON coa_subgroups.id = coa.subgroup_id").
			Joins("JOIN coa_groups ON coa_groups.id = coa_subgroups.group_id").
			Where("coa_groups.type IN ?", values)
	}
	return db
}

func (c *COA) GetAbsoluteNormalBalance() string {
	if c.SubGroup == nil || c.SubGroup.Group == nil {
		return ""
	}

	isGroupDebit := c.SubGroup.Group.Type.IsDebitNormal()

	isAccountDebit := isGroupDebit
	if c.IsContra {
		isAccountDebit = !isGroupDebit
	}

	if isAccountDebit {
		return "DEBIT"
	}
	return "CREDIT"
}
