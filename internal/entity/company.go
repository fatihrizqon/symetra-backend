package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (Company) TableName() string {
	return "companies"
}

type Company struct {
	Id        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string     `gorm:"type:varchar(150);not null" json:"name"`
	LegalName string     `gorm:"type:varchar(200)" json:"legal_name,omitempty"`
	TaxID     string     `gorm:"type:varchar(50)" json:"tax_id,omitempty"`
	Address   string     `gorm:"type:text" json:"address,omitempty"`
	Phone     string     `gorm:"type:varchar(30);uniqueIndex" json:"phone,omitempty"`
	Email     string     `gorm:"type:varchar(150);uniqueIndex" json:"email,omitempty"`
	Industry  string     `gorm:"type:varchar(100)" json:"industry,omitempty"`
	Currency  string     `gorm:"type:varchar(10);not null;default:'IDR'" json:"currency"`
	CreatedBy uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	Creator   *User      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

}

func (Company) SearchableFields() []string {
	return []string{"name", "legalname", "email"}
}

func (Company) ApplyFilters(db *gorm.DB, filters map[string][]string) *gorm.DB {
	if values, ok := filters["currency"]; ok {
		db = db.Where("currency IN ?", values)
	}

	return db
}

// ─── CompanyMember ────────────────────────────────────────────────────────────
func (CompanyMember) TableName() string { return "company_members" }

type CompanyMember struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	CompanyId uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:company_user_unique;" json:"company_id"`
	UserId    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:company_user_unique;" json:"user_id"`
	Role      string    `gorm:"type:varchar(50);not null;default:'viewer'" json:"role"`
	InvitedBy uuid.UUID `gorm:"type:uuid;" json:"invited_by"`
	JoinedAt  time.Time `gorm:"autoCreateTime;" json:"joined_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;" json:"updated_at"`
	Company   Company   `gorm:"foreignKey:CompanyId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"company,omitempty"`
	User      User      `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
}
