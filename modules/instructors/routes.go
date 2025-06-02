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
// @Description This endpoint retrieves paginated responses of the authenticated instructor courses
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

// @Summary Retrieve Instructor Course Details
// @Description This endpoint retrieves the details of a particular course for the authenticated instructor
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
