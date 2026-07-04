package request

import "github.com/google/uuid"

type UserCreateRequest struct {
	Username string `validate:"required,min=1,max=20" json:"username"`
	Name     string `validate:"required,min=1" json:"name"`
	Email    string `validate:"required,min=1" json:"email"`
	Password string `validate:"required,min=8" json:"password"`
}

type UserUpdateRequest struct {
	Id       uuid.UUID
	Username string `validate:"required,min=1,max=20" json:"username"`
	Name     string `validate:"required,min=1,max=20" json:"name"`
	Email    string `validate:"required,min=1" json:"email"`
	Password string `validate:"min=8" json:"password"`
}
