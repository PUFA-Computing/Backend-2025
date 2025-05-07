package configs

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	RedisURL  string
	RedisPass string

	ApiPort      string
	JWTSecretKey string

	CloudflareAccountId   string
	CloudflareR2AccessId  string
	CloudflareR2AccessKey string

	S3Endpoint         string
	AWSAccessKeyId     string
	AWSSecretAccessKey string
	AWSRegion          string
	S3UsePathStyle     bool
	S3Bucket           string

	// Email service toggle
	UseSmtp bool

	// Legacy SMTP settings
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SenderEmail  string

	// SendGrid settings
	SendGridAPIKey     string
	SendGridSender     string
	SendGridSenderName string

	GithubAccessToken string
	HunterApiKey      string

	BaseURL string
}

func LoadConfig() *Config {
	env := os.Getenv("ENV")

	var baseURl string
<<<<<<< HEAD
	if env == "production" {
		baseURl = "https://compsci.president.ac.id"
	} else if env == "staging" {
=======
	if env == "local" || env == "test" {
>>>>>>> 6d28e549bd0c6115e365a5402caec0ff3f844d69
		baseURl = "http://localhost:3000"
	} else {
		// Default to production URL if ENV is not explicitly set to local/test
		baseURl = "https://compsci.president.ac.id"
	}

	cfg := &Config{
		DBHost:                os.Getenv("DB_HOST"),
		DBPort:                os.Getenv("DB_PORT"),
		DBUser:                os.Getenv("DB_USER"),
		DBPassword:            os.Getenv("DB_PASSWORD"),
		DBName:                os.Getenv("DB_NAME"),
		RedisURL:              os.Getenv("REDIS_URL"),
		RedisPass:             os.Getenv("REDIS_PASS"),
		ApiPort:               os.Getenv("API_PORT"),
		JWTSecretKey:          os.Getenv("JWT_SECRET_KEY"),
		CloudflareAccountId:   os.Getenv("CLOUDFLARE_ACCOUNT_ID"),
		CloudflareR2AccessId:  os.Getenv("CLOUDFLARE_R2_ACCESS_ID"),
		CloudflareR2AccessKey: os.Getenv("CLOUDFLARE_R2_ACCESS_KEY"),
		S3Endpoint:            os.Getenv("S3_ENDPOINT"),
		S3UsePathStyle:        os.Getenv("S3_USE_PATH_STYLE") == "true",
		AWSAccessKeyId:        os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey:    os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWSRegion:             os.Getenv("AWS_REGION"),
		S3Bucket:              os.Getenv("S3_BUCKET"),
		// Email service toggle
		UseSmtp: 			   os.Getenv("USE_SMTP") == "true",
		// Legacy SMTP settings
		SMTPHost:      		   os.Getenv("SMTP_HOST"),
		SMTPPort:      		   os.Getenv("SMTP_PORT"),
		SMTPUsername:  		   os.Getenv("SMTP_USERNAME"),
		SMTPPassword:  		   os.Getenv("SMTP_PASSWORD"),
		SenderEmail:   		   os.Getenv("SMTP_SENDER_EMAIL"),
		// SendGrid settings
		SendGridAPIKey:        os.Getenv("SENDGRID_API_KEY"),
		SendGridSender:        os.Getenv("SENDGRID_SENDER_EMAIL"),
		SendGridSenderName:    os.Getenv("SENDGRID_SENDER_NAME"),
		BaseURL:               baseURl,
		GithubAccessToken:     os.Getenv("GH_ACCESS_TOKEN"),
		HunterApiKey:          os.Getenv("HUNTER_API_KEY"),
	}

	fmt.Printf("Loaded Config: %+v\n", cfg)
	return cfg
}
