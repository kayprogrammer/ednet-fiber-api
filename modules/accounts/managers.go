package accounts

import (
	"context"

	"github.com/google/uuid"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/user"
)

// ----------------------------------
// USER MANAGEMENT
// --------------------------------
type UserManager struct {
}

func (obj UserManager) GetById(db *ent.Client, ctx context.Context, id uuid.UUID) *ent.User {
	u, _ := db.User.
		Query().
		Where(user.ID(id)).
		Only(ctx)
	return u
}

func (obj UserManager) GetByRefreshToken(db *ent.Client, ctx context.Context, token string) *ent.User {
	u, _ := db.User.
		Query().
		Where(user.Refresh(token)).
		Only(ctx)
	return u
}

func (obj UserManager) GetByEmail(db *ent.Client, ctx context.Context, email string) *ent.User {
	u, _ := db.User.
		Query().
		Where(user.Email(email)).
		Only(ctx)
	return u
}

func (obj UserManager) GetByUsername(db *ent.Client, ctx context.Context, username string) *ent.User {
	u, _ := db.User.
		Query().
		Where(user.Username(username)).
		Only(ctx)
	return u
}

func (obj UserManager) GetByUsernames(db *ent.Client, ctx context.Context, usernames []string, excludeOpts ...uuid.UUID) []*ent.User {
	usersQ := db.User.
		Query().
		Where(user.UsernameIn(usernames...))
	if len(excludeOpts) > 0 {
		usersQ = usersQ.Where(user.IDNEQ(excludeOpts[0]))
	}
	users := usersQ.AllX(ctx)
	return users
}

func (obj UserManager) Create(db *ent.Client, ctx context.Context, userData RegisterSchema, isStaff bool, isVerified bool) *ent.User {
	password := config.HashPassword(userData.Password)
	u := db.User.Create().
		SetName(userData.Name).
		SetEmail(userData.Email).
		SetUsername(userData.Username).
		SetPassword(password).
		SetIsStaff(isStaff).
		SetIsVerified(isVerified).
		SaveX(ctx)
	return u
}

func (obj UserManager) GetOrCreate(db *ent.Client, ctx context.Context, userData RegisterSchema, isVerified bool, isStaff bool) *ent.User {
	user := obj.GetByEmail(db, ctx, userData.Email)
	if user == nil {
		// Create user
		user = obj.Create(db, ctx, userData, isStaff, isVerified)
	}
	return user
}

func (obj UserManager) UpdateTokens(ctx context.Context, user *ent.User, access string, refresh string) *ent.User {
	u := user.Update().SetAccess(access).SetRefresh(refresh).SaveX(ctx)
	return u
}

func (obj UserManager) DropData(db *ent.Client, ctx context.Context) {
	db.User.Delete().ExecX(ctx)
}