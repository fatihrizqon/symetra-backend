package response

import (
	"time"

	"github.com/google/uuid"
)

type FileResponse struct {
	Id           uuid.UUID `json:"id"`
	OriginalName string    `json:"original_name"`
	MimeType     string    `json:"mime_type"`
	Size         int64     `json:"size"`
	URL          string    `json:"url"`
	UploadedBy   uuid.UUID `json:"uploaded_by"`
	CreatedAt    time.Time `json:"created_at"`
}
