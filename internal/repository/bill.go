package repository

import (
	"fmt"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IBillRepository interface {
	Create(bill *entity.Bill) error
	FindById(companyId, id uuid.UUID) (entity.Bill, error)
	FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]entity.Bill, int64, error)
	Update(bill *entity.Bill) error
	Delete(companyId, id uuid.UUID) error
	BulkDestroy(companyId uuid.UUID, ids []uuid.UUID) error
	UpdateStatus(companyId, id uuid.UUID, status entity.BillStatus, journalEntryId *uuid.UUID) error
	AddPayment(payment *entity.BillPayment) error
	UpdatePaymentStatus(companyId, id uuid.UUID, amountPaid float64, status entity.PaymentStatus) error
}

type BillRepository struct {
	db *gorm.DB
}

func NewBillRepository(db *gorm.DB) IBillRepository {
	return &BillRepository{db: db}
}

func (r *BillRepository) Create(bill *entity.Bill) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(bill).Error
	})
}

func (r *BillRepository) FindById(companyId, id uuid.UUID) (entity.Bill, error) {
	var bill entity.Bill
	err := r.db.Preload("Vendor").Preload("Items").Preload("Items.Account").Preload("Payments").Preload("Payments.PaymentAccount").
		Where("id = ? AND company_id = ?", id, companyId).
		First(&bill).Error
	return bill, err
}

var billSortColumns = map[string]string{
	"bill_number":    "bills.bill_number",
	"bill_date":      "bills.bill_date",
	"due_date":       "bills.due_date",
	"bill_status":    "bills.bill_status",
	"payment_status": "bills.payment_status",
	"created_at":     "bills.created_at",
	"updated_at":     "bills.updated_at",
}

func (r *BillRepository) FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]entity.Bill, int64, error) {
	var bills []entity.Bill
	var total int64
	query := r.db.Model(&entity.Bill{}).Where("company_id = ?", companyId)

	if qp.Search != "" {
		searchLike := "%" + qp.Search + "%"
		query = query.Where("bill_number ILIKE ? OR notes ILIKE ?", searchLike, searchLike)
	}

	if statusValues, ok := qp.Filters["bill_status"]; ok && len(statusValues) > 0 && statusValues[0] != "" {
		query = query.Where("bill_status = ?", statusValues[0])
	}
	if statusValues, ok := qp.Filters["payment_status"]; ok && len(statusValues) > 0 && statusValues[0] != "" {
		query = query.Where("payment_status = ?", statusValues[0])
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sort := "created_at"
	order := "desc"
	if qp.SortBy != "" {
		if col, valid := billSortColumns[qp.SortBy]; valid {
			sort = col
		}
	}
	if qp.SortDir == "asc" || qp.SortDir == "desc" {
		order = qp.SortDir
	}
	query = query.Order(fmt.Sprintf("%s %s", sort, order))

	query = util.ApplyPagination(query, qp)

	err := query.Preload("Vendor").Find(&bills).Error
	return bills, total, err
}

func (r *BillRepository) Update(bill *entity.Bill) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("bill_id = ?", bill.Id).Delete(&entity.BillItem{}).Error; err != nil {
			return err
		}
		return tx.Save(bill).Error
	})
}

func (r *BillRepository) Delete(companyId, id uuid.UUID) error {
	return r.db.Where("id = ? AND company_id = ? AND bill_status = ?", id, companyId, entity.BillStatusDraft).Delete(&entity.Bill{}).Error
}

func (r *BillRepository) BulkDestroy(companyId uuid.UUID, ids []uuid.UUID) error {
	return r.db.Where("company_id = ? AND id IN ? AND bill_status = ?", companyId, ids, entity.BillStatusDraft).Delete(&entity.Bill{}).Error
}

func (r *BillRepository) UpdateStatus(companyId, id uuid.UUID, status entity.BillStatus, journalEntryId *uuid.UUID) error {
	updates := map[string]interface{}{
		"bill_status": status,
	}
	if journalEntryId != nil {
		updates["journal_entry_id"] = *journalEntryId
	}
	return r.db.Model(&entity.Bill{}).
		Where("id = ? AND company_id = ?", id, companyId).
		Updates(updates).Error
}

func (r *BillRepository) AddPayment(payment *entity.BillPayment) error {
	return r.db.Create(payment).Error
}

func (r *BillRepository) UpdatePaymentStatus(companyId, id uuid.UUID, amountPaid float64, status entity.PaymentStatus) error {
	return r.db.Model(&entity.Bill{}).
		Where("id = ? AND company_id = ?", id, companyId).
		Updates(map[string]interface{}{
			"amount_paid":    amountPaid,
			"amount_due":     gorm.Expr("grand_total - ?", amountPaid),
			"payment_status": status,
		}).Error
}
