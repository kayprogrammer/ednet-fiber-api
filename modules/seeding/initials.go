package seeding

import (
	"context"
	"log"

	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
)

func CreateInitialData(db *ent.Client, ctx context.Context, cfg config.Config) {
	log.Println("Creating Initial Data....")
	createAdmin(db, ctx, cfg)
	createStudent(db, ctx, cfg)
	instructor := createInstructor(db, ctx, cfg)
	categories := createCategories(db, ctx)
	createCourses(db, ctx, instructor, categories)
	log.Println("Initial Data Created")
}
