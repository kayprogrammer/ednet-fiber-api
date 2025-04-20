package seeding

type CourseData struct {
	Title     string
	Slug      string
	Desc      string
	ThumbnailUrl string
	IntroVideoUrl string
	Duration  uint
	IsFree    bool
	Price float64
	DiscountPrice float64
	Lessons   []LessonData
	Quiz      *QuizData
}

type LessonData struct {
	Title     string
	Slug      string
	Desc      string
	ThumbnailUrl string
	VideoUrl string
	Content string
	Order     uint
	Duration  uint
	IsPublished bool
	IsFreePreview bool
}

type QuizData struct {
	Title       string
	Slug      string
	Description string
	Duration    int
	Questions   []QuestionData
}

type QuestionData struct {
	Text    string
	Order   int
	Options []OptionData
}

type OptionData struct {
	Text      string
	IsCorrect bool
}

var CategoriesToCreate = []string{"Programming", "API", "Software Development"}

var CoursesToCreate = []CourseData{
	// Course 1: Go for Beginners
	{
		Title:         "Go for Beginners",
		Slug:          "go-for-beginners",
		Desc:          "An introductory course to the Go programming language.",
		ThumbnailUrl:  "https://placehold.co/600x400",
		IntroVideoUrl: "https://videos.example.com/go-beginners-intro.mp4",
		Duration:      120,
		IsFree:        true,
		Price:         0,
		DiscountPrice: 0,
		Lessons: []LessonData{
			{"Getting Started", "lesson-1-go", "Learn how to set up Go.", "https://placehold.co/300x200", "https://videos.example.com/lesson1.mp4", "Installing Go and writing your first program.", 1, 10, true, true},
			{"Variables", "lesson-2-go", "Understanding variables in Go.", "https://placehold.co/300x200", "https://videos.example.com/lesson2.mp4", "Different types of variables and how to use them.", 2, 12, true, false},
			{"Functions", "lesson-3-go", "Using functions in Go.", "https://placehold.co/300x200", "https://videos.example.com/lesson3.mp4", "Define and invoke functions effectively.", 3, 15, true, false},
			{"Structs", "lesson-4-go", "Introduction to structs.", "https://placehold.co/300x200", "https://videos.example.com/lesson4.mp4", "Learn to model data using structs.", 4, 14, true, false},
			{"Concurrency", "lesson-5-go", "Go's concurrency model.", "https://placehold.co/300x200", "https://videos.example.com/lesson5.mp4", "Understanding goroutines and channels.", 5, 20, true, false},
		},
		Quiz: &QuizData{
			Title:       "Go Basics Quiz",
			Slug: "go-basics-quiz",
			Description: "Assess your understanding of Go fundamentals.",
			Duration:    10,
			Questions: []QuestionData{
				{
					Text: "Which keyword is used to declare a variable in Go?",
					Order: 1,
					Options: []OptionData{
						{"var", true}, {"let", false}, {"define", false}, {"declare", false},
					},
				},
				{
					Text: "What does ':=' do in Go?",
					Order: 2,
					Options: []OptionData{
						{"Declare and assign", true}, {"Only declare", false}, {"Only assign", false}, {"Compare values", false},
					},
				},
			},
		},
	},

	// Course 2: Advanced Go
	{
		Title:         "Advanced Go",
		Slug:          "advanced-go",
		Desc:          "Explore Go's powerful advanced features.",
		ThumbnailUrl:  "https://placehold.co/600x400",
		IntroVideoUrl: "https://videos.example.com/adv-go-intro.mp4",
		Duration:      150,
		IsFree:        false,
		Price:         49.99,
		DiscountPrice: 29.99,
		Lessons: []LessonData{
			{"Interfaces", "lesson-1-adv-go", "Learn about interfaces.", "https://placehold.co/300x200", "https://videos.example.com/adv-1.mp4", "How interfaces work in Go.", 1, 18, true, false},
			{"Error Handling", "lesson-2-adv-go", "Advanced error management.", "https://placehold.co/300x200", "https://videos.example.com/adv-2.mp4", "Custom errors, panic, and recover.", 2, 16, true, false},
			{"Reflection", "lesson-3-adv-go", "Using reflection in Go.", "https://placehold.co/300x200", "https://videos.example.com/adv-3.mp4", "The reflect package and dynamic types.", 3, 17, true, false},
			{"Channels", "lesson-4-adv-go", "Deep dive into channels.", "https://placehold.co/300x200", "https://videos.example.com/adv-4.mp4", "Buffered vs unbuffered, select statement.", 4, 22, true, false},
			{"Generics", "lesson-5-adv-go", "Go's generics system.", "https://placehold.co/300x200", "https://videos.example.com/adv-5.mp4", "How to use generics in Go 1.18+.", 5, 25, true, false},
		},
		Quiz: &QuizData{
			Title:       "Advanced Go Quiz",
			Slug: "advanced-go-quiz",
			Description: "Test your advanced knowledge of Go.",
			Duration:    15,
			Questions: []QuestionData{
				{
					Text: "Which package enables reflection in Go?",
					Order: 1,
					Options: []OptionData{
						{"reflect", true}, {"mirror", false}, {"runtime", false}, {"meta", false},
					},
				},
				{
					Text: "What is the main use of interfaces in Go?",
					Order: 2,
					Options: []OptionData{
						{"Abstraction and polymorphism", true}, {"UI rendering", false}, {"Memory allocation", false}, {"File IO", false},
					},
				},
			},
		},
	},

	// Course 3: Web Development with Go
	{
		Title:         "Web Development with Go",
		Slug:          "web-dev-go",
		Desc:          "Build scalable web apps using Go.",
		ThumbnailUrl:  "https://placehold.co/600x400",
		IntroVideoUrl: "https://videos.example.com/web-dev-intro.mp4",
		Duration:      180,
		IsFree:        false,
		Price:         59.99,
		DiscountPrice: 39.99,
		Lessons: []LessonData{
			{"HTTP Basics", "lesson-1-web-go", "Handling HTTP requests.", "https://placehold.co/300x200", "https://videos.example.com/web-1.mp4", "Using net/http package.", 1, 15, true, false},
			{"Routing", "lesson-2-web-go", "Setting up routes.", "https://placehold.co/300x200", "https://videos.example.com/web-2.mp4", "Basic and nested routing.", 2, 18, true, false},
			{"Templates", "lesson-3-web-go", "HTML templates in Go.", "https://placehold.co/300x200", "https://videos.example.com/web-3.mp4", "Dynamic rendering with html/template.", 3, 20, true, false},
			{"Middleware", "lesson-4-web-go", "Creating middleware.", "https://placehold.co/300x200", "https://videos.example.com/web-4.mp4", "Logging, auth, and error handling.", 4, 20, true, false},
			{"Sessions", "lesson-5-web-go", "User sessions in Go.", "https://placehold.co/300x200", "https://videos.example.com/web-5.mp4", "Session management with cookies.", 5, 22, true, false},
		},
		Quiz: &QuizData{
			Title:       "Go Web Quiz",
			Slug: "go-web-quiz",
			Description: "Evaluate your web dev skills in Go.",
			Duration:    12,
			Questions: []QuestionData{
				{
					Text: "Which package is used to serve HTTP in Go?",
					Order: 1,
					Options: []OptionData{
						{"net/http", true}, {"web/http", false}, {"http/server", false}, {"go/web", false},
					},
				},
				{
					Text: "What is the default HTTP method when a form is submitted?",
					Order: 2,
					Options: []OptionData{
						{"GET", true}, {"POST", false}, {"PUT", false}, {"PATCH", false},
					},
				},
			},
		},
	},

	// Course 4: Building REST APIs with Go
	{
		Title:         "Building REST APIs with Go",
		Slug:          "rest-api-go",
		Desc:          "Design and develop RESTful APIs using Go and Gin.",
		ThumbnailUrl:  "https://placehold.co/600x400",
		IntroVideoUrl: "https://videos.example.com/rest-api-intro.mp4",
		Duration:      130,
		IsFree:        true,
		Price:         0,
		DiscountPrice: 0,
		Lessons: []LessonData{
			{"Intro to REST", "lesson-1-api-go", "What are REST APIs?", "https://placehold.co/300x200", "https://videos.example.com/api-1.mp4", "Overview and principles of REST.", 1, 10, true, true},
			{"Using Gin", "lesson-2-api-go", "Setup and usage of Gin framework.", "https://placehold.co/300x200", "https://videos.example.com/api-2.mp4", "Building APIs with Gin.", 2, 15, true, false},
			{"CRUD Operations", "lesson-3-api-go", "Create, Read, Update, Delete.", "https://placehold.co/300x200", "https://videos.example.com/api-3.mp4", "CRUD in practice with handlers.", 3, 20, true, false},
			{"Middleware", "lesson-4-api-go", "Custom middleware.", "https://placehold.co/300x200", "https://videos.example.com/api-4.mp4", "Authentication and logging.", 4, 18, true, false},
			{"Swagger Docs", "lesson-5-api-go", "Auto-generated documentation.", "https://placehold.co/300x200", "https://videos.example.com/api-5.mp4", "Using swaggo to generate docs.", 5, 20, true, false},
		},
		Quiz: &QuizData{
			Title:       "REST API Quiz",
			Slug: "rest-api-quiz",
			Description: "Review your REST API knowledge.",
			Duration:    10,
			Questions: []QuestionData{
				{
					Text: "Which HTTP method is used to delete a resource?",
					Order: 1,
					Options: []OptionData{
						{"DELETE", true}, {"REMOVE", false}, {"DESTROY", false}, {"ERASE", false},
					},
				},
				{
					Text: "Which status code represents a successful POST request?",
					Order: 2,
					Options: []OptionData{
						{"201 Created", true}, {"200 OK", false}, {"204 No Content", false}, {"400 Bad Request", false},
					},
				},
			},
		},
	},

	// Course 5: Testing in Go
	{
		Title:         "Testing in Go",
		Slug:          "testing-in-go",
		Desc:          "Learn how to write tests and use Go's testing tools.",
		ThumbnailUrl:  "https://placehold.co/600x400",
		IntroVideoUrl: "https://videos.example.com/testing-intro.mp4",
		Duration:      90,
		IsFree:        false,
		Price:         24.99,
		DiscountPrice: 19.99,
		Lessons: []LessonData{
			{"Why Testing?", "lesson-1-testing-go", "Importance of testing.", "https://placehold.co/300x200", "https://videos.example.com/test-1.mp4", "Benefits of automated tests.", 1, 10, true, true},
			{"Unit Tests", "lesson-2-testing-go", "Write unit tests in Go.", "https://placehold.co/300x200", "https://videos.example.com/test-2.mp4", "Using the `testing` package.", 2, 15, true, false},
			{"Table-Driven Tests", "lesson-3-testing-go", "Efficient test cases.", "https://placehold.co/300x200", "https://videos.example.com/test-3.mp4", "Using data-driven test design.", 3, 14, true, false},
			{"Mocks and Interfaces", "lesson-4-testing-go", "Dependency injection for testing.", "https://placehold.co/300x200", "https://videos.example.com/test-4.mp4", "Mocking services with interfaces.", 4, 20, true, false},
			{"Integration Tests", "lesson-5-testing-go", "End-to-end testing strategies.", "https://placehold.co/300x200", "https://videos.example.com/test-5.mp4", "Full system testing with databases.", 5, 18, true, false},
		},
		Quiz: &QuizData{
			Title:       "Go Testing Quiz",
			Slug: "go-testing-quiz",
			Description: "Review your understanding of testing in Go.",
			Duration:    8,
			Questions: []QuestionData{
				{
					Text: "Which package is used for testing in Go?",
					Order: 1,
					Options: []OptionData{
						{"testing", true}, {"unittest", false}, {"gotest", false}, {"gounit", false},
					},
				},
				{
					Text: "What is a table-driven test?",
					Order: 2,
					Options: []OptionData{
						{"Test using input-output tables", true}, {"GUI-based testing", false}, {"Database testing", false}, {"Code formatting", false},
					},
				},
			},
		},
	},
}
