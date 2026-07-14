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
	FindById(reqId uuid.UUID) (response.UserResponse, error)
	Update(req request.UserUpdateRequest) (response.UserResponse, error)
	Delete(reqId uuid.UUID) (response.UserResponse, error)
	Lock(reqId uuid.UUID) (response.UserResponse, error)
	Unlock(reqId uuid.UUID) (response.UserResponse, error)
	UploadAvatar(ctx context.Context, userId uuid.UUID, file *multipart.FileHeader) (response.UserResponse, error)
}

type UserService struct {
	IUserRepository repository.IUserRepository
	validate        *validator.Validate
	fileService     IFileService
}

func NewUserService(repo repository.IUserRepository, validate *validator.Validate, fileService IFileService) IUserService {
	return &UserService{
		IUserRepository: repo,
		validate:        validate,
		fileService:     fileService,
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

	u := entity.User{
		Username: strings.ToLower(req.Username),
		Name:     req.Name,
		Email:    strings.ToLower(strings.TrimSpace(req.Email)),
		Password: string(hashed),
	}

	err = s.IUserRepository.WithTransaction(func(txRepo repository.IUserRepository) error {
		var txErr error
		u, txErr = txRepo.Create(u)
		return txErr
	})

	return u, err
}

func (s *UserService) FindAll(qp *util.QueryParams) ([]response.UserResponse, int, error) {
	entities, totalCount, err := s.IUserRepository.FindAll(qp)
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

	resps := make([]response.UserResponse, 0, len(entities))
	for _, u := range entities {
		resp := response.UserResponse{
			Id:        u.Id,
			Username:  u.Username,
			Name:      u.Name,
			Email:     u.Email,
			Status:    u.Status,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
			DeletedAt: u.DeletedAt,
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

func (s *UserService) FindById(reqId uuid.UUID) (response.UserResponse, error) {
	u, err := s.IUserRepository.FindById(reqId)
	if err != nil {
		return response.UserResponse{}, err
	}

	resp := response.UserResponse{
		Id:        u.Id,
		Username:  u.Username,
		Name:      u.Name,
		Email:     u.Email,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,
	}

	if u.Avatar != nil {
		resp.AvatarURL = "/uploads/" + u.Avatar.Path
	}

	permissionSet := make(map[string]struct{})

	for _, r := range u.Roles {
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

	var u entity.User
	err := s.IUserRepository.WithTransaction(func(txRepo repository.IUserRepository) error {
		var err error
		u, err = txRepo.FindById(req.Id)
		if err != nil {
			return err
		}

		u.Username = strings.ToLower(req.Username)
		u.Name = req.Name
		u.Email = req.Email

		if hashedPassword != "" {
			u.Password = hashedPassword
		}

		if err := txRepo.Update(u); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
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

func (s *UserService) UploadAvatar(ctx context.Context, userId uuid.UUID, file *multipart.FileHeader) (response.UserResponse, error) {
	fileResp, err := s.fileService.Upload(ctx, file, userId)
	if err != nil {
		return response.UserResponse{}, err
	}

	if err := s.IUserRepository.UpdateAvatar(userId, fileResp.Id); err != nil {
		return response.UserResponse{}, err
	}

	u, err := s.IUserRepository.FindById(userId)
	if err != nil {
		return response.UserResponse{}, err
	}

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
		resp.AvatarURL = "/uploads/" + u.Avatar.Path // Ideally use storage.GetURL but keeping it simple, or inject base url. wait, FileService returns URL in FileResponse!
	}

	return resp, nil
}

func (s *UserService) Delete(reqId uuid.UUID) (response.UserResponse, error) {
	u, err := s.IUserRepository.FindById(reqId)
	if err != nil {
		return response.UserResponse{}, err
	}

	if err := s.IUserRepository.Delete(reqId); err != nil {
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

func (s *UserService) Lock(reqId uuid.UUID) (response.UserResponse, error) {
	u, err := s.IUserRepository.FindById(reqId)
	if err != nil {
		return response.UserResponse{}, err
	}

	if err := s.IUserRepository.UpdateStatus(reqId, 0); err != nil {
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
func (s *UserService) Unlock(reqId uuid.UUID) (response.UserResponse, error) {
	u, err := s.IUserRepository.FindById(reqId)
	if err != nil {
		return response.UserResponse{}, err
	}

	if err := s.IUserRepository.UpdateStatus(reqId, 1); err != nil {
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
