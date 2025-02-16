package accounts

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/token"
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

func (obj UserManager) GetByRefreshToken(db *ent.Client, ctx context.Context, tokenStr string) *ent.User {
	u, _ := db.User.
		Query().
		Where(user.HasTokensWith(token.Refresh(tokenStr))).
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

func (obj UserManager) GetByEmailOrUsername(db *ent.Client, ctx context.Context, emailOrUsername string) *ent.User {
	u, _ := db.User.
		Query().
		Where(
			user.Or(
				user.Email(emailOrUsername), 
				user.Username(emailOrUsername),
			),
		).
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

func (obj UserManager) Create(db *ent.Client, ctx context.Context, userData RegisterSchema, role user.Role, isVerified bool) *ent.User {
	password := config.HashPassword(userData.Password)
	otp, otpExp := obj.GetOtp()

	u := db.User.Create().
		SetName(userData.Name).
		SetEmail(userData.Email).
		SetUsername(userData.Username).
		SetPassword(password).
		SetRole(role).
		SetIsVerified(isVerified).
		SetOtp(otp).
		SetOtpExpiry(otpExp).
		SaveX(ctx)
	return u
}

func (obj UserManager) GetOtp () (uint32, time.Time) {
	cfg := config.GetConfig()
	otp := config.GetRandomInt(6)
	otpExpiry := time.Now().UTC().Add(time.Duration(cfg.EmailOtpExpireMinutes) * time.Minute)
	return otp, otpExpiry
}

func (obj UserManager) IsOtpExpired (user *ent.User) bool {
	return time.Now().UTC().After(user.OtpExpiry.UTC())
}

func (obj UserManager) GetOrCreate(db *ent.Client, ctx context.Context, userData RegisterSchema, isVerified bool, role user.Role) *ent.User {
	user := obj.GetByEmail(db, ctx, userData.Email)
	if user == nil {
		// Create user
		user = obj.Create(db, ctx, userData, role, isVerified)
	}
	return user
}

func (obj UserManager) AddTokens(db *ent.Client, ctx context.Context, user *ent.User, access string, refresh string) {
	log.Println(db == nil, user.ID, access, refresh)
	db.Token.
        Create().
		SetUserID(user.ID).
        SetAccess(access).
        SetRefresh(refresh).
        SaveX(ctx)
}

func (obj UserManager) UpdateTokens(db *ent.Client, ctx context.Context, access string, refresh string, oldRefresh string) {
	db.Token.
        Update().
        Where(token.Refresh(oldRefresh)).
        SetAccess(access).
        SetRefresh(refresh).
        SaveX(ctx)
}

func (obj UserManager) DeleteToken(db *ent.Client, ctx context.Context, access string) {
	db.Token.
        Delete().
        Where(token.Access(access)).
		ExecX(ctx)
}

func (obj UserManager) ClearTokens(db *ent.Client, ctx context.Context, userID uuid.UUID) {
	db.Token.
        Delete().
        Where(token.HasUserWith(user.ID(userID))).
		ExecX(ctx)
}

func (obj UserManager) DropData(db *ent.Client, ctx context.Context) {
	db.User.Delete().ExecX(ctx)
}