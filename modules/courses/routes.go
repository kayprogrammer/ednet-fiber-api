package courses

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/enrollment"
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
// @Param sortByRating query string false "Sort By Rating (asc or desc)"
// @Success 200 {object} CoursesResponseSchema
// @Router /courses [get]
func GetLatestCourses(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		courses := courseManager.GetAllPaginated(db, c)
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
		course := courseManager.GetCourseBySlug(db, ctx, c.Params("slug"), nil, true)
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
		course := courseManager.GetCourseBySlug(db, c.Context(), c.Params("slug"), nil, false)
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
// @Param course_slug path string true "Course Slug"
// @Param lesson_slug path string true "Lesson Slug"
// @Success 200 {object} LessonResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Router /courses/{course_slug}/lessons/{lesson_slug} [get]
func GetCourseLessonDetails(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		lesson := courseManager.GetCourseLessonBySlug(db, ctx, c.Params("lesson_slug"), nil, true)
		if lesson == nil {
			return config.APIError(c, 404, config.NotFoundErr("Lesson Not Found"))
		}
		if lesson.Edges.Course.Slug != c.Params("course_slug") {
			return config.APIError(c, 404, config.NotFoundErr("Lesson Not Found for specified course"))
		}
		response := LessonResponseSchema{
			ResponseSchema: base.ResponseMessage("Lesson Details Fetched Successfully"),
			Data:           LessonDetailSchema{}.Assign(lesson),
		}
		return c.Status(200).JSON(response)
	}
}

// @Summary Enroll for a course
// @Description This endpoint allows a user to enroll for a specific course
// @Tags Courses
// @Param slug path string true "Course Slug"
// @Param enrollment body EnrollForACourseSchema true "Enrollment object"
// @Success 200 {object} EnrollmentResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Success 400 {object} base.InvalidErrorExample
// @Failure 422 {object} base.ValidationErrorExample
// @Router /courses/{slug}/enroll [post]
// @Security BearerAuth
func EnrollForACourse(db *ent.Client, cfg config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		user := base.RequestUser(c)
		course := courseManager.GetCourseBySlug(db, ctx, c.Params("slug"), nil, true)
		if course == nil {
			return config.APIError(c, 404, config.NotFoundErr("Course Not Found"))
		}
		data := EnrollForACourseSchema{}
		// Validate request
		if errCode, errData := config.ValidateRequest(c, &data); errData != nil {
			return config.APIError(c, *errCode, *errData)
		}

		enrollment, err := courseManager.CreateEnrollment(db, ctx, user, course)
		if err != nil {
			return config.APIError(c, 400, *err)
		}

		checkoutUrl, err := CreateCheckoutSession(cfg, course, data.SuccessUrl, data.CancelUrl, enrollment)
		if err != nil {
			return config.APIError(c, 500, *err)
		}
		enrollment.Update().SetCheckoutURL(*checkoutUrl).SaveX(ctx)
		enrollment.CheckoutURL = *checkoutUrl

		response := EnrollmentResponseSchema{
			ResponseSchema: base.ResponseMessage("Enrollment Created Successfully"),
			Data:           EnrollmentSchema{}.Assign(enrollment),
		}
		return c.Status(200).JSON(response)
	}
}

// @Summary Retrieve Lesson Quizzes
// @Description This endpoint retrieves paginated responses of a lesson quizzes
// @Tags Courses
// @Param slug path string true "Lesson Slug"
// @Param page query int false "Current Page" default(1)
// @Param limit query int false "Page Limit" default(100)
// @Param title query string false "Filter By Title"
// @Success 404 {object} base.NotFoundErrorExample
// @Success 200 {object} QuizzesResponseSchema
// @Router /courses/lessons/{slug}/quizzes [get]
// @Security BearerAuth
func GetLessonQuizzes(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		user := base.RequestUser(c)
		lesson := courseManager.GetCourseLessonBySlug(db, ctx, c.Params("slug"), nil, true)
		if lesson == nil {
			return config.APIError(c, 404, config.NotFoundErr("Lesson Not Found"))
		}
		// Check if user is enrolled for this course
		enrollmentObj := courseManager.GetExistentEnrollmentByUserAndCourse(db, ctx, user, lesson.Edges.Course, false)
		if enrollmentObj == nil || enrollmentObj.PaymentStatus != enrollment.PaymentStatusSuccessful {
			return config.APIError(c, 403, config.ForbiddenErr("Only for enrolled users"))
		}
		quizzes := courseManager.GetQuizzes(db, lesson, c)

		response := QuizzesResponseSchema{
			ResponseSchema: base.ResponseMessage("Quizzes Fetched Successfully"),
		}.Assign(quizzes)
		return c.Status(200).JSON(response)
	}
}

