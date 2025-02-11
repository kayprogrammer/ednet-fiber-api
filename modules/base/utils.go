package base

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
)

func RequestUser(c *fiber.Ctx) *ent.User {
	return c.Locals("user").(*ent.User)
}