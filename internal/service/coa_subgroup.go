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
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]response.COASubGroupResponse, int, error)
	FindById(companyID, reqId uuid.UUID) (response.COASubGroupResponse, error)
	Update(companyID uuid.UUID, req request.COASubGroupUpdateRequest) (entity.COASubGroup, error)
	Delete(companyID, reqId uuid.UUID) (entity.COASubGroup, error)
	SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error)
	Destroy(companyID uuid.UUID, ids []uuid.UUID) error
}

type COASubGroupService struct {
	ICOASubGroupRepository repository.ICOASubGroupRepository
	validate               *validator.Validate
}

func NewCOASubGroupService(repo repository.ICOASubGroupRepository, validate *validator.Validate) ICOASubGroupService {
	return &COASubGroupService{ICOASubGroupRepository: repo, validate: validate}
}

func (s *COASubGroupService) Create(companyID uuid.UUID, req request.COASubGroupCreateRequest) (entity.COASubGroup, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.COASubGroup{}, err
	}
	sg := entity.COASubGroup{CompanyId: companyID, GroupId: req.GroupId, Code: req.Code, Name: req.Name}
	return s.ICOASubGroupRepository.Create(sg)
}

func (s *COASubGroupService) FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]response.COASubGroupResponse, int, error) {
	entities, totalCount, err := s.ICOASubGroupRepository.FindAll(companyID, qp)
	if err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return []response.COASubGroupResponse{}, 0, nil
	}
	totalPages := (totalCount + qp.PageSize - 1) / qp.PageSize
	if qp.Page > totalPages {
		return nil, totalCount, nil
	}
	resps := make([]response.COASubGroupResponse, 0, len(entities))
	for _, v := range entities {
		resps = append(resps, mapCOASubGroup(v))
	}
	return resps, totalCount, nil
}

func (s *COASubGroupService) FindById(companyID, reqId uuid.UUID) (response.COASubGroupResponse, error) {
	result, err := s.ICOASubGroupRepository.FindById(companyID, reqId)
	if err != nil {
		return response.COASubGroupResponse{}, err
	}
	return mapCOASubGroup(result), nil
}

func (s *COASubGroupService) Update(companyID uuid.UUID, req request.COASubGroupUpdateRequest) (entity.COASubGroup, error) {
	sg, err := s.ICOASubGroupRepository.FindById(companyID, req.Id)
	if err != nil {
		return sg, err
	}
	sg.GroupId = req.GroupId
	sg.Code = req.Code
	sg.Name = req.Name
	return sg, s.ICOASubGroupRepository.Update(sg)
}

func (s *COASubGroupService) Delete(companyID, reqId uuid.UUID) (entity.COASubGroup, error) {
	sg, err := s.ICOASubGroupRepository.FindById(companyID, reqId)
	if err != nil {
		return sg, err
	}
	hasCOAs, err := s.ICOASubGroupRepository.HasCOAs(companyID, reqId)
	if err != nil {
		return sg, err
	}
	if hasCOAs {
		return sg, fmt.Errorf("HAS_CHILD_RECORDS")
	}
	return sg, s.ICOASubGroupRepository.Delete(companyID, reqId)
}

func (s *COASubGroupService) SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error) {
	entities, total, err := s.ICOASubGroupRepository.SelectDropdownList(companyID, qp)
	if err != nil {
		return nil, 0, err
	}
	resps := make([]response.SelectDropdownListResponse, 0, len(entities))
	for _, v := range entities {
		resps = append(resps, response.SelectDropdownListResponse{Value: v.Id, Label: v.Name})
	}
	return resps, total, nil
}

func mapCOASubGroup(v entity.COASubGroup) response.COASubGroupResponse {
	r := response.COASubGroupResponse{
		Id:        v.Id,
		GroupId:   v.GroupId,
		Code:      v.Code,
		Name:      v.Name,
		Status:    v.Status,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
	if v.Group != nil {
		r.Group = &response.COAGroupResponse{
			Id:            v.Group.Id,
			Code:          v.Group.Code,
			Name:          v.Group.Name,
			Status:        v.Group.Status,
			CreatedAt:     v.Group.CreatedAt,
			UpdatedAt:     v.Group.UpdatedAt,
		}
	}
	return r
}

func (s *COASubGroupService) Destroy(companyID uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return fmt.Errorf("no ids provided")
	}
	return s.ICOASubGroupRepository.BulkDestroy(companyID, ids)
}

