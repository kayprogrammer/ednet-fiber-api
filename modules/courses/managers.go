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

func (c CourseManager) GetCategoryByName(db *ent.Client, ctx context.Context, name string) *ent.Category {
	category, _ := db.Category.Query().Where(category.NameEQ(name)).Only(ctx)
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
