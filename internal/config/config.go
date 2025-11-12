package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Jwt      JwtConfig      `mapstructure:"jwt"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`

	// Connection Pool Settings
	MaxOpenConns    int `mapstructure:"max_open_conns"`
	MaxIdleConns    int `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int `mapstructure:"conn_max_lifetime"`  // in seconds
	ConnMaxIdleTime int `mapstructure:"conn_max_idle_time"` // in seconds
}

type JwtConfig struct {
	SecretKey                  string `mapstructure:"secret_key"` // must be at least 32 characters
	Issuer                     string `mapstructure:"issuer"`
	ExpirationHours            int    `mapstructure:"expiration_hours"`
	RefreshTokenExpirationDays int    `mapstructure:"refresh_token_expiration_days"`
}

type RedisConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Password        string `mapstructure:"password"`
	DB              int    `mapstructure:"db"`
	MessageTTLHours int    `mapstructure:"message_ttl_hours"`
}

// Validate validates the JWT configuration parameters.
func (j JwtConfig) Validate() error {
	if j.SecretKey == "" {
		return fmt.Errorf("jwt secret key cannot be empty")
	}

	if len(j.SecretKey) < 32 {
		return fmt.Errorf("jwt secret key must be at least 32 characters for security, got %d", len(j.SecretKey))
	}

	if j.ExpirationHours <= 0 {
		return fmt.Errorf("jwt expiration hours must be positive, got %d", j.ExpirationHours)
	}

	if j.ExpirationHours > 8760 { // 1 year
		return fmt.Errorf("jwt expiration hours too long (max 8760 hours/1 year), got %d", j.ExpirationHours)
	}

	if j.Issuer == "" {
		return fmt.Errorf("jwt issuer cannot be empty")
	}

	if j.RefreshTokenExpirationDays <= 0 {
		return fmt.Errorf("JWT refresh token expiration days must be positive, got %d", j.RefreshTokenExpirationDays)
	}

	return nil
}

func (r RedisConfig) Validate() error {
	if r.Host == "" {
		return fmt.Errorf("redis host cannot be empty")
	}

	if r.Port <= 0 {
		return fmt.Errorf("redis port must be positive, got %d", r.Port)
	}

	if r.DB < 0 {
		return fmt.Errorf("redis DB must be non-negative, got %d", r.DB)
	}

	return nil
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate JWT configuration
	if err := cfg.Jwt.Validate(); err != nil {
		return nil, fmt.Errorf("invalid JWT configuration: %w", err)
	}

	// Validate Redis configuration
	if err := cfg.Redis.Validate(); err != nil {
		return nil, fmt.Errorf("invalid Redis configuration: %w", err)
	}

	return &cfg, nil
}
