package response

import (
	"time"

	"github.com/google/uuid"
)

type COAGroupResponse struct {
	Id            uuid.UUID `json:"id"`
	Code          string    `json:"code"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	NormalBalance string    `json:"normal_balance"`
	Category      *string   `json:"category"`
	Status        int       `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
