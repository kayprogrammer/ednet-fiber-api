package instructors

import (
	"mime/multipart"

	"github.com/kayprogrammer/ednet-fiber-api/ent/course"
)

type CourseCreateSchema struct {
	Title          string                `form:"title"`
	Desc           string                `form:"desc"`
	ThumbnailUrl   multipart.FileHeader  `form:"thumbnail_url"`
	IntroVideoUrl  multipart.FileHeader  `form:"intro_video_url"`
	CategorySlug   string                `form:"category_slug"`
	Language       string                `form:"language"`
	Difficulty     course.Difficulty     `form:"difficulty"`
	Duration       uint                  `form:"duration"`
	IsFree         bool                  `form:"is_free"`
	Price          float64               `form:"price"`
	DiscountPrice  float64               `form:"discount_price"`
	EnrollmentType course.EnrollmentType `form:"enrollment_type"`
	Certification  bool                  `form:"certification"`
}
