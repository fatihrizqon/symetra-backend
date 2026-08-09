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
		companyID, ok := ctx.Locals("company_id").(uuid.UUID)
		if !ok || companyID == uuid.Nil {
			util.HandleError(ctx, fiber.StatusBadRequest, fmt.Errorf(
				"no active company selected. Call POST /api/v1/companies/{id}/select first",
			))
			return nil
		}

		userID, err := util.GetAuthorID(ctx)
		if err != nil {
			util.HandleError(ctx, fiber.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return nil
		}

		member, err := companyRepo.FindMember(ctx.Context(), companyID, userID)
		if err != nil {
			util.HandleError(ctx, fiber.StatusForbidden, fmt.Errorf("you are not a member of this company"))
			return nil
		}

		ctx.Locals("role_id", member.Role)

		return ctx.Next()
	}
}
