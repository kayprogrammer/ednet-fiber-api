package admin

import (
	"context"

	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
)

type AdminManager struct{}

func (a AdminManager) CreateCategory (db *ent.Client, ctx context.Context, name string) *ent.Category {
	category := db.Category.Create().SetName(name).SetSlug(config.Slugify(name)).SaveX(ctx)
	return category
}