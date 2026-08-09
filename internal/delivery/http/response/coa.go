package response

import (
	"time"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/google/uuid"
)

type COAResponse struct {
	Id            uuid.UUID            `json:"id"`
	SubgroupId    uuid.UUID            `json:"subgroup_id"`
	SubGroup      *COASubGroupResponse `json:"subgroup,omitempty"`
	Code          string               `json:"code"`
	Group         *COAGroupResponse    `json:"group"`
	Name          string               `json:"name"`
	IsContra      bool                 `json:"is_contra"`
	NormalBalance string               `json:"normal_balance"`
	ControlType   string               `json:"control_type"`
	Status        int                  `json:"status"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

func NewCOAResponse(e entity.COA) COAResponse {
	resp := COAResponse{
		Id:            e.Id,
		SubgroupId:    e.SubgroupId,
		Code:          e.Code,
		Name:          e.Name,
		IsContra:      e.IsContra,
		NormalBalance: e.GetAbsoluteNormalBalance(),
		ControlType:   e.ControlType,
		Status:        e.Status,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}

	if e.SubGroup != nil {
		subGroupResp := NewCOASubGroupResponse(*e.SubGroup)
		resp.SubGroup = &subGroupResp
		if e.SubGroup.Group != nil {
			groupResp := NewCOAGroupResponse(*e.SubGroup.Group)
			resp.Group = &groupResp
		}
	}

	return resp
}

func NewCOAResponses(entities []entity.COA) []COAResponse {
	var responses = make([]COAResponse, 0, len(entities))
	for _, e := range entities {
		responses = append(responses, NewCOAResponse(e))
	}
	return responses
}
