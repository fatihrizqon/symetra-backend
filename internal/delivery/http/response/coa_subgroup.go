package response

import (
	"time"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/google/uuid"
)

type COASubGroupResponse struct {
	Id        uuid.UUID         `json:"id"`
	Code      string            `json:"code"`
	Name      string            `json:"name"`
	Status    int               `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	GroupId   uuid.UUID         `json:"group_id"`
	Group     *COAGroupResponse `json:"group,omitempty"`
}

func NewCOASubGroupResponse(e entity.COASubGroup) COASubGroupResponse {
	resp := COASubGroupResponse{
		Id:        e.Id,
		GroupId:   e.GroupId,
		Code:      e.Code,
		Name:      e.Name,
		Status:    e.Status,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}

	if e.Group != nil {
		groupResp := NewCOAGroupResponse(*e.Group)
		resp.Group = &groupResp
	}

	return resp
}

func NewCOASubGroupResponses(entities []entity.COASubGroup) []COASubGroupResponse {
	var responses = make([]COASubGroupResponse, 0, len(entities))
	for _, e := range entities {
		responses = append(responses, NewCOASubGroupResponse(e))
	}
	return responses
}
