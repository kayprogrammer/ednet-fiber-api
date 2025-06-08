package instructors

import (
	"github.com/kayprogrammer/ednet-fiber-api/ent/course"
)

type CourseCreateSchema struct {
	Title          string                `form:"title" validate:"required,max=50,min=10"`
	Desc           string                `form:"desc" validate:"required,max=10000,min=10"`
	CategorySlug   string                `form:"category_slug" validate:"required"`
	Language       string                `form:"language" validate:"required" example:"English"`
	Difficulty     course.Difficulty     `form:"difficulty" validate:"required,difficulty_type_validator"`
	Duration       uint                  `form:"duration" validate:"required"`
	IsFree         bool                  `form:"is_free"`
	Price          float64               `form:"price" validate:"required"`
	DiscountPrice  float64               `form:"discount_price" validate:"required"`
	EnrollmentType course.EnrollmentType `form:"enrollment_type" validate:"required,enrollment_type_validator"`
	Certification  bool                  `form:"certification"`
}
