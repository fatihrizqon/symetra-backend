package request

import "github.com/google/uuid"

type COACreateRequest struct {
	SubgroupId    uuid.UUID `validate:"required" json:"subgroup_id"`
	Code          string    `validate:"required,min=1" json:"code"`
	Name          string    `validate:"required,min=1" json:"name"`
	IsContra      bool      `json:"is_contra"`
	ControlType   string    `validate:"omitempty,oneof=ar ap cash bank" json:"control_type"`
}

type COAUpdateRequest struct {
	Id            uuid.UUID
	SubgroupId    uuid.UUID `validate:"required" json:"subgroup_id"`
	Code          string    `validate:"required,min=1,max=20" json:"code"`
	Name          string    `validate:"required,min=1,max=20" json:"name"`
	IsContra      bool      `json:"is_contra"`
	ControlType   string    `validate:"omitempty,oneof=ar ap cash bank" json:"control_type"`
}
