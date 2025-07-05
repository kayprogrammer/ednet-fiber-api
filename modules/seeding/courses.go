package seeding

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/kayprogrammer/ednet-fiber-api/ent"
)

func createCategories(db *ent.Client, ctx context.Context) []*ent.Category {
	categories := courseManager.GetCategories(db, ctx)
	if len(categories) < 1 {
		for _, categoryName := range CategoriesToCreate {
			category_ := courseManager.GetCategoryByName(db, ctx, categoryName)
			if category_ == nil {
				adminManager.CreateCategory(db, ctx, categoryName) 
			}
		}
		categories = courseManager.GetCategories(db, ctx)
	}
	return categories
}

func createCourses(db *ent.Client, ctx context.Context, instructor *ent.User, categories []*ent.Category) []*ent.Course {
	log.Println("Seeding Courses Data...")
	courses := courseManager.GetAll(db, ctx)
	if len(courses) < 1 {
		// Loop through each course in your Courses slice and create it
		for _, course := range CoursesToCreate {
			// Shuffle the list
			rand.New(rand.NewSource(time.Now().UnixNano()))
			rand.Shuffle(len(categories), func(i, j int) {
				categories[i], categories[j] = categories[j], categories[i]
			})
			// Create the Course record
			courseRecord, err := db.Course.Create().
				SetInstructor(instructor).
				SetTitle(course.Title).
				SetSlug(course.Slug).
				SetCategory(categories[0]).
				SetDesc(course.Desc).
				SetThumbnailURL(course.ThumbnailUrl).
				SetIntroVideoURL(course.IntroVideoUrl).
				SetDuration(course.Duration).
				SetIsFree(course.IsFree).
				SetPrice(course.Price).
				SetDiscountPrice(course.DiscountPrice).
				SetIsPublished(rand.Intn(2) == 1).
				Save(ctx)

			if err != nil {
				log.Fatalf("Failed to create course '%s': %v", course.Title, err)
			}

			// Create Lessons for this Course
			for _, lesson := range course.Lessons {
				lessonRecord, err := db.Lesson.Create().
					SetCourseID(courseRecord.ID).
					SetTitle(lesson.Title).
					SetSlug(lesson.Slug).
					SetDesc(lesson.Desc).
					SetThumbnailURL(lesson.ThumbnailUrl).
					SetVideoURL(lesson.VideoUrl).
					SetContent(lesson.Content).
					SetOrder(lesson.Order).
					SetDuration(lesson.Duration).
					SetIsPublished(lesson.IsPublished).
					SetIsFreePreview(lesson.IsFreePreview).
					Save(ctx)

				if err != nil {
					log.Fatalf("Failed to create lesson '%s': %v", lesson.Title, err)
				}

				// If the lesson has a quiz, create it
				if lesson.Quiz != nil {
					quizRecord, err := db.Quiz.Create().
						SetLesson(lessonRecord).
						SetTitle(lesson.Quiz.Title).
						SetSlug(lesson.Quiz.Slug).
						SetDescription(lesson.Quiz.Description).
						SetDuration(lesson.Quiz.Duration).
						SetIsPublished(true).
						Save(ctx)

					if err != nil {
						log.Fatalf("Failed to create quiz for lesson '%s': %v", lesson.Title, err)
					}

					// Create Questions for this Quiz
					for _, questionData := range lesson.Quiz.Questions {
						questionRecord, err := db.Question.Create().
							SetQuiz(quizRecord).
							SetText(questionData.Text).
							SetOrder(questionData.Order).
							Save(ctx)

						if err != nil {
							log.Fatalf("Failed to create question '%s': %v", questionData.Text, err)
						}

						// Create Options for this Question
						for _, optionData := range questionData.Options {
							_, err := db.QuestionOption.Create().
								SetQuestion(questionRecord).
								SetText(optionData.Text).
								SetIsCorrect(optionData.IsCorrect).
								Save(ctx)

							if err != nil {
								log.Fatalf("Failed to create option for question '%s': %v", questionData.Text, err)
							}
						}
					}
				}
			}
		}
		courses = courseManager.GetAll(db, ctx)
	}
	log.Println("Courses Data Seeded Successfully.")
	return courses
}