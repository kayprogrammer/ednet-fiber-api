package courses

import (
	"time"

	"github.com/google/uuid"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/course"
	"github.com/kayprogrammer/ednet-fiber-api/ent/enrollment"
	"github.com/kayprogrammer/ednet-fiber-api/modules/base"
)

type CategoryOrTagSchema struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (c CategoryOrTagSchema) Assign(category *ent.Category, tag *ent.Tag) CategoryOrTagSchema {
	if category != nil {
		c.Name = category.Name
		c.Slug = category.Slug
	} else {
		c.Name = tag.Name
		c.Slug = tag.Slug
	}
	return c
}

// CourseListSchema - Summary of a course for listings
type CourseListSchema struct {
	Instructor    base.UserDataSchema `json:"instructor"`
	Title         string              `json:"title" example:"Go Programming for Beginners"`
	Slug          string              `json:"slug" example:"go-programming-for-beginners"`
	Desc          string              `json:"desc"`
	ThumbnailURL  string              `json:"thumbnail_url" example:"https://ednet-images.com/courses/go.jpg"`
	Language      string              `json:"language" example:"English"`
	Difficulty    course.Difficulty   `json:"difficulty" example:"Beginner"`
	DiscountPrice *float64            `json:"discount_price,omitempty"`
	Price         float64             `json:"price" example:"19.99"`
	IsFree        bool                `json:"is_free" example:"false"`
	IsPublished   bool                `json:"is_published" example:"false"`
	Rating        float64             `json:"rating" example:"4.8"`
	StudentsCount int                 `json:"students_count" example:"1200"`
	LessonsCount  int                 `json:"lessons_count" example:"20"`
	Category      CategoryOrTagSchema `json:"category"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
}

// Assign values from Course to CourseListSchema
func (c CourseListSchema) Assign(course *ent.Course) CourseListSchema {
	c.Instructor = c.Instructor.Assign(course.Edges.Instructor)
	c.Title = course.Title
	c.Slug = course.Slug
	c.Desc = course.Desc
	c.ThumbnailURL = course.ThumbnailURL
	c.Language = course.Language
	c.Difficulty = course.Difficulty
	c.DiscountPrice = &course.DiscountPrice
	c.Price = course.Price
	c.IsFree = course.IsFree
	c.IsPublished = course.IsPublished
	c.Rating = courseManager.GetAverageRating(course.Edges.Reviews)
	c.StudentsCount = len(course.Edges.Enrollments)
	c.LessonsCount = len(course.Edges.Lessons)
	c.Category = c.Category.Assign(course.Edges.Category, nil)
	c.CreatedAt = course.CreatedAt
	c.UpdatedAt = course.CreatedAt
	return c
}

type CoursesResponseSchema struct {
	base.ResponseSchema
	Data config.PaginationResponse[CourseListSchema] `json:"data"`
}

func (c CoursesResponseSchema) Assign(coursesData *config.PaginationResponse[*ent.Course]) CoursesResponseSchema {
	items := make([]CourseListSchema, 0)
	for _, course := range coursesData.Items {
		items = append(items, CourseListSchema{}.Assign(course))
	}
	c.Data.Items = items
	c.Data.ItemsCount = coursesData.ItemsCount
	c.Data.Page = coursesData.Page
	c.Data.TotalPages = coursesData.TotalPages
	c.Data.Limit = coursesData.Limit
	return c
}

// CourseDetailSchema - Full details of a course
type CourseDetailSchema struct {
	CourseListSchema
	IntroVideoURL  *string               `json:"intro_video_url,omitempty"`
	Duration       uint                  `json:"duration"` // in minutes
	EnrollmentType course.EnrollmentType `json:"enrollment_type"`
	Certification  bool                  `json:"certification"`
	ReviewsCount   int                   `json:"reviews_count"`
}

// Assign values from Course to CourseDetailSchema
func (c CourseDetailSchema) Assign(course *ent.Course) CourseDetailSchema {
	c.CourseListSchema = c.CourseListSchema.Assign(course)
	c.IntroVideoURL = &course.IntroVideoURL
	c.Duration = course.Duration
	c.EnrollmentType = course.EnrollmentType
	c.Certification = course.Certification
	c.ReviewsCount = len(course.Edges.Reviews)
	return c
}

type CourseResponseSchema struct {
	base.ResponseSchema
	Data CourseDetailSchema `json:"data"`
}

type LessonListSchema struct {
	Title         string `json:"title"`
	Slug          string `json:"slug"`
	Desc          string `json:"desc"`
	Order         uint   `json:"order"`
	Duration      uint   `json:"duration"`
	IsPublished   bool   `json:"is_published"`
	IsFreePreview bool   `json:"is_free_preview"`
	ThumbnailURL  string `json:"thumbnail_url" example:"https://ednet-images.com/lessons/go.jpg"`
}

// Assign values from Lesson to LessonSchema
func (l LessonListSchema) Assign(lesson *ent.Lesson) LessonListSchema {
	l.Title = lesson.Title
	l.Slug = lesson.Slug
	l.Desc = lesson.Desc
	l.Order = lesson.Order
	l.Duration = lesson.Duration
	l.IsPublished = lesson.IsPublished
	l.IsFreePreview = lesson.IsFreePreview
	l.ThumbnailURL = lesson.ThumbnailURL
	return l
}

type LessonsResponseSchema struct {
	base.ResponseSchema
	Data config.PaginationResponse[LessonListSchema] `json:"data"`
}

func (l LessonsResponseSchema) Assign(lessonsData *config.PaginationResponse[*ent.Lesson]) LessonsResponseSchema {
	items := make([]LessonListSchema, 0)
	for _, lesson := range lessonsData.Items {
		items = append(items, LessonListSchema{}.Assign(lesson))
	}
	l.Data.Items = items
	l.Data.ItemsCount = lessonsData.ItemsCount
	l.Data.Page = lessonsData.Page
	l.Data.TotalPages = lessonsData.TotalPages
	l.Data.Limit = lessonsData.Limit
	return l
}

type LessonDetailSchema struct {
	LessonListSchema
	QuizzesCount int    `json:"quizzes_count"`
	VideoUrl     string `json:"video_url"`
	Content      string `json:"content"`
}

// Assign values from Lesson to LessonDetailSchema
func (l LessonDetailSchema) Assign(lesson *ent.Lesson) LessonDetailSchema {
	l.LessonListSchema = l.LessonListSchema.Assign(lesson)
	l.QuizzesCount = len(lesson.Edges.Quizzes)
	l.VideoUrl = lesson.VideoURL
	l.Content = lesson.Content
	return l
}

type LessonResponseSchema struct {
	base.ResponseSchema
	Data LessonDetailSchema `json:"data"`
}

type QuizzesResponseSchema struct {
	base.ResponseSchema
	Data config.PaginationResponse[QuizListSchema] `json:"data"`
}

func (q QuizzesResponseSchema) Assign(quizzesData *config.PaginationResponse[*ent.Quiz]) QuizzesResponseSchema {
	items := make([]QuizListSchema, 0)
	for _, quiz := range quizzesData.Items {
		items = append(items, QuizListSchema{}.Assign(quiz))
	}
	q.Data.Items = items
	q.Data.ItemsCount = quizzesData.ItemsCount
	q.Data.Page = quizzesData.Page
	q.Data.TotalPages = quizzesData.TotalPages
	q.Data.Limit = quizzesData.Limit
	return q
}

type QuizResponseSchema struct {
	base.ResponseSchema
	Data QuizDetailSchema `json:"data"`
}

type QuizListSchema struct {
	ID             uuid.UUID `json:"id"`
	Title          string    `json:"title"`
	Slug           string    `json:"slug"`
	Description    string    `json:"description"`
	TotalQuestions int       `json:"total_questions"`
	Duration       int       `json:"duration"`
	IsPublished    bool      `json:"is_published"`
}

// Assign values from Quix to QuizListSchema
func (q QuizListSchema) Assign(quiz *ent.Quiz) QuizListSchema {
	q.ID = quiz.ID
	q.Title = quiz.Title
	q.Slug = quiz.Slug
	q.Description = quiz.Description
	q.TotalQuestions = len(quiz.Edges.Questions)
	q.Duration = quiz.Duration
	q.IsPublished = quiz.IsPublished
	return q
}

type QuizDetailSchema struct {
	QuizListSchema
	Questions []QuestionSchema `json:"questions"`
}

// Assign values from Quiz to QuizDetailSchema
func (q QuizDetailSchema) Assign(quiz *ent.Quiz) QuizDetailSchema {
	q.QuizListSchema = q.QuizListSchema.Assign(quiz)
	questions := quiz.Edges.Questions
	parsedQuestions := make([]QuestionSchema, 0)
	for _, question := range questions {
		parsedOptions := make([]QuestionOptionSchema, 0)
		for _, option := range question.Edges.Options {
			parsedOptions = append(parsedOptions, QuestionOptionSchema{Text: option.Text, IsCorrect: option.IsCorrect})
		}
		parsedQuestions = append(parsedQuestions, QuestionSchema{
			Text:    question.Text,
			Order:   question.Order,
			Options: parsedOptions,
		})
	}
	q.Questions = parsedQuestions
	return q
}

type QuestionSchema struct {
	ID      uuid.UUID              `json:"id"`
	Order   int                    `json:"order"`
	Text    string                 `json:"text"`
	Options []QuestionOptionSchema `json:"options"`
}

type QuestionOptionSchema struct {
	ID        uuid.UUID `json:"id"`
	Text      string    `json:"text"`
	IsCorrect bool      `json:"is_correct"`
}

type EnrollForACourseSchema struct {
	SuccessUrl string `json:"success_url" validate:"required,url" example:"https://domain-example.com/payment-success"`
	CancelUrl  string `json:"cancel_url" validate:"required,url" example:"https://domain-example.com/payment-cancelled"`
}

type EnrollmentSchema struct {
	User          base.UserDataSchema      `json:"user"`
	Course        CourseListSchema         `json:"course"`
	Status        enrollment.Status        `json:"status"`
	PaymentStatus enrollment.PaymentStatus `json:"payment_status"`
	CheckoutURL   string                   `json:"checkout_url"`
	Progress      int                      `json:"progress"`
}

func (e EnrollmentSchema) Assign(enrollmentObj *ent.Enrollment) EnrollmentSchema {
	e.User = e.User.Assign(enrollmentObj.Edges.User)
	e.Course = e.Course.Assign(enrollmentObj.Edges.Course)
	e.Status = enrollmentObj.Status
	e.PaymentStatus = enrollmentObj.PaymentStatus
	e.CheckoutURL = enrollmentObj.CheckoutURL
	e.Progress = enrollmentObj.Progress
	return e
}

type EnrollmentResponseSchema struct {
	base.ResponseSchema
	Data EnrollmentSchema `json:"data"`
}

type SubmitAnswerSchema struct {
	QuestionID       uuid.UUID `json:"question_id" validate:"required"`
	SelectedOptionID uuid.UUID `json:"selected_option_id" validate:"required"`
}

type QuizSubmissionSchema struct {
	Answers   []SubmitAnswerSchema `json:"answers" validate:"required,min=1"`
	TimeTaken int                  `json:"time_taken" validate:"required,min=1"` // in seconds
}

type AnswerSchema struct {
	QuestionID       uuid.UUID `json:"question_id" validate:"required"`
	SelectedOptionID uuid.UUID `json:"selected_option_id" validate:"required"`
}

type QuizResultSchema struct {
	ID          uuid.UUID      `json:"id"`
	Score       float64        `json:"score"`
	TimeTaken   int            `json:"time_taken"`
	StartedAt   time.Time      `json:"started_at"`
	CompletedAt time.Time      `json:"completed_at"`
	Answers     []AnswerSchema `json:"answers"`
}

func (q QuizResultSchema) Assign(quizResult *ent.QuizResult) QuizResultSchema {
	q.ID = quizResult.ID
	q.Score = quizResult.Score
	q.TimeTaken = quizResult.TimeTaken
	q.StartedAt = quizResult.StartedAt
	q.CompletedAt = quizResult.CompletedAt
	answers := make([]AnswerSchema, 0)
	for _, answer := range quizResult.Edges.Answers {
		answers = append(answers, AnswerSchema{
			QuestionID:       answer.QuestionID,
			SelectedOptionID: answer.SelectedOptionID,
		})
	}
	q.Answers = answers
	return q
}

type QuizResultResponseSchema struct {
	base.ResponseSchema
	Data QuizResultSchema `json:"data"`
}
