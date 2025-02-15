package config

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cld *cloudinary.Cloudinary
var err error

func initializeCloudinary () Config {
	cfg := GetConfig()
	// Initialize Cloudinary client
	cld, err = cloudinary.NewFromParams(cfg.CloudinaryCloudName, cfg.CloudinaryApiKey, cfg.CloudinaryApiSecret)
	if err != nil {
		fmt.Println("failed to initialize Cloudinary client: %w", err)
	}
	return cfg
}

func UploadFile(file *multipart.FileHeader, folder string) string {
	cfg := initializeCloudinary()
	if cfg.Environment == "test" {
		return "https://testfile.com"
	}

	folder = fmt.Sprintf("%s/%s", cfg.Environment, folder)

	// Open the file
	src, err := file.Open()
	if err != nil {
		fmt.Println("failed to open file: %w", err)
		return ""
	}
	defer src.Close()

	// Upload the file to Cloudinary
	uploadResult, err := cld.Upload.Upload(context.Background(), src, uploader.UploadParams{Folder: folder})
	if err != nil {
		fmt.Println("failed to upload to Cloudinary: %w", err)
		return ""
	}

	// Return the secure URL of the uploaded file
	return uploadResult.SecureURL
}

type FILE_FOLDER_CHOICES string

const (
	FF_AVATARS = "avatars"
	FF_COURSES = "courses"
)