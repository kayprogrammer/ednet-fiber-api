
package courses

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/category"
	"github.com/kayprogrammer/ednet-fiber-api/ent/course"
	"github.com/kayprogrammer/ednet-fiber-api/ent/enrollment"
	"github.com/kayprogrammer/ednet-fiber-api/ent/lesson"
	"github.com/kayprogrammer/ednet-fiber-api/ent/questionoption"
	"github.com/kayprogrammer/ednet-fiber-api/ent/quiz"
	"github.com/kayprogrammer/ednet-fiber-api/ent/quizresult"
	"github.com/kayprogrammer/ednet-fiber-api/ent/review"
	"github.com/kayprogrammer/ednet-fiber-api/ent/user"
	"github.com/kayprogrammer/ednet-fiber-api/modules/courses/certs"
)

type CourseManager struct{}

func (c CourseManager) GetCategories(db *ent.Client, ctx context.Context) []*ent.Category {
	categories := db.Category.Query().AllX(ctx)
	return categories
}

func (c CourseManager) GetCategoryByName(db *ent.Client, ctx context.Context, name string) *ent.Category {
	category, _ := db.Category.Query().Where(category.NameEQ(name)).Only(ctx)
	return category
}

func (c CourseManager) GetCategoryBySlug(db *ent.Client, ctx context.Context, slug string) *ent.Category {
	category, _ := db.Category.Query().Where(category.SlugEQ(slug)).Only(ctx)
	return category
}

func (c CourseManager) ApplyCourseFilters(fibCtx *fiber.Ctx, query *ent.CourseQuery) *ent.CourseQuery {
	filters := map[string]func(string){
		"title": func(value string) { query.Where(course.TitleContainsFold(value)) },
		"instructor": func(value string) {
			query.Where(course.HasInstructorWith(user.Or(user.NameContainsFold(value), user.UsernameContainsFold(value))))
		},
		"isFree": func(value string) {
			if freeStatus, err := strconv.ParseBool(value); err == nil {
				query.Where(course.IsFreeEQ(freeStatus))
			}
		},
		"isPublished": func(value string) {
			if publishedStatus, err := strconv.ParseBool(value); err == nil {
				query.Where(course.IsPublishedEQ(publishedStatus))
			}
		},
	}
	// Apply filters dynamically
	for param, apply := range filters {
		if value := fibCtx.Query(param); value != "" {
			apply(value)
		}
	}

	sortBy := fibCtx.Query("sortByRating")
	if sortBy == "asc" || sortBy == "desc" {
		query = query.Order(func(s *sql.Selector) {
			// Use COALESCE to handle NULL values and avoid syntax errors
			avgExpr := fmt.Sprintf(
				"COALESCE((SELECT AVG(%s) FROM %s WHERE %s = %s.%s), 0)",
				review.FieldRating,
				review.Table,
				review.CourseColumn,
				course.Table,
				course.FieldID,
			)

			// Apply sort order based on variable
			if strings.ToLower(sortBy) == "desc" {
				s.OrderBy(sql.Desc(avgExpr))
			} else {
				s.OrderBy(avgExpr)
			}
		})
	}
	return query
}

func (c CourseManager) GetAll(db *ent.Client, ctx context.Context) []*ent.Course {
	courses := db.Course.Query().
		WithInstructor().
		WithCategory().
		WithTags().
		AllX(ctx)
	return courses
}

func (c CourseManager) GetAllPaginated(db *ent.Client, fibCtx *fiber.Ctx) *config.PaginationResponse[*ent.Course] {
	query := db.Course.Query().
		Where(course.IsPublishedEQ(true)).
		WithInstructor().
		WithCategory().
		WithTags().
		WithReviews().
		WithEnrollments().
		WithLessons()

	query = c.ApplyCourseFilters(fibCtx, query)
	return config.PaginateModel(fibCtx, query)
}

func (c CourseManager) GetCourseByName(db *ent.Client, ctx context.Context, name string) *ent.Course {
	course, _ := db.Course.Query().Where(course.TitleEQ(name)).First(ctx)
	return course
}

