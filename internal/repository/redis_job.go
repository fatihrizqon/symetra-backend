package repository

import (
	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"gorm.io/gorm"
)

type IRedisJobRepository interface {
	Create(job entity.RedisJob) error
	GetPendingJobs(limit int) ([]entity.RedisJob, error)
	UpdateStatus(id string, status string, errStr string) error
}

type RedisJobRepository struct {
	Db *gorm.DB
}

func NewRedisJobRepository(Db *gorm.DB) IRedisJobRepository {
	return &RedisJobRepository{Db: Db}
}

func (r *RedisJobRepository) Create(entity entity.RedisJob) error {
	return r.Db.Create(&entity).Error
}

func (r *RedisJobRepository) GetPendingJobs(limit int) ([]entity.RedisJob, error) {
	var entities []entity.RedisJob
	err := r.Db.Where("status = ?", "PENDING").Order("created_at asc").Limit(limit).Find(&entities).Error
	return entities, err
}

func (r *RedisJobRepository) UpdateStatus(id string, status string, errStr string) error {
	return r.Db.Model(&entity.RedisJob{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status": status,
		"error":  errStr,
	}).Error
}
