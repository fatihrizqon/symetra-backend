package response

import (
	"time"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/google/uuid"
)

type COAGroupResponse struct {
	Id        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Category  *string   `json:"category"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewCOAGroupResponse(e entity.COAGroup) COAGroupResponse {
	return COAGroupResponse{
		Id:        e.Id,
		Code:      e.Code,
		Name:      e.Name,
		Type:      string(e.Type),
		Category:  e.Category,
		Status:    e.Status,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func NewCOAGroupResponses(entities []entity.COAGroup) []COAGroupResponse {
	var responses = make([]COAGroupResponse, 0, len(entities))
	for _, e := range entities {
		responses = append(responses, NewCOAGroupResponse(e))
	}
	return responses
}