func (c CourseManager) FilterCoursesByInstructor(db *ent.Client, fibCtx *fiber.Ctx, instructor *ent.User) *config.PaginationResponse[*ent.Course] {
	query := db.Course.Query().
		Where(course.InstructorIDEQ(instructor.ID)).
		Where(course.IsPublishedEQ(true)).
		WithInstructor().
		WithCategory().
		WithTags()
	query = c.ApplyCourseFilters(fibCtx, query)
	courses := config.PaginateModel(fibCtx, query)
	return courses
}

func (c CourseManager) GetCourseBySlug(db *ent.Client, ctx context.Context, slug string, instructor *ent.User, loaded bool) *ent.Course {
	query := db.Course.Query().
		Where(course.SlugEQ(slug))
	if instructor != nil {
		query = query.Where(course.InstructorIDEQ(instructor.ID))
	}
	if loaded {
		query = query.
			WithInstructor().
			WithCategory().
			WithTags().
			WithEnrollments().
			WithReviews()
	}
	course, _ := query.Only(ctx)
	return course
}

func (c CourseManager) ApplyLessonFilters(fibCtx *fiber.Ctx, query *ent.LessonQuery) *ent.LessonQuery {
	filters := map[string]func(string){
		"title": func(value string) { query.Where(lesson.TitleContainsFold(value)) },
		"isFreePreview": func(value string) {
			if freeStatus, err := strconv.ParseBool(value); err == nil {
				query.Where(lesson.IsFreePreviewEQ(freeStatus))
			}
		},
	}
	// Apply filters dynamically
	for param, apply := range filters {
		if value := fibCtx.Query(param); value != "" {
			apply(value)
		}
	}
	return query
}

func (c CourseManager) GetLessons(db *ent.Client, course *ent.Course, fibCtx *fiber.Ctx) *config.PaginationResponse[*ent.Lesson] {
	query := db.Lesson.Query().Where(lesson.CourseID(course.ID)).Order(ent.Asc(lesson.FieldOrder))
	query = c.ApplyLessonFilters(fibCtx, query)
	lessons := config.PaginateModel(fibCtx, query)
	return lessons
}

func (c CourseManager) ApplyQuizFilters(fibCtx *fiber.Ctx, query *ent.QuizQuery) *ent.QuizQuery {
	filters := map[string]func(string){
		"title": func(value string) { query.Where(quiz.TitleContainsFold(value)) },
		"isPublished": func(value string) {
			if publishedStatus, err := strconv.ParseBool(value); err == nil {
				query.Where(quiz.IsPublishedEQ(publishedStatus))
			}
		},
	}
	// Apply filters dynamically
	for param, apply := range filters {
		if value := fibCtx.Query(param); value != "" {
			apply(value)
		}
	}
	return query
}

func (c CourseManager) GetQuizzes(db *ent.Client, lesson *ent.Lesson, fibCtx *fiber.Ctx) *config.PaginationResponse[*ent.Quiz] {
	query := db.Quiz.Query().Where(quiz.LessonID(lesson.ID)).Order(ent.Asc(quiz.FieldTitle)).WithQuestions()
	query = c.ApplyQuizFilters(fibCtx, query)
	quizzes := config.PaginateModel(fibCtx, query)
	return quizzes
}

func (c CourseManager) GetQuizBySlug(db *ent.Client, ctx context.Context, slug string, instructor *ent.User, loaded bool) *ent.Quiz {
	query := db.Quiz.Query().
		Where(quiz.SlugEQ(slug))

	if instructor != nil {
		query = query.Where(quiz.HasLessonWith(lesson.HasCourseWith(course.InstructorIDEQ(instructor.ID))))
	}
	if loaded {
		query = query.
			WithLesson(
				func(q *ent.LessonQuery) {
					q.WithCourse()
				},
			).
			WithQuestions(
				func(q *ent.QuestionQuery) {
					q.WithOptions()
				},
			)
	}
	quiz, _ := query.Only(ctx)
	return quiz
}

func (c CourseManager) GetCourseLessonBySlug(db *ent.Client, ctx context.Context, slug string, instructor *ent.User, loaded bool) *ent.Lesson {
	query := db.Lesson.Query().
		Where(lesson.SlugEQ(slug))

	if instructor != nil {
		query = query.Where(lesson.HasCourseWith(course.InstructorID(instructor.ID)))
	}
	if loaded {
		query = query.
			WithQuizzes().
			WithCourse()
	}
	lesson, _ := query.Only(ctx)
	return lesson
}

