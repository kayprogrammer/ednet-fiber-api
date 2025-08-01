package profiles

import (
	"context"
	"math"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/course"
	"github.com/kayprogrammer/ednet-fiber-api/ent/enrollment"
	"github.com/kayprogrammer/ednet-fiber-api/ent/lesson"
	"github.com/kayprogrammer/ednet-fiber-api/ent/lessonprogress"
	"github.com/kayprogrammer/ednet-fiber-api/ent/predicate"
	"github.com/kayprogrammer/ednet-fiber-api/ent/quizresult"
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

func (p ProfileManager) GetAllPaginatedEnrolledCourses(db *ent.Client, fibCtx *fiber.Ctx, user *ent.User, status string) *config.PaginationResponse[*ent.Course] {
	enrollmentPredicates := []predicate.Enrollment{
		enrollment.UserID(user.ID),
		enrollment.PaymentStatusEQ(enrollment.PaymentStatusSuccessful),
	}

	if status != "" {
		enrollmentPredicates = append(enrollmentPredicates, enrollment.StatusEQ(enrollment.Status(status)))
	}

	query := db.Course.Query().
		Where(
			course.HasEnrollmentsWith(enrollmentPredicates...),
		).
		WithInstructor().
		WithCategory().
		WithTags()

	query = courseManager.ApplyCourseFilters(fibCtx, query)
	courses := config.PaginateModel(fibCtx, query)
	return courses
}

// ----------------------------------
// LESSON PROGRESS MANAGEMENT
// --------------------------------

func (p ProfileManager) CreateOrUpdateLessonProgress(db *ent.Client, ctx context.Context, user *ent.User, lesson *ent.Lesson, isCompleted bool) (*ent.LessonProgress, string) {
	lessonProgress := p.GetLessonProgress(db, ctx, user, lesson.ID)
	message := "created"
	if lessonProgress == nil {
		// Create lesson progress
		lessonProgressCreateQ := db.LessonProgress.Create().
			SetUserID(user.ID).
			SetLessonID(lesson.ID)
		if isCompleted {
			lessonProgressCreateQ = lessonProgressCreateQ.SetCompletedAt(time.Now())
		}
		lessonProgress = lessonProgressCreateQ.SaveX(ctx)
	} else {
		if isCompleted {
			lessonProgress.Update().SetCompletedAt(time.Now()).SaveX(ctx)
		}
		message = "updated"
	}
	return lessonProgress, message
}

func (p ProfileManager) GetLessonProgress(db *ent.Client, ctx context.Context, user *ent.User, lessonId uuid.UUID) *ent.LessonProgress {
	lessonProgress, _ := db.LessonProgress.Query().
		Where(
			lessonprogress.UserID(user.ID),
			lessonprogress.LessonID(lessonId),
		).Only(ctx)
	return lessonProgress
}

func (p ProfileManager) GetCourseProgress(
	db *ent.Client,
	ctx context.Context,
	user *ent.User,
	courseObj *ent.Course,
) float64 {
	lessonProgressCount := db.LessonProgress.Query().
		Where(
			lessonprogress.UserID(user.ID),
			lessonprogress.CompletedAtNotNil(),
			lessonprogress.HasLessonWith(lesson.CourseIDEQ(courseObj.ID)),
		).CountX(ctx)

	lessonsCount := courseObj.QueryLessons().CountX(ctx)

	if lessonsCount == 0 {
		return 0
	}
	percentage := (float64(lessonProgressCount) / float64(lessonsCount)) * 100
	rounded := math.Round(percentage*100) / 100
	return rounded
}

func (p ProfileManager) GetLeaderboard(db *ent.Client, ctx context.Context) []*LeaderboardEntry {
	leaderboard := []*LeaderboardEntry{}

	db.QuizResult.Query().
		GroupBy(quizresult.UserColumn).
		Aggregate(ent.Sum(quizresult.FieldScore)).
		ScanX(ctx, &leaderboard)

	// Sort manually since Ent doesn't allow ordering on groupBy aggregate directly
	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].TotalScore > leaderboard[j].TotalScore
	})

	// Fetch user info
	userIDs := make([]uuid.UUID, len(leaderboard))
	for i, entry := range leaderboard {
		userIDs[i] = entry.UserID
	}

	users := db.User.Query().Where(user.IDIn(userIDs...)).AllX(ctx)

	userMap := make(map[uuid.UUID]*ent.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	for _, entry := range leaderboard {
		if u, ok := userMap[entry.UserID]; ok {
			entry.Name = u.Name
			entry.Username = u.Username
		}
	}
	return leaderboard
}