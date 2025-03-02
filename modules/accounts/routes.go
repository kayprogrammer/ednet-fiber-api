package accounts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/user"
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
			return config.APIError(c, *errCode, *errData)
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
		newUser := userManager.Create(db, ctx, data, user.RoleStudent, false)

		// Send Email
		go config.SendEmail(newUser, config.ET_ACTIVATE, newUser.Otp)

		response := RegisterResponseSchema{
			ResponseSchema: base.ResponseMessage("Registration successful"),
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
			return config.APIError(c, *errCode, *errData)
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
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		data := EmailRequestSchema{}

		// Validate request
		if errCode, errData := config.ValidateRequest(c, &data); errData != nil {
			return config.APIError(c, *errCode, *errData)
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

// @Summary Send Password Reset Otp
// @Description `This endpoint sends new password reset otp to the user's email.`
// @Tags Auth
// @Param email body EmailRequestSchema true "Email object"
// @Success 200 {object} base.ResponseSchema
// @Failure 422 {object} base.ValidationErrorExample
// @Failure 404 {object} base.NotFoundErrorExample
// @Router /auth/send-password-reset-otp [post]
func SendPasswordResetOtp(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		data := EmailRequestSchema{}

		// Validate request
		if errCode, errData := config.ValidateRequest(c, &data); errData != nil {
			return config.APIError(c, *errCode, *errData)
		}

		user := userManager.GetByEmail(db, ctx, data.Email)
		if user == nil {
			return config.APIError(c, 404, config.RequestErr(config.ERR_INCORRECT_EMAIL, "Incorrect Email"))
		}

		// Send Email
		otp, otpExp := userManager.GetOtp()
		user.Update().SetOtp(otp).SetOtpExpiry(otpExp).Save(ctx)
		go config.SendEmail(user, config.ET_RESET, &otp)
		return c.Status(200).JSON(base.ResponseMessage("Password otp sent"))
	}
}

// @Summary Set New Password
// @Description `This endpoint verifies the password reset otp.`
// @Tags Auth
// @Param email body SetNewPasswordSchema true "Password reset object"
// @Success 200 {object} base.ResponseSchema
// @Failure 422 {object} base.ValidationErrorExample
// @Failure 404 {object} base.NotFoundErrorExample
// @Failure 400 {object} base.InvalidErrorExample
// @Router /auth/set-new-password [post]
func SetNewPassword(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		data := SetNewPasswordSchema{}

		// Validate request
		if errCode, errData := config.ValidateRequest(c, &data); errData != nil {
			return config.APIError(c, *errCode, *errData)
		}
		user := userManager.GetByEmail(db, ctx, data.Email)
		if user == nil {
			return config.APIError(c, 404, config.RequestErr(config.ERR_INCORRECT_EMAIL, "Incorrect Email"))
		}

		if user.Otp == nil || *user.Otp != data.Otp {
			return config.APIError(c, 404, config.RequestErr(config.ERR_INCORRECT_OTP, "Incorrect Otp"))
		}

		if userManager.IsOtpExpired(user) {
			return config.APIError(c, 400, config.RequestErr(config.ERR_EXPIRED_OTP, "Expired Otp"))
		}

		// Set Password
		user.Update().SetPassword(config.HashPassword(data.Password)).Save(ctx)

		// Send Email
		go config.SendEmail(user, config.ET_RESET_SUCC, nil)
		return c.Status(200).JSON(base.ResponseMessage("Password reset successful"))
	}
}

// @Summary Login a user
// @Description `This endpoint generates new access and refresh tokens for authentication`
// @Tags Auth
// @Param user body LoginSchema true "User login"
// @Success 201 {object} LoginResponseSchema
// @Failure 422 {object} base.ValidationErrorExample
// @Failure 401 {object} base.UnauthorizedErrorExample
// @Router /auth/login [post]
func Login(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		data := LoginSchema{}

		// Validate request
		if errCode, errData := config.ValidateRequest(c, &data); errData != nil {
			return config.APIError(c, *errCode, *errData)
		}

		user := userManager.GetByEmailOrUsername(db, ctx, data.EmailOrUsername)
		if user == nil {
			return config.APIError(c, 401, config.RequestErr(config.ERR_INVALID_CREDENTIALS, "Invalid Credentials"))
		}
		if !user.IsVerified {
			return config.APIError(c, 401, config.RequestErr(config.ERR_UNVERIFIED_USER, "Verify your email first"))
		}

		// Create Auth Tokens
		access := GenerateAccessToken(user.ID, user.Username)
		refresh := GenerateRefreshToken()
		userManager.AddTokens(db, ctx, user, access, refresh)

		response := LoginResponseSchema{
			ResponseSchema: base.ResponseMessage("Login successful"),
			Data:           TokensResponseSchema{Access: access, Refresh: refresh},
		}
		return c.Status(201).JSON(response)
	}
}

