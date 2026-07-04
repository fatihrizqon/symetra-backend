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

type ICompanyService interface {
	Create(req request.CompanyCreateRequest, userId uuid.UUID) (entity.Company, error)
	FindAll(qp *util.QueryParams) ([]response.CompanyResponse, int, error)
	FindById(reqId uuid.UUID) (response.CompanyResponse, error)
	Update(req request.CompanyUpdateRequest) (response.CompanyResponse, error)
	Delete(reqId uuid.UUID) (response.CompanyResponse, error)
}

type CompanyService struct {
	ICompanyRepository repository.ICompanyRepository
	validate           *validator.Validate
}

func NewCompanyService(repo repository.ICompanyRepository, validate *validator.Validate) ICompanyService {
	return &CompanyService{
		ICompanyRepository: repo,
		validate:           validate,
	}
}

func (s *CompanyService) Create(req request.CompanyCreateRequest, userId uuid.UUID) (entity.Company, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.Company{}, err
	}

	c := entity.Company{
		Name:      req.Name,
		LegalName: req.LegalName,
		TaxID:     req.TaxID,
		Address:   req.Address,
		Phone:     req.Phone,
		Email:     req.Email,
		Industry:  req.Industry,
		Currency:  req.Currency,
		CreatedBy: userId,
	}

	err := s.ICompanyRepository.WithTransaction(func(txRepo repository.ICompanyRepository) error {
		var txErr error
		c, txErr = txRepo.Create(c)
		return txErr
	})

	return c, err
}

func (s *CompanyService) FindAll(qp *util.QueryParams) ([]response.CompanyResponse, int, error) {
	entities, totalCount, err := s.ICompanyRepository.FindAll(qp)
	if err != nil {
		return nil, 0, err
	}
	fmt.Println(entities)
	if totalCount == 0 {
		return []response.CompanyResponse{}, 0, nil
	}

	totalPages := (totalCount + qp.PageSize - 1) / qp.PageSize
	if qp.Page > totalPages {
		return nil, totalCount, nil
	}

	resps := make([]response.CompanyResponse, 0, len(entities))
	for _, c := range entities {
		resp := response.CompanyResponse{
			Id:        c.Id,
			Name:      c.Name,
			LegalName: c.LegalName,
			TaxID:     c.TaxID,
			Address:   c.Address,
			Phone:     c.Phone,
			Email:     c.Email,
			Industry:  c.Industry,
			Currency:  c.Currency,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			DeletedAt: c.DeletedAt,
		}
		resps = append(resps, resp)
	}

	return resps, totalCount, nil
}

func (s *CompanyService) FindById(reqId uuid.UUID) (response.CompanyResponse, error) {
	c, err := s.ICompanyRepository.FindById(reqId)
	if err != nil {
		return response.CompanyResponse{}, err
	}

	resp := response.CompanyResponse{
		Id:        c.Id,
		Name:      c.Name,
		LegalName: c.LegalName,
		TaxID:     c.TaxID,
		Address:   c.Address,
		Phone:     c.Phone,
		Email:     c.Email,
		Industry:  c.Industry,
		Currency:  c.Currency,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		DeletedAt: c.DeletedAt,
	}

	return resp, nil
}

func (s *CompanyService) Update(req request.CompanyUpdateRequest) (response.CompanyResponse, error) {
	var c entity.Company
	err := s.ICompanyRepository.WithTransaction(func(txRepo repository.ICompanyRepository) error {
		var err error
		c, err = txRepo.FindById(req.Id)
		if err != nil {
			return err
		}

		c.Name = req.Name
		c.LegalName = req.LegalName
		c.TaxID = req.TaxID
		c.Address = req.Address
		c.Phone = req.Phone
		c.Email = req.Email
		c.Industry = req.Industry
		c.Currency = req.Currency

		if err := txRepo.Update(c); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return response.CompanyResponse{}, err
	}

	return response.CompanyResponse{
		Id:        c.Id,
		Name:      c.Name,
		LegalName: c.LegalName,
		TaxID:     c.TaxID,
		Address:   c.Address,
		Phone:     c.Phone,
		Email:     c.Email,
		Industry:  c.Industry,
		Currency:  c.Currency,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		DeletedAt: c.DeletedAt,
	}, nil
}

func (s *CompanyService) Delete(reqId uuid.UUID) (response.CompanyResponse, error) {
	c, err := s.ICompanyRepository.FindById(reqId)
	if err != nil {
		return response.CompanyResponse{}, err
	}

	if err := s.ICompanyRepository.Delete(reqId); err != nil {
		return response.CompanyResponse{}, err
	}

	return response.CompanyResponse{
		Id:        c.Id,
		Name:      c.Name,
		LegalName: c.LegalName,
		TaxID:     c.TaxID,
		Address:   c.Address,
		Phone:     c.Phone,
		Email:     c.Email,
		Industry:  c.Industry,
		Currency:  c.Currency,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		DeletedAt: c.DeletedAt,
	}, nil
}
