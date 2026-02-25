package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Env  string
		Port int
	}
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
		SSLMode  string `mapstructure:"sslmode"`
	}
	Redis struct {
		Host     string
		Port     int
		Password string
		DB       int
	}
}

func LoadConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	// Load .env file if it exists (for dev mode)
	if env == "dev" {
		if err := godotenv.Load(".env"); err != nil {
			// Log the error but continue, as environment variables might still be set
			slog.Error("Warning: Could not load .env file: %v", err)
		}
	}

	// Set defaults first
	viper.SetDefault("app.env", env)
	viper.SetDefault("app.port", 8080)
	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "pass")
	viper.SetDefault("database.name", "postgres")
	viper.SetDefault("database.sslmode", "disable")
	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// Enable reading from environment variables
	viper.AutomaticEnv()

	// Bind environment variables to viper keys
	viper.BindEnv("app.env", "APP_ENV")
	viper.BindEnv("app.port", "APP_PORT")
	// Database env vars
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.name", "DB_NAME")
	viper.BindEnv("database.sslmode", "DB_SSLMODE")
	// Redis env vars
	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("redis.db", "REDIS_DB")

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
