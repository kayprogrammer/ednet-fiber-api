package profiles

import (
	"context"

	"github.com/google/uuid"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/user"
)

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