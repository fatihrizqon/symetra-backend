package route

import (
	"time"

	_ "github.com/fatihrizqon/gofiber-microservice/docs"
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/handler"
	"github.com/fatihrizqon/gofiber-microservice/internal/rbac"
	swagger "github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

type RouteConfig struct {
	App                         *fiber.App
	AuthMiddleware              fiber.Handler
	CompanyMiddleware           fiber.Handler
	RbacEngine                  *rbac.RBAC
	UserHandler                 *handler.UserHandler
	AuthHandler                 *handler.AuthHandler
	FileHandler                 *handler.FileHandler
	RbacHandler                 *handler.RbacHandler
	Production                  bool
	CompanyHandler              *handler.CompanyHandler
	COAGroupHandler             *handler.COAGroupHandler
	COASubGroupHandler          *handler.COASubGroupHandler
	COAHandler                  *handler.COAHandler
	CompanyConfigurationHandler *handler.CompanyConfigurationHandler
	FiscalYearHandler           *handler.FiscalYearHandler
	JournalEntryHandler         *handler.JournalEntryHandler
}

func (rc *RouteConfig) Setup() {
	rc.SetupGuestRoute()
	rc.SetupAuthRoute()
}

func (rc *RouteConfig) SetupGuestRoute() {
	authLimiter := limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
	})

	rc.App.Post("/api/v1/auth/register", authLimiter, rc.AuthHandler.Register)
	rc.App.Post("/api/v1/auth/login", authLimiter, rc.AuthHandler.Login)
	rc.App.Post("/api/v1/auth/refresh", authLimiter, rc.AuthHandler.Refresh)
	rc.App.Get("/swagger/*", swagger.HandlerDefault)
}

