package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret      []byte
	JWTMaximumAge  time.Duration
	MinioAccessKey string
	MinioEndpoint  string
	MinioSecretKey string
	APIKey         string
}

var Env *Config

func init() {
	// Load .env file if present
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	Env = &Config{
		JWTSecret:      getEnvWithTransformer("JWT_SECRET", jwtSecretTransformer),
		JWTMaximumAge:  getEnvWithTransformerAndDefaultValue("JWT_MAXIMUM_AGE", maximumAgeTransformer, "1d"),
		MinioAccessKey: getEnv("MINIO_ACCESS_KEY"),
		MinioEndpoint:  getEnv("MINIO_ENDPOINT"),
		MinioSecretKey: getEnv("MINIO_SECRET_KEY"),
		APIKey:         getEnv("API_KEY"),
	}
}

// Helper function to read an environment variable or return a default value
func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalln("Environment variable ", key, " is required but not set")
	}
	return value
}

// Helper function to read an environment variable or return a default value
func getEnvWithDefaultValue(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return value
}

// Helper function to read an environment variable or return a default value
func getEnvWithTransformer[T any](key string, transformer transformer[T]) T {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalln("Environment variable ", key, " is required but not set")
	}
	return transformer(key, value)
}

// Helper function to read an environment variable or return a default value
func getEnvWithTransformerAndDefaultValue[T any](key string, transformer transformer[T], defaultValue string) T {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return transformer(key, value)
}
