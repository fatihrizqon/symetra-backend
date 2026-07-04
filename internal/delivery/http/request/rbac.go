package request

type AssignRoleRequest struct {
	RoleId string `validate:"required,uuid" json:"role_id"`
}

type AssignPermissionRequest struct {
	PermissionId string `validate:"required,uuid" json:"permission_id"`
}

type CreateRoleRequest struct {
	Name        string `validate:"required,min=3,max=100" json:"name"`
	Description string `validate:"max=255" json:"description"`
}

type UpdateRoleRequest struct {
	Name        string `validate:"required,min=3,max=100" json:"name"`
	Description string `validate:"max=255" json:"description"`
}

type CreatePermissionRequest struct {
	Name        string `validate:"required,min=3,max=100" json:"name"`
	Description string `validate:"max=255" json:"description"`
}

type UpdatePermissionRequest struct {
	Name        string `validate:"required,min=3,max=100" json:"name"`
	Description string `validate:"max=255" json:"description"`
}
