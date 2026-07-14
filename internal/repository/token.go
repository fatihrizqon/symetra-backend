package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ITokenRepository interface {
	CreateSession(session entity.Session) (entity.Session, error)
	FindSessionByID(ctx context.Context, sessionID uuid.UUID) (entity.Session, error)
	RevokeSession(sessionID uuid.UUID) error
	CreateCredential(credential entity.Credential) error
	FindCredentialByToken(token string) (entity.Credential, error)
	RevokeCredentialByID(id uuid.UUID) error
	RevokeCredentialBySession(sessionID uuid.UUID) error
	SetActiveCompany(sessionID uuid.UUID, companyID uuid.UUID) error
	GetActiveCompany(sessionID uuid.UUID) (*uuid.UUID, error)
}

type TokenRepository struct {
	Db    *gorm.DB
	redis *redis.Client
}

func NewTokenRepository(Db *gorm.DB, redis *redis.Client) ITokenRepository {
	return &TokenRepository{Db: Db, redis: redis}
}

func (r *TokenRepository) CreateSession(session entity.Session) (entity.Session, error) {
	if err := r.Db.Create(&session).Error; err != nil {
		return session, errors.New("failed to create session")
	}
	return session, nil
}

func (r *TokenRepository) FindSessionByID(ctx context.Context, sessionID uuid.UUID) (entity.Session, error) {
	var session entity.Session
	cacheKey := fmt.Sprintf("session:%s", sessionID.String())

	// Try Redis first
	if cached, err := r.redis.Get(ctx, cacheKey).Result(); err == nil {
		if err := json.Unmarshal([]byte(cached), &session); err == nil {
			return session, nil
		}
	}

	// Fallback to DB
	if err := r.Db.WithContext(ctx).
		Preload("User").
		Where("id = ? AND revoked_at IS NULL", sessionID).
		First(&session).Error; err != nil {

		return session, errors.New("session not found or revoked")
	}

	// Set cache
	if b, err := json.Marshal(session); err == nil {
		r.redis.Set(ctx, cacheKey, b, 60*time.Second)
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

	ctx := context.Background()
	cacheKey := fmt.Sprintf("session:%s", sessionID.String())
	r.redis.Del(ctx, cacheKey)

	return nil
}

func (r *TokenRepository) CreateCredential(entity entity.Credential) error {
	if err := r.Db.Create(&entity).Error; err != nil {
		return errors.New("failed to create credential : " + err.Error())
	}
	return nil
}

func (r *TokenRepository) FindCredentialByToken(refreshToken string) (entity.Credential, error) {
	var entity entity.Credential

	err := r.Db.
		Where("refresh_token = ? AND revoked_at IS NULL AND expires_at > ?", refreshToken, time.Now()).
		First(&entity).Error

	if err != nil {
		return entity, fmt.Errorf("invalid, expired, or revoked refresh token")
	}

	return entity, nil
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

func (r *TokenRepository) SetActiveCompany(sessionID uuid.UUID, companyID uuid.UUID) error {
	err := r.Db.Model(&entity.Session{}).
		Where("id = ? AND revoked_at IS NULL", sessionID).
		Update("active_company_id", companyID).Error

	if err == nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("session:%s", sessionID.String())
		r.redis.Del(ctx, cacheKey)
	}
	return err
}

func (r *TokenRepository) GetActiveCompany(sessionID uuid.UUID) (*uuid.UUID, error) {
	var session entity.Session
	err := r.Db.
		Select("active_company_id").
		Where("id = ? AND revoked_at IS NULL", sessionID).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return session.ActiveCompanyID, nil
}
