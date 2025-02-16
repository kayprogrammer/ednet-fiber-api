package schemas

import (
	"context"

	"github.com/gosimple/slug"
	"entgo.io/ent"
)

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