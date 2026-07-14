package repository

import (
	"errors"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"gorm.io/gorm"
)

type IAuthRepository interface {
	Register(ent entity.User) (entity.User, error)
	RegisterWithJobs(ent entity.User, jobs []entity.RedisJob) (entity.User, error)
	Login(email string) (entity.User, error)
}

type AuthRepository struct {
	Db *gorm.DB
}

func NewAuthRepository(Db *gorm.DB) IAuthRepository {
	return &AuthRepository{Db: Db}
}

func (r *AuthRepository) Register(ent entity.User) (entity.User, error) {
	if err := r.Db.Create(&ent).Error; err != nil {
		return ent, err
	}
	return ent, nil
}

func (r *AuthRepository) RegisterWithJobs(ent entity.User, jobs []entity.RedisJob) (entity.User, error) {
	err := r.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&ent).Error; err != nil {
			return err
		}

		for _, job := range jobs {
			if err := tx.Create(&job).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return entity.User{}, err
	}

	return ent, nil
}

func (r *AuthRepository) Login(email string) (entity.User, error) {
	var entity entity.User
	if err := r.Db.Where("email = ?", email).First(&entity).Error; err != nil {
		return entity, errors.New("credentials does not matches our record")
	}
	return entity, nil
}
