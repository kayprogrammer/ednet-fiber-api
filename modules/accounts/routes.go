package accounts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/modules/base"
)

var userManager = UserManager{}

// @Summary Register a new user
// @Description `This endpoint registers new users into our application.`
// @Tags Auth
// @Param user body RegisterSchema true "User object"
// @Success 201 {object} RegisterResponseSchema
// @Failure 422 {object} base.ValidationErrorExample
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

// @Summary Verify a user's email
// @Description `This endpoint verifies a user's email.`
// @Tags Auth
// @Param email_data body VerifyEmailRequestSchema true "Email object"
// @Success 200 {object} base.ResponseSchema
// @Failure 422 {object} base.ValidationErrorExample
// @Failure 404 {object} base.NotFoundErrorExample
// @Failure 400 {object} base.InvalidErrorExample
// @Router /auth/verify-email [post]
func VerifyEmail(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		data := VerifyEmailRequestSchema{}
		// Validate request
		if errCode, errData := config.ValidateRequest(c, &data); errData != nil {
			return c.Status(*errCode).JSON(errData)
		}
		user := userManager.GetByEmail(db, ctx, data.Email)
		if user == nil {
			return config.APIError(c, 404, config.RequestErr(config.ERR_INCORRECT_EMAIL, "Incorrect Email"))
		}

		if user.IsVerified {
			return c.Status(200).JSON(base.ResponseMessage("Email already verified"))
		}
		if user.Otp == nil || *user.Otp != data.Otp {
			return config.APIError(c, 404, config.RequestErr(config.ERR_INCORRECT_OTP, "Incorrect Otp"))
		}

		if userManager.IsOtpExpired(user) {
			return config.APIError(c, 400, config.RequestErr(config.ERR_EXPIRED_OTP, "Expired Otp"))
		}

		// Update User
		user.Update().SetIsVerified(true).Save(ctx)

		// Send Welcome Email
		go config.SendEmail(user, config.ET_WELCOME, nil)
		return c.Status(200).JSON(base.ResponseMessage("Account verification successful"))
	}
}

// @Summary Resend Verification Email
// @Description `This endpoint resends new otp to the user's email.`
// @Tags Auth
// @Param email body EmailRequestSchema true "Email object"
// @Success 200 {object} base.ResponseSchema
// @Failure 422 {object} base.ValidationErrorExample
// @Failure 404 {object} base.NotFoundErrorExample
// @Router /auth/resend-verification-email [post]
func ResendVerificationEmail(db *ent.Client) fiber.Handler {
	return func (c *fiber.Ctx) error {
		ctx := c.Context()
		data := EmailRequestSchema{}

		// Validate request
		if errCode, errData := config.ValidateRequest(c, &data); errData != nil {
			return c.Status(*errCode).JSON(errData)
		}

		user := userManager.GetByEmail(db, ctx, data.Email)
		if user == nil {
			return config.APIError(c, 404, config.RequestErr(config.ERR_INCORRECT_EMAIL, "Incorrect Email"))
		}

		if user.IsVerified {
			return c.Status(200).JSON(base.ResponseMessage("Email already verified"))
		}

		// Send Email
		otp, otpExp := userManager.GetOtp()
		user.Update().SetOtp(otp).SetOtpExpiry(otpExp).Save(ctx)
		go config.SendEmail(user, config.ET_ACTIVATE, &otp)
		return c.Status(200).JSON(base.ResponseMessage("Verification email sent"))
	}
}
