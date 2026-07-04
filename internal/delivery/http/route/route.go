package route

import (
	_ "github.com/fatihrizqon/gofiber-microservice/docs"
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/handler"
	"github.com/fatihrizqon/gofiber-microservice/internal/rbac"
	swagger "github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
)

type RouteConfig struct {
	App            *fiber.App
	AuthMiddleware fiber.Handler
	RbacEngine     *rbac.RBAC
	UserHandler    *handler.UserHandler
	AuthHandler    *handler.AuthHandler
	FileHandler    *handler.FileHandler
	RbacHandler    *handler.RbacHandler
	Production     bool
	CompanyHandler *handler.CompanyHandler
}

func (rc *RouteConfig) Setup() {
	rc.SetupGuestRoute()
	rc.SetupAuthRoute()
}

func (rc *RouteConfig) SetupGuestRoute() {
	rc.App.Post("/api/v1/auth/register", rc.AuthHandler.Register)
	rc.App.Post("/api/v1/auth/login", rc.AuthHandler.Login)
	rc.App.Post("/api/v1/auth/refresh", rc.AuthHandler.Refresh)
	rc.App.Get("/swagger/*", swagger.HandlerDefault)
}

func (rc *RouteConfig) SetupAuthRoute() {
	rc.App.Use(rc.AuthMiddleware)

	rc.App.Post("/api/v1/auth/logout", rc.AuthHandler.Logout)

	// User Profile
	rc.App.Post("/api/v1/users/me/avatar", rc.UserHandler.UploadAvatar)

	rc.App.Get("/api/v1/users", rc.RbacEngine.Require("users.read"), rc.UserHandler.FindAll)
	rc.App.Get("/api/v1/users/:id", rc.RbacEngine.Require("users.read"), rc.UserHandler.FindById)
	rc.App.Put("/api/v1/users/:id", rc.RbacEngine.Require("users.write"), rc.UserHandler.Update)
	rc.App.Delete("/api/v1/users/:id", rc.RbacEngine.Require("users.write"), rc.UserHandler.Delete)
	rc.App.Patch("/api/v1/users/:id/lock", rc.RbacEngine.Require("users.write"), rc.UserHandler.Lock)
	rc.App.Patch("/api/v1/users/:id/unlock", rc.RbacEngine.Require("users.write"), rc.UserHandler.Unlock)

	rc.App.Post("/api/v1/files/upload", rc.FileHandler.Upload)

	// User - Role Management
	rc.App.Post("/api/v1/users/:id/roles", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.AssignRole)
	rc.App.Delete("/api/v1/users/:id/roles/:role_id", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.RevokeRole)

	// Role - Permission Management
	rc.App.Post("/api/v1/roles/:id/permissions", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.AssignPermission)
	rc.App.Delete("/api/v1/roles/:id/permissions/:permission_id", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.RevokePermission)

	// Role Management
	rc.App.Post("/api/v1/roles", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.CreateRole)
	rc.App.Get("/api/v1/roles", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.GetRoles)
	rc.App.Get("/api/v1/roles/:id", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.GetRoleById)
	rc.App.Put("/api/v1/roles/:id", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.UpdateRole)
	rc.App.Delete("/api/v1/roles/:id", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.DeleteRole)

	// Permission Management
	rc.App.Post("/api/v1/permissions", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.CreatePermission)
	rc.App.Get("/api/v1/permissions", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.GetPermissions)
	rc.App.Get("/api/v1/permissions/:id", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.GetPermissionById)
	rc.App.Put("/api/v1/permissions/:id", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.UpdatePermission)
	rc.App.Delete("/api/v1/permissions/:id", rc.RbacEngine.Require("rbac.manage"), rc.RbacHandler.DeletePermission)

	// Company Management
	rc.App.Get("/api/v1/companies", rc.RbacEngine.Require("companies.read"), rc.CompanyHandler.FindAll)
	rc.App.Post("/api/v1/companies", rc.RbacEngine.Require("companies.write"), rc.CompanyHandler.Create)
	rc.App.Get("/api/v1/companies/:id", rc.RbacEngine.Require("companies.read"), rc.CompanyHandler.FindById)
	rc.App.Put("/api/v1/companies/:id", rc.RbacEngine.Require("companies.write"), rc.CompanyHandler.Update)
	rc.App.Delete("/api/v1/companies/:id", rc.RbacEngine.Require("companies.write"), rc.CompanyHandler.Delete)
}