func (c CourseManager) GetExistentEnrollmentByUserAndCourse(db *ent.Client, ctx context.Context, user *ent.User, course *ent.Course, loaded bool) *ent.Enrollment {
	query := db.Enrollment.Query().
		Where(
			enrollment.UserIDEQ(user.ID),
			enrollment.CourseIDEQ(course.ID),
		)
	if loaded {
		query = query.
			WithCourse().
			WithUser()
	}
	enrollmentObj, _ := query.Only(ctx)
	return enrollmentObj
}

func (c CourseManager) CreateEnrollment(db *ent.Client, ctx context.Context, user *ent.User, course *ent.Course) (*ent.Enrollment, *config.ErrorResponse) {
	existentEnrollment := c.GetExistentEnrollmentByUserAndCourse(db, ctx, user, course, false)
	if existentEnrollment != nil {
		err := config.RequestErr(config.ERR_NOT_ALLOWED, "Enrollment has been created already")
		return nil, &err
	}
			enrollmentQuery := db.Enrollment.
		Create().
		SetCourse(course).
		SetUser(user)

	if course.IsFree {
		enrollmentQuery = enrollmentQuery.SetStatus(enrollment.StatusActive).
			SetPaymentStatus(enrollment.PaymentStatusSuccessful)
	}
	enrollmentObj := enrollmentQuery.SaveX(ctx)
	enrollmentObj.Edges.User = user
	enrollmentObj.Edges.Course = course
	return enrollmentObj, nil
}

func (c CourseManager) UpdateEnrollment(db *ent.Client, ctx context.Context, enrollmentID uuid.UUID, paymentStatus enrollment.PaymentStatus) {
	enrollmentObj, err := db.Enrollment.Query().Where(enrollment.ID(enrollmentID)).Only(ctx)
	if err != nil {
		log.Printf("Error fetching enrollment: %v", err)
		return
	}
	enrollmentStatus := enrollment.StatusInactive
	if paymentStatus == enrollment.PaymentStatusSuccessful {
		enrollmentStatus = enrollment.StatusActive
	}
	enrollmentObj.Update().SetPaymentStatus(paymentStatus).SetStatus(enrollment.Status(enrollmentStatus)).SaveX(ctx)
}

func (c CourseManager) GetAverageRating(reviews []*ent.Review) float64 {
	if len(reviews) == 0 {
		return 0.0
	}

	var total float64
	for _, review := range reviews {
		total += float64(review.Rating) // Assuming Rating is int or float
	}

	return total / float64(len(reviews))
}

func (c CourseManager) CreateQuizResultData(
	db *ent.Client, ctx context.Context, user *ent.User, quiz *ent.Quiz,
) (*ent.QuizResult, *config.ErrorResponse) {
	// Check if the user has already created a result for this quiz
	existingResult, _ := db.QuizResult.Query().
		Where(quizresult.UserIDEQ(user.ID), quizresult.QuizIDEQ(quiz.ID)).
		Only(ctx)
	if existingResult != nil {
		err := config.RequestErr(config.ERR_NOT_ALLOWED, "You have already done this quiz before")
		return nil, &err
	}
	quizResult := db.QuizResult.Create().SetUserID(user.ID).SetQuizID(quiz.ID).SaveX(ctx)
	return quizResult, nil
}

func (c CourseManager) SaveQuizResult(
	db *ent.Client, ctx context.Context, user *ent.User, quiz *ent.Quiz, result *ent.QuizResult, data QuizSubmissionSchema,
) (*ent.QuizResult, *config.ErrorResponse) {
	// Prepare option IDs to batch fetch them
	optionIDs := make([]uuid.UUID, 0, len(data.Answers))
	for _, ans := range data.Answers {
		optionIDs = append(optionIDs, ans.SelectedOptionID)
	}

	optionsMap := make(map[uuid.UUID]*ent.QuestionOption)
	options := db.QuestionOption.Query().
		Where(questionoption.IDIn(optionIDs...)).
		AllX(ctx)
	for _, opt := range options {
		optionsMap[opt.ID] = opt
	}

	// Score calculation
	total := len(data.Answers)
	correct := 0
	for _, ans := range data.Answers {
		if opt, ok := optionsMap[ans.SelectedOptionID]; ok && opt.IsCorrect {
			correct++
		}
	}
	score := (float64(correct) / float64(total)) * 100

	// Update QuizResult
	result = result.Update().
		SetScore(score).
		SetTimeTaken(data.TimeTaken).
		SetCompletedAt(time.Now()).
		SaveX(ctx)

	// Batch create answers
	bulk := make([]*ent.AnswerCreate, 0, total)
	for _, ans := range data.Answers {
		opt := optionsMap[ans.SelectedOptionID]
		isCorrect := false
		if opt != nil && opt.IsCorrect {
			isCorrect = true
		}
		bulk = append(bulk, db.Answer.Create().
			SetResult(result).
			SetQuestionID(ans.QuestionID).
			SetSelectedOptionID(ans.SelectedOptionID).
			SetIsCorrect(isCorrect),
		)
	}
	db.Answer.CreateBulk(bulk...).ExecX(ctx)
	result.Edges.Answers = result.QueryAnswers().AllX(ctx)
	return result, nil
}

