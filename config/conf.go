package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Environment               string `mapstructure:"ENVIRONMENT"`
	EmailOtpExpireMinutes     int64  `mapstructure:"EMAIL_OTP_EXPIRE_MINUTES"`
	AccessTokenExpireMinutes  int    `mapstructure:"ACCESS_TOKEN_EXPIRE_MINUTES"`
	RefreshTokenExpireMinutes int    `mapstructure:"REFRESH_TOKEN_EXPIRE_MINUTES"`
	Port                      string `mapstructure:"PORT"`
	SecretKey                 string `mapstructure:"SECRET_KEY"`
	SecretKeyByte             []byte
	FirstAdminEmail           string `mapstructure:"FIRST_ADMIN_EMAIL"`
	FirstAdminPassword        string `mapstructure:"FIRST_ADMIN_PASSWORD"`
	FirstInstructorEmail      string `mapstructure:"FIRST_INSTRUCTOR_EMAIL"`
	FirstInstructorPassword   string `mapstructure:"FIRST_INSTRUCTOR_PASSWORD"`
	FirstStudentEmail         string `mapstructure:"FIRST_STUDENT_EMAIL"`
	FirstStudentPassword      string `mapstructure:"FIRST_STUDENT_PASSWORD"`
	PostgresUser              string `mapstructure:"POSTGRES_USER"`
	PostgresPassword          string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresServer            string `mapstructure:"POSTGRES_SERVER"`
	PostgresPort              string `mapstructure:"POSTGRES_PORT"`
	PostgresDB                string `mapstructure:"POSTGRES_DB"`
	TestPostgresDB            string `mapstructure:"TEST_POSTGRES_DB"`
	MailSenderEmail           string `mapstructure:"MAIL_SENDER_EMAIL"`
	MailFrom                  string `mapstructure:"MAIL_FROM"`
	MailSenderPassword        string `mapstructure:"MAIL_SENDER_PASSWORD"`
	MailSenderHost            string `mapstructure:"MAIL_SENDER_HOST"`
	MailSenderPort            int    `mapstructure:"MAIL_SENDER_PORT"`
	CORSAllowedOrigins        string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	CORSAllowCredentials      bool   `mapstructure:"CORS_ALLOW_CREDENTIALS"`
	GoogleClientID            string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret        string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	SocialsPassword           string `mapstructure:"SOCIALS_PASSWORD"`
	StripePublicKey           string `mapstructure:"STRIPE_PUBLIC_KEY"`
	StripeSecretKey           string `mapstructure:"STRIPE_SECRET_KEY"`
	StripeWebhookSecret       string `mapstructure:"STRIPE_WEBHOOK_SECRET"`
	SocketSecret              string `mapstructure:"SOCKET_SECRET"`
	CloudinaryCloudName       string `mapstructure:"CLOUDINARY_CLOUD_NAME"`
	CloudinaryApiKey          string `mapstructure:"CLOUDINARY_API_KEY"`
	CloudinaryApiSecret       string `mapstructure:"CLOUDINARY_API_SECRET"`
	GeminiApiKey              string `mapstructure:"GEMINI_API_KEY"`
}

// bindEnvs explicitly binds environment variables to viper keys using struct tags.
// It expects a pointer to a struct, which is a more idiomatic way to handle reflection.
func bindEnvs(cfgPtr interface{}) {
	v := reflect.ValueOf(cfgPtr)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		panic("bindEnvs requires a non-nil pointer to a struct")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		panic("bindEnvs requires a pointer to a struct")
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("mapstructure")
		if tag != "" {
			// Bind the environment variable to the viper key.
			_ = viper.BindEnv(tag)
		}
	}
}

func GetConfig() (config Config) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}
	if env == "local" || env == "development" {
		if err := godotenv.Load(); err != nil {
			fmt.Println("No .env file found (skipping godotenv)")
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	bindEnvs(&config)

	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("unable to decode into struct: %w", err))
	}

	fmt.Println("--- Loaded Configuration ---")
	fmt.Printf("Environment: %s\n", config.Environment)
	fmt.Printf("Port: %s\n", config.Port)
	fmt.Println("----------------------------")

	config.SecretKeyByte = []byte(config.SecretKey)
	return
}