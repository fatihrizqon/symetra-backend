package rbac

import (
	"github.com/gofiber/fiber/v3"
)

// PermissionStore is an interface that must be implemented by the consumer
// to retrieve user permissions from their database, cache, or other store.
type PermissionStore interface {
	// GetUserPermissions returns a list of permissions (strings) that the user has.
	GetUserPermissions(userID string) ([]string, error)
}

// Config defines the configuration for the RBAC middleware.
type Config struct {
	// Store is the implementation of PermissionStore interface. Required.
	Store PermissionStore

	// UserLookup defines how to retrieve the UserID from the Fiber Context.
	// Default: looks up c.Locals("user_id") and returns it if it's a string.
	UserLookup func(c fiber.Ctx) string

	// UnauthorizedHandler defines the response handler when access is denied.
	// Default: returns a 403 Forbidden status with a JSON message.
	UnauthorizedHandler fiber.Handler
}

// RBAC represents the RBAC middleware instance.
type RBAC struct {
	config Config
}

// New creates a new RBAC middleware instance with the provided config.
func New(cfg Config) *RBAC {
	// Set default UserLookup if not provided
	if cfg.UserLookup == nil {
		cfg.UserLookup = func(c fiber.Ctx) string {
			val := c.Locals("user_id")
			if val == nil {
				return ""
			}
			str, ok := val.(string)
			if !ok {
				return ""
			}
			return str
		}
	}

	// Set default UnauthorizedHandler if not provided
	if cfg.UnauthorizedHandler == nil {
		cfg.UnauthorizedHandler = func(c fiber.Ctx) error {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Forbidden",
			})
		}
	}

	return &RBAC{config: cfg}
}

// Require is a middleware generator that checks if the logged-in user
// has the specified permission.
func (r *RBAC) Require(permission string) fiber.Handler {
	return func(c fiber.Ctx) error {
		// 1. Get userID from Context
		userID := r.config.UserLookup(c)
		if userID == "" {
			return r.config.UnauthorizedHandler(c)
		}

		// 2. Get permissions for the user
		if r.config.Store == nil {
			return fiber.ErrInternalServerError
		}
		permissions, err := r.config.Store.GetUserPermissions(userID)
		if err != nil {
			return err
		}

		// 3. Check if user has the required permission
		hasPermission := false
		for _, p := range permissions {
			if p == permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			return r.config.UnauthorizedHandler(c)
		}

		// 4. Continue to the next handler
		return c.Next()
	}
}
