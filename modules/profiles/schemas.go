package profiles

import (
	"time"

	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/user"
	"github.com/kayprogrammer/ednet-fiber-api/modules/base"
)

type ProfileSchema struct {
	Name     string     `json:"name" example:"John Doe"`
	Username string     `json:"username" example:"johndoe"`
	Email    string     `json:"email" example:"johndoe@example.com"`
	Bio      *string    `json:"bio" example:"I'm the boss"`
	Dob      *time.Time `json:"dob" example:"2000-09-12"`
	Avatar   *string    `json:"avatar" example:"https://ednet-images.com/users/john-doe"`
	Role     user.Role  `json:"role" example:"student"`
}

func (p ProfileSchema) Assign(u *ent.User) ProfileSchema {
	p.Name = u.Name
	p.Username = u.Username
	p.Email = u.Email
	p.Bio = u.Bio
	p.Dob = u.Dob
	p.Avatar = u.Avatar
	p.Role = u.Role
	return p
}

type ProfileResponseSchema struct {
	base.ResponseSchema
	Data ProfileSchema `json:"data"`
}

type ProfileUpdateSchema struct {
	Name     string     `json:"name" validate:"required,max=150,min=10" example:"John Doe"`
	Username string     `json:"username" validate:"required,max=50,min=2" example:"john-doe"`
	Bio      *string    `json:"bio" validate:"max=300,min=10" example:"I'm the boss"`
	Dob      *time.Time `json:"dob" example:"2000-09-12"`
}
