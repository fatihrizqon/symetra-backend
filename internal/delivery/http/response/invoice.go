package response

import (
	"time"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/google/uuid"
)

type InvoiceItemResponse struct {
	Id            uuid.UUID `json:"id"`
	Description   string    `json:"description"`
	Qty           float64   `json:"qty"`
	Price         float64   `json:"price"`
	Discount      float64   `json:"discount"`
	TaxApplicable bool      `json:"tax_applicable"`
	Amount        float64   `json:"amount"`
}

type InvoiceResponse struct {
	Id               uuid.UUID             `json:"id"`
	CompanyId        uuid.UUID             `json:"company_id"`
	InvoiceNumber    string                `json:"invoice_number"`
	QuotationId      *uuid.UUID            `json:"quotation_id,omitempty"`
	CustomerId       uuid.UUID             `json:"customer_id"`
	Customer         *CustomerResponse     `json:"customer,omitempty"`
	InvoiceDate      time.Time             `json:"invoice_date"`
	DueDate          time.Time             `json:"due_date"`
	Subtotal         float64               `json:"subtotal"`
	DiscountTotal    float64               `json:"discount_total"`
	Dpp              float64               `json:"dpp"`
	TaxRate          float64               `json:"tax_rate"`
	TaxAmount        float64               `json:"tax_amount"`
	GrandTotal       float64               `json:"grand_total"`
	AmountPaid       float64               `json:"amount_paid"`
	AmountDue        float64               `json:"amount_due"`
	InvoiceStatus    string                `json:"invoice_status"`
	PaymentStatus    string                `json:"payment_status"`
	TaxStatus        string                `json:"tax_status"`
	Notes            string                `json:"notes"`
	JournalEntryId   *uuid.UUID            `json:"journal_entry_id,omitempty"`
	PaymentJournalId *uuid.UUID            `json:"payment_journal_id,omitempty"`
	CreatedBy        uuid.UUID             `json:"created_by"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
	Items            []InvoiceItemResponse `json:"items,omitempty"`
}

// FromInvoiceEntity converts an entity.Invoice to an InvoiceResponse.
func FromInvoiceEntity(inv entity.Invoice) InvoiceResponse {
	var items []InvoiceItemResponse
	for _, item := range inv.Items {
		items = append(items, InvoiceItemResponse{
			Id:            item.Id,
			Description:   item.Description,
			Qty:           item.Qty,
			Price:         item.Price,
			Discount:      item.Discount,
			TaxApplicable: item.TaxApplicable,
			Amount:        item.Amount,
		})
	}

	resp := InvoiceResponse{
		Id:               inv.Id,
		CompanyId:        inv.CompanyId,
		InvoiceNumber:    inv.InvoiceNumber,
		CustomerId:       inv.CustomerId,
		QuotationId:      inv.QuotationId,
		InvoiceDate:      inv.InvoiceDate,
		DueDate:          inv.DueDate,
		Subtotal:         inv.Subtotal,
		DiscountTotal:    inv.DiscountTotal,
		Dpp:              inv.Dpp,
		TaxRate:          inv.TaxRate,
		TaxAmount:        inv.TaxAmount,
		GrandTotal:       inv.GrandTotal,
		AmountPaid:       inv.AmountPaid,
		AmountDue:        inv.AmountDue,
		InvoiceStatus:    string(inv.InvoiceStatus),
		PaymentStatus:    string(inv.PaymentStatus),
		TaxStatus:        string(inv.TaxStatus),
		Notes:            inv.Notes,
		JournalEntryId:   inv.JournalEntryId,
		PaymentJournalId: inv.PaymentJournalId,
		CreatedBy:        inv.CreatedBy,
		CreatedAt:        inv.CreatedAt,
		UpdatedAt:        inv.UpdatedAt,
		Items:            items,
	}

	if inv.Customer != nil {
		custResp := CustomerResponse{
			Id:        inv.Customer.Id,
			CompanyId: inv.Customer.CompanyId,
			Code:      inv.Customer.Code,
			Name:      inv.Customer.Name,
			Email:     inv.Customer.Email,
			Phone:     inv.Customer.Phone,
			Address:   inv.Customer.Address,
			CoaId:     inv.Customer.CoaId,
			Status:    inv.Customer.Status,
			CreatedAt: inv.Customer.CreatedAt,
			UpdatedAt: inv.Customer.UpdatedAt,
		}
		if inv.Customer.Coa != nil {
			custResp.Coa = &COAResponse{
				Id:            inv.Customer.Coa.Id,
				Code:          inv.Customer.Coa.Code,
				Name:          inv.Customer.Coa.Name,
				IsContra:      inv.Customer.Coa.IsContra,
				NormalBalance: inv.Customer.Coa.GetAbsoluteNormalBalance(),
			}
		}
		resp.Customer = &custResp
	}

	return resp
}
