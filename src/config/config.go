package config

import (
	"fmt"
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
	ChunkSize      int64
	BackendUrl     string
}

var Env *Config

func init() {
	// Load .env file if present
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	Env = &Config{
		JWTSecret:      getEnvWithTransformer("JWT_SECRET", jwtSecretTransformer),
		JWTMaximumAge:  getEnvWithTransformerAndDefaultValue("JWT_MAXIMUM_AGE", maximumAgeTransformer, "1d"),
		MinioAccessKey: getEnv("MINIO_ACCESS_KEY"),
		MinioEndpoint:  getEnv("MINIO_ENDPOINT"),
		MinioSecretKey: getEnv("MINIO_SECRET_KEY"),
		APIKey:         getEnv("API_KEY"),
		ChunkSize:      getEnvWithTransformerAndDefaultValue("CHUNK_SIZE", chunkSizeTransformer, fmt.Sprintf("%d", 1024*1024*5)),
		BackendUrl:     getEnv("BACKEND_URL"),
	}
}

// Helper function to read an environment variable
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

// Helper function to read an environment variable and then transform it
func getEnvWithTransformer[T any](key string, transformer transformer[T]) T {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalln("Environment variable ", key, " is required but not set")
	}
	return transformer(key, value)
}

// Helper function to read an environment variable and then transform it or the default value
func getEnvWithTransformerAndDefaultValue[T any](key string, transformer transformer[T], defaultValue string) T {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return transformer(key, value)
}
