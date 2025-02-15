package schemas

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Payment schema.
type Payment struct {
	ent.Schema
}

// Fields of Payment.
func (Payment) Fields() []ent.Field {
	return append(
		CommonFields,
		field.UUID("user_id", uuid.UUID{}),
		field.UUID("course_id", uuid.UUID{}),
		field.Float("amount"),
		field.Enum("status").Values("pending", "successful", "failed").Default("pending"),
		field.String("payment_method"),
		field.String("transaction_id").Unique(),
	)
}

// Edges of Payment.
func (Payment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("payments").Field("user_id").Unique().Required(),
		edge.From("course", Course.Type).Ref("payments").Field("course_id").Unique().Required(),
	}
}
