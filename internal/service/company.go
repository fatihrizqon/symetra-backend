package service

import (
	"context"
	"errors"

	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/request"
	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ICompanyService interface {
	Create(req request.CompanyCreateRequest, userId uuid.UUID) (entity.Company, error)
	FindAll(qp *util.QueryParams) ([]entity.Company, int, error)
	FindById(id uuid.UUID) (entity.Company, error)
	FindMyCompanies(userID uuid.UUID) ([]entity.Company, error)
	FindMembersByCompany(companyID uuid.UUID) ([]entity.CompanyMember, error)
	AssignMember(req request.AssignMemberRequest, invitedBy uuid.UUID) (entity.CompanyMember, error)
	UpdateMemberRole(req request.UpdateMemberRoleRequest) error
	RemoveMember(companyID uuid.UUID, userId uuid.UUID) error
	Update(req request.CompanyUpdateRequest) (entity.Company, error)
	Delete(id uuid.UUID) error
	SelectCompany(sessionID uuid.UUID, companyID uuid.UUID, userID uuid.UUID) (entity.Company, error)
	GetActiveCompany(sessionID uuid.UUID) (entity.Company, error)
	Destroy(ids []uuid.UUID) error
}

type CompanyService struct {
	validate           *validator.Validate
	ICompanyRepository repository.ICompanyRepository
	ITokenRepository   repository.ITokenRepository
}

func NewCompanyService(validate *validator.Validate, ICompanyRepository repository.ICompanyRepository, ITokenRepository repository.ITokenRepository) ICompanyService {
	return &CompanyService{
		validate:           validate,
		ICompanyRepository: ICompanyRepository,
		ITokenRepository:   ITokenRepository,
	}
}

func (s *CompanyService) Create(req request.CompanyCreateRequest, userId uuid.UUID) (entity.Company, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.Company{}, err
	}

	company := entity.Company{
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
		company, txErr = txRepo.Create(company)
		if txErr != nil {
			return txErr
		}

		_, txErr = txRepo.AssignMember(entity.CompanyMember{
			CompanyId: company.Id,
			UserId:    userId,
			Role:      "owner",
		})

		return txErr
	})

	return company, err
}

func (s *CompanyService) FindAll(qp *util.QueryParams) ([]entity.Company, int, error) {
	companies, totalCount, err := s.ICompanyRepository.FindAll(qp)
	if err != nil {
		return nil, 0, err
	}

	return companies, totalCount, nil
}

func (s *CompanyService) FindById(id uuid.UUID) (entity.Company, error) {
	company, err := s.ICompanyRepository.FindById(id)
	if err != nil {
		return entity.Company{}, err
	}

	return company, nil
}

func (s *CompanyService) FindMyCompanies(userId uuid.UUID) ([]entity.Company, error) {
	companies, err := s.ICompanyRepository.FindMyCompanies(userId)
	if err != nil {
		return nil, err
	}

	results := make([]entity.Company, 0, len(companies))

	for _, myCompany := range companies {
		company := myCompany.Company

		results = append(results, entity.Company{
			Id:        company.Id,
			Name:      company.Name,
			LegalName: company.LegalName,
			TaxID:     company.TaxID,
			Address:   company.Address,
			Phone:     company.Phone,
			Email:     company.Email,
			Industry:  company.Industry,
			Currency:  company.Currency,
			CreatedBy: company.CreatedBy,
			CreatedAt: company.CreatedAt,
			UpdatedAt: company.UpdatedAt,
		})
	}

	return results, nil

}

func (s *CompanyService) FindMembersByCompany(companyID uuid.UUID) ([]entity.CompanyMember, error) {
	companies, err := s.ICompanyRepository.FindMembersByCompany(companyID)
	if err != nil {
		return nil, err
	}

	results := make([]entity.CompanyMember, 0, len(companies))
	for _, member := range companies {
		results = append(results, entity.CompanyMember{
			Id:        member.Id,
			CompanyId: member.CompanyId,
			UserId:    member.UserId,
			Role:      member.Role,
			UpdatedAt: member.UpdatedAt,
			Company: entity.Company{
				Id:        member.Company.Id,
				Name:      member.Company.Name,
				LegalName: member.Company.LegalName,
				TaxID:     member.Company.TaxID,
				Address:   member.Company.Address,
				Phone:     member.Company.Phone,
				Email:     member.Company.Email,
				Industry:  member.Company.Industry,
				Currency:  member.Company.Currency,
				CreatedBy: member.Company.CreatedBy,
				CreatedAt: member.Company.CreatedAt,
				UpdatedAt: member.Company.UpdatedAt,
			},
		})
	}
	return results, nil
}

