package seeding

import (
	"context"

	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/course"
	"github.com/kayprogrammer/ednet-fiber-api/modules/instructors"
)

func createCategory(db *ent.Client, ctx context.Context) *ent.Category {
	name := "Programming"
	category_ := courseManager.GetCategoryByName(db, ctx, name)
	if category_ == nil {
		category_ = adminManager.CreateCategory(db, ctx, name) 
	}
	return category_
}

func createCourse(db *ent.Client, ctx context.Context, instructor *ent.User, category *ent.Category) *ent.Course {
	courseData := instructors.CourseCreateSchema{
		Title: "Test Course", Desc: "This is a test course", Language: "English",
		Difficulty: course.DifficultyBeginner, Duration: 240, IsFree: false, Price: 100,
		DiscountPrice: 80, EnrollmentType: course.EnrollmentTypeOpen, Certification: true,
	} 

	course_ := courseManager.GetCourseByName(db, ctx, courseData.Title)
	if course_ == nil {
		thumbnailUrl := "https://ednet-images.com/150"
		introVideoUrl := "https://ednet-videos.com/watch?v=introvideo"
		course_ = instructorManager.CreateCourse(db, ctx, instructor, category, thumbnailUrl, &introVideoUrl, courseData)
	}
	return course_
}