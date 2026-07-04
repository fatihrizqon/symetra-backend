package repository

import (
	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IRbacRepository interface {
	GetUserPermissions(userId string) ([]string, error)
	AssignRoleToUser(userId string, roleId string) error
	RevokeRoleFromUser(userId string, roleId string) error
	AssignPermissionToRole(roleId string, permissionId string) error
	RevokePermissionFromRole(roleId string, permissionId string) error
	RoleExists(roleId string) (bool, error)
	PermissionExists(permissionId string) (bool, error)
	CreateRole(role *entity.Role) error
	GetRoles() ([]entity.Role, error)
	GetRoleById(roleId string) (*entity.Role, error)
	UpdateRole(role *entity.Role) error
	DeleteRole(roleId string) error
	CreatePermission(permission *entity.Permission) error
	GetPermissions() ([]entity.Permission, error)
	GetPermissionById(permissionId string) (*entity.Permission, error)
	UpdatePermission(permission *entity.Permission) error
	DeletePermission(permissionId string) error
}

type RbacRepository struct {
	Db *gorm.DB
}

func NewRbacRepository(db *gorm.DB) IRbacRepository {
	return &RbacRepository{Db: db}
}

func (r *RbacRepository) GetUserPermissions(userId string) ([]string, error) {
	var permissions []string

	// Ambil semua name permission berdasarkan relasi many-to-many user -> role -> permission
	err := r.Db.Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userId).
		Pluck("permissions.name", &permissions).Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (r *RbacRepository) AssignRoleToUser(userId string, roleId string) error {
	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		return err
	}
	parsedRoleId, err := uuid.Parse(roleId)
	if err != nil {
		return err
	}
	user := entity.User{Id: parsedUserId}
	role := entity.Role{Id: parsedRoleId}
	return r.Db.Model(&user).Association("Roles").Append(&role)
}

func (r *RbacRepository) RevokeRoleFromUser(userId string, roleId string) error {
	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		return err
	}
	parsedRoleId, err := uuid.Parse(roleId)
	if err != nil {
		return err
	}
	user := entity.User{Id: parsedUserId}
	role := entity.Role{Id: parsedRoleId}
	return r.Db.Model(&user).Association("Roles").Delete(&role)
}

func (r *RbacRepository) AssignPermissionToRole(roleId string, permissionId string) error {
	parsedRoleId, err := uuid.Parse(roleId)
	if err != nil {
		return err
	}
	parsedPermissionId, err := uuid.Parse(permissionId)
	if err != nil {
		return err
	}
	role := entity.Role{Id: parsedRoleId}
	permission := entity.Permission{Id: parsedPermissionId}
	return r.Db.Model(&role).Association("Permissions").Append(&permission)
}

func (r *RbacRepository) RevokePermissionFromRole(roleId string, permissionId string) error {
	parsedRoleId, err := uuid.Parse(roleId)
	if err != nil {
		return err
	}
	parsedPermissionId, err := uuid.Parse(permissionId)
	if err != nil {
		return err
	}
	role := entity.Role{Id: parsedRoleId}
	permission := entity.Permission{Id: parsedPermissionId}
	return r.Db.Model(&role).Association("Permissions").Delete(&permission)
}

func (r *RbacRepository) RoleExists(roleId string) (bool, error) {
	parsedId, err := uuid.Parse(roleId)
	if err != nil {
		return false, err
	}
	var count int64
	err = r.Db.Model(&entity.Role{}).Where("id = ?", parsedId).Count(&count).Error
	return count > 0, err
}

func (r *RbacRepository) PermissionExists(permissionId string) (bool, error) {
	parsedId, err := uuid.Parse(permissionId)
	if err != nil {
		return false, err
	}
	var count int64
	err = r.Db.Model(&entity.Permission{}).Where("id = ?", parsedId).Count(&count).Error
	return count > 0, err
}

func (r *RbacRepository) CreateRole(role *entity.Role) error {
	return r.Db.Create(role).Error
}

func (r *RbacRepository) GetRoles() ([]entity.Role, error) {
	var roles []entity.Role
	err := r.Db.Find(&roles).Error
	return roles, err
}

func (r *RbacRepository) GetRoleById(roleId string) (*entity.Role, error) {
	parsedId, err := uuid.Parse(roleId)
	if err != nil {
		return nil, err
	}
	var role entity.Role
	err = r.Db.Where("id = ?", parsedId).First(&role).Error
	return &role, err
}

func (r *RbacRepository) UpdateRole(role *entity.Role) error {
	return r.Db.Save(role).Error
}

func (r *RbacRepository) DeleteRole(roleId string) error {
	parsedId, err := uuid.Parse(roleId)
	if err != nil {
		return err
	}
	return r.Db.Delete(&entity.Role{}, parsedId).Error
}

func (r *RbacRepository) CreatePermission(permission *entity.Permission) error {
	return r.Db.Create(permission).Error
}

func (r *RbacRepository) GetPermissions() ([]entity.Permission, error) {
	var permissions []entity.Permission
	err := r.Db.Find(&permissions).Error
	return permissions, err
}

func (r *RbacRepository) GetPermissionById(permissionId string) (*entity.Permission, error) {
	parsedId, err := uuid.Parse(permissionId)
	if err != nil {
		return nil, err
	}
	var permission entity.Permission
	err = r.Db.Where("id = ?", parsedId).First(&permission).Error
	return &permission, err
}

func (r *RbacRepository) UpdatePermission(permission *entity.Permission) error {
	return r.Db.Save(permission).Error
}

func (r *RbacRepository) DeletePermission(permissionId string) error {
	parsedId, err := uuid.Parse(permissionId)
	if err != nil {
		return err
	}
	return r.Db.Delete(&entity.Permission{}, parsedId).Error
}
