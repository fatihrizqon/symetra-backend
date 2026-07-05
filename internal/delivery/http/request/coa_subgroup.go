package request

import "github.com/google/uuid"

type COASubGroupCreateRequest struct {
	GroupId uuid.UUID `validate:"required" json:"group_id"`
	Code    string    `validate:"required,min=1" json:"code"`
	Name    string    `validate:"required,min=1" json:"name"`
}

type COASubGroupUpdateRequest struct {
	Id      uuid.UUID
	GroupId uuid.UUID `validate:"required" json:"group_id"`
	Code    string    `validate:"required,min=1,max=20" json:"code"`
	Name    string    `validate:"required,min=1,max=20" json:"name"`
}
