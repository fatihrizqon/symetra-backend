package request

import "github.com/google/uuid"

type CustomerCreateRequest struct {
	Code    string     `json:"code" validate:"required,min=2"`
	Name    string     `json:"name" validate:"required,min=2"`
	Email   string     `json:"email" validate:"omitempty,email"`
	Phone   string     `json:"phone"`
	Address string     `json:"address"`
	CoaId   *uuid.UUID `json:"coa_id"`
}

type CustomerUpdateRequest struct {
	Id      uuid.UUID  `json:"-"`
	Code    string     `json:"code" validate:"required,min=2"`
	Name    string     `json:"name" validate:"required,min=2"`
	Email   string     `json:"email" validate:"omitempty,email"`
	Phone   string     `json:"phone"`
	Address string     `json:"address"`
	CoaId   *uuid.UUID `json:"coa_id"`
}
