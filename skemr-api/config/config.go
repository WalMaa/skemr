package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	defaultAppPort      = 8080
	defaultDatabasePort = 5432
	defaultRedisPort    = 6379
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
			slog.Error("Warning: Could not load .env file", "err", err)
		}
	}

	// Set defaults first
	viper.SetDefault("app.env", env)
	viper.SetDefault("app.port", defaultAppPort)
	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", defaultDatabasePort)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "pass")
	viper.SetDefault("database.name", "postgres")
	viper.SetDefault("database.sslmode", "disable")
	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", defaultRedisPort)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// Enable reading from environment variables
	viper.AutomaticEnv()

	// Bind environment variables to viper keys
	if err := viper.BindEnv("app.env", "APP_ENV"); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("app.port", "APP_PORT"); err != nil {
		return nil, err
	}
	// Database env vars
	if err := viper.BindEnv("database.host", "DB_HOST"); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("database.port", "DB_PORT"); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("database.user", "DB_USER"); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("database.password", "DB_PASSWORD"); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("database.name", "DB_NAME"); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("database.sslmode", "DB_SSLMODE"); err != nil {
		return nil, err
	}
	// Redis env vars
	if err := viper.BindEnv("redis.host", "REDIS_HOST"); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("redis.port", "REDIS_PORT"); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("redis.password", "REDIS_PASSWORD"); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("redis.db", "REDIS_DB"); err != nil {
		return nil, err
	}

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
