package request

import "github.com/google/uuid"

type CompanyCreateRequest struct {
	Name      string `validate:"required,min=1,max=150"       json:"name"`
	LegalName string `validate:"max=200"                      json:"legal_name"`
	TaxID     string `validate:"max=50"                       json:"tax_id"`
	Address   string `json:"address"`
	Phone     string `validate:"required,min=1,max=30"        json:"phone"`
	Email     string `validate:"required,min=1,max=150,email" json:"email"`
	Industry  string `validate:"max=100"                      json:"industry"`
	Currency  string `validate:"max=10"                       json:"currency"`
}

type CompanyUpdateRequest struct {
	Id        uuid.UUID
	Name      string `validate:"required,min=1,max=150" json:"name"`
	LegalName string `validate:"max=200"                 json:"legal_name"`
	TaxID     string `validate:"max=50"                  json:"tax_id"`
	Address   string `json:"address"`
	Phone     string `validate:"max=30"                  json:"phone"`
	Email     string `validate:"email,max=150" json:"email"`
	Industry  string `validate:"max=100"                 json:"industry"`
	Currency  string `validate:"max=10"                  json:"currency"`
}

type AssignMemberRequest struct {
	CompanyId uuid.UUID
	UserID    uuid.UUID `validate:"required" json:"user_id"`
	Role      string    `validate:"required,oneof=owner admin accountant viewer" json:"role"`
}

type UpdateMemberRoleRequest struct {
	CompanyId uuid.UUID
	UserID    uuid.UUID `validate:"required" json:"user_id"`
	Role      string    `validate:"required,oneof=owner admin accountant viewer" json:"role"`
}
