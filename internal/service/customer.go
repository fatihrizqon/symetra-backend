package service

import (
	"errors"

	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/request"
	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ICustomerService interface {
	Create(companyID uuid.UUID, req request.CustomerCreateRequest) (entity.Customer, error)
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.Customer, int, error)
	FindById(companyID, id uuid.UUID) (entity.Customer, error)
	Update(companyID uuid.UUID, req request.CustomerUpdateRequest) (entity.Customer, error)
	Delete(companyID uuid.UUID, id uuid.UUID) error
	Destroy(companyID uuid.UUID, ids []uuid.UUID) error
}

type CustomerService struct {
	validate            *validator.Validate
	ICustomerRepository repository.ICustomerRepository
}

func NewCustomerService(validate *validator.Validate, repo repository.ICustomerRepository) ICustomerService {
	return &CustomerService{
		validate:            validate,
		ICustomerRepository: repo,
	}
}

func (s *CustomerService) Create(companyID uuid.UUID, req request.CustomerCreateRequest) (entity.Customer, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.Customer{}, err
	}

	customer := entity.Customer{
		CompanyId: companyID,
		Code:      req.Code,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		CoaId:     req.CoaId,
		Status:    1,
	}

	if err := s.ICustomerRepository.Create(&customer); err != nil {
		return entity.Customer{}, err
	}

	createdCustomer, err := s.ICustomerRepository.FindById(companyID, customer.Id)
	if err != nil {
		return entity.Customer{}, err
	}

	result := entity.Customer{
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
		result.Coa = &entity.COA{
			Id:          createdCustomer.Coa.Id,
			Code:        createdCustomer.Coa.Code,
			Name:        createdCustomer.Coa.Name,
			IsContra:    createdCustomer.Coa.IsContra,
			ControlType: createdCustomer.Coa.ControlType,
			Status:      createdCustomer.Coa.Status,
			CreatedAt:   createdCustomer.Coa.CreatedAt,
			UpdatedAt:   createdCustomer.Coa.UpdatedAt,
		}
	}
	return result, nil
}

func (s *CustomerService) FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.Customer, int, error) {
	customers, totalCount, err := s.ICustomerRepository.FindAll(companyID, qp)
	if err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return []entity.Customer{}, 0, nil
	}
	totalPages := (totalCount + qp.PageSize - 1) / qp.PageSize
	if qp.Page > totalPages {
		return nil, totalCount, nil
	}
	results := make([]entity.Customer, 0, len(customers))
	for _, customer := range customers {
		result := entity.Customer{
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
			result.Coa = &entity.COA{
				Id:       customer.Coa.Id,
				Code:     customer.Coa.Code,
				Name:     customer.Coa.Name,
				IsContra: customer.Coa.IsContra,
			}
		}
		results = append(results, result)
	}

	return results, totalCount, nil
}

func (s *CustomerService) FindById(companyID uuid.UUID, id uuid.UUID) (entity.Customer, error) {
	customer, err := s.ICustomerRepository.FindById(companyID, id)
	if err != nil {
		return entity.Customer{}, errors.New("customer not found")
	}
	result := entity.Customer{
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
		result.Coa = &entity.COA{
			Id:       customer.Coa.Id,
			Code:     customer.Coa.Code,
			Name:     customer.Coa.Name,
			IsContra: customer.Coa.IsContra,
		}
	}
	return result, nil
}

func (s *CustomerService) Update(companyID uuid.UUID, req request.CustomerUpdateRequest) (entity.Customer, error) {
	customer, err := s.ICustomerRepository.FindById(companyID, req.Id)
	if err != nil {
		return entity.Customer{}, errors.New("customer not found")
	}

	customer.Code = req.Code
	customer.Name = req.Name
	customer.Email = req.Email
	customer.Phone = req.Phone
	customer.Address = req.Address
	customer.CoaId = req.CoaId

	if err := s.ICustomerRepository.Update(&customer); err != nil {
		return entity.Customer{}, err
	}

	updatedCustomer, _ := s.ICustomerRepository.FindById(companyID, customer.Id)
	result := entity.Customer{
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
	return result, nil
}

func (s *CustomerService) Delete(companyID uuid.UUID, id uuid.UUID) error {
	return s.ICustomerRepository.Delete(companyID, id)
}

func (s *CustomerService) Destroy(companyID uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return errors.New("no ids provided")
	}
	return s.ICustomerRepository.BulkDestroy(companyID, ids)
}
