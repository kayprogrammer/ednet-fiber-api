package schemas

import (
	"time"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
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
		edge.From("courses", Course.Type).Ref("tags"),
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
		field.String("slug").Unique().NotEmpty(),
		field.Text("desc").NotEmpty(),
		field.String("thumbnail_url").NotEmpty(),
		field.String("intro_video_url").Optional(),
		field.UUID("category_id", uuid.UUID{}),
		field.String("language").Default("English"),
		field.Enum("difficulty").Values("beginner", "intermediate", "advanced").Default("beginner"),
		field.UUID("instructor_id", uuid.UUID{}),
		field.Uint("duration").Default(0), // in minutes
		field.Bool("is_published").Default(false),
		field.Bool("is_free").Default(false),
		field.Float("price").Default(0.0),
		field.Float("discount_price").Default(0.0),
		field.Enum("enrollment_type").Values("open", "restricted", "invite_only").Default("open"),
		field.Bool("certification").Default(true),
	)
}

// Edges of the Course.
func (Course) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("instructor", User.Type).Ref("courses").Field("instructor_id").Unique().Required(),
		edge.From("category", Category.Type).Ref("courses").Field("category_id").Unique().Required(),
		edge.To("tags", Tag.Type),
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
		field.Text("desc"),
		field.String("thumbnail_url").NotEmpty(),
		field.String("video_url").Optional(),
		field.Text("content").Optional(),
		field.Uint("order"),
		field.Uint("duration").Default(1),
		field.Bool("is_published").Default(false),
		field.Bool("is_free_preview").Default(false),
	)
}

// Edges of the Lesson.
func (Lesson) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("course", Course.Type).Ref("lessons").Field("course_id").Unique().Required(),
		edge.To("quizzes", Quiz.Type),
		edge.To("progress", LessonProgress.Type),
	}
}

// LessonProgress schema.
type LessonProgress struct {
	ent.Schema
}

// Fields of the LessonProgress.
func (LessonProgress) Fields() []ent.Field {
	return append(
		CommonFields,
		field.UUID("user_id", uuid.UUID{}),
		field.UUID("lesson_id", uuid.UUID{}),
		field.Time("completed_at").Optional(),
	)
}

// Edges of the LessonProgress.
func (LessonProgress) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("progress").Field("user_id").Unique().Required(),
		edge.From("lesson", Lesson.Type).Ref("progress").Field("lesson_id").Unique().Required(),
	}
}

func (LessonProgress) Indexes() []ent.Index {
	return []ent.Index{
		// Unique constraint on user_id + lesson_id to prevent duplicate progress entries
		index.Fields("user_id", "lesson_id").Unique(),
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
		field.Enum("status").Values("inactive", "active", "completed", "dropped").Default("inactive"),
		field.Enum("payment_status").Values("successful", "cancelled", "pending", "failed").Default("pending"),
		field.String("checkout_url").Optional(),
		field.Int("progress").Default(0), // Percentage (0-100)
		field.String("cert").Optional(),
	)
}

// Edges of Enrollment.
func (Enrollment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("enrollments").Field("user_id").Unique().Required(),
		edge.From("course", Course.Type).Ref("enrollments").Field("course_id").Unique().Required(),
	}
}

func (Enrollment) Indexes() []ent.Index {
	return []ent.Index{
		// Unique constraint on user_id + course_id to prevent duplicate enrollments
		index.Fields("user_id", "course_id").Unique(),
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

func (Review) Indexes() []ent.Index {
	return []ent.Index{
		// Unique constraint on user_id + course_id to prevent duplicate reviews
		index.Fields("user_id", "course_id").Unique(),
	}
}

// Quiz schema.
type Quiz struct {
	ent.Schema
}

// Fields of the Quiz.
func (Quiz) Fields() []ent.Field {
	return append(
		CommonFields,
		field.UUID("lesson_id", uuid.UUID{}),
		field.String("title").NotEmpty(),
		field.String("slug").Unique(),
		field.Text("description").Optional(),
		field.Int("duration").Default(0), // in minutes
		field.Bool("is_published").Default(false),
	)
}

// Edges of the Quiz.
func (Quiz) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("lesson", Lesson.Type).Ref("quizzes").Field("lesson_id").Unique().Required(),
		edge.To("questions", Question.Type),
		edge.To("results", QuizResult.Type),
	}
}

// Question schema.
type Question struct {
	ent.Schema
}

// Fields of the Question.
func (Question) Fields() []ent.Field {
	return append(
		CommonFields,
		field.UUID("quiz_id", uuid.UUID{}),
		field.Text("text").NotEmpty(), // Question text
		field.Int("order"),            // Order in the quiz
	)
}

// Edges of the Question.
func (Question) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("quiz", Quiz.Type).Ref("questions").Field("quiz_id").Unique().Required(),
		edge.To("options", QuestionOption.Type),
		edge.To("answers", Answer.Type),
	}
}

// Option schema.
type QuestionOption struct {
	ent.Schema
}

// Fields of the Option.
func (QuestionOption) Fields() []ent.Field {
	return append(
		CommonFields,
		field.UUID("question_id", uuid.UUID{}),
		field.Text("text").NotEmpty(),           // Option text
		field.Bool("is_correct").Default(false), // Correct answer flag
	)
}

// Edges of the Option.
func (QuestionOption) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("question", Question.Type).Ref("options").Field("question_id").Unique().Required(),
		edge.To("answers", Answer.Type),
	}
}

// QuizResult schema.
type QuizResult struct {
	ent.Schema
}

// Fields of the QuizResult.
func (QuizResult) Fields() []ent.Field {
	return append(
		CommonFields,
		field.UUID("user_id", uuid.UUID{}),
		field.UUID("quiz_id", uuid.UUID{}),

		field.Float("score").Min(0).Max(100).Default(0),
		field.Int("time_taken").Default(0), // in seconds

		field.Time("started_at").Optional().Default(time.Now),
		field.Time("completed_at").Optional().Nillable(),
	)
}

// Edges of the QuizResult.
func (QuizResult) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("quiz_results").Field("user_id").Unique().Required(),
		edge.From("quiz", Quiz.Type).Ref("results").Field("quiz_id").Unique().Required(),
		edge.To("answers", Answer.Type),
	}
}



// Answer schema.
type Answer struct {
	ent.Schema
}

// Fields of the Answer.
func (Answer) Fields() []ent.Field {
	return append(
		CommonFields,
		field.UUID("result_id", uuid.UUID{}),
		field.UUID("question_id", uuid.UUID{}),
		field.UUID("selected_option_id", uuid.UUID{}),
		field.Bool("is_correct").Default(false),
	)
}

// Edges of the Answer.
func (Answer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("result", QuizResult.Type).Ref("answers").Field("result_id").Unique().Required(),
		edge.From("question", Question.Type).Ref("answers").Field("question_id").Unique().Required(),
		edge.From("selected_option", QuestionOption.Type).Ref("answers").Field("selected_option_id").Unique().Required(),
	}
}
