package middleware

import (
	"fmt"

	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func NewCompany(companyRepo repository.ICompanyRepository) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		// 1. Read X-Company-ID header
		companyHeader := ctx.Get("X-Company-ID", "")
		if companyHeader == "" {
			util.HandleError(ctx, fiber.StatusBadRequest, fmt.Errorf("missing X-Company-ID header"))
			return nil
		}

		// 2. Parse the company ID
		companyID, err := uuid.Parse(companyHeader)
		if err != nil {
			util.HandleError(ctx, fiber.StatusBadRequest, fmt.Errorf("invalid X-Company-ID format"))
			return nil
		}

		// 3. Get the authenticated user from c.Locals (set by AuthMiddleware)
		userID, err := util.GetAuthor(ctx)
		if err != nil {
			util.HandleError(ctx, fiber.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return nil
		}

		// 4. Validate user is a member of the company and retrieve the role
		member, err := companyRepo.FindMember(companyID, userID)
		if err != nil {
			util.HandleError(ctx, fiber.StatusForbidden, fmt.Errorf("you are not a member of this company"))
			return nil
		}

		// 5. Set company_id and role_id in Locals for downstream handlers
		ctx.Locals("company_id", companyID)
		ctx.Locals("role_id", member.Role)

		return ctx.Next()
	}
}
