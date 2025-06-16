package config

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
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

func ValidateFile(c *fiber.Ctx, name string, required bool, isVideo bool) (*multipart.FileHeader, *ErrorResponse) {
	file, err := c.FormFile(name)
	errData := ValidationErr(name, "Invalid file type")

	if required && err != nil {
		errData = ValidationErr(name, "File is required")
		return nil, &errData
	}

	// Open the file
	if file != nil {
		fileHandle, err := file.Open()
		if err != nil {
			return nil, &errData
		}

		defer fileHandle.Close()

		// Read the first 512 bytes for content type detection
		buffer := make([]byte, 512)
		_, err = fileHandle.Read(buffer)
		if err != nil {
			return nil, &errData
		}

		// Detect the content type
		contentType := http.DetectContentType(buffer)
		if isVideo && contentType == "video/mp4" {
			return file, nil
		} else {
			switch contentType {
			case "image/jpeg", "image/png", "image/gif":
				return file, nil
			}
		}
		return nil, &errData
	}
	return nil, nil
}

type FILE_FOLDER_CHOICES string

const (
	FF_AVATARS = "avatars"
	FF_THUMBNAIL = "thumbnails"
	FF_INTRO_VIDEOS = "intro_videos"
	FF_LESSON_VIDEOS = "lesson_videos"
)