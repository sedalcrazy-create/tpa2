package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig
	CORS     CORSConfig
	External ExternalConfig
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name        string
	Environment string // development, staging, production
	Port        int
	Debug       bool
	BaseURL     string
	APIPrefix   string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	Name         string
	SSLMode      string
	MaxIdleConns int
	MaxOpenConns int
	MaxLifetime  time.Duration
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret           string
	AccessExpiresIn  time.Duration
	RefreshExpiresIn time.Duration
	Issuer           string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	MaxAge         int
}

// ExternalConfig holds external API configurations
type ExternalConfig struct {
	Tamin  TaminConfig
	Sepas  SepasConfig
	IRC    IRCConfig
}

// TaminConfig - تنظیمات API تامین اجتماعی
type TaminConfig struct {
	BaseURL  string
	Username string
	Password string
	Timeout  time.Duration
}

// SepasConfig - تنظیمات API سپاس
type SepasConfig struct {
	BaseURL string
	APIKey  string
	Timeout time.Duration
}

// IRCConfig - تنظیمات IRC برای داروها
type IRCConfig struct {
	BaseURL string
	APIKey  string
	Timeout time.Duration
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "TPA System"),
			Environment: getEnv("APP_ENV", "development"),
			Port:        getEnvAsInt("APP_PORT", 8080),
			Debug:       getEnvAsBool("APP_DEBUG", true),
			BaseURL:     getEnv("APP_BASE_URL", "http://localhost:8080"),
			APIPrefix:   getEnv("API_PREFIX", "/api/v1"),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnvAsInt("DB_PORT", 5432),
			User:         getEnv("DB_USER", "postgres"),
			Password:     getEnv("DB_PASSWORD", ""),
			Name:         getEnv("DB_NAME", "tpa"),
			SSLMode:      getEnv("DB_SSLMODE", "disable"),
			MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
			MaxLifetime:  getEnvAsDuration("DB_MAX_LIFETIME", 1*time.Hour),
		},
		JWT: JWTConfig{
			Secret:           getEnv("JWT_SECRET", "change-this-secret-in-production"),
			AccessExpiresIn:  getEnvAsDuration("JWT_ACCESS_EXPIRES", 1*time.Hour),
			RefreshExpiresIn: getEnvAsDuration("JWT_REFRESH_EXPIRES", 7*24*time.Hour),
			Issuer:           getEnv("JWT_ISSUER", "tpa-system"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvAsSlice("CORS_ORIGINS", []string{"*"}),
			AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Tenant-ID"},
			MaxAge:         86400,
		},
		External: ExternalConfig{
			Tamin: TaminConfig{
				BaseURL:  getEnv("TAMIN_BASE_URL", ""),
				Username: getEnv("TAMIN_USERNAME", ""),
				Password: getEnv("TAMIN_PASSWORD", ""),
				Timeout:  getEnvAsDuration("TAMIN_TIMEOUT", 30*time.Second),
			},
			Sepas: SepasConfig{
				BaseURL: getEnv("SEPAS_BASE_URL", ""),
				APIKey:  getEnv("SEPAS_API_KEY", ""),
				Timeout: getEnvAsDuration("SEPAS_TIMEOUT", 30*time.Second),
			},
			IRC: IRCConfig{
				BaseURL: getEnv("IRC_BASE_URL", ""),
				APIKey:  getEnv("IRC_API_KEY", ""),
				Timeout: getEnvAsDuration("IRC_TIMEOUT", 30*time.Second),
			},
		},
	}

	return cfg, nil
}

// DSN returns the database connection string
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// RedisAddr returns the Redis address
func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		// Simple comma-separated parsing
		var result []string
		start := 0
		for i := 0; i <= len(value); i++ {
			if i == len(value) || value[i] == ',' {
				if start < i {
					result = append(result, value[start:i])
				}
				start = i + 1
			}
		}
		return result
	}
	return defaultValue
}
