package general

import (
	"context"

	"github.com/kayprogrammer/ednet-fiber-api/ent"
)

type SiteDetailService struct {
}

func (obj SiteDetailService) Get(db *ent.Client, ctx context.Context) *ent.SiteDetail {
	s, _ := db.SiteDetail.
		Query().
		First(ctx)
	return s
}

func (obj SiteDetailService) Create(client *ent.Client, ctx context.Context) *ent.SiteDetail {
	s := client.SiteDetail.
		Create().
		SaveX(ctx)
	return s
}

func (obj SiteDetailService) GetOrCreate(db *ent.Client, ctx context.Context) *ent.SiteDetail {
	sitedetail := obj.Get(db, ctx)
	if sitedetail == nil {
		sitedetail = obj.Create(db, ctx)
	}
	return sitedetail
}
