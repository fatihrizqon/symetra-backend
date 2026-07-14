package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (COAGroup) TableName() string { return "coa_groups" }

type COAGroupType string

const (
	COAGroupTypeAsset     COAGroupType = "asset"
	COAGroupTypeLiability COAGroupType = "liability"
	COAGroupTypeEquity    COAGroupType = "equity"
	COAGroupTypeRevenue   COAGroupType = "revenue"
	COAGroupTypeCOGS      COAGroupType = "cogs"
	COAGroupTypeExpense   COAGroupType = "expense"
)

var ValidCOAGroupTypes = map[COAGroupType]bool{
	COAGroupTypeAsset:     true,
	COAGroupTypeLiability: true,
	COAGroupTypeEquity:    true,
	COAGroupTypeRevenue:   true,
	COAGroupTypeCOGS:      true,
	COAGroupTypeExpense:   true,
}

func (t COAGroupType) IsDebitNormal() bool {
	return t == COAGroupTypeAsset || t == COAGroupTypeExpense || t == COAGroupTypeCOGS
}

type COAGroup struct {
	Id            uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	CompanyId     uuid.UUID    `gorm:"type:uuid;not null;index;" json:"company_id"`
	Code          string       `gorm:"type:character varying;not null;" json:"code"`
	Name          string       `gorm:"type:character varying;not null;" json:"name"`
	Type          COAGroupType `gorm:"type:character varying;not null;default:'asset';" json:"type"`
	NormalBalance string       `gorm:"type:character varying;not null;" json:"normal_balance"`
	Category      *string      `gorm:"type:character varying;" json:"category"`
	Status        int          `gorm:"type:int;not null;default:1;" json:"status"`
	CreatedAt     time.Time    `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt     time.Time    `gorm:"autoUpdateTime;" json:"updated_at"`
}

func (COAGroup) SearchableFields() []string { return []string{"code", "name"} }

func (COAGroup) ApplyFilters(db *gorm.DB, filters map[string][]string) *gorm.DB {
	if values, ok := filters["code"]; ok {
		db = db.Where("coa_groups.code IN ?", values)
	}
	if values, ok := filters["type"]; ok {
		db = db.Where("coa_groups.type IN ?", values)
	}
	if values, ok := filters["status"]; ok {
		db = db.Where("coa_groups.status IN ?", values)
	}
	return db
}
