package request

import (
	"time"

	"github.com/google/uuid"
)

type JournalLineRequest struct {
	CoaId       uuid.UUID `validate:"required" json:"coa_id"`
	Description *string   `json:"description"`
	Debit       float64   `validate:"min=0" json:"debit"`
	Credit      float64   `validate:"min=0" json:"credit"`
}

type JournalEntryCreateRequest struct {
	Date        time.Time            `validate:"required" json:"date"`
	Description string               `validate:"required,min=3" json:"description"`
	Lines       []JournalLineRequest `validate:"required,min=2,dive" json:"lines"`
	FileIds     []uuid.UUID          `json:"file_ids"`
}

type JournalEntryUpdateRequest struct {
	Id          uuid.UUID            `json:"-"`
	Date        time.Time            `validate:"required" json:"date"`
	Description string               `validate:"required,min=3" json:"description"`
	Lines       []JournalLineRequest `validate:"required,min=2,dive" json:"lines"`
	FileIds     []uuid.UUID          `json:"file_ids"`
}
