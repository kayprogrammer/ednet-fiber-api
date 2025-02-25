package courses

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/category"
	"github.com/kayprogrammer/ednet-fiber-api/ent/course"
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

func (c CourseManager) GetAll(db *ent.Client, fibCtx *fiber.Ctx) *config.PaginationResponse[*ent.Course] {
	courses := config.PaginateModel(
		fibCtx,
		db.Course.Query().
			WithInstructor().
			WithCategory().
			WithTags(),
	)
	return courses
}

func (c CourseManager) GetCourseByName(db *ent.Client, ctx context.Context, name string) *ent.Course {
	course, _ := db.Course.Query().Where(course.TitleEQ(name)).First(ctx)
	return course
}

func (c CourseManager) FilterCoursesByName(db *ent.Client, fibCtx *fiber.Ctx, name string) *config.PaginationResponse[*ent.Course] {
	courses := config.PaginateModel(
		fibCtx,
		db.Course.Query().
			Where(course.TitleContains(name)).
			WithInstructor().
			WithCategory().
			WithTags(),
	)
	return courses
}

func (c CourseManager) FilterFreeOrPaidCourses(db *ent.Client, fibCtx *fiber.Ctx, isFree bool) *config.PaginationResponse[*ent.Course] {
	courses := config.PaginateModel(
		fibCtx,
		db.Course.Query().
			Where(course.IsFreeEQ(isFree)).
			WithInstructor().
			WithCategory().
			WithTags(),
	)
	return courses
}

func (c CourseManager) FilterCoursesByInstructor(db *ent.Client, fibCtx *fiber.Ctx, instructor *ent.User) *config.PaginationResponse[*ent.Course] {
	courses := config.PaginateModel(
		fibCtx,
		db.Course.Query().
			Where(course.InstructorIDEQ(instructor.ID)).
			WithInstructor().
			WithCategory().
			WithTags(),
	)
	return courses
}

func (c CourseManager) GetCourseBySlug(db *ent.Client, ctx context.Context, slug string) *ent.Course {
	course, _ := db.Course.Query().Where(course.SlugEQ(slug)).Only(ctx)
	return course
}
