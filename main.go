package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/modules/base/routes"
	"github.com/kayprogrammer/ednet-fiber-api/modules/seeding"
)

// Custom panic recovery middleware
func RecoveryMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// Log the panic
				log.Printf("[PANIC] %v\n", r)
				// Return a safe JSON response
				_ = config.APIError(c, 500, config.ServerErr("Something went wrong!"))
			}
		}()
		return c.Next()
	}
}

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

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 15MB
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err == fiber.ErrRequestEntityTooLarge {
				return config.APIError(c, fiber.StatusRequestEntityTooLarge, config.ServerErr("Request body is too large. Maximum allowed size is 10MB."))
			}
			// Default error handler
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return config.APIError(c, code, config.ServerErr(err.Error()))
		},
	})

	// Register custom panic recovery middleware
	app.Use(RecoveryMiddleware())

	app.Get("/metrics", monitor.New())

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

	routes.SetupRoutes(app, db, cfg)
	app.Use(func(c *fiber.Ctx) error {
		return config.APIError(c, 404, config.NotFoundErr("Path not found"))
	})
	defer db.Close()
	log.Fatal(app.Listen(fmt.Sprintf(":%s", cfg.Port)))
}
