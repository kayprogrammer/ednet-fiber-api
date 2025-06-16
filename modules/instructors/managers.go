package instructors

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/course"
	"github.com/kayprogrammer/ednet-fiber-api/ent/enrollment"
	"github.com/kayprogrammer/ednet-fiber-api/ent/lesson"
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

func (i InstructorManager) UpdateCourse(db *ent.Client, ctx context.Context, course *ent.Course, category *ent.Category, thumbnailUrl *string, introVideoUrl *string, data CourseCreateSchema) *ent.Course {
	slug := course.Slug
	if data.Title != course.Title {
		slug = i.GenerateCourseSlug(db, ctx, data.Title)
	}
	updatedCourseQuery := course.Update().SetTitle(data.Title).SetSlug(slug).SetDesc(data.Desc).
		SetCategoryID(category.ID).SetLanguage(data.Language).
		SetDifficulty(data.Difficulty).SetDuration(data.Duration).SetIsFree(data.IsFree).
		SetPrice(data.Price).SetDiscountPrice(data.DiscountPrice).SetEnrollmentType(data.EnrollmentType).
		SetCertification(data.Certification)
	if thumbnailUrl != nil {
		updatedCourseQuery = updatedCourseQuery.SetThumbnailURL(*thumbnailUrl)
	}
	if introVideoUrl != nil {
		updatedCourseQuery = updatedCourseQuery.SetIntroVideoURL(*introVideoUrl)
	}
	updatedCourse := updatedCourseQuery.SaveX(ctx)

	// Edges reassignment to prevent reload
	updatedCourse.Edges.Instructor = course.Edges.Instructor
	updatedCourse.Edges.Category = category
	updatedCourse.Edges.Reviews = course.Edges.Reviews
	updatedCourse.Edges.Quizzes = course.Edges.Quizzes
	updatedCourse.Edges.Enrollments = course.Edges.Enrollments
	updatedCourse.Edges.Lessons = course.Edges.Lessons
	return updatedCourse
}

func (i InstructorManager) DeleteCourse(db *ent.Client, ctx context.Context, courseObj *ent.Course) *string {
	// Prevent deletion if there's a paid enrollment
	enrollmentExists := db.Enrollment.Query().Where(enrollment.CourseIDEQ(courseObj.ID), enrollment.PaymentStatusEQ(enrollment.PaymentStatusSuccessful)).ExistX(ctx)
	if enrollmentExists {
		errMsg := "Cannot delete a course that has at least one paid enrollment"
		return &errMsg
	}
	db.Course.DeleteOne(courseObj).ExecX(ctx)
	return nil
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

func (i InstructorManager) GetCourseLessonBySlug(db *ent.Client, ctx context.Context, instructor *ent.User, slug string, loaded bool) *ent.Lesson {
	query := db.Lesson.Query().
		Where(
			lesson.SlugEQ(slug),
			lesson.HasCourseWith(course.InstructorIDEQ(instructor.ID)),
		)
	if loaded {
		query = query.
			WithCourse()
	}
	lessonObj, _ := query.Only(ctx)
	return lessonObj
}

func (i InstructorManager) GenerateLessonSlug(db *ent.Client, ctx context.Context, title string) string {
	baseSlug := config.Slugify(title)
	uniqueSlug := baseSlug
	for {
		exists, _ := db.Lesson.Query().Where(lesson.SlugEQ(uniqueSlug)).Exist(ctx)
		if !exists {
			break
		}
		uniqueSlug = baseSlug + "-" + config.GetRandomString(7)
	}
	return uniqueSlug
}

func (i InstructorManager) CreateLesson(db *ent.Client, ctx context.Context, course *ent.Course, thumbnailUrl string, videoUrl *string, data LessonCreateSchema) *ent.Lesson {
	slug := i.GenerateLessonSlug(db, ctx, data.Title)
	lessonObj := db.Lesson.Create().SetTitle(data.Title).SetSlug(slug).SetDesc(data.Desc).
		SetCourse(course).SetNillableContent(data.Content).SetOrder(data.Order).
		SetIsPublished(data.IsPublished).SetDuration(data.Duration).SetIsFreePreview(data.IsFreePreview).
		SetThumbnailURL(thumbnailUrl).SetNillableVideoURL(videoUrl).
		SaveX(ctx)
	return lessonObj
}

func (i InstructorManager) UpdateLesson(db *ent.Client, ctx context.Context, lesson *ent.Lesson, thumbnailUrl *string, videoUrl *string, data LessonCreateSchema) *ent.Lesson {
	slug := lesson.Slug
	if data.Title != lesson.Title {
		slug = i.GenerateLessonSlug(db, ctx, data.Title)
	}

	updateLessonQuery := lesson.Update().SetTitle(data.Title).SetSlug(slug).SetDesc(data.Desc).
		SetOrder(data.Order).SetIsPublished(data.IsPublished).SetDuration(data.Duration).SetIsFreePreview(data.IsFreePreview)

	if thumbnailUrl != nil {
		updateLessonQuery = updateLessonQuery.SetNillableThumbnailURL(thumbnailUrl)
	}
	if videoUrl != nil {
		updateLessonQuery = updateLessonQuery.SetNillableVideoURL(videoUrl)
	}
	if data.Content != nil {
		updateLessonQuery = updateLessonQuery.SetNillableContent(data.Content)
	}
	lessonObj := updateLessonQuery.SaveX(ctx)
	return lessonObj
}

func (i InstructorManager) DeleteLesson(db *ent.Client, ctx context.Context, lessonObj *ent.Lesson) *string {
	// Prevent deletion if there's a paid enrollment for a published lesson
	enrollmentExists := db.Enrollment.Query().Where(enrollment.CourseIDEQ(lessonObj.CourseID), enrollment.PaymentStatusEQ(enrollment.PaymentStatusSuccessful)).ExistX(ctx)
	if enrollmentExists && lessonObj.IsPublished {
		errMsg := "Cannot delete a published lesson which has at least one paid enrollment"
		return &errMsg
	}
	db.Lesson.DeleteOne(lessonObj).ExecX(ctx)
	return nil
}