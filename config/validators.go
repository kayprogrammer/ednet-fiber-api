package config

import (
	"mime/multipart"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Validates if a account type value is the correct one
// func AccountTypeValidator(fl validator.FieldLevel) bool {
// 	return fl.Field().Interface().(choices.AccType).IsValid()
// }

func ValidateImage(c *fiber.Ctx, name string, required bool) (*multipart.FileHeader, *ErrorResponse) {
	file, err := c.FormFile(name)
	errData := ValidationErr(name, "Invalid image type")

	if required && err != nil {
		errData = ValidationErr(name, "Image is required")
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
		switch contentType {
		case "image/jpeg", "image/png", "image/gif":
			return file, nil
		}
		return nil, &errData
	}
	return nil, nil
}