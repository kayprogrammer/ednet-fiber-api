package courses

import (
	"time"

	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/course"
	"github.com/kayprogrammer/ednet-fiber-api/modules/base"
)

type CategoryOrTagSchema struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
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
	c.Rating = course.Rating
	c.StudentsCount = course.StudentsCount
	c.LessonsCount = course.LessonsCount
	category := course.Edges.Category
	c.Category = CategoryOrTagSchema{Name: category.Name, Slug: category.Slug}
	c.CreatedAt = course.CreatedAt
	c.UpdatedAt = course.CreatedAt
	return c
}

type CoursesResponseSchema struct {
	base.ResponseSchema
	Data  config.PaginationResponse[CourseListSchema] `json:"data"` 
}

func (c CoursesResponseSchema) Assign(coursesData *config.PaginationResponse[*ent.Course]) CoursesResponseSchema {
	items := c.Data.Items
	for i, course := range coursesData.Items {
		items[i] = items[i].Assign(course)
	}
	c.Data.Items = items
	return c
}


// CourseDetailSchema - Full details of a course
type CourseDetailSchema struct {
	CourseListSchema
	IntroVideoURL  *string               `json:"intro_video_url,omitempty"`
	IsPublished    bool                  `json:"is_published"`
	QuizzesCount   int                   `json:"quizzes_count"`
	Duration       uint                   `json:"duration"` // in minutes
	EnrollmentType course.EnrollmentType `json:"enrollment_type"`
	Certification  bool                  `json:"certification"`
	ReviewsCount   int                   `json:"reviews_count"`
}

// Assign values from Course to CourseDetailSchema
func (c CourseDetailSchema) Assign(course *ent.Course) CourseDetailSchema {
	c.CourseListSchema = c.CourseListSchema.Assign(course)
	c.IntroVideoURL = &course.IntroVideoURL
	c.IsPublished = course.IsPublished
	c.QuizzesCount = course.QuizzesCount
	c.Duration = course.Duration
	c.EnrollmentType = course.EnrollmentType
	c.Certification = course.Certification
	c.ReviewsCount = course.ReviewsCount
	return c
}
