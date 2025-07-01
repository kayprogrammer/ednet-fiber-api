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

// @Summary Retrieve Courses
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
// @Success 201 {object} courses.CourseResponseSchema
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
		thumbnail, err := config.ValidateFile(c, "thumbnail", true, false)
		if err != nil {
			return c.Status(422).JSON(err)
		}
		introVideo, err := config.ValidateFile(c, "intro_video", true, true)
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
		return c.Status(201).JSON(response)
	}
}

// @Summary Retrieve Course Details
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
		thumbnail, err := config.ValidateFile(c, "thumbnail", false, false)
		if err != nil {
			return c.Status(422).JSON(err)
		}
		introVideo, err := config.ValidateFile(c, "intro_video", false, true)
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
// @Success 200 {object} base.ResponseSchema
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

// @Summary Retrieve Course Lessons
// @Description `This endpoint retrieves the lessons of a particular course for the authenticated instructor`
// @Tags Instructor
// @Param slug path string true "Course Slug"
// @Success 200 {object} courses.LessonsResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Param page query int false "Current Page" default(1)
// @Param limit query int false "Page Limit" default(100)
// @Param title query string false "Filter By Title"
// @Param isFreePreview query bool false "Filter By Free Preview"
// @Router /instructor/courses/{slug}/lessons [get]
// @Security BearerAuth
func GetInstructorCourseLessons(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		user := base.RequestUser(c)
		course := courseManager.GetCourseBySlug(db, ctx, c.Params("slug"), user, true)
		if course == nil {
			return config.APIError(c, 404, config.NotFoundErr("Instructor has no course with that slug"))
		}
		lessons := courseManager.GetLessons(db, course, c)

		response := courses.LessonsResponseSchema{
			ResponseSchema: base.ResponseMessage("Lessons Fetched Successfully"),
		}.Assign(lessons)
		return c.Status(200).JSON(response)
	}
}

// @Summary Create Course Lesson
// @Description `This endpoint creates a lesson of a particular course for the authenticated instructor`
// @Tags Instructor
// @Param slug path string true "Course Slug"
// @Param lesson formData LessonCreateSchema true "Lesson object"
// @Param thumbnail formData file true "Thumbnail to upload"
// @Param video formData file false "Video to upload"
// @Success 201 {object} courses.LessonResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Router /instructor/courses/{slug}/lessons [post]
// @Security BearerAuth
func CreateInstructorCourseLesson(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		user := base.RequestUser(c)
		course := courseManager.GetCourseBySlug(db, ctx, c.Params("slug"), user, true)
		if course == nil {
			return config.APIError(c, 404, config.NotFoundErr("Instructor has no course with that slug"))
		}

		data := LessonCreateSchema{}
		if errCode, errData := config.ValidateFormRequest(c, &data); errData != nil {
			return config.APIError(c, *errCode, *errData)
		}
		// Check and validate files
		thumbnail, err := config.ValidateFile(c, "thumbnail", true, false)
		if err != nil {
			return c.Status(422).JSON(err)
		}
		video, err := config.ValidateFile(c, "video", false, true)
		if err != nil {
			return c.Status(422).JSON(err)
		}
		thumbnailUrl := config.UploadFile(thumbnail, string(config.FF_THUMBNAIL))
		var videoUrl *string
		if video != nil {
			url := config.UploadFile(video, string(config.FF_LESSON_VIDEOS))
			videoUrl = &url
		}

		lesson := instructorManager.CreateLesson(db, ctx, course, thumbnailUrl, videoUrl, data)

		response := courses.LessonResponseSchema{
			ResponseSchema: base.ResponseMessage("Lesson Created Successfully"),
			Data: courses.LessonDetailSchema{}.Assign(lesson),
		}
		return c.Status(201).JSON(response)
	}
}

// @Summary Retrieve Instructor Lesson Details
// @Description This endpoint retrieves the details of a particular lesson belonging to an instructor
// @Tags Instructor
// @Param slug path string true "Lesson Slug"
// @Success 200 {object} courses.LessonResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Router /instructor/lessons/{slug} [get]
// @Security BearerAuth
func GetInstructorCourseLessonDetails(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		user := base.RequestUser(c)
		lesson := instructorManager.GetCourseLessonBySlug(db, ctx, user, c.Params("slug"), true)
		if lesson == nil {
			return config.APIError(c, 404, config.NotFoundErr("Instructor lesson not found"))
		}
		response := courses.LessonResponseSchema{
			ResponseSchema: base.ResponseMessage("Lesson Details Fetched Successfully"),
			Data:           courses.LessonDetailSchema{}.Assign(lesson),
		}
		return c.Status(200).JSON(response)
	}
}