// @Summary Retrieve Quiz Details
// @Description This endpoint retrieves the details of a particular quiz
// @Tags Courses
// @Param quiz_slug path string true "Quiz Slug"
// @Success 200 {object} QuizResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Router /courses/quizzes/{quiz_slug} [get]
// @Security BearerAuth
func GetLessonQuizDetails(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		user := base.RequestUser(c)
		quiz := courseManager.GetQuizBySlug(db, ctx, c.Params("quiz_slug"), nil, true)
		if quiz == nil {
			return config.APIError(c, 404, config.NotFoundErr("Quiz Not Found"))
		}
		// Check if user is enrolled for this course
		enrollmentObj := courseManager.GetExistentEnrollmentByUserAndCourse(db, ctx, user, quiz.Edges.Lesson.Edges.Course, false)
		if enrollmentObj == nil || enrollmentObj.PaymentStatus != enrollment.PaymentStatusSuccessful {
			return config.APIError(c, 403, config.ForbiddenErr("Only for enrolled users"))
		}
		response := QuizResponseSchema{
			ResponseSchema: base.ResponseMessage("Quiz Details Fetched Successfully"),
			Data:           QuizDetailSchema{}.Assign(quiz),
		}
		return c.Status(200).JSON(response)
	}
}

// @Summary Start Quiz
// @Description `This endpoint allows a user to start a quiz`
// @Tags Courses
// @Param quiz_slug path string true "Quiz Slug"
// @Success 200 {object} base.ResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Success 400 {object} base.InvalidErrorExample
// @Router /courses/quizzes/{quiz_slug}/start [get]
// @Security BearerAuth
func StartQuiz(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		user := base.RequestUser(c)
		quiz := courseManager.GetQuizBySlug(db, ctx, c.Params("quiz_slug"), nil, true)
		if quiz == nil {
			return config.APIError(c, 404, config.NotFoundErr("Quiz Not Found"))
		}
		// Check if user is enrolled for this course
		enrollmentObj := courseManager.GetExistentEnrollmentByUserAndCourse(db, ctx, user, quiz.Edges.Lesson.Edges.Course, false)
		if enrollmentObj == nil || enrollmentObj.PaymentStatus != enrollment.PaymentStatusSuccessful {
			return config.APIError(c, 403, config.ForbiddenErr("Only for enrolled users"))
		}

		_, err := courseManager.CreateQuizResultData(db, ctx, user, quiz)
		if err != nil {
			return config.APIError(c, 400, *err)
		}
		return c.Status(200).JSON(base.ResponseMessage("Quiz started successfully"))
	}
}

// @Summary Submit Quiz
// @Description `This endpoint allows a user to submit their answers for a quiz`
// @Description `If this submission is for the last quiz in the course, a certificate will be generated`
// @Tags Courses
// @Param quiz_slug path string true "Quiz Slug"
// @Param result body QuizSubmissionSchema true "Submission object"
// @Success 200 {object} QuizResultResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Success 400 {object} base.InvalidErrorExample
// @Failure 422 {object} base.ValidationErrorExample
// @Router /courses/quizzes/{quiz_slug}/results [post]
// @Security BearerAuth
func SubmitQuizResult(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		user := base.RequestUser(c)
		quiz := courseManager.GetQuizBySlug(db, ctx, c.Params("quiz_slug"), nil, true)
		if quiz == nil {
			return config.APIError(c, 404, config.NotFoundErr("Quiz Not Found"))
		}

		quizResult := courseManager.GetQuizResult(db, ctx, user, quiz.ID)
		if quizResult == nil {
			return config.APIError(c, 404, config.NotFoundErr("You cannot submit a quiz you didn't start"))
		}
		data := QuizSubmissionSchema{}
		// Validate request
		if errCode, errData := config.ValidateRequest(c, &data); errData != nil {
			return config.APIError(c, *errCode, *errData)
		}

		if quizResult.CompletedAt != nil {
			return config.APIError(c, 400, config.RequestErr(config.ERR_NOT_ALLOWED, "You have already submitted this quiz"))
		}

		quizResult, err := courseManager.SaveQuizResult(db, ctx, user, quiz, quizResult, data)
		if err != nil {
			return config.APIError(c, 400, *err)
		}

		// Generate cert if this is the last quiz in the course
		if courseManager.IsLastQuizInCourse(db, ctx, quiz) {
			courseManager.GenerateCertificate(db, ctx, user, quiz.Edges.Lesson.QueryCourse().WithInstructor().OnlyX(ctx))
		}
		
		response := QuizResultResponseSchema{
			ResponseSchema: base.ResponseMessage("Quiz Submitted Successfully"),
			Data:           QuizResultSchema{}.Assign(quizResult),
		}
		return c.Status(200).JSON(response)
	}
}

