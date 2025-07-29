package seeding

import (
	"context"
	"log"

	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
)

func CreateInitialData(db *ent.Client, ctx context.Context, cfg config.Config) {
	log.Println("Creating Initial Data....")
	admin := createAdmin(db, ctx, cfg)
	student := createStudent(db, ctx, cfg)
	instructor := createInstructor(db, ctx, cfg)
	users := []*ent.User{admin, student, instructor}
	categories := createCategories(db, ctx)
	courses := createCourses(db, ctx, instructor, categories)
	createReviews(db, ctx, users, courses)
	log.Println("Initial Data Created")
}
