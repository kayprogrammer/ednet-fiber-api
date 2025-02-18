package hooks

import (
	"context"

	"github.com/gosimple/slug"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
)

// Generates a slug from the name field.
func GenerateCategoryOrTagSlug(next ent.Mutator) ent.Mutator {
	type SlugSetter interface {
        SetSlug(value string)
    }
	return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
		// Get the name field from mutation
		name, _ := m.Field("name")

		// Generate the slug
		slugified := slug.Make(name.(string))

		if s, ok := m.(SlugSetter); ok {
			s.SetSlug(slugified)
		}
		return next.Mutate(ctx, m)
		
	})
}