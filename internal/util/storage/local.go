package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// LocalStorage menyimpan file ke disk lokal server.
type LocalStorage struct {
	BaseDir string // Contoh: "./public/uploads"
	BaseURL string // Contoh: "http://localhost:3000/uploads"
}

// NewLocalStorage membuat instance LocalStorage dan memastikan BaseDir ada.
func NewLocalStorage(baseDir string, baseURL string) *LocalStorage {
	os.MkdirAll(baseDir, os.ModePerm)
	return &LocalStorage{
		BaseDir: baseDir,
		BaseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (s *LocalStorage) Save(_ context.Context, file *multipart.FileHeader, destinationFolder string) (string, error) {
	destDir := filepath.Join(s.BaseDir, destinationFolder)
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	ext := filepath.Ext(file.Filename)
	uniqueName := fmt.Sprintf("%d-%s%s", time.Now().UnixMilli(), uuid.NewString(), ext)
	relativePath := filepath.Join(destinationFolder, uniqueName)
	fullPath := filepath.Join(s.BaseDir, relativePath)

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return relativePath, nil
}

func (s *LocalStorage) Delete(_ context.Context, path string) error {
	fullPath := filepath.Join(s.BaseDir, path)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (s *LocalStorage) GetURL(path string) string {
	return s.BaseURL + "/" + strings.ReplaceAll(path, "\\", "/")
}
