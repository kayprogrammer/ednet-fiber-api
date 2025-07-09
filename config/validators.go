package config

import (
	"mime/multipart"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/ednet-fiber-api/ent/course"
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

func DifficultyTypeValidator(fl validator.FieldLevel) bool {
	fieldVal := fl.Field().Interface().(course.Difficulty)
	return fieldVal == course.DifficultyBeginner || fieldVal == course.DifficultyAdvanced || fieldVal == course.DifficultyIntermediate
}

func EnrollmentTypeValidator(fl validator.FieldLevel) bool {
	fieldVal := fl.Field().Interface().(course.EnrollmentType)
	return fieldVal == course.EnrollmentTypeOpen || fieldVal == course.EnrollmentTypeInviteOnly || fieldVal == course.EnrollmentTypeRestricted
}
