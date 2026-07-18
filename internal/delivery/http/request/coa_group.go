package request

import "github.com/google/uuid"

type COAGroupCreateRequest struct {
	Code          string  `validate:"required,min=1"         json:"code"`
	Name          string  `validate:"required,min=1"         json:"name"`
	Type          string  `validate:"required,oneof=asset liability equity revenue expense" json:"type"`
	Category      *string `json:"category"`
}

type COAGroupUpdateRequest struct {
	Id            uuid.UUID
	Code          string  `validate:"required,min=1,max=20"  json:"code"`
	Name          string  `validate:"required,min=1,max=20"  json:"name"`
	Type          string  `validate:"required,oneof=asset liability equity revenue expense" json:"type"`
	Category      *string `json:"category"`
}
