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

type IVendorService interface {
	Create(companyID uuid.UUID, req request.VendorCreateRequest) (entity.Vendor, error)
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.Vendor, int, error)
	FindById(companyID, id uuid.UUID) (entity.Vendor, error)
	Update(companyID uuid.UUID, req request.VendorUpdateRequest) (entity.Vendor, error)
	Delete(companyID, id uuid.UUID) error
	Destroy(companyID uuid.UUID, ids []uuid.UUID) error
}

type VendorService struct {
	validate          *validator.Validate
	IVendorRepository repository.IVendorRepository
}

func NewVendorService(validate *validator.Validate, repo repository.IVendorRepository) IVendorService {
	return &VendorService{
		validate:          validate,
		IVendorRepository: repo,
	}
}

func (s *VendorService) Create(companyID uuid.UUID, req request.VendorCreateRequest) (entity.Vendor, error) {
	vendor := entity.Vendor{
		CompanyId: companyID,
		Code:      req.Code,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		CoaId:     req.CoaId,
		Status:    1,
	}
	if err := s.IVendorRepository.Create(&vendor); err != nil {
		return entity.Vendor{}, err
	}

	createdVendor, _ := s.IVendorRepository.FindById(companyID, vendor.Id)

	resp := entity.Vendor{
		Id:        createdVendor.Id,
		CompanyId: createdVendor.CompanyId,
		Code:      createdVendor.Code,
		Name:      createdVendor.Name,
		Email:     createdVendor.Email,
		Phone:     createdVendor.Phone,
		Address:   createdVendor.Address,
		CoaId:     createdVendor.CoaId,
		Status:    createdVendor.Status,
		CreatedAt: createdVendor.CreatedAt,
		UpdatedAt: createdVendor.UpdatedAt,
	}

	if createdVendor.Coa != nil {
		resp.Coa = &entity.COA{
			Id:       createdVendor.Coa.Id,
			Code:     createdVendor.Coa.Code,
			Name:     createdVendor.Coa.Name,
			IsContra: createdVendor.Coa.IsContra,
		}
	}
	return resp, nil
}

func (s *VendorService) FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.Vendor, int, error) {
	vendors, total, err := s.IVendorRepository.FindAll(companyID, qp)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []entity.Vendor{}, 0, nil
	}

	totalPages := (int(total) + qp.PageSize - 1) / qp.PageSize
	if qp.Page > totalPages {
		return nil, int(total), nil
	}

	resps := make([]entity.Vendor, 0, len(vendors))
	for _, v := range vendors {
		resp := entity.Vendor{
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
			resp.Coa = &entity.COA{
				Id:       v.Coa.Id,
				Code:     v.Coa.Code,
				Name:     v.Coa.Name,
				IsContra: v.Coa.IsContra,
			}
		}
		resps = append(resps, resp)
	}

	return resps, int(total), nil
}

func (s *VendorService) FindById(companyID, id uuid.UUID) (entity.Vendor, error) {
	vendor, err := s.IVendorRepository.FindById(companyID, id)
	if err != nil {
		return entity.Vendor{}, errors.New("vendor not found")
	}
	resp := entity.Vendor{
		Id:        vendor.Id,
		CompanyId: vendor.CompanyId,
		Code:      vendor.Code,
		Name:      vendor.Name,
		Email:     vendor.Email,
		Phone:     vendor.Phone,
		Address:   vendor.Address,
		CoaId:     vendor.CoaId,
		Status:    vendor.Status,
		CreatedAt: vendor.CreatedAt,
		UpdatedAt: vendor.UpdatedAt,
	}

	if vendor.Coa != nil {
		resp.Coa = &entity.COA{
			Id:       vendor.Coa.Id,
			Code:     vendor.Coa.Code,
			Name:     vendor.Coa.Name,
			IsContra: vendor.Coa.IsContra,
		}
	}
	return resp, nil
}

func (s *VendorService) Update(companyID uuid.UUID, req request.VendorUpdateRequest) (entity.Vendor, error) {
	vendor, err := s.IVendorRepository.FindById(companyID, req.Id)
	if err != nil {
		return entity.Vendor{}, errors.New("vendor not found")
	}

	vendor.Code = req.Code
	vendor.Name = req.Name
	vendor.Email = req.Email
	vendor.Phone = req.Phone
	vendor.Address = req.Address
	vendor.CoaId = req.CoaId

	if err := s.IVendorRepository.Update(&vendor); err != nil {
		return entity.Vendor{}, err
	}
	updatedVendor, _ := s.IVendorRepository.FindById(companyID, vendor.Id)
	resp := entity.Vendor{
		Id:        updatedVendor.Id,
		CompanyId: updatedVendor.CompanyId,
		Code:      updatedVendor.Code,
		Name:      updatedVendor.Name,
		Email:     updatedVendor.Email,
		Phone:     updatedVendor.Phone,
		Address:   updatedVendor.Address,
		CoaId:     updatedVendor.CoaId,
		Status:    updatedVendor.Status,
		CreatedAt: updatedVendor.CreatedAt,
		UpdatedAt: updatedVendor.UpdatedAt,
	}

	if updatedVendor.Coa != nil {
		resp.Coa = &entity.COA{
			Id:       updatedVendor.Coa.Id,
			Code:     updatedVendor.Coa.Code,
			Name:     updatedVendor.Coa.Name,
			IsContra: updatedVendor.Coa.IsContra,
		}
	}
	return resp, nil
}

func (s *VendorService) Delete(companyID, id uuid.UUID) error {
	return s.IVendorRepository.Delete(companyID, id)
}

func (s *VendorService) Destroy(companyID uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return errors.New("no ids provided")
	}
	return s.IVendorRepository.BulkDestroy(companyID, ids)
}
