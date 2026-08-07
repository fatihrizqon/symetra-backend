package database

import (
	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	_ = db.AutoMigrate(
		// ── Auth & Users ──────────────────────────────────────────────────────
		&entity.User{},
		&entity.Session{},
		&entity.Credential{},
		&entity.Role{},
		&entity.Permission{},
		&entity.File{},
		&entity.RedisJob{},
		&entity.Company{},
		&entity.CompanyMember{},
		// ── COA ───────────────────────────────────────────────────────────────
		&entity.COAGroup{},
		&entity.COASubGroup{},
		&entity.COA{},
		// ── Phase 1 ───────────────────────────────────────────────────────────
		&entity.CompanyConfiguration{},
		&entity.FiscalYear{},
		&entity.FiscalPeriod{},
		&entity.FiscalPeriodLog{},
		&entity.JournalEntry{},
		&entity.JournalLine{},
		// ── Procurement ───────────────────────────────────────────────────────
		&entity.Vendor{},
		&entity.PurchaseOrder{},
		&entity.PurchaseOrderItem{},
		&entity.Bill{},
		&entity.BillItem{},
		&entity.BillPayment{},
		// ── Sales ─────────────────────────────────────────────────────────────
		&entity.Customer{},
		&entity.Quotation{},
		&entity.QuotationItem{},
		&entity.Invoice{},
		&entity.InvoiceItem{},
		&entity.InvoicePayment{},
	)

	// Performance indexes
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_redis_jobs_pending ON redis_jobs (created_at) WHERE status = 'PENDING';`)
	db.Exec(`CREATE EXTENSION IF NOT EXISTS pg_trgm;`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_username_trgm ON users USING gin (username gin_trgm_ops);`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_email_trgm ON users USING gin (email gin_trgm_ops);`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_companies_name_trgm ON companies USING gin (name gin_trgm_ops);`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_coa_name_trgm ON coa USING gin (name gin_trgm_ops);`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_coa_code_trgm ON coa USING gin (code gin_trgm_ops);`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_coa_subgroups_name_trgm ON coa_subgroups USING gin (name gin_trgm_ops);`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_coa_groups_name_trgm ON coa_groups USING gin (name gin_trgm_ops);`)

	seedDefaultData(db)
}

func seedDefaultData(db *gorm.DB) {
	SeedData(db)
}
