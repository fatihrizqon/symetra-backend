package repository

import (
	"fmt"
	"time"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IJournalEntryRepository interface {
	Create(journal *entity.JournalEntry) error
	FindById(companyId, id uuid.UUID) (entity.JournalEntry, error)
	FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]entity.JournalEntry, int64, error)
	Update(journal *entity.JournalEntry) error
	Delete(companyId, id uuid.UUID) error
	Post(companyId, id, fiscalPeriodId uuid.UUID) error
	Void(companyId, id uuid.UUID) error
	GenerateJournalNumber(companyId uuid.UUID, date time.Time, typePrefix string) (string, error)
	BulkDestroy(companyId uuid.UUID, ids []uuid.UUID) error
}

type JournalEntryRepository struct {
	db *gorm.DB
}

func NewJournalEntryRepository(db *gorm.DB) IJournalEntryRepository {
	return &JournalEntryRepository{db: db}
}

func (r *JournalEntryRepository) GenerateJournalNumber(companyId uuid.UUID, date time.Time, typePrefix string) (string, error) {
	if typePrefix == "" {
		typePrefix = "JE"
	}
	monthStr := date.Format("200601")
	prefix := fmt.Sprintf("%s-%s-", typePrefix, monthStr)

	var lastJournal entity.JournalEntry
	err := r.db.Where("company_id = ? AND journal_number LIKE ?", companyId, prefix+"%").
		Order("journal_number desc").
		First(&lastJournal).Error

	var sequence int
	if err == gorm.ErrRecordNotFound {
		sequence = 1
	} else if err != nil {
		return "", err
	} else {
		// parse sequence from lastJournal.JournalNumber e.g. JE-202501-0001
		fmt.Sscanf(lastJournal.JournalNumber, prefix+"%04d", &sequence)
		sequence++
	}

	return fmt.Sprintf("%s%04d", prefix, sequence), nil
}

func (r *JournalEntryRepository) Create(journal *entity.JournalEntry) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(journal).Error
	})
}

func (r *JournalEntryRepository) FindById(companyId, id uuid.UUID) (entity.JournalEntry, error) {
	var journal entity.JournalEntry
	err := r.db.Preload("Lines").Preload("Lines.Coa").Preload("Files").
		Where("id = ? AND company_id = ?", id, companyId).
		First(&journal).Error
	return journal, err
}

var journalSortColumns = map[string]string{
	"journal_number": "journal_entries.journal_number",
	"date":           "journal_entries.date",
	"status":         "journal_entries.status",
	"created_at":     "journal_entries.created_at",
	"updated_at":     "journal_entries.updated_at",
}

func (r *JournalEntryRepository) FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]entity.JournalEntry, int64, error) {
	var journals []entity.JournalEntry
	var total int64
	query := r.db.Model(&entity.JournalEntry{}).Where("company_id = ?", companyId)

	if qp.Search != "" {
		searchLike := "%" + qp.Search + "%"
		query = query.Where(
			"journal_number ILIKE ? OR description ILIKE ?",
			searchLike, searchLike,
		)
	}
	if statusValues, ok := qp.Filters["status"]; ok && len(statusValues) > 0 && statusValues[0] != "" {
		query = query.Where("status = ?", statusValues[0])
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sort := "created_at"
	order := "desc"
	if qp.SortBy != "" {
		if col, valid := journalSortColumns[qp.SortBy]; valid {
			sort = col
		}
	}
	if qp.SortDir == "asc" || qp.SortDir == "desc" {
		order = qp.SortDir
	}
	query = query.Order(fmt.Sprintf("%s %s", sort, order))

	query = util.ApplyPagination(query, qp)

	err := query.Preload("Lines").Preload("Lines.Coa").Preload("Files").Find(&journals).Error
	return journals, total, err
}

func (r *JournalEntryRepository) Update(journal *entity.JournalEntry) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("journal_entry_id = ?", journal.Id).Delete(&entity.JournalLine{}).Error; err != nil {
			return err
		}
		if err := tx.Model(journal).Association("Files").Replace(journal.Files); err != nil {
			return err
		}
		return tx.Save(journal).Error
	})
}

func (r *JournalEntryRepository) Delete(companyId, id uuid.UUID) error {
	return r.db.Where("id = ? AND company_id = ? AND status = ?", id, companyId, entity.JournalStatusDraft).Delete(&entity.JournalEntry{}).Error
}

func (r *JournalEntryRepository) Post(companyId, id, fiscalPeriodId uuid.UUID) error {
	return r.db.Model(&entity.JournalEntry{}).
		Where("id = ? AND company_id = ? AND status = ?", id, companyId, entity.JournalStatusDraft).
		Updates(map[string]interface{}{
			"status":           entity.JournalStatusPosted,
			"fiscal_period_id": fiscalPeriodId,
		}).Error
}

func (r *JournalEntryRepository) Void(companyId, id uuid.UUID) error {
	return r.db.Model(&entity.JournalEntry{}).
		Where("id = ? AND company_id = ? AND status = ?", id, companyId, entity.JournalStatusPosted).
		Update("status", entity.JournalStatusVoid).Error
}

func (r *JournalEntryRepository) BulkDestroy(companyId uuid.UUID, ids []uuid.UUID) error {
	return r.db.Where("company_id = ? AND id IN ? AND status = ?", companyId, ids, entity.JournalStatusDraft).Delete(&entity.JournalEntry{}).Error
}

