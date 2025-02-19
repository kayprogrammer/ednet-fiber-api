package profiles

import (
	"context"
	"github.com/google/uuid"
	"github.com/kayprogrammer/ednet-fiber-api/config"
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

func (obj ProfileManager) Update(db *ent.Client, ctx context.Context, user *ent.User, data ProfileUpdateSchema, avatar *string) *ent.User {
	updatedUser := user.Update().
		SetName(data.Name).
		SetUsername(data.Username).
		SetNillableBio(data.Bio).
		SetNillableDob(config.ParseDate(data.Dob)).
		SetNillableAvatar(avatar).
		SaveX(ctx)
	return updatedUser
}