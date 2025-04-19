package courses

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/modules/base"
)

var courseManager = CourseManager{}

// @Summary Retrieve Latest Courses
// @Description This endpoint retrieves paginated responses of latest courses
// @Tags Courses
// @Param page query int false "Current Page" default(1)
// @Param limit query int false "Page Limit" default(100)
// @Param title query string false "Filter By Title"
// @Param instructor query string false "Filter By Instructor's Name Or Username"
// @Param isFree query bool false "Filter By Free Status"
// @Success 200 {object} CoursesResponseSchema
// @Router /courses [get]
func GetLatestCourses(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		courses := courseManager.GetAll(db, c)
		response := CoursesResponseSchema{
			ResponseSchema: base.ResponseMessage("Courses Fetched Successfully"),
		}.Assign(courses)
		return c.Status(200).JSON(response)
	}
}

// @Summary Retrieve Course Details
// @Description This endpoint retrieves the details of a particular course
// @Tags Courses
// @Param slug path string true "Course Slug"
// @Success 200 {object} CourseResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Router /courses/{slug} [get]
func GetCourseDetails(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		course := courseManager.GetCourseBySlug(db, ctx, c.Params("slug"), true)
		if course == nil {
			return config.APIError(c, 404, config.NotFoundErr("Course Not Found"))
		}
		response := CourseResponseSchema{
			ResponseSchema: base.ResponseMessage("Course Details Fetched Successfully"),
			Data:           CourseDetailSchema{}.Assign(course),
		}
		return c.Status(200).JSON(response)
	}
}

// @Summary Retrieve Course Lessons
// @Description This endpoint retrieves paginated responses of a course lessons
// @Tags Courses
// @Param slug path string true "Course Slug"
// @Param page query int false "Current Page" default(1)
// @Param limit query int false "Page Limit" default(100)
// @Param title query string false "Filter By Title"
// @Param isFreePreview query bool false "Filter By Free Preview"
// @Success 404 {object} base.NotFoundErrorExample
// @Success 200 {object} LessonsResponseSchema
// @Router /courses/{slug}/lessons [get]
func GetCourseLessons(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		course := courseManager.GetCourseBySlug(db, c.Context(), c.Params("slug"), false)
		if course == nil {
			return config.APIError(c, 404, config.NotFoundErr("Course Not Found"))
		}
		lessons := courseManager.GetLessons(db, course, c)

		response := LessonsResponseSchema{
			ResponseSchema: base.ResponseMessage("Lessons Fetched Successfully"),
		}.Assign(lessons)
		return c.Status(200).JSON(response)
	}
}

// @Summary Retrieve Lesson Details
// @Description This endpoint retrieves the details of a particular lesson
// @Tags Courses
// @Param slug path string true "Course Slug"
// @Success 200 {object} LessonResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Router /courses/{course_slug}/lessons/{lesson_slug} [get]
func GetCourseLessonDetails(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		lesson := courseManager.GetCourseLessonBySlug(db, ctx, c.Params("lesson_slug"), true)
		if lesson == nil {
			return config.APIError(c, 404, config.NotFoundErr("Lesson Not Found"))
		}
		if lesson.Edges.Course.Slug != c.Params("course_slug") {
			return config.APIError(c, 404, config.NotFoundErr("Lesson Not Found for specified course"))
		}
		response := LessonResponseSchema{
			ResponseSchema: base.ResponseMessage("Course Details Fetched Successfully"),
			Data:           LessonDetailSchema{}.Assign(lesson),
		}
		return c.Status(200).JSON(response)
	}
}