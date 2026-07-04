package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITokenRepository interface {
	CreateSession(session entity.Session) (entity.Session, error)
	FindSessionByID(sessionID uuid.UUID) (entity.Session, error)
	RevokeSession(sessionID uuid.UUID) error
	CreateCredential(credential entity.Credential) error
	FindCredentialByToken(token string) (entity.Credential, error)
	RevokeCredentialByID(id uuid.UUID) error
	RevokeCredentialBySession(sessionID uuid.UUID) error
}

type TokenRepository struct {
	Db *gorm.DB
}

func NewTokenRepository(Db *gorm.DB) ITokenRepository {
	return &TokenRepository{Db: Db}
}

func (r *TokenRepository) CreateSession(session entity.Session) (entity.Session, error) {
	if err := r.Db.Create(&session).Error; err != nil {
		return session, errors.New("failed to create session")
	}
	return session, nil
}

func (r *TokenRepository) FindSessionByID(sessionID uuid.UUID) (entity.Session, error) {
	var session entity.Session

	if err := r.Db.
		Preload("User").
		Where("id = ? AND revoked_at IS NULL", sessionID).
		First(&session).Error; err != nil {

		return session, errors.New("session not found or revoked")
	}

	return session, nil
}

func (r *TokenRepository) RevokeSession(sessionID uuid.UUID) error {
	now := time.Now()

	if err := r.Db.Model(&entity.Session{}).
		Where("id = ?", sessionID).
		Update("revoked_at", now).Error; err != nil {

		return errors.New("failed to revoke session")
	}

	return nil
}

func (r *TokenRepository) CreateCredential(credential entity.Credential) error {
	if err := r.Db.Create(&credential).Error; err != nil {
		return errors.New("failed to create credential : " + err.Error())
	}
	return nil
}

func (r *TokenRepository) FindCredentialByToken(refreshToken string) (entity.Credential, error) {
	var credential entity.Credential

	err := r.Db.
		Where("refresh_token = ? AND revoked_at IS NULL AND expires_at > ?", refreshToken, time.Now()).
		First(&credential).Error

	if err != nil {
		return credential, fmt.Errorf("invalid, expired, or revoked refresh token")
	}

	return credential, nil
}

func (r *TokenRepository) RevokeCredentialByID(id uuid.UUID) error {
	now := time.Now()

	if err := r.Db.Model(&entity.Credential{}).
		Where("id = ?", id).
		Update("revoked_at", now).Error; err != nil {

		return errors.New("failed to revoke credential")
	}

	return nil
}

func (r *TokenRepository) RevokeCredentialBySession(sessionID uuid.UUID) error {
	now := time.Now()

	if err := r.Db.Model(&entity.Credential{}).
		Where("session_id = ?", sessionID).
		Update("revoked_at", now).Error; err != nil {

		return errors.New("failed to revoke credentials")
	}

	return nil
}