// @Summary Update Course Lesson
// @Description `This endpoint updates a lesson of a particular course for the authenticated instructor`
// @Tags Instructor
// @Param slug path string true "Lesson Slug"
// @Param lesson formData LessonCreateSchema true "Lesson object"
// @Param thumbnail formData file false "Thumbnail to upload"
// @Param video formData file false "Video to upload"
// @Success 200 {object} courses.LessonResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Router /instructor/lessons/{slug} [put]
// @Security BearerAuth
func UpdateCourseLesson(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		user := base.RequestUser(c)
		lesson := instructorManager.GetCourseLessonBySlug(db, ctx, user, c.Params("slug"), false)
		if lesson == nil {
			return config.APIError(c, 404, config.NotFoundErr("Instructor lesson not found"))
		}

		data := LessonCreateSchema{}
		if errCode, errData := config.ValidateFormRequest(c, &data); errData != nil {
			return config.APIError(c, *errCode, *errData)
		}
		// Check and validate files
		thumbnail, err := config.ValidateFile(c, "thumbnail", false, false)
		if err != nil {
			return c.Status(422).JSON(err)
		}
		video, err := config.ValidateFile(c, "video", false, true)
		if err != nil {
			return c.Status(422).JSON(err)
		}
		var thumbnailUrl *string
		if thumbnail != nil {
			url := config.UploadFile(thumbnail, string(config.FF_THUMBNAIL))
			thumbnailUrl = &url
		}
		var videoUrl *string
		if video != nil {
			url := config.UploadFile(video, string(config.FF_LESSON_VIDEOS))
			videoUrl = &url
		}

		lesson = instructorManager.UpdateLesson(db, ctx, lesson, thumbnailUrl, videoUrl, data)

		response := courses.LessonResponseSchema{
			ResponseSchema: base.ResponseMessage("Lesson Updated Successfully"),
			Data: courses.LessonDetailSchema{}.Assign(lesson),
		}
		return c.Status(200).JSON(response)
	}
}

// @Summary Delete Instructor Lesson Details
// @Description This endpoint deletes a particular lesson belonging to an instructor
// @Tags Instructor
// @Param slug path string true "Lesson Slug"
// @Success 200 {object} base.ResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Router /instructor/lessons/{slug} [delete]
// @Security BearerAuth
func DeleteCourseLesson(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		user := base.RequestUser(c)
		lesson := instructorManager.GetCourseLessonBySlug(db, ctx, user, c.Params("slug"), false)
		if lesson == nil {
			return config.APIError(c, 404, config.NotFoundErr("Instructor lesson not found"))
		}
		err := instructorManager.DeleteLesson(db, ctx, lesson)
		if err != nil {
			return config.APIError(c, 403, config.RequestErr(config.ERR_NOT_ALLOWED, *err))
		}
		return c.Status(200).JSON(base.ResponseMessage("Lesson deleted successfully"))
	}
}

// @Summary Retrieve Course Quizzes
// @Description `This endpoint retrieves the quizzes of a particular course for the authenticated instructor`
// @Tags Instructor
// @Param slug path string true "Course Slug"
// @Success 200 {object} courses.QuizzesResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Param page query int false "Current Page" default(1)
// @Param limit query int false "Page Limit" default(100)
// @Param title query string false "Filter By Title"
// @Param isPublished query bool false "Filter By Published Status"
// @Router /instructor/courses/{slug}/quizzes [get]
// @Security BearerAuth
func GetInstructorCourseQuizzes(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		user := base.RequestUser(c)
		course := courseManager.GetCourseBySlug(db, ctx, c.Params("slug"), user, true)
		if course == nil {
			return config.APIError(c, 404, config.NotFoundErr("Instructor has no course with that slug"))
		}
		quizzes := courseManager.GetQuizzes(db, course, c)

		response := courses.QuizzesResponseSchema{
			ResponseSchema: base.ResponseMessage("Quizzes Fetched Successfully"),
		}.Assign(quizzes)
		return c.Status(200).JSON(response)
	}
}

