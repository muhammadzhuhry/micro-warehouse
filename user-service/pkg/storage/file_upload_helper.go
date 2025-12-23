package storage

import (
	"context"
	"fmt"
	"micro-warehouse/user-service/configs"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2/log"
)

const (
	MaxImageSize = 2 * 1024 * 1024 // 2 MB

	AllowedImageExtensions = ".jpg,.jpeg,.png,.webp,.svg"

	BucketUser = "users"
)

type FileUploadHelper struct {
	storage SupabaseInterface
	config  configs.Config
}

func NewFileUploadHelper(storage SupabaseInterface, config configs.Config) *FileUploadHelper {
	return &FileUploadHelper{
		storage: storage,
		config:  config,
	}
}

func (f *FileUploadHelper) UploadPhoto(ctx context.Context, file *multipart.FileHeader) (*UploadResult, error) {
	if err := f.validateImageFile(file, MaxImageSize); err != nil {
		log.Errorf("failed to validate image: %v", err)
		return nil, err
	}

	result, err := f.storage.UploadFile(ctx, file, BucketUser)
	if err != nil {
		log.Errorf("failed to upload file: %v", err)
		return nil, err
	}

	return result, nil
}

func (f *FileUploadHelper) validateImageFile(file *multipart.FileHeader, maxSize int64) error {
	if !validateFileSize(file.Size, maxSize) {
		return fmt.Errorf("file size exceeds the maximum allowed size")
	}

	if !validateFileExtension(file.Filename, AllowedImageExtensions) {
		return fmt.Errorf("file extension is not allowed")
	}

	return nil
}

func validateFileSize(size, maxSize int64) bool {
	return size <= maxSize
}

func getFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

func validateFileExtension(ext string, allowedExts string) bool {
	allowed := strings.Split(allowedExts, ",")
	for _, allow := range allowed {
		if strings.TrimSpace(allow) == ext {
			return true
		}
	}
	return false
}
