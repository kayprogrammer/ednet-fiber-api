package profiles

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/modules/base"
	"github.com/kayprogrammer/ednet-fiber-api/modules/courses"
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
		ctx := c.Context()
		data := ProfileUpdateSchema{}
		if errCode, errData := config.ValidateFormRequest(c, &data); errData != nil {
			return config.APIError(c, *errCode, *errData)
		}

		existingUser := userManager.GetByUsername(db, ctx, data.Username)
		if existingUser != nil && existingUser.ID != user.ID {
			return config.APIError(c, 422, config.ValidationErr("username", "Username already used"))
		}
		
		// Check and validate image
		file, err := config.ValidateImage(c, "avatar", false)
		if err != nil {
			return c.Status(422).JSON(err)
		}
		
		// Upload File
		var avatar *string
		if file != nil {
			avatarStr := config.UploadFile(file, string(config.FF_AVATARS))
			avatar = &avatarStr
		}
		updatedUser := profileManager.Update(db, ctx, user, data, avatar)
		response := ProfileResponseSchema{
			ResponseSchema: base.ResponseMessage("Profile fetched"),
			Data:           ProfileSchema{}.Assign(updatedUser),
		}
		return c.Status(200).JSON(response)
	}
}

// @Summary Get Your Enrolled Courses
// @Description `This endpoint allows a user to view his/her enrolled courses`
// @Tags Profiles
// @Param page query int false "Current Page" default(1)
// @Param limit query int false "Page Limit" default(100)
// @Param title query string false "Filter By Title"
// @Param status query string false "Filter By Status (inactive, active, completed, dropped)"
// @Param instructor query string false "Filter By Instructor's Name Or Username"
// @Param isFree query bool false "Filter By Free Status"
// @Success 200 {object} courses.CoursesResponseSchema
// @Failure 401 {object} base.UnauthorizedErrorExample
// @Router /profiles/courses [get]
// @Security BearerAuth
func GetEnrolledCourses(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := base.RequestUser(c)
		coursesData := profileManager.GetAllPaginatedEnrolledCourses(db, c, user, c.Query("status", ""))
		response := courses.CoursesResponseSchema{
			ResponseSchema: base.ResponseMessage("Courses Fetched Successfully"),
		}.Assign(coursesData)
		return c.Status(200).JSON(response)
	}
}