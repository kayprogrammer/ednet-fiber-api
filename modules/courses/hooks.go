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

// Generates a slug from the name field.
func GenerateCategoryOrTagSlug() ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			// Ensure it's a create operation
			if !m.Op().Is(ent.OpCreate) {
				return next.Mutate(ctx, m)
			}

			// Get the name field from mutation
			name, _ := m.Field("name")

			// Generate the slug
			slugified := slug.Make(name.(string))

			// Set the slug field
			if err := m.SetField("slug", slugified); err != nil {
				return nil, err
			}

			return next.Mutate(ctx, m)
		})
	}
}