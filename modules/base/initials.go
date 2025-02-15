package base

import (
	"context"
	"log"

	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/user"
)

func createAdmin(db *ent.Client, ctx context.Context, cfg config.Config) *ent.User {
	email := cfg.FirstAdminEmail
	user_, _ := db.User.Query().Where(user.Email(email)).Only(ctx)
	if user_ == nil {
		user_ = db.User.Create().SetName("Test Admin").SetEmail(email).SetUsername("testadmin").
			SetRole(user.RoleAdmin).SetIsVerified(true).SetPassword(config.HashPassword(cfg.FirstAdminPassword)).
			SaveX(ctx)
	}
	return user_
}

func createStudent(db *ent.Client, ctx context.Context, cfg config.Config) *ent.User {
	email := cfg.FirstStudentEmail
	user_, _ := db.User.Query().Where(user.Email(email)).Only(ctx)
	if user_ == nil {
		user_ = db.User.Create().SetName("Test Student").SetEmail(email).SetUsername("teststudent").
			SetRole(user.RoleStudent).SetIsVerified(true).SetPassword(config.HashPassword(cfg.FirstStudentPassword)).
			SaveX(ctx)
	}
	return user_
}

func createInstructor(db *ent.Client, ctx context.Context, cfg config.Config) *ent.User {
	email := cfg.FirstInstructorEmail
	user_, _ := db.User.Query().Where(user.Email(email)).Only(ctx)
	if user_ == nil {
		user_ = db.User.Create().SetName("Test Instructor").SetEmail(email).SetUsername("testinstructor").
			SetRole(user.RoleInstructor).SetIsVerified(true).SetPassword(config.HashPassword(cfg.FirstInstructorPassword)).
			SaveX(ctx)
	}
	return user_
}

func CreateInitialData(db *ent.Client, ctx context.Context, cfg config.Config) {
	log.Println("Creating Initial Data....")
	createAdmin(db, ctx, cfg)
	createStudent(db, ctx, cfg)
	createInstructor(db, ctx, cfg)
	log.Println("Initial Data Created....")
}
