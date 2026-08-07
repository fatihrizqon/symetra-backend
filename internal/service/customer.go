package service

import (
	"errors"

	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/request"
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/response"
	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/google/uuid"
)

type ICustomerService interface {
	Create(companyId uuid.UUID, req request.CustomerCreateRequest) (response.CustomerResponse, error)
	FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]response.CustomerResponse, int, error)
	FindById(companyId, id uuid.UUID) (response.CustomerResponse, error)
	Update(companyId, id uuid.UUID, req request.CustomerUpdateRequest) (response.CustomerResponse, error)
	Delete(companyId, id uuid.UUID) error
	Destroy(companyId uuid.UUID, ids []uuid.UUID) error
}

type CustomerService struct {
	repo repository.ICustomerRepository
}

func NewCustomerService(repo repository.ICustomerRepository) ICustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) Create(companyId uuid.UUID, req request.CustomerCreateRequest) (response.CustomerResponse, error) {
	customer := entity.Customer{
		CompanyId: companyId,
		Code:      req.Code,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		CoaId:     req.CoaId,
		Status:    1,
	}
	if err := s.repo.Create(&customer); err != nil {
		return response.CustomerResponse{}, err
	}
	createdCustomer, _ := s.repo.FindById(companyId, customer.Id)
		resp := response.CustomerResponse{
		Id:        createdCustomer.Id,
		CompanyId: createdCustomer.CompanyId,
		Code:      createdCustomer.Code,
		Name:      createdCustomer.Name,
		Email:     createdCustomer.Email,
		Phone:     createdCustomer.Phone,
		Address:   createdCustomer.Address,
		CoaId:     createdCustomer.CoaId,
		Status:    createdCustomer.Status,
		CreatedAt: createdCustomer.CreatedAt,
		UpdatedAt: createdCustomer.UpdatedAt,
	}

	if createdCustomer.Coa != nil {
		resp.Coa = &response.COAResponse{
			Id:            createdCustomer.Coa.Id,
			Code:          createdCustomer.Coa.Code,
			Name:          createdCustomer.Coa.Name,
			IsContra:      createdCustomer.Coa.IsContra,
			NormalBalance: createdCustomer.Coa.GetAbsoluteNormalBalance(),
		}
	}
	return resp, nil
}

func (s *CustomerService) FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]response.CustomerResponse, int, error) {
	customers, total, err := s.repo.FindAll(companyId, qp)
	if err != nil {
		return nil, 0, err
	}
		resps := make([]response.CustomerResponse, 0, len(customers))
	for _, v := range customers {
	resp := response.CustomerResponse{
				Id:        v.Id,
				CompanyId: v.CompanyId,
				Code:      v.Code,
				Name:      v.Name,
				Email:     v.Email,
				Phone:     v.Phone,
				Address:   v.Address,
				CoaId:     v.CoaId,
				Status:    v.Status,
				CreatedAt: v.CreatedAt,
				UpdatedAt: v.UpdatedAt,
			}
		
			if v.Coa != nil {
				resp.Coa = &response.COAResponse{
					Id:            v.Coa.Id,
					Code:          v.Coa.Code,
					Name:          v.Coa.Name,
					IsContra:      v.Coa.IsContra,
					NormalBalance: v.Coa.GetAbsoluteNormalBalance(),
				}
			}
		resps = append(resps, resp)
	}

	return resps, int(total), nil
}

func (s *CustomerService) FindById(companyId, id uuid.UUID) (response.CustomerResponse, error) {
	customer, err := s.repo.FindById(companyId, id)
	if err != nil {
		return response.CustomerResponse{}, errors.New("customer not found")
	}
		resp := response.CustomerResponse{
		Id:        customer.Id,
		CompanyId: customer.CompanyId,
		Code:      customer.Code,
		Name:      customer.Name,
		Email:     customer.Email,
		Phone:     customer.Phone,
		Address:   customer.Address,
		CoaId:     customer.CoaId,
		Status:    customer.Status,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}

	if customer.Coa != nil {
		resp.Coa = &response.COAResponse{
			Id:            customer.Coa.Id,
			Code:          customer.Coa.Code,
			Name:          customer.Coa.Name,
			IsContra:      customer.Coa.IsContra,
			NormalBalance: customer.Coa.GetAbsoluteNormalBalance(),
		}
	}
	return resp, nil
}

func (s *CustomerService) Update(companyId, id uuid.UUID, req request.CustomerUpdateRequest) (response.CustomerResponse, error) {
	customer, err := s.repo.FindById(companyId, id)
	if err != nil {
		return response.CustomerResponse{}, errors.New("customer not found")
	}

	customer.Code = req.Code
	customer.Name = req.Name
	customer.Email = req.Email
	customer.Phone = req.Phone
	customer.Address = req.Address
	customer.CoaId = req.CoaId

	if err := s.repo.Update(&customer); err != nil {
		return response.CustomerResponse{}, err
	}
	updatedCustomer, _ := s.repo.FindById(companyId, customer.Id)
		resp := response.CustomerResponse{
		Id:        updatedCustomer.Id,
		CompanyId: updatedCustomer.CompanyId,
		Code:      updatedCustomer.Code,
		Name:      updatedCustomer.Name,
		Email:     updatedCustomer.Email,
		Phone:     updatedCustomer.Phone,
		Address:   updatedCustomer.Address,
		CoaId:     updatedCustomer.CoaId,
		Status:    updatedCustomer.Status,
		CreatedAt: updatedCustomer.CreatedAt,
		UpdatedAt: updatedCustomer.UpdatedAt,
	}

	if updatedCustomer.Coa != nil {
		resp.Coa = &response.COAResponse{
			Id:            updatedCustomer.Coa.Id,
			Code:          updatedCustomer.Coa.Code,
			Name:          updatedCustomer.Coa.Name,
			IsContra:      updatedCustomer.Coa.IsContra,
			NormalBalance: updatedCustomer.Coa.GetAbsoluteNormalBalance(),
		}
	}
	return resp, nil
}

func (s *CustomerService) Delete(companyId, id uuid.UUID) error {
	return s.repo.Delete(companyId, id)
}

func (s *CustomerService) Destroy(companyId uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return errors.New("no ids provided")
	}
	return s.repo.BulkDestroy(companyId, ids)
}
