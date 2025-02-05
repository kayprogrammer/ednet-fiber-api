package schemas

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return append(
		CommonFields,
		field.String("name").
			NotEmpty(),
		field.String("username").
			NotEmpty().Unique(),
		field.String("email").
			NotEmpty(),
		field.String("password").
			NotEmpty(),
		field.Bool("is_verified").
			Default(false),
		field.Bool("is_staff").
			Default(false),
		field.Bool("is_active").
			Default(true),
		field.String("bio").
			Optional().Nillable(),
		field.Time("dob").
			Optional().Nillable(),
		field.String("avatar").
			Optional().Nillable(),
		field.String("access").
			Optional().Nillable(),
		field.String("refresh").
			Optional().Nillable(),
		field.Int("otp").
			Optional().Nillable(),
		field.Time("otpExpiry").
			Optional().Nillable(),
	)
}
