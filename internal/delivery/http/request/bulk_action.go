package request

import "github.com/google/uuid"

// BulkActionRequest is a generic request for bulk operations (destroy, approve, reject, etc).
type BulkActionRequest struct {
	Ids []uuid.UUID `json:"ids" validate:"required,min=1"`
}
