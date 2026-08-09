package repository

import (
	"fmt"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IInvoiceRepository interface {
	Create(invoice *entity.Invoice) error
	FindById(companyID, id uuid.UUID) (entity.Invoice, error)
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.Invoice, int64, error)
	Update(invoice *entity.Invoice) error
	Delete(companyID, id uuid.UUID) error
	BulkDestroy(companyID uuid.UUID, ids []uuid.UUID) error
	UpdateStatus(companyID, id uuid.UUID, status entity.InvoiceStatus, journalEntryId *uuid.UUID) error
	AddPayment(payment *entity.InvoicePayment) error
	UpdatePaymentStatus(companyID, id uuid.UUID, amountPaid float64, status entity.PaymentStatus) error
}

type InvoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) IInvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (r *InvoiceRepository) Create(invoice *entity.Invoice) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(invoice).Error
	})
}

func (r *InvoiceRepository) FindById(companyID, id uuid.UUID) (entity.Invoice, error) {
	var invoice entity.Invoice
	err := r.db.Preload("Customer").Preload("Items").Preload("Payments").Preload("Payments.PaymentAccount").
		Where("id = ? AND company_id = ?", id, companyID).
		First(&invoice).Error
	return invoice, err
}

var invoiceSortColumns = map[string]string{
	"invoice_number": "invoices.invoice_number",
	"invoice_date":   "invoices.invoice_date",
	"due_date":       "invoices.due_date",
	"invoice_status": "invoices.invoice_status",
	"created_at":     "invoices.created_at",
	"updated_at":     "invoices.updated_at",
}

func (r *InvoiceRepository) FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.Invoice, int64, error) {
	var invoices []entity.Invoice
	var total int64
	query := r.db.Model(&entity.Invoice{}).Where("company_id = ?", companyID)

	if qp.Search != "" {
		searchLike := "%" + qp.Search + "%"
		query = query.Where("invoice_number ILIKE ? OR notes ILIKE ?", searchLike, searchLike)
	}

	if statusValues, ok := qp.Filters["invoice_status"]; ok && len(statusValues) > 0 && statusValues[0] != "" {
		query = query.Where("invoice_status = ?", statusValues[0])
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sort := "created_at"
	order := "desc"
	if qp.SortBy != "" {
		if col, valid := invoiceSortColumns[qp.SortBy]; valid {
			sort = col
		}
	}
	if qp.SortDir == "asc" || qp.SortDir == "desc" {
		order = qp.SortDir
	}
	query = query.Order(fmt.Sprintf("%s %s", sort, order))

	query = util.ApplyPagination(query, qp)

	err := query.Preload("Customer").Find(&invoices).Error
	return invoices, total, err
}

func (r *InvoiceRepository) Update(invoice *entity.Invoice) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("invoice_id = ?", invoice.Id).Delete(&entity.InvoiceItem{}).Error; err != nil {
			return err
		}
		return tx.Save(invoice).Error
	})
}

func (r *InvoiceRepository) Delete(companyID, id uuid.UUID) error {
	return r.db.Where("id = ? AND company_id = ? AND invoice_status = ?", id, companyID, entity.InvoiceStatusDraft).Delete(&entity.Invoice{}).Error
}

func (r *InvoiceRepository) BulkDestroy(companyID uuid.UUID, ids []uuid.UUID) error {
	return r.db.Where("company_id = ? AND id IN ? AND invoice_status = ?", companyID, ids, entity.InvoiceStatusDraft).Delete(&entity.Invoice{}).Error
}

func (r *InvoiceRepository) UpdateStatus(companyID, id uuid.UUID, status entity.InvoiceStatus, journalEntryId *uuid.UUID) error {
	updates := map[string]interface{}{
		"invoice_status": status,
	}
	if journalEntryId != nil {
		updates["journal_entry_id"] = *journalEntryId
	}
	return r.db.Model(&entity.Invoice{}).
		Where("id = ? AND company_id = ?", id, companyID).
		Updates(updates).Error
}

func (r *InvoiceRepository) AddPayment(payment *entity.InvoicePayment) error {
	return r.db.Create(payment).Error
}

func (r *InvoiceRepository) UpdatePaymentStatus(companyID, id uuid.UUID, amountPaid float64, status entity.PaymentStatus) error {
	return r.db.Model(&entity.Invoice{}).
		Where("id = ? AND company_id = ?", id, companyID).
		Updates(map[string]interface{}{
			"amount_paid":    amountPaid,
			"amount_due":     gorm.Expr("grand_total - ?", amountPaid),
			"payment_status": status,
		}).Error
}
