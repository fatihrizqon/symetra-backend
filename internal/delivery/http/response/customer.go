package response

import (
	"time"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/google/uuid"
)

type CustomerResponse struct {
	Id        uuid.UUID    `json:"id"`
	CompanyId uuid.UUID    `json:"company_id"`
	Code      string       `json:"code"`
	Name      string       `json:"name"`
	Email     string       `json:"email"`
	Phone     string       `json:"phone"`
	Address   string       `json:"address"`
	CoaId     *uuid.UUID   `json:"coa_id,omitempty"`
	Coa       *COAResponse `json:"coa,omitempty"`
	Status    int          `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func NewCustomerResponse(e entity.Customer) CustomerResponse {
	resp := CustomerResponse{
		Id:        e.Id,
		CompanyId: e.CompanyId,
		Code:      e.Code,
		Name:      e.Name,
		Email:     e.Email,
		Phone:     e.Phone,
		Address:   e.Address,
		CoaId:     e.CoaId,
		Status:    e.Status,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}

	if e.Coa != nil {
		coaResp := NewCOAResponse(*e.Coa)
		resp.Coa = &coaResp
	}

	return resp
}

func NewCustomerResponses(entities []entity.Customer) []CustomerResponse {
	var responses = make([]CustomerResponse, 0, len(entities))
	for _, e := range entities {
		responses = append(responses, NewCustomerResponse(e))
	}
	return responses
}
