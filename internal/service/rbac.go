package service

import (
	"errors"

	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/request"
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/response"
	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type IRbacService interface {
	AssignRole(userId string, req request.AssignRoleRequest) error
	RevokeRole(userId string, roleId string) error
	AssignPermission(roleId string, req request.AssignPermissionRequest) error
	RevokePermission(roleId string, permissionId string) error
	CreateRole(req request.CreateRoleRequest) error
	GetRoles() ([]response.RoleResponse, error)
	GetRoleById(roleId string) (*response.RoleResponse, error)
	UpdateRole(roleId string, req request.UpdateRoleRequest) error
	DeleteRole(roleId string) error
	CreatePermission(req request.CreatePermissionRequest) error
	GetPermissions() ([]response.PermissionResponse, error)
	GetPermissionById(permissionId string) (*response.PermissionResponse, error)
	UpdatePermission(permissionId string, req request.UpdatePermissionRequest) error
	DeletePermission(permissionId string) error
}

type RbacService struct {
	RbacRepo repository.IRbacRepository
	UserRepo repository.IUserRepository
	Validate *validator.Validate
}

func NewRbacService(rbacRepo repository.IRbacRepository, userRepo repository.IUserRepository, validate *validator.Validate) IRbacService {
	return &RbacService{
		RbacRepo: rbacRepo,
		UserRepo: userRepo,
		Validate: validate,
	}
}

func (s *RbacService) AssignRole(userId string, req request.AssignRoleRequest) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		return errors.New("invalid user id")
	}

	_, err = s.UserRepo.FindById(parsedUserId)
	if err != nil {
		return errors.New("user not found")
	}

	exists, err := s.RbacRepo.RoleExists(req.RoleId)
	if err != nil || !exists {
		return errors.New("role not found")
	}

	return s.RbacRepo.AssignRoleToUser(userId, req.RoleId)
}

func (s *RbacService) RevokeRole(userId string, roleId string) error {
	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		return errors.New("invalid user id")
	}

	_, err = s.UserRepo.FindById(parsedUserId)
	if err != nil {
		return errors.New("user not found")
	}

	exists, err := s.RbacRepo.RoleExists(roleId)
	if err != nil || !exists {
		return errors.New("role not found")
	}

	return s.RbacRepo.RevokeRoleFromUser(userId, roleId)
}

func (s *RbacService) AssignPermission(roleId string, req request.AssignPermissionRequest) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	exists, err := s.RbacRepo.RoleExists(roleId)
	if err != nil || !exists {
		return errors.New("role not found")
	}

	existsPerm, err := s.RbacRepo.PermissionExists(req.PermissionId)
	if err != nil || !existsPerm {
		return errors.New("permission not found")
	}

	return s.RbacRepo.AssignPermissionToRole(roleId, req.PermissionId)
}

func (s *RbacService) RevokePermission(roleId string, permissionId string) error {
	exists, err := s.RbacRepo.RoleExists(roleId)
	if err != nil || !exists {
		return errors.New("role not found")
	}

	existsPerm, err := s.RbacRepo.PermissionExists(permissionId)
	if err != nil || !existsPerm {
		return errors.New("permission not found")
	}

	return s.RbacRepo.RevokePermissionFromRole(roleId, permissionId)
}

func (s *RbacService) CreateRole(req request.CreateRoleRequest) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}
	role := &entity.Role{
		Name:        req.Name,
		Description: req.Description,
	}
	return s.RbacRepo.CreateRole(role)
}

func (s *RbacService) GetRoles() ([]response.RoleResponse, error) {
	roles, err := s.RbacRepo.GetRoles()
	if err != nil {
		return nil, err
	}
	var res []response.RoleResponse
	for _, r := range roles {
		res = append(res, response.RoleResponse{
			Id:          r.Id.String(),
			Name:        r.Name,
			Description: r.Description,
		})
	}
	return res, nil
}

func (s *RbacService) GetRoleById(roleId string) (*response.RoleResponse, error) {
	r, err := s.RbacRepo.GetRoleById(roleId)
	if err != nil {
		return nil, err
	}
	return &response.RoleResponse{
		Id:          r.Id.String(),
		Name:        r.Name,
		Description: r.Description,
	}, nil
}

func (s *RbacService) UpdateRole(roleId string, req request.UpdateRoleRequest) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}
	r, err := s.RbacRepo.GetRoleById(roleId)
	if err != nil {
		return errors.New("role not found")
	}
	r.Name = req.Name
	r.Description = req.Description
	return s.RbacRepo.UpdateRole(r)
}

func (s *RbacService) DeleteRole(roleId string) error {
	_, err := s.RbacRepo.GetRoleById(roleId)
	if err != nil {
		return errors.New("role not found")
	}
	return s.RbacRepo.DeleteRole(roleId)
}

func (s *RbacService) CreatePermission(req request.CreatePermissionRequest) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}
	perm := &entity.Permission{
		Name:        req.Name,
		Description: req.Description,
	}
	return s.RbacRepo.CreatePermission(perm)
}

func (s *RbacService) GetPermissions() ([]response.PermissionResponse, error) {
	perms, err := s.RbacRepo.GetPermissions()
	if err != nil {
		return nil, err
	}
	var res []response.PermissionResponse
	for _, p := range perms {
		res = append(res, response.PermissionResponse{
			Id:          p.Id.String(),
			Name:        p.Name,
			Description: p.Description,
		})
	}
	return res, nil
}

func (s *RbacService) GetPermissionById(permissionId string) (*response.PermissionResponse, error) {
	p, err := s.RbacRepo.GetPermissionById(permissionId)
	if err != nil {
		return nil, err
	}
	return &response.PermissionResponse{
		Id:          p.Id.String(),
		Name:        p.Name,
		Description: p.Description,
	}, nil
}

func (s *RbacService) UpdatePermission(permissionId string, req request.UpdatePermissionRequest) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}
	p, err := s.RbacRepo.GetPermissionById(permissionId)
	if err != nil {
		return errors.New("permission not found")
	}
	p.Name = req.Name
	p.Description = req.Description
	return s.RbacRepo.UpdatePermission(p)
}

func (s *RbacService) DeletePermission(permissionId string) error {
	_, err := s.RbacRepo.GetPermissionById(permissionId)
	if err != nil {
		return errors.New("permission not found")
	}
	return s.RbacRepo.DeletePermission(permissionId)
}
