package accounts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/modules/base"
)

var userManager = UserManager{}

// @Summary Register a new user
// @Description This endpoint registers new users into our application.
// @Tags Auth
// @Param user body RegisterSchema true "User object"
// @Success 201 {object} RegisterResponseSchema
// @Failure 422 {object} config.ErrorResponse
// @Router /auth/register [post]
func Register(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		data := RegisterSchema{}
		// Validate request
		if errCode, errData := config.ValidateRequest(c, &data); errData != nil {
			return c.Status(*errCode).JSON(errData)
		}

		userByEmail := userManager.GetByEmail(db, ctx, data.Email)
		if userByEmail != nil {
			return config.APIError(c, 422, config.ValidationErr("email", "Email already registered!"))
		}
		userByUsername := userManager.GetByUsername(db, ctx, data.Username)
		if userByUsername != nil {
			return config.APIError(c, 422, config.ValidationErr("username", "Username already used!"))
		}

		// Create User
		newUser := userManager.Create(db, ctx, data, false, false)

		// Send Email
		go config.SendEmail(newUser, config.ET_ACTIVATE, newUser.Otp)

		response := RegisterResponseSchema{
			ResponseSchema: base.ResponseSchema{Message: "Registration successful"},
			Data:           EmailRequestSchema{Email: newUser.Email},
		}
		return c.Status(201).JSON(response)
	}
}
