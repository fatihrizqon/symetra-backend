package service

import (
	"context"
	"errors"
	"mime/multipart"
	"strings"

	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/request"
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/response"
	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	Create(req request.UserCreateRequest) (entity.User, error)
	FindAll(qp *util.QueryParams) ([]response.UserResponse, int, error)
	FindById(id uuid.UUID) (response.UserResponse, error)
	Update(req request.UserUpdateRequest) (response.UserResponse, error)
	Delete(id uuid.UUID) error
	Lock(id uuid.UUID) (response.UserResponse, error)
	Unlock(id uuid.UUID) (response.UserResponse, error)
	UploadAvatar(ctx context.Context, userId uuid.UUID, file *multipart.FileHeader) (response.UserResponse, error)
	Destroy(ids []uuid.UUID) error
}

type UserService struct {
	validate        *validator.Validate
	IFileService    IFileService
	IUserRepository repository.IUserRepository
}

func NewUserService(validate *validator.Validate, IFileService IFileService, IUserRepository repository.IUserRepository) IUserService {
	return &UserService{
		validate:        validate,
		IFileService:    IFileService,
		IUserRepository: IUserRepository,
	}
}

func (s *UserService) Create(req request.UserCreateRequest) (entity.User, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.User{}, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), util.PasswordHashCost)
	if err != nil {
		return entity.User{}, errors.New("failed to hash password")
	}

	user := entity.User{
		Username: strings.ToLower(strings.TrimSpace(req.Username)),
		Name:     strings.TrimSpace(req.Name),
		Email:    strings.ToLower(strings.TrimSpace(req.Email)),
		Password: string(hashed),
	}

	err = s.IUserRepository.WithTransaction(func(txRepo repository.IUserRepository) error {
		var txErr error
		user, txErr = txRepo.Create(user)
		return txErr
	})

	return user, err
}

func (s *UserService) FindAll(qp *util.QueryParams) ([]response.UserResponse, int, error) {
	users, totalCount, err := s.IUserRepository.FindAll(qp)
	if err != nil {
		return nil, 0, err
	}

	if totalCount == 0 {
		return []response.UserResponse{}, 0, nil
	}

	totalPages := (totalCount + qp.PageSize - 1) / qp.PageSize
	if qp.Page > totalPages {
		return nil, totalCount, nil
	}

	resps := make([]response.UserResponse, 0, len(users))
	for _, u := range users {
		resp := response.UserResponse{
			Id:        u.Id,
			Username:  u.Username,
			Name:      u.Name,
			Email:     u.Email,
			Status:    u.Status,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		}
		if u.Avatar != nil {
			resp.AvatarURL = "/uploads/" + u.Avatar.Path
		}
		for _, r := range u.Roles {
			resp.Roles = append(resp.Roles, r.Name)
		}
		resps = append(resps, resp)
	}

	return resps, totalCount, nil
}

func (s *UserService) FindById(id uuid.UUID) (response.UserResponse, error) {
	user, err := s.IUserRepository.FindById(id)
	if err != nil {
		return response.UserResponse{}, err
	}

	resp := response.UserResponse{
		Id:        user.Id,
		Username:  user.Username,
		Name:      user.Name,
		Email:     user.Email,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.Avatar != nil {
		resp.AvatarURL = "/uploads/" + user.Avatar.Path
	}

	permissionSet := make(map[string]struct{})

	for _, r := range user.Roles {
		resp.Roles = append(resp.Roles, r.Name)

		for _, p := range r.Permissions {
			if _, exists := permissionSet[p.Name]; !exists {
				permissionSet[p.Name] = struct{}{}
				resp.Permissions = append(resp.Permissions, p.Name)
			}
		}
	}

	return resp, nil
}

func (s *UserService) Update(req request.UserUpdateRequest) (response.UserResponse, error) {
	var hashedPassword string
	if req.Password != "" {
		hashed, errHash := bcrypt.GenerateFromPassword([]byte(req.Password), util.PasswordHashCost)
		if errHash != nil {
			return response.UserResponse{}, errors.New("failed to generate password")
		}
		hashedPassword = string(hashed)
	}

	var user entity.User
	err := s.IUserRepository.WithTransaction(func(txRepo repository.IUserRepository) error {
		var err error
		user, err = txRepo.FindById(req.Id)
		if err != nil {
			return err
		}

		user.Username = strings.ToLower(req.Username)
		user.Name = req.Name
		user.Email = req.Email

		if hashedPassword != "" {
			user.Password = hashedPassword
		}

		if err := txRepo.Update(user); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return response.UserResponse{}, err
	}

	return response.UserResponse{
		Id:        user.Id,
		Username:  user.Username,
		Name:      user.Name,
		Email:     user.Email,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *UserService) UploadAvatar(ctx context.Context, userId uuid.UUID, file *multipart.FileHeader) (response.UserResponse, error) {
	fileResp, err := s.IFileService.Upload(ctx, file, userId)
	if err != nil {
		return response.UserResponse{}, err
	}

	if err := s.IUserRepository.UpdateAvatar(userId, fileResp.Id); err != nil {
		return response.UserResponse{}, err
	}

	user, err := s.IUserRepository.FindById(userId)
	if err != nil {
		return response.UserResponse{}, err
	}

	resp := response.UserResponse{
		Id:        user.Id,
		Username:  user.Username,
		Name:      user.Name,
		Email:     user.Email,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.Avatar != nil {
		resp.AvatarURL = "/uploads/" + user.Avatar.Path // Ideally use storage.GetURL but keeping it simple, or inject base url. wait, FileService returns URL in FileResponse!
	}

	return resp, nil
}

func (s *UserService) Delete(id uuid.UUID) error {
	return s.IUserRepository.Delete(id)
}

func (s *UserService) Lock(id uuid.UUID) (response.UserResponse, error) {
	u, err := s.IUserRepository.FindById(id)
	if err != nil {
		return response.UserResponse{}, err
	}

	if err := s.IUserRepository.UpdateStatus(id, 0); err != nil {
		return response.UserResponse{}, err
	}

	return response.UserResponse{
		Id:        u.Id,
		Username:  u.Username,
		Name:      u.Name,
		Email:     u.Email,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}
func (s *UserService) Unlock(id uuid.UUID) (response.UserResponse, error) {
	u, err := s.IUserRepository.FindById(id)
	if err != nil {
		return response.UserResponse{}, err
	}

	if err := s.IUserRepository.UpdateStatus(id, 1); err != nil {
		return response.UserResponse{}, err
	}

	return response.UserResponse{
		Id:        u.Id,
		Username:  u.Username,
		Name:      u.Name,
		Email:     u.Email,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

func (s *UserService) Destroy(ids []uuid.UUID) error {
	if len(ids) == 0 {
		return errors.New("no ids provided")
	}
	return s.IUserRepository.BulkDestroy(ids)
}
