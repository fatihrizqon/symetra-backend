package response

import "github.com/google/uuid"

type SelectDropdownListResponse struct {
	Value uuid.UUID `json:"value"`
	Label string    `json:"label"`
}
