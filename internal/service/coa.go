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

type ICOAService interface {
	Create(companyID uuid.UUID, req request.COACreateRequest) (entity.COA, error)
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.COA, int, error)
	FindById(companyID, id uuid.UUID) (entity.COA, error)
	Update(companyID uuid.UUID, req request.COAUpdateRequest) (entity.COA, error)
	Delete(companyID uuid.UUID, id uuid.UUID) error
	SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error)
	Destroy(companyID uuid.UUID, ids []uuid.UUID) error
}

type COAService struct {
	validate       *validator.Validate
	ICOARepository repository.ICOARepository
}

func NewCOAService(validate *validator.Validate, ICOARepository repository.ICOARepository) ICOAService {
	return &COAService{
		validate:       validate,
		ICOARepository: ICOARepository,
	}
}

func (s *COAService) Create(companyID uuid.UUID, req request.COACreateRequest) (entity.COA, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.COA{}, err
	}

	coa := entity.COA{
		CompanyId:   companyID,
		SubgroupId:  req.SubgroupId,
		Code:        req.Code,
		Name:        req.Name,
		IsContra:    req.IsContra,
		ControlType: req.ControlType,
		Active:      true,
	}
	return s.ICOARepository.Create(coa)
}

func (s *COAService) FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.COA, int, error) {
	coas, totalCount, err := s.ICOARepository.FindAll(companyID, qp)
	if err != nil {
		return nil, 0, err
	}
	return coas, totalCount, nil
}

func (s *COAService) FindById(companyID uuid.UUID, id uuid.UUID) (entity.COA, error) {
	coa, err := s.ICOARepository.FindById(companyID, id)
	if err != nil {
		return entity.COA{}, err
	}

	return coa, nil
}

func (s *COAService) Update(companyID uuid.UUID, req request.COAUpdateRequest) (entity.COA, error) {
	coa, err := s.ICOARepository.FindById(companyID, req.Id)
	if err != nil {
		return coa, err
	}
	coa.SubgroupId = req.SubgroupId
	coa.Code = req.Code
	coa.Name = req.Name
	coa.IsContra = req.IsContra
	coa.ControlType = req.ControlType
	return coa, s.ICOARepository.Update(coa)
}

func (s *COAService) Delete(companyID uuid.UUID, id uuid.UUID) error {
	return s.ICOARepository.Delete(companyID, id)
}

func (s *COAService) SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error) {
	coas, total, err := s.ICOARepository.SelectDropdownList(companyID, qp)
	if err != nil {
		return nil, 0, err
	}
	resps := make([]response.SelectDropdownListResponse, 0, len(coas))
	for _, v := range coas {
		resps = append(resps, response.SelectDropdownListResponse{Value: v.Id, Label: fmt.Sprintf("%s - %s", v.Code, v.Name)})
	}
	return resps, total, nil
}

func (s *COAService) Destroy(companyID uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return fmt.Errorf("no ids provided")
	}
	return s.ICOARepository.BulkDestroy(companyID, ids)
}
