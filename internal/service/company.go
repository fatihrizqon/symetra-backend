package service

import (
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
	FindMyCompanies(userID uuid.UUID) ([]response.MyCompanyResponse, error)
	FindMembersByCompany(companyId uuid.UUID) ([]response.CompanyMemberResponse, error)
	AssignMember(req request.AssignMemberRequest, invitedBy uuid.UUID) (response.CompanyMemberResponse, error)
	UpdateMemberRole(req request.UpdateMemberRoleRequest) error
	RemoveMember(companyId uuid.UUID, userId uuid.UUID) error
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

		_, txErr = txRepo.AssignMember(entity.CompanyMember{
			CompanyId: c.Id,
			UserId:    userId,
			Role:      "owner",
		})

		return txErr
	})

	return c, err
}

func (s *CompanyService) FindAll(qp *util.QueryParams) ([]response.CompanyResponse, int, error) {
	entities, totalCount, err := s.ICompanyRepository.FindAll(qp)
	if err != nil {
		return nil, 0, err
	}

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

func (s *CompanyService) FindMyCompanies(userId uuid.UUID) ([]response.MyCompanyResponse, error) {
	entities, err := s.ICompanyRepository.FindMyCompanies(userId)
	if err != nil {
		return nil, err
	}

	resps := make([]response.MyCompanyResponse, 0, len(entities))
	for _, m := range entities {
		c := m.Company
		resps = append(resps, response.MyCompanyResponse{
			CompanyResponse: response.CompanyResponse{
				Id:        c.Id,
				Name:      c.Name,
				LegalName: c.LegalName,
				TaxID:     c.TaxID,
				Address:   c.Address,
				Phone:     c.Phone,
				Email:     c.Email,
				Industry:  c.Industry,
				Currency:  c.Currency,
				CreatedBy: c.CreatedBy,
				CreatedAt: c.CreatedAt,
				UpdatedAt: c.UpdatedAt,
				DeletedAt: c.DeletedAt,
			},
			Role: m.Role,
		})
	}
	return resps, nil

}

func (s *CompanyService) FindMembersByCompany(companyId uuid.UUID) ([]response.CompanyMemberResponse, error) {
	entities, err := s.ICompanyRepository.FindMembersByCompany(companyId)
	if err != nil {
		return nil, err
	}

	resps := make([]response.CompanyMemberResponse, 0, len(entities))
	for _, m := range entities {
		resps = append(resps, response.CompanyMemberResponse{
			Id:        m.Id,
			CompanyId: m.CompanyId,
			UserId:    m.UserId,
			Username:  m.User.Username,
			Name:      m.User.Name,
			Email:     m.User.Email,
			Role:      m.Role,
			CreatedAt: m.JoinedAt,
			UpdatedAt: m.UpdatedAt,
			Company: response.CompanyResponse{
				Id:        m.Company.Id,
				Name:      m.Company.Name,
				LegalName: m.Company.LegalName,
				TaxID:     m.Company.TaxID,
				Address:   m.Company.Address,
				Phone:     m.Company.Phone,
				Email:     m.Company.Email,
				Industry:  m.Company.Industry,
				Currency:  m.Company.Currency,
				CreatedBy: m.Company.CreatedBy,
				CreatedAt: m.Company.CreatedAt,
				UpdatedAt: m.Company.UpdatedAt,
				DeletedAt: m.Company.DeletedAt,
			},
		})
	}
	return resps, nil
}

func (s *CompanyService) AssignMember(req request.AssignMemberRequest, invitedBy uuid.UUID) (response.CompanyMemberResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.CompanyMemberResponse{}, err
	}

	member := entity.CompanyMember{
		CompanyId: req.CompanyId,
		UserId:    req.UserID,
		Role:      req.Role,
		InvitedBy: invitedBy,
	}

	created, err := s.ICompanyRepository.AssignMember(member)
	if err != nil {
		return response.CompanyMemberResponse{}, err
	}

	return response.CompanyMemberResponse{
		Id:        created.Id,
		CompanyId: created.CompanyId,
		UserId:    created.UserId,
		Username:  created.User.Username,
		Name:      created.User.Name,
		Email:     created.User.Email,
		Role:      created.Role,
		CreatedAt: created.JoinedAt,
		UpdatedAt: created.UpdatedAt,
		Company: response.CompanyResponse{
			Id:        created.Company.Id,
			Name:      created.Company.Name,
			LegalName: created.Company.LegalName,
			TaxID:     created.Company.TaxID,
			Address:   created.Company.Address,
			Phone:     created.Company.Phone,
			Email:     created.Company.Email,
			Industry:  created.Company.Industry,
			Currency:  created.Company.Currency,
			CreatedBy: created.Company.CreatedBy,
			CreatedAt: created.Company.CreatedAt,
			UpdatedAt: created.Company.UpdatedAt,
			DeletedAt: created.Company.DeletedAt,
		},
	}, nil
}

func (s *CompanyService) UpdateMemberRole(req request.UpdateMemberRoleRequest) error {
	if err := s.validate.Struct(req); err != nil {
		return err
	}
	return s.ICompanyRepository.UpdateMemberRole(req.CompanyId, req.UserID, req.Role)
}

func (s *CompanyService) RemoveMember(companyId uuid.UUID, userId uuid.UUID) error {
	return s.ICompanyRepository.RemoveMember(companyId, userId)
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
