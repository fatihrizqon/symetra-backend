package service

import (
	"errors"
	"math"

	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/request"
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/response"
	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/google/uuid"
)

type IJournalEntryService interface {
	Create(companyId, userId uuid.UUID, req request.JournalEntryCreateRequest) (response.JournalEntryResponse, error)
	FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]response.JournalEntryResponse, int, error)
	FindById(companyId, id uuid.UUID) (response.JournalEntryResponse, error)
	Update(companyId, id uuid.UUID, req request.JournalEntryUpdateRequest) (response.JournalEntryResponse, error)
	Delete(companyId, id uuid.UUID) error
	Post(companyId, id uuid.UUID) error
	Void(companyId, id uuid.UUID) error
}

type JournalEntryService struct {
	repo         repository.IJournalEntryRepository
	coaRepo      repository.ICOARepository
	fiscalRepo   repository.IFiscalYearRepository
}

func NewJournalEntryService(
	repo repository.IJournalEntryRepository,
	coaRepo repository.ICOARepository,
	fiscalRepo repository.IFiscalYearRepository,
) IJournalEntryService {
	return &JournalEntryService{repo, coaRepo, fiscalRepo}
}

func (s *JournalEntryService) Create(companyId, userId uuid.UUID, req request.JournalEntryCreateRequest) (response.JournalEntryResponse, error) {
	if len(req.Lines) < 2 {
		return response.JournalEntryResponse{}, errors.New("journal entry must have at least 2 lines")
	}

	var totalDebit, totalCredit float64
	var lines []entity.JournalLine

	for _, lineReq := range req.Lines {
		if (lineReq.Debit > 0 && lineReq.Credit > 0) || (lineReq.Debit == 0 && lineReq.Credit == 0) {
			return response.JournalEntryResponse{}, errors.New("each line must have either debit or credit, not both or neither")
		}

		coa, err := s.coaRepo.FindById(companyId, lineReq.CoaId)
		if err != nil {
			return response.JournalEntryResponse{}, errors.New("invalid coa_id")
		}
		if !coa.Active {
			return response.JournalEntryResponse{}, errors.New("coa is inactive")
		}

		totalDebit += lineReq.Debit
		totalCredit += lineReq.Credit

		lines = append(lines, entity.JournalLine{
			CoaId:       lineReq.CoaId,
			Description: lineReq.Description,
			Debit:       lineReq.Debit,
			Credit:      lineReq.Credit,
		})
	}

	if math.Abs(totalDebit-totalCredit) > 0.0001 {
		return response.JournalEntryResponse{}, errors.New("journal entry must be balanced (total debit = total credit)")
	}

	journalNum, err := s.repo.GenerateJournalNumber(companyId, req.Date, string(entity.JournalTypeGeneral))
	if err != nil {
		return response.JournalEntryResponse{}, err
	}

	var files []entity.File
	for _, fileId := range req.FileIds {
		files = append(files, entity.File{Id: fileId})
	}

	journal := entity.JournalEntry{
		CompanyId:     companyId,
		JournalNumber: journalNum,
		Type:          entity.JournalTypeGeneral,
		Date:          req.Date,
		Description:   req.Description,
		Status:        entity.JournalStatusDraft,
		TotalDebit:    totalDebit,
		TotalCredit:   totalCredit,
		CreatedBy:     userId,
		Lines:         lines,
		Files:         files,
	}

	if err := s.repo.Create(&journal); err != nil {
		return response.JournalEntryResponse{}, err
	}

	// Reload to get preloaded COAs
	createdJournal, _ := s.repo.FindById(companyId, journal.Id)
	return response.FromJournalEntryEntity(createdJournal), nil
}

func (s *JournalEntryService) FindAll(companyId uuid.UUID, qp *util.QueryParams) ([]response.JournalEntryResponse, int, error) {
	journals, total, err := s.repo.FindAll(companyId, qp)
	if err != nil {
		return nil, 0, err
	}
	return response.FromJournalEntryEntities(journals), int(total), nil
}

