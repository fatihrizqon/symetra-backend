package service

import (
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/request"
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/response"
	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ICompanyConfigurationService interface {
	GetByCompanyId(companyID uuid.UUID) (response.CompanyConfigurationResponse, error)
	Upsert(companyID uuid.UUID, req request.CompanyConfigurationUpdateRequest) (response.CompanyConfigurationResponse, error)
}

type CompanyConfigurationService struct {
	Repo     repository.ICompanyConfigurationRepository
	validate *validator.Validate
}

func NewCompanyConfigurationService(repo repository.ICompanyConfigurationRepository, validate *validator.Validate) ICompanyConfigurationService {
	return &CompanyConfigurationService{Repo: repo, validate: validate}
}

func (s *CompanyConfigurationService) GetByCompanyId(companyID uuid.UUID) (response.CompanyConfigurationResponse, error) {
	config, err := s.Repo.GetByCompanyId(companyID)
	if err != nil {
		return response.CompanyConfigurationResponse{}, err
	}
	return mapCompanyConfiguration(config), nil
}

func (s *CompanyConfigurationService) Upsert(companyID uuid.UUID, req request.CompanyConfigurationUpdateRequest) (response.CompanyConfigurationResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.CompanyConfigurationResponse{}, err
	}
	config := entity.CompanyConfiguration{
		CompanyId:              companyID,
		EnableTax:              req.EnableTax,
		TaxRate:                req.TaxRate,
		ARAccountId:            req.ARAccountId,
		APAccountId:            req.APAccountId,
		SalesRevenueAccountId:  req.SalesRevenueAccountId,
		TaxPayableAccountId:    req.TaxPayableAccountId,
		TaxReceivableAccountId: req.TaxReceivableAccountId,
		BankAccountId:          req.BankAccountId,
		CashAccountId:          req.CashAccountId,
		RetainedEarningsCOAId:  req.RetainedEarningsCOAId,
		InvoicePrefix:          req.InvoicePrefix,
		QuotationPrefix:        req.QuotationPrefix,
		InvoiceDueDays:         req.InvoiceDueDays,
	}
	// Fallbacks if prefixes are empty
	if config.InvoicePrefix == "" {
		config.InvoicePrefix = "INV"
	}
	if config.QuotationPrefix == "" {
		config.QuotationPrefix = "QUO"
	}
	if config.InvoiceDueDays <= 0 {
		config.InvoiceDueDays = 30
	}

	result, err := s.Repo.Upsert(companyID, config)
	if err != nil {
		return response.CompanyConfigurationResponse{}, err
	}
	return mapCompanyConfiguration(result), nil
}

func mapCompanyConfiguration(v entity.CompanyConfiguration) response.CompanyConfigurationResponse {
	resp := response.CompanyConfigurationResponse{
		Id:                     v.Id,
		CompanyId:              v.CompanyId,
		EnableTax:              v.EnableTax,
		TaxRate:                v.TaxRate,
		ARAccountId:            v.ARAccountId,
		APAccountId:            v.APAccountId,
		SalesRevenueAccountId:  v.SalesRevenueAccountId,
		TaxPayableAccountId:    v.TaxPayableAccountId,
		TaxReceivableAccountId: v.TaxReceivableAccountId,
		BankAccountId:          v.BankAccountId,
		CashAccountId:          v.CashAccountId,
		RetainedEarningsCOAId:  v.RetainedEarningsCOAId,
		InvoicePrefix:          v.InvoicePrefix,
		QuotationPrefix:        v.QuotationPrefix,
		InvoiceDueDays:         v.InvoiceDueDays,
	}

	if v.ARAccount != nil {
		resp.ARAccount = &response.COAResponse{Id: v.ARAccount.Id, Code: v.ARAccount.Code, Name: v.ARAccount.Name}
	}
	if v.APAccount != nil {
		resp.APAccount = &response.COAResponse{Id: v.APAccount.Id, Code: v.APAccount.Code, Name: v.APAccount.Name}
	}
	return resp
}
