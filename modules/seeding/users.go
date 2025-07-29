package seeding

import (
	"context"

	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/user"
)

func createAdmin(db *ent.Client, ctx context.Context, cfg config.Config) *ent.User {
	email := cfg.FirstAdminEmail
	username := "testadmin"
	user_, _ := db.User.Query().Where(user.Or(user.Email(email), user.Username(username))).Only(ctx)
	if user_ == nil {
		user_ = db.User.Create().SetName("Test Admin").SetEmail(email).SetUsername(username).
			SetRole(user.RoleAdmin).SetIsVerified(true).SetPassword(config.HashPassword(cfg.FirstAdminPassword)).
			SaveX(ctx)
	}
	return user_
}

func createStudent(db *ent.Client, ctx context.Context, cfg config.Config) *ent.User {
	email := cfg.FirstStudentEmail
	username := "teststudent"
	user_, _ := db.User.Query().Where(user.Or(user.Email(email), user.Username(username))).Only(ctx)
	if user_ == nil {
		user_ = db.User.Create().SetName("Test Student").SetEmail(email).SetUsername(username).
			SetRole(user.RoleStudent).SetIsVerified(true).SetPassword(config.HashPassword(cfg.FirstStudentPassword)).
			SaveX(ctx)
	}
	return user_
}

func createInstructor(db *ent.Client, ctx context.Context, cfg config.Config) *ent.User {
	email := cfg.FirstInstructorEmail
	username := "testinstructor"
	user_, _ := db.User.Query().Where(user.Or(user.Email(email), user.Username(username))).Only(ctx)
	if user_ == nil {
		user_ = db.User.Create().SetName("Test Instructor").SetEmail(email).SetUsername(username).
			SetRole(user.RoleInstructor).SetIsVerified(true).SetPassword(config.HashPassword(cfg.FirstInstructorPassword)).
			SaveX(ctx)
	}
	return user_
}
