package base

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
)

func RequestUser(c *fiber.Ctx) *ent.User {
	user := c.Locals("user")
	if user != nil {
		return user.(*ent.User)
	}
	return nil
}