// @Summary Retrieve Quiz Result
// @Description `This endpoint retrieves the result of a particular quiz for a user`
// @Tags Courses
// @Param quiz_slug path string true "Quiz Slug"
// @Success 200 {object} QuizResultResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Router /courses/quizzes/{quiz_slug}/results [get]
// @Security BearerAuth
func GetQuizResult(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		user := base.RequestUser(c)
		quiz := courseManager.GetQuizBySlug(db, ctx, c.Params("quiz_slug"), nil, true)
		if quiz == nil {
			return config.APIError(c, 404, config.NotFoundErr("Quiz Not Found"))
		}
		quizResult := courseManager.GetQuizResult(db, ctx, user, quiz.ID)
		if quizResult == nil {
			return config.APIError(c, 404, config.NotFoundErr("You cannot submit a quiz you didn't start"))
		}
		response := QuizResultResponseSchema{
			ResponseSchema: base.ResponseMessage("Quiz Result Fetched Successfully"),
			Data:           QuizResultSchema{}.Assign(quizResult),
		}
		return c.Status(200).JSON(response)
	}
}

// @Summary Summarize a PDF document
// @Description `This endpoint accepts a PDF file and returns a summarized version of its content.`
// @Tags Courses
// @Accept multipart/form-data
// @Param file formData file true "PDF file to summarize"
// @Param max_points query int false "Maximum number of summary points" default(30)
// @Success 200 {object} PDFSummaryResponseSchema
// @Failure 400 {object} base.InvalidErrorExample
// @Router /courses/pdf/summarize [post]
// @Security BearerAuth
func PostSummarizePDF(db *ent.Client, cfg config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := base.RequestUser(c)
		// Get max_points from query, default to 30
		maxPoints, _ := strconv.Atoi(c.Query("max_points", "30"))
		if maxPoints > 100 {
			maxPoints = 100
		}

		// Check daily limit
		now := time.Now()
		if user.LastSummaryDate != nil && user.LastSummaryDate.Year() == now.Year() && user.LastSummaryDate.YearDay() == now.YearDay() {
			if user.SummaryCount >= 10 {
				return config.APIError(c, fiber.StatusTooManyRequests, config.RequestErr(config.ERR_TOO_MANY_REQUESTS, "You have reached your daily limit of 10 PDF summaries."))
			}
		} else {
			// Reset count for a new day
			user.Update().SetSummaryCount(0).ExecX(c.Context())
		}

		summary, status, errData := SummarizePDF(c, cfg, maxPoints)
		if errData != nil {
			return config.APIError(c, status, *errData)
		}

		// Update user's summary count and date
		if user.LastSummaryDate == nil || user.LastSummaryDate.Year() != now.Year() || user.LastSummaryDate.YearDay() != now.YearDay() {
			user.Update().SetLastSummaryDate(now).SetSummaryCount(1).ExecX(c.Context())
		} else {
			user.Update().SetSummaryCount(user.SummaryCount + 1).ExecX(c.Context())
		}

		return c.Status(200).JSON(
			PDFSummaryResponseSchema{
				ResponseSchema: base.ResponseMessage("PDF Summarized Successfully"),
				Data: PDFSummarySchema{
					Summary: summary,
					Points:  maxPoints,
				},
			},
		)
	}
}

// @Summary Retrieve Course Reviews
// @Description This endpoint retrieves paginated responses of a course reviews
// @Tags Courses
// @Param slug path string true "Course Slug"
// @Param page query int false "Current Page" default(1)
// @Param limit query int false "Page Limit" default(100)
// @Success 404 {object} base.NotFoundErrorExample
// @Success 200 {object} ReviewsResponseSchema
// @Router /courses/{slug}/reviews [get]
func GetCourseReviews(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		course := courseManager.GetCourseBySlug(db, c.Context(), c.Params("slug"), nil, false)
		if course == nil {
			return config.APIError(c, 404, config.NotFoundErr("Course Not Found"))
		}
		reviews := courseManager.GetReviews(db, course, c)

		response := ReviewsResponseSchema{
			ResponseSchema: base.ResponseMessage("Reviews Fetched Successfully"),
		}.Assign(reviews)
		return c.Status(200).JSON(response)
	}
}

