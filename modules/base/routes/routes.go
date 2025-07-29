package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/user"
	"github.com/kayprogrammer/ednet-fiber-api/modules/accounts"
	"github.com/kayprogrammer/ednet-fiber-api/modules/courses"
	"github.com/kayprogrammer/ednet-fiber-api/modules/general"
	"github.com/kayprogrammer/ednet-fiber-api/modules/instructors"
	"github.com/kayprogrammer/ednet-fiber-api/modules/profiles"
)

// All Endpoints (51)
func SetupRoutes(app *fiber.App, db *ent.Client, cfg config.Config) {

	api := app.Group("/api/v1")
	// HealthCheck Route (1)
	api.Get("/healthcheck", HealthCheck)

	// General Routes (1)
	generalRouter := api.Group("/general")
	generalRouter.Get("/site-detail", general.GetSiteDetails(db))

	// Auth Routes (10)
	authRouter := api.Group("/auth")
	authRouter.Post("/register", accounts.Register(db))
	authRouter.Post("/verify-email", accounts.VerifyEmail(db))
	authRouter.Post("/resend-verification-email", accounts.ResendVerificationEmail(db))
	authRouter.Post("/send-password-reset-otp", accounts.SendPasswordResetOtp(db))
	authRouter.Post("/set-new-password", accounts.SetNewPassword(db))
	authRouter.Post("/login", accounts.Login(db))
	authRouter.Post("/google-login", accounts.GoogleLogin(db))
	authRouter.Post("/refresh", accounts.Refresh(db))
	authRouter.Get("/logout", accounts.AuthMiddleware(db), accounts.Logout(db))
	authRouter.Get("/logout/all", accounts.AuthMiddleware(db), accounts.LogoutAll(db))

	// Profiles Routes (7)
	profilesRouter := api.Group("/profiles")
	profilesRouter.Get("", accounts.AuthMiddleware(db), profiles.GetProfile(db))
	profilesRouter.Put("", accounts.AuthMiddleware(db), profiles.UpdateProfile(db))
	profilesRouter.Get("/courses", accounts.AuthMiddleware(db), profiles.GetEnrolledCourses(db))
	profilesRouter.Get("/courses/:slug/progress", accounts.AuthMiddleware(db), profiles.GetCourseProgress(db))

	profilesRouter.Post("/lessons/:slug/progress", accounts.AuthMiddleware(db), profiles.CreateOrUpdateLessonProgress(db))
	profilesRouter.Get("/lessons/:slug/progress", accounts.AuthMiddleware(db), profiles.GetLessonProgress(db))
	profilesRouter.Get("/leaderboard", accounts.AuthMiddleware(db), profiles.GetLeaderboard(db))

	// Courses Routes (17)
	coursesRouter := api.Group("/courses")
	coursesRouter.Get("", courses.GetLatestCourses(db))
	coursesRouter.Post("/pdf/summarize", accounts.AuthMiddleware(db), courses.PostSummarizePDF(db, cfg))
	coursesRouter.Get("/:slug", courses.GetCourseDetails(db))
	coursesRouter.Get("/:slug/lessons", courses.GetCourseLessons(db))
	coursesRouter.Get("/:course_slug/lessons/:lesson_slug", courses.GetCourseLessonDetails(db))
	coursesRouter.Post("/:slug/enroll", accounts.AuthMiddleware(db), courses.EnrollForACourse(db, cfg))
	coursesRouter.Get("/lessons/:slug/quizzes", accounts.AuthMiddleware(db), courses.GetLessonQuizzes(db))
	coursesRouter.Get("/quizzes/:quiz_slug", accounts.AuthMiddleware(db), courses.GetLessonQuizDetails(db))
	coursesRouter.Get("/quizzes/:quiz_slug/start", accounts.AuthMiddleware(db), courses.StartQuiz(db))
	coursesRouter.Get("/quizzes/:quiz_slug/results", accounts.AuthMiddleware(db), courses.GetQuizResult(db))
	coursesRouter.Post("/quizzes/:quiz_slug/results", accounts.AuthMiddleware(db), courses.SubmitQuizResult(db))
	coursesRouter.Post("/webhook/stripe", courses.StripeWebhook(db, cfg))

	// Reviews
	coursesRouter.Get("/:slug/reviews", courses.GetCourseReviews(db))
	coursesRouter.Post("/:slug/reviews", accounts.AuthMiddleware(db), courses.CreateCourseReview(db))
	coursesRouter.Get("/reviews/:id", courses.GetCourseReview(db))
	coursesRouter.Put("/reviews/:id", accounts.AuthMiddleware(db), courses.UpdateCourseReview(db))
	coursesRouter.Delete("/reviews/:id", accounts.AuthMiddleware(db), courses.DeleteCourseReview(db))


	// Instructor Routes (15)
	instructorsRouter := api.Group("/instructor", accounts.AuthMiddleware(db, user.RoleInstructor))
	instructorsRouter.Get("/courses", instructors.GetInstructorCourses(db))
	instructorsRouter.Post("/courses", instructors.CreateCourse(db))
	instructorsRouter.Get("/courses/:slug", instructors.GetInstructorCourseDetails(db))
	instructorsRouter.Put("/courses/:slug", instructors.UpdateCourse(db))
	instructorsRouter.Delete("/courses/:slug", instructors.DeleteACourse(db))
	instructorsRouter.Get("/courses/:slug/lessons", instructors.GetInstructorCourseLessons(db))
	instructorsRouter.Post("/courses/:slug/lessons", instructors.CreateInstructorCourseLesson(db))

	instructorsRouter.Get("/lessons/:slug", instructors.GetInstructorCourseLessonDetails(db))
	instructorsRouter.Put("/lessons/:slug", instructors.UpdateCourseLesson(db))
	instructorsRouter.Delete("/lessons/:slug", instructors.DeleteCourseLesson(db))

	instructorsRouter.Get("/courses/:slug/quizzes", instructors.GetInstructorLessonQuizzes(db))
	instructorsRouter.Post("/courses/:slug/quizzes", instructors.CreateInstructorLessonQuiz(db))
	instructorsRouter.Get("/quizzes/:slug", instructors.GetInstructorLessonQuizDetails(db))
	instructorsRouter.Put("/quizzes/:slug", instructors.UpdateLessonQuiz(db))
	instructorsRouter.Delete("/quizzes/:slug", instructors.DeleteLessonQuiz(db))
}

type HealthCheckSchema struct {
	Success string `json:"success" example:"pong"`
}

// @Summary HealthCheck
// @Description This endpoint checks the health of our application.
// @Tags HealthCheck
// @Success 200 {object} HealthCheckSchema
// @Router /healthcheck [get]
func HealthCheck(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{"success": "pong"})
}
