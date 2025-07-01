package courses

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/category"
	"github.com/kayprogrammer/ednet-fiber-api/ent/course"
	"github.com/kayprogrammer/ednet-fiber-api/ent/enrollment"
	"github.com/kayprogrammer/ednet-fiber-api/ent/lesson"
	"github.com/kayprogrammer/ednet-fiber-api/ent/quiz"
	"github.com/kayprogrammer/ednet-fiber-api/ent/review"
	"github.com/kayprogrammer/ednet-fiber-api/ent/user"
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
			WithQuizzes().
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

func (c CourseManager) GetQuizzes(db *ent.Client, course *ent.Course, fibCtx *fiber.Ctx) *config.PaginationResponse[*ent.Quiz] {
	query := db.Quiz.Query().Where(quiz.CourseID(course.ID)).Order(ent.Asc(quiz.FieldTitle)).WithQuestions()
	query = c.ApplyQuizFilters(fibCtx, query)
	quizzes := config.PaginateModel(fibCtx, query)
	return quizzes
}

func (c CourseManager) GetQuizBySlug(db *ent.Client, ctx context.Context, slug string, loaded bool) *ent.Quiz {
	query := db.Quiz.Query().
		Where(quiz.SlugEQ(slug))
	if loaded {
		query = query.
			WithCourse().
			WithQuestions(
				func(q *ent.QuestionQuery) {
					q.WithOptions()
				},
			)
	}
	quiz, _ := query.Only(ctx)
	return quiz
}

func (c CourseManager) GetCourseLessonBySlug(db *ent.Client, ctx context.Context, slug string, loaded bool) *ent.Lesson {
	query := db.Lesson.Query().
		Where(lesson.SlugEQ(slug))
	if loaded {
		query = query.
			WithCourse()
	}
	lesson, _ := query.Only(ctx)
	return lesson
}

func (c CourseManager) GetExistentEnrollmentByUserAndCourse (db *ent.Client, ctx context.Context, user *ent.User, course *ent.Course, loaded bool) *ent.Enrollment {
	query := db.Enrollment.Query().
		Where(
			enrollment.UserID(user.ID),
			enrollment.CourseID(course.ID),
		)
	if loaded {
		query = query.
			WithCourse().
			WithUser()
	}
	enrollmentObj, _ := query.Only(ctx)
	return enrollmentObj
}

func (c CourseManager) CreateEnrollment (db *ent.Client, ctx context.Context, user *ent.User, course *ent.Course, checkoutUrl string) (*ent.Enrollment, *config.ErrorResponse) {
	existentEnrollment := c.GetExistentEnrollmentByUserAndCourse(db, ctx, user, course, false)
	if existentEnrollment != nil {
		err := config.RequestErr(config.ERR_NOT_ALLOWED, "Enrollment has been created already")
		return nil, &err
	}
	enrollmentQuery := db.Enrollment.
		Create().
		SetCourse(course).
		SetUser(user).
		SetCheckoutURL(checkoutUrl)
	
		if course.IsFree {
			enrollmentQuery = enrollmentQuery.SetStatus(enrollment.StatusActive).
			SetPaymentStatus(enrollment.PaymentStatusSuccessful)
		}
	enrollmentObj := enrollmentQuery.SaveX(ctx)
	enrollmentObj.Edges.User = user
	enrollmentObj.Edges.Course = course
	return enrollmentObj, nil
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