package schemas

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

var CommonFields = []ent.Field{
	field.UUID("id", uuid.UUID{}).
		Default(uuid.New),
	field.Time("created_at").
		Default(time.Now),
	field.Time("updated_at").
		Default(time.Now).
		UpdateDefault(time.Now),
}

type SiteDetail struct {
	ent.Schema
}

// Fields of the SiteDetail.
func (SiteDetail) Fields() []ent.Field {
	return append(
		CommonFields,
		field.String("name").
			Default("EDNET"),
		field.String("email").
			Default("kayprogrammer1@gmail.com"),
		field.String("phone").
			Default("+2348133831036"),
		field.String("address").
			Default("234, Lagos, Nigeria"),
		field.String("fb").
			Default("https://facebook.com"),
		field.String("tw").
			Default("https://twitter.com"),
		field.String("wh").
			Default("https://wa.me/2348133831036"),
		field.String("ig").
			Default("https://instagram.com"),
	)
}
