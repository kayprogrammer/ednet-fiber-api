package courses

import (
	"context"
	"fmt"

	"github.com/gosimple/slug"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/course"
)

// GenerateCourseSlug creates a unique slug with a 7-character random suffix.
func GenerateCourseSlug(db *ent.Client, ctx context.Context, name string) string {
	baseSlug := slug.Make(name)
	courseSlug := baseSlug

	const maxAttempts = 5
	for i := 0; i < maxAttempts; i++ {
		// Check if slug exists
		course, _ := db.Course.Query().Where(course.Slug(courseSlug)).Only(ctx)
		if course == nil {
			break
		}
		// Generate a new slug with random 7-character suffix
		courseSlug = fmt.Sprintf("%s-%s", baseSlug, config.GetRandomString(7))
	}
	return courseSlug
}