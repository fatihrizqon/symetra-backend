package repository

import (
	"errors"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"gorm.io/gorm"
)

type IAuthRepository interface {
	Register(user entity.User) (entity.User, error)
	RegisterWithJobs(user entity.User, jobs []entity.RedisJob) (entity.User, error)
	Login(username string) (entity.User, error)
}

type AuthRepository struct {
	Db *gorm.DB
}

func NewAuthRepository(Db *gorm.DB) IAuthRepository {
	return &AuthRepository{Db: Db}
}

func (r *AuthRepository) Register(user entity.User) (entity.User, error) {
	if err := r.Db.Create(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *AuthRepository) RegisterWithJobs(user entity.User, jobs []entity.RedisJob) (entity.User, error) {
	err := r.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
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
	
	return user, nil
}

func (r *AuthRepository) Login(email string) (entity.User, error) {
	var entity entity.User
	if err := r.Db.Where("email = ?", email).First(&entity).Error; err != nil {
		return entity, errors.New("credentials does not matches our record")
	}
	return entity, nil
}