func (s *JournalEntryService) FindById(companyId, id uuid.UUID) (response.JournalEntryResponse, error) {
	journal, err := s.repo.FindById(companyId, id)
	if err != nil {
		return response.JournalEntryResponse{}, errors.New("journal entry not found")
	}
	return response.FromJournalEntryEntity(journal), nil
}

func (s *JournalEntryService) Update(companyId, id uuid.UUID, req request.JournalEntryUpdateRequest) (response.JournalEntryResponse, error) {
	journal, err := s.repo.FindById(companyId, id)
	if err != nil {
		return response.JournalEntryResponse{}, errors.New("journal entry not found")
	}

	if journal.Status != entity.JournalStatusDraft {
		return response.JournalEntryResponse{}, errors.New("only draft journal entries can be edited")
	}

	var totalDebit, totalCredit float64
	var lines []entity.JournalLine

	for _, lineReq := range req.Lines {
		if (lineReq.Debit > 0 && lineReq.Credit > 0) || (lineReq.Debit == 0 && lineReq.Credit == 0) {
			return response.JournalEntryResponse{}, errors.New("each line must have either debit or credit")
		}

		coa, err := s.coaRepo.FindById(companyId, lineReq.CoaId)
		if err != nil || !coa.Active {
			return response.JournalEntryResponse{}, errors.New("invalid or inactive coa_id")
		}

		totalDebit += lineReq.Debit
		totalCredit += lineReq.Credit

		lines = append(lines, entity.JournalLine{
			JournalEntryId: id,
			CoaId:          lineReq.CoaId,
			Description:    lineReq.Description,
			Debit:          lineReq.Debit,
			Credit:         lineReq.Credit,
		})
	}

	if math.Abs(totalDebit-totalCredit) > 0.0001 {
		return response.JournalEntryResponse{}, errors.New("journal entry must be balanced")
	}

	var files []entity.File
	for _, fileId := range req.FileIds {
		files = append(files, entity.File{Id: fileId})
	}

	journal.Date = req.Date
	journal.Description = req.Description
	journal.TotalDebit = totalDebit
	journal.TotalCredit = totalCredit
	journal.Lines = lines
	journal.Files = files

	if err := s.repo.Update(&journal); err != nil {
		return response.JournalEntryResponse{}, err
	}

	updatedJournal, _ := s.repo.FindById(companyId, journal.Id)
	return response.FromJournalEntryEntity(updatedJournal), nil
}

func (s *JournalEntryService) Delete(companyId, id uuid.UUID) error {
	journal, err := s.repo.FindById(companyId, id)
	if err != nil {
		return errors.New("journal entry not found")
	}
	if journal.Status != entity.JournalStatusDraft {
		return errors.New("only draft journal entries can be deleted")
	}
	return s.repo.Delete(companyId, id)
}

func (s *JournalEntryService) Post(companyId, id uuid.UUID) error {
	journal, err := s.repo.FindById(companyId, id)
	if err != nil {
		return errors.New("journal entry not found")
	}
	if journal.Status != entity.JournalStatusDraft {
		return errors.New("only draft journal entries can be posted")
	}
	if math.Abs(journal.TotalDebit-journal.TotalCredit) > 0.0001 {
		return errors.New("journal entry must be balanced to be posted")
	}

	period, err := s.fiscalRepo.GetOpenPeriodByDate(companyId, journal.Date.Format("2006-01-02"))
	if err != nil {
		return errors.New("date does not fall within an open fiscal period")
	}

	return s.repo.Post(companyId, id, period.Id)
}

func (s *JournalEntryService) Void(companyId, id uuid.UUID) error {
	journal, err := s.repo.FindById(companyId, id)
	if err != nil {
		return errors.New("journal entry not found")
	}
	if journal.Status != entity.JournalStatusPosted {
		return errors.New("only posted journal entries can be voided")
	}
	return s.repo.Void(companyId, id)
}
