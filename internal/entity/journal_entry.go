package entity

import (
	"time"

	"github.com/google/uuid"
)

type JournalEntryStatus string

const (
	JournalStatusDraft  JournalEntryStatus = "draft"
	JournalStatusPosted JournalEntryStatus = "posted"
	JournalStatusVoid   JournalEntryStatus = "void"
)

type JournalEntryType string

const (
	JournalTypeGeneral    JournalEntryType = "general"
	JournalTypeRevenue    JournalEntryType = "revenue"
	JournalTypeExpense    JournalEntryType = "expense"
	JournalTypePayable    JournalEntryType = "payable"
	JournalTypeReceivable JournalEntryType = "receivable"
	JournalTypePayment    JournalEntryType = "payment"
	JournalTypeClosing    JournalEntryType = "closing"
	JournalTypeOpening    JournalEntryType = "opening"
)

type JournalEntry struct {
	Id             uuid.UUID          `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	CompanyId      uuid.UUID          `gorm:"type:uuid;not null;index;" json:"company_id"`
	FiscalPeriodId *uuid.UUID         `gorm:"type:uuid;index;" json:"fiscal_period_id"` // Terisi saat Posted
	JournalNumber  string             `gorm:"type:varchar;not null;index;" json:"journal_number"`
	Type           JournalEntryType   `gorm:"type:varchar;not null;default:'general';" json:"type"`
	Date           time.Time          `gorm:"type:date;not null;" json:"date"`
	Description    string             `gorm:"type:text;not null;" json:"description"`
	Status         JournalEntryStatus `gorm:"type:varchar;not null;default:'draft';" json:"status"`
	TotalDebit     float64            `gorm:"type:numeric(20,4);not null;default:0;" json:"total_debit"`
	TotalCredit    float64            `gorm:"type:numeric(20,4);not null;default:0;" json:"total_credit"`
	CreatedBy      uuid.UUID          `gorm:"type:uuid;" json:"created_by"`
	CreatedAt      time.Time          `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt      time.Time          `gorm:"autoUpdateTime;" json:"updated_at"`

	Lines []JournalLine `gorm:"foreignKey:JournalEntryId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"lines"`
	Files []File        `gorm:"many2many:journal_entry_files;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"files,omitempty"`
}

func (JournalEntry) TableName() string { return "journal_entries" }

func (JournalEntry) SearchableFields() []string { return []string{"journal_number", "description"} }
