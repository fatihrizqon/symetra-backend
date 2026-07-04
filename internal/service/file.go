package service

import (
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/response"
	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/fatihrizqon/gofiber-microservice/internal/util/storage"
	"github.com/google/uuid"
)

const maxFileSize = 5 * 1024 * 1024 // 5MB

var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"application/pdf": true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
}

type IFileService interface {
	Upload(ctx context.Context, file *multipart.FileHeader, uploadedBy uuid.UUID) (response.FileResponse, error)
	Delete(ctx context.Context, fileId uuid.UUID) error
}

type FileService struct {
	fileRepo repository.IFileRepository
	storage  storage.IStorage
}

func NewFileService(fileRepo repository.IFileRepository, store storage.IStorage) IFileService {
	return &FileService{
		fileRepo: fileRepo,
		storage:  store,
	}
}

func (s *FileService) Upload(ctx context.Context, file *multipart.FileHeader, uploadedBy uuid.UUID) (response.FileResponse, error) {
	if file.Size > maxFileSize {
		return response.FileResponse{}, errors.New("file size exceeds maximum limit of 5MB")
	}

	mimeType, err := detectMimeType(file)
	if err != nil {
		return response.FileResponse{}, errors.New("failed to detect file type")
	}

	if !allowedMimeTypes[mimeType] {
		ext := filepath.Ext(file.Filename)
		return response.FileResponse{}, errors.New("file type not allowed: " + ext)
	}

	relativePath, err := s.storage.Save(ctx, file, "documents")
	if err != nil {
		return response.FileResponse{}, err
	}

	record := entity.File{
		OriginalName: file.Filename,
		Path:         relativePath,
		MimeType:     mimeType,
		Size:         file.Size,
		UploadedBy:   uploadedBy,
	}

	record, err = s.fileRepo.Create(record)
	if err != nil {
		_ = s.storage.Delete(ctx, relativePath)
		return response.FileResponse{}, errors.New("failed to save file metadata")
	}

	return response.FileResponse{
		Id:           record.Id,
		OriginalName: record.OriginalName,
		MimeType:     record.MimeType,
		Size:         record.Size,
		URL:          s.storage.GetURL(record.Path),
		UploadedBy:   record.UploadedBy,
		CreatedAt:    record.CreatedAt,
	}, nil
}

func (s *FileService) Delete(ctx context.Context, fileId uuid.UUID) error {
	record, err := s.fileRepo.FindById(fileId)
	if err != nil {
		return errors.New("file not found")
	}

	if err := s.storage.Delete(ctx, record.Path); err != nil {
		return err
	}

	return s.fileRepo.Delete(fileId)
}

func detectMimeType(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	buf := make([]byte, 512)
	n, err := src.Read(buf)
	if err != nil {
		return "", err
	}

	mime := http.DetectContentType(buf[:n])
	mime = strings.Split(mime, ";")[0]

	// http.DetectContentType cannot detect xlsx, so fallback to extension check
	if mime == "application/zip" || mime == "application/octet-stream" {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext == ".xlsx" {
			return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", nil
		}
	}

	return mime, nil
}
