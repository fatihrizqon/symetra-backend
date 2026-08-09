package response

import (
	"time"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/google/uuid"
)

type JournalLineResponse struct {
	Id             uuid.UUID    `json:"id"`
	JournalEntryId uuid.UUID    `json:"journal_entry_id"`
	CoaId          uuid.UUID    `json:"coa_id"`
	Coa            *COAResponse `json:"coa,omitempty"`
	Description    *string      `json:"description"`
	Debit          float64      `json:"debit"`
	Credit         float64      `json:"credit"`
}

func NewJournalLineResponse(e entity.JournalLine) JournalLineResponse {
	resp := JournalLineResponse{
		Id:             e.Id,
		JournalEntryId: e.JournalEntryId,
		CoaId:          e.CoaId,
		Description:    e.Description,
		Debit:          e.Debit,
		Credit:         e.Credit,
	}

	if e.Coa != nil {
		coaResp := NewCOAResponse(*e.Coa)
		resp.Coa = &coaResp
	}

	return resp
}

type JournalEntryResponse struct {
	Id             uuid.UUID             `json:"id"`
	CompanyId      uuid.UUID             `json:"company_id"`
	FiscalPeriodId *uuid.UUID            `json:"fiscal_period_id,omitempty"`
	JournalNumber  string                `json:"journal_number"`
	Type           string                `json:"type"`
	Date           time.Time             `json:"date"`
	Description    string                `json:"description"`
	Status         string                `json:"status"`
	TotalDebit     float64               `json:"total_debit"`
	TotalCredit    float64               `json:"total_credit"`
	CreatedBy      uuid.UUID             `json:"created_by"`
	Lines          []JournalLineResponse `json:"lines,omitempty"`
	Files          []FileResponse        `json:"files,omitempty"`
}

func NewJournalEntryResponse(e entity.JournalEntry) JournalEntryResponse {
	resp := JournalEntryResponse{
		Id:             e.Id,
		CompanyId:      e.CompanyId,
		FiscalPeriodId: e.FiscalPeriodId,
		JournalNumber:  e.JournalNumber,
		Type:           string(e.Type),
		Date:           e.Date,
		Description:    e.Description,
		Status:         string(e.Status),
		TotalDebit:     e.TotalDebit,
		TotalCredit:    e.TotalCredit,
		CreatedBy:      e.CreatedBy,
	}

	lines := make([]JournalLineResponse, 0, len(e.Lines))
	for _, line := range e.Lines {
		lines = append(lines, NewJournalLineResponse(line))
	}
	resp.Lines = lines

	files := make([]FileResponse, 0, len(e.Files))
	for _, file := range e.Files {
		files = append(files, FileResponse{
			Id:           file.Id,
			OriginalName: file.OriginalName,
			MimeType:     file.MimeType,
			Size:         file.Size,
			URL:          file.Path,
			UploadedBy:   file.UploadedBy,
			CreatedAt:    file.CreatedAt,
		})
	}
	resp.Files = files

	return resp
}

func NewJournalEntryResponses(
	entities []entity.JournalEntry,
) []JournalEntryResponse {

	responses := make([]JournalEntryResponse, 0, len(entities))

	for _, e := range entities {
		responses = append(responses, NewJournalEntryResponse(e))
	}

	return responses
}
