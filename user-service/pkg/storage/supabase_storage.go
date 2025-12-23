package storage

import (
	"context"
	"fmt"
	"micro-warehouse/user-service/configs"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	storage_go "github.com/supabase-community/storage-go"
)

type SupabaseInterface interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (*UploadResult, error)
}

type SupabaseStorage struct {
	config configs.Config
	client *storage_go.Client
}

type UploadResult struct {
	URL      string `json:"url"`
	Path     string `json:"path"`
	Filename string `json:"filename"`
}

func NewSupabaseStorage(config configs.Config) SupabaseInterface {
	client := storage_go.NewClient(config.Supabase.Url, config.Supabase.Key, nil)

	return &SupabaseStorage{
		config: config,
		client: client,
	}
}

// UploadFile implements [SupabaseInterface].
func (s *SupabaseStorage) UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (*UploadResult, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%s_%d%s", strings.TrimSuffix(file.Filename, ext), timestamp, ext)

	// Create file path
	filePath := fmt.Sprintf("%s/%s", folder, filename)

	// Content type handler
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		// switch content type based on file extension
		switch strings.ToLower(ext) {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".webp":
			contentType = "image/webp"
		case ".svg":
			contentType = "image/svg+xml"
		default:
			contentType = "application/octet-stream"
		}
	}

	// Create client with proper Content-Type
	client := storage_go.NewClient(s.config.Supabase.Url, s.config.Supabase.Key, map[string]string{
		"Content-Type": contentType,
	})

	// Upload file
	_, err = client.UploadFile(s.config.Supabase.Bucket, filePath, src)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Generate public URL
	publicURL := client.GetPublicUrl(s.config.Supabase.Bucket, filePath)

	return &UploadResult{
		URL:      publicURL.SignedURL,
		Path:     filePath,
		Filename: filename,
	}, nil
}