func (rc *RouteConfig) SetupAuthRoute() {
	// Initialize AuthMiddleware
	rc.App.Use(rc.AuthMiddleware)

	rc.App.Post("/api/v1/auth/logout", rc.AuthHandler.Logout)

	// User Profile
	rc.App.Post("/api/v1/users/me/avatar", rc.UserHandler.UploadAvatar)

	rc.App.Get("/api/v1/users", rc.RbacEngine.Require("users.read"), rc.UserHandler.FindAll)
	rc.App.Get("/api/v1/users/:id", rc.RbacEngine.Require("users.read"), rc.UserHandler.FindById)
	rc.App.Put("/api/v1/users/:id", rc.RbacEngine.Require("users.manage"), rc.UserHandler.Update)
	rc.App.Delete("/api/v1/users/:id", rc.RbacEngine.Require("users.manage"), rc.UserHandler.Delete)
	rc.App.Patch("/api/v1/users/:id/lock", rc.RbacEngine.Require("users.manage"), rc.UserHandler.Lock)
	rc.App.Patch("/api/v1/users/:id/unlock", rc.RbacEngine.Require("users.manage"), rc.UserHandler.Unlock)

	rc.App.Post("/api/v1/files/upload", rc.FileHandler.Upload)

	// User - Role Management
	rc.App.Post("/api/v1/users/:id/roles", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.AssignRole)
	rc.App.Delete("/api/v1/users/:id/roles/:role_id", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.RevokeRole)

	// Role - Permission Management
	rc.App.Post("/api/v1/roles/:id/permissions", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.AssignPermission)
	rc.App.Delete("/api/v1/roles/:id/permissions/:permission_id", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.RevokePermission)

	// Role Management
	rc.App.Post("/api/v1/roles", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.CreateRole)
	rc.App.Get("/api/v1/roles", cache.New(cache.Config{Expiration: 5 * time.Minute}), rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.GetRoles)
	rc.App.Get("/api/v1/roles/:id", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.GetRoleById)
	rc.App.Put("/api/v1/roles/:id", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.UpdateRole)
	rc.App.Delete("/api/v1/roles/:id", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.DeleteRole)

	// Permission Management
	rc.App.Post("/api/v1/permissions", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.CreatePermission)
	rc.App.Get("/api/v1/permissions", cache.New(cache.Config{Expiration: 5 * time.Minute}), rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.GetPermissions)
	rc.App.Get("/api/v1/permissions/:id", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.GetPermissionById)
	rc.App.Put("/api/v1/permissions/:id", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.UpdatePermission)
	rc.App.Delete("/api/v1/permissions/:id", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.DeletePermission)

	// Company Management
	rc.App.Get("/api/v1/companies", rc.RbacEngine.Require("companies.read"), rc.CompanyHandler.FindAll)
	rc.App.Post("/api/v1/companies", rc.RbacEngine.Require("companies.manage"), rc.CompanyHandler.Create)
	rc.App.Get("/api/v1/companies/mine", rc.RbacEngine.Require("companies.read"), rc.CompanyHandler.FindMyCompanies)
	rc.App.Get("/api/v1/companies/active", rc.CompanyHandler.GetActiveCompany)
	rc.App.Get("/api/v1/companies/:id", rc.RbacEngine.Require("companies.read"), rc.CompanyHandler.FindById)
	rc.App.Put("/api/v1/companies/:id", rc.RbacEngine.Require("companies.manage"), rc.CompanyHandler.Update)
	rc.App.Post("/api/v1/companies/:id/select", rc.CompanyHandler.SelectCompany)
	rc.App.Delete("/api/v1/companies/:id", rc.RbacEngine.Require("companies.manage"), rc.CompanyHandler.Delete)

	// Company Members Management
	rc.App.Get("/api/v1/companies/:id/members", rc.RbacEngine.Require("companies.manage"), rc.CompanyHandler.FindMembersByCompany)
	rc.App.Post("/api/v1/companies/:id/members", rc.RbacEngine.Require("companies.manage"), rc.CompanyHandler.AssignMember)
	rc.App.Put("/api/v1/companies/:id/members/:user_id/role", rc.RbacEngine.Require("companies.manage"), rc.CompanyHandler.UpdateMemberRole)
	rc.App.Delete("/api/v1/companies/:id/members/:user_id", rc.RbacEngine.Require("companies.manage"), rc.CompanyHandler.RemoveMember)

	// COA Group Management
	rc.App.Get("/api/v1/coa-groups", rc.RbacEngine.Require("coa_groups.read"), rc.COAGroupHandler.FindAll)
	rc.App.Post("/api/v1/coa-groups", rc.RbacEngine.Require("coa_groups.manage"), rc.COAGroupHandler.Create)
	rc.App.Get("/api/v1/coa-groups/select", cache.New(cache.Config{Expiration: 5 * time.Minute}), rc.RbacEngine.Require("coa_groups.read"), rc.COAGroupHandler.SelectDropdownList)
	rc.App.Get("/api/v1/coa-groups/:id", rc.RbacEngine.Require("coa_groups.read"), rc.COAGroupHandler.FindById)
	rc.App.Put("/api/v1/coa-groups/:id", rc.RbacEngine.Require("coa_groups.manage"), rc.COAGroupHandler.Update)
	rc.App.Delete("/api/v1/coa-groups/:id", rc.RbacEngine.Require("coa_groups.manage"), rc.COAGroupHandler.Delete)

	// COA Subgroup Management
	rc.App.Get("/api/v1/coa-subgroups", rc.RbacEngine.Require("coa_subgroups.read"), rc.COASubGroupHandler.FindAll)
	rc.App.Post("/api/v1/coa-subgroups", rc.RbacEngine.Require("coa_subgroups.manage"), rc.COASubGroupHandler.Create)
	rc.App.Get("/api/v1/coa-subgroups/select", cache.New(cache.Config{Expiration: 5 * time.Minute}), rc.RbacEngine.Require("coa_subgroups.read"), rc.COASubGroupHandler.SelectDropdownList)
	rc.App.Get("/api/v1/coa-subgroups/:id", rc.RbacEngine.Require("coa_subgroups.read"), rc.COASubGroupHandler.FindById)
	rc.App.Put("/api/v1/coa-subgroups/:id", rc.RbacEngine.Require("coa_subgroups.manage"), rc.COASubGroupHandler.Update)
	rc.App.Delete("/api/v1/coa-subgroups/:id", rc.RbacEngine.Require("coa_subgroups.manage"), rc.COASubGroupHandler.Delete)

	// COA Management
	rc.App.Get("/api/v1/coa", rc.RbacEngine.Require("coa.read"), rc.COAHandler.FindAll)
	rc.App.Post("/api/v1/coa", rc.RbacEngine.Require("coa.manage"), rc.COAHandler.Create)
	rc.App.Get("/api/v1/coa/:id", rc.RbacEngine.Require("coa.read"), rc.COAHandler.FindById)
	rc.App.Put("/api/v1/coa/:id", rc.RbacEngine.Require("coa.manage"), rc.COAHandler.Update)
	rc.App.Delete("/api/v1/coa/:id", rc.RbacEngine.Require("coa.manage"), rc.COAHandler.Delete)

	// Company Configuration Management
	rc.App.Get("/api/v1/company-configuration", rc.RbacEngine.Require("company_config.read"), rc.CompanyConfigurationHandler.Get)
	rc.App.Put("/api/v1/company-configuration", rc.RbacEngine.Require("company_config.manage"), rc.CompanyConfigurationHandler.Upsert)

	// Fiscal Year Management
	rc.App.Get("/api/v1/fiscal-years", rc.RbacEngine.Require("fiscal_years.read"), rc.FiscalYearHandler.FindAll)
	rc.App.Post("/api/v1/fiscal-years", rc.RbacEngine.Require("fiscal_years.manage"), rc.FiscalYearHandler.Create)
	rc.App.Get("/api/v1/fiscal-years/:id", rc.RbacEngine.Require("fiscal_years.read"), rc.FiscalYearHandler.FindById)
	rc.App.Post("/api/v1/fiscal-years/:id/activate", rc.RbacEngine.Require("fiscal_years.manage"), rc.FiscalYearHandler.Activate)
	rc.App.Post("/api/v1/fiscal-years/periods/:period_id/close", rc.RbacEngine.Require("fiscal_years.manage"), rc.FiscalYearHandler.ClosePeriod)

	// Journal Entries Management
	rc.App.Post("/api/v1/journal_entries", rc.RbacEngine.Require("transactions.manage"), rc.JournalEntryHandler.Create)
	rc.App.Get("/api/v1/journal_entries", rc.RbacEngine.Require("transactions.read"), rc.JournalEntryHandler.FindAll)
	rc.App.Get("/api/v1/journal_entries/:id", rc.RbacEngine.Require("transactions.read"), rc.JournalEntryHandler.FindById)
	rc.App.Put("/api/v1/journal_entries/:id", rc.RbacEngine.Require("transactions.manage"), rc.JournalEntryHandler.Update)
	rc.App.Delete("/api/v1/journal_entries/:id", rc.RbacEngine.Require("transactions.manage"), rc.JournalEntryHandler.Delete)
	rc.App.Put("/api/v1/journal_entries/:id/post", rc.RbacEngine.Require("transactions.manage"), rc.JournalEntryHandler.Post)
	rc.App.Put("/api/v1/journal_entries/:id/void", rc.RbacEngine.Require("transactions.manage"), rc.JournalEntryHandler.Void)
}
