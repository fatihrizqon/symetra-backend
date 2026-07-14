package repository

import (
	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IFileRepository interface {
	Create(entity entity.File) (entity.File, error)
	FindById(id uuid.UUID) (entity.File, error)
	Delete(id uuid.UUID) error
}

type FileRepository struct {
	Db *gorm.DB
}

func NewFileRepository(db *gorm.DB) IFileRepository {
	return &FileRepository{Db: db}
}

func (r *FileRepository) Create(entity entity.File) (entity.File, error) {
	if err := r.Db.Create(&entity).Error; err != nil {
		return entity, err
	}
	return entity, nil
}

func (r *FileRepository) FindById(id uuid.UUID) (entity.File, error) {
	var entity entity.File
	if err := r.Db.Where("id = ?", id).First(&entity).Error; err != nil {
		return entity, err
	}
	return entity, nil
}

func (r *FileRepository) Delete(id uuid.UUID) error {
	if err := r.Db.Where("id = ?", id).Delete(&entity.File{}).Error; err != nil {
		return err
	}
	return nil
}
