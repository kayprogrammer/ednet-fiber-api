package instructors

import (
	"context"

	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/course"
)

type InstructorManager struct{}

func (i InstructorManager) CreateCourse(db *ent.Client, ctx context.Context, instructor *ent.User, category *ent.Category, thumbnailUrl string, introVideoUrl *string, data CourseCreateSchema) *ent.Course {
	slug := i.GenerateCourseSlug(db, ctx, data.Title)
	course := db.Course.Create().SetTitle(data.Title).SetSlug(slug).SetDesc(data.Desc).
		SetInstructor(instructor).SetCategoryID(category.ID).SetLanguage(data.Language).
		SetDifficulty(data.Difficulty).SetDuration(data.Duration).SetIsFree(data.IsFree).
		SetThumbnailURL(thumbnailUrl).SetNillableIntroVideoURL(introVideoUrl).
		SetPrice(data.Price).SetDiscountPrice(data.DiscountPrice).SetEnrollmentType(data.EnrollmentType).
		SetCertification(data.Certification).SaveX(ctx)
	return course
}

func (i InstructorManager) GenerateCourseSlug(db *ent.Client, ctx context.Context, title string) string {
	baseSlug := config.Slugify(title)
	uniqueSlug := baseSlug
	for {
		exists, _ := db.Course.Query().Where(course.SlugEQ(uniqueSlug)).Exist(ctx)
		if !exists {
			break
		}
		uniqueSlug = baseSlug + "-" + config.GetRandomString(7)
	}
	return uniqueSlug
}