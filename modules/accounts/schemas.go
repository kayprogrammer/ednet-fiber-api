package accounts

import "github.com/kayprogrammer/ednet-fiber-api/modules/base"

// REQUEST BODY SCHEMAS
type RegisterSchema struct {
	Name     string `json:"name" validate:"required,max=50" example:"John Doe"`
	Username string `json:"username" validate:"required,max=50" example:"johndoe"`
	Email    string `json:"email" validate:"required,min=5,email" example:"johndoe@example.com"`
	Password string `json:"password" validate:"required,min=8,max=50" example:"strongpassword"`
}

type EmailRequestSchema struct {
	Email string `json:"email" validate:"required,min=5,email" example:"johndoe@email.com"`
}

type VerifyEmailRequestSchema struct {
	EmailRequestSchema
	Otp uint32 `json:"otp" validate:"required" example:"123456"`
}

type SetNewPasswordSchema struct {
	VerifyEmailRequestSchema
	Password string `json:"password" validate:"required,min=8,max=50" example:"newstrongpassword"`
}

type LoginSchema struct {
	EmailOrUsername string `json:"email_or_username" validate:"required" example:"johndoe"`
	Password        string `json:"password" validate:"required" example:"password"`
}

type TokenSchema struct {
	Token string `json:"token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InNpbXBsZWlkIiwiZXhwIjoxMjU3ODk0MzAwfQ.Ys_jP70xdxch32hFECfJQuvpvU5_IiTIN2pJJv68EqQ"`
}

// RESPONSE BODY SCHEMAS
type RegisterResponseSchema struct {
	base.ResponseSchema
	Data EmailRequestSchema `json:"data"`
}
type TokensResponseSchema struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type LoginResponseSchema struct {
	base.ResponseSchema
	Data TokensResponseSchema `json:"data"`
}
