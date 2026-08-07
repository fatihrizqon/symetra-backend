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

type IVendorService interface {
	Create(companyId uuid.UUID, req request.VendorCreateRequest) (response.VendorResponse, error)
	FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]response.VendorResponse, int, error)
	FindById(companyId, id uuid.UUID) (response.VendorResponse, error)
	Update(companyId, id uuid.UUID, req request.VendorUpdateRequest) (response.VendorResponse, error)
	Delete(companyId, id uuid.UUID) error
	Destroy(companyId uuid.UUID, ids []uuid.UUID) error
}

type VendorService struct {
	repo repository.IVendorRepository
}

func NewVendorService(repo repository.IVendorRepository) IVendorService {
	return &VendorService{repo: repo}
}

func (s *VendorService) Create(companyId uuid.UUID, req request.VendorCreateRequest) (response.VendorResponse, error) {
	vendor := entity.Vendor{
		CompanyId: companyId,
		Code:      req.Code,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		CoaId:     req.CoaId,
		Status:    1,
	}
	if err := s.repo.Create(&vendor); err != nil {
		return response.VendorResponse{}, err
	}
	createdVendor, _ := s.repo.FindById(companyId, vendor.Id)
	
	resp := response.VendorResponse{
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
		resp.Coa = &response.COAResponse{
			Id:            createdVendor.Coa.Id,
			Code:          createdVendor.Coa.Code,
			Name:          createdVendor.Coa.Name,
			IsContra:      createdVendor.Coa.IsContra,
			NormalBalance: createdVendor.Coa.GetAbsoluteNormalBalance(),
		}
	}
	return resp, nil
}

func (s *VendorService) FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]response.VendorResponse, int, error) {
	vendors, total, err := s.repo.FindAll(companyId, qp)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []response.VendorResponse{}, 0, nil
	}

	totalPages := (int(total) + qp.PageSize - 1) / qp.PageSize
	if qp.Page > totalPages {
		return nil, int(total), nil
	}

	resps := make([]response.VendorResponse, 0, len(vendors))
	for _, v := range vendors {
		resp := response.VendorResponse{
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

func (s *VendorService) FindById(companyId, id uuid.UUID) (response.VendorResponse, error) {
	vendor, err := s.repo.FindById(companyId, id)
	if err != nil {
		return response.VendorResponse{}, errors.New("vendor not found")
	}
	resp := response.VendorResponse{
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
		resp.Coa = &response.COAResponse{
			Id:            vendor.Coa.Id,
			Code:          vendor.Coa.Code,
			Name:          vendor.Coa.Name,
			IsContra:      vendor.Coa.IsContra,
			NormalBalance: vendor.Coa.GetAbsoluteNormalBalance(),
		}
	}
	return resp, nil
}

func (s *VendorService) Update(companyId, id uuid.UUID, req request.VendorUpdateRequest) (response.VendorResponse, error) {
	vendor, err := s.repo.FindById(companyId, id)
	if err != nil {
		return response.VendorResponse{}, errors.New("vendor not found")
	}

	vendor.Code = req.Code
	vendor.Name = req.Name
	vendor.Email = req.Email
	vendor.Phone = req.Phone
	vendor.Address = req.Address
	vendor.CoaId = req.CoaId

	if err := s.repo.Update(&vendor); err != nil {
		return response.VendorResponse{}, err
	}
	updatedVendor, _ := s.repo.FindById(companyId, vendor.Id)
	resp := response.VendorResponse{
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
		resp.Coa = &response.COAResponse{
			Id:            updatedVendor.Coa.Id,
			Code:          updatedVendor.Coa.Code,
			Name:          updatedVendor.Coa.Name,
			IsContra:      updatedVendor.Coa.IsContra,
			NormalBalance: updatedVendor.Coa.GetAbsoluteNormalBalance(),
		}
	}
	return resp, nil
}

func (s *VendorService) Delete(companyId, id uuid.UUID) error {
	return s.repo.Delete(companyId, id)
}

func (s *VendorService) Destroy(companyId uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return errors.New("no ids provided")
	}
	return s.repo.BulkDestroy(companyId, ids)
}
