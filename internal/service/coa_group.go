package service

import (
	"fmt"

	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/request"
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/response"
	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ICOAGroupService interface {
	Create(companyID uuid.UUID, req request.COAGroupCreateRequest) (entity.COAGroup, error)
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]response.COAGroupResponse, int, error)
	FindById(companyID, reqId uuid.UUID) (response.COAGroupResponse, error)
	Update(companyID uuid.UUID, req request.COAGroupUpdateRequest) (entity.COAGroup, error)
	Delete(companyID, reqId uuid.UUID) (entity.COAGroup, error)
	SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error)
}

type COAGroupService struct {
	ICOAGroupRepository repository.ICOAGroupRepository
	validate            *validator.Validate
}

func NewCOAGroupService(repo repository.ICOAGroupRepository, validate *validator.Validate) ICOAGroupService {
	return &COAGroupService{ICOAGroupRepository: repo, validate: validate}
}

func (s *COAGroupService) Create(companyID uuid.UUID, req request.COAGroupCreateRequest) (entity.COAGroup, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.COAGroup{}, err
	}
	groupType := entity.COAGroupType(req.Type)
	if !entity.ValidCOAGroupTypes[groupType] {
		return entity.COAGroup{}, fmt.Errorf("invalid type: %s", req.Type)
	}
	g := entity.COAGroup{
		CompanyId:     companyID,
		Code:          req.Code,
		Name:          req.Name,
		Type:          groupType,
		NormalBalance: req.NormalBalance,
		Category:      req.Category,
	}
	return s.ICOAGroupRepository.Create(g)
}

func (s *COAGroupService) FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]response.COAGroupResponse, int, error) {
	entities, totalCount, err := s.ICOAGroupRepository.FindAll(companyID, qp)
	if err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return []response.COAGroupResponse{}, 0, nil
	}
	totalPages := (totalCount + qp.PageSize - 1) / qp.PageSize
	if qp.Page > totalPages {
		return nil, totalCount, nil
	}
	resps := make([]response.COAGroupResponse, 0, len(entities))
	for _, v := range entities {
		resps = append(resps, mapCOAGroup(v))
	}
	return resps, totalCount, nil
}

func (s *COAGroupService) FindById(companyID, reqId uuid.UUID) (response.COAGroupResponse, error) {
	result, err := s.ICOAGroupRepository.FindById(companyID, reqId)
	if err != nil {
		return response.COAGroupResponse{}, err
	}
	return mapCOAGroup(result), nil
}

func (s *COAGroupService) Update(companyID uuid.UUID, req request.COAGroupUpdateRequest) (entity.COAGroup, error) {
	g, err := s.ICOAGroupRepository.FindById(companyID, req.Id)
	if err != nil {
		return g, err
	}
	groupType := entity.COAGroupType(req.Type)
	if !entity.ValidCOAGroupTypes[groupType] {
		return entity.COAGroup{}, fmt.Errorf("invalid type: %s", req.Type)
	}
	g.Code = req.Code
	g.Name = req.Name
	g.Type = groupType
	g.NormalBalance = req.NormalBalance
	g.Category = req.Category
	return g, s.ICOAGroupRepository.Update(g)
}

func (s *COAGroupService) Delete(companyID, reqId uuid.UUID) (entity.COAGroup, error) {
	g, err := s.ICOAGroupRepository.FindById(companyID, reqId)
	if err != nil {
		return g, err
	}
	hasSubgroups, err := s.ICOAGroupRepository.HasSubgroups(companyID, reqId)
	if err != nil {
		return g, err
	}
	if hasSubgroups {
		return g, fmt.Errorf("HAS_CHILD_RECORDS")
	}
	return g, s.ICOAGroupRepository.Delete(companyID, reqId)
}

func (s *COAGroupService) SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error) {
	entities, total, err := s.ICOAGroupRepository.SelectDropdownList(companyID, qp)
	if err != nil {
		return nil, 0, err
	}
	resps := make([]response.SelectDropdownListResponse, 0, len(entities))
	for _, v := range entities {
		resps = append(resps, response.SelectDropdownListResponse{Value: v.Id, Label: v.Name})
	}
	return resps, total, nil
}

func mapCOAGroup(v entity.COAGroup) response.COAGroupResponse {
	return response.COAGroupResponse{
		Id:            v.Id,
		Code:          v.Code,
		Name:          v.Name,
		Type:          string(v.Type),
		NormalBalance: v.NormalBalance,
		Category:      v.Category,
		Status:        v.Status,
		CreatedAt:     v.CreatedAt,
		UpdatedAt:     v.UpdatedAt,
	}
}
