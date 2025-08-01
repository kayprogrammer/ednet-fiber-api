package schemas

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return append(
		CommonFields,
		field.String("name").NotEmpty().MaxLen(200).MinLen(2),
		field.String("username").NotEmpty().Unique().MaxLen(100).MinLen(2),
		field.String("email").NotEmpty(),
		field.String("password").MinLen(8).NotEmpty(),
		field.Bool("is_verified").Default(false),
		field.Bool("is_active").Default(true),
		field.String("bio").Optional().Nillable().MaxLen(300).MinLen(10),
		field.Time("dob").Optional().Nillable(),
		field.String("avatar").Optional().Nillable(),
		field.Uint32("otp").Optional().Nillable(),
		field.Time("otp_expiry").Optional().Nillable(),
		field.Bool("social_login").Default(false),
		field.Time("last_summary_date").Optional().Nillable(),
		field.Int("summary_count").Default(0),
		field.Enum("role").Values("student", "instructor", "admin").Default("student"),
	)
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tokens", Token.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("courses", Course.Type).Annotations(entsql.OnDelete(entsql.SetNull)),
		edge.To("enrollments", Enrollment.Type).Annotations(entsql.OnDelete(entsql.SetNull)),
		edge.To("reviews", Review.Type).Annotations(entsql.OnDelete(entsql.SetNull)),
		edge.To("payments", Payment.Type).Annotations(entsql.OnDelete(entsql.SetNull)),
		edge.To("quiz_results", QuizResult.Type).Annotations(entsql.OnDelete(entsql.SetNull)),
        edge.To("progress", LessonProgress.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

// Token holds authentication tokens for users.
type Token struct {
	ent.Schema
}

// Fields of the Token.
func (Token) Fields() []ent.Field {
	return append(
		CommonFields,
		field.UUID("user_id", uuid.UUID{}),
		field.String("access").NotEmpty(),
		field.String("refresh").NotEmpty(),
	)
}

// Edges of the Token.
func (Token) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("tokens").Field("user_id").Unique().Required(),
	}
}
