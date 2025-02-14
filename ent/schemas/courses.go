package schemas

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Category schema.
type Category struct {
	ent.Schema
}

// Fields of the Category.
func (Category) Fields() []ent.Field {
	return append(
		CommonFields,
		field.String("name").Unique(),
		field.String("slug").Unique(),
	)

}

// Edges of the Category.
func (Category) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("courses", Course.Type),
	}
}

// Tag schema.
type Tag struct {
	ent.Schema
}

// Fields of the Tag.
func (Tag) Fields() []ent.Field {
	return append(
		CommonFields,
		field.String("name").Unique(),
		field.String("slug").Unique(),
	)
}

// Edges of the Tag.
func (Tag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("courses", Course.Type),
	}
}

// Course schema.
type Course struct {
	ent.Schema
}

// Fields of the Course.
func (Course) Fields() []ent.Field {
	return append(
		CommonFields,
		field.String("title").NotEmpty(),
		field.String("slug").Unique(),
		field.Text("desc"),
		field.String("thumbnail_url"),
		field.String("intro_video_url").Optional(),
		field.UUID("category_id", uuid.UUID{}),
		field.String("language").Default("English"),
		field.Enum("difficulty").Values("Beginner", "Intermediate", "Advanced"),
		field.UUID("instructor_id", uuid.UUID{}),
		field.Int("students_count").Default(0),
		field.Int("total_lessons").Default(0),
		field.Int("total_quizzes").Default(0),
		field.Int("duration").Default(0), // in minutes
		field.Bool("is_published").Default(false),
		field.Bool("is_free").Default(false),
		field.Float("price").Default(0.0),
		field.Float("discount_price").Optional(),
		field.Enum("enrollment_type").Values("Open", "Restricted", "InviteOnly"),
		field.Bool("certification").Default(false),
		field.Float("rating").Default(0.0),
		field.Int("reviews_count").Default(0),
	)
}

// Edges of the Course.
func (Course) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("instructor", User.Type).Ref("courses").Field("instructor_id").Unique().Required(),
		edge.From("category", Category.Type).Ref("courses").Field("category_id").Unique().Required(),
		edge.From("tags", Tag.Type).Ref("courses"),
		edge.To("lessons", Lesson.Type),
		edge.To("enrollments", Enrollment.Type),
		edge.To("reviews", Review.Type),
		edge.To("payments", Payment.Type),
	}
}

// Lesson schema.
type Lesson struct {
	ent.Schema
}

// Fields of the Lesson.
func (Lesson) Fields() []ent.Field {
	return append(
		CommonFields,
		field.UUID("course_id", uuid.UUID{}),
		field.String("title").NotEmpty(),
		field.String("slug").Unique(),
		field.Text("description"),
		field.String("video_url").Optional(),
		field.Text("content").Optional(),
		field.Int("order"),
		field.Int("duration").Default(0),
		field.Bool("is_published").Default(false),
		field.Bool("is_free_preview").Default(false),
	)
}

// Edges of the Lesson.
func (Lesson) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("course", Course.Type).Ref("lessons").Field("course_id").Unique().Required(),
	}
}

// Enrollment schema.
type Enrollment struct {
	ent.Schema
}

// Fields of Enrollment.
func (Enrollment) Fields() []ent.Field {
	return append(
		CommonFields,
		field.UUID("user_id", uuid.UUID{}),
		field.UUID("course_id", uuid.UUID{}),
		field.Enum("status").Values("active", "completed", "dropped").Default("active"),
		field.Int("progress").Default(0), // Percentage (0-100)
	)
}

// Edges of Enrollment.
func (Enrollment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("enrollments").Field("user_id").Unique().Required(),
		edge.From("course", Course.Type).Ref("enrollments").Field("course_id").Unique().Required(),
	}
}

// Review schema.
type Review struct {
	ent.Schema
}

// Fields of Review.
func (Review) Fields() []ent.Field {
	return append(
		CommonFields,
		field.UUID("user_id", uuid.UUID{}),
		field.UUID("course_id", uuid.UUID{}),
		field.Float("rating").Min(1.0).Max(5.0),
		field.Text("comment").Optional(),
	)
}

// Edges of Review.
func (Review) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("reviews").Field("user_id").Unique().Required(),
		edge.From("course", Course.Type).Ref("reviews").Field("course_id").Unique().Required(),
	}
}

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