// @Summary Create a review for a course
// @Description This endpoint allows a user to create a review for a specific course
// @Tags Courses
// @Param slug path string true "Course Slug"
// @Param review body ReviewSchema true "Review object"
// @Success 201 {object} ReviewResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Success 400 {object} base.InvalidErrorExample
// @Failure 422 {object} base.ValidationErrorExample
// @Router /courses/{slug}/reviews [post]
// @Security BearerAuth
func CreateCourseReview(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		user := base.RequestUser(c)
		course := courseManager.GetCourseBySlug(db, ctx, c.Params("slug"), nil, false)
		if course == nil {
			return config.APIError(c, 404, config.NotFoundErr("Course Not Found"))
		}
		data := ReviewSchema{}
		// Validate request
		if errCode, errData := config.ValidateRequest(c, &data); errData != nil {
			return config.APIError(c, *errCode, *errData)
		}

		review, err := courseManager.CreateReview(db, ctx, user, course, data)
		if err != nil {
			return config.APIError(c, 400, *err)
		}

		response := ReviewResponseSchema{
			ResponseSchema: base.ResponseMessage("Review Created Successfully"),
			Data:           ReviewResponseData{}.Assign(review),
		}
		return c.Status(201).JSON(response)
	}
}

// @Summary Retrieve a Review
// @Description This endpoint retrieves a specific review
// @Tags Courses
// @Param id path string true "Review ID"
// @Success 200 {object} ReviewResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Router /courses/reviews/{id} [get]
func GetCourseReview(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		reviewID, _ := uuid.Parse(c.Params("id"))
		review := courseManager.GetReview(db, ctx, reviewID)
		if review == nil {
			return config.APIError(c, 404, config.NotFoundErr("Review Not Found"))
		}
		response := ReviewResponseSchema{
			ResponseSchema: base.ResponseMessage("Review Fetched Successfully"),
			Data:           ReviewResponseData{}.Assign(review),
		}
		return c.Status(200).JSON(response)
	}
}

// @Summary Update a review
// @Description This endpoint allows a user to update their review for a specific course
// @Tags Courses
// @Param id path string true "Review ID"
// @Param review body ReviewSchema true "Review object"
// @Success 200 {object} ReviewResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Failure 422 {object} base.ValidationErrorExample
// @Router /courses/reviews/{id} [put]
// @Security BearerAuth
func UpdateCourseReview(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		user := base.RequestUser(c)
		reviewID, _ := uuid.Parse(c.Params("id"))
		review := courseManager.GetReview(db, ctx, reviewID)
		if review == nil {
			return config.APIError(c, 404, config.NotFoundErr("Review Not Found"))
		}
		if review.Edges.User.ID != user.ID {
			return config.APIError(c, 403, config.ForbiddenErr("You are not authorized to perform this action"))
		}
		data := ReviewSchema{}
		// Validate request
		if errCode, errData := config.ValidateRequest(c, &data); errData != nil {
			return config.APIError(c, *errCode, *errData)
		}

		review = courseManager.UpdateReview(db, ctx, review, data)

		response := ReviewResponseSchema{
			ResponseSchema: base.ResponseMessage("Review Updated Successfully"),
			Data:           ReviewResponseData{}.Assign(review),
		}
		return c.Status(200).JSON(response)
	}
}

// @Summary Delete a review
// @Description This endpoint allows a user to delete their review for a specific course
// @Tags Courses
// @Param id path string true "Review ID"
// @Success 200 {object} base.ResponseSchema
// @Success 404 {object} base.NotFoundErrorExample
// @Router /courses/reviews/{id} [delete]
// @Security BearerAuth
func DeleteCourseReview(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		user := base.RequestUser(c)
		reviewID, _ := uuid.Parse(c.Params("id"))
		review := courseManager.GetReview(db, ctx, reviewID)
		if review == nil {
			return config.APIError(c, 404, config.NotFoundErr("Review Not Found"))
		}
		if review.Edges.User.ID != user.ID {
			return config.APIError(c, 403, config.ForbiddenErr("You are not authorized to perform this action"))
		}
		db.Review.DeleteOne(review).ExecX(ctx)
		return c.Status(200).JSON(base.ResponseMessage("Review Deleted Successfully"))
	}
}