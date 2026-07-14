package repository

import (
	"context"
	"time"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// userSortColumns is the whitelist of sortable columns for the users table.
// Key = public-facing ?sort= value, Value = safe SQL column expression.
var userSortColumns = map[string]string{
	"name":       "users.name",
	"username":   "users.username",
	"email":      "users.email",
	"status":     "users.status",
	"created_at": "users.created_at",
	"updated_at": "users.updated_at",
}

type IUserRepository interface {
	WithTransaction(fn func(txRepo IUserRepository) error) error
	Create(entity.User) (entity.User, error)
	FindAll(qp *util.QueryParams) ([]entity.User, int, error)
	FindById(id uuid.UUID) (entity.User, error)
	Update(entity.User) error
	Delete(id uuid.UUID) error
	UpdateAvatar(userId uuid.UUID, fileId uuid.UUID) error
	UpdateStatus(id uuid.UUID, status int) error
	CountAll(ctx context.Context) (int64, error)
	CountBetween(ctx context.Context, from time.Time, to time.Time) (int64, error)
	CountVerified(ctx context.Context) (int64, error)
}

type UserRepository struct {
	Db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{Db: db}
}

func (r *UserRepository) WithTransaction(fn func(txRepo IUserRepository) error) error {
	return r.Db.Transaction(func(tx *gorm.DB) error {
		txRepo := &UserRepository{Db: tx}
		return fn(txRepo)
	})
}

func (r *UserRepository) Create(ent entity.User) (entity.User, error) {
	if err := r.Db.Create(&ent).Error; err != nil {
		return ent, err
	}
	return ent, nil
}

func (r *UserRepository) FindAll(qp *util.QueryParams) ([]entity.User, int, error) {
	var entities []entity.User
	var totalCount int64

	query := r.Db.Model(&entity.User{}).Preload("Avatar").Preload("Roles").Where("deleted_at IS NULL")
	query = util.ApplySearch(query, qp)
	query = entity.User{}.ApplyFilters(query, qp.Filters)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return entities, 0, nil
	}

	query = util.ApplySort(query, qp, userSortColumns, "users.created_at")
	query = util.ApplyPagination(query, qp)

	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, int(totalCount), nil
}

func (r *UserRepository) FindById(id uuid.UUID) (entity.User, error) {
	var ent entity.User
	if err := r.Db.Preload("Avatar").Preload("Roles").Preload("Roles.Permissions").Where("id = ?", id).First(&ent).Error; err != nil {
		return ent, err
	}
	return ent, nil
}

func (r *UserRepository) Update(ent entity.User) error {
	if err := r.Db.Model(&ent).Updates(ent).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	if err := r.Db.Model(&entity.User{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UpdateAvatar(userId uuid.UUID, fileId uuid.UUID) error {
	if err := r.Db.Model(&entity.User{}).Where("id = ?", userId).Update("avatar_file_id", fileId).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UpdateStatus(id uuid.UUID, status int) error {
	if err := r.Db.Model(&entity.User{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

// CountAll implements IUserRepository.
func (r *UserRepository) CountAll(ctx context.Context) (int64, error) {
	var count int64
	if err := r.Db.WithContext(ctx).Model(&entity.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountBetween(ctx, from time.Time, to time.Time) implements IUserRepository.
func (r *UserRepository) CountBetween(ctx context.Context, from time.Time, to time.Time) (int64, error) {
	var count int64
	if err := r.Db.WithContext(ctx).Model(&entity.User{}).Where("created_at BETWEEN ? AND ?", from, to).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountVerified implements IUserRepository.
func (r *UserRepository) CountVerified(ctx context.Context) (int64, error) {
	var count int64
	if err := r.Db.WithContext(ctx).Model(&entity.User{}).Where("email_verified_at IS NOT NULL").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
