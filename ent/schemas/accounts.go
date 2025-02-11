package schemas

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
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
		field.Uint32("otp").
			Optional().Nillable(),
		field.Time("otp_expiry").
			Optional().Nillable(),
		field.Bool("social_login").
			Default(false),
	)
}

// Edges of the User.
func (User) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("tokens", Token.Type),
    }
}

// Token holds authentication tokens for users.
type Token struct {
    ent.Schema
}

// Fields of the Token.
func (Token) Fields() []ent.Field {
    return []ent.Field{
        field.String("access").NotEmpty(),
        field.String("refresh").NotEmpty(),
    }
}

// Edges of the Token.
func (Token) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("user", User.Type).
            Ref("tokens").
            Unique().
            Required(),
    }
}