package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBURL      string

	JWTSecret       string
	JWTAccessHours  int
	JWTRememberDays int
	JWTRefreshDays  int

	SMTPHost                     string
	SMTPPort                     int
	SMTPUser                     string
	SMTPPassword                 string
	AppURL                       string
	VerificationTokenExpiryHours int

	StripeSecretKey     string
	StripeWebhookSecret string
	StripeCurrency      string

	AdminEmail string
	Port       string
}

var App *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	App = &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "marketplace"),
		DBURL:      getEnv("DATABASE_URL", ""),

		JWTSecret:       getEnv("JWT_SECRET", "secret"),
		JWTAccessHours:  getEnvInt("JWT_ACCESS_HOURS", 24),
		JWTRememberDays: getEnvInt("JWT_REMEMBER_DAYS", 30),
		JWTRefreshDays:  getEnvInt("JWT_REFRESH_DAYS", 60),

		SMTPHost:                     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:                     getEnvInt("SMTP_PORT", 587),
		SMTPUser:                     getEnv("SMTP_USER", ""),
		SMTPPassword:                 getEnv("SMTP_PASSWORD", ""),
		AppURL:                       getEnv("APP_URL", "http://localhost:5000"),
		VerificationTokenExpiryHours: getEnvInt("VERIFICATION_TOKEN_EXPIRY_HOURS", 72),

		StripeSecretKey:     getEnv("STRIPE_SECRET_KEY", ""),
		StripeWebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),
		StripeCurrency:      getEnv("STRIPE_CURRENCY", "usd"),

		AdminEmail: getEnv("ADMIN_EMAIL", "ulvinmzv43@gmail.com"),
		Port:       getEnv("PORT", "5000"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}
