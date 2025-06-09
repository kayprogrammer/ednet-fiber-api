package instructors

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/modules/base"
	"github.com/kayprogrammer/ednet-fiber-api/modules/courses"
)

var instructorManager = InstructorManager{}
var courseManager = courses.CourseManager{}

// @Summary Retrieve Instructor Courses
// @Description `This endpoint retrieves paginated responses of the authenticated instructor courses`
// @Tags Instructor
// @Param page query int false "Current Page" default(1)
// @Param limit query int false "Page Limit" default(100)
// @Param title query string false "Filter By Title"
// @Param isFree query bool false "Filter By Free Status"
// @Param isPublished query bool false "Filter By Published Status"
// @Param sortByRating query string false "Sort By Rating (asc or desc)"
// @Success 200 {object} courses.CoursesResponseSchema
// @Router /instructor/courses [get]
// @Security BearerAuth
func GetInstructorCourses(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := base.RequestUser(c)
		courseList := instructorManager.GetCoursesPaginated(db, c, user)
		response := courses.CoursesResponseSchema{
			ResponseSchema: base.ResponseMessage("Courses Fetched Successfully"),
		}.Assign(courseList)
		return c.Status(200).JSON(response)
	}
}

// @Summary Create A Course
// @Description `This endpoint allows an instructor to create a course`
// @Tags Instructor
// @Param course formData CourseCreateSchema true "Course object"
// @Param thumbnail formData file true "Thumbnail to upload"
// @Param intro_video formData file true "Intro video to upload"
// @Success 200 {object} courses.CourseResponseSchema
// @Router /instructor/courses [post]
// @Security BearerAuth
func CreateCourse(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		user := base.RequestUser(c)
		data := CourseCreateSchema{}
		if errCode, errData := config.ValidateFormRequest(c, &data); errData != nil {
			return config.APIError(c, *errCode, *errData)
		}
		category := courseManager.GetCategoryBySlug(db, ctx, data.CategorySlug)
		if category == nil {
			return config.APIError(c, 422, config.ValidationErr("categorySlug", "Invalid category slug"))
		}

		// Check and validate files
		thumbnail, err := config.ValidateFile(c, "thumbnail", true)
		if err != nil {
			return c.Status(422).JSON(err)
		}
		introVideo, err := config.ValidateFile(c, "intro_video", true)
		if err != nil {
			return c.Status(422).JSON(err)
		}
		thumbnailUrl := config.UploadFile(thumbnail, string(config.FF_THUMBNAIL))
		introVideoUrl := config.UploadFile(introVideo, string(config.FF_INTRO_VIDEOS))

		course := instructorManager.CreateCourse(db, ctx, user, category, thumbnailUrl, &introVideoUrl, data)
		response := courses.CourseResponseSchema{
			ResponseSchema: base.ResponseMessage("Course Created Successfully"),
			Data:           courses.CourseDetailSchema{}.Assign(course),
		}
		return c.Status(200).JSON(response)
	}
}

// @Summary Retrieve Instructor Course Details
// @Description `This endpoint retrieves the details of a particular course for the authenticated instructor`
// @Tags Instructor
// @Param slug path string true "Course Slug"
// @Success 200 {object} courses.CourseResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Router /instructor/courses/{slug} [get]
// @Security BearerAuth
func GetInstructorCourseDetails(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		user := base.RequestUser(c)
		course := courseManager.GetCourseBySlug(db, ctx, c.Params("slug"), user, true)
		if course == nil {
			return config.APIError(c, 404, config.NotFoundErr("Instructor has no course with that slug"))
		}
		response := courses.CourseResponseSchema{
			ResponseSchema: base.ResponseMessage("Course Details Fetched Successfully"),
			Data:           courses.CourseDetailSchema{}.Assign(course),
		}
		return c.Status(200).JSON(response)
	}
}

// @Summary Update A Course
// @Description `This endpoint allows an instructor to update a course`
// @Tags Instructor
// @Param slug path string true "Course Slug"
// @Param course formData CourseCreateSchema true "Course object"
// @Param thumbnail formData file false "Thumbnail to upload"
// @Param intro_video formData file false "Intro video to upload"
// @Success 200 {object} courses.CourseResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Router /instructor/courses/{slug} [put]
// @Security BearerAuth
func UpdateCourse(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		user := base.RequestUser(c)
		course := courseManager.GetCourseBySlug(db, ctx, c.Params("slug"), user, true)
		if course == nil {
			return config.APIError(c, 404, config.NotFoundErr("Instructor has no course with that slug"))
		}

		data := CourseCreateSchema{}
		if errCode, errData := config.ValidateFormRequest(c, &data); errData != nil {
			return config.APIError(c, *errCode, *errData)
		}
		category := courseManager.GetCategoryBySlug(db, ctx, data.CategorySlug)
		if category == nil {
			return config.APIError(c, 422, config.ValidationErr("categorySlug", "Invalid category slug"))
		}

		// Check and validate files
		thumbnail, err := config.ValidateFile(c, "thumbnail", false)
		if err != nil {
			return c.Status(422).JSON(err)
		}
		introVideo, err := config.ValidateFile(c, "intro_video", false)
		if err != nil {
			return c.Status(422).JSON(err)
		}
		var (
			thumbnailUrl  *string
			introVideoUrl *string
		)
		if thumbnail != nil {
			url := config.UploadFile(thumbnail, string(config.FF_THUMBNAIL))
			thumbnailUrl = &url
		}
		if introVideo != nil {
			url := config.UploadFile(introVideo, string(config.FF_INTRO_VIDEOS))
			introVideoUrl = &url
		}
		updatedCourse := instructorManager.UpdateCourse(db, ctx, course, category, thumbnailUrl, introVideoUrl, data)
		response := courses.CourseResponseSchema{
			ResponseSchema: base.ResponseMessage("Course Updated Successfully"),
			Data:           courses.CourseDetailSchema{}.Assign(updatedCourse),
		}
		return c.Status(200).JSON(response)
	}
}

// @Summary Delete A Course
// @Description `This endpoint allows an authenticated instructor to delete a course`
// @Tags Instructor
// @Param slug path string true "Course Slug"
// @Success 200 {object} courses.CourseResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Router /instructor/courses/{slug} [delete]
// @Security BearerAuth
func DeleteACourse(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		user := base.RequestUser(c)
		course := courseManager.GetCourseBySlug(db, ctx, c.Params("slug"), user, false)
		if course == nil {
			return config.APIError(c, 404, config.NotFoundErr("Instructor has no course with that slug"))
		}
		err := instructorManager.DeleteCourse(db, ctx, course)
		if err != nil {
			return config.APIError(c, 403, config.RequestErr(config.ERR_NOT_ALLOWED, *err))
		}
		return c.Status(200).JSON(base.ResponseMessage("Course Deleted successfully"))
	}
}