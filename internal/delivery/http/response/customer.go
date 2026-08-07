package response

import (
	"time"
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


