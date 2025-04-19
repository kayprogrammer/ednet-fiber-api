package courses

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/category"
	"github.com/kayprogrammer/ednet-fiber-api/ent/course"
	"github.com/kayprogrammer/ednet-fiber-api/ent/lesson"
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
		"instructor": func(value string) {
			query.Where(course.HasInstructorWith(user.Or(user.NameContains(value), user.UsernameContains(value))))
		},
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

func (c CourseManager) GetCourseBySlug(db *ent.Client, ctx context.Context, slug string, loaded bool) *ent.Course {
	query := db.Course.Query().
		Where(course.SlugEQ(slug))
	if loaded {
		query = query.
			WithInstructor().
			WithCategory().
			WithTags()
	}
	course, _ := query.Only(ctx)
	return course
}

func (c CourseManager) ApplyLessonFilters(fibCtx *fiber.Ctx, query *ent.LessonQuery) *ent.LessonQuery {
	filters := map[string]func(string){
		"title": func(value string) { query.Where(lesson.TitleContains(value)) },
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
	query := db.Lesson.Query().Where(lesson.CourseID(course.ID))
	query = c.ApplyLessonFilters(fibCtx, query)
	lessons := config.PaginateModel(fibCtx, query)
	return lessons
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