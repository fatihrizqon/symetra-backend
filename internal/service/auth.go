package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/request"
	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type AuthResult struct {
	AccessToken  string
	RefreshToken string
	User         entity.User
}

type IAuthService interface {
	Register(req request.RegisterRequest) (entity.User, error)
	Login(req request.LoginRequest) (AuthResult, error)
	RefreshToken(refreshToken string) (AuthResult, error)
	Logout(refreshToken string) error
}

type AuthService struct {
	IAuthRepository  repository.IAuthRepository
	ITokenRepository repository.ITokenRepository
	validate         *validator.Validate
	emailService     IEmailService
	appURL           string
}

func NewAuthService(
	authRepo repository.IAuthRepository,
	tokenRepo repository.ITokenRepository,
	validate *validator.Validate,
	emailService IEmailService,
	appURL string,
) IAuthService {
	return &AuthService{
		IAuthRepository:  authRepo,
		ITokenRepository: tokenRepo,
		validate:         validate,
		emailService:     emailService,
		appURL:           appURL,
	}
}

func (s *AuthService) Register(req request.RegisterRequest) (entity.User, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.User{}, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := entity.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: string(passwordHash),
	}

	verificationLink := fmt.Sprintf("%s/api/v1/auth/verify?email=%s", s.appURL, user.Email)
	job, err := s.emailService.CreateVerificationEmailJob(user.Email, user.Email, verificationLink)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to create email job: %w", err)
	}

	newUser, err := s.IAuthRepository.RegisterWithJobs(user, []entity.RedisJob{*job})
	if err != nil {
		return entity.User{}, err
	}

	return newUser, nil
}

func (s *AuthService) Login(req request.LoginRequest) (AuthResult, error) {
	var res AuthResult

	user, err := s.IAuthRepository.Login(req.Email)
	if err != nil {
		return res, err
	}

	if err := ValidatePassword(req.Password, user.Password); err != nil {
		return res, errors.New("credentials does not matches our record")
	}

	session := entity.Session{
		ID:     util.GenerateUUID(),
		UserID: user.Id,
	}

	session, err = s.ITokenRepository.CreateSession(session)
	if err != nil {
		return res, err
	}

	accessToken, err := util.CreateAccessToken(user, session.ID)
	if err != nil {
		return res, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := util.CreateRefreshToken(user, session.ID)
	if err != nil {
		return res, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	credential := entity.Credential{
		ID:           util.GenerateUUID(),
		SessionID:    session.ID,
		Type:         "REFRESH_TOKEN",
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.ITokenRepository.CreateCredential(credential); err != nil {
		return res, err
	}

	return AuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (AuthResult, error) {
	oldCredential, err := s.ITokenRepository.FindCredentialByToken(refreshToken)
	if err != nil {
		return AuthResult{}, errors.New("invalid or revoked refresh token")
	}

	session, err := s.ITokenRepository.FindSessionByID(oldCredential.SessionID)
	if err != nil {
		return AuthResult{}, errors.New("session revoked or not found")
	}

	accessToken, err := util.CreateAccessToken(session.User, session.ID)
	if err != nil {
		return AuthResult{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := util.CreateRefreshToken(session.User, session.ID)
	if err != nil {
		return AuthResult{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	newCredential := entity.Credential{
		ID:           util.GenerateUUID(),
		SessionID:    session.ID,
		Type:         "REFRESH_TOKEN",
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.ITokenRepository.RevokeCredentialByID(oldCredential.ID); err != nil {
		return AuthResult{}, fmt.Errorf("failed to revoke old refresh token: %w", err)
	}

	if err := s.ITokenRepository.CreateCredential(newCredential); err != nil {
		return AuthResult{}, fmt.Errorf("failed to save new refresh token: %w", err)
	}

	return AuthResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         session.User,
	}, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	credential, err := s.ITokenRepository.FindCredentialByToken(refreshToken)
	if err != nil {
		return err
	}
	if err := s.ITokenRepository.RevokeCredentialByID(credential.ID); err != nil {
		return err
	}
	if err := s.ITokenRepository.RevokeSession(credential.SessionID); err != nil {
		return err
	}
	return nil
}

func ValidatePassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
}
