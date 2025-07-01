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

// All Endpoints (50)
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

	// Profiles Routes (2)
	profilesRouter := api.Group("/profiles")
	profilesRouter.Get("", accounts.AuthMiddleware(db), profiles.GetProfile(db))
	profilesRouter.Put("", accounts.AuthMiddleware(db), profiles.UpdateProfile(db))
	profilesRouter.Get("/courses", accounts.AuthMiddleware(db), profiles.GetEnrolledCourses(db))

	// Courses Routes (2)
	coursesRouter := api.Group("/courses")
	coursesRouter.Get("", courses.GetLatestCourses(db))
	coursesRouter.Get("/:slug", courses.GetCourseDetails(db))
	coursesRouter.Get("/:slug/lessons", courses.GetCourseLessons(db))
	coursesRouter.Get("/:course_slug/lessons/:lesson_slug", courses.GetCourseLessonDetails(db))
	coursesRouter.Post("/:slug/enroll", accounts.AuthMiddleware(db), courses.EnrollForACourse(db, cfg))
	coursesRouter.Get("/:slug/quizzes", accounts.AuthMiddleware(db), courses.GetCourseQuizzes(db))
	coursesRouter.Get("/:course_slug/quizzes/:quiz_slug", accounts.AuthMiddleware(db), courses.GetCourseQuizDetails(db))

	// Instructor Routes (2)
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

	instructorsRouter.Get("/courses/:slug/quizzes", instructors.GetInstructorCourseQuizzes(db))
	instructorsRouter.Post("/courses/:slug/quizzes", instructors.CreateInstructorCourseQuiz(db))
	instructorsRouter.Get("/quizzes/:slug", instructors.GetInstructorCourseQuizDetails(db))
	instructorsRouter.Put("/quizzes/:slug", instructors.UpdateCourseQuiz(db))
	instructorsRouter.Delete("/quizzes/:slug", instructors.DeleteCourseQuiz(db))
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