func (s *CompanyService) AssignMember(req request.AssignMemberRequest, invitedBy uuid.UUID) (entity.CompanyMember, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.CompanyMember{}, err
	}

	if req.Role == "owner" {
		return entity.CompanyMember{}, errors.New("cannot assign owner role. A company can only have one owner")
	}

	member := entity.CompanyMember{
		CompanyId: req.CompanyId,
		UserId:    req.UserID,
		Role:      req.Role,
		InvitedBy: invitedBy,
	}

	created, err := s.ICompanyRepository.AssignMember(member)
	if err != nil {
		return entity.CompanyMember{}, err
	}

	return entity.CompanyMember{
		Id:        created.Id,
		CompanyId: created.CompanyId,
		UserId:    created.UserId,
		Role:      created.Role,
		UpdatedAt: created.UpdatedAt,
		Company: entity.Company{
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
		},
	}, nil
}

func (s *CompanyService) UpdateMemberRole(req request.UpdateMemberRoleRequest) error {
	if err := s.validate.Struct(req); err != nil {
		return err
	}

	if req.Role == "owner" {
		return errors.New("cannot assign owner role. A company can only have one owner")
	}

	member, err := s.ICompanyRepository.FindMember(context.Background(), req.CompanyId, req.UserID)
	if err != nil {
		return err
	}

	if member.Role == "owner" {
		return errors.New("cannot update the role of the company owner")
	}

	return s.ICompanyRepository.UpdateMemberRole(req.CompanyId, req.UserID, req.Role)
}

func (s *CompanyService) RemoveMember(companyID uuid.UUID, userId uuid.UUID) error {
	member, err := s.ICompanyRepository.FindMember(context.Background(), companyID, userId)
	if err != nil {
		return err
	}

	if member.Role == "owner" {
		return errors.New("cannot remove the company owner")
	}

	return s.ICompanyRepository.RemoveMember(companyID, userId)
}

func (s *CompanyService) Update(req request.CompanyUpdateRequest) (entity.Company, error) {
	var company entity.Company
	err := s.ICompanyRepository.WithTransaction(func(txRepo repository.ICompanyRepository) error {
		var err error
		company, err = txRepo.FindById(req.Id)
		if err != nil {
			return err
		}

		company.Name = req.Name
		company.LegalName = req.LegalName
		company.TaxID = req.TaxID
		company.Address = req.Address
		company.Phone = req.Phone
		company.Email = req.Email
		company.Industry = req.Industry
		company.Currency = req.Currency

		if err := txRepo.Update(company); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return entity.Company{}, err
	}

	return entity.Company{
		Id:        company.Id,
		Name:      company.Name,
		LegalName: company.LegalName,
		TaxID:     company.TaxID,
		Address:   company.Address,
		Phone:     company.Phone,
		Email:     company.Email,
		Industry:  company.Industry,
		Currency:  company.Currency,
		CreatedAt: company.CreatedAt,
		UpdatedAt: company.UpdatedAt,
	}, nil
}

func (s *CompanyService) Delete(id uuid.UUID) error {
	return s.ICompanyRepository.Delete(id)
}

func (s *CompanyService) SelectCompany(sessionID, companyID, userID uuid.UUID) (entity.Company, error) {
	_, err := s.ICompanyRepository.FindMember(context.Background(), companyID, userID)
	if err != nil {
		return entity.Company{}, errors.New("you are not a member of this company")
	}

	company, err := s.ICompanyRepository.FindById(companyID)
	if err != nil {
		return entity.Company{}, errors.New("company not found")
	}

	if err := s.ITokenRepository.SetActiveCompany(sessionID, companyID); err != nil {
		return entity.Company{}, errors.New("failed to set active company")
	}

	return entity.Company{
		Id:        company.Id,
		Name:      company.Name,
		LegalName: company.LegalName,
		TaxID:     company.TaxID,
		Address:   company.Address,
		Phone:     company.Phone,
		Email:     company.Email,
		Industry:  company.Industry,
		Currency:  company.Currency,
		CreatedBy: company.CreatedBy,
		CreatedAt: company.CreatedAt,
		UpdatedAt: company.UpdatedAt,
	}, nil
}

func (s *CompanyService) GetActiveCompany(sessionID uuid.UUID) (entity.Company, error) {
	companyID, err := s.ITokenRepository.GetActiveCompany(sessionID)
	if err != nil {
		return entity.Company{}, errors.New("failed to retrieve session")
	}
	if companyID == nil {
		return entity.Company{}, errors.New("no active company selected")
	}

	company, err := s.ICompanyRepository.FindById(*companyID)
	if err != nil {
		return entity.Company{}, errors.New("active company not found")
	}

	return entity.Company{
		Id:        company.Id,
		Name:      company.Name,
		LegalName: company.LegalName,
		TaxID:     company.TaxID,
		Address:   company.Address,
		Phone:     company.Phone,
		Email:     company.Email,
		Industry:  company.Industry,
		Currency:  company.Currency,
		CreatedBy: company.CreatedBy,
		CreatedAt: company.CreatedAt,
		UpdatedAt: company.UpdatedAt,
	}, nil
}

func (s *CompanyService) Destroy(ids []uuid.UUID) error {
	if len(ids) == 0 {
		return errors.New("no ids provided")
	}
	return s.ICompanyRepository.BulkDestroy(ids)
}
