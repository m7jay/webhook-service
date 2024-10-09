package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	// Server
	Port int

	// Database
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	// Redis
	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int

	// Kafka
	KafkaBrokers []string
	KafkaTopic   string
	KafkaGroupID string

	// JWT
	JWTSecret string

	// Webhook
	WebhookWorkers int
	RetryAttempts  int
	RetryDelay     int

	// File Processing
	UploadDirectory    string
	ProcessedDirectory string

	// Logging
	LogLevel string
	Logger   *logrus.Logger
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: Error loading .env file")
	}

	config := &Config{
		Port:               getEnvAsInt("PORT", 8080),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnvAsInt("DB_PORT", 5432),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", ""),
		DBName:             getEnv("DB_NAME", "webhook_service"),
		RedisHost:          getEnv("REDIS_HOST", "localhost"),
		RedisPort:          getEnvAsInt("REDIS_PORT", 6379),
		RedisPassword:      getEnv("REDIS_PASSWORD", ""),
		RedisDB:            getEnvAsInt("REDIS_DB", 0),
		KafkaBrokers:       getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
		KafkaTopic:         getEnv("KAFKA_TOPIC", "webhook_events"),
		KafkaGroupID:       getEnv("KAFKA_GROUP_ID", "webhook_group"),
		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key"),
		WebhookWorkers:     getEnvAsInt("WEBHOOK_WORKERS", 10),
		RetryAttempts:      getEnvAsInt("RETRY_ATTEMPTS", 3),
		RetryDelay:         getEnvAsInt("RETRY_DELAY", 5),
		UploadDirectory:    getEnv("UPLOAD_DIRECTORY", "./uploads"),
		ProcessedDirectory: getEnv("PROCESSED_DIRECTORY", "./processed"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
	}

	config.Logger = configureLogger(config.LogLevel)

	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, ",")
}

func configureLogger(logLevel string) *logrus.Logger {
	logger := logrus.New()

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	return logger
}
