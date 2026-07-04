package storage

import (
	"context"
	"mime/multipart"
)

// IStorage mendefinisikan standar cara aplikasi berinteraksi dengan media penyimpanan.
// Implementasikan interface ini untuk setiap jenis storage (Local, S3, GCS, dll).
type IStorage interface {
	// Save menyimpan file ke direktori tujuan dan mengembalikan relative path-nya.
	Save(ctx context.Context, file *multipart.FileHeader, destinationFolder string) (string, error)
	// Delete menghapus file berdasarkan relative path.
	Delete(ctx context.Context, path string) error
	// GetURL mengubah relative path menjadi URL publik yang bisa diakses via browser.
	GetURL(path string) string
}
