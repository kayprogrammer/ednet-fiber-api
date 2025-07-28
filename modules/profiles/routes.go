package profiles

import (
	"fmt"
	"log"

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

// @Summary Create/Update Lesson Progress
// @Description `This endpoint allows a user to create or update a lesson progress`
// @Tags Profiles
// @Param slug path string true "Lesson Slug"
// @Param lesson_progress body LessonProgressInputSchema true "Lesson Progress object"
// @Success 201 {object} LessonProgressResponseSchema
// @Failure 422 {object} base.ValidationErrorExample
// @Failure 404 {object} base.NotFoundErrorExample
// @Failure 401 {object} base.UnauthorizedErrorExample
// @Router /profiles/lessons/{slug}/progress [post]
// @Security BearerAuth
func CreateOrUpdateLessonProgress(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := base.RequestUser(c)
		ctx := c.Context()
		data := LessonProgressInputSchema{}
		if errCode, errData := config.ValidateRequest(c, &data); errData != nil {
			return config.APIError(c, *errCode, *errData)
		}

		lesson := courseManager.GetCourseLessonBySlug(db, ctx, c.Params("slug"), nil, true)
		if lesson == nil {
			return config.APIError(c, 404, config.NotFoundErr("Lesson not found"))
		}

		enrollment := courseManager.GetExistentEnrollmentByUserAndCourse(db, ctx, user, lesson.Edges.Course, false)
		if enrollment == nil {
			return config.APIError(c, 403, config.RequestErr(config.ERR_NOT_ALLOWED, "You are not enrolled in this lesson"))
		}

		lessonProgress, message := profileManager.CreateOrUpdateLessonProgress(db, ctx, user, lesson, data.IsCompleted)
		response := LessonProgressResponseSchema{
			ResponseSchema: base.ResponseMessage(fmt.Sprintf("Lesson progress %s successfully", message)),
			Data:           LessonProgressResponseData{}.Assign(lessonProgress),
		}
		return c.Status(201).JSON(response)
	}
}

// @Summary Get Lesson Progress
// @Description `This endpoint allows a user to get his/her lesson progress`
// @Tags Profiles
// @Param slug path string true "Lesson Slug"
// @Success 200 {object} LessonProgressResponseSchema
// @Failure 404 {object} base.NotFoundErrorExample
// @Failure 401 {object} base.UnauthorizedErrorExample
// @Router /profiles/lessons/{slug}/progress [get]
// @Security BearerAuth
func GetLessonProgress(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := base.RequestUser(c)
		ctx := c.Context()
		lesson := courseManager.GetCourseLessonBySlug(db, ctx, c.Params("slug"), nil, true)
		if lesson == nil {
			return config.APIError(c, 404, config.NotFoundErr("Lesson not found"))
		}

		enrollment := courseManager.GetExistentEnrollmentByUserAndCourse(db, ctx, user, lesson.Edges.Course, false)
		if enrollment == nil {
			return config.APIError(c, 403, config.RequestErr(config.ERR_NOT_ALLOWED, "You are not enrolled in this lesson"))
		}

		lessonProgress := profileManager.GetLessonProgress(db, ctx, user, lesson.ID)
		if lessonProgress == nil {
			return config.APIError(c, 404, config.NotFoundErr("No progress recorded for this lesson"))
		}
		response := LessonProgressResponseSchema{
			ResponseSchema: base.ResponseMessage("Lesson progress fetched successfully"),
			Data:           LessonProgressResponseData{}.Assign(lessonProgress),
		}
		return c.Status(200).JSON(response)
	}
}

// @Summary Get Course Progress
// @Description `This endpoint allows a user to get his/her course progress`
// @Tags Profiles
// @Param slug path string true "Course Slug"
// @Success 200 {object} CourseProgressResponseSchema
// @Failure 404 {object} base.NotFoundErrorExample
// @Failure 401 {object} base.UnauthorizedErrorExample
// @Router /profiles/courses/{slug}/progress [get]
// @Security BearerAuth
func GetCourseProgress(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := base.RequestUser(c)
		ctx := c.Context()
		course := courseManager.GetCourseBySlug(db, ctx, c.Params("slug"), nil, false)
		if course == nil {
			return config.APIError(c, 404, config.NotFoundErr("Course not found"))
		}

		enrollment := courseManager.GetExistentEnrollmentByUserAndCourse(db, ctx, user, course, false)
		if enrollment == nil {
			return config.APIError(c, 403, config.RequestErr(config.ERR_NOT_ALLOWED, "You are not enrolled in this course"))
		}

		courseProgressPercentage := profileManager.GetCourseProgress(db, ctx, user, course)
		response := CourseProgressResponseSchema{
			ResponseSchema: base.ResponseMessage("Course progress fetched successfully"),
			Data:           CourseProgressResponseData{Percentage: courseProgressPercentage},
		}
		return c.Status(200).JSON(response)
	}
}

// @Summary Get Leaderboard
// @Description `This endpoint retrieves the top 100 students by quiz score.`
// @Tags Profiles
// @Failure 200 {object} LeaderboardResponseSchema
// @Router /profiles/leaderboard [get]
// @Security BearerAuth
func GetLeaderboard(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		leaderboard := profileManager.GetLeaderboard(db, c.Context())
		log.Println(leaderboard)
		return c.Status(200).JSON(LeaderboardResponseSchema{
			ResponseSchema: base.ResponseMessage("Leaderboard fetched successfully"),
			Data:           leaderboard,
		})
	}
}