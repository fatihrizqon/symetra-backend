package config

import (
	"fmt"

	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"

	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/handler"
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/middleware"
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/route"
	"github.com/fatihrizqon/gofiber-microservice/internal/rbac"
	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/fatihrizqon/gofiber-microservice/internal/service"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/fatihrizqon/gofiber-microservice/internal/util/storage"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	Cors       *CORSConfig
	DB         *gorm.DB
	App        *fiber.App
	Log        *logrus.Logger
	Validate   *validator.Validate
	Config     *viper.Viper
	JWT        *JWTService
	Redis      *redis.Client
	Production bool
}

func Bootstrap(deps *BootstrapConfig) {
	cfg := deps.Config

	deps.App.Use(recover.New())
	deps.App.Use(compress.New())
	deps.App.Use(deps.Cors.Handler())

	// ── Static Files ─────────────────────────────────────────────────────────
	deps.App.Get("/uploads*", static.New("./public/uploads"))

	// ── Storage ──────────────────────────────────────────────────────────────
	baseURL := fmt.Sprintf("%s/uploads", cfg.GetString("web.base_url"))
	localStorage := storage.NewLocalStorage("./public/uploads", baseURL)

	appURL := cfg.GetString("web.base_url")

	// ── Repositories ──────────────────────────────────────────────────────────
	tokenRepository := repository.NewTokenRepository(deps.DB, deps.Redis)
	fileRepository := repository.NewFileRepository(deps.DB)
	authRepository := repository.NewAuthRepository(deps.DB)
	rbacRepository := repository.NewRbacRepository(deps.DB, deps.Redis)
	userRepository := repository.NewUserRepository(deps.DB)
	companyRepository := repository.NewCompanyRepository(deps.DB)
	coaGroupRepository := repository.NewCOAGroupRepository(deps.DB)
	coaSubGroupRepository := repository.NewCOASubGroupRepository(deps.DB)
	coaRepository := repository.NewCOARepository(deps.DB)
	companyConfigurationRepository := repository.NewCompanyConfigurationRepository(deps.DB)
	fiscalYearRepository := repository.NewFiscalYearRepository(deps.DB)
	journalEntryRepository := repository.NewJournalEntryRepository(deps.DB)
	reportRepository := repository.NewReportRepository(deps.DB)
	vendorRepository := repository.NewVendorRepository(deps.DB)
	customerRepository := repository.NewCustomerRepository(deps.DB)
	purchaseOrderRepository := repository.NewPurchaseOrderRepository(deps.DB)
	billRepository := repository.NewBillRepository(deps.DB)
	quotationRepository := repository.NewQuotationRepository(deps.DB)
	invoiceRepository := repository.NewInvoiceRepository(deps.DB)

	// ── Services ──────────────────────────────────────────────────────────────
	emailService := service.NewEmailService(deps.Log)
	fileService := service.NewFileService(fileRepository, localStorage)
	authService := service.NewAuthService(authRepository, tokenRepository, deps.Validate, emailService, appURL)
	rbacService := service.NewRbacService(rbacRepository, userRepository, deps.Validate)
	userService := service.NewUserService(deps.Validate, fileService, userRepository)
	companyService := service.NewCompanyService(deps.Validate, companyRepository, tokenRepository)
	coaGroupService := service.NewCOAGroupService(deps.Validate, coaGroupRepository)
	coaSubGroupService := service.NewCOASubGroupService(deps.Validate, coaSubGroupRepository)
	coaService := service.NewCOAService(deps.Validate, coaRepository)
	companyConfigurationService := service.NewCompanyConfigurationService(deps.Validate, companyConfigurationRepository)
	fiscalYearService := service.NewFiscalYearService(deps.Validate, fiscalYearRepository)
	journalEntryService := service.NewJournalEntryService(journalEntryRepository, coaRepository, fiscalYearRepository)
	reportService := service.NewReportService(reportRepository, companyConfigurationRepository)
	vendorService := service.NewVendorService(deps.Validate, vendorRepository)
	customerService := service.NewCustomerService(deps.Validate, customerRepository)
	purchaseOrderService := service.NewPurchaseOrderService(purchaseOrderRepository, vendorRepository)
	billService := service.NewBillService(billRepository, vendorRepository, companyConfigurationRepository, journalEntryService, deps.Validate)
	quotationService := service.NewQuotationService(quotationRepository, customerRepository)
	invoiceService := service.NewInvoiceService(deps.Validate, invoiceRepository, customerRepository, companyConfigurationRepository, journalEntryService)

	// ── Handlers ──────────────────────────────────────────────────────────────
	fileHandler := handler.NewFileHandler(fileService)
	authHandler := handler.NewAuthHandler(authService, deps.Production)
	rbacHandler := handler.NewRbacHandler(rbacService)
	userHandler := handler.NewUserHandler(userService)
	companyHandler := handler.NewCompanyHandler(companyService)
	coaGroupHandler := handler.NewCOAGroupHandler(coaGroupService)
	coaSubGroupHandler := handler.NewCOASubGroupHandler(coaSubGroupService)
	coaHandler := handler.NewCOAHandler(coaService)
	companyConfigurationHandler := handler.NewCompanyConfigurationHandler(companyConfigurationService)
	fiscalYearHandler := handler.NewFiscalYearHandler(fiscalYearService)
	journalEntryHandler := handler.NewJournalEntryHandler(journalEntryService)
	reportHandler := handler.NewReportHandler(reportService)
	vendorHandler := handler.NewVendorHandler(vendorService)
	customerHandler := handler.NewCustomerHandler(customerService)
	purchaseOrderHandler := handler.NewPurchaseOrderHandler(purchaseOrderService)
	billHandler := handler.NewBillHandler(billService)
	quotationHandler := handler.NewQuotationHandler(quotationService)
	invoiceHandler := handler.NewInvoiceHandler(invoiceService)

	// ── Middleware ────────────────────────────────────────────────────────────
	authMiddleware := middleware.NewAuth(tokenRepository)
	companyMiddleware := middleware.NewCompany(companyRepository)

	rbacEngine := rbac.New(rbac.Config{
		Store: rbacRepository,
		UserLookup: func(c fiber.Ctx) string {
			claims, ok := c.Locals("auth").(*util.Claims)
			if !ok {
				return ""
			}
			return claims.UserID.String()
		},
	})

	// ── Routes ────────────────────────────────────────────────────────────────
	routeConfig := route.RouteConfig{
		App:                         deps.App,
		AuthMiddleware:              authMiddleware,
		CompanyMiddleware:           companyMiddleware,
		RbacEngine:                  rbacEngine,
		UserHandler:                 userHandler,
		AuthHandler:                 authHandler,
		FileHandler:                 fileHandler,
		RbacHandler:                 rbacHandler,
		Production:                  deps.Production,
		CompanyHandler:              companyHandler,
		COAGroupHandler:             coaGroupHandler,
		COASubGroupHandler:          coaSubGroupHandler,
		COAHandler:                  coaHandler,
		CompanyConfigurationHandler: companyConfigurationHandler,
		FiscalYearHandler:           fiscalYearHandler,
		JournalEntryHandler:         journalEntryHandler,
		ReportHandler:               reportHandler,
		VendorHandler:               vendorHandler,
		CustomerHandler:             customerHandler,
		PurchaseOrderHandler:        purchaseOrderHandler,
		BillHandler:                 billHandler,
		QuotationHandler:            quotationHandler,
		InvoiceHandler:              invoiceHandler,
	}

	routeConfig.Setup()
}
