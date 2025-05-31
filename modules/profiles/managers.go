package profiles

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/course"
	"github.com/kayprogrammer/ednet-fiber-api/ent/enrollment"
	"github.com/kayprogrammer/ednet-fiber-api/ent/user"
	"github.com/kayprogrammer/ednet-fiber-api/modules/courses"
)

var courseManager = courses.CourseManager{}
// ----------------------------------
// PROFILES MANAGEMENT
// --------------------------------
type ProfileManager struct {
}

func (obj ProfileManager) GetById(db *ent.Client, ctx context.Context, id uuid.UUID) *ent.User {
	u, _ := db.User.
		Query().
		Where(user.ID(id)).
		Only(ctx)
	return u
}

func (p ProfileManager) Update(db *ent.Client, ctx context.Context, user *ent.User, data ProfileUpdateSchema, avatar *string) *ent.User {
	updatedUser := user.Update().
		SetName(data.Name).
		SetUsername(data.Username).
		SetNillableBio(data.Bio).
		SetNillableDob(config.ParseDate(data.Dob)).
		SetNillableAvatar(avatar).
		SaveX(ctx)
	return updatedUser
}

func (p ProfileManager) GetAllPaginatedEnrolledCourses(db *ent.Client, fibCtx *fiber.Ctx, user *ent.User) *config.PaginationResponse[*ent.Course] {
	query := db.Course.Query().
		Where(
			course.HasEnrollmentsWith(
				enrollment.UserID(user.ID),
				enrollment.PaymentStatusEQ(enrollment.PaymentStatusSuccessful),
			),
		).
		WithInstructor().
		WithCategory().
		WithTags()

	query = courseManager.ApplyCourseFilters(fibCtx, query)
	courses := config.PaginateModel(fibCtx, query)
	return courses
}
