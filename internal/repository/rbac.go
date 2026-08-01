package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type IRbacRepository interface {
	GetUserPermissions(ctx context.Context, userId string) ([]string, error)
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
	BulkDestroyRoles(ids []string) error
	BulkDestroyPermissions(ids []string) error
}

type RbacRepository struct {
	Db    *gorm.DB
	redis *redis.Client
}

func NewRbacRepository(db *gorm.DB, redis *redis.Client) IRbacRepository {
	return &RbacRepository{Db: db, redis: redis}
}

func (r *RbacRepository) invalidateUserCache(userId string) {
	if r.redis != nil {
		r.redis.Del(context.Background(), fmt.Sprintf("permissions:%s", userId))
	}
}

func (r *RbacRepository) invalidateAllUsersCache() {
	if r.redis != nil {
		ctx := context.Background()
		var cursor uint64
		for {
			var keys []string
			var err error
			keys, cursor, err = r.redis.Scan(ctx, cursor, "permissions:*", 100).Result()
			if err != nil {
				break
			}
			if len(keys) > 0 {
				r.redis.Del(ctx, keys...)
			}
			if cursor == 0 {
				break
			}
		}
	}
}

func (r *RbacRepository) GetUserPermissions(ctx context.Context, userId string) ([]string, error) {
	cacheKey := fmt.Sprintf("permissions:%s", userId)

	if r.redis != nil {
		if cached, err := r.redis.Get(ctx, cacheKey).Result(); err == nil {
			var permissions []string
			if err := json.Unmarshal([]byte(cached), &permissions); err == nil {
				return permissions, nil
			}
		}
	}

	var permissions []string
	err := r.Db.WithContext(ctx).Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userId).
		Pluck("permissions.name", &permissions).Error

	if err != nil {
		return nil, err
	}

	if r.redis != nil {
		if b, err := json.Marshal(permissions); err == nil {
			r.redis.Set(ctx, cacheKey, b, 120*time.Second)
		}
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
	err = r.Db.Model(&user).Association("Roles").Append(&role)
	if err == nil {
		r.invalidateUserCache(userId)
	}
	return err
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
	err = r.Db.Model(&user).Association("Roles").Delete(&role)
	if err == nil {
		r.invalidateUserCache(userId)
	}
	return err
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
	err = r.Db.Model(&role).Association("Permissions").Append(&permission)
	if err == nil {
		r.invalidateAllUsersCache()
	}
	return err
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
	err = r.Db.Model(&role).Association("Permissions").Delete(&permission)
	if err == nil {
		r.invalidateAllUsersCache()
	}
	return err
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

func (r *RbacRepository) CreateRole(entity *entity.Role) error {
	return r.Db.Create(entity).Error
}

func (r *RbacRepository) GetRoles() ([]entity.Role, error) {
	var entities []entity.Role
	err := r.Db.Find(&entities).Error
	return entities, err
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

func (r *RbacRepository) UpdateRole(entity *entity.Role) error {
	err := r.Db.Save(entity).Error
	if err == nil {
		r.invalidateAllUsersCache()
	}
	return err
}

func (r *RbacRepository) DeleteRole(roleId string) error {
	parsedId, err := uuid.Parse(roleId)
	if err != nil {
		return err
	}
	err = r.Db.Delete(&entity.Role{}, parsedId).Error
	if err == nil {
		r.invalidateAllUsersCache()
	}
	return err
}

func (r *RbacRepository) CreatePermission(entity *entity.Permission) error {
	return r.Db.Create(entity).Error
}

func (r *RbacRepository) GetPermissions() ([]entity.Permission, error) {
	var entities []entity.Permission
	err := r.Db.Find(&entities).Error
	return entities, err
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

func (r *RbacRepository) UpdatePermission(entity *entity.Permission) error {
	return r.Db.Save(entity).Error
}

func (r *RbacRepository) DeletePermission(permissionId string) error {
	parsedId, err := uuid.Parse(permissionId)
	if err != nil {
		return err
	}
	return r.Db.Delete(&entity.Permission{}, parsedId).Error
}

func (r *RbacRepository) BulkDestroyRoles(ids []string) error {
	var uuids []uuid.UUID
	for _, id := range ids {
		if parsedId, err := uuid.Parse(id); err == nil {
			uuids = append(uuids, parsedId)
		}
	}
	if len(uuids) == 0 {
		return nil
	}
	err := r.Db.Where("id IN ?", uuids).Delete(&entity.Role{}).Error
	if err == nil {
		r.invalidateAllUsersCache()
	}
	return err
}

func (r *RbacRepository) BulkDestroyPermissions(ids []string) error {
	var uuids []uuid.UUID
	for _, id := range ids {
		if parsedId, err := uuid.Parse(id); err == nil {
			uuids = append(uuids, parsedId)
		}
	}
	if len(uuids) == 0 {
		return nil
	}
	return r.Db.Where("id IN ?", uuids).Delete(&entity.Permission{}).Error
}

