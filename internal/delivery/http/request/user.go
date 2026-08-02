package request

import (
	"github.com/google/uuid"
)

type UserCreateRequest struct {
	Username string `validate:"required,min=3,max=20,alphanum" json:"username"`
	Name     string `validate:"required,min=1,max=100" json:"name"`
	Email    string `validate:"required,email,min=1,max=254" json:"email"`
	Password string `validate:"required,min=8,max=100" json:"password"`
}

type UserUpdateRequest struct {
	Id       uuid.UUID
	Username string `validate:"required,min=3,max=20,alphanum" json:"username"`
	Name     string `validate:"required,min=1,max=100" json:"name"`
	Email    string `validate:"required,email,min=1,max=254" json:"email"`
	Password string `validate:"omitempty,min=8,max=100" json:"password"`
}
