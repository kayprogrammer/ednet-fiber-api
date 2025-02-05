package config

import (
	"os"

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
	FirstSuperuserEmail       string `mapstructure:"FIRST_SUPERUSER_EMAIL"`
	FirstSuperUserPassword    string `mapstructure:"FIRST_SUPERUSER_PASSWORD"`
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
	RedisUrl                  string `mapstructure:"REDIS_URL"`
}

func GetConfig() (config Config) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "." // Default to current directory if not set
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	var err error
	if err = viper.ReadInConfig(); err != nil {
		panic(err)
	}
	viper.Unmarshal(&config)
	config.SecretKeyByte = []byte(config.SecretKey)
	return
}
