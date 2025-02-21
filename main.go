package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/modules/base/routes"
	"github.com/kayprogrammer/ednet-fiber-api/modules/seeding"
)

// @title EDNET API
// @version 1.0
// @description.markdown api
// @Accept json
// @Produce json
// @BasePath  /api/v1
// @Security BearerAuth
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type 'Bearer jwt_string' to correctly set the API Key
func main() {
	cfg := config.GetConfig()
	ctx := context.Background()
	db := config.ConnectDb(cfg, ctx)
	seeding.CreateInitialData(db, ctx, cfg)

	app := fiber.New()

	// CORS config
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowedOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Access-Control-Allow-Origin, Content-Disposition",
		AllowCredentials: true,
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	// Swagger Config
	swaggerCfg := swagger.Config{
		FilePath: "./docs/swagger.json",
		Path:     "/",
		Title:    "EDNET API Documentation",
		CacheAge: 1,
	}

	app.Use(swagger.New(swaggerCfg))

	routes.SetupRoutes(app, db)
	defer db.Close()
	log.Fatal(app.Listen(fmt.Sprintf(":%s", cfg.Port)))
}
