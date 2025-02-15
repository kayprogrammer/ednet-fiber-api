package profiles

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/modules/base"
)

var profileManager = ProfileManager{}

// @Summary Get Your Profile
// @Description `This endpoint allows a user to view his/her profile`
// @Tags Profiles
// @Success 200 {object} ProfileResponseSchema
// @Failure 401 {object} base.UnauthorizedErrorExample
// @Router /profiles [get]
// @Security BearerAuth
func GetProfile(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := base.RequestUser(c)
		response := ProfileResponseSchema{
			ResponseSchema: base.ResponseMessage("Profile fetched"),
			Data:           ProfileSchema{}.Assign(user),
		}
		return c.Status(200).JSON(response)
	}
}

// @Summary Update Your Profile
// @Description `This endpoint allows a user to update his/her profile`
// @Tags Profiles
// @Param profile formData ProfileUpdateSchema true "Profile object"
// @Param avatar formData file false "Profile picture to upload"
// @Success 200 {object} ProfileResponseSchema
// @Failure 422 {object} base.ValidationErrorExample
// @Failure 401 {object} base.UnauthorizedErrorExample
// @Router /profiles [put]
// @Security BearerAuth
func UpdateProfile(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := base.RequestUser(c)
		data := ProfileUpdateSchema{}
		if errCode, errData := config.ValidateRequest(c, &data); errData != nil {
			return config.APIError(c, *errCode, *errData)
		}

		// Check and validate image
		file, err := config.ValidateImage(c, "cover_image", false)
		if err != nil {
			return c.Status(422).JSON(err)
		}
		
		// Upload File
		var avatar *string
		if file != nil {
			avatarStr := config.UploadFile(file, string(config.FF_AVATARS))
			avatar = &avatarStr
		}
		updatedUser := profileManager.Update(db, c.Context(), user, data, avatar)
		response := ProfileResponseSchema{
			ResponseSchema: base.ResponseMessage("Profile fetched"),
			Data:           ProfileSchema{}.Assign(updatedUser),
		}
		return c.Status(200).JSON(response)
	}
}
