package response

import (
	"time"

	"github.com/google/uuid"
)

type CompanyResponse struct {
	Id        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	LegalName string     `json:"legal_name"`
	TaxID     string     `json:"tax_id"`
	Address   string     `json:"address"`
	Phone     string     `json:"phone"`
	Email     string     `json:"email"`
	Industry  string     `json:"industry"`
	Currency  string     `json:"currency"`
	CreatedBy uuid.UUID  `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

}

type CompanyMemberResponse struct {
	Id        uuid.UUID       `json:"id"`
	CompanyId uuid.UUID       `json:"company_id"`
	UserId    uuid.UUID       `json:"user_id"`
	Username  string          `json:"username"`
	Name      string          `json:"name"`
	Email     string          `json:"email"`
	Role      string          `json:"role"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Company   CompanyResponse `json:"company"`
}

type MyCompanyResponse struct {
	CompanyResponse
	Role string `json:"role"`
}