// @Summary Create Course Quiz
// @Description `This endpoint creates a quiz of a particular course for the authenticated instructor`
// @Tags Instructor
// @Param slug path string true "Course Slug"
// @Param quiz body QuizCreateSchema true "Quiz object"
// @Success 201 {object} courses.QuizResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Router /instructor/courses/{slug}/quizzes [post]
// @Security BearerAuth
func CreateInstructorCourseQuiz(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		user := base.RequestUser(c)
		course := courseManager.GetCourseBySlug(db, ctx, c.Params("slug"), user, true)
		if course == nil {
			return config.APIError(c, 404, config.NotFoundErr("Instructor has no course with that slug"))
		}

		data := QuizCreateSchema{}
		if errCode, errData := config.ValidateRequest(c, &data); errData != nil {
			return config.APIError(c, *errCode, *errData)
		}

		quiz := instructorManager.CreateQuiz(db, ctx, course, data)

		response := courses.QuizResponseSchema{
			ResponseSchema: base.ResponseMessage("Quiz Created Successfully"),
			Data:           courses.QuizDetailSchema{}.Assign(quiz),
		}
		return c.Status(201).JSON(response)
	}
}

// @Summary Retrieve Instructor Quiz Details
// @Description This endpoint retrieves the details of a particular quiz belonging to an instructor
// @Tags Instructor
// @Param slug path string true "Quiz Slug"
// @Success 200 {object} courses.QuizResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Router /instructor/quizzes/{slug} [get]
// @Security BearerAuth
func GetInstructorCourseQuizDetails(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		user := base.RequestUser(c)
		quiz := instructorManager.GetCourseQuizBySlug(db, ctx, user, c.Params("slug"), true)
		if quiz == nil {
			return config.APIError(c, 404, config.NotFoundErr("Instructor quiz not found"))
		}
		response := courses.QuizResponseSchema{
			ResponseSchema: base.ResponseMessage("Quiz Details Fetched Successfully"),
			Data:           courses.QuizDetailSchema{}.Assign(quiz),
		}
		return c.Status(200).JSON(response)
	}
}

// @Summary Update Course Quiz
// @Description `This endpoint updates a quiz of a particular course for the authenticated instructor`
// @Tags Instructor
// @Param slug path string true "Quiz Slug"
// @Param quiz body QuizCreateSchema true "Quiz object"
// @Success 200 {object} courses.QuizResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Router /instructor/quizzes/{slug} [put]
// @Security BearerAuth
func UpdateCourseQuiz(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		user := base.RequestUser(c)
		quiz := instructorManager.GetCourseQuizBySlug(db, ctx, user, c.Params("slug"), false)
		if quiz == nil {
			return config.APIError(c, 404, config.NotFoundErr("Instructor quiz not found"))
		}

		data := QuizCreateSchema{}
		if errCode, errData := config.ValidateRequest(c, &data); errData != nil {
			return config.APIError(c, *errCode, *errData)
		}

		quiz = instructorManager.UpdateQuiz(db, ctx, quiz, user, data)

		response := courses.QuizResponseSchema{
			ResponseSchema: base.ResponseMessage("Quiz Updated Successfully"),
			Data:           courses.QuizDetailSchema{}.Assign(quiz),
		}
		return c.Status(200).JSON(response)
	}
}

// @Summary Delete Instructor Quiz
// @Description This endpoint deletes a particular quiz belonging to an instructor
// @Tags Instructor
// @Param slug path string true "Quiz Slug"
// @Success 200 {object} base.ResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Router /instructor/quizzes/{slug} [delete]
// @Security BearerAuth
func DeleteCourseQuiz(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		user := base.RequestUser(c)
		quiz := instructorManager.GetCourseQuizBySlug(db, ctx, user, c.Params("slug"), false)
		if quiz == nil {
			return config.APIError(c, 404, config.NotFoundErr("Instructor quiz not found"))
		}
		err := instructorManager.DeleteQuiz(db, ctx, quiz)
		if err != nil {
			return config.APIError(c, 403, config.RequestErr(config.ERR_NOT_ALLOWED, *err))
		}
		return c.Status(200).JSON(base.ResponseMessage("Quiz deleted successfully"))
	}
}
