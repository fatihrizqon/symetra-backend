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
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.COAGroup, int, error)
	FindById(companyID uuid.UUID, reqId uuid.UUID) (entity.COAGroup, error)
	Update(companyID uuid.UUID, req request.COAGroupUpdateRequest) (entity.COAGroup, error)
	Delete(companyID uuid.UUID, reqId uuid.UUID) error
	SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error)
	Destroy(companyID uuid.UUID, ids []uuid.UUID) error
}

type COAGroupService struct {
	validate            *validator.Validate
	ICOAGroupRepository repository.ICOAGroupRepository
}

func NewCOAGroupService(validate *validator.Validate, repo repository.ICOAGroupRepository) ICOAGroupService {
	return &COAGroupService{
		validate:            validate,
		ICOAGroupRepository: repo,
	}
}

func (s *COAGroupService) Create(companyID uuid.UUID, req request.COAGroupCreateRequest) (entity.COAGroup, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.COAGroup{}, err
	}

	groupType := entity.COAGroupType(req.Type)
	if !entity.ValidCOAGroupTypes[groupType] {
		return entity.COAGroup{}, fmt.Errorf("invalid type: %s", req.Type)
	}

	coa_group := entity.COAGroup{
		CompanyId: companyID,
		Code:      req.Code,
		Name:      req.Name,
		Type:      groupType,
		Category:  req.Category,
	}

	return s.ICOAGroupRepository.Create(coa_group)
}

func (s *COAGroupService) FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.COAGroup, int, error) {
	coa_groups, totalCount, err := s.ICOAGroupRepository.FindAll(companyID, qp)
	if err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return []entity.COAGroup{}, 0, nil
	}
	totalPages := (totalCount + qp.PageSize - 1) / qp.PageSize
	if qp.Page > totalPages {
		return nil, totalCount, nil
	}
	return coa_groups, totalCount, nil
}

func (s *COAGroupService) FindById(companyID uuid.UUID, reqId uuid.UUID) (entity.COAGroup, error) {
	coa_group, err := s.ICOAGroupRepository.FindById(companyID, reqId)
	if err != nil {
		return entity.COAGroup{}, err
	}
	return coa_group, nil
}

func (s *COAGroupService) Update(companyID uuid.UUID, req request.COAGroupUpdateRequest) (entity.COAGroup, error) {
	coa_group, err := s.ICOAGroupRepository.FindById(companyID, req.Id)
	if err != nil {
		return entity.COAGroup{}, err
	}
	groupType := entity.COAGroupType(req.Type)
	if !entity.ValidCOAGroupTypes[groupType] {
		return entity.COAGroup{}, fmt.Errorf("invalid type: %s", req.Type)
	}

	coa_group.Code = req.Code
	coa_group.Name = req.Name
	coa_group.Type = groupType
	coa_group.Category = req.Category

	return coa_group, s.ICOAGroupRepository.Update(coa_group)
}

func (s *COAGroupService) Delete(companyID uuid.UUID, reqId uuid.UUID) error {
	return s.ICOAGroupRepository.Delete(companyID, reqId)
}

func (s *COAGroupService) SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error) {
	coa_groups, total, err := s.ICOAGroupRepository.SelectDropdownList(companyID, qp)
	if err != nil {
		return nil, 0, err
	}
	results := make([]response.SelectDropdownListResponse, 0, len(coa_groups))
	for _, v := range coa_groups {
		results = append(results, response.SelectDropdownListResponse{Value: v.Id, Label: v.Name})
	}
	return results, total, nil
}

func (s *COAGroupService) Destroy(companyID uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return fmt.Errorf("no ids provided")
	}
	return s.ICOAGroupRepository.BulkDestroy(companyID, ids)
}
