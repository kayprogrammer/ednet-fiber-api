package instructors

import (
	"context"
	"github.com/gofiber/fiber/v2"
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

	// Edges reassignment to prevent reload
	course.Edges.Instructor = instructor
	course.Edges.Category = category
	course.Edges.Reviews = []*ent.Review{}
	course.Edges.Quizzes = []*ent.Quiz{}
	course.Edges.Enrollments = []*ent.Enrollment{}
	course.Edges.Lessons = []*ent.Lesson{}
	return course
}

func (i InstructorManager) UpdateCourse(db *ent.Client, ctx context.Context, course *ent.Course, category *ent.Category, thumbnailUrl string, introVideoUrl *string, data CourseCreateSchema) *ent.Course {
	slug := course.Slug
	if data.Title != course.Title {
		slug = i.GenerateCourseSlug(db, ctx, data.Title)
	}
	updatedCourse := course.Update().SetTitle(data.Title).SetSlug(slug).SetDesc(data.Desc).
		SetCategoryID(category.ID).SetLanguage(data.Language).
		SetDifficulty(data.Difficulty).SetDuration(data.Duration).SetIsFree(data.IsFree).
		SetThumbnailURL(thumbnailUrl).SetNillableIntroVideoURL(introVideoUrl).
		SetPrice(data.Price).SetDiscountPrice(data.DiscountPrice).SetEnrollmentType(data.EnrollmentType).
		SetCertification(data.Certification).SaveX(ctx)
	return updatedCourse
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

func (i InstructorManager) GetCoursesPaginated(db *ent.Client, fibCtx *fiber.Ctx, instructor *ent.User) *config.PaginationResponse[*ent.Course] {
	query := db.Course.Query().
		Where(course.InstructorIDEQ(instructor.ID)).
		WithInstructor().
		WithCategory().
		WithTags().
		WithReviews().
		WithEnrollments().
		WithLessons()
	query = courseManager.ApplyCourseFilters(fibCtx, query)
	courses := config.PaginateModel(fibCtx, query)
	return courses
}