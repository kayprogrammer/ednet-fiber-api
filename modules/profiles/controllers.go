package profiles

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/modules/base"
)

// @Summary Get Your Profile
// @Description `This endpoint allows a user to view his/her profile`
// @Tags Profiles
// @Success 200 {object} ProfileResponseSchema
// @Failure 401 {object} base.UnauthorizedErrorExample
// @Router /profiles [get]
// @Security BearerAuth
func GetProfile(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := base.RequestUser(c)
		response := ProfileResponseSchema{
			ResponseSchema: base.ResponseMessage("Profiles fetched"),
			Data: ProfileSchema{}.Assign(user),
		}
		return c.Status(200).JSON(response)
	}
}
