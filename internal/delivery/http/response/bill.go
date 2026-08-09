package response

import (
	"time"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/google/uuid"
)

type BillItemResponse struct {
	Id            uuid.UUID    `json:"id"`
	Description   string       `json:"description"`
	Qty           float64      `json:"qty"`
	Price         float64      `json:"price"`
	Discount      float64      `json:"discount"`
	TaxApplicable bool         `json:"tax_applicable"`
	Amount        float64      `json:"amount"`
	AccountId     *uuid.UUID   `json:"account_id,omitempty"`
	Account       *COAResponse `json:"account,omitempty"`
}

func NewBillItemResponse(e entity.BillItem) BillItemResponse {
	resp := BillItemResponse{
		Id:            e.Id,
		Description:   e.Description,
		Qty:           e.Qty,
		Price:         e.Price,
		Discount:      e.Discount,
		TaxApplicable: e.TaxApplicable,
		Amount:        e.Amount,
		AccountId:     e.AccountId,
	}
	if e.Account != nil {
		accountResp := NewCOAResponse(*e.Account)
		resp.Account = &accountResp
	}
	return resp
}

type BillPaymentResponse struct {
	Id               uuid.UUID    `json:"id"`
	Amount           float64      `json:"amount"`
	PaymentDate      time.Time    `json:"payment_date"`
	PaymentAccountId uuid.UUID    `json:"payment_account_id"`
	PaymentAccount   *COAResponse `json:"payment_account,omitempty"`
	JournalEntryId   *uuid.UUID   `json:"journal_entry_id,omitempty"`
	Notes            string       `json:"notes"`
	CreatedBy        uuid.UUID    `json:"created_by"`
	CreatedAt        time.Time    `json:"created_at"`
}

func NewBillPaymentResponse(e entity.BillPayment) BillPaymentResponse {
	resp := BillPaymentResponse{
		Id:               e.Id,
		Amount:           e.Amount,
		PaymentDate:      e.PaymentDate,
		PaymentAccountId: e.PaymentAccountId,
		JournalEntryId:   e.JournalEntryId,
		Notes:            e.Notes,
		CreatedBy:        e.CreatedBy,
		CreatedAt:        e.CreatedAt,
	}
	if e.PaymentAccount != nil {
		accountResp := NewCOAResponse(*e.PaymentAccount)
		resp.PaymentAccount = &accountResp
	}
	return resp
}

type BillResponse struct {
	Id              uuid.UUID             `json:"id"`
	CompanyId       uuid.UUID             `json:"company_id"`
	BillNumber      string                `json:"bill_number"`
	PurchaseOrderId *uuid.UUID            `json:"purchase_order_id,omitempty"`
	VendorId        uuid.UUID             `json:"vendor_id"`
	Vendor          *VendorResponse       `json:"vendor,omitempty"`
	BillDate        time.Time             `json:"bill_date"`
	DueDate         time.Time             `json:"due_date"`
	Subtotal        float64               `json:"subtotal"`
	DiscountTotal   float64               `json:"discount_total"`
	Dpp             float64               `json:"dpp"`
	TaxRate         float64               `json:"tax_rate"`
	TaxAmount       float64               `json:"tax_amount"`
	GrandTotal      float64               `json:"grand_total"`
	AmountPaid      float64               `json:"amount_paid"`
	AmountDue       float64               `json:"amount_due"`
	BillStatus      string                `json:"bill_status"`
	PaymentStatus   string                `json:"payment_status"`
	Notes           string                `json:"notes"`
	JournalEntryId  *uuid.UUID            `json:"journal_entry_id,omitempty"`
	CreatedBy       uuid.UUID             `json:"created_by"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
	Items           []BillItemResponse    `json:"items,omitempty"`
	Payments        []BillPaymentResponse `json:"payments,omitempty"`
}

func NewBillResponse(e entity.Bill) BillResponse {
	resp := BillResponse{
		Id:              e.Id,
		CompanyId:       e.CompanyId,
		BillNumber:      e.BillNumber,
		PurchaseOrderId: e.PurchaseOrderId,
		VendorId:        e.VendorId,
		BillDate:        e.BillDate,
		DueDate:         e.DueDate,
		Subtotal:        e.Subtotal,
		DiscountTotal:   e.DiscountTotal,
		Dpp:             e.Dpp,
		TaxRate:         e.TaxRate,
		TaxAmount:       e.TaxAmount,
		GrandTotal:      e.GrandTotal,
		AmountPaid:      e.AmountPaid,
		AmountDue:       e.AmountDue,
		BillStatus:      string(e.BillStatus),
		PaymentStatus:   string(e.PaymentStatus),
		Notes:           e.Notes,
		JournalEntryId:  e.JournalEntryId,
		CreatedBy:       e.CreatedBy,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}

	if e.Vendor != nil {
		vendorResp := NewVendorResponse(*e.Vendor)
		resp.Vendor = &vendorResp
	}

	var items []BillItemResponse
	for _, item := range e.Items {
		items = append(items, NewBillItemResponse(item))
	}
	resp.Items = items

	var payments []BillPaymentResponse
	for _, p := range e.Payments {
		payments = append(payments, NewBillPaymentResponse(p))
	}
	resp.Payments = payments

	return resp
}

func NewBillResponses(entities []entity.Bill) []BillResponse {
	var responses = make([]BillResponse, 0, len(entities))
	for _, e := range entities {
		responses = append(responses, NewBillResponse(e))
	}
	return responses
}
