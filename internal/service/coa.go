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
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]response.COAResponse, int, error)
	FindById(companyID, reqId uuid.UUID) (response.COAResponse, error)
	Update(companyID uuid.UUID, req request.COAUpdateRequest) (entity.COA, error)
	Delete(companyID, reqId uuid.UUID) (entity.COA, error)
	SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error)
}

type COAService struct {
	ICOARepository          repository.ICOARepository
	// IJournalEntryRepository repository.IJournalEntryRepository
	validate                *validator.Validate
}

func NewCOAService(repo repository.ICOARepository, validate *validator.Validate) ICOAService {
	return &COAService{ICOARepository: repo, validate: validate}
}

func (s *COAService) Create(companyID uuid.UUID, req request.COACreateRequest) (entity.COA, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.COA{}, err
	}
	c := entity.COA{
		CompanyId:  companyID,
		SubgroupId: req.SubgroupId,
		Code:       req.Code,
		Name:       req.Name,
		// CurrencyCode: req.CurrencyCode,
		// BUG NOTE 06032026: Undefined CurrencyCode di COACreateRequest
		Active: true,
	}
	if c.CurrencyCode == "" {
		c.CurrencyCode = "IDR"
	}
	return s.ICOARepository.Create(c)
}

func (s *COAService) FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]response.COAResponse, int, error) {
	entities, totalCount, err := s.ICOARepository.FindAll(companyID, qp)
	if err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return []response.COAResponse{}, 0, nil
	}
	totalPages := (totalCount + qp.PageSize - 1) / qp.PageSize
	if qp.Page > totalPages {
		return nil, totalCount, nil
	}
	resps := make([]response.COAResponse, 0, len(entities))
	for _, v := range entities {
		resps = append(resps, mapCOA(v))
	}
	return resps, totalCount, nil
}

func (s *COAService) FindById(companyID, reqId uuid.UUID) (response.COAResponse, error) {
	result, err := s.ICOARepository.FindById(companyID, reqId)
	if err != nil {
		return response.COAResponse{}, err
	}
	return mapCOA(result), nil
}

func (s *COAService) Update(companyID uuid.UUID, req request.COAUpdateRequest) (entity.COA, error) {
	c, err := s.ICOARepository.FindById(companyID, req.Id)
	if err != nil {
		return c, err
	}
	c.SubgroupId = req.SubgroupId
	c.Code = req.Code
	c.Name = req.Name
	// BUG NOTE 06032026: Undefined CurrencyCode di COACreateRequest
	// if req.CurrencyCode != "" {
	// 	c.CurrencyCode = req.CurrencyCode
	// }
	return c, s.ICOARepository.Update(c)
}

func (s *COAService) Delete(companyID, reqId uuid.UUID) (entity.COA, error) {
	c, err := s.ICOARepository.FindById(companyID, reqId)
	if err != nil {
		return c, err
	}
	// Guard: cannot delete if used in transactions
	// if s.IJournalEntryRepository != nil {
	// 	hasTransactions, err := s.IJournalEntryRepository.HasTransactions(companyID, reqId)
	// 	if err != nil {
	// 		return c, err
	// 	}
	// 	if hasTransactions {
	// 		return c, errors.New("cannot delete an account that has been used in journal entries")
	// 	}
	// }
	return c, s.ICOARepository.Delete(companyID, reqId)
}

func (s *COAService) SelectDropdownList(companyID uuid.UUID, qp *util.QueryParams) ([]response.SelectDropdownListResponse, int, error) {
	entities, total, err := s.ICOARepository.SelectDropdownList(companyID, qp)
	if err != nil {
		return nil, 0, err
	}
	resps := make([]response.SelectDropdownListResponse, 0, len(entities))
	for _, v := range entities {
		resps = append(resps, response.SelectDropdownListResponse{Value: v.Id, Label: fmt.Sprintf("%s - %s", v.Code, v.Name)})
	}
	return resps, total, nil
}

func mapCOA(v entity.COA) response.COAResponse {
	r := response.COAResponse{
		Id:         v.Id,
		SubgroupId: v.SubgroupId,
		Code:       v.Code,
		Name:       v.Name,
		Status:     v.Status,
		CreatedAt:  v.CreatedAt,
		UpdatedAt:  v.UpdatedAt,
	}
	if v.SubGroup != nil {
		r.SubGroup = &response.COASubGroupResponse{
			Id:        v.SubGroup.Id,
			Code:      v.SubGroup.Code,
			Name:      v.SubGroup.Name,
			Status:    v.SubGroup.Status,
			CreatedAt: v.SubGroup.CreatedAt,
			UpdatedAt: v.SubGroup.UpdatedAt,
			GroupId:   v.SubGroup.GroupId,
		}
		if v.SubGroup.Group != nil {
			r.SubGroup.Group = &response.COAGroupResponse{
				Id:            v.SubGroup.Group.Id,
				Code:          v.SubGroup.Group.Code,
				Name:          v.SubGroup.Group.Name,
				Type:          string(v.SubGroup.Group.Type),
				NormalBalance: v.SubGroup.Group.NormalBalance,
				Status:        v.SubGroup.Group.Status,
				CreatedAt:     v.SubGroup.Group.CreatedAt,
				UpdatedAt:     v.SubGroup.Group.UpdatedAt,
			}
		}
	}
	return r
}