// @Summary Login a user via google
// @Description `This endpoint generates new access and refresh tokens for authentication via google`
// @Description `Pass in token gotten from gsi client authentication here in payload to retrieve tokens for authorization`
// @Tags Auth
// @Param user body TokenSchema true "Google auth"
// @Success 201 {object} LoginResponseSchema
// @Failure 422 {object} base.ValidationErrorExample
// @Failure 401 {object} base.UnauthorizedErrorExample
// @Router /auth/google [post]
func GoogleLogin(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		data := TokenSchema{}

		// Validate request
		if errCode, errData := config.ValidateRequest(c, &data); errData != nil {
			return c.Status(*errCode).JSON(errData)
		}

		userGoogleData, errData := ConvertGoogleToken(ctx, data.Token)
		if errData != nil {
			return config.APIError(c, 401, *errData)
		}

		email := userGoogleData.Email
		name := userGoogleData.Name
		avatar := userGoogleData.Picture

		access, refresh, err := RegisterSocialUser(db, ctx, email, name, &avatar)
		if err != nil {
			return c.Status(401).JSON(err)
		}
		response := LoginResponseSchema{
			ResponseSchema: base.ResponseMessage("Login successful"),
			Data:           TokensResponseSchema{Access: *access, Refresh: *refresh},
		}
		return c.Status(201).JSON(response)
	}
}

// @Summary Refresh tokens
// @Description `This endpoint refresh tokens by generating new access and refresh tokens for a user`
// @Tags Auth
// @Param refresh body TokenSchema true "Refresh token"
// @Success 201 {object} LoginResponseSchema
// @Failure 422 {object} base.ValidationErrorExample
// @Failure 401 {object} base.UnauthorizedErrorExample
// @Router /auth/refresh [post]
func Refresh(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		data := TokenSchema{}

		// Validate request
		if errCode, errData := config.ValidateRequest(c, &data); errData != nil {
			return config.APIError(c, *errCode, *errData)
		}

		token := data.Token
		user := DecodeRefreshToken(db, ctx, token)
		if user == nil {
			return config.APIError(c, 401, config.RequestErr(config.ERR_INVALID_TOKEN, "Refresh token is invalid or expired"))
		}

		// Create and Update Auth Tokens
		access := GenerateAccessToken(user.ID, user.Username)
		refresh := GenerateRefreshToken()
		userManager.UpdateTokens(db, ctx, access, refresh, token)

		response := LoginResponseSchema{
			ResponseSchema: base.ResponseMessage("Tokens refresh successful"),
			Data:           TokensResponseSchema{Access: access, Refresh: refresh},
		}
		return c.Status(201).JSON(response)
	}
}

// @Summary Logout a user
// @Description `This endpoint logs a user out from our application from a single device`
// @Tags Auth
// @Success 200 {object} base.ResponseSchema
// @Failure 401 {object} base.UnauthorizedErrorExample
// @Router /auth/logout [get]
// @Security BearerAuth
func Logout(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userManager.DeleteToken(db, c.Context(), c.Get("Authorization")[7:])
		return c.Status(200).JSON(base.ResponseMessage("Logout successful"))
	}
}

// @Summary Logout a user from all devices
// @Description `This endpoint logs a user out from our application from all devices`
// @Tags Auth
// @Success 200 {object} base.ResponseSchema
// @Failure 401 {object} base.UnauthorizedErrorExample
// @Router /auth/logout/all [get]
// @Security BearerAuth
func LogoutAll(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := base.RequestUser(c)
		userManager.ClearTokens(db, c.Context(), user.ID)
		return c.Status(200).JSON(base.ResponseMessage("Logout successful"))
	}
}
