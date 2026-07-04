package handler

import (
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/request"
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/response"
	"github.com/fatihrizqon/gofiber-microservice/internal/service"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/gofiber/fiber/v3"
)

type RbacHandler struct {
	IRbacService service.IRbacService
}

func NewRbacHandler(serv service.IRbacService) *RbacHandler {
	return &RbacHandler{IRbacService: serv}
}

func (h *RbacHandler) AssignRole(ctx fiber.Ctx) error {
	userId := ctx.Params("id")
	req := request.AssignRoleRequest{}
	if err := ctx.Bind().Body(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	if err := h.IRbacService.AssignRole(userId, req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Role successfully assigned to user.",
	})
}

func (h *RbacHandler) RevokeRole(ctx fiber.Ctx) error {
	userId := ctx.Params("id")
	roleId := ctx.Params("role_id")

	if err := h.IRbacService.RevokeRole(userId, roleId); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Role successfully revoked from user.",
	})
}

func (h *RbacHandler) AssignPermission(ctx fiber.Ctx) error {
	roleId := ctx.Params("id")
	req := request.AssignPermissionRequest{}
	if err := ctx.Bind().Body(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}

	if err := h.IRbacService.AssignPermission(roleId, req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Permission successfully assigned to role.",
	})
}

func (h *RbacHandler) RevokePermission(ctx fiber.Ctx) error {
	roleId := ctx.Params("id")
	permissionId := ctx.Params("permission_id")

	if err := h.IRbacService.RevokePermission(roleId, permissionId); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Permission successfully revoked from role.",
	})
}

func (h *RbacHandler) CreateRole(ctx fiber.Ctx) error {
	req := request.CreateRoleRequest{}
	if err := ctx.Bind().Body(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	if err := h.IRbacService.CreateRole(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(response.JSON{
		Status:  fiber.StatusCreated,
		Message: "Role successfully created.",
	})
}

func (h *RbacHandler) GetRoles(ctx fiber.Ctx) error {
	roles, err := h.IRbacService.GetRoles()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status: fiber.StatusOK,
		Data:   roles,
	})
}

func (h *RbacHandler) GetRoleById(ctx fiber.Ctx) error {
	roleId := ctx.Params("id")
	role, err := h.IRbacService.GetRoleById(roleId)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status:  fiber.StatusNotFound,
			Message: err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status: fiber.StatusOK,
		Data:   role,
	})
}

func (h *RbacHandler) UpdateRole(ctx fiber.Ctx) error {
	roleId := ctx.Params("id")
	req := request.UpdateRoleRequest{}
	if err := ctx.Bind().Body(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	if err := h.IRbacService.UpdateRole(roleId, req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Role successfully updated.",
	})
}

func (h *RbacHandler) DeleteRole(ctx fiber.Ctx) error {
	roleId := ctx.Params("id")
	if err := h.IRbacService.DeleteRole(roleId); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Role successfully deleted.",
	})
}

func (h *RbacHandler) CreatePermission(ctx fiber.Ctx) error {
	req := request.CreatePermissionRequest{}
	if err := ctx.Bind().Body(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	if err := h.IRbacService.CreatePermission(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(response.JSON{
		Status:  fiber.StatusCreated,
		Message: "Permission successfully created.",
	})
}

func (h *RbacHandler) GetPermissions(ctx fiber.Ctx) error {
	perms, err := h.IRbacService.GetPermissions()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.JSON{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status: fiber.StatusOK,
		Data:   perms,
	})
}

func (h *RbacHandler) GetPermissionById(ctx fiber.Ctx) error {
	permId := ctx.Params("id")
	perm, err := h.IRbacService.GetPermissionById(permId)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(response.JSON{
			Status:  fiber.StatusNotFound,
			Message: err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status: fiber.StatusOK,
		Data:   perm,
	})
}

func (h *RbacHandler) UpdatePermission(ctx fiber.Ctx) error {
	permId := ctx.Params("id")
	req := request.UpdatePermissionRequest{}
	if err := ctx.Bind().Body(&req); err != nil {
		util.HandleError(ctx, fiber.StatusBadRequest, err)
		return nil
	}
	if err := h.IRbacService.UpdatePermission(permId, req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Permission successfully updated.",
	})
}

func (h *RbacHandler) DeletePermission(ctx fiber.Ctx) error {
	permId := ctx.Params("id")
	if err := h.IRbacService.DeletePermission(permId); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(response.JSON{
		Status:  fiber.StatusOK,
		Message: "Permission successfully deleted.",
	})
}
