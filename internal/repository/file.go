package repository

import (
	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IFileRepository interface {
	Create(file entity.File) (entity.File, error)
	FindById(id uuid.UUID) (entity.File, error)
	Delete(id uuid.UUID) error
}

type FileRepository struct {
	Db *gorm.DB
}

func NewFileRepository(db *gorm.DB) IFileRepository {
	return &FileRepository{Db: db}
}

func (r *FileRepository) Create(file entity.File) (entity.File, error) {
	if err := r.Db.Create(&file).Error; err != nil {
		return file, err
	}
	return file, nil
}

func (r *FileRepository) FindById(id uuid.UUID) (entity.File, error) {
	var file entity.File
	if err := r.Db.Where("id = ?", id).First(&file).Error; err != nil {
		return file, err
	}
	return file, nil
}

func (r *FileRepository) Delete(id uuid.UUID) error {
	if err := r.Db.Where("id = ?", id).Delete(&entity.File{}).Error; err != nil {
		return err
	}
	return nil
}
