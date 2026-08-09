package service

import (
	"errors"
	"math"

	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/request"
	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/response"
	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/fatihrizqon/gofiber-microservice/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type IJournalEntryService interface {
	Create(companyID, userId uuid.UUID, req request.JournalEntryCreateRequest) (entity.JournalEntry, error)
	FindAll(companyID uuid.UUID, qp *util.QueryParams) ([]entity.JournalEntry, int, error)
	FindById(companyID, id uuid.UUID) (entity.JournalEntry, error)
	Update(companyID uuid.UUID, req request.JournalEntryUpdateRequest) (entity.JournalEntry, error)
	Delete(companyID, id uuid.UUID) error
	Post(companyID, id uuid.UUID) error
	Void(companyID, id uuid.UUID) error
	Destroy(companyID uuid.UUID, ids []uuid.UUID) error
}

type JournalEntryService struct {
	validate         *validator.Validate
	journalEntryRepo repository.IJournalEntryRepository
	coaRepo          repository.ICOARepository
	fiscalYearRepo   repository.IFiscalYearRepository
}

func NewJournalEntryService(
	validate *validator.Validate,
	journalEntryRepo repository.IJournalEntryRepository,
	coaRepo repository.ICOARepository,
	fiscalYearRepo repository.IFiscalYearRepository,
) IJournalEntryService {
	return &JournalEntryService{
		validate:         validate,
		journalEntryRepo: journalEntryRepo,
		coaRepo:          coaRepo,
		fiscalYearRepo:   fiscalYearRepo,
	}
}

func (s *JournalEntryService) Create(
	companyID, userId uuid.UUID,
	req request.JournalEntryCreateRequest,
) (entity.JournalEntry, error) {

	if len(req.Lines) < 2 {
		return entity.JournalEntry{}, errors.New(
			"journal entry must have at least 2 lines",
		)
	}

	var totalDebit float64
	var totalCredit float64

	lines := make([]entity.JournalLine, 0, len(req.Lines))

	for _, lineReq := range req.Lines {
		if (lineReq.Debit > 0 && lineReq.Credit > 0) ||
			(lineReq.Debit == 0 && lineReq.Credit == 0) {
			return entity.JournalEntry{}, errors.New(
				"each line must have either debit or credit, not both or neither",
			)
		}

		coa, err := s.coaRepo.FindById(companyID, lineReq.CoaId)
		if err != nil {
			return entity.JournalEntry{}, errors.New("invalid coa_id")
		}

		if !coa.Active {
			return entity.JournalEntry{}, errors.New("coa is inactive")
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
		return entity.JournalEntry{}, errors.New(
			"journal entry must be balanced (total debit = total credit)",
		)
	}

	journalNum, err := s.journalEntryRepo.GenerateJournalNumber(
		companyID,
		req.Date,
		string(entity.JournalTypeGeneral),
	)
	if err != nil {
		return entity.JournalEntry{}, err
	}

	files := make([]entity.File, 0, len(req.FileIds))

	for _, fileId := range req.FileIds {
		files = append(files, entity.File{
			Id: fileId,
		})
	}

	journal := entity.JournalEntry{
		CompanyId:     companyID,
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

	if err := s.journalEntryRepo.Create(&journal); err != nil {
		return entity.JournalEntry{}, err
	}

	return response.NewJournalEntryResponse(journal), nil
}

func (s *JournalEntryService) FindAll(
	companyID uuid.UUID,
	qp *util.QueryParams,
) ([]entity.JournalEntry, int, error) {

	journals, total, err := s.journalEntryRepo.FindAll(companyID, qp)
	if err != nil {
		return nil, 0, err
	}

	return response.NewJournalEntryResponses(journals), int(total), nil
}

func (s *JournalEntryService) FindById(
	companyID, id uuid.UUID,
) (entity.JournalEntry, error) {

	journal, err := s.journalEntryRepo.FindById(companyID, id)
	if err != nil {
		return entity.JournalEntry{}, errors.New("journal entry not found")
	}

	return response.NewJournalEntryResponse(journal), nil
}

func (s *JournalEntryService) Update(
	companyID uuid.UUID,
	req request.JournalEntryUpdateRequest,
) (entity.JournalEntry, error) {

	journal, err := s.journalEntryRepo.FindById(companyID, req.Id)
	if err != nil {
		return entity.JournalEntry{}, errors.New("journal entry not found")
	}

	if journal.Status != entity.JournalStatusDraft {
		return entity.JournalEntry{}, errors.New(
			"only draft journal entries can be edited",
		)
	}

	var totalDebit float64
	var totalCredit float64

	lines := make([]entity.JournalLine, 0, len(req.Lines))

	for _, lineReq := range req.Lines {
		if (lineReq.Debit > 0 && lineReq.Credit > 0) ||
			(lineReq.Debit == 0 && lineReq.Credit == 0) {
			return entity.JournalEntry{}, errors.New(
				"each line must have either debit or credit",
			)
		}

		coa, err := s.coaRepo.FindById(companyID, lineReq.CoaId)
		if err != nil || !coa.Active {
			return entity.JournalEntry{}, errors.New(
				"invalid or inactive coa_id",
			)
		}

		totalDebit += lineReq.Debit
		totalCredit += lineReq.Credit

		lines = append(lines, entity.JournalLine{
			JournalEntryId: req.Id,
			CoaId:          lineReq.CoaId,
			Description:    lineReq.Description,
			Debit:          lineReq.Debit,
			Credit:         lineReq.Credit,
		})
	}

	if math.Abs(totalDebit-totalCredit) > 0.0001 {
		return entity.JournalEntry{}, errors.New(
			"journal entry must be balanced",
		)
	}

	files := make([]entity.File, 0, len(req.FileIds))

	for _, fileId := range req.FileIds {
		files = append(files, entity.File{
			Id: fileId,
		})
	}

	journal.Date = req.Date
	journal.Description = req.Description
	journal.TotalDebit = totalDebit
	journal.TotalCredit = totalCredit
	journal.Lines = lines
	journal.Files = files

	if err := s.journalEntryRepo.Update(&journal); err != nil {
		return entity.JournalEntry{}, err
	}

	updatedJournal, err := s.journalEntryRepo.FindById(
		companyID,
		journal.Id,
	)
	if err != nil {
		return entity.JournalEntry{}, err
	}

	return response.NewJournalEntryResponse(updatedJournal), nil
}

func (s *JournalEntryService) Delete(
	companyID, id uuid.UUID,
) error {

	journal, err := s.journalEntryRepo.FindById(companyID, id)
	if err != nil {
		return errors.New("journal entry not found")
	}

	if journal.Status != entity.JournalStatusDraft {
		return errors.New(
			"only draft journal entries can be deleted",
		)
	}

	return s.journalEntryRepo.Delete(companyID, id)
}

func (s *JournalEntryService) Post(
	companyID, id uuid.UUID,
) error {

	journal, err := s.journalEntryRepo.FindById(companyID, id)
	if err != nil {
		return errors.New("journal entry not found")
	}

	if journal.Status != entity.JournalStatusDraft {
		return errors.New(
			"only draft journal entries can be posted",
		)
	}

	if math.Abs(journal.TotalDebit-journal.TotalCredit) > 0.0001 {
		return errors.New(
			"journal entry must be balanced to be posted",
		)
	}

	period, err := s.fiscalYearRepo.GetOpenPeriodByDate(
		companyID,
		journal.Date.Format("2006-01-02"),
	)
	if err != nil {
		return errors.New(
			"date does not fall within an open fiscal period",
		)
	}

	return s.journalEntryRepo.Post(
		companyID,
		id,
		period.Id,
	)
}

func (s *JournalEntryService) Void(
	companyID, id uuid.UUID,
) error {

	journal, err := s.journalEntryRepo.FindById(companyID, id)
	if err != nil {
		return errors.New("journal entry not found")
	}

	if journal.Status != entity.JournalStatusPosted {
		return errors.New(
			"only posted journal entries can be voided",
		)
	}

	return s.journalEntryRepo.Void(companyID, id)
}

func (s *JournalEntryService) Destroy(
	companyID uuid.UUID,
	ids []uuid.UUID,
) error {

	if len(ids) == 0 {
		return errors.New("no ids provided")
	}

	return s.journalEntryRepo.BulkDestroy(companyID, ids)
}
