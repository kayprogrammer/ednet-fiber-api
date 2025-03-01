package courses

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/category"
	"github.com/kayprogrammer/ednet-fiber-api/ent/course"
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
		"title": func(value string) { query.Where(course.TitleContains(value)) },
		"instructor": func(value string) { query.Where(course.HasInstructorWith(user.NameContains(value))) },
		"instructorUsername": func(value string) { query.Where(course.HasInstructorWith(user.UsernameContains(value))) },
		"exactInstructorUsername": func(value string) { query.Where(course.HasInstructorWith(user.UsernameEQ(value))) },
		"isFree": func(value string) {
			if freeStatus, err := strconv.ParseBool(value); err == nil {
				query.Where(course.IsFreeEQ(freeStatus))
			}
		},
	}
	// Apply filters dynamically
	for param, apply := range filters {
		if value := fibCtx.Query(param); value != "" {
			apply(value)
		}
	}

	// Sorting by rating
	switch fibCtx.Query("sortByRating") {
	case "desc":
		query.Order(ent.Desc(course.FieldRating))
	case "asc":
		query.Order(ent.Asc(course.FieldRating))
	}
	return query
}

func (c CourseManager) GetAll(db *ent.Client, fibCtx *fiber.Ctx) *config.PaginationResponse[*ent.Course] {
	query := db.Course.Query().
		WithInstructor().
		WithCategory().
		WithTags()

	query = c.ApplyCourseFilters(fibCtx, query)
	courses := config.PaginateModel(fibCtx, query)
	return courses
}

func (c CourseManager) GetCourseByName(db *ent.Client, ctx context.Context, name string) *ent.Course {
	course, _ := db.Course.Query().Where(course.TitleEQ(name)).First(ctx)
	return course
}

func (c CourseManager) FilterCoursesByInstructor(db *ent.Client, fibCtx *fiber.Ctx, instructor *ent.User) *config.PaginationResponse[*ent.Course] {
	query := db.Course.Query().
		Where(course.InstructorIDEQ(instructor.ID)).
		WithInstructor().
		WithCategory().
		WithTags()
	query = c.ApplyCourseFilters(fibCtx, query)
	courses := config.PaginateModel(fibCtx, query)
	return courses
}

func (c CourseManager) GetCourseBySlug(db *ent.Client, ctx context.Context, slug string) *ent.Course {
	course, _ := db.Course.Query().Where(course.SlugEQ(slug)).Only(ctx)
	return course
}
