package database

import (
	"log"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) {
	log.Println("[SEED] Starting seed data...")

	seedAdminUser(db)

	log.Println("[SEED] Seed complete.")
}

func seedAdminUser(db *gorm.DB) {
	// 1. Seed Permissions
	var (
		pUserRead, pUserManage, pRbacManage      entity.Permission
		pCompanyRead, pCompanyManage             entity.Permission
		pCoaGroupRead, pCoaGroupManage           entity.Permission
		pCoaSubgroupRead, pCoaSubgroupManage     entity.Permission
		pCoaRead, pCoaManage                     entity.Permission
		pCompanyConfigRead, pCompanyConfigManage entity.Permission
		pFiscalYearRead, pFiscalYearManage       entity.Permission
		pTransactionsRead, pTransactionsManage   entity.Permission
		pReportsRead                             entity.Permission
		
		pPurchaseOrderRead, pPurchaseOrderManage entity.Permission
		pBillRead, pBillManage                   entity.Permission
		pQuotationRead, pQuotationManage         entity.Permission
		pInvoiceRead, pInvoiceManage             entity.Permission
		pCustomerRead, pCustomerManage           entity.Permission
		pVendorRead, pVendorManage               entity.Permission
	)

	if err := db.Where("name = ?", "users.read").First(&pUserRead).Error; err != nil {
		pUserRead = entity.Permission{Name: "users.read", Description: "Read user data"}
		db.Create(&pUserRead)
	}

	if err := db.Where("name = ?", "users.manage").First(&pUserManage).Error; err != nil {
		pUserManage = entity.Permission{Name: "users.manage", Description: "Manage user data"}
		db.Create(&pUserManage)
	}

	if err := db.Where("name = ?", "rbac.manage").First(&pRbacManage).Error; err != nil {
		pRbacManage = entity.Permission{Name: "rbac.manage", Description: "Manage RBAC (Roles and Permissions)"}
		db.Create(&pRbacManage)
	}

	if err := db.Where("name = ?", "companies.read").First(&pCompanyRead).Error; err != nil {
		pCompanyRead = entity.Permission{Name: "companies.read", Description: "Read company data"}
		db.Create(&pCompanyRead)
	}

	if err := db.Where("name = ?", "companies.manage").First(&pCompanyManage).Error; err != nil {
		pCompanyManage = entity.Permission{Name: "companies.manage", Description: "Manage company data"}
		db.Create(&pCompanyManage)
	}

	if err := db.Where("name = ?", "coa_groups.read").First(&pCoaGroupRead).Error; err != nil {
		pCoaGroupRead = entity.Permission{Name: "coa_groups.read", Description: "Read COA group data"}
		db.Create(&pCoaGroupRead)
	}

	if err := db.Where("name = ?", "coa_groups.manage").First(&pCoaGroupManage).Error; err != nil {
		pCoaGroupManage = entity.Permission{Name: "coa_groups.manage", Description: "Manage COA group data"}
		db.Create(&pCoaGroupManage)
	}

	if err := db.Where("name = ?", "coa_subgroups.read").First(&pCoaSubgroupRead).Error; err != nil {
		pCoaSubgroupRead = entity.Permission{Name: "coa_subgroups.read", Description: "Read COA subgroup data"}
		db.Create(&pCoaSubgroupRead)
	}

	if err := db.Where("name = ?", "coa_subgroups.manage").First(&pCoaSubgroupManage).Error; err != nil {
		pCoaSubgroupManage = entity.Permission{Name: "coa_subgroups.manage", Description: "Manage COA subgroup data"}
		db.Create(&pCoaSubgroupManage)
	}

	if err := db.Where("name = ?", "coa.read").First(&pCoaRead).Error; err != nil {
		pCoaRead = entity.Permission{Name: "coa.read", Description: "Read COA data"}
		db.Create(&pCoaRead)
	}

	if err := db.Where("name = ?", "coa.manage").First(&pCoaManage).Error; err != nil {
		pCoaManage = entity.Permission{Name: "coa.manage", Description: "Manage COA data"}
		db.Create(&pCoaManage) // Fixed previous typo here from pCoaSubgroupManage
	}

	if err := db.Where("name = ?", "company_config.read").First(&pCompanyConfigRead).Error; err != nil {
		pCompanyConfigRead = entity.Permission{Name: "company_config.read", Description: "Read company configuration"}
		db.Create(&pCompanyConfigRead)
	}

	if err := db.Where("name = ?", "company_config.manage").First(&pCompanyConfigManage).Error; err != nil {
		pCompanyConfigManage = entity.Permission{Name: "company_config.manage", Description: "Manage company configuration"}
		db.Create(&pCompanyConfigManage)
	}

	if err := db.Where("name = ?", "fiscal_years.read").First(&pFiscalYearRead).Error; err != nil {
		pFiscalYearRead = entity.Permission{Name: "fiscal_years.read", Description: "Read fiscal years data"}
		db.Create(&pFiscalYearRead)
	}

	if err := db.Where("name = ?", "fiscal_years.manage").First(&pFiscalYearManage).Error; err != nil {
		pFiscalYearManage = entity.Permission{Name: "fiscal_years.manage", Description: "Manage fiscal years data"}
		db.Create(&pFiscalYearManage)
	}

	if err := db.Where("name = ?", "transactions.read").First(&pTransactionsRead).Error; err != nil {
		pTransactionsRead = entity.Permission{Name: "transactions.read", Description: "Read journal entries and transactions"}
		db.Create(&pTransactionsRead)
	}

	if err := db.Where("name = ?", "transactions.manage").First(&pTransactionsManage).Error; err != nil {
		pTransactionsManage = entity.Permission{Name: "transactions.manage", Description: "Create, edit, post and void journal entries"}
		db.Create(&pTransactionsManage)
	}

	if err := db.Where("name = ?", "reports.read").First(&pReportsRead).Error; err != nil {
		pReportsRead = entity.Permission{Name: "reports.read", Description: "Read reports"}
		db.Create(&pReportsRead)
	}

	// Procurement Permissions
	if err := db.Where("name = ?", "purchase_orders.read").First(&pPurchaseOrderRead).Error; err != nil {
		pPurchaseOrderRead = entity.Permission{Name: "purchase_orders.read", Description: "Read purchase orders"}
		db.Create(&pPurchaseOrderRead)
	}
	if err := db.Where("name = ?", "purchase_orders.manage").First(&pPurchaseOrderManage).Error; err != nil {
		pPurchaseOrderManage = entity.Permission{Name: "purchase_orders.manage", Description: "Manage purchase orders"}
		db.Create(&pPurchaseOrderManage)
	}
	if err := db.Where("name = ?", "bills.read").First(&pBillRead).Error; err != nil {
		pBillRead = entity.Permission{Name: "bills.read", Description: "Read bills"}
		db.Create(&pBillRead)
	}
	if err := db.Where("name = ?", "bills.manage").First(&pBillManage).Error; err != nil {
		pBillManage = entity.Permission{Name: "bills.manage", Description: "Manage bills"}
		db.Create(&pBillManage)
	}
	if err := db.Where("name = ?", "vendors.read").First(&pVendorRead).Error; err != nil {
		pVendorRead = entity.Permission{Name: "vendors.read", Description: "Read vendors"}
		db.Create(&pVendorRead)
	}
	if err := db.Where("name = ?", "vendors.manage").First(&pVendorManage).Error; err != nil {
		pVendorManage = entity.Permission{Name: "vendors.manage", Description: "Manage vendors"}
		db.Create(&pVendorManage)
	}

	// Sales Permissions
	if err := db.Where("name = ?", "quotations.read").First(&pQuotationRead).Error; err != nil {
		pQuotationRead = entity.Permission{Name: "quotations.read", Description: "Read quotations"}
		db.Create(&pQuotationRead)
	}
	if err := db.Where("name = ?", "quotations.manage").First(&pQuotationManage).Error; err != nil {
		pQuotationManage = entity.Permission{Name: "quotations.manage", Description: "Manage quotations"}
		db.Create(&pQuotationManage)
	}
	if err := db.Where("name = ?", "invoices.read").First(&pInvoiceRead).Error; err != nil {
		pInvoiceRead = entity.Permission{Name: "invoices.read", Description: "Read invoices"}
		db.Create(&pInvoiceRead)
	}
	if err := db.Where("name = ?", "invoices.manage").First(&pInvoiceManage).Error; err != nil {
		pInvoiceManage = entity.Permission{Name: "invoices.manage", Description: "Manage invoices"}
		db.Create(&pInvoiceManage)
	}
	if err := db.Where("name = ?", "customers.read").First(&pCustomerRead).Error; err != nil {
		pCustomerRead = entity.Permission{Name: "customers.read", Description: "Read customers"}
		db.Create(&pCustomerRead)
	}
	if err := db.Where("name = ?", "customers.manage").First(&pCustomerManage).Error; err != nil {
		pCustomerManage = entity.Permission{Name: "customers.manage", Description: "Manage customers"}
		db.Create(&pCustomerManage)
	}

	// 2. Seed Role
	var roleAdmin entity.Role
	if err := db.Where("name = ?", "admin").First(&roleAdmin).Error; err != nil {
		roleAdmin = entity.Role{
			Name:        "admin",
			Description: "Administrator Role",
			Permissions: []entity.Permission{
				pUserRead, pUserManage, pRbacManage,
				pCompanyRead, pCompanyManage,
				pCoaGroupRead, pCoaGroupManage,
				pCoaSubgroupRead, pCoaSubgroupManage,
				pCoaRead, pCoaManage,
				pCompanyConfigRead, pCompanyConfigManage,
				pFiscalYearRead, pFiscalYearManage,
				pTransactionsRead, pTransactionsManage,
				pReportsRead,
				pPurchaseOrderRead, pPurchaseOrderManage,
				pBillRead, pBillManage,
				pQuotationRead, pQuotationManage,
				pInvoiceRead, pInvoiceManage,
				pCustomerRead, pCustomerManage,
				pVendorRead, pVendorManage,
			},
		}
		db.Create(&roleAdmin)
	} else {
		_ = db.Model(&roleAdmin).Association("Permissions").Append(
			&pRbacManage, &pCompanyConfigRead, &pCompanyConfigManage, 
			&pFiscalYearRead, &pFiscalYearManage, &pTransactionsRead, &pTransactionsManage, &pReportsRead,
			&pPurchaseOrderRead, &pPurchaseOrderManage,
			&pBillRead, &pBillManage,
			&pQuotationRead, &pQuotationManage,
			&pInvoiceRead, &pInvoiceManage,
			&pCustomerRead, &pCustomerManage,
			&pVendorRead, &pVendorManage,
		)
	}

	var user entity.User
	if db.Where("email = ?", "admin@example.com").First(&user).Error != nil {
		hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		db.Create(&entity.User{
			Username: "admin", Name: "System Administrator",
			Email:    "admin@example.com",
			Password: string(hash), Status: 1,
			Roles: []entity.Role{roleAdmin},
		})
		log.Println("[SEED] Admin user created: admin@example.com / password")
	} else {
		_ = db.Model(&user).Association("Roles").Append(&roleAdmin)
	}
	if db.Where("email = ?", "operator@example.com").First(&user).Error != nil {
		hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		db.Create(&entity.User{
			Username: "operator", Name: "System Operator",
			Email:    "operator@example.com",
			Password: string(hash), Status: 0,
		})
		log.Println("[SEED] operator user created: operator@example.com / password")
	}
}
