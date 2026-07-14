package response

import (
	"time"

	"github.com/google/uuid"
)

type COAResponse struct {
	Id            uuid.UUID            `json:"id"`
	SubgroupId    uuid.UUID            `json:"subgroup_id"`
	SubGroup      *COASubGroupResponse `json:"subgroup,omitempty"`
	Code          string               `json:"code"`
	Group         *COAGroupResponse    `json:"group"`
	Name          string               `json:"name"`
	NormalBalance *string              `json:"normal_balance"`
	ControlType   string               `json:"control_type"`
	Status        int                  `json:"status"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}
