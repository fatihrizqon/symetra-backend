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

type ICOASubGroupService interface {
	Create(companyID uuid.UUID, req request.COASubGroupCreateRequest) (entity.COASubGroup, error)
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.COASubGroup, int, error)
	FindById(companyID uuid.UUID, id uuid.UUID) (entity.COASubGroup, error)
	Update(companyID uuid.UUID, req request.COASubGroupUpdateRequest) (entity.COASubGroup, error)
	Delete(companyID uuid.UUID, id uuid.UUID) error
	SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error)
	Destroy(companyID uuid.UUID, ids []uuid.UUID) error
}

type COASubGroupService struct {
	validate               *validator.Validate
	ICOASubGroupRepository repository.ICOASubGroupRepository
}

func NewCOASubGroupService(validate *validator.Validate, ICOASubGroupRepository repository.ICOASubGroupRepository) ICOASubGroupService {
	return &COASubGroupService{
		validate:               validate,
		ICOASubGroupRepository: ICOASubGroupRepository,
	}
}

func (s *COASubGroupService) Create(companyID uuid.UUID, req request.COASubGroupCreateRequest) (entity.COASubGroup, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.COASubGroup{}, err
	}
	coa_subgroup := entity.COASubGroup{CompanyId: companyID, GroupId: req.GroupId, Code: req.Code, Name: req.Name}
	return s.ICOASubGroupRepository.Create(coa_subgroup)
}

func (s *COASubGroupService) FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.COASubGroup, int, error) {
	coa_subgroups, totalCount, err := s.ICOASubGroupRepository.FindAll(companyID, qp)
	if err != nil {
		return nil, 0, err
	}

	return coa_subgroups, totalCount, nil
}

func (s *COASubGroupService) FindById(companyID uuid.UUID, id uuid.UUID) (entity.COASubGroup, error) {
	coa_subgroup, err := s.ICOASubGroupRepository.FindById(companyID, id)
	if err != nil {
		return entity.COASubGroup{}, err
	}

	return coa_subgroup, nil
}

func (s *COASubGroupService) Update(companyID uuid.UUID, req request.COASubGroupUpdateRequest) (entity.COASubGroup, error) {
	coa_subgroup, err := s.ICOASubGroupRepository.FindById(companyID, req.Id)
	if err != nil {
		return coa_subgroup, err
	}
	coa_subgroup.GroupId = req.GroupId
	coa_subgroup.Code = req.Code
	coa_subgroup.Name = req.Name
	return coa_subgroup, s.ICOASubGroupRepository.Update(coa_subgroup)
}

func (s *COASubGroupService) Delete(companyID uuid.UUID, id uuid.UUID) error {
	return s.ICOASubGroupRepository.Delete(companyID, id)
}

func (s *COASubGroupService) SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error) {
	coa_subgroups, total, err := s.ICOASubGroupRepository.SelectDropdownList(companyID, qp)
	if err != nil {
		return nil, 0, err
	}
	resps := make([]response.SelectDropdownListResponse, 0, len(coa_subgroups))
	for _, v := range coa_subgroups {
		resps = append(resps, response.SelectDropdownListResponse{Value: v.Id, Label: v.Name})
	}
	return resps, total, nil
}

func (s *COASubGroupService) Destroy(companyID uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return fmt.Errorf("no ids provided")
	}
	return s.ICOASubGroupRepository.BulkDestroy(companyID, ids)
}