func (c CourseManager) GetQuizResult(
	db *ent.Client, ctx context.Context, user *ent.User, quizID uuid.UUID,
) *ent.QuizResult {
	result, _ := db.QuizResult.Query().
		Where(
			quizresult.UserIDEQ(user.ID),
			quizresult.QuizIDEQ(quizID),
		).
		WithAnswers().
		Only(ctx)
	return result
}

func (c CourseManager) IsLastQuizInCourse(db *ent.Client, ctx context.Context, quizObj *ent.Quiz) bool {
	currentLesson := quizObj.Edges.Lesson
	courseID := currentLesson.CourseID

	lastQuiz := db.Quiz.Query().
		Where(
			quiz.HasLessonWith(lesson.CourseID(courseID)),
		).Order(ent.Desc(quiz.FieldCreatedAt)).
		FirstX(ctx)
	return quizObj.ID == lastQuiz.ID
}

func (c CourseManager) GenerateCertificate(db *ent.Client, ctx context.Context, user *ent.User, course *ent.Course) {
	cert := certs.GenerateCertificate(user, course, course.Edges.Instructor.Name)
	db.Enrollment.Update().Where(enrollment.CourseID(course.ID), enrollment.UserID(user.ID)).
		SetCert(cert).
		SaveX(ctx)
}

func (c CourseManager) GetReviews(db *ent.Client, course *ent.Course, fibCtx *fiber.Ctx) *config.PaginationResponse[*ent.Review] {
	query := db.Review.Query().Where(review.CourseID(course.ID)).Order(ent.Desc(review.FieldCreatedAt)).WithUser()
	reviews := config.PaginateModel(fibCtx, query)
	return reviews
}

func (c CourseManager) GetReview(db *ent.Client, ctx context.Context, reviewID uuid.UUID) *ent.Review {
	review, _ := db.Review.Query().Where(review.IDEQ(reviewID)).WithUser().Only(ctx)
	return review
}

func (c CourseManager) CreateReview(db *ent.Client, ctx context.Context, user *ent.User, course *ent.Course, data ReviewSchema) (*ent.Review, *config.ErrorResponse) {
	// Check if user has already reviewed this course
	existentReview, _ := db.Review.Query().Where(review.UserIDEQ(user.ID), review.CourseIDEQ(course.ID)).Only(ctx)
	if existentReview != nil {
		err := config.RequestErr(config.ERR_NOT_ALLOWED, "You have already reviewed this course")
		return nil, &err
	}

	// Check if user is enrolled for this course
	enrollmentObj := c.GetExistentEnrollmentByUserAndCourse(db, ctx, user, course, false)
	if enrollmentObj == nil || enrollmentObj.PaymentStatus != enrollment.PaymentStatusSuccessful {
		err := config.RequestErr(config.ERR_FORBIDDEN, "Only enrolled users can review this course")
		return nil, &err
	}

	review := db.Review.Create().
		SetUser(user).
		SetCourse(course).
		SetRating(data.Rating).
		SetComment(data.Comment).
		SaveX(ctx)
	review.Edges.User = user
	return review, nil
}

func (c CourseManager) UpdateReview(db *ent.Client, ctx context.Context, reviewObj *ent.Review, data ReviewSchema) *ent.Review {
	updatedReviewObj := reviewObj.Update().
		SetRating(data.Rating).
		SetComment(data.Comment).
		SaveX(ctx)
	updatedReviewObj.Edges.User = reviewObj.Edges.User
	return updatedReviewObj
}
