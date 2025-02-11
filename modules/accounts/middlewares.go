package accounts

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
)

func GetUser(db *ent.Client, ctx context.Context, token string) (*ent.User, *string) {
	if !strings.HasPrefix(token, "Bearer ") {
		err := "Auth Bearer Not Provided!"
		return nil, &err
	}
	user, err := DecodeAccessToken(db, ctx, token[7:])
	if err != nil {
		return nil, err
	}
	return user, nil
}

func AuthMiddleware(db *ent.Client) fiber.Handler {
	return func (c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if len(token) < 1 {
			return config.APIError(c, 401, config.RequestErr(config.ERR_UNAUTHORIZED_USER, "Unauthorized User!"))
		}
		user, err := GetUser(db, c.Context(), token)
		if err != nil {
			return config.APIError(c, 401, config.RequestErr(config.ERR_INVALID_TOKEN, *err))
		}
		c.Locals("user", user)
		return c.Next()
	}
}