package courses

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
)

type CourseManager struct{}

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
