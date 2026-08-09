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
	FindById(companyID, reqId uuid.UUID) (entity.COA, error)
	Update(companyID uuid.UUID, req request.COAUpdateRequest) (entity.COA, error)
	Delete(companyID uuid.UUID, reqId uuid.UUID) error
	SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error)
	Destroy(companyID uuid.UUID, ids []uuid.UUID) error
}

type COAService struct {
	validate       *validator.Validate
	ICOARepository repository.ICOARepository
	// IJournalEntryRepository repository.IJournalEntryRepository
}

func NewCOAService(validate *validator.Validate, repo repository.ICOARepository) ICOAService {
	return &COAService{
		validate:       validate,
		ICOARepository: repo,
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
	if totalCount == 0 {
		return []entity.COA{}, 0, nil
	}
	totalPages := (totalCount + qp.PageSize - 1) / qp.PageSize
	if qp.Page > totalPages {
		return nil, totalCount, nil
	}
	results := make([]entity.COA, 0, len(coas))
	for _, coa := range coas {
		result := entity.COA{
			Id:          coa.Id,
			SubgroupId:  coa.SubgroupId,
			Code:        coa.Code,
			Name:        coa.Name,
			IsContra:    coa.IsContra,
			ControlType: coa.ControlType,
			Status:      coa.Status,
			CreatedAt:   coa.CreatedAt,
			UpdatedAt:   coa.UpdatedAt,
		}
		if coa.SubGroup != nil {
			result.SubGroup = &entity.COASubGroup{
				Id:        coa.SubGroup.Id,
				Code:      coa.SubGroup.Code,
				Name:      coa.SubGroup.Name,
				Status:    coa.SubGroup.Status,
				CreatedAt: coa.SubGroup.CreatedAt,
				UpdatedAt: coa.SubGroup.UpdatedAt,
				GroupId:   coa.SubGroup.GroupId,
			}
			if coa.SubGroup.Group != nil {
				result.SubGroup.Group = &entity.COAGroup{
					Id:        coa.SubGroup.Group.Id,
					Code:      coa.SubGroup.Group.Code,
					Name:      coa.SubGroup.Group.Name,
					Type:      coa.SubGroup.Group.Type,
					Status:    coa.SubGroup.Group.Status,
					CreatedAt: coa.SubGroup.Group.CreatedAt,
					UpdatedAt: coa.SubGroup.Group.UpdatedAt,
				}
			}
		}
		results = append(results, result)
	}
	return results, totalCount, nil
}

func (s *COAService) FindById(companyID, reqId uuid.UUID) (entity.COA, error) {
	coa, err := s.ICOARepository.FindById(companyID, reqId)
	if err != nil {
		return entity.COA{}, err
	}

	result := entity.COA{
		Id:          coa.Id,
		SubgroupId:  coa.SubgroupId,
		Code:        coa.Code,
		Name:        coa.Name,
		IsContra:    coa.IsContra,
		ControlType: coa.ControlType,
		Status:      coa.Status,
		CreatedAt:   coa.CreatedAt,
		UpdatedAt:   coa.UpdatedAt,
	}
	if coa.SubGroup != nil {
		result.SubGroup = &entity.COASubGroup{
			Id:        coa.SubGroup.Id,
			Code:      coa.SubGroup.Code,
			Name:      coa.SubGroup.Name,
			Status:    coa.SubGroup.Status,
			CreatedAt: coa.SubGroup.CreatedAt,
			UpdatedAt: coa.SubGroup.UpdatedAt,
			GroupId:   coa.SubGroup.GroupId,
		}
		if coa.SubGroup.Group != nil {
			result.SubGroup.Group = &entity.COAGroup{
				Id:        coa.SubGroup.Group.Id,
				Code:      coa.SubGroup.Group.Code,
				Name:      coa.SubGroup.Group.Name,
				Type:      coa.SubGroup.Group.Type,
				Status:    coa.SubGroup.Group.Status,
				CreatedAt: coa.SubGroup.Group.CreatedAt,
				UpdatedAt: coa.SubGroup.Group.UpdatedAt,
			}
		}
	}
	return result, nil
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

func (s *COAService) Delete(companyID uuid.UUID, reqId uuid.UUID) error {
	return s.ICOARepository.Delete(companyID, reqId)
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
